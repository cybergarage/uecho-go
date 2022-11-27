// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

const (
	NodeManufacturerExperimental = ObjectManufacturerExperimental
)

// Node is an interface for Echonet node.
type Node interface {
	// Objects returns all objects.
	Objects() []*Object
	// FindObject returns the specified object.
	FindObject(code ObjectCode) (*Object, error)

	// AddDevice adds a new device into the node.
	AddDevice(dev *Device)
	// Devices returns all device objects.
	Devices() []*Device
	// FindDevice returns the specified device object.
	FindDevice(code ObjectCode) (*Device, error)

	// AddProfile adds a new profile object into the node.
	AddProfile(prof *Profile)
	// Profiles returns all profile objects.
	Profiles() []*Profile
	// FindProfile returns the specified profile object.
	FindProfile(code ObjectCode) (*Profile, error)

	// Address returns the bound address.
	Address() string
	// GetPort returns the bound address.
	Port() int

	// Equals returns true whether the specified node is same, otherwise false.
	Equals(Node) bool
}

// nodeEquals returns true whether the specified node is same, otherwise false.
func nodeEquals(node1, node2 Node) bool {
	if node1.Port() != node2.Port() {
		return false
	}
	if node1.Address() != node2.Address() {
		return false
	}
	return true
}
