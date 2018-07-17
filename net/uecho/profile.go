// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

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
	PropertyMapFormat1MaxSize = 15
	PropertyMapFormat2Size    = 18
	PropertyMapFormatMaxSize  = PropertyMapFormat2Size
)

const (
	ProfileFaultStatusSize         = 1
	ProfileManufacturerCodeSize    = ObjectManufacturerCodeSize
	ProfilePlaceOfBusinessCodeSize = 3
	ProfileProductCodeSize         = 12
	ProfileSerialNumberSize        = 12
	ProfileDateOfManufactureSize   = 4
	ProfileAnnoPropertyMapMaxSize  = ObjectAnnoPropertyMapMaxSize
	ProfileSetPropertyMapMaxSize   = ObjectSetPropertyMap
	ProfileGetPropertyMapMaxSize   = ObjectGetPropertyMap
)

const (
	ProfileFaultEncountered   = 0x41
	ProfileNoFaultEncountered = 0x42
)

const (
	NodeProfileClassOperatingStatus           = 0x80
	NodeProfileClassVersionInformation        = 0x82
	NodeProfileClassIdentificationNumber      = 0x83
	NodeProfileClassFaultContent              = 0x89
	NodeProfileClassUniqueIdentifierData      = 0xBF
	NodeProfileClassNumberOfSelfNodeInstances = 0xD3
	NodeProfileClassNumberOfSelfNodeClasses   = 0xD4
	NodeProfileClassInstanceListNotification  = 0xD5
	NodeProfileClassSelfNodeInstanceListS     = 0xD6
	NodeProfileClassSelfNodeClassListS        = 0xD7
)

const (
	NodeProfileClassOperatingStatusSize                = 1
	NodeProfileClassVersionInformationSize             = 4
	NodeProfileClassIdentificationManufacturerCodeSize = 3
	NodeProfileClassIdentificationUniqueIdSize         = 13
	NodeProfileClassIdentificationNumberSize           = 1 + NodeProfileClassIdentificationManufacturerCodeSize + NodeProfileClassIdentificationUniqueIdSize
	NodeProfileClassFaultContentSize                   = 2
	NodeProfileClassUniqueIdentifierDataSize           = 2
	NodeProfileClassNumberOfSelfNodeInstancesSize      = 3
	NodeProfileClassNumberOfSelfNodeClassesSize        = 2
	NodeProfileClassSelfNodeInstanceListSMax          = 0xFF
	NodeProfileClassSelfNodeClassListSMax             = 0xFF
	NodeProfileClassInstanceListNotificationMax       = NodeProfileClassSelfNodeInstanceListSMax
)

const (
	NodeProfileClassBooting             = 0x30
	NodeProfileClassNotBooting          = 0x31
	LowerCommunicationLayerProtocolType = 0xFE
)

// Profile is an instance for Echonet profile object.
type Profile = Object

// NewProfile returns a new profile object.
func NewProfile() *Profile {
	prop := NewObject()
	return prop
}
