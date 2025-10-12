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
	errInvalidNotificationMessage = "%w: notification message (%s)"
)

// remoteNode is an instance for Echonet node.
type remoteNode struct {
	*baseNode

	address string
	port    int
}

// newRemoteNode returns a new remote node.
func newRemoteNode() *remoteNode {
	node := &remoteNode{
		baseNode: newBaseNode(),
		address:  "",
		port:     0,
	}

	return node
}

// newRemoteNodeWithRequestMessage returns a new node with the specified request message.
func newRemoteNodeWithRequestMessage(msg *protocol.Message) *remoteNode {
	node := newRemoteNode()
	node.SetAddress(msg.From.IP.String())
	node.SetPort(msg.From.Port)
	return node
}

// newRemoteNodeWithInstanceListMessage returns a new node with the specified notification message.
func newRemoteNodeWithInstanceListMessage(msg *protocol.Message) (*remoteNode, error) {
	if msgOPC := msg.OPC(); msgOPC < 1 {
		return nil, fmt.Errorf(errInvalidNotificationMessage, ErrInvalid, msg)
	}

	prop := msg.Property(0)
	if prop == nil {
		return nil, fmt.Errorf(errInvalidNotificationMessage, ErrInvalid, msg)
	}

	if prop.Code() != NodeProfileClassInstanceListNotification && prop.Code() != NodeProfileClassSelfNodeInstanceListS {
		return nil, fmt.Errorf(errInvalidNotificationMessage, ErrInvalid, msg)
	}

	propData := prop.Data()
	propSize := len(propData)
	if propSize == 0 {
		return nil, fmt.Errorf(errInvalidNotificationMessage, ErrInvalid, msg)
	}

	instanceCount := int(propData[0])
	if propSize < ((instanceCount * ObjectCodeSize) + 1) {
		return nil, fmt.Errorf(errInvalidNotificationMessage, ErrInvalid, msg)
	}

	// Create a new remote Node

	node := newRemoteNode()
	node.SetAddress(msg.SourceAddress())
	node.SetPort(msg.SourcePort())
	node.AddProfile(NewNodeProfile())

	newObjectWithCodeBytes := func(codes []byte) (any, error) {
		objCode, err := NewObjectCodeFromBytes(codes)
		if err != nil {
			return nil, err
		}
		if isProfileObjectCode(codes[0]) {
			obj := NewProfileWithCode(objCode)
			return obj, nil
		}
		return NewDeviceWithCode(objCode)
	}

	for n := range instanceCount {
		objCodes := make([]byte, ObjectCodeSize)
		copy(objCodes, propData[((n*ObjectCodeSize)+1):])
		obj, err := newObjectWithCodeBytes(objCodes)
		if err != nil {
			return nil, err
		}
		switch v := obj.(type) {
		case Profile:
			node.AddProfile(v)
		case Device:
			node.AddDevice(v)
		}
	}

	node.AddProfile(NewNodeProfile())

	return node, nil
}

// newRemoteNodeWithInstanceListMessageAndConfig returns a new node with the specified notification message and configuration.
func newRemoteNodeWithInstanceListMessageAndConfig(msg *protocol.Message, conf *transport.Config) (*remoteNode, error) {
	remoteNode, err := newRemoteNodeWithInstanceListMessage(msg)
	if err != nil {
		return nil, err
	}

	// Set a default static UDP port number when the auto port binding option is disabled
	if !conf.AutoPortBindingEnabled() {
		remoteNode.SetPort(transport.UDPPort)
	}

	return remoteNode, nil
}

// SetAddress set the address to the node.
func (node *remoteNode) SetAddress(addr string) {
	node.address = addr
}

// Address returns the address of the node.
func (node *remoteNode) Address() string {
	return node.address
}

// SetPort set a port to the node.
func (node *remoteNode) SetPort(port int) {
	node.port = port
}

// Port returns the port of the node.
func (node *remoteNode) Port() int {
	return node.port
}

// AddDevice adds a new device into the node, and set the node profile and manufacture code.
func (node *remoteNode) AddDevice(dev Device) {
	node.baseNode.AddDevice(dev)
	dev.SetParentNode(node)
}

// AddProfile adds a new profile object into the node, and set the node profile and manufacture code.
func (node *remoteNode) AddProfile(prof Profile) {
	node.baseNode.AddProfile(prof)
	prof.SetParentNode(node)
}

// Equals returns true whether the specified node is same, otherwise false.
func (node *remoteNode) Equals(otherNode Node) bool {
	return nodeEquals(node, otherNode)
}

// String returns the node string representation.
func (node *remoteNode) String() string {
	return net.JoinHostPort(node.Address(), strconv.Itoa(node.Port()))
}
