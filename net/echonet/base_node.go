// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
)

const (
	errorObjectNotFound              = "Object (%X) not found"
	errorObjectProfileObjectNotFound = "Object profile object not found"
)

// baseNode is an instance for Echonet node.
type baseNode struct {
	devices  []*Device
	profiles []*Profile
}

// NewbaseNode returns a new node.
func newBaseNode() *baseNode {
	node := &baseNode{
		devices:  make([]*Device, 0),
		profiles: make([]*Profile, 0),
	}
	return node
}

// AddDevice adds a new device into the node.
func (node *baseNode) AddDevice(dev *Device) error {
	node.devices = append(node.devices, dev)
	return nil
}

// GetDevices returns all device objects.
func (node *baseNode) GetDevices() []*Device {
	return node.devices
}

// GetDevice returns a specified device object.
func (node *baseNode) GetDevice(code ObjectCode) (*Device, error) {
	for _, dev := range node.devices {
		objCode := dev.GetCode()
		if objCode == code {
			return dev, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// AddProfile adds a new profile object into the node.
func (node *baseNode) AddProfile(prof *Profile) error {
	node.profiles = append(node.profiles, prof)
	return nil
}

// GetProfiles returns all profile objects.
func (node *baseNode) GetProfiles() []*Profile {
	return node.profiles
}

// GetProfile returns a specified profile object.
func (node *baseNode) GetProfile(code ObjectCode) (*Profile, error) {
	for _, prof := range node.profiles {
		objCode := prof.GetCode()
		if objCode == code {
			return prof, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// GetNodeProfile returns a node profile in the node.
func (node *baseNode) GetNodeProfile() (*Profile, error) {
	prof, err := node.GetProfile(NodeProfileObject)
	if err == nil {
		return prof, nil
	}
	return node.GetProfile(NodeProfileObjectReadOnly)
}

// GetObjects returns all objects.
func (node *baseNode) GetObjects() []*Object {
	objs := make([]*Object, 0)

	devs := node.GetDevices()
	for _, dev := range devs {
		objs = append(objs, dev.Object)
	}

	profs := node.GetProfiles()
	for _, prof := range profs {
		objs = append(objs, prof.Object)
	}

	return objs
}

// GetObject returns a specified object.
func (node *baseNode) GetObject(code ObjectCode) (*Object, error) {
	dev, err := node.GetDevice(code)
	if err == nil {
		return dev.Object, nil
	}

	prof, err := node.GetProfile(code)
	if err == nil {
		return prof.Object, nil
	}

	return nil, fmt.Errorf(errorObjectNotFound, code)
}
