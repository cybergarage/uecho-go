// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

const (
	errorInvalidNotificationMessage = "Invalid Notification Message : %s"
)

// RemoteNode is an instance for Echonet node.
type RemoteNode struct {
	*baseNode
	Address string
	Port    int
}

// NewRemoteNode returns a new remote node.
func NewRemoteNode() *RemoteNode {
	node := &RemoteNode{
		baseNode: newBaseNode(),
		Address:  "",
		Port:     0,
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
	msgOPC := msg.GetOPC()
	if msgOPC < 1 {
		return nil, fmt.Errorf(errorInvalidNotificationMessage, msg)
	}

	prop := msg.GetProperty(0)
	if prop == nil {
		return nil, fmt.Errorf(errorInvalidNotificationMessage, msg)
	}

	if prop.GetCode() != NodeProfileClassInstanceListNotification && prop.GetCode() != NodeProfileClassSelfNodeInstanceListS {
		return nil, fmt.Errorf(errorInvalidNotificationMessage, msg)
	}

	propData := prop.GetData()
	propSize := len(propData)
	if propSize <= 0 {
		return nil, fmt.Errorf(errorInvalidNotificationMessage, msg)
	}

	instanceCount := int(propData[0])
	if propSize < ((instanceCount * ObjectCodeSize) + 1) {
		return nil, fmt.Errorf(errorInvalidNotificationMessage, msg)
	}

	// Create a new remote Node

	node := NewRemoteNode()
	node.SetAddress(msg.GetSourceAddress())
	node.SetPort(msg.GetSourcePort())

	for n := 0; n < instanceCount; n++ {
		objCodes := make([]byte, ObjectCodeSize)
		copy(objCodes, propData[((n*ObjectCodeSize)+1):])
		obj, err := NewObjectWithCodes(objCodes)
		if err != nil {
			return nil, err
		}
		switch obj.(type) {
		case (*Device):
			dev, _ := obj.(*Device)
			node.AddDevice(dev)
		case (*Profile):
			prof, _ := obj.(*Profile)
			node.AddProfile(prof)
		}
	}

	return node, nil
}

// SetAddress set the address to the node.
func (node *RemoteNode) SetAddress(addr string) {
	node.Address = addr
}

// GetAddress returns the address of the node.
func (node *RemoteNode) GetAddress() string {
	return node.Address
}

// SetPort set a port to the node.
func (node *RemoteNode) SetPort(port int) {
	node.Port = port
}

// GetPort returns the port of the node.
func (node *RemoteNode) GetPort() int {
	return node.Port
}

// AddDevice adds a new device into the node, and set the node profile and manufacture code.
func (node *RemoteNode) AddDevice(dev *Device) error {
	err := node.baseNode.AddDevice(dev)
	if err != nil {
		return err
	}
	dev.SetParentNode(node)
	return nil
}

// AddProfile adds a new profile object into the node, and set the node profile and manufacture code.
func (node *RemoteNode) AddProfile(prof *Profile) error {
	err := node.baseNode.AddProfile(prof)
	if err != nil {
		return err
	}
	prof.SetParentNode(node)
	return nil
}

// Equals returns true whether the specified node is same, otherwise false.
func (node *RemoteNode) Equals(otherNode Node) bool {
	return nodeEquals(node, otherNode)
}

// String returns the node string representation.
func (node *RemoteNode) String() string {
	return fmt.Sprintf("%s:%d", node.GetAddress(), node.GetPort())
}
