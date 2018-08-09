// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"testing"
	"time"
)

const (
	testNodeManufacturerCode   = NodeManufacturerUnknown + 1
	testLightDeviceCode        = 0x029101
	testLightPropertyPowerCode = 0x80
	testLightPropertyPowerOn   = 0x30
	testLightPropertyPowerOff  = 0x31
)

const (
	errorNodeNotFound = "Node Not Found : %s:%d"
)

type testSampleNode struct {
	*LocalNode
}

// newTestSampleNode returns a test node.
func newTestSampleNode() (*testSampleNode, error) {
	node := &testSampleNode{
		LocalNode: NewLocalNode(),
	}

	node.SetManufacturerCode(testNodeManufacturerCode)

	dev := NewDevice()
	dev.SetCode(testLightDeviceCode)
	dev.CreateProperty(testLightPropertyPowerCode, PropertyAttributeReadWrite)
	powerData := []byte{testLightPropertyPowerOff}
	err := dev.SetPropertyData(testLightPropertyPowerCode, powerData)
	if err != nil {
		return nil, err
	}
	err = node.AddDevice(dev)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func TestNewSampleNode(t *testing.T) {
	ctrl := NewController()

	node, err := newTestSampleNode()
	if err != nil {
		t.Error(err)
		return
	}

	// Start

	err = ctrl.Start()
	if err != nil {
		t.Error(err)
		return
	}

	err = node.Start()
	if err != nil {
		t.Error(err)
		err = ctrl.Stop()
		if err != nil {
			t.Error(err)
		}
		return
	}
	// Check device

	time.Sleep(time.Second)

	foundNodes := ctrl.GetNodes()
	if len(foundNodes) <= 0 {
		t.Errorf(errorNodeNotFound, node.GetAddress(), node.GetPort())
	}

	foundNode := foundNodes[0]
	if !node.Equals(foundNode) {
		t.Errorf(errorNodeNotFound, foundNode.GetAddress(), foundNode.GetPort())
	}

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
