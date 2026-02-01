// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
)

const (
	errObjectNotFound              = "%w: object (%X)"
	errObjectProfileObjectNotFound = "%w: object profile"
	errUnknownObjectType           = "%w: unknown object type (%v)"
)

// baseNode is an instance for Echonet node.
type baseNode struct {
	devices  []Device
	profiles []Profile
}

// NewbaseNode returns a new node.
func newBaseNode() *baseNode {
	node := &baseNode{
		devices:  make([]Device, 0),
		profiles: make([]Profile, 0),
	}
	return node
}

// AddDevice adds a new device into the node.
func (node *baseNode) AddDevice(dev Device) {
	node.devices = append(node.devices, dev)
}

// Devices returns all device objects.
func (node *baseNode) Devices() []Device {
	return node.devices
}

// LookupDevice returns a specified device object.
func (node *baseNode) LookupDevice(code ObjectCode) (Device, error) {
	for _, dev := range node.devices {
		objCode := dev.Code()
		if objCode == code {
			return dev, nil
		}
	}
	return nil, fmt.Errorf(errObjectNotFound, ErrNotFound, code)
}

// AddProfile adds a new profile object into the node.
func (node *baseNode) AddProfile(prof Profile) {
	node.profiles = append(node.profiles, prof)
}

// Profiles returns all profile objects.
func (node *baseNode) Profiles() []Profile {
	return node.profiles
}

// LookupProfile returns a specified profile object.
func (node *baseNode) LookupProfile(code ObjectCode) (Profile, error) {
	for _, prof := range node.profiles {
		objCode := prof.Code()
		if objCode == code {
			return prof, nil
		}
	}
	return nil, fmt.Errorf(errObjectNotFound, ErrNotFound, code)
}

// NodeProfile returns a node profile in the node.
func (node *baseNode) NodeProfile() (NodeProfile, error) {
	prof, err := node.LookupProfile(NodeProfileObjectCode)
	if err == nil {
		return NewNodeProfileWith(prof), nil
	}
	prof, err = node.LookupProfile(NodeProfileObjectReadOnlyCode)
	if err == nil {
		return NewNodeProfileWith(prof), nil
	}
	return nil, fmt.Errorf(errObjectProfileObjectNotFound, ErrNotFound)
}

// AddObject adds a new object into the node.
func (node *baseNode) AddObject(obj any) error {
	switch v := obj.(type) {
	case Profile:
		node.AddProfile(v)
		return nil
	case Device:
		node.AddDevice(v)
		return nil
	}
	return fmt.Errorf(errUnknownObjectType, ErrUnknown, obj)
}

// Objects returns all objects.
func (node *baseNode) Objects() []Object {
	devs := node.Devices()
	profs := node.Profiles()
	objs := make([]Object, 0, len(devs)+len(profs))

	for _, dev := range devs {
		objs = append(objs, dev)
	}
	for _, prof := range profs {
		objs = append(objs, prof)
	}

	return objs
}

// LookupObject returns a specified object.
func (node *baseNode) LookupObject(code ObjectCode) (Object, error) {
	dev, err := node.LookupDevice(code)
	if err == nil {
		return dev, nil
	}
	prof, err := node.LookupProfile(code)
	if err == nil {
		return prof, nil
	}
	return nil, fmt.Errorf(errObjectNotFound, ErrNotFound, code)
}
