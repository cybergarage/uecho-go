// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/uecho/encoding"
	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

const (
	errorParentNodeNotFound = "Parent node not found"
)

const (
	ObjectCodeMin     = 0x000000
	ObjectCodeMax     = 0xFFFFFF
	ObjectCodeUnknown = ObjectCodeMin

	NodeProfileObject         = 0x0EF001
	NodeProfileObjectReadOnly = 0x0EF002
)

const (
	ObjectManufacturerCode = 0x8A
	ObjectAnnoPropertyMap  = 0x9D
	ObjectSetPropertyMap   = 0x9E
	ObjectGetPropertyMap   = 0x9F
)

const (
	ObjectManufacturerCodeLen   = 3
	ObjectPropertyMapMaxLen     = 16
	ObjectAnnoPropertyMapMaxLen = (ObjectPropertyMapMaxLen + 1)
	ObjectSetPropertyMapMaxLen  = (ObjectPropertyMapMaxLen + 1)
	ObjectGetPropertyMapMaxLen  = (ObjectPropertyMapMaxLen + 1)
)

const (
	ManufacturerCodeMin    = 0xFFFFF0
	ManufacturerCodeMax    = 0xFFFFFF
	ManufactureCodeDefault = ManufacturerCodeMin
)

// Object is an instance for Echonet object.
type Object struct {
	Code             []byte
	Properties       []*Property
	annoPropMapSize  int
	annoPropMapBytes []byte
	setPropMapSize   int
	setPropMapBytes  []byte
	getPropMapSize   int
	getPropMapBytes  []byte

	parentNode *Node
}

// NewObject returns a new object.
func NewObject() *Object {
	obj := &Object{
		Code:       make([]byte, 3),
		Properties: make([]*Property, 0),
		parentNode: nil,
	}
	return obj
}

// SetParentNode sets a parent node
func (obj *Object) SetParentNode(node *Node) error {
	obj.parentNode = node
	return nil
}

// GetParentNode returns a parent node
func (obj *Object) GetParentNode() *Node {
	return obj.parentNode
}

// SetCode sets a code to the object
func (obj *Object) Setcode(code uint) error {
	encoding.IntegerToByte(code, obj.Code)
	return nil
}

// GetCode returns the code.
func (obj *Object) GetCode() uint {
	return encoding.ByteToInteger(obj.Code)
}

// IsCode returns true when the object code is the specified code, otherwise false
func (obj *Object) IsCode(code uint) bool {
	if code != obj.GetCode() {
		return false
	}
	return true
}

// SetClassGroupCode sets a class group code to the object
func (obj *Object) SetClassGroupCode(code byte) error {
	obj.Code[0] = code
	return nil
}

// GetClassGroupCode returns the class group code
func (obj *Object) GetClassGroupCode() byte {
	return obj.Code[0]
}

// SetClassCode sets a class code to the object
func (obj *Object) SetClassCode(code byte) error {
	obj.Code[1] = code
	return nil
}

// GetClassCode returns the class group code
func (obj *Object) GetClassCode() byte {
	return obj.Code[1]
}

// SetInstanceCode sets a instance code to the object
func (obj *Object) SetInstanceCode(code byte) error {
	obj.Code[2] = code
	return nil
}

// GetInstanceCode returns the instance code
func (obj *Object) GetInstanceCode() byte {
	return obj.Code[2]
}

// IsDevice returns true when the class group code is the device code, otherwise false
func (obj *Object) IsDevice() bool {
	if ClassGroupDeviceMax < obj.Code[0] {
		return false
	}
	return true
}

// IsProfile returns true when the class group code is the profile code, otherwise false
func (obj *Object) IsProfile() bool {
	if ClassGroupProfile != obj.Code[0] {
		return false
	}
	return true
}

