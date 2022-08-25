// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
	"net"
	"strconv"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
	"github.com/cybergarage/uecho-go/net/echonet/transport"
)

const (
	errorInvalidNotificationMessage = "invalid notification message : %s"
)

// RemoteNode is an instance for Echonet node.
type RemoteNode struct {
	*baseNode
	address string
	port    int
}

// NewRemoteNode returns a new remote node.
func NewRemoteNode() *RemoteNode {
	node := &RemoteNode{
		baseNode: newBaseNode(),
		address:  "",
		port:     0,
	}

	return node
}

// NewRemoteNodeWithRequestMessage returns a new node with the specified request message.
func NewRemoteNodeWithRequestMessage(msg *protocol.Message) *RemoteNode {
	node := NewRemoteNode()
	node.SetAddress(msg.From.IP.String())
	node.SetPort(msg.From.Port)
	return node
}

// NewRemoteNodeWithInstanceListMessage returns a new node with the specified notification message.
func NewRemoteNodeWithInstanceListMessage(msg *protocol.Message) (*RemoteNode, error) {
	if msgOPC := msg.OPC(); msgOPC < 1 {
		return nil, fmt.Errorf(errorInvalidNotificationMessage, msg)
	}

	prop := msg.PropertyAt(0)
	if prop == nil {
		return nil, fmt.Errorf(errorInvalidNotificationMessage, msg)
	}

	if prop.Code() != NodeProfileClassInstanceListNotification && prop.Code() != NodeProfileClassSelfNodeInstanceListS {
		return nil, fmt.Errorf(errorInvalidNotificationMessage, msg)
	}

	propData := prop.Data()
	propSize := len(propData)
	if propSize == 0 {
		return nil, fmt.Errorf(errorInvalidNotificationMessage, msg)
	}

	instanceCount := int(propData[0])
	if propSize < ((instanceCount * ObjectCodeSize) + 1) {
		return nil, fmt.Errorf(errorInvalidNotificationMessage, msg)
	}

	// Create a new remote Node

	node := NewRemoteNode()
	node.SetAddress(msg.SourceAddress())
	node.SetPort(msg.SourcePort())

	for n := 0; n < instanceCount; n++ {
		objCodes := make([]byte, ObjectCodeSize)
		copy(objCodes, propData[((n*ObjectCodeSize)+1):])
		obj, err := NewObjectWithCodes(objCodes)
		if err != nil {
			return nil, err
		}
		switch v := obj.(type) {
		case *Device:
			node.AddDevice(v)
		case *Profile:
			node.AddProfile(v)
		}
	}

	return node, nil
}

// NewRemoteNodeWithInstanceListMessageAndConfig returns a new node with the specified notification message and configuration.
func NewRemoteNodeWithInstanceListMessageAndConfig(msg *protocol.Message, conf *transport.Config) (*RemoteNode, error) {
	remoteNode, err := NewRemoteNodeWithInstanceListMessage(msg)
	if err != nil {
		return nil, err
	}

	// Set a default static UDP port number when the auto port binding option is disabled
	if !conf.IsAutoPortBindingEnabled() {
		remoteNode.SetPort(transport.UDPPort)
	}

	return remoteNode, nil
}

// SetAddress set the address to the node.
func (node *RemoteNode) SetAddress(addr string) {
	node.address = addr
}

// Address returns the address of the node.
func (node *RemoteNode) Address() string {
	return node.address
}

// SetPort set a port to the node.
func (node *RemoteNode) SetPort(port int) {
	node.port = port
}

// GetPort returns the port of the node.
func (node *RemoteNode) Port() int {
	return node.port
}

// AddDevice adds a new device into the node, and set the node profile and manufacture code.
func (node *RemoteNode) AddDevice(dev *Device) {
	node.baseNode.AddDevice(dev)
	dev.SetParentNode(node)
}

// AddProfile adds a new profile object into the node, and set the node profile and manufacture code.
func (node *RemoteNode) AddProfile(prof *Profile) {
	node.baseNode.AddProfile(prof)
	prof.SetParentNode(node)
}

// Equals returns true whether the specified node is same, otherwise false.
func (node *RemoteNode) Equals(otherNode Node) bool {
	return nodeEquals(node, otherNode)
}

// String returns the node string representation.
func (node *RemoteNode) String() string {
	return net.JoinHostPort(node.Address(), strconv.Itoa(node.Port()))
}
