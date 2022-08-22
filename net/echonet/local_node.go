// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"net"
	"strconv"
	"sync"

	"github.com/cybergarage/go-logger/log"
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

// Address returns the bound address.
func (node *LocalNode) Address() string {
	for _, server := range node.GeUnicastManager().Servers {
		addr, err := server.UDPSocket.Address()
		if err == nil {
			return addr
		}
	}
	return ""
}

// GetPort returns the bound address.
func (node *LocalNode) Port() int {
	return node.GeUnicastManager().Port()
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

// Listener returns the listener of the node.
func (node *LocalNode) Listener() NodeListener {
	return node.listener
}

// LastTID returns a last sent TID.
func (node *LocalNode) LastTID() uint {
	return node.lastTID
}

// NextTID returns a next TID.
func (node *LocalNode) NextTID() uint {
	if TIDMax <= node.lastTID {
		node.lastTID = TIDMin
	} else {
		node.lastTID++
	}
	return node.lastTID
}

// AddDevice adds a new device into the node, and set the node profile and manufacture code.
func (node *LocalNode) AddDevice(dev *Device) {
	node.baseNode.AddDevice(dev)
	dev.SetManufacturerCode(node.manufacturerCode)
	dev.SetParentNode(node)
	node.updateNodeProfile()
}

// AddProfile adds a new profile object into the node, and set the node profile and manufacture code.
func (node *LocalNode) AddProfile(prof *Profile) {
	node.baseNode.AddProfile(prof)
	prof.SetManufacturerCode(node.manufacturerCode)
	prof.SetParentNode(node)
	node.updateNodeProfile()
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
func (node *LocalNode) updateNodeProfile() {
	nodeProf, err := node.NodeProfile()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Check the current all objects

	classes := make([]*Class, 0)

	for _, dev := range node.devices {
		devClass := dev.Class()
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
		profClass := prof.Class()
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
}

// String returns the node string representation.
func (node *LocalNode) String() string {
	return net.JoinHostPort(node.Address(), strconv.Itoa(node.Port()))
}
