// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"fmt"
	"net"

	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

const (
	errorObjectNotFound              = "Object (%d) not found"
	errorObjectProfileObjectNotFound = "Object profile object not found"
)

// LocalNode is an instance for Echonet node.
type LocalNode struct {
	devices  []*Device
	profiles []*Profile
	Address  net.IP
	server   *server
}

// NewLocalNode returns a new node.
func NewLocalNode() *Node {
	node := &Node{
		devices:  make([]*Device, 0),
		profiles: make([]*Profile, 0),
		server:   newServer(),
	}

	node.AddProfile(NewLocalNodeProfile())
	node.server.SetMessageListener(node)

	return node
}

// AddDevice adds a new device into the node.
func (node *LocalNode) AddDevice(dev *Device) error {
	node.devices = append(node.devices, dev)
	dev.SetParentNode(node)
	return node.updateNodeProfile()
}

// GetDevices returns all device objects.
func (node *LocalNode) GetDevices() []*Device {
	return node.devices
}

// GetDeviceByCode returns a specified device object.
func (node *LocalNode) GetDeviceByCode(code uint) (*Device, error) {
	for _, obj := range node.devices {
		if obj.GetCode() == code {
			return obj, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// AddProfile adds a new profile object into the node.
func (node *LocalNode) AddProfile(prof *Profile) error {
	node.profiles = append(node.profiles, prof)
	prof.SetParentNode(node)
	return node.updateNodeProfile()
}

// GetProfiles returns all profile objects.
func (node *LocalNode) GetProfiles() []*Profile {
	return node.profiles
}

// GetProfileByCode returns a specified profile object.
func (node *LocalNode) GetProfileByCode(code uint) (*Profile, error) {
	for _, prof := range node.profiles {
		if prof.GetCode() == code {
			return prof, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// GetNodeProfile returns a node profile in the node.
func (node *LocalNode) GetNodeProfile() (*Profile, error) {
	return node.GetProfileByCode(NodeProfileObject)
}

// GetObjectByCode returns a specified object.
func (node *LocalNode) GetObjectByCode(code uint) (*Object, error) {
	dev, err := node.GetDeviceByCode(code)
	if err != nil {
		return dev.Object, nil
	}

	prof, err := node.GetProfileByCode(code)
	if err != nil {
		return prof.Object, nil
	}

	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// AnnounceMessage announces a message.
func (node *LocalNode) AnnounceMessage(msg *protocol.Message) error {
	return node.server.NotifyMessage(msg)
}

// SetAddress set an address to the node.
func (node *LocalNode) SetAddress(addr net.IP) {
	node.Address = addr
}

// GetAddress returns a IP address of the node.
func (node *LocalNode) GetAddress() net.IP {
	return node.Address
}

// AnnounceProperty announces a specified property.
func (node *LocalNode) AnnounceProperty(prop *Property) error {
	msg := protocol.NewMessage()
	msg.SetESV(protocol.ESVNotification)
	msg.SetSourceObjectCode(NodeProfileObject)
	msg.SetDestinationObjectCode(NodeProfileObject)
	msg.AddProperty(prop.toProtocolProperty())

	return node.AnnounceMessage(msg)
}

// Announce announces the node
func (node *LocalNode) Announce() error {
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
func (node *LocalNode) SendMessage(dstnode *LocalNode, msg *protocol.Message) error {
	return node.server.SendMessage(string(dstNode.GetAddress()), dstNode.GetPort(), msg)
}

// Start starts the node.
func (node *LocalNode) Start() error {
	err := node.server.Start()
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
