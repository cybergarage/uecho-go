// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

const (
	ProfileClassGroupCode = 0x0E
)

const (
	ProfileFaultStatus         = 0x88
	ProfileManufacturerCode    = ObjectManufacturerCode
	ProfilePlaceOfBusinessCode = 0x8B
	ProfileProductCode         = 0x8C
	ProfileSerialNumber        = 0x8D
	ProfileDateOfManufacture   = 0x8E
	ProfileAnnoPropertyMap     = ObjectAnnoPropertyMap
	ProfileSetPropertyMap      = ObjectSetPropertyMap
	ProfileGetPropertyMap      = ObjectGetPropertyMap
)

const (
	ProfileFaultStatusLen         = 1
	ProfileManufacturerCodeLen    = ObjectManufacturerCodeSize
	ProfilePlaceOfBusinessCodeLen = 3
	ProfileProductCodeLen         = 12
	ProfileSerialNumberLen        = 12
	ProfileDateOfManufactureLen   = 4
)

const (
	ProfileFaultEncountered    = 0x41
	ProfileNoFaultEncountered  = 0x42
	ProfileManufacturerUnknown = ObjectManufacturerExperimental
)

// Profile represents an instance for a profile object of Echonet.
type Profile interface {
	SuperObject
	// IsNodeProfile indicates whether the profile is a node profile or not.
	IsNodeProfile() bool
}

type profile struct {
	SuperObject
}

// isProfileObjectCode returns true when the class group code is the profile code, otherwise false.
func isProfileObjectCode(code byte) bool {
	return (code == ProfileClassGroupCode)
}

// NewProfile returns a new profile object.
func NewProfile() Profile {
	prof := &profile{
		SuperObject: NewSuperObject(),
	}
	prof.SetClassGroupCode(ProfileClassGroupCode)
	return prof
}

// IsNodeProfile indicates whether the profile is a node profile or not.
func (prof *profile) IsNodeProfile() bool {
	return isNodeProfileObjectCode(prof.Code())
}
