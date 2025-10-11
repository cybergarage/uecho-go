// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	testNodeManufacturerCode            = NodeManufacturerExperimental + 1
	testLightDeviceCode                 = 0x029101
	testLightPropertyPowerCode          = 0x80
	testLightPropertyPowerOn            = 0x30
	testLightPropertyPowerOff           = 0x31
	testLightPropertyInitialPowerStatus = testLightPropertyPowerOff
)

const (
	errTestNodeNotFound = "%v: node (%s:%d)"
)

type testLocalNode struct {
	LocalNode
}

func newTestSampleNodeWithConfig(config *Config) (*testLocalNode, error) {
	dev, err := NewDevice(
		WithDeviceCode(testLightDeviceCode),
	)
	if err != nil {
		return nil, err
	}
	powerData := []byte{testLightPropertyInitialPowerStatus}
	err = dev.SetPropertyData(testLightPropertyPowerCode, powerData)
	if err != nil {
		return nil, err
	}

	node := &testLocalNode{
		LocalNode: NewLocalNode(
			WithLocalNodeConfig(config),
			WithLocalNodeDevices(dev),
			WithLocalNodeManufacturerCode(testNodeManufacturerCode),
		),
	}
	node.SetListener(node)

	return node, nil
}

func newTestSampleNode(config *Config) (*testLocalNode, error) {
	return newTestSampleNodeWithConfig(config)
}

func (node *testLocalNode) OnMessage(msg *protocol.Message) error {
	dev, err := node.LookupDevice(testLightDeviceCode)
	if err != nil {
		return err
	}

	if msg.ESV().IsWriteRequest() {
		for _, msgProp := range msg.Properties() {
			if msgProp.Code() == testLightPropertyPowerCode {
				prop, ok := dev.LookupProperty(testLightPropertyPowerCode)
				if ok {
					prop.SetData(msgProp.Data())
				}
			}
		}
	}

	return nil
}
