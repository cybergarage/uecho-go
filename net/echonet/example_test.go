// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet_test

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	unknown = "unknown"
)

func ExampleNewProperty() {
	prop := echonet.NewProperty(
		echonet.WithPropertyCode(0x80),
		echonet.WithPropertyData([]byte{0x30}),
	)
	_ = prop
}

func ExampleNewProperty_noData() {
	prop := echonet.NewProperty(
		echonet.WithPropertyCode(0x80),
	)
	_ = prop
}

func ExampleNewDevice() {
	// Creates a standard mono functional device.
	dev, err := echonet.NewDevice(
		echonet.WithDeviceCode(0x029101),             // Mono functional lighting
		echonet.WithDeviceManufacturerCode(0xFFFFFF), // Experimental
	)
	if err != nil {
		return
	}
	// Sets the operation status to "on".
	dev.SetPropertyInteger(0x80, 0x30, 1)
}

func ExampleNewMessage_readRequest() {
	msg := echonet.NewMessage(
		echonet.WithMessageESV(echonet.ESVReadRequest),
		echonet.WithMessageDEOJ(0x029101), // Mono functional lighting
		echonet.WithMessageProperties(
			echonet.NewProperty(
				echonet.WithPropertyCode(0x8A), // Manufacturer code
			),
		),
	)
	_ = msg
}

func ExampleNewMessage_writeRequest() {
	msg := echonet.NewMessage(
		echonet.WithMessageESV(echonet.ESVWriteRequest),
		echonet.WithMessageDEOJ(0x029101), // Mono functional lighting
		echonet.WithMessageProperties(
			echonet.NewProperty(
				echonet.WithPropertyCode(0x80),         // Operation status
				echonet.WithPropertyData([]byte{0x30}), // ON
			),
		),
	)
	_ = msg
}

func ExampleNewMessage_search() {
	msg := echonet.NewMessage(
		echonet.WithMessageESV(echonet.ESVReadRequest),
		echonet.WithMessageDEOJ(0x0EF001), // Node profile
		echonet.WithMessageProperties(
			echonet.NewProperty(
				echonet.WithPropertyCode(0xD6), // Self-node instance list S
			),
		),
	)
	_ = msg
}

func ExampleNewDefaultConfig() {
	conf := echonet.NewDefaultConfig(
		echonet.WithConfigTCPEnabled(true),
	)
	_ = conf
}

func ExampleNewController() {
	ctrl := echonet.NewController()
	err := ctrl.Start()
	if err != nil {
		return
	}
	defer ctrl.Stop()
}

func ExampleSharedStandardDatabase_object() {
	db := echonet.SharedStandardDatabase()
	obj, ok := db.LookupObject(0x029101) // Mono functional lighting
	if ok {
		fmt.Printf("Object: %s\n", obj.Name())
	}
}

func ExampleStandardDatabase_manufacture() {
	db := echonet.SharedStandardDatabase()
	man, ok := db.LookupManufacture(0x000005)
	if ok {
		fmt.Printf("Manufacture: %s\n", man.Name())
	}
}

