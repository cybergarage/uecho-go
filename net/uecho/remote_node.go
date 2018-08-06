// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"github.com/cybergarage/uecho-go/net/uecho/protocol"
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
	node.SetAddress(msg.From.String())
	return node
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
