// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type LightNode struct {
	echonet.LocalNode
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
	// Check whether the property request is a write request. Basically, the developer should handle only write requests.

	if !protocol.IsWriteRequest(esv) {
		return nil
	}

	// Check whether the local object (device) has the requested property

	reqPropCode := reqProp.Code()
	reqPropData := reqProp.Data()

	switch reqPropCode {
	case 0x80: // Operation status
		// Check whether the request data is 0x30(ON) or 0x31(OFF)
		if (len(reqPropData) != 1) || (reqPropData[0] != 0x30) || (reqPropData[0] != 0x31) {
			err := fmt.Errorf("Invalid Request : %02X %s", reqPropCode, hex.EncodeToString(reqPropData))
			OutputError(err)
			return err
		}
	default:
		err := fmt.Errorf("Invalid Request : %02X %s", reqPropCode, hex.EncodeToString(reqPropData))
		OutputError(err)
		return err
	}

	// Output the update message
	// NOTE : Object::GetProperty() can get the specified property always because the PropertyRequestReceived is not called when the object has no the specified property

	targetProp, _ := obj.FindProperty(reqPropCode)
	OutputMessage("0x%02X : 0x%s -> 0x%s", esv, hex.EncodeToString(targetProp.Data()), hex.EncodeToString(reqPropData))

	// Set the requested data to the local object (device)

	targetProp.SetData(reqPropData)

	return nil
}
