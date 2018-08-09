// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

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
