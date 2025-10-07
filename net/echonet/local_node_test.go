// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cybergarage/go-logger/log"
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
	if _, err := node.NodeProfile(); err != nil {
		t.Error(err)
	}
}

func localNodeCheckResponseMessagePowerStatus(reqMsg Message, resMsg Message, powerStatus byte) error {
	if resOpc := resMsg.OPC(); resOpc != 1 {
		return fmt.Errorf(errorLocalNodeTestInvalidResponse, resMsg)
	}

	resProp, ok := resMsg.Property(0)
	if !ok {
		return fmt.Errorf(errorLocalNodeTestInvalidResponse, resMsg)
	}
	if resProp.Code() != testLightPropertyPowerCode {
		return fmt.Errorf(errorLocalNodeTestInvalidResponse, resMsg)
	}

	resData := resProp.Data()
	if len(resData) != 1 {
		return fmt.Errorf(errorLocalNodeTestInvalidResponse, resMsg)
	}
	if resData[0] != powerStatus {
		return fmt.Errorf(errorLocalNodeTestInvalidPropertyData, resData[0], powerStatus, reqMsg, resMsg)
	}

	return nil
}

// nolint ifshort
func testLocalNodeWithConfig(t *testing.T, config *Config) {
	// Setup

	ctrl := NewController()
	ctrl.SetConfig(config)

	node, err := newTestSampleNode(config)
	if err != nil {
		t.Error(err)
		return
	}

	startTID := ctrl.LastTID()

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

	ctrlPort := ctrl.Port()
	nodePort := node.Port()
	if ctrlPort == nodePort {
		t.Errorf("%d == %d", ctrlPort, nodePort)
		return
	}

	// log.Trace("controller:%d, node:%d", ctrlPort, nodePort)

	time.Sleep(time.Second)

	// Check a found node

	var foundNode Node
	for _, ctrlNode := range ctrl.Nodes() {
		// Skip deprecated nodes
		if !ctrlNode.Equals(node) {
			continue
		}
		// Skip other Echonet nodes
		_, err := ctrlNode.LookupDevice(testLightDeviceCode)
		if err != nil {
			continue
		}
		foundNode = ctrlNode
		break
	}

	if foundNode == nil {
		t.Errorf(errorTestNodeNotFound, node.Address(), node.Port())
		return
	}
	// Check a found device

	dev, err := ctrl.LookupObject(testLightDeviceCode)
	if err != nil {
		t.Error(err)
		return
	}

	// Send read request

	prop := NewPropertyWithCode(testLightPropertyPowerCode)
	for n := 0; n < testNodeRequestCount; n++ {
		err = ctrl.SendRequest(
			context.Background(),
			dev.ParentNode(),
			testLightDeviceCode,
			protocol.ESVReadRequest,
			prop)
		if err != nil {
			t.Error(err)
			return
		}
	}

	// Send read request (post)

	prop = NewPropertyWithCode(testLightPropertyPowerCode)
	for n := 0; n < testNodeRequestCount; n++ {
		time.Sleep(testNodeRequestSleep)
		reqMsg := NewMessage(
			WithMessageDEOJ(testLightDeviceCode),
			WithMessageESV(protocol.ESVReadRequest),
			WithMessageProperties(prop),
		)
		resMsg, err := ctrl.PostMessage(
			context.Background(),
			dev.ParentNode(),
			reqMsg)
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
		err = ctrl.SendRequest(
			context.Background(),
			dev.ParentNode(),
			testLightDeviceCode,
			protocol.ESVWriteRequest,
			prop)
		if err != nil {
			t.Error(err)
			return
		}

		// Read

		time.Sleep(testNodeRequestSleep)

		prop = NewPropertyWithCode(testLightPropertyPowerCode)
		reqMsg := NewMessage(
			WithMessageDEOJ(testLightDeviceCode),
			WithMessageESV(protocol.ESVReadRequest),
			WithMessageProperties(prop),
		)
		resMsg, err := ctrl.PostMessage(
			context.Background(),
			dev.ParentNode(),
			reqMsg)
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
		reqMsg := NewMessage(
			WithMessageDEOJ(testLightDeviceCode),
			WithMessageESV(protocol.ESVWriteReadRequest),
			WithMessageProperties(prop),
		)
		resMsg, err := ctrl.PostMessage(context.Background(), dev.ParentNode(), reqMsg)
		if err != nil {
			t.Error(err)
			return
		}
		if err := localNodeCheckResponseMessagePowerStatus(reqMsg, resMsg, lastLightPowerStatus); err != nil {
			t.Error(err)
			return
		}
	}

	lastTID := ctrl.LastTID()
	if lastTID < startTID {
		t.Errorf("%d < %d", lastTID, startTID)
	}
}

func TestLocalNode(t *testing.T) {
	log.EnableStdoutDebug(true)
	defer log.EnableStdoutDebug(false)

	t.Run("Default", func(t *testing.T) {
		conf := newTestDefaultConfig()
		conf.SetConnectionTimeout(testNodeRequestTimeout)
		testLocalNodeWithConfig(t, conf)
	})
	t.Run("TCPEnabled", func(t *testing.T) {
		conf := newTestDefaultConfig()
		conf.SetConnectionTimeout(testNodeRequestTimeout)
		conf.SetTCPEnabled(true)
		testLocalNodeWithConfig(t, conf)
	})
}
