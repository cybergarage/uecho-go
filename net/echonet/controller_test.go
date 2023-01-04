// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"testing"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	testControllerNodeCount = 10
)

type testController struct {
	*Controller
	foundTestNodeCount int
}

func newTestController() *testController {
	ctrl := &testController{
		Controller:         NewController(),
		foundTestNodeCount: 0,
	}
	ctrl.SetListener(ctrl)
	return ctrl
}

func (ctrl *testController) ControllerMessageReceived(*protocol.Message) {

}

func (ctrl *testController) ControllerNewNodeFound(node *RemoteNode) {
	_, err := node.FindObject(testLightDeviceCode)
	if err != nil {
		return
	}

	ctrl.foundTestNodeCount++
}

func TestNewController(t *testing.T) {
	ctrl := NewController()
	ctrl.SetConfig(newTestDefaultConfig())

	err := ctrl.Start()
	if err != nil {
		t.Error(err)
	}

	err = ctrl.SearchAllObjects()
	if err != nil {
		t.Error(err)
	}

	// time.Sleep(time.Second * 120)

	err = ctrl.Stop()
	if err != nil {
		t.Error(err)
	}
}

func testControllerSearchWithConfig(t *testing.T, config *Config) {
	t.Helper()

	// Create test nodes

	nodes := make([]*testLocalNode, testControllerNodeCount)
	for n := 0; n < testControllerNodeCount; n++ {
		node, err := newTestSampleNode()
		if err != nil {
			t.Error(err)
			return
		}
		nodes[n] = node
	}

	// Start a test node

	for _, node := range nodes {
		node.SetConfig(config)
		err := node.Start()
		if err != nil {
			t.Error(err)
			return
		}
		defer node.Stop()
	}

	// Start a controller

	ctrl := newTestController()
	ctrl.SetConfig(config)
	err := ctrl.Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer ctrl.Stop()

	err = ctrl.SearchAllObjects()
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(time.Millisecond * 200 * testControllerNodeCount)

	// Check a found device by the listener

	if ctrl.foundTestNodeCount < testControllerNodeCount {
		if ctrl.foundTestNodeCount == 0 {
			t.Errorf("Any nodes are not found (%d < %d)", ctrl.foundTestNodeCount, testControllerNodeCount)
			return
		}
		t.Skipf("%d < %d", ctrl.foundTestNodeCount, testControllerNodeCount)
	}

	if ctrl.foundTestNodeCount != testControllerNodeCount {
		for foundNodeIdx, foundNode := range ctrl.Nodes() {
			isTestNode := false
			for _, node := range nodes {
				if node.Equals(foundNode) {
					isTestNode = true
					break
				}
			}

			if !isTestNode {
				t.Skipf("[%d] %s:%d is an unknow node", foundNodeIdx, foundNode.Address(), foundNode.Port())
			}
		}
	}

	// Send read / write request (post)

	for foundNodeIdx, foundNode := range ctrl.Nodes() {
		isTestNode := false
		for _, node := range nodes {
			if node.Equals(foundNode) {
				isTestNode = true
				break
			}
		}

		if !isTestNode {
			t.Skipf("[%d] %s:%d is an unknow node", foundNodeIdx, foundNode.Address(), foundNode.Port())
			continue
		}

		// Skip other Echonet nodes
		_, err := foundNode.FindDevice(testLightDeviceCode)
		if err != nil {
			continue
		}
		for n := 0; n < testNodeRequestCount; n++ {
			var lastLightPowerStatus byte
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
			resMsg, err := ctrl.PostMessage(foundNode, reqMsg)
			if err != nil {
				t.Errorf("[%d] %s:%d is not responding", foundNodeIdx, foundNode.Address(), foundNode.Port())
				t.Error(err)
				return
			}
			if err := localNodeCheckResponseMessagePowerStatus(reqMsg, resMsg, lastLightPowerStatus); err != nil {
				t.Error(err)
				return
			}
		}
	}
}

func TestController(t *testing.T) {
	log.SetStdoutDebugEnbled(true)
	defer log.SetStdoutDebugEnbled(false)

	t.Run("Default", func(t *testing.T) {
		conf := newTestDefaultConfig()
		conf.SetConnectionTimeout(testNodeRequestTimeout)
		testControllerSearchWithConfig(t, conf)
	})

	t.Run("TCPEnabled", func(t *testing.T) {
		conf := newTestDefaultConfig()
		conf.SetConnectionTimeout(testNodeRequestTimeout)
		conf.SetTCPEnabled(true)
		testControllerSearchWithConfig(t, conf)
	})
}
