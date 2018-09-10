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

func TestControllerSearch(t *testing.T) {
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
		err := node.Start()
		if err != nil {
			t.Error(err)
			return
		}
		defer node.Stop()
	}

	// Start a controller

	ctrl := newTestController()
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
