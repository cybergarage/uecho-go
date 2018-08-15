// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"testing"
	"time"
)

func TestNewController(t *testing.T) {
	ctrl := NewController()

	err := ctrl.Start()
	if err != nil {
		t.Error(err)
	}

	err = ctrl.Stop()
	if err != nil {
		t.Error(err)
	}
}

func TestControllerSearch(t *testing.T) {
	// Start a test node

	node, err := newTestSampleNode()
	if err != nil {
		t.Error(err)
		return
	}

	err = node.Start()
	if err != nil {
		t.Error(err)
		return
	}

	// Start a controller

	ctrl := NewController()
	err = ctrl.Start()
	if err != nil {
		node.Stop()
		t.Error(err)
		return
	}

	err = ctrl.SearchAllObjects()
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Second)

	// Check a found node

	foundNodes := ctrl.GetNodes()
	if 0 < len(foundNodes) {
		foundNode := foundNodes[0]
		if !node.Equals(foundNode) {
			t.Errorf(errorNodeNotFound, foundNode.GetAddress(), foundNode.GetPort())
		}
	} else {
		t.Errorf(errorNodeNotFound, node.GetAddress(), node.GetPort())
	}

	// Check a found device

	_, err = ctrl.GetObject(testLightDeviceCode)
	if err != nil {
		t.Error(err)
	}

	// Finalize

	err = node.Stop()
	if err != nil {
		t.Error(err)
	}

	err = ctrl.Stop()
	if err != nil {
		t.Error(err)
	}

}
