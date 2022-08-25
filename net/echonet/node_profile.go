// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

const (
	NodeProfileObjectCode                   = ObjectCode(0x0EF001)
	NodeProfileObjectReadOnlyCode           = ObjectCode(0x0EF002)
	NodeProfileClassCode                    = 0xF0
	NodeProfileInstanceGeneralCode          = 0x01
	NodeProfileInstanceTransmissionOnlyCode = 0x02
)

const (
	NodeProfileClassOperatingStatus           = ObjectOperatingStatus
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
	NodeProfileClassIdentificationUniqueIDSize         = 13
	NodeProfileClassIdentificationNumberSize           = 1 + NodeProfileClassIdentificationManufacturerCodeSize + NodeProfileClassIdentificationUniqueIDSize
	NodeProfileClassFaultContentSize                   = 2
	NodeProfileClassUniqueIdentifierDataSize           = 2
	NodeProfileClassNumberOfSelfNodeInstancesSize      = 3
	NodeProfileClassNumberOfSelfNodeClassesSize        = 2
	NodeProfileClassSelfNodeInstanceListSMax           = 0xFF
	NodeProfileClassSelfNodeClassListSMax              = 0xFF
	NodeProfileClassInstanceListNotificationMax        = NodeProfileClassSelfNodeInstanceListSMax
)

const (
	NodeProfileClassOperatingStatusOn   = ObjectOperatingStatusOn
	NodeProfileClassOperatingStatusOff  = ObjectOperatingStatusOff
	NodeProfileClassBooting             = 0x30
	NodeProfileClassNotBooting          = 0x31
	LowerCommunicationLayerProtocolType = 0xFE
)

// NewNodeProfile returns a new node profile object.
func NewNodeProfile() *Profile {
	prof := NewProfile()
	prof.SetCode(NodeProfileObjectCode)
	return prof
}

// SetVersion sets a version to the profile.
func (prof *Profile) SetVersion(major int, minor int) error {
	verBytes := []byte{
		byte(major),
	}
	return prof.SetPropertyData(NodeProfileClassVersionInformation, verBytes)
}

// SetID sets a ID to the profile.
func (prof *Profile) SetID(manufactureCode uint) error {
	manufactureCodeBytes := make([]byte, NodeProfileClassIdentificationManufacturerCodeSize)
	encoding.IntegerToByte(manufactureCode, manufactureCodeBytes)

	// TODO : Set a unique number
	uniqID := make([]byte, NodeProfileClassIdentificationUniqueIDSize)

	return prof.SetPropertyData(NodeProfileClassIdentificationNumber, append(manufactureCodeBytes, uniqID...))
}

// SetInstanceCount sets a instance count in a node.
func (prof *Profile) SetInstanceCount(count uint) error {
	return prof.SetPropertyIntegerData(NodeProfileClassNumberOfSelfNodeInstances, count, NodeProfileClassNumberOfSelfNodeInstancesSize)
}

// SetInstanceList sets a instance list in a node.
func (prof *Profile) SetInstanceList(devices []*Device) error {
	instanceList := make([]byte, 1)
	if instanceCount := len(devices); instanceCount <= (NodeProfileClassSelfNodeInstanceListSMax - 1) {
		instanceList[0] = byte(instanceCount)
	} else {
		instanceList[0] = NodeProfileClassSelfNodeInstanceListSMax
	}

	for _, dev := range devices {
		instanceList = append(instanceList, dev.Codes()...)
	}

	err := prof.SetPropertyData(NodeProfileClassInstanceListNotification, instanceList)
	if err != nil {
		return err
	}

	err = prof.SetPropertyData(NodeProfileClassSelfNodeInstanceListS, instanceList)
	if err != nil {
		return err
	}

	return nil
}

// SetClassCount sets a class count in a node.
func (prof *Profile) SetClassCount(count uint) error {
	return prof.SetPropertyIntegerData(NodeProfileClassNumberOfSelfNodeClasses, count, NodeProfileClassNumberOfSelfNodeClassesSize)
}

// SetClassList sets a class list in a node.
func (prof *Profile) SetClassList(classes []*Class) error {
	classList := make([]byte, 1)
	if classCount := len(classes); classCount <= (NodeProfileClassSelfNodeClassListSMax - 1) {
		classList[0] = byte(classCount)
	} else {
		classList[0] = NodeProfileClassSelfNodeClassListSMax
	}

	for _, class := range classes {
		classList = append(classList, class.Codes()...)
	}
	return prof.SetPropertyData(NodeProfileClassSelfNodeClassListS, classList)
}
