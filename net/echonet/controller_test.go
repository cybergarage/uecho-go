// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"context"
	"testing"
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	testControllerNodeCount = 10
)

type testController struct {
	Controller

	foundTestNodeCount int
}

func newTestController(opts ...ControllerOption) *testController {
	ctrl := &testController{
		Controller:         NewController(opts...),
		foundTestNodeCount: 0,
	}
	ctrl.SetListener(ctrl)
	return ctrl
}

func (ctrl *testController) ControllerMessageReceived(*protocol.Message) {

}

func (ctrl *testController) ControllerNewNodeFound(node Node) {
	_, err := node.LookupObject(testLightDeviceCode)
	if err != nil {
		return
	}

	ctrl.foundTestNodeCount++
}

func TestNewController(t *testing.T) {
	ctrl := NewController(
		WithControllerConfig(newTestDefaultConfig()),
	)

	err := ctrl.Start()
	if err != nil {
		t.Error(err)
	}

	err = ctrl.Search(context.Background())
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
	for n := range testControllerNodeCount {
		node, err := newTestSampleNode(config)
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

	ctrl := newTestController(
		WithControllerConfig(config),
	)
	err := ctrl.Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer ctrl.Stop()

	err = ctrl.Search(context.Background())
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
		_, err := foundNode.LookupDevice(testLightDeviceCode)
		if err != nil {
			continue
		}
		for n := range testNodeRequestCount {
			var lastLightPowerStatus byte
			if (n % 2) == 0 {
				lastLightPowerStatus = testLightPropertyPowerOn
			} else {
				lastLightPowerStatus = testLightPropertyPowerOff
			}

			// WriteRead

			time.Sleep(testNodeRequestSleep)

			prop := NewProperty(
				WithPropertyCode(testLightPropertyPowerCode),
			)
			prop.SetData([]byte{lastLightPowerStatus})
			reqMsg := NewMessage(
				WithMessageDEOJ(testLightDeviceCode),
				WithMessageESV(protocol.ESVWriteReadRequest),
				WithMessageProperties(prop),
			)
			resMsg, err := ctrl.PostMessage(context.Background(), foundNode, reqMsg)
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
	log.EnableStdoutDebug(true)
	defer log.EnableStdoutDebug(false)

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
