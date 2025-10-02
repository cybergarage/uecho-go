// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet_test

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

const (
	unknown = "unknown"
)

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

	err = ctrl.SearchAllObjects()
	if err != nil {
		log.Errorf("%s", err)
		return
	}

	// Waits node responses in the local network

	time.Sleep(time.Second * 1)

	// Outputs all found nodes

	db := echonet.GetStandardDatabase()

	for i, node := range ctrl.Nodes() {
		// Gets manufacture code.

		manufactureName := unknown
		req := echonet.NewMessage()
		req.SetESV(echonet.ESVReadRequest)
		req.SetDEOJ(0x0EF001)
		req.AddProperty(echonet.NewProperty().SetCode(0x8A))
		res, err := ctrl.PostMessage(node, req)
		if err == nil {
			if props := res.Properties(); len(props) == 1 {
				manufacture, ok := db.FindManufacture(echonet.ManufactureCode(encoding.ByteToInteger(props[0].Data())))
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
				req := echonet.NewMessage()
				req.SetESV(echonet.ESVReadRequest)
				req.SetDEOJ(obj.Code())
				req.AddProperty(echonet.NewProperty().SetCode(prop.Code()))
				res, err := ctrl.PostMessage(node, req)
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
