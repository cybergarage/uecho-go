// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

// NewRemoteNode returns a new node.
func NewRemoteNode() *Node {
	node := &Node{
		devices:  make([]*Device, 0),
		profiles: make([]*Profile, 0),
		server:   nil,
	}

	return node
}

// NewRemoteNodeWithRequestMessage returns a new node with the specified request message.
func NewRemoteNodeWithRequestMessage(msg *protocol.Message) *Node {
	node := NewRemoteNode()
	node.SetAddress(msg.From.IP)
	return node
}
