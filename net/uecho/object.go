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
)

// Object is an instance for Echonet object.
type Object struct {
	*PropertyMap
	Code       []byte
	listeners  []ObjectListener
	parentNode *Node
}

// NewObject returns a new object.
func NewObject() *Object {
	obj := &Object{
		PropertyMap: NewPropertyMap(),
		Code:        make([]byte, 3),
		listeners:   make([]ObjectListener, 0),
		parentNode:  nil,
	}

	obj.SetParentObject(obj)

	return obj
}

// SetCode sets a code to the object
func (obj *Object) SetCode(code uint) {
	encoding.IntegerToByte(code, obj.Code)
}

// GetCode returns the code.
func (obj *Object) GetCode() uint {
	return encoding.ByteToInteger(obj.Code)
}

// GetCodes returns the code byte array.
func (obj *Object) GetCodes() []byte {
	return obj.Code
}

// IsCode returns true when the object code is the specified code, otherwise false.
func (obj *Object) IsCode(code uint) bool {
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
	if obj.Code[0] != ProfileClassGroupCode {
		return false
	}
	return true
}

// SetParentNode sets a parent node.
func (obj *Object) SetParentNode(node *Node) {
	obj.parentNode = node
}

// GetParentNode returns a parent node.
func (obj *Object) GetParentNode() *Node {
	return obj.parentNode
}

// AddListener add the specified listener to the node.
func (obj *Object) AddListener(l ObjectListener) {
	obj.listeners = append(obj.listeners, l)
}

// GetListeners returns all listeners of the node.
func (obj *Object) GetListeners() []ObjectListener {
	return obj.listeners
}

// notifyPropertyRequest notifies a request to the listeners.
func (obj *Object) notifyPropertyRequest(esv protocol.ESV, prop *protocol.Property) {
	for _, l := range obj.listeners {
		l.PropertyRequestReceived(obj, esv, prop)
	}
}

// AnnounceMessage announces a message.
func (obj *Object) AnnounceMessage(msg *protocol.Message) error {
	if obj.parentNode == nil {
		return fmt.Errorf(errorParentNodeNotFound)
	}
	msg.SetSourceObjectCode(obj.GetCode())
	return obj.parentNode.AnnounceMessage(msg)
}

// SendMessage send a message to the node of the destination object.
func (obj *Object) SendMessage(dstObj *Object, msg *protocol.Message) error {
	parentNode := obj.GetParentNode()
	dstParentNode := dstObj.GetParentNode()

	if parentNode == nil || dstParentNode == nil {
		return nil
	}

	msg.SetSourceObjectCode(obj.GetCode())
	msg.SetDestinationObjectCode(dstObj.GetCode())

	return parentNode.SendMessage(dstParentNode, msg)
}
