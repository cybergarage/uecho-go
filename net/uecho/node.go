// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

const (
	NodeManufacturerUnknown = ObjectManufacturerUnknown
)

// Node is an interface for Echonet node.
type Node interface {
	// GetObject returns the specified object.
	GetObject(code uint) (*Object, error)

	// AddDevice adds a new device into the node.
	AddDevice(dev *Device) error
	// GetDevices returns all device objects.
	GetDevices() []*Device
	// GetDevice returns the specified device object.
	GetDevice(code uint) (*Device, error)

	// AddProfile adds a new profile object into the node.
	AddProfile(prof *Profile) error
	// GetProfiles returns all profile objects.
	GetProfiles() []*Profile
	// GetProfile returns the specified profile object.
	GetProfile(code uint) (*Profile, error)

	// GetAddress returns the bound address.
	GetAddress() string
	// GetPort returns the bound address.
	GetPort() int
}
