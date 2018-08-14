// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

// MessageReceived is a listener for the server
func (node *LocalNode) MessageReceived(msg *protocol.Message) {
	if node.isResponseMessageWaiting() {
		if node.isResponseMessage(msg) {
			node.setResponseMessage(msg)
		}
	}
	node.executeObjectControl(msg)

	l := node.GetListener()
	if l != nil {
		l.MessageReceived(msg)
	}
}

// postImpossibleResponse returns an individual response to the source node.
func (node *LocalNode) postImpossibleResponse(msg *protocol.Message) {
	if !msg.IsResponseRequired() {
		return
	}
	resMsg := protocol.NewImpossibleMessageWithMessage(msg)
	node.SendMessage(NewRemoteNodeWithRequestMessage(msg), resMsg)
}

// executeObjectControl executes the specified message based on the Echonet specification (4.2.2 Basic Sequences for Object Control in General)
func (node *LocalNode) executeObjectControl(msg *protocol.Message) {
	//4.2.2 Basic Sequences for Object Control in General

	msgDstObjCode := msg.GetDestinationObjectCode()
	msgESV := msg.GetESV()
	msgOPC := msg.GetOPC()

	// (A) Processing when the controlled object does not exist

	dstObj, err := node.GetObject(msgDstObjCode)
	if err != nil {
		return
	}

	// (B) Processing when the controlled object exists, except when ESV = 0x60 to 0x63, 0x6E and 0x74

	switch msgESV {
	case protocol.ESVWriteRequest:
	case protocol.ESVWriteRequestResponseRequired:
	case protocol.ESVReadRequest:
	case protocol.ESVNotificationRequest:
	case protocol.ESVWriteReadRequest:
	case protocol.ESVNotificationResponseRequired:
	default:
		return
	}

	for n := 0; n < msgOPC; n++ {
		msgProp := msg.GetProperty(n)
		if msgProp == nil {
			continue
		}
		// (C) Processing when the controlled object exists but the controlled property does not exist or can be processed only partially
		prop, ok := dstObj.GetProperty(PropertyCode(msgProp.GetCode()))
		if !ok {
			node.postImpossibleResponse(msg)
			return
		}
		// (D) Processing when the controlled property exists but the stipulated service processing functions are not available
		if !prop.IsAvailableService(msgESV) {
			node.postImpossibleResponse(msg)
			return
		}
		// (E) Processing when the controlled property exists and the stipulated service processing functions are available but the EDT size does not match
		if protocol.IsWriteRequest(msgESV) {
			if !prop.IsWritable() {
				node.postImpossibleResponse(msg)
			}
			if msgProp.Size() != prop.Size() {
				node.postImpossibleResponse(msg)
				return
			}
			// FIXME : Check whether user approve the write request
			prop.SetData(msgProp.GetData())
		}
	}

	// (F) Processing when the controlled property exists, the stipulated service processing functions are available and also the EDT size matches

	for n := 0; n < msgOPC; n++ {
		msgProp := msg.GetProperty(n)
		if msgProp == nil {
			continue
		}
		dstObj.notifyPropertyRequest(msgESV, msgProp)
	}

	if !msg.IsResponseRequired() {
		return
	}

	resMsg := protocol.NewResponseMessageWithMessage(msg)
	for n := 0; n < msgOPC; n++ {
		msgProp := msg.GetProperty(n)
		if msgProp == nil {
			continue
		}
		prop, ok := dstObj.GetProperty(PropertyCode(msgProp.GetCode()))
		if !ok {
			continue
		}
		resMsg.AddProperty(prop.toProtocolProperty())
	}
	node.responseMessage(NewRemoteNodeWithRequestMessage(msg), resMsg)
}