// SetProperty sets a property to the object
/*
func (obj *Object) SetProperty(uEchoPropertyCode code, uEchoPropertyAttr attr)
 {
   uEchoProperty *prop;

   if (!obj)
	 return false;

   if (!uecho_propertylist_set(obj.properties, code, attr, NULL, 0))
	 return false;

   prop = func (obj *Object) _getproperty(obj, code);
   if (!prop)
	 return false;

   if (!uecho_property_setparentobject(prop, obj))
	 return false;

   return func (obj *Object) _updatepropertymaps(obj);
 }


 bool func (obj *Object) _setpropertydata(uEchoObject *obj, uEchoPropertyCode code, byte *data, size_t dataLen)
 {
   if (!obj)
	 return false;

   return uecho_propertylist_setdata(obj.properties, code, data, dataLen);
 }

 bool func (obj *Object) _setpropertyintegerdata(uEchoObject *obj, uEchoPropertyCode code, int data, size_t dataLen)
 {
   if (!obj)
	 return false;

   return uecho_propertylist_setintegerdata(obj.properties, code, data, dataLen);
 }

 bool func (obj *Object) _setpropertybytedata(uEchoObject *obj, uEchoPropertyCode code, byte data)
 {
   if (!obj)
	 return false;

   return uecho_propertylist_setbytedata(obj.properties, code, data);
 }

 bool func (obj *Object) _setpropertyattribute(uEchoObject *obj, uEchoPropertyCode code, uEchoPropertyAttr attr)
 {
   if (!obj)
	 return false;

   return uecho_propertylist_setattribute(obj.properties, code, attr);
 }

 uEchoProperty *func (obj *Object) _getproperties(uEchoObject *obj)
 {
   if (!obj)
	 return NULL;

   return uecho_propertylist_gets(obj.properties);
 }

 bool func (obj *Object) _hasproperty(uEchoObject *obj, uEchoPropertyCode code)
 {
   if (!obj)
	 return false;

   return (func (obj *Object) _getproperty(obj, code) != NULL) ? true : false;
 }
*/

// GetProperty returns the specified property
func (obj *Object) GetProperty(code int) (*Property, error) {
	return nil, nil
}

/*
 uEchoProperty *func (obj *Object) _getpropertywait(uEchoObject *obj, uEchoPropertyCode code, clock_t waitMiliTime)
 {
   uEchoProperty *prop;
   int n;

   for (n=0; n<uEchoWaitRetryCount; n++) {
	 uecho_sleep(waitMiliTime / uEchoWaitRetryCount);
	 prop = func (obj *Object) _getproperty(obj, code);
	 if (prop)
	   return prop;
   }

   return NULL;
 }


 size_t func (obj *Object) _getpropertycount(uEchoObject *obj)
 {
   if (!obj)
	 return 0;

   return uecho_propertylist_size(obj.properties);
 }
 void func (obj *Object) _clearproperties(uEchoObject *obj)
 {
   if (!obj)
	 return;

   uecho_propertylist_clear(obj.properties);
   func (obj *Object) _clearpropertymapcaches(obj);
 }

 int func (obj *Object) _getpropertydatasize(uEchoObject *obj, uEchoPropertyCode code)
 {
   if (!obj)
	 return 0;

   return uecho_propertylist_getdatasize(obj.properties, code);
 }

 byte *func (obj *Object) _getpropertydata(uEchoObject *obj, uEchoPropertyCode code)
 {
   if (!obj)
	 return NULL;

   return uecho_propertylist_getdata(obj.properties, code);
 }

 bool func (obj *Object) _getpropertyintegerdata(uEchoObject *obj, uEchoPropertyCode code, size_t dataLen, int *data)
 {
   if (!obj)
	 return false;

   return uecho_propertylist_getintegerdata(obj.properties, code, dataLen, data);
 }

 bool func (obj *Object) _getpropertybytedata(uEchoObject *obj, uEchoPropertyCode code, byte *data)
 {
   if (!obj)
	 return false;

   return uecho_propertylist_getbytedata(obj.properties, code, data);
 }


 void func (obj *Object) _setmessagelistener(uEchoObject *obj, uEchoObjectMessageListener listener)
 {
   if (!obj)
	 return;

   obj.allMsgListener = listener;
 }

 uEchoObjectMessageListener func (obj *Object) _getmessagelistener(uEchoObject *obj)
 {
   if (!obj)
	 return NULL;

   return obj.allMsgListener;
 }


 bool func (obj *Object) _hasmessagelistener(uEchoObject *obj)
 {
   if (!obj)
	 return false;

   return obj.allMsgListener ? true : false;
 }

 bool func (obj *Object) _setpropertyrequestlistener(uEchoObject *obj, uEchoEsv esv, uEchoPropertyCode code, uEchoPropertyRequestListener listener)
 {
   if (!obj)
	 return false;

   return func (obj *Object) _property_observer_manager_setobserver(obj.propListenerMgr, esv, code, listener);
 }

 uEchoPropertyRequestListener func (obj *Object) _getpropertyrequestlistener(uEchoObject *obj, uEchoEsv esv, uEchoPropertyCode code)
 {
   uEchoObjectPropertyObserver *obs;

   if (!obj)
	 return NULL;

   obs = func (obj *Object) _property_observer_manager_getobserver(obj.propListenerMgr, esv, code);
   if (!obs)
	 return NULL;

   return obs.listener;
 }

 bool func (obj *Object) _haspropertyrequestlistener(uEchoObject *obj, uEchoEsv esv, uEchoPropertyCode code)
 {
   return (func (obj *Object) _getpropertyrequestlistener(obj, esv, code) != NULL) ? true : false;
 }

 bool func (obj *Object) _setpropertywriterequestlistener(uEchoObject *obj, uEchoPropertyCode code, uEchoPropertyRequestListener listener)
 {
   bool isSeccess = true;

   isSeccess &= func (obj *Object) _setpropertyrequestlistener(obj, uEchoEsvWriteRequest, code, listener);
   isSeccess &= func (obj *Object) _setpropertyrequestlistener(obj, uEchoEsvWriteRequestResponseRequired, code, listener);
   isSeccess &= func (obj *Object) _setpropertyrequestlistener(obj, uEchoEsvWriteReadRequest, code, listener);

   return isSeccess;
 }

 bool func (obj *Object) _setpropertyreadlistener(uEchoObject *obj, uEchoPropertyCode code, uEchoPropertyRequestListener listener)
 {
   bool isSeccess = true;

   isSeccess &= func (obj *Object) _setpropertyrequestlistener(obj, uEchoEsvReadRequest, code, listener);
   isSeccess &= func (obj *Object) _setpropertyrequestlistener(obj, uEchoEsvWriteReadRequest, code, listener);

   return isSeccess;

 }

  /****************************************
  * func (obj *Object) _setpropertymap
  ****************************************/

