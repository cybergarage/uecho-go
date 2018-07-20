// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import "github.com/cybergarage/uecho-go/net/uecho/encoding"

const (
	NodeProfileObject                       = 0x0EF001
	NodeProfileObjectReadOnly               = 0x0EF002
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
	NodeProfileClassIdentificationUniqueIdSize         = 13
	NodeProfileClassIdentificationNumberSize           = 1 + NodeProfileClassIdentificationManufacturerCodeSize + NodeProfileClassIdentificationUniqueIdSize
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

	prof.SetCode(NodeProfileObjectReadOnly)

	prof.addNodeProfileMandatoryProperties()

	return prof
}

// addNodeProfileMandatoryProperties sets mandatory properties for node profile
func (prof *Profile) addNodeProfileMandatoryProperties() error {

	// Operation Status
	prof.CreateProperty(NodeProfileClassOperatingStatus, PropertyAttributeReadAnno)
	prof.SetOperatingStatus(true)

	// Version Information
	prof.CreateProperty(NodeProfileClassVersionInformation, PropertyAttributeRead)
	prof.SetVersion(MajorVersion, MinorVersion)

	// Identification Number
	prof.CreateProperty(NodeProfileClassIdentificationNumber, PropertyAttributeRead)
	prof.SetID(ProfileManufacturerUnknown)

	// Number Of Self Node Instances
	prof.CreateProperty(NodeProfileClassNumberOfSelfNodeInstances, PropertyAttributeRead)

	// Number Of Self Node Classes
	prof.CreateProperty(NodeProfileClassNumberOfSelfNodeClasses, PropertyAttributeRead)

	// Instance List Notification
	prof.CreateProperty(NodeProfileClassInstanceListNotification, PropertyAttributeAnno)

	// Self Node Instance ListS
	prof.CreateProperty(NodeProfileClassSelfNodeInstanceListS, PropertyAttributeRead)

	// Self Node Class List S
	prof.CreateProperty(NodeProfileClassSelfNodeClassListS, PropertyAttributeRead)

	return nil
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
	uniqID := make([]byte, NodeProfileClassIdentificationUniqueIdSize)

	return prof.SetPropertyData(NodeProfileClassIdentificationNumber, append(manufactureCodeBytes, uniqID...))
}

// SetInstanceCount sets a instance count in a node.
func (prof *Profile) SetInstanceCount(count uint) error {
	return prof.SetPropertyIntegerData(NodeProfileClassNumberOfSelfNodeInstances, count, NodeProfileClassNumberOfSelfNodeInstancesSize)
}

// SetInstanceList sets a instance list in a node.
func (prof *Profile) SetInstanceList(devices []*Device) error {
	instanceList := make([]byte, 1)
	instanceCount := len(devices)
	if instanceCount <= (NodeProfileClassSelfNodeInstanceListSMax - 1) {
		instanceList[0] = byte(instanceCount)
	} else {
		instanceList[0] = NodeProfileClassSelfNodeInstanceListSMax
	}

	for _, dev := range devices {
		instanceList = append(instanceList, dev.GetCodes()...)
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
	classCount := len(classes)
	if classCount <= (NodeProfileClassSelfNodeClassListSMax - 1) {
		classList[0] = byte(classCount)
	} else {
		classList[0] = NodeProfileClassSelfNodeClassListSMax
	}

	for _, class := range classes {
		classList = append(classList, class.GetCodes()...)
	}

	return prof.SetPropertyData(NodeProfileClassSelfNodeClassListS, classList)
}

/*
  bool _nodeprofileclass_setclasscount(Object *obj, int count)
 {
   return _object_setpropertyintegerdata(obj, NodeProfileClassNumberOfSelfNodeClasses, count, NodeProfileClassNumberOfSelfNodeClassesLen);
 }


 bool _nodeprofileclass_setclasslist(Object *obj, int listCnt, byte *listBytes)
 {
   if (NodeProfileClassSelfNodeClassListSMax < listCnt) {
	 listCnt = NodeProfileClassSelfNodeClassListSMax;
   }
   listBytes[0] = listCnt;
   return _object_setpropertydata(obj, NodeProfileClassSelfNodeClassListS, listBytes, ((listCnt * 2) + 1));
 }

 bool _nodeprofileclass_setinstancelist(Object *obj, int listCnt, byte *listBytes)
 {
   bool isSuccess;

   if (NodeProfileClassSelfNodeInstanceListSMax < listCnt) {
	 listCnt = NodeProfileClassSelfNodeInstanceListSMax;
   }
   listBytes[0] = listCnt;

   isSuccess = true;
   isSuccess &= _object_setpropertydata(obj, NodeProfileClassSelfNodeInstanceListS, listBytes, ((listCnt * 3) + 1));
   isSuccess &= _object_setpropertydata(obj, NodeProfileClassInstanceListNotification, listBytes, ((listCnt * 3) + 1));

   return isSuccess;
 }

 bool _nodeprofileclass_isoperatingstatus(Object *obj)
 {
   byte statsByte;

   if (!_object_getpropertybytedata(obj, NodeProfileClassOperatingStatus, &statsByte))
	 return false;

   return (statsByte == NodeProfileClassBooting) ? true : false;
 }

 int _nodeprofileclass_getinstancecount(Object *obj)
 {
   int count;

   if (!_object_getpropertyintegerdata(obj, NodeProfileClassNumberOfSelfNodeInstances, NodeProfileClassNumberOfSelfNodeInstancesLen, &count))
	 return 0;

   return count;
 }

 int _nodeprofileclass_getclasscount(Object *obj)
 {
   int count;

   if (!_object_getpropertyintegerdata(obj, NodeProfileClassNumberOfSelfNodeClasses, NodeProfileClassNumberOfSelfNodeClassesLen, &count))
	 return 0;

   return count;
 }

 byte *_nodeprofileclass_getnotificationinstancelist(Object *obj)
 {
   return _object_getpropertydata(obj, NodeProfileClassInstanceListNotification);
 }

 byte *_nodeprofileclass_getinstancelist(Object *obj)
 {
   return _object_getpropertydata(obj, NodeProfileClassSelfNodeInstanceListS);
 }

 byte *_nodeprofileclass_getclasslist(Object *obj)
 {
   return _object_getpropertydata(obj, NodeProfileClassSelfNodeClassListS);
 }
*/
