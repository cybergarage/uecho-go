// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// Object is an interface to represent an object.
type Object interface {
	// ParentNode returns the parent node of the object.
	ParentNode() Node
	// ClassName returns the class name.
	ClassName() string
	// Name returns the object name.
	Name() string
	// Code returns the object code.
	Code() ObjectCode
	// IsCode returns true if the object code is the specified code.
	IsCode(code ObjectCode) bool
	// Codes returns the object codes.
	Codes() []byte
	// Class returns the class of the object.
	Class() Class
	// ClassGroupCode returns the class group code.
	ClassGroupCode() byte
	// ClassCode returns the class code.
	ClassCode() byte
	// InstanceCode returns the instance code.
	InstanceCode() byte
	// IsDevice returns true if the object is a device.
	IsDevice() bool
	// IsProfile returns true if the object is a profile.
	IsProfile() bool
	// Properties returns the properties of the object.
	Properties() []Property
	// HasProperty returns true if the object has the specified property.
	HasProperty(code PropertyCode) bool
	// LookupProperty returns the specified property of the object.
	LookupProperty(code PropertyCode) (Property, bool)
	// SetListener sets the listener of the object.
	SetListener(l ObjectListener)
	// ObjectMutator returns the object mutator.
	ObjectMutator
	// ObjectHelper returns the object helper.
	ObjectHelper
	// objectInternal returns the object internal interface.
	objectInternal
}

// ObjectMutator is an interface to mutate the object.
type ObjectMutator interface {
	// SetClassName sets the class name of the object.
	SetClassName(name string)
	// SetName sets the name of the object.
	SetName(name string)
	// SetCode sets the code of the object.
	SetCode(code ObjectCode)
	// SetCodes sets the codes of the object.
	SetCodes(codes []byte)
	// SetClassGroupCode sets the class group code of the object.
	SetClassGroupCode(code byte)
	// SetClassCode sets the class code of the object.
	SetClassCode(code byte)
	// SetInstanceCode sets the instance code of the object.
	SetInstanceCode(code byte)
	// SetParentNode sets the parent node of the object.
	SetParentNode(node Node)
	// AddProperty adds a property to the object.
	AddProperty(prop Property)
}

// ObjectHelper is an interface to help the object.
type ObjectHelper interface {
	// LookupPropertyData returns the specified property data in the object.
	LookupPropertyData(propCode PropertyCode) ([]byte, error)
	// LookupPropertyByte returns the specified property byte data in the object.
	LookupPropertyByte(propCode PropertyCode) (byte, error)
	// LookupPropertyInteger returns the specified property integer data in the object.
	LookupPropertyInteger(propCode PropertyCode) (uint, error)
	// SetPropertyData sets a data to the existing property.
	SetPropertyData(propCode PropertyCode, propData []byte) error
	// SetPropertyByte sets a byte to the existing property.
	SetPropertyByte(propCode PropertyCode, propData byte) error
	// SetPropertyInteger sets a integer to the existing property.
	SetPropertyInteger(propCode PropertyCode, propData uint, propSize uint) error
}

// objectInternal is an interface for internal use of the object.
type objectInternal interface {
	// notifyPropertyRequest notifies a request to the object listener.
	notifyPropertyRequest(esv protocol.ESV, prop *protocol.Property) error
}

type object struct {
	*propertyMap

	clsName    string
	name       string
	codes      []byte
	listener   ObjectListener
	parentNode Node
}

func newObject() *object {
	obj := &object{
		propertyMap: newPropertyMap(),
		clsName:     "",
		name:        "",
		codes:       make([]byte, ObjectCodeSize),
		listener:    nil,
		parentNode:  nil,
	}

	obj.SetParentObject(obj)

	return obj
}

// NewObject returns a new object.
func NewObject() Object {
	return newObject()
}

