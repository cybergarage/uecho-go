// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	logNodeListenerFormat = "LocalNode::ProtocolMessageReceived : %s"
)

// ProtocolMessageReceived is a listener for the server.
func (node *LocalNode) ProtocolMessageReceived(msg *protocol.Message) (*protocol.Message, error) {
	// Ignore own messages
	msgNode := NewRemoteNodeWithRequestMessage(msg)
	if msgNode.Equals(node) {
		return nil, nil
	}

	if node.isResponseMessageWaiting() {
		if node.isResponseMessage(msg) {
			node.setResponseMessage(msg)
		}
	}

	if !node.validateReceivedMessage(msg) {
		return protocol.NewImpossibleMessageWithMessage(msg), nil
	}

	err := node.executeMessageListeners(msg)
	if err != nil {
		return nil, err
	}

	if !msg.IsResponseRequired() {
		return nil, nil
	}

	resMsg, err := node.createResponseMessageForRequestMessage(msg)
	if err != nil {
		log.Error(err.Error())
	}

	return resMsg, err
}

// validateReceivedMessage checks whether the received message is a valid message.
func (node *LocalNode) validateReceivedMessage(msg *protocol.Message) bool {
	// 4.2.2 Basic Sequences for Object Control in General

	msgDstObjCode := msg.GetDestinationObjectCode()
	msgESV := msg.GetESV()
	msgOPC := msg.GetOPC()

	// (A) Processing when the controlled object does not exist

	dstObj, err := node.FindObject(msgDstObjCode)
	if err != nil {
		// TODO : Check the DEOJ code based on Echonet specification
		return false
	}

	// (B) Processing when the controlled object exists, except when ESV = 0x60 to 0x63, 0x6E and 0x74

	if !msg.IsValidESV() { // Check only whether the ESV is valid
		return false
	}

	// (C), (D), (E)

	if msg.IsReadRequest() || msg.IsWriteRequest() {
		for n := 0; n < msgOPC; n++ {
			msgProp := msg.GetProperty(n)
			if msgProp == nil {
				continue
			}
			// (C) Processing when the controlled object exists but the controlled property does not exist or can be processed only partially
			prop, ok := dstObj.GetProperty(msgProp.Code())
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
	}

	return true
}

// executeMessageListeners post the received message to the listeners.
func (node *LocalNode) executeMessageListeners(msg *protocol.Message) error {
	msgDstObjCode := msg.GetDestinationObjectCode()
	dstObj, err := node.FindObject(msgDstObjCode)
	if err != nil {
		return err
	}

	msgESV := msg.GetESV()
	msgOPC := msg.GetOPC()

	var lastErr error

	// Message Listener

	if l := node.GetListener(); l != nil {
		err := l.NodeMessageReceived(msg)
		if err != nil {
			lastErr = err
		}
	}

	// Object Listener

	for n := 0; n < msgOPC; n++ {
		msgProp := msg.GetProperty(n)
		if msgProp == nil {
			continue
		}
		err := dstObj.notifyPropertyRequest(msgESV, msgProp)
		if err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// createResponseMessageForRequestMessage retunrs the response message for the specified request message.
func (node *LocalNode) createResponseMessageForRequestMessage(reqMsg *protocol.Message) (*protocol.Message, error) {
	msgDstObjCode := reqMsg.GetDestinationObjectCode()
	dstObj, err := node.FindObject(msgDstObjCode)
	if err != nil {
		return nil, err
	}

	msgOPC := reqMsg.GetOPC()

	resMsg := protocol.NewResponseMessageWithMessage(reqMsg)
	for n := 0; n < msgOPC; n++ {
		msgProp := reqMsg.GetProperty(n)
		if msgProp == nil {
			continue
		}
		prop, ok := dstObj.GetProperty(msgProp.Code())
		if !ok {
			continue
		}
		resMsg.AddProperty(prop.toProtocolProperty())
	}

	return resMsg, nil
}
