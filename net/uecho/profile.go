// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

const (
	NodeProfileObject                       = 0x0EF001
	NodeProfileObjectReadOnly               = 0x0EF002
	NodeProfileClassGroupCode               = 0x0E
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

// Profile is an instance for Echonet profile object.
type Profile = Object

// NewProfile returns a new profile object.
func NewProfile() *Profile {
	prop := NewObject()

	prop.SetClassGroupCode(NodeProfileClassGroupCode)
	prop.SetClassCode(NodeProfileClassCode)
	prop.SetInstanceCode(NodeProfileInstanceGeneralCode)

	prop.addProfileMandatoryProperties()

	return prop
}

// addProfileMandatoryProperties sets mandatory properties for node profile
func (prop *Profile) addProfileMandatoryProperties() error {
	// Operation Status
	prop.CreateProperty(NodeProfileClassOperatingStatus, PropertyAttributeReadAnno)
	prop.SetOperatingStatus(true)

	/*
	   // Version Information

	   uecho_object_setproperty(obj, uEchoNodeProfileClassVersionInformation, uEchoPropertyAttrRead);
	   uecho_nodeprofileclass_setversion(obj, uEchoMajorVersion, uEchoMinorVersion, uEchoSpecifiedMessageFormat);

	   // Identification Number

	   uecho_object_setproperty(obj, uEchoNodeProfileClassIdentificationNumber, uEchoPropertyAttrRead);
	   uecho_nodeprofileclass_setdefaultid(obj);

	   // Number Of Self Node Instances

	   uecho_object_setproperty(obj, uEchoNodeProfileClassNumberOfSelfNodeInstances, uEchoPropertyAttrRead);

	   // Number Of Self Node Classes

	   uecho_object_setproperty(obj, uEchoNodeProfileClassNumberOfSelfNodeClasses, uEchoPropertyAttrRead);

	   // Instance List Notification

	   uecho_object_setproperty(obj, uEchoNodeProfileClassInstanceListNotification, uEchoPropertyAttrAnno);

	   // Self Node Instance ListS

	   uecho_object_setproperty(obj, uEchoNodeProfileClassSelfNodeInstanceListS, uEchoPropertyAttrRead);

	   // Self Node Class List S

	   uecho_object_setproperty(obj, uEchoNodeProfileClassSelfNodeClassListS, uEchoPropertyAttrRead);
	*/

	return nil
}

/*
 bool uecho_nodeprofileclass_setversion(uEchoObject *obj, int majorVer, int minorVer, uEchoMessageFormatType msgType)
 {
   byte verBytes[uEchoNodeProfileClassVersionInformationLen];

   verBytes[0] = uEchoMajorVersion;
   verBytes[1] = uEchoMinorVersion;
   verBytes[2] = uEchoSpecifiedMessageFormat;
   verBytes[3] = 0x00;

   return uecho_object_setpropertydata(obj, uEchoNodeProfileClassVersionInformation, verBytes, uEchoNodeProfileClassVersionInformationLen);
 }

 bool uecho_nodeprofileclass_setid(uEchoObject *obj, byte *manCode, byte *uniqId)
 {
   byte propData[uEchoNodeProfileClassIdentificationNumberLen];
   byte *prop;

   propData[0] = uEchoLowerCommunicationLayerProtocolType;

   prop = propData + 1;
   memcpy(prop, manCode, uEchoNodeProfileClassIdentificationManufacturerCodeLen);

   prop += uEchoNodeProfileClassIdentificationManufacturerCodeLen;
   memcpy(prop, uniqId, uEchoNodeProfileClassIdentificationUniqueIdLen);

   return uecho_object_setpropertydata(obj, uEchoNodeProfileClassIdentificationNumber, propData, uEchoNodeProfileClassIdentificationNumberLen);
 }

 bool uecho_nodeprofileclass_setdefaultid(uEchoObject *obj)
 {
   byte manCode[uEchoNodeProfileClassIdentificationManufacturerCodeLen];
   byte uniqId[uEchoNodeProfileClassIdentificationUniqueIdLen];

   memset(manCode, 0, sizeof(manCode));
   memset(uniqId, 0, sizeof(uniqId));

   return uecho_nodeprofileclass_setid(obj, manCode, uniqId);
 }

  bool uecho_nodeprofileclass_setclasscount(uEchoObject *obj, int count)
 {
   return uecho_object_setpropertyintegerdata(obj, uEchoNodeProfileClassNumberOfSelfNodeClasses, count, uEchoNodeProfileClassNumberOfSelfNodeClassesLen);
 }

 bool uecho_nodeprofileclass_setinstancecount(uEchoObject *obj, int count)
 {
   return uecho_object_setpropertyintegerdata(obj, uEchoNodeProfileClassNumberOfSelfNodeInstances, count, uEchoNodeProfileClassNumberOfSelfNodeInstancesLen);
 }

 bool uecho_nodeprofileclass_setclasslist(uEchoObject *obj, int listCnt, byte *listBytes)
 {
   if (uEchoNodeProfileClassSelfNodeClassListSMax < listCnt) {
	 listCnt = uEchoNodeProfileClassSelfNodeClassListSMax;
   }
   listBytes[0] = listCnt;
   return uecho_object_setpropertydata(obj, uEchoNodeProfileClassSelfNodeClassListS, listBytes, ((listCnt * 2) + 1));
 }

 bool uecho_nodeprofileclass_setinstancelist(uEchoObject *obj, int listCnt, byte *listBytes)
 {
   bool isSuccess;

   if (uEchoNodeProfileClassSelfNodeInstanceListSMax < listCnt) {
	 listCnt = uEchoNodeProfileClassSelfNodeInstanceListSMax;
   }
   listBytes[0] = listCnt;

   isSuccess = true;
   isSuccess &= uecho_object_setpropertydata(obj, uEchoNodeProfileClassSelfNodeInstanceListS, listBytes, ((listCnt * 3) + 1));
   isSuccess &= uecho_object_setpropertydata(obj, uEchoNodeProfileClassInstanceListNotification, listBytes, ((listCnt * 3) + 1));

   return isSuccess;
 }

 bool uecho_nodeprofileclass_isoperatingstatus(uEchoObject *obj)
 {
   byte statsByte;

   if (!uecho_object_getpropertybytedata(obj, uEchoNodeProfileClassOperatingStatus, &statsByte))
	 return false;

   return (statsByte == uEchoNodeProfileClassBooting) ? true : false;
 }

 int uecho_nodeprofileclass_getinstancecount(uEchoObject *obj)
 {
   int count;

   if (!uecho_object_getpropertyintegerdata(obj, uEchoNodeProfileClassNumberOfSelfNodeInstances, uEchoNodeProfileClassNumberOfSelfNodeInstancesLen, &count))
	 return 0;

   return count;
 }

 int uecho_nodeprofileclass_getclasscount(uEchoObject *obj)
 {
   int count;

   if (!uecho_object_getpropertyintegerdata(obj, uEchoNodeProfileClassNumberOfSelfNodeClasses, uEchoNodeProfileClassNumberOfSelfNodeClassesLen, &count))
	 return 0;

   return count;
 }

 byte *uecho_nodeprofileclass_getnotificationinstancelist(uEchoObject *obj)
 {
   return uecho_object_getpropertydata(obj, uEchoNodeProfileClassInstanceListNotification);
 }

 byte *uecho_nodeprofileclass_getinstancelist(uEchoObject *obj)
 {
   return uecho_object_getpropertydata(obj, uEchoNodeProfileClassSelfNodeInstanceListS);
 }

 byte *uecho_nodeprofileclass_getclasslist(uEchoObject *obj)
 {
   return uecho_object_getpropertydata(obj, uEchoNodeProfileClassSelfNodeClassListS);
 }

 bool uecho_nodeprofileclass_updateinstanceproperties(uEchoObject *obj)
 {
   uEchoNode *node;
   uEchoClass *nodeCls;
   uEchoObject *nodeObj;
   byte *nodeClassList, *nodeInstanceList;
   int nodeClassListCnt, nodeInstanceListCnt;
   int nodeClassCnt, nodeInstanceCnt;
   int idx;

   if(!obj)
	 return false;

   node = uecho_object_getparentnode(obj);
   if (!node)
	 return false;

   // Class Properties

   nodeClassList = (byte *)realloc(NULL, 1);
   nodeClassListCnt = 0;
   nodeClassCnt = 0;

   for (nodeCls = uecho_node_getclasses(node); nodeCls; nodeCls = uecho_class_next(nodeCls)) {
	 nodeClassCnt++;

	 if (uecho_class_isprofile(nodeCls))
	   continue;

	 nodeClassListCnt++;
	 nodeClassList = (byte *)realloc(nodeClassList, ((2 * nodeClassListCnt) + 1));
	 idx = (2 * (nodeClassListCnt - 1)) + 1;
	 nodeClassList[idx + 0] = uecho_class_getclassgroupcode(nodeCls);
	 nodeClassList[idx + 1] = uecho_class_getclasscode(nodeCls);
   }

   uecho_nodeprofileclass_setclasscount(obj, nodeClassCnt);
   uecho_nodeprofileclass_setclasslist(obj, nodeClassListCnt, nodeClassList);

   free(nodeClassList);

   // Instance Properties

   nodeInstanceList = (byte *)realloc(NULL, 1);
   nodeInstanceListCnt = 0;
   nodeInstanceCnt = 0;

   for (nodeObj = uecho_node_getobjects(node); nodeObj; nodeObj = uecho_object_next(nodeObj)) {
	 if (uecho_object_isprofile(nodeObj))
	   continue;

	 nodeInstanceCnt++;

	 nodeInstanceListCnt++;
	 nodeInstanceList = (byte *)realloc(nodeInstanceList, ((3 * nodeInstanceListCnt) + 1));
	 idx = (3 * (nodeInstanceListCnt - 1)) + 1;
	 nodeInstanceList[idx + 0] = uecho_object_getclassgroupcode(nodeObj);
	 nodeInstanceList[idx + 1] = uecho_object_getclasscode(nodeObj);
	 nodeInstanceList[idx + 2] = uecho_object_getinstancecode(nodeObj);
   }

   uecho_nodeprofileclass_setinstancecount(obj, nodeInstanceCnt);
   uecho_nodeprofileclass_setinstancelist(obj, nodeInstanceListCnt, nodeInstanceList);

   free(nodeInstanceList);

   return true;
 }

*/
