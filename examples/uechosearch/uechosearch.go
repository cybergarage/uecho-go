// Copyright (C) 2018 Satoshi Konno. All rights reserved.
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
	"encoding/hex"
	"flag"
	"fmt"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

func main() {
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

	err = ctrl.SearchAllObjects()
	if err != nil {
		return
	}

	// Wait node responses in the local network

	time.Sleep(time.Second * 1)

	// Output all found nodes

	db := echonet.GetStandardDatabase()

	for i, node := range ctrl.Nodes() {
		manufactureName := "unknown"
		req := echonet.NewMessage()
		req.SetESV(echonet.ESVReadRequest)
		req.SetDEOJ(0x0EF001)
		req.AddProperty(echonet.NewProperty().SetCode(0x8A))
		res, err := ctrl.PostMessage(node, req)
		if err == nil {
			if props := res.Properties(); len(props) == 1 {
				manufacture, ok := db.FindManufacture(echonet.ManufactureCode(encoding.ByteToInteger(props[0].Data())))
				if ok {
					manufactureName = manufacture.Name
				}
			}
		}

		fmt.Printf("[%d] %-15s:%d (%s)\n", i, node.Address(), node.Port(), manufactureName)

		for j, obj := range node.Objects() {
			fmt.Printf("  [%d] %06X (%s)\n", j, obj.Code(), obj.ClassName())
			for _, prop := range obj.Properties() {
				propData := ""
				if prop.IsReadable() {
					req := echonet.NewMessage()
					req.SetESV(echonet.ESVReadRequest)
					req.SetDEOJ(obj.Code())
					req.AddProperty(echonet.NewProperty().SetCode(prop.Code()))
					res, err := ctrl.PostMessage(node, req)
					if err == nil {
						if props := res.Properties(); len(props) == 1 {
							propData = hex.EncodeToString(props[0].Data())
						}
					}
				}
				fmt.Printf("    [%02X] %s: %s\n", prop.Code(), prop.Name(), propData)
			}
		}
	}

	// Stop the controller

	err = ctrl.Stop()
	if err != nil {
		return
	}
}
