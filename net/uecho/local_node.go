// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"fmt"
	"sync"
	"time"

	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

const (
	DefaultLocalNodeRequestTimeout = time.Second * 30
)

// LocalNode is an instance for Echonet node.
type LocalNode struct {
	*baseNode
	*server
	*sync.Mutex
	manufacturerCode uint
	lastTID          uint
	postResponseCh   chan *protocol.Message
	postRequestMsg   *protocol.Message
	requestTimeout   time.Duration
}

// NewLocalNode returns a new node.
func NewLocalNode() *LocalNode {
	node := &LocalNode{
		baseNode:         newBaseNode(),
		server:           newServer(),
		Mutex:            new(sync.Mutex),
		manufacturerCode: NodeManufacturerUnknown,
		lastTID:          TIDMin,
		postResponseCh:   nil,
		postRequestMsg:   nil,
		requestTimeout:   DefaultLocalNodeRequestTimeout,
	}

	node.AddProfile(NewLocalNodeProfile())
	node.server.SetMessageListener(node)

	return node
}

// SetManufacturerCode sets a manufacture codes to the node.
func (node *LocalNode) SetManufacturerCode(code uint) {
	node.manufacturerCode = code
}

// GetManufacturerCode return the manufacture codes of the node.
func (node *LocalNode) GetManufacturerCode() uint {
	return node.manufacturerCode
}

// SetRequestTimeout sets a request timeoutto the node.
func (node *LocalNode) SetRequestTimeout(d time.Duration) {
	node.requestTimeout = d
}

// GetRequestTimeout return the request timeout of the node.
func (node *LocalNode) GetRequestTimeout() time.Duration {
	return node.requestTimeout
}

// getNextTID returns a next TID.
func (node *LocalNode) getNextTID() uint {
	if TIDMax <= node.lastTID {
		node.lastTID = TIDMin
	} else {
		node.lastTID++
	}
	return node.lastTID
}

// AddDevice adds a new device into the node, and set the node profile and manufacture code.
func (node *LocalNode) AddDevice(dev *Device) error {
	err := node.baseNode.AddDevice(dev)
	if err != nil {
		return err
	}

	dev.SetManufacturerCode(node.manufacturerCode)
	dev.SetParentNode(node)

	return node.updateNodeProfile()
}

// AddProfile adds a new profile object into the node, and set the node profile and manufacture code.
func (node *LocalNode) AddProfile(prof *Profile) error {
	err := node.baseNode.AddProfile(prof)
	if err != nil {
		return err
	}

	prof.SetManufacturerCode(node.manufacturerCode)
	prof.SetParentNode(node)

	return node.updateNodeProfile()
}

// GetAddress returns the bound address.
func (node *LocalNode) GetAddress() string {
	addrs := node.server.GetBoundAddresses()
	if len(addrs) <= 0 {
		return ""
	}
	return addrs[0]
}

// GetPort returns the bound address.
func (node *LocalNode) GetPort() int {
	return node.server.GetPort()
}

// Start starts the node.
func (node *LocalNode) Start() error {
	err := node.server.Start()
	if err != nil {
		return err
	}

	err = node.Announce()
	if err != nil {
		return err
	}

	return nil
}

// Stop stop the node.
func (node *LocalNode) Stop() error {
	err := node.server.Stop()
	if err != nil {
		return err
	}

	return nil
}

// Equals returns true whether the specified node is same, otherwise false.
func (node *LocalNode) Equals(otherNode Node) bool {
	return nodeEquals(node, otherNode)
}

// updateNodeProfile updates the node profile in the node.
func (node *LocalNode) updateNodeProfile() error {
	nodeProf, err := node.GetNodeProfile()
	if err != nil {
		return err
	}

	// Check the current all objects

	classes := make([]*Class, 0)

	for _, dev := range node.devices {
		devClass := dev.GetClass()
		hasSameClass := false
		for _, class := range classes {
			if class.Equals(devClass) {
				hasSameClass = true
				break
			}
		}
		if hasSameClass {
			continue
		}
		classes = append(classes, devClass)
	}

	for _, prof := range node.profiles {
		profClass := prof.GetClass()
		hasSameClass := false
		for _, class := range classes {
			if class.Equals(profClass) {
				hasSameClass = true
				break
			}
		}
		if hasSameClass {
			continue
		}
		classes = append(classes, profClass)
	}

	// Number of self-node instances

	instanceCount := uint(len(node.devices))
	nodeProf.SetInstanceCount(instanceCount)

	// Number of self-node classes

	nodeProf.SetClassCount(uint(len(classes)))

	// Self-node instance list S and Instance list notification

	nodeProf.SetInstanceList(node.devices)

	// Self-node class list S

	nodeProf.SetClassList(classes)

	return nil
}

// String returns the node string representation.
func (node *LocalNode) String() string {
	return fmt.Sprintf("%s:%d", node.GetAddress(), node.GetPort())
}
