// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

const ()

// Node is an instance for Echonet node.
type Node struct {
	Classes []Class
	Objects []Object
	Address string
}

// NewNode returns a new object.
func NewNode() *Node {
	node := &Node{}
	return node
}
