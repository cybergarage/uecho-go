// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"testing"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	testNodeRequestCount = 10
)

const (
	errorLocalNodeTestInvalidResponse     = "Invalid Respose : %s"
	errorLocalNodeTestInvalidPropertyData = "Invalid Respose Status : %X != %X"
)

func TestNewLocalNode(t *testing.T) {
	node := NewLocalNode()

	_, err := node.GetNodeProfile()
	if err != nil {
		t.Error(err)
	}
}

func localNodeCheckResponseMessagePowerStatus(t *testing.T, resMsg *protocol.Message, powerStatus byte) {
	resOpc := resMsg.GetOPC()
	if resOpc == 1 {
		resProp := resMsg.GetProperty(0)
		if resProp != nil || (resProp.GetCode() == testLightPropertyPowerCode) {
			resData := resProp.GetData()
			if len(resData) == 1 {
				if resData[0] != powerStatus {
					t.Errorf(errorLocalNodeTestInvalidPropertyData, resData[0], powerStatus)
				}
			} else {
				t.Errorf(errorLocalNodeTestInvalidResponse, resMsg)
			}
		} else {
			t.Errorf(errorLocalNodeTestInvalidResponse, resMsg)
		}
	} else {
		t.Errorf(errorLocalNodeTestInvalidResponse, resMsg)
	}
}

func testLocalNodeWithConfig(t *testing.T, config *Config) {
	ctrl := NewController()

	node, err := newTestSampleNode()
	if err != nil {
		t.Error(err)
		return
	}

	startTID := ctrl.GetLastTID()

	// Start

	err = ctrl.Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer ctrl.Stop()

	ctrlPort := ctrl.GetPort()

	err = node.Start()
	if err != nil {
		t.Error(err)
		err = ctrl.Stop()
		if err != nil {
			t.Error(err)
		}
		return
	}
	defer node.Stop()

	nodePort := node.GetPort()
	if ctrlPort == nodePort {
		t.Errorf("%d == %d", ctrlPort, nodePort)
	}

	time.Sleep(time.Second)

	// Check a found node

	var foundNode *RemoteNode
	for _, ctrlNode := range ctrl.GetNodes() {
		// Skip other Echonet nodes
		_, err := ctrlNode.GetDevice(testLightDeviceCode)
		if err != nil {
			continue
		}
		foundNode = ctrlNode
		break
	}

	if foundNode == nil {
		t.Errorf(errorNodeNotFound, node.GetAddress(), node.GetPort())
		return
	}
	// Check a found device

	dev, err := ctrl.GetObject(testLightDeviceCode)
	if err != nil {
		t.Error(err)
		return
	}

	// Send read request

	prop := NewPropertyWithCode(testLightPropertyPowerCode)
	for n := 0; n < testNodeRequestCount; n++ {
		err = ctrl.SendRequest(dev.GetParentNode(), testLightDeviceCode, protocol.ESVReadRequest, []*Property{prop})
		if err != nil {
			t.Error(err)
		}
	}

	// Send read request (post)

	prop = NewPropertyWithCode(testLightPropertyPowerCode)
	for n := 0; n < testNodeRequestCount; n++ {
		resMsg, err := ctrl.PostRequest(dev.GetParentNode(), testLightDeviceCode, protocol.ESVReadRequest, []*Property{prop})
		if err == nil {
			localNodeCheckResponseMessagePowerStatus(t, resMsg, testLightPropertyInitialPowerStatus)
		} else {
			t.Error(err)
		}
	}

	// Send write request (off <-> on)

	var lastLightPowerStatus byte

	for n := 0; n < testNodeRequestCount; n++ {
		if (n % 2) == 0 {
			lastLightPowerStatus = testLightPropertyPowerOn
		} else {

			lastLightPowerStatus = testLightPropertyPowerOff
		}

		// Write

		prop := NewPropertyWithCode(testLightPropertyPowerCode)
		prop.SetData([]byte{lastLightPowerStatus})
		err = ctrl.SendRequest(dev.GetParentNode(), testLightDeviceCode, protocol.ESVWriteRequest, []*Property{prop})
		if err != nil {
			t.Error(err)
		}

		// Read

		prop = NewPropertyWithCode(testLightPropertyPowerCode)
		resMsg, err := ctrl.PostRequest(dev.GetParentNode(), testLightDeviceCode, protocol.ESVReadRequest, []*Property{prop})
		if err == nil {
			localNodeCheckResponseMessagePowerStatus(t, resMsg, lastLightPowerStatus)
		} else {
			t.Error(err)
		}
	}

	// Send read / write request (post)

	prop = NewPropertyWithCode(testLightPropertyPowerCode)
	for n := 0; n < testNodeRequestCount; n++ {
		if (n % 2) == 0 {
			lastLightPowerStatus = testLightPropertyPowerOn
		} else {

			lastLightPowerStatus = testLightPropertyPowerOff
		}

		// Write / Read

		prop := NewPropertyWithCode(testLightPropertyPowerCode)
		prop.SetData([]byte{lastLightPowerStatus})
		resMsg, err := ctrl.PostRequest(dev.GetParentNode(), testLightDeviceCode, protocol.ESVWriteReadRequest, []*Property{prop})
		if err == nil {
			localNodeCheckResponseMessagePowerStatus(t, resMsg, lastLightPowerStatus)
		} else {
			t.Error(err)
		}
	}

	lastTID := ctrl.GetLastTID()
	if lastTID < startTID {
		t.Errorf("%d < %d", lastTID, startTID)
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

func TestLocalNodeWithDefaultConfig(t *testing.T) {
	testLocalNodeWithConfig(t, NewDefaultConfig())
}

func TestLocalNodeWithOnlyUDPConfig(t *testing.T) {
	conf := NewDefaultConfig()
	conf.SetTCPEnabled(false)
	conf.SetUDPEnabled(true)
	testLocalNodeWithConfig(t, conf)
}

/*
func TestLocalNodeWithEnableTCPConfig(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	conf := NewDefaultConfig()
	conf.SetTCPEnabled(true)
	conf.SetUDPEnabled(true)
	testLocalNodeWithConfig(t, conf)
}
*/
