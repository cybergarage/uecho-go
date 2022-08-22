// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	errorParentNodeNotFound = "parent node not found"
)

// Object is an instance for Echonet object.
type Object struct {
	*PropertyMap
	clsName    string
	name       string
	codes      []byte
	listener   ObjectListener
	parentNode Node
}

// NewObject returns a new object.
func NewObject() *Object {
	obj := &Object{
		PropertyMap: NewPropertyMap(),
		clsName:     "",
		name:        "",
		codes:       make([]byte, ObjectCodeSize),
		listener:    nil,
		parentNode:  nil,
	}

	obj.SetParentObject(obj)

	return obj
}

// NewObjectWithCodes returns a new object of the specified object codes.
func NewObjectWithCodes(codes []byte) (interface{}, error) {
	objCode, err := BytesToObjectCode(codes)
	if err != nil {
		return nil, err
	}

	if isProfileObjectCode(codes[0]) {
		obj := NewProfile()
		obj.SetCode(objCode)
		return obj, nil
	}

	obj := NewDevice()
	obj.SetCode(objCode)
	return obj, nil
}

// SetClassName sets a class name to the object.
func (obj *Object) SetClassName(name string) {
	obj.clsName = name
}

// ClassName returns the class name.
func (obj *Object) ClassName() string {
	return obj.clsName
}

// SetName sets a name to the object.
func (obj *Object) SetName(name string) {
	obj.name = name
}

// Name returns the object name.
func (obj *Object) Name() string {
	return obj.name
}

// SetCode sets a code to the object.
func (obj *Object) SetCode(code ObjectCode) {
	encoding.IntegerToByte(uint(code), obj.codes)
}

// Code returns the code.
func (obj *Object) Code() ObjectCode {
	return ObjectCode(encoding.ByteToInteger(obj.codes))
}

// SetCodes sets codes to the object.
func (obj *Object) SetCodes(codes []byte) {
	copy(obj.codes, codes)
}

// Codes returns the code byte array.
func (obj *Object) Codes() []byte {
	return obj.codes
}

// IsCode returns true when the object code is the specified code, otherwise false.
func (obj *Object) IsCode(code ObjectCode) bool {
	return (code == obj.Code())
}

// Class returns the class of the object.
func (obj *Object) Class() *Class {
	return NewClassWithCodes(obj.codes)
}

// SetClassGroupCode sets a class group code to the object.
func (obj *Object) SetClassGroupCode(code byte) {
	obj.codes[0] = code
}

// ClassGroupCode returns the class group code.
func (obj *Object) ClassGroupCode() byte {
	return obj.codes[0]
}

// SetClassCode sets a class code to the object.
func (obj *Object) SetClassCode(code byte) {
	obj.codes[1] = code
}

// ClassCode returns the class group code.
func (obj *Object) ClassCode() byte {
	return obj.codes[1]
}

// SetInstanceCode sets a instance code to the object.
func (obj *Object) SetInstanceCode(code byte) {
	obj.codes[2] = code
}

// InstanceCode returns the instance code.
func (obj *Object) InstanceCode() byte {
	return obj.codes[2]
}

// IsDevice returns true when the class group code is the device code, otherwise false.
func (obj *Object) IsDevice() bool {
	return !obj.IsProfile()
}

// IsProfile returns true when the class group code is the profile code, otherwise false.
func (obj *Object) IsProfile() bool {
	return isProfileObjectCode(obj.codes[0])
}

// SetParentNode sets a parent node.
func (obj *Object) SetParentNode(node Node) {
	obj.parentNode = node
}

// ParentNode returns a parent node.
func (obj *Object) ParentNode() Node {
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

// Copy copies the object instance without the data.
func (obj *Object) Copy() *Object {
	newObj := &Object{
		PropertyMap: NewPropertyMap(),
		clsName:     obj.ClassName(),
		name:        obj.Name(),
		codes:       make([]byte, ObjectCodeSize),
		listener:    nil,
		parentNode:  nil,
	}

	newObj.SetCode(newObj.Code())
	newObj.SetParentObject(newObj)
	for _, prop := range obj.GetProperties() {
		newObj.AddProperty(prop.Copy())
	}

	return newObj
}
