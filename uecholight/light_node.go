// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type LightNode struct {
	echonet.Node
}

// NewLightNode returns a new light device.
func NewLightNode() *LightNode {

	node := &LightNode{
		Node: echonet.NewLocalNode(),
	}

	dev := NewLightDevice()
	node.AddDevice(dev)
	dev.AddListener(node)

	return node
}

// PropertyRequestReceived is a listener for Echonet requests.
func (node *LightNode) PropertyRequestReceived(obj *echonet.Object, esv protocol.ESV, reqProp *protocol.Property) error {
	if !protocol.IsWriteRequest(esv) {
		return nil
	}

	propCode := reqProp.GetCode()
	prop, ok := obj.GetProperty(propCode)
	if !ok {
		return nil
	}

	prop.SetData(reqProp.GetData())

	return nil
}
