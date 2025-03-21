// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	logLocalNodeSendMessageFormat = "LocalNode::SendMessage : %s (%d)"
	logLocalNodePostMessageFormat = "LocalNode::PostMessage : %s"
)

const (
	errorNodeRequestTimeout = "request timeout : %v"
	errorNodeIsNotRunning   = "Node (%s) is not running "
)

// AnnounceMessage announces a message.
func (node *LocalNode) AnnounceMessage(msg *protocol.Message) error {
	if !node.IsRunning() {
		return fmt.Errorf(errorNodeIsNotRunning, node)
	}
	msg.SetTID(node.NextTID())
	msg.SetDEOJ(NodeProfileObjectCode)
	return node.server.AnnounceMessage(msg)
}

// AnnounceProperty announces a specified property.
func (node *LocalNode) AnnounceProperty(prop *Property) error {
	msg := protocol.NewMessage()
	msg.SetESV(protocol.ESVNotification)
	msg.SetSEOJ(prop.ParentObject().Code())
	msg.AddProperty(prop.toProtocolProperty())
	return node.AnnounceMessage(msg)
}

// Announce announces the node.
func (node *LocalNode) Announce() error {
	// 4.3.1 Basic Sequence for ECHONET Lite Node Startup

	nodePropObj, err := node.NodeProfile()
	if err != nil {
		return err
	}

	nodeProp, ok := nodePropObj.FindProperty(NodeProfileClassInstanceListNotification)
	if !ok {
		return fmt.Errorf(errorObjectProfileObjectNotFound)
	}

	return node.AnnounceProperty(nodeProp)
}

// updateMessageDestinationHeader update the message header using the local node status.
func (node *LocalNode) updateMessageDestinationHeader(msg *protocol.Message) error {
	msg.SetTID(node.NextTID())

	// SEOJ
	nodeProp, err := node.NodeProfile()
	if err != nil {
		return err
	}
	msg.SetSEOJ(nodeProp.ParentObject().Code())

	return err
}

// SendMessage sends a message to the destination node.
func (node *LocalNode) SendMessage(dstNode Node, msg *Message) error {
	if !node.IsRunning() {
		return fmt.Errorf(errorNodeIsNotRunning, node)
	}

	err := node.updateMessageDestinationHeader(msg.protocolMessage())
	if err != nil {
		return err
	}

	_, err = node.server.SendMessage(dstNode.Address(), dstNode.Port(), msg.protocolMessage())

	// log.Trace(logLocalNodeSendMessageFormat, msg.String(), n))

	return err
}

// postMessageSynchronously posts a message to the destination node using a TCP connection and gets the response message.
func (node *LocalNode) postMessageSynchronously(dstNode Node, reqMsg *protocol.Message) (*protocol.Message, error) {
	if !node.IsRunning() {
		return nil, fmt.Errorf(errorNodeIsNotRunning, node)
	}

	err := node.updateMessageDestinationHeader(reqMsg)
	if err != nil {
		return nil, err
	}

	resMsg, err := node.server.PostMessage(dstNode.Address(), dstNode.Port(), reqMsg)

	// log.Trace(logLocalNodeSendMessageFormat, msg.String(), n))

	return resMsg, err
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
	if msg.Equals(node.postRequestMsg) {
		return false
	}
	if msg.TID() != node.postRequestMsg.TID() {
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
func (node *LocalNode) PostMessage(dstNode Node, msg *Message) (*Message, error) {
	// Use TCP connection when the function is enabled

	if node.IsTCPEnabled() {
		resMsg, err := node.postMessageSynchronously(dstNode, msg.protocolMessage())
		if err == nil {
			return newMessageWithProtocolMessage(resMsg), nil
		}
	}

	// Part V ECHONET Lite System Design Guidelines v1.12
	// Chapter 5 - Guidelines on TCP
	// A node sending a request message to another node should send the message again by UDP unicast when necessary
	//  in case of a TCP connection failure since the remote party may not be able to use TCP.

	node.Lock()
	defer node.Unlock()

	defer node.closeResponseChannel()

	node.postResponseCh = make(chan *protocol.Message)
	node.postRequestMsg = msg.protocolMessage()

	// log.Trace(logLocalNodePostMessageFormat, msg.String()))

	err := node.SendMessage(dstNode, msg)
	if err != nil {
		return nil, err
	}

	var resMsg *protocol.Message
	select {
	case resMsg = <-node.postResponseCh:
	case <-time.After(node.RequestTimeout()):
		err = fmt.Errorf(errorNodeRequestTimeout, msg)
	}

	return newMessageWithProtocolMessage(resMsg), err
}

// SendRequest sends a specified request to the object.
func (node *LocalNode) SendRequest(dstNode Node, objCode ObjectCode, esv protocol.ESV, props []*Property) error {
	return node.SendMessage(dstNode, NewMessageWithParameters(objCode, esv, props))
}

// PostRequest posts a message to the node, and wait the response message.
func (node *LocalNode) PostRequest(dstNode Node, objCode ObjectCode, esv protocol.ESV, props []*Property) (*Message, error) {
	return node.PostMessage(dstNode, NewMessageWithParameters(objCode, esv, props))
}
