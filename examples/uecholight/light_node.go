// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"encoding/hex"

	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type LightNode struct {
	*echonet.LocalNode
}

// NewLightNode returns a new light device.
func NewLightNode() *LightNode {

	node := &LightNode{
		LocalNode: echonet.NewLocalNode(),
	}

	dev := NewLightDevice()
	node.AddDevice(dev)
	dev.SetListener(node)

	return node
}

func (node *LightNode) PropertyRequestReceived(obj *echonet.Object, esv protocol.ESV, reqProp *protocol.Property) error {
	// Check whether the property request is a write request
	if !protocol.IsWriteRequest(esv) {
		return nil
	}

	// Check whether the local object (device) has the requested property
	propCode := reqProp.GetCode()
	prop, ok := obj.GetProperty(propCode)
	if !ok {
		return nil
	}

	OutputMessage("%02X : %s -> %s", esv, hex.EncodeToString(prop.GetData()), hex.EncodeToString(reqProp.GetData()))

	// Set the requested data to the local object (device)
	prop.SetData(reqProp.GetData())

	return nil
}