// NewObjectWithCodeBytes returns a new object of the specified object codes.
func NewObjectWithCodeBytes(codes []byte) (any, error) {
	objCode, err := NewObjectCodeFromBytes(codes)
	if err != nil {
		return nil, err
	}

	if isProfileObjectCode(codes[0]) {
		obj := NewProfileWithCode(objCode)
		return obj, nil
	}

	return NewDeviceWithCode(objCode), nil
}

// NewObjectWithCode returns a new object of the specified object code.
func NewObjectWithCode(code ObjectCode) (any, error) {
	return NewObjectWithCodeBytes(code.Bytes())
}

// SetClassName sets a class name to the object.
func (obj *object) SetClassName(name string) {
	obj.clsName = name
}

// ClassName returns the class name.
func (obj *object) ClassName() string {
	return obj.clsName
}

// SetName sets a name to the object.
func (obj *object) SetName(name string) {
	obj.name = name
}

// Name returns the object name.
func (obj *object) Name() string {
	return obj.name
}

// SetCode sets a code to the object.
func (obj *object) SetCode(code ObjectCode) {
	encoding.IntegerToByte(uint(code), obj.codes)
}

// Code returns the code.
func (obj *object) Code() ObjectCode {
	return ObjectCode(encoding.ByteToInteger(obj.codes))
}

// SetCodes sets codes to the object.
func (obj *object) SetCodes(codes []byte) {
	copy(obj.codes, codes)
}

// Codes returns the code byte array.
func (obj *object) Codes() []byte {
	return obj.codes
}

// IsCode returns true when the object code is the specified code, otherwise false.
func (obj *object) IsCode(code ObjectCode) bool {
	return (code == obj.Code())
}

// Class returns the class of the object.
func (obj *object) Class() Class {
	class, _ := NewClass(WithClassBytes(obj.codes[:2]))
	return class
}

// SetClassGroupCode sets a class group code to the object.
func (obj *object) SetClassGroupCode(code byte) {
	obj.codes[0] = code
}

// ClassGroupCode returns the class group code.
func (obj *object) ClassGroupCode() byte {
	return obj.codes[0]
}

// SetClassCode sets a class code to the object.
func (obj *object) SetClassCode(code byte) {
	obj.codes[1] = code
}

// ClassCode returns the class group code.
func (obj *object) ClassCode() byte {
	return obj.codes[1]
}

// SetInstanceCode sets a instance code to the object.
func (obj *object) SetInstanceCode(code byte) {
	obj.codes[2] = code
}

// InstanceCode returns the instance code.
func (obj *object) InstanceCode() byte {
	return obj.codes[2]
}

// IsDevice returns true when the class group code is the device code, otherwise false.
func (obj *object) IsDevice() bool {
	return !obj.IsProfile()
}

// IsProfile returns true when the class group code is the profile code, otherwise false.
func (obj *object) IsProfile() bool {
	return isProfileObjectCode(obj.codes[0])
}

// SetParentNode sets a parent node.
func (obj *object) SetParentNode(node Node) {
	obj.parentNode = node
}

// ParentNode returns a parent node.
func (obj *object) ParentNode() Node {
	return obj.parentNode
}

// SetListener sets a listener to the object.
func (obj *object) SetListener(l ObjectListener) {
	obj.listener = l
}

// notifyPropertyRequest notifies a request to the object listener.
func (obj *object) notifyPropertyRequest(esv protocol.ESV, prop *protocol.Property) error {
	if obj.listener == nil {
		return nil
	}
	return obj.listener.PropertyRequestReceived(obj, esv, prop)
}

// Copy copies the object instance without the data.
func (obj *object) Copy() *object {
	newObj := &object{
		propertyMap: newPropertyMap(),
		clsName:     obj.ClassName(),
		name:        obj.Name(),
		codes:       make([]byte, ObjectCodeSize),
		listener:    nil,
		parentNode:  nil,
	}

	newObj.SetCode(newObj.Code())
	newObj.SetParentObject(newObj)
	for _, prop := range obj.Properties() {
		newObj.AddProperty(prop.Copy())
	}

	return newObj
}
