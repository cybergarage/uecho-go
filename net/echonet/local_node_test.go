// Copyright (C) 2018 Satoshi Konno. All rights reserved.
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

func TestNewSampleNode(t *testing.T) {
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

	ctrlAddr := ctrl.GetAddress()
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

	nodeAddr := node.GetAddress()
	nodePort := node.GetPort()

	if ctrlAddr == nodeAddr {
		//t.Errorf("%s == %s", ctrlAddr, nodeAddr)
	}

	if ctrlPort == nodePort {
		t.Errorf("%d == %d", ctrlPort, nodePort)
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
		err = ctrl.SendRequest(dev.GetParentNode(), dev, protocol.ESVReadRequest, []*Property{prop})
		if err != nil {
			t.Error(err)
		}
	}

	// Send read request (post)

	prop = NewPropertyWithCode(testLightPropertyPowerCode)
	for n := 0; n < testNodeRequestCount; n++ {
		resMsg, err := ctrl.PostRequest(dev.GetParentNode(), dev, protocol.ESVReadRequest, []*Property{prop})
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
		err = ctrl.SendRequest(dev.GetParentNode(), dev, protocol.ESVWriteRequest, []*Property{prop})
		if err != nil {
			t.Error(err)
		}

		// Read

		prop = NewPropertyWithCode(testLightPropertyPowerCode)
		resMsg, err := ctrl.PostRequest(dev.GetParentNode(), dev, protocol.ESVReadRequest, []*Property{prop})
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
		resMsg, err := ctrl.PostRequest(dev.GetParentNode(), dev, protocol.ESVWriteReadRequest, []*Property{prop})
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
