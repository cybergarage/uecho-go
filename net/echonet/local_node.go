// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"net"
	"strconv"
	"sync"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// LocalNode is an instance for Echonet node.
type LocalNode struct {
	*baseNode
	*server
	*sync.Mutex
	*Config
	manufacturerCode uint
	lastTID          uint
	postResponseCh   chan *protocol.Message
	postRequestMsg   *protocol.Message
	listener         NodeListener
}

// NewLocalNode returns a new node.
func NewLocalNode() *LocalNode {
	node := &LocalNode{
		baseNode:         newBaseNode(),
		server:           newServer(),
		Mutex:            new(sync.Mutex),
		manufacturerCode: NodeManufacturerUnknown,
		Config:           NewDefaultConfig(),
		lastTID:          TIDMin,
		postResponseCh:   nil,
		postRequestMsg:   nil,
		listener:         nil,
	}

	node.AddProfile(NewLocalNodeProfile())
	node.server.SetMessageHandler(node)

	return node
}

// SetConfig sets all configuration flags.
func (node *LocalNode) SetConfig(newConfig *Config) {
	node.Config = newConfig
	node.server.MessageManager.SetConfig(newConfig.TransportConfig)
}

// SetManufacturerCode sets a manufacture codes to the node.
func (node *LocalNode) SetManufacturerCode(code uint) {
	node.manufacturerCode = code

	for _, dev := range node.devices {
		dev.SetManufacturerCode(code)
	}

	for _, prop := range node.profiles {
		prop.SetManufacturerCode(code)
	}
}

// GetManufacturerCode return the manufacture codes of the node.
func (node *LocalNode) GetManufacturerCode() uint {
	return node.manufacturerCode
}

// SetListener set a listener to the node.
func (node *LocalNode) SetListener(l NodeListener) {
	node.listener = l
}

// GetListener returns the listener of the node.
func (node *LocalNode) GetListener() NodeListener {
	return node.listener
}

// GetLastTID returns a last sent TID.
func (node *LocalNode) GetLastTID() uint {
	return node.lastTID
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
	node.baseNode.AddDevice(dev)
	dev.SetManufacturerCode(node.manufacturerCode)
	dev.SetParentNode(node)
	return node.updateNodeProfile()
}

// AddProfile adds a new profile object into the node, and set the node profile and manufacture code.
func (node *LocalNode) AddProfile(prof *Profile) error {
	node.baseNode.AddProfile(prof)
	prof.SetManufacturerCode(node.manufacturerCode)
	prof.SetParentNode(node)
	return node.updateNodeProfile()
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
	if err := node.server.Stop(); err != nil {
		return err
	}
	return nil
}

// Restart starts the node.
func (node *LocalNode) Restart() error {
	if err := node.Stop(); err != nil {
		return err
	}
	return node.Start()
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
	return net.JoinHostPort(node.GetAddress(), strconv.Itoa(node.GetPort()))
}
