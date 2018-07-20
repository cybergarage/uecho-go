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
	devices  []*Device
	profiles []*Profile
	Address  string
	server   *server
}

// NewNode returns a new object.
func NewNode() *Node {
	node := &Node{
		devices: make([]*Device, 0),
		server:  newServer(),
	}

	node.AddProfile(NewNodeProfile())

	return node
}

// AddDevice adds a new device into the node.
func (node *Node) AddDevice(dev *Device) error {
	node.devices = append(node.devices, dev)
	return node.updateNodeProfile()
}

// GetDevices returns all device objects.
func (node *Node) GetDevices() []*Device {
	return node.devices
}

// GetDeviceByCode returns a specified device object.
func (node *Node) GetDeviceByCode(code uint) (*Device, error) {
	for _, obj := range node.devices {
		if obj.GetCode() == code {
			return obj, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// AddProfile adds a new profile object into the node.
func (node *Node) AddProfile(prof *Profile) error {
	node.profiles = append(node.profiles, prof)
	return node.updateNodeProfile()
}

// GetProfiles returns all profile objects.
func (node *Node) GetProfiles() []*Profile {
	return node.profiles
}

// GetProfilesByCode returns a specified profile object.
func (node *Node) GetProfilesByCode(code uint) (*Profile, error) {
	for _, prof := range node.profiles {
		if prof.GetCode() == code {
			return prof, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// GetNodeProfile returns a node profile in the node.
func (node *Node) GetNodeProfile() (*Profile, error) {
	return node.GetProfilesByCode(NodeProfileObject)
}

// AnnounceMessage announces a message.
func (node *Node) AnnounceMessage(msg *protocol.Message) error {
	return node.server.SendMessageAll(msg)
}

// GetAddress returns a IP address of the node.
func (node *Node) GetAddress() string {
	return node.Address
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
	nodePropObj, err := node.GetNodeProfile()
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

// updateNodeProfile updates the node profile in the node.
func (node *Node) updateNodeProfile() error {
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