// SetPropertyMap sets a instance code to the object
/*
func (obj *Object) SetPropertyMap(mapCode byte, propCodes byte[]) error
{
   propsCodeSize := len(propCodes)
  propMapData := make(byte, PropertyMapMaxLen + 1)
  propMapData[0] =  len(propCodes);
  propMap = propMapData + 1;

  // propsCodeSize <= uEchoPropertyMapMaxLen

  if (propsCodeSize <= uEchoPropertyMapMaxLen) {
	memcpy(propMap, propCodes, propsCodeSize);
	func (obj *Object) _setproperty(obj, mapCode, uEchoPropertyAttrRead);
	func (obj *Object) _setpropertydata(obj, mapCode, propMapData, (propsCodeSize + 1));
	return true;
  }

  // uEchoPropertyMapMaxLen < propsCodeSize

  for (n=0; n<propsCodeSize; n++) {
	byte propCode;
	propCode = propCodes[n];
	if ((propCode < uEchoPropertyCodeMin) || (uEchoPropertyCodeMax < propCode))
	  continue;
	propByteIdx = (propCode - uEchoPropertyCodeMin) & 0x0F;
	propMap[propByteIdx] |= (((propCode - uEchoPropertyCodeMin) & 0xF0) >> 8) & 0x0F;
  }

  func (obj *Object) _setproperty(obj, mapCode, uEchoPropertyAttrRead);
  func (obj *Object) _setpropertydata(obj, mapCode, propMapData, (uEchoPropertyMapMaxLen + 1));

  return nil
}
*/

// AnnounceMessage announces a message.
func (obj *Object) AnnounceMessage(msg *protocol.Message) error {
	if obj.parentNode == nil {
		return fmt.Errorf(errorParentNodeNotFound)
	}
	msg.SetSourceObjectCode(obj.GetCode())
	return obj.parentNode.AnnounceMessage(msg)
}

/****************************************
 * func (obj *Object) _sendmessage
 ****************************************/

/*
 bool func (obj *Object) _sendmessage(uEchoObject *obj, uEchoObject *dstObj, uEchoMessage *msg)
 {
   uEchoNode *parentNode, *dstParentNode;

   if (!obj || !dstObj)
	 return false;

   parentNode = func (obj *Object) _getparentnode(obj);
   dstParentNode = func (obj *Object) _getparentnode(dstObj);
   if (!parentNode || !dstParentNode)
	 return false;

   uecho_message_setsourceobjectcode(msg, func (obj *Object) _getcode(obj));
   uecho_message_setdestinationobjectcode(msg, func (obj *Object) _getcode(dstObj));

   return uecho_node_sendmessage(parentNode, dstParentNode, msg);
 }

*/
