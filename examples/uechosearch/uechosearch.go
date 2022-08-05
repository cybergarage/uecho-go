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
	"flag"
	"fmt"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/log"
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

	for _, node := range ctrl.GetNodes() {
		objs := node.GetObjects()
		if len(objs) == 0 {
			fmt.Printf("%-15s\n", node.GetAddress())
			continue
		}
		for _, obj := range objs {
			fmt.Printf("%-15s : %06X\n", node.GetAddress(), obj.GetCode())
		}
	}

	// Stop the controller

	err = ctrl.Stop()
	if err != nil {
		return
	}
}
