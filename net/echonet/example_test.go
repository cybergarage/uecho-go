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

func ExampleNewMessage() {
	msg := echonet.NewMessage(
		echonet.WithMessageESV(echonet.ESVReadRequest),
		echonet.WithMessageDEOJ(0x0EF001),
		echonet.WithMessageProperties(
			echonet.NewProperty(
				echonet.WithPropertyCode(0x8A),
			),
		),
	)
	_ = msg
}

func ExampleNewLocalNode() {
	dev, err := echonet.NewDevice(
		echonet.WithDeviceCode(0x029101), // Mono functional lighting
	)
	if err != nil {
		return
	}

	node := echonet.NewLocalNode(
		echonet.WithLocalNodeDevices(dev),
	)

	err = node.Start()
	if err != nil {
		return
	}
	defer node.Stop()
}

func ExampleNewController() {
	ctrl := echonet.NewController()
	err := ctrl.Start()
	if err != nil {
		return
	}
	defer ctrl.Stop()
}

func ExampleSharedStandardDatabase() {
	db := echonet.SharedStandardDatabase()

	man, ok := db.LookupManufacture(0x000005)
	if ok {
		fmt.Printf("Manufacture: %s\n", man.Name())
	}

	obj, ok := db.LookupObjectByCode(0x029101) // Mono functional lighting
	if ok {
		fmt.Printf("Object: %s\n", obj.Name())
	}
}

// Example demonstrates how to use the echonet package to discover ECHONET Lite nodes on the local network, retrieve their manufacturer information, enumerate their objects, and read the required properties of each object.
// The function initializes a controller, starts network discovery, and prints out details of each found node, including address, port, manufacturer, object codes, object names, and property values.
func Example() {
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

	// Searches echonet nodes in the local network with context  until the context is done.
	err = ctrl.Search(context.Background())
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	// Outputs all found nodes

	db := echonet.SharedStandardDatabase()

	for i, node := range ctrl.Nodes() {
		// Gets manufacture code.

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

		// Prints node data.

		fmt.Printf("[%d] %-15s:%d (%s)\n", i, node.Address(), node.Port(), manufactureName)

		for j, obj := range node.Objects() {
			// Prints object data.

			objName := obj.ClassName()
			if len(objName) == 0 {
				objName = unknown
			}
			fmt.Printf("    [%d] %06X (%s)\n", j, obj.Code(), objName)

			// Prints only read required properties with the current property data.

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
