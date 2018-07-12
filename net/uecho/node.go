// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
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
	Objects []*Object
	Address string
	server  *server
}

// NewNode returns a new object.
func NewNode() *Node {
	node := &Node{
		server: newServer(),
	}
	return node
}

// GetObjectByCode returns a specified object.
func (node *Node) GetObjectByCode(code uint) (*Object, error) {
	for _, obj := range node.Objects {
		if obj.GetCode() == code {
			return obj, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// GetNodeProfileObject returns a specified object.
func (node *Node) GetNodeProfileObject() (*Object, error) {
	return node.GetObjectByCode(NodeProfileObject)
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

// AnnounceMessage announces a message.
func (node *Node) AnnounceMessage(msg *protocol.Message) error {
	return node.server.SendMulticastMessage(msg)
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

	nodeProp, err := nodePropObj.GetProperty(NodeProfileClassInstanceListNotification)
	if err != nil {
		return err
	}

	return node.AnnounceProperty(nodeProp)
}
