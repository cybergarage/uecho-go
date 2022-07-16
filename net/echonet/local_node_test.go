// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
	"testing"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	testNodeRequestCount   = 4
	testNodeRequestTimeout = time.Millisecond * 1000
	testNodeRequestSleep   = time.Millisecond * 200
)

const (
	errorLocalNodeTestInvalidResponse     = "Invalid Respose : %s"
	errorLocalNodeTestInvalidPropertyData = "Invalid Respose Status : %X != %X (%s != %s)"
)

func TestNewLocalNode(t *testing.T) {
	node := NewLocalNode()

	if _, err := node.GetNodeProfile(); err != nil {
		t.Error(err)
	}
}

func localNodeCheckResponseMessagePowerStatus(reqMsg *protocol.Message, resMsg *protocol.Message, powerStatus byte) error {
	if resOpc := resMsg.GetOPC(); resOpc != 1 {
		return fmt.Errorf(errorLocalNodeTestInvalidResponse, resMsg)
	}

	resProp := resMsg.GetProperty(0)
	if resProp == nil {
		return fmt.Errorf(errorLocalNodeTestInvalidResponse, resMsg)
	}
	if resProp.GetCode() != testLightPropertyPowerCode {
		return fmt.Errorf(errorLocalNodeTestInvalidResponse, resMsg)
	}

	resData := resProp.GetData()
	if len(resData) != 1 {
		return fmt.Errorf(errorLocalNodeTestInvalidResponse, resMsg)
	}
	if resData[0] != powerStatus {
		return fmt.Errorf(errorLocalNodeTestInvalidPropertyData, resData[0], powerStatus, reqMsg, resMsg)
	}

	return nil
}

//nolint ifshort
func testLocalNodeWithConfig(t *testing.T, config *Config) {
	ctrl := NewController()
	ctrl.SetConfig(config)

	node, err := newTestSampleNode()
	if err != nil {
		t.Error(err)
		return
	}
	node.SetConfig(config)

	startTID := ctrl.GetLastTID()

	// Start

	err = ctrl.Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer ctrl.Stop()

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

	ctrlPort := ctrl.GetPort()
	nodePort := node.GetPort()
	if ctrlPort == nodePort {
		t.Errorf("%d == %d", ctrlPort, nodePort)
		return
	}

	time.Sleep(time.Second)

	// Check a found node

	var foundNode *RemoteNode
	for _, ctrlNode := range ctrl.GetNodes() {
		// Skip deprecated nodes
		if !ctrlNode.Equals(node) {
			continue
		}
		// Skip other Echonet nodes
		_, err := ctrlNode.GetDevice(testLightDeviceCode)
		if err != nil {
			continue
		}
		foundNode = ctrlNode
		break
	}

	if foundNode == nil {
		t.Errorf(errorTestNodeNotFound, node.GetAddress(), node.GetPort())
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
			return
		}
	}

	// Send read request (post)

	prop = NewPropertyWithCode(testLightPropertyPowerCode)
	for n := 0; n < testNodeRequestCount; n++ {
		time.Sleep(testNodeRequestSleep)
		reqMsg := NewMessageWithParameters(testLightDeviceCode, protocol.ESVReadRequest, []*Property{prop})
		resMsg, err := ctrl.PostMessage(dev.GetParentNode(), reqMsg)
		if err != nil {
			t.Error(err)
			return
		}
		if err := localNodeCheckResponseMessagePowerStatus(reqMsg, resMsg, testLightPropertyInitialPowerStatus); err != nil {
			t.Error(err)
			return
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

		time.Sleep(testNodeRequestSleep)

		prop := NewPropertyWithCode(testLightPropertyPowerCode)
		prop.SetData([]byte{lastLightPowerStatus})
		err = ctrl.SendRequest(dev.GetParentNode(), testLightDeviceCode, protocol.ESVWriteRequest, []*Property{prop})
		if err != nil {
			t.Error(err)
			return
		}

		// Read

		time.Sleep(testNodeRequestSleep)

		prop = NewPropertyWithCode(testLightPropertyPowerCode)
		reqMsg := NewMessageWithParameters(testLightDeviceCode, protocol.ESVReadRequest, []*Property{prop})
		resMsg, err := ctrl.PostMessage(dev.GetParentNode(), reqMsg)
		if err != nil {
			t.Error(err)
			return
		}
		if err := localNodeCheckResponseMessagePowerStatus(reqMsg, resMsg, lastLightPowerStatus); err != nil {
			t.Error(err)
			return
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

		// WriteRead

		time.Sleep(testNodeRequestSleep)

		prop := NewPropertyWithCode(testLightPropertyPowerCode)
		prop.SetData([]byte{lastLightPowerStatus})
		reqMsg := NewMessageWithParameters(testLightDeviceCode, protocol.ESVWriteReadRequest, []*Property{prop})
		resMsg, err := ctrl.PostMessage(dev.GetParentNode(), reqMsg)
		if err != nil {
			t.Error(err)
			return
		}
		if err := localNodeCheckResponseMessagePowerStatus(reqMsg, resMsg, lastLightPowerStatus); err != nil {
			t.Error(err)
			return
		}
	}

	lastTID := ctrl.GetLastTID()
	if lastTID < startTID {
		t.Errorf("%d < %d", lastTID, startTID)
	}
}

func TestLocalNodeWithDefaultConfig(t *testing.T) {
	log.SetStdoutDebugEnbled(true)
	defer log.SetStdoutDebugEnbled(false)
	conf := newTestDefaultConfig()
	conf.SetConnectionTimeout(testNodeRequestTimeout)
	testLocalNodeWithConfig(t, conf)
}

func TestLocalNodeWithOnlyUDPConfig(t *testing.T) {
	log.SetStdoutDebugEnbled(true)
	defer log.SetStdoutDebugEnbled(false)
	conf := newTestDefaultConfig()
	conf.SetConnectionTimeout(testNodeRequestTimeout)
	conf.SetTCPEnabled(false)
	testLocalNodeWithConfig(t, conf)
}

func TestLocalNodeWithEnableTCPConfig(t *testing.T) {
	log.SetStdoutDebugEnbled(true)
	defer log.SetStdoutDebugEnbled(false)
	conf := newTestDefaultConfig()
	conf.SetConnectionTimeout(testNodeRequestTimeout)
	conf.SetTCPEnabled(true)
	testLocalNodeWithConfig(t, conf)
}
