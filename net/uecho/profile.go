// Copyright (C) 2018 Satoshi Konno. All rights reserved.
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
	PropertyMapMaxLen = 16
)

const (
	ProfileFaultStatusLen         = 1
	ProfileManufacturerCodeLen    = ObjectManufacturerCodeLen
	ProfilePlaceOfBusinessCodeLen = 3
	ProfileProductCodeLen         = 12
	ProfileSerialNumberLen        = 12
	ProfileDateOfManufactureLen   = 4
	ProfileAnnoPropertyMapMaxLen  = ObjectAnnoPropertyMapMaxLen
	ProfileSetPropertyMapMaxLen   = ObjectSetPropertyMap
	ProfileGetPropertyMapMaxLen   = ObjectGetPropertyMap
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
	NodeProfileClassOperatingStatusLen                = 1
	NodeProfileClassVersionInformationLen             = 4
	NodeProfileClassIdentificationManufacturerCodeLen = 3
	NodeProfileClassIdentificationUniqueIdLen         = 13
	NodeProfileClassIdentificationNumberLen           = 1 + NodeProfileClassIdentificationManufacturerCodeLen + NodeProfileClassIdentificationUniqueIdLen
	NodeProfileClassFaultContentLen                   = 2
	NodeProfileClassUniqueIdentifierDataLen           = 2
	NodeProfileClassNumberOfSelfNodeInstancesLen      = 3
	NodeProfileClassNumberOfSelfNodeClassesLen        = 2
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
