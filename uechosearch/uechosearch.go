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
	"fmt"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet"
)

func main() {
	ctrl := echonet.NewController()

	err := ctrl.Start()
	if err != nil {
		return
	}

	err = ctrl.SearchAllObjects()
	if err != nil {
		return
	}

	time.Sleep(time.Second)

	for _, node := range ctrl.GetNodes() {
		objs := node.GetObjects()
		if len(objs) <= 0 {
			fmt.Printf("%s\n", node.GetAddress())
			continue
		}
		for _, obj := range objs {
			fmt.Printf("%s %06X\n", node.GetAddress(), obj.GetCode())
		}
	}

	err = ctrl.Stop()
	if err != nil {
		return
	}
}
