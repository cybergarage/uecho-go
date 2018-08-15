// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// MessageReceived is a listener for the server
func (node *LocalNode) MessageReceived(msg *protocol.Message) {
	// Ignore own messages
	msgNode := NewRemoteNodeWithRequestMessage(msg)
	if msgNode.Equals(node) {
		return
	}

	if node.isResponseMessageWaiting() {
		if node.isResponseMessage(msg) {
			node.setResponseMessage(msg)
		}
	}

	if !node.validateRequestMessage(msg) {
		if msg.IsResponseRequired() {
			node.postImpossibleResponse(msg)
		}
		return
	}

	node.execiteMessageListeners(msg)

	if msg.IsResponseRequired() {
		node.postResponseMessage(msg)
	}
}

// postImpossibleResponse returns an individual response to the source node.
func (node *LocalNode) postImpossibleResponse(msg *protocol.Message) {
	resMsg := protocol.NewImpossibleMessageWithMessage(msg)
	node.SendMessage(NewRemoteNodeWithRequestMessage(msg), resMsg)
}

// validateRequestMessage checks whether the message is a valid request.
func (node *LocalNode) validateRequestMessage(msg *protocol.Message) bool {
	//4.2.2 Basic Sequences for Object Control in General

	msgDstObjCode := msg.GetDestinationObjectCode()
	msgESV := msg.GetESV()
	msgOPC := msg.GetOPC()

	// (A) Processing when the controlled object does not exist

	dstObj, err := node.GetObject(msgDstObjCode)
	if err != nil {
		return false
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
		return false
	}

	for n := 0; n < msgOPC; n++ {
		msgProp := msg.GetProperty(n)
		if msgProp == nil {
			continue
		}
		// (C) Processing when the controlled object exists but the controlled property does not exist or can be processed only partially
		prop, ok := dstObj.GetProperty(PropertyCode(msgProp.GetCode()))
		if !ok {
			return false
		}
		// (D) Processing when the controlled property exists but the stipulated service processing functions are not available
		if !prop.IsAvailableService(msgESV) {
			return false
		}
		// (E) Processing when the controlled property exists and the stipulated service processing functions are available but the EDT size does not match
		if protocol.IsWriteRequest(msgESV) {
			if !prop.IsWritable() {
				return false
			}
			if msgProp.Size() != prop.Size() {
				return false
			}
		}
	}

	return true
}

// execiteMessageListeners post the received message to the listeners.
func (node *LocalNode) execiteMessageListeners(msg *protocol.Message) bool {
	msgDstObjCode := msg.GetDestinationObjectCode()
	dstObj, err := node.GetObject(msgDstObjCode)
	if err != nil {
		return false
	}

	msgESV := msg.GetESV()
	msgOPC := msg.GetOPC()

	// Message Listener

	l := node.GetListener()
	if l != nil {
		l.MessageReceived(msg)
	}

	// Object Listener

	for n := 0; n < msgOPC; n++ {
		msgProp := msg.GetProperty(n)
		if msgProp == nil {
			continue
		}
		dstObj.notifyPropertyRequest(msgESV, msgProp)
	}

	return true
}

// postResponseMessage posts the response message to the destination node.
func (node *LocalNode) postResponseMessage(msg *protocol.Message) bool {
	msgDstObjCode := msg.GetDestinationObjectCode()
	dstObj, err := node.GetObject(msgDstObjCode)
	if err != nil {
		return false
	}

	msgOPC := msg.GetOPC()

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

	return true
}
