// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"testing"
	"time"

	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

func TestNewLocalNode(t *testing.T) {
	node := NewLocalNode()

	_, err := node.GetNodeProfile()
	if err != nil {
		t.Error(err)
	}
}

func TestNewSampleNode(t *testing.T) {
	ctrl := NewController()

	node, err := newTestSampleNode()
	if err != nil {
		t.Error(err)
		return
	}

	// Start

	err = ctrl.Start()
	if err != nil {
		t.Error(err)
		return
	}

	err = node.Start()
	if err != nil {
		t.Error(err)
		err = ctrl.Stop()
		if err != nil {
			t.Error(err)
		}
		return
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
	}

	// Check a found device

	dev, err := ctrl.GetObject(testLightDeviceCode)
	if err != nil {
		t.Error(err)
	}

	// Send requests

	prop := NewPropertyWithCode(testLightPropertyPowerCode)
	err = ctrl.SendRequest(dev.GetParentNode(), dev, protocol.ESVReadRequest, []*Property{prop})
	if err != nil {
		t.Error(err)
	}

	_, err = ctrl.PostRequest(dev.GetParentNode(), dev, protocol.ESVReadRequest, []*Property{prop})
	if err != nil {
		t.Error(err)
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
