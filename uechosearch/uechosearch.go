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
	"fmt"

	"github.com/cybergarage/uecho-go/net/uecho"
)

func main() {
	ctrl := uecho.NewController()

	err := ctrl.Start()
	if err != nil {
		return
	}

	err = ctrl.SearchAllObjects()
	if err != nil {
		return
	}

	for _, node := range ctrl.GetNodes() {
		objs := node.GetObjects()
		if len(objs) <= 0 {
			fmt.Printf("%s\n", node.GetAddress())
			continue
		}
		for _, obj := range objs {
			fmt.Printf("%s %06X\n", node.GetAddress(), obj.getCode())
		}
	}

	err = ctrl.Stop()
	if err != nil {
		return
	}
}
