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
	_, err := node.GetObject(testLightDeviceCode)
	if err != nil {
		return
	}

	ctrl.foundTestNodeCount++
}

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

func testControllerSearchWithConfig(t *testing.T, config *Config) {
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

	err = ctrl.SearchAllObjects()
	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Millisecond * 100 * testControllerNodeCount)

	// Check a found device by the listener

	if ctrl.foundTestNodeCount < testControllerNodeCount {
		t.Errorf("%d < %d", ctrl.foundTestNodeCount, testControllerNodeCount)
	}

	// Send read / write request (post)

	for _, node := range ctrl.GetNodes() {
		// Skip other Echonet nodes
		_, err := node.GetDevice(testLightDeviceCode)
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

			// Write / Read

			prop := NewPropertyWithCode(testLightPropertyPowerCode)
			prop.SetData([]byte{lastLightPowerStatus})
			resMsg, err := ctrl.PostRequest(node, testLightDeviceCode, protocol.ESVWriteReadRequest, []*Property{prop})
			if err == nil {
				localNodeCheckResponseMessagePowerStatus(t, resMsg, lastLightPowerStatus)
			} else {
				t.Error(err)
			}
		}
	}

	// Finalize

	err = ctrl.Stop()
	if err != nil {
		t.Error(err)
	}

}

func TestControllerSearchithWithDefaultConfig(t *testing.T) {
	//log.SetStdoutDebugEnbled(true)
	//defer log.SetStdoutDebugEnbled(false)
	conf := newTestDefaultConfig()
	conf.SetConnectionTimeout(testNodeRequestTimeout)
	testControllerSearchWithConfig(t, conf)
}

func TestControllerSearchWithOnlyUDPConfig(t *testing.T) {
	//log.SetStdoutDebugEnbled(true)
	//defer log.SetStdoutDebugEnbled(false)
	conf := newTestDefaultConfig()
	conf.SetConnectionTimeout(testNodeRequestTimeout)
	conf.SetTCPEnabled(false)
	testControllerSearchWithConfig(t, conf)
}

func TestControllerSearchWithEnableTCPConfig(t *testing.T) {
	//log.SetStdoutDebugEnbled(true)
	//defer log.SetStdoutDebugEnbled(false)
	conf := newTestDefaultConfig()
	conf.SetConnectionTimeout(testNodeRequestTimeout)
	conf.SetTCPEnabled(true)
	testControllerSearchWithConfig(t, conf)
}
