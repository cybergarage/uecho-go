// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
uechosearch is a search utility for Echonet Lite.

	NAME
	uecholight

	SYNOPSIS
	uecholight [OPTIONS]

	DESCRIPTION
	uecholight is a sample implementation of mono functional lighting device

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"os"
	"github.com/cybergarage/uecho-go/net/uecho"
)

// See : APPENDIX Detailed Requirements for ECHONET Device objects
//       3.3.29 Requirements for mono functional lighting class


import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/uecho"
)

func main() {

	node := uecho.NewNode()

	dev := NewLight()
	err := node.AddDevice(dev)
	if err != nil {
		os.Exit(EXIT_FAILURE)
	}

	err := node.Start()
	if err != nil {
		os.Exit(EXIT_FAILURE)
	}

	err = node.Stop()
	if err != nil {
		os.Exit(EXIT_FAILURE)
	}

	os.Exit(EXIT_SUCCESS)
}

