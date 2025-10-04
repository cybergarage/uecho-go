// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
uechobench is a search utility for Echonet Lite.

	NAME
	uechobench

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
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

const (
	unknown = "unknown"
)

func main() {
	numRepeat := flag.Int("n", 1, "Number of repeat")
	verbose := flag.Bool("v", false, "Enable verbose output")
	flag.Parse()

	// Setup logger

	if *verbose {
		log.SetSharedLogger(log.NewStdoutLogger(log.LevelTrace))
	}

	// Start a controller for Echonet Lite node

	ctrl := NewSearchController()

	if *verbose {
		ctrl.SetListener(ctrl)
	}

	err := ctrl.Start()
	if err != nil {
		return
	}

	err = ctrl.Search(context.Background())
	if err != nil {
		return
	}

	// Wait node responses in the local network

	time.Sleep(time.Second * 1)

	// Output all found nodes

	db := echonet.SharedStandardDatabase()

	for n := 0; n < *numRepeat; n++ {
		for i, node := range ctrl.Nodes() {
			manufactureName := unknown
			req := echonet.NewMessage()
			req.SetESV(echonet.ESVReadRequest)
			req.SetDEOJ(0x0EF001)
			req.AddProperty(echonet.NewProperty().SetCode(0x8A))
			res, err := ctrl.PostMessage(context.Background(), node, req)
			if err == nil {
				if props := res.Properties(); len(props) == 1 {
					manufacture, ok := db.LookupManufacture(echonet.ManufactureCode(encoding.ByteToInteger(props[0].Data())))
					if ok {
						manufactureName = manufacture.Name()
					}
				}
			}

			fmt.Printf("[%d] %-15s:%d (%s)\n", i, node.Address(), node.Port(), manufactureName)

			for j, obj := range node.Objects() {
				objName := obj.ClassName()
				if len(objName) == 0 {
					objName = unknown
				}
				fmt.Printf("    [%d] %06X (%s)\n", j, obj.Code(), objName)

				for _, prop := range obj.Properties() {
					if !prop.IsReadable() {
						continue
					}
					propName := prop.Name()
					if len(propName) == 0 {
						propName = "(" + unknown + ")"
					}
					propData := ""
					req := echonet.NewMessage()
					req.SetESV(echonet.ESVReadRequest)
					req.SetDEOJ(obj.Code())
					req.AddProperty(echonet.NewProperty().SetCode(prop.Code()))
					res, err := ctrl.PostMessage(context.Background(), node, req)
					if err == nil {
						if props := res.Properties(); len(props) == 1 {
							propData = hex.EncodeToString(props[0].Data())
						}
					} else {
						propData = err.Error()
					}
					fmt.Printf("        [%02X] %s (%s)\n", prop.Code(), propData, propName)
				}
			}
		}
	}

	// Stop the controller

	err = ctrl.Stop()
	if err != nil {
		return
	}
}
