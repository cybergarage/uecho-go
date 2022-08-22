// Copyright (C) 2018 Satoshi Konno. All rights reserved.
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
	ClassName  string
	Name       string
	codes      []byte
	listener   ObjectListener
	parentNode Node
}

// NewObject returns a new object.
func NewObject() *Object {
	obj := &Object{
		PropertyMap: NewPropertyMap(),
		ClassName:   "",
		Name:        "",
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
	obj.ClassName = name
}

// GetClassName returns the class name.
func (obj *Object) GetClassName() string {
	return obj.ClassName
}

// SetName sets a name to the object.
func (obj *Object) SetName(name string) {
	obj.Name = name
}

// GetName returns the object name.
func (obj *Object) GetName() string {
	return obj.Name
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

// GetClass returns the class of the object.
func (obj *Object) GetClass() *Class {
	return NewClassWithCodes(obj.codes)
}

// SetClassGroupCode sets a class group code to the object.
func (obj *Object) SetClassGroupCode(code byte) {
	obj.codes[0] = code
}

// GetClassGroupCode returns the class group code.
func (obj *Object) GetClassGroupCode() byte {
	return obj.codes[0]
}

// SetClassCode sets a class code to the object.
func (obj *Object) SetClassCode(code byte) {
	obj.codes[1] = code
}

// GetClassCode returns the class group code.
func (obj *Object) GetClassCode() byte {
	return obj.codes[1]
}

// SetInstanceCode sets a instance code to the object.
func (obj *Object) SetInstanceCode(code byte) {
	obj.codes[2] = code
}

// GetInstanceCode returns the instance code.
func (obj *Object) GetInstanceCode() byte {
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

// Copy copies the object instance without the data.
func (obj *Object) Copy() *Object {
	newObj := &Object{
		PropertyMap: NewPropertyMap(),
		ClassName:   obj.GetClassName(),
		Name:        obj.GetName(),
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
