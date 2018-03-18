// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

const ()

// Node is an instance for Echonet node.
type Node struct {
	Classes []Class
	Objects []Object
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

// Start starts the node.
func (node *Node) Start() error {
	err := node.server.Start()
	if err != nil {
		return nil
	}

	return nil
}

// Stop stop the node.
func (node *Node) Stop() error {
	err := node.server.Stop()
	if err != nil {
		return nil
	}

	return nil
}
