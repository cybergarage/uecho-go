// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

const (
	errorObjectNotFound              = "Object (%d) not found"
	errorObjectProfileObjectNotFound = "Object profile object not found"
)

// Node is an instance for Echonet node.
type Node struct {
	Classes []*Class
	devices []*Device
	Objects []*Object
	Address string
	server  *server
}

// NewNode returns a new object.
func NewNode() *Node {
	node := &Node{
		devices: make([]*Device, 0),
		server:  newServer(),
	}
	return node
}

// GetDevices returns all devices.
func (node *Node) GetDevices() []*Device {
	return node.devices
}

// GetDeviceByCode returns a specified device.
func (node *Node) GetDeviceByCode(code uint) (*Device, error) {
	for _, obj := range node.devices {
		if obj.GetCode() == code {
			return obj, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// GetObjectByCode returns a specified device.
func (node *Node) GetObjectByCode(code uint) (*Object, error) {
	for _, obj := range node.Objects {
		if obj.GetCode() == code {
			return obj, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// GetAddress returns a IP address of the node.
func (node *Node) GetAddress() string {
	return node.Address
}

// GetNodeProfileObject returns a specified object.
func (node *Node) GetNodeProfileObject() (*Object, error) {
	return node.GetObjectByCode(NodeProfileObject)
}

// AnnounceMessage announces a message.
func (node *Node) AnnounceMessage(msg *protocol.Message) error {
	return node.server.SendMessageAll(msg)
}

// AnnounceProperty announces a specified property.
func (node *Node) AnnounceProperty(prop *Property) error {
	msg := protocol.NewMessage()
	msg.SetESV(protocol.ESVNotification)
	msg.SetSourceObjectCode(NodeProfileObject)
	msg.SetDestinationObjectCode(NodeProfileObject)
	msg.AddProperty(prop.toProtocolProperty())

	return node.AnnounceMessage(msg)
}

// Announce announces the node
func (node *Node) Announce() error {
	nodePropObj, err := node.GetNodeProfileObject()
	if err != nil {
		return err
	}

	nodeProp, ok := nodePropObj.GetProperty(NodeProfileClassInstanceListNotification)
	if !ok {
		return fmt.Errorf(errorObjectProfileObjectNotFound)
	}

	return node.AnnounceProperty(nodeProp)
}

// SendMessage send a message to the node
func (node *Node) SendMessage(dstNode *Node, msg *protocol.Message) error {
	return node.server.SendMessage(dstNode.GetAddress(), msg)
}

// Start starts the node.
func (node *Node) Start() error {
	err := node.server.Start()
	if err != nil {
		return err
	}

	return nil
}

// Stop stop the node.
func (node *Node) Stop() error {
	err := node.server.Stop()
	if err != nil {
		return err
	}

	return nil
}
