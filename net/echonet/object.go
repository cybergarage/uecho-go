// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	errorInvalidObjectCodes = "Invalid Object Code : %s"
	errorParentNodeNotFound = "Parent node not found"
)

const (
	ObjectCodeMin     = 0x000000
	ObjectCodeMax     = 0xFFFFFF
	ObjectCodeSize    = 3
	ObjectCodeUnknown = ObjectCodeMin
)

// ObjectCode is a type for object code.
type ObjectCode = protocol.ObjectCode

// Object is an instance for Echonet object.
type Object struct {
	*PropertyMap
	Code       []byte
	listener   ObjectListener
	parentNode Node
}

// NewObject returns a new object.
func NewObject() *Object {
	obj := &Object{
		PropertyMap: NewPropertyMap(),
		Code:        make([]byte, ObjectCodeSize),
		listener:    nil,
		parentNode:  nil,
	}

	obj.SetParentObject(obj)

	return obj
}

// NewObjectWithCodes returns a new object of the specified object codes.
func NewObjectWithCodes(codes []byte) (interface{}, error) {
	if len(codes) != ObjectCodeSize {
		return nil, fmt.Errorf(errorInvalidObjectCodes, string(codes))
	}

	if isProfileObjectCode(codes[0]) {
		obj := NewProfile()
		obj.SetCodes(codes)
		return obj, nil
	}

	obj := NewDevice()
	obj.SetCodes(codes)
	return obj, nil
}

// SetCode sets a code to the object
func (obj *Object) SetCode(code ObjectCode) {
	encoding.IntegerToByte(uint(code), obj.Code)
}

// GetCode returns the code.
func (obj *Object) GetCode() ObjectCode {
	return ObjectCode(encoding.ByteToInteger(obj.Code))
}

// SetCodes sets codes to the object
func (obj *Object) SetCodes(codes []byte) {
	copy(obj.Code, codes)
}

// GetCodes returns the code byte array.
func (obj *Object) GetCodes() []byte {
	return obj.Code
}

// IsCode returns true when the object code is the specified code, otherwise false.
func (obj *Object) IsCode(code ObjectCode) bool {
	if code != obj.GetCode() {
		return false
	}
	return true
}

// GetClass returns the class of the object.
func (obj *Object) GetClass() *Class {
	return NewClassWithCodes(obj.Code)
}

// SetClassGroupCode sets a class group code to the object.
func (obj *Object) SetClassGroupCode(code byte) {
	obj.Code[0] = code
}

// GetClassGroupCode returns the class group code.
func (obj *Object) GetClassGroupCode() byte {
	return obj.Code[0]
}

// SetClassCode sets a class code to the object.
func (obj *Object) SetClassCode(code byte) {
	obj.Code[1] = code
}

// GetClassCode returns the class group code.
func (obj *Object) GetClassCode() byte {
	return obj.Code[1]
}

// SetInstanceCode sets a instance code to the object.
func (obj *Object) SetInstanceCode(code byte) {
	obj.Code[2] = code
}

// GetInstanceCode returns the instance code.
func (obj *Object) GetInstanceCode() byte {
	return obj.Code[2]
}

// IsDevice returns true when the class group code is the device code, otherwise false.
func (obj *Object) IsDevice() bool {
	if obj.IsProfile() {
		return false
	}
	return true
}

// IsProfile returns true when the class group code is the profile code, otherwise false.
func (obj *Object) IsProfile() bool {
	return isProfileObjectCode(obj.Code[0])
}

// SetParentNode sets a parent node.
func (obj *Object) SetParentNode(node Node) {
	obj.parentNode = node
}

// GetParentNode returns a parent node.
func (obj *Object) GetParentNode() Node {
	return obj.parentNode
}

// SetListener sets a listener to the object.
func (obj *Object) SetListener(l ObjectListener) {
	obj.listener = l
}

// notifyPropertyRequest notifies a request to the object listener.
func (obj *Object) notifyPropertyRequest(esv protocol.ESV, prop *protocol.Property) error {
	if obj.listener == nil {
		return nil
	}
	return obj.listener.PropertyRequestReceived(obj, esv, prop)
}
