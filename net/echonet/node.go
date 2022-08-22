// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

const (
	NodeManufacturerUnknown = ObjectManufacturerUnknown
)

// Node is an interface for Echonet node.
type Node interface {
	// GetObjects returns all objects.
	GetObjects() []*Object
	// GetObject returns the specified object.
	GetObject(code ObjectCode) (*Object, error)

	// AddDevice adds a new device into the node.
	AddDevice(dev *Device)
	// Devices returns all device objects.
	Devices() []*Device
	// GetDevice returns the specified device object.
	GetDevice(code ObjectCode) (*Device, error)

	// AddProfile adds a new profile object into the node.
	AddProfile(prof *Profile)
	// GetProfiles returns all profile objects.
	GetProfiles() []*Profile
	// GetProfile returns the specified profile object.
	GetProfile(code ObjectCode) (*Profile, error)

	// GetAddress returns the bound address.
	GetAddress() string
	// GetPort returns the bound address.
	GetPort() int

	// Equals returns true whether the specified node is same, otherwise false.
	Equals(Node) bool
}

// nodeEquals returns true whether the specified node is same, otherwise false.
func nodeEquals(node1, node2 Node) bool {
	if node1.GetPort() != node2.GetPort() {
		return false
	}
	if node1.GetAddress() != node2.GetAddress() {
		return false
	}
	return true
}
