// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
uechosearch is a search utility for Echonet Lite.

	NAME
	uechosearch

	SYNOPSIS
	uechosearch [OPTIONS]

	DESCRIPTION
	uechosearch is a search utility for Echonet Lite.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

const (
	unknown = "unknown"
)

func main() {
	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	// Enables protocol logger

	if *verbose {
		log.SetSharedLogger(log.NewStdoutLogger(log.LevelTrace))
	}

	// Starts a controller for Echonet Lite node

	ctrl := NewSearchController()

	if *verbose {
		ctrl.SetListener(ctrl)
	}

	err := ctrl.Start()
	if err != nil {
		log.Fatalf("%s", err)
		return
	}

	err = ctrl.Search(context.Background())
	if err != nil {
		log.Fatalf("%s", err)
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
					echonet.WithPropertyCode(0x8A),
				),
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

	// Stops the controller

	err = ctrl.Stop()
	if err != nil {
		log.Fatalf("%s", err)
		return
	}
}
