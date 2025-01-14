// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"testing"
)

func TestLightNode(t *testing.T) {
	node := NewLightNode()

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
