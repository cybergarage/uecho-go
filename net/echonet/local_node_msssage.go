// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
	"reflect"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	errorNodeRequestInvalidDestinationNode   = "Invalid Destination Node : %v"
	errorNodeRequestInvalidDestinationObject = "Invalid Destination Object : %v"
	errorNodeRequestTimeout                  = "Request Timeout : %v"
	errorNodeIsNotRunning                    = "Node (%s) is not running "
)

// AnnounceMessage announces a message.
func (node *LocalNode) AnnounceMessage(msg *protocol.Message) error {
	if !node.IsRunning() {
		return fmt.Errorf(errorNodeIsNotRunning, node)
	}
	msg.SetTID(node.getNextTID())
	msg.SetDestinationObjectCode(NodeProfileObject)
	return node.server.NotifyMessage(msg)
}

// AnnounceProperty announces a specified property.
func (node *LocalNode) AnnounceProperty(prop *Property) error {
	msg := protocol.NewMessage()
	msg.SetESV(protocol.ESVNotification)
	msg.SetSourceObjectCode(prop.GetParentObject().GetCode())
	msg.AddProperty(prop.toProtocolProperty())
	return node.AnnounceMessage(msg)
}

// Announce announces the node
func (node *LocalNode) Announce() error {
	//4.3.1 Basic Sequence for ECHONET Lite Node Startup

	nodePropObj, err := node.GetNodeProfile()
	if err != nil {
		return err
	}

	nodeProp, ok := nodePropObj.GetProperty(NodeProfileClassInstanceListNotification)
	if !ok {
		return fmt.Errorf(errorObjectProfileObjectNotFound)
	}

	return node.AnnounceProperty(nodeProp)
}

// SendMessage sends a new specified message to the node
func (node *LocalNode) SendMessage(dstNode Node, msg *protocol.Message) error {
	if !node.IsRunning() {
		return fmt.Errorf(errorNodeIsNotRunning, node)
	}

	msg.SetTID(node.getNextTID())

	// SEOJ
	nodeProp, err := node.GetNodeProfile()
	if err != nil {
		return err
	}
	msg.SetSourceObjectCode(nodeProp.GetParentObject().GetCode())

	_, err = node.server.SendMessage(string(dstNode.GetAddress()), dstNode.GetPort(), msg)
	return err
}

// responseMessage sends a specified response message to the node
func (node *LocalNode) responseMessage(dstNode Node, msg *protocol.Message) error {
	_, err := node.server.SendMessage(dstNode.GetAddress(), dstNode.GetPort(), msg)
	return err
}

// isResponseMessageWaiting returns true when the node is waiting the response message, otherwise false.
func (node *LocalNode) isResponseMessageWaiting() bool {
	if node.postRequestMsg == nil {
		return false
	}
	if node.postResponseCh == nil {
		return false
	}
	return true
}

// isResponseMessage returns true when it is the response message, otherwise false.
func (node *LocalNode) isResponseMessage(msg *protocol.Message) bool {
	// TODO : Check the response message more strictly
	if !node.isResponseMessageWaiting() {
		return false
	}
	if msg.GetTID() != node.postRequestMsg.GetTID() {
		return false
	}
	return true
}

// setResponseMessage sets a message to the response channel.
func (node *LocalNode) setResponseMessage(msg *protocol.Message) bool {
	if !node.isResponseMessageWaiting() {
		return false
	}
	node.postResponseCh <- msg
	return true
}

// closeResponseChannel closes the response channel.
func (node *LocalNode) closeResponseChannel() {
	if node.postResponseCh == nil {
		return
	}
	close(node.postResponseCh)
	node.postResponseCh = nil
	node.postRequestMsg = nil
}

// PostMessage posts a message to the node, and wait the response message.
func (node *LocalNode) PostMessage(dstNode Node, msg *protocol.Message) (*protocol.Message, error) {
	node.Lock()
	defer node.Unlock()
	defer node.closeResponseChannel()

	node.postResponseCh = make(chan *protocol.Message)
	node.postRequestMsg = msg

	err := node.SendMessage(dstNode, msg)
	if err != nil {
		return nil, err
	}

	var resMsg *protocol.Message
	select {
	case resMsg = <-node.postResponseCh:
	case <-time.After(node.GetRequestTimeout()):
		err = fmt.Errorf(errorNodeRequestTimeout, msg)
	}

	return resMsg, err
}

// createRequestMessage creates a message with the specified parameters
func (node *LocalNode) createRequestMessage(dstNode Node, objCode ObjectCode, esv protocol.ESV, props []*Property) (*protocol.Message, error) {
	if dstNode == nil || reflect.ValueOf(dstNode).IsNil() {
		return nil, fmt.Errorf(errorNodeRequestInvalidDestinationNode, dstNode)
	}

	msg := protocol.NewMessage()
	msg.SetESV(esv)
	msg.SetDestinationObjectCode(objCode)
	for _, prop := range props {
		msg.AddProperty(prop.toProtocolProperty())
	}
	return msg, nil
}

// SendRequest sends a specified request to the object.
func (node *LocalNode) SendRequest(dstNode Node, objCode ObjectCode, esv protocol.ESV, props []*Property) error {
	msg, err := node.createRequestMessage(dstNode, objCode, esv, props)
	if err != nil {
		return err
	}
	return node.SendMessage(dstNode, msg)
}

// PostRequest posts a message to the node, and wait the response message.
func (node *LocalNode) PostRequest(dstNode Node, objCode ObjectCode, esv protocol.ESV, props []*Property) (*protocol.Message, error) {
	msg, err := node.createRequestMessage(dstNode, objCode, esv, props)
	if err != nil {
		return nil, err
	}
	return node.PostMessage(dstNode, msg)
}
