// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
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
	// LookupObject returns the specified object.
	LookupObject(code ObjectCode) (*Object, error)

	// Devices returns all device objects.
	Devices() []*Device
	// LookupDevice returns the specified device object.
	LookupDevice(code ObjectCode) (*Device, error)

	// Profiles returns all profile objects.
	Profiles() []*Profile
	// LookupProfile returns the specified profile object.
	LookupProfile(code ObjectCode) (*Profile, error)
	// NodeProfile returns the node profile object.
	NodeProfile() (*Profile, error)

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