// The example demonstrates how to create and start a local Echonet node
// with a single device (Mono functional lighting) and a custom property request handler.
// The handler processes only write requests for the operation status property (0x80),
// accepting values 0x30 (ON) and 0x31 (OFF). Invalid requests are rejected with an error.
// The example shows device and node creation, handler registration, node startup, and cleanup.
func ExampleNewLocalNode() {
	// onRequest handles property requests for the local device.
	// It processes only write requests, updating the property if valid.
	onRequest := func(obj echonet.Object, esv protocol.ESV, reqProp protocol.Property) error {
		// Only handle write requests.
		if !esv.IsWriteRequest() {
			return nil
		}
		reqPropCode := reqProp.Code()
		reqPropData := reqProp.Data()
		// NOTE: Object::LookupProperty() will always find the property here
		// because onRequest is only called if the object has the specified property.
		targetProp, _ := obj.LookupProperty(reqPropCode)
		switch reqPropCode {
		case 0x80: // Operation status
			// Accept only 0x30 (ON) or 0x31 (OFF) as valid data.
			reqData, err := reqProp.AsByte()
			if err != nil {
				return err
			}
			switch reqData {
			case 0x30, 0x31:
				targetProp.SetData(reqPropData)
				return nil
			default:
			}
		}
		return fmt.Errorf("invalid request : %02X %s", reqPropCode, hex.EncodeToString(reqPropData))
	}

	dev, err := echonet.NewDevice(
		echonet.WithDeviceCode(0x029101), // Mono functional lighting
		echonet.WithDeviceRequestHandler(onRequest),
	)
	if err != nil {
		return
	}

	node := echonet.NewLocalNode(
		echonet.WithLocalNodeConfig(
			echonet.NewDefaultConfig(
				echonet.WithConfigTCPEnabled(false),
			),
		),
		echonet.WithLocalNodeDevices(dev),
	)

	err = node.Start()
	if err != nil {
		return
	}
	defer node.Stop()
}

// Example demonstrates how to use the echonet package to discover ECHONET Lite nodes
// on the local network, retrieve their manufacturer information, enumerate their objects,
// and read the required properties of each object.
// The function initializes a controller, starts network discovery, and prints out details
// of each found node, including address, port, manufacturer, object codes, object names,
// and property values.
func Example() {
	// Search for ECHONET Lite nodes on the local network.

	ctrl := echonet.NewController()

	err := ctrl.Start()
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	defer func() {
		err := ctrl.Stop()
		if err != nil {
			log.Errorf("%s", err)
		}
	}()

	err = ctrl.Search(context.Background())
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	// Output details of all discovered nodes

	db := echonet.SharedStandardDatabase()

	for i, node := range ctrl.Nodes() {
		// Get manufacturer code.

		manufactureName := unknown
		req := echonet.NewMessage(
			echonet.WithMessageESV(echonet.ESVReadRequest),
			echonet.WithMessageDEOJ(0x0EF001),
			echonet.WithMessageProperties(
				echonet.NewProperty(
					echonet.WithPropertyCode(0x8A)),
			),
		)
		res, err := ctrl.PostMessage(context.Background(), node, req)
		if err == nil {
			if props := res.Properties(); len(props) == 1 {
				manufacture, ok := db.LookupManufacture(echonet.ManufactureCode(encoding.ByteToInteger(props[0].Data())))
				if ok {
					manufactureName = manufacture.Name()
				}
			}
		}

		// Print manufacturer information for the node.

		fmt.Printf("[%d] %-15s:%d (%s)\n", i, node.Address(), node.Port(), manufactureName)

		for j, obj := range node.Objects() {
			// Print object information.

			objName := obj.ClassName()
			if len(objName) == 0 {
				objName = unknown
			}
			fmt.Printf("    [%d] %06X (%s)\n", j, obj.Code(), objName)

			// Iterate over all properties that are required to be readable and print their current data.

			for _, prop := range obj.Properties() {
				if !prop.IsReadRequired() {
					continue
				}
				propName := prop.Name()
				if len(propName) == 0 {
					propName = "(" + unknown + ")"
				}
				propData := "--"
				req := echonet.NewMessage(
					echonet.WithMessageESV(echonet.ESVReadRequest),
					echonet.WithMessageDEOJ(obj.Code()),
					echonet.WithMessageProperties(
						echonet.NewProperty(
							echonet.WithPropertyCode(prop.Code()),
						),
					),
				)
				res, err := ctrl.PostMessage(context.Background(), node, req)
				if err == nil {
					if props := res.Properties(); len(props) == 1 {
						propData = hex.EncodeToString(props[0].Data())
					}
				} else {
					propData = err.Error()
				}
				fmt.Printf("        [%02X] %s: %s\n", prop.Code(), propName, propData)
			}
		}
	}
}
