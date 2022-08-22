// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"bytes"
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	PropertyCodeMin           = 0x80
	PropertyCodeMax           = 0xFF
	PropertyMapFormat1MaxSize = 15
	PropertyMapFormat2Size    = 18
	PropertyMapFormatMaxSize  = PropertyMapFormat2Size
)

const (
	errorPropertyNoParentNode = "Property has no parent node"
	errorPropertyNoData       = "Property has no data"
)

// PropertyCode is a type for property code.
type PropertyCode = protocol.PropertyCode

// PropertyAttr is a type for property attribute.
type PropertyAttr uint

const (
	PropertyAttrNone    = PropertyAttr(0x00)
	PropertyAttrGet     = PropertyAttr(0x01)
	PropertyAttrSet     = PropertyAttr(0x02)
	PropertyAttrAnno    = PropertyAttr(0x10)
	PropertyAttrGetSet  = PropertyAttrGet | PropertyAttrSet
	PropertyAttrGetAnno = PropertyAttrGet | PropertyAttrAnno
)

// Property is an instance for Echonet property.
type Property struct {
	Name         string
	code         PropertyCode
	Attr         PropertyAttr
	Data         []byte
	ParentObject *Object
}

// NewProperty returns a new property.
func NewProperty() *Property {
	return &Property{
		Name:         "",
		code:         0,
		Attr:         PropertyAttrNone,
		Data:         make([]byte, 0),
		ParentObject: nil,
	}
}

// NewPropertyWithCode returns a new property with the specified property code.
func NewPropertyWithCode(code PropertyCode) *Property {
	prop := NewProperty()
	prop.SetCode(code)
	return prop
}

// NewPropertiesWithCodes returns a new properties with the specified property codes.
func NewPropertiesWithCodes(codes []PropertyCode) []*Property {
	props := make([]*Property, len(codes))
	for n, code := range codes {
		props[n] = NewPropertyWithCode(code)
	}
	return props
}

// SetParentObject sets a parent object into the property.
func (prop *Property) SetParentObject(obj *Object) {
	prop.ParentObject = obj
}

// GetParentObject returns the parent object.
func (prop *Property) GetParentObject() *Object {
	return prop.ParentObject
}

// GetNode returns a parent node of the parent object.
func (prop *Property) GetNode() Node {
	parentObj := prop.GetParentObject()
	if parentObj == nil {
		return nil
	}
	return parentObj.ParentNode()
}

// SetName sets a name to the property.
func (prop *Property) SetName(name string) {
	prop.Name = name
}

// GetName returns the property name.
func (prop *Property) GetName() string {
	return prop.Name
}

// SetCode sets a specified code to the property.
func (prop *Property) SetCode(code PropertyCode) {
	prop.code = code
}

// Code returns the property code.
func (prop *Property) Code() PropertyCode {
	return prop.code
}

// ClearData clears the property data.
func (prop *Property) ClearData() {
	prop.Data = make([]byte, 0)
}

// Size return the property data size.
func (prop *Property) Size() int {
	return len(prop.Data)
}

// SetAttribute sets an attribute to the property.
func (prop *Property) SetAttribute(attr PropertyAttr) {
	prop.Attr = attr
}

// GetAttribute returns the property attribute.
func (prop *Property) GetAttribute() PropertyAttr {
	return prop.Attr
}

// IsReadable returns true when the property attribute is readable, otherwise false.
func (prop *Property) IsReadable() bool {
	return ((prop.Attr & PropertyAttrGet) != 0)
}

// IsWritable returns true when the property attribute is writable, otherwise false.
func (prop *Property) IsWritable() bool {
	return ((prop.Attr & PropertyAttrSet) != 0)
}

// IsReadOnly returns true when the property attribute is read only, otherwise false.
func (prop *Property) IsReadOnly() bool {
	if (prop.Attr & PropertyAttrGet) == 0 {
		return false
	}

	if (prop.Attr & PropertyAttrSet) != 0 {
		return false
	}

	return true
}

// IsWriteOnly returns true when the property attribute is write only, otherwise false.
func (prop *Property) IsWriteOnly() bool {
	if (prop.Attr & PropertyAttrSet) == 0 {
		return false
	}
	if (prop.Attr & PropertyAttrGet) != 0 {
		return false
	}
	return true
}

// isAnnounceable returns true when the property attribute is announcement, otherwise false.
func (prop *Property) isAnnounceable() bool {
	return ((prop.Attr & PropertyAttrAnno) != 0)
}

// IsAvailableService returns true whether the specified service can execute, otherwise false.
func (prop *Property) IsAvailableService(esv protocol.ESV) bool {
	switch esv {
	case protocol.ESVWriteRequest:
		if prop.IsWritable() {
			return true
		}
		return false
	case protocol.ESVWriteRequestResponseRequired:
		if prop.IsWritable() {
			return true
		}
		return false
	case protocol.ESVReadRequest:
		if prop.IsReadable() {
			return true
		}
		return false
	case protocol.ESVNotificationRequest:
		if prop.isAnnounceable() {
			return true
		}
		return false
	case protocol.ESVWriteReadRequest:
		if prop.IsWritable() && prop.IsReadable() {
			return true
		}
		return false
	case protocol.ESVNotificationResponseRequired:
		if prop.isAnnounceable() {
			return true
		}
		return false
	}
	return false
}

// SetData sets a specified data to the property.
func (prop *Property) SetData(data []byte) {
	prop.Data = make([]byte, len(data))
	copy(prop.Data, data)

	// (D) Basic sequence for autonomous notification.

	if prop.isAnnounceable() {
		prop.announce()
	}
}

// SetByteData is an alias of SetData.
func (prop *Property) SetByteData(data []byte) {
	prop.SetData(data)
}

// SetIntegerData sets a specified integer data to the property.
func (prop *Property) SetIntegerData(data uint, size uint) {
	binData := make([]byte, size)
	encoding.IntegerToByte(data, binData)
	prop.SetData(binData)
}

// GetData returns the property data.
func (prop *Property) GetData() []byte {
	return prop.Data
}

// GetByteData returns a byte value of the property data.
func (prop *Property) GetByteData() (byte, error) {
	if len(prop.Data) == 0 {
		return 0, fmt.Errorf(errorPropertyNoData)
	}
	return prop.Data[0], nil
}

// GetStringData returns a byte value of the property string data.
func (prop *Property) GetStringData() (string, error) {
	return string(prop.Data), nil
}

// GetIntegerData returns a integer value of the property integer data.
func (prop *Property) GetIntegerData() (uint, error) {
	if len(prop.Data) == 0 {
		return 0, fmt.Errorf(errorPropertyNoData)
	}
	return encoding.ByteToInteger(prop.GetData()), nil
}

// announce announces the property.
func (prop *Property) announce() error {
	parentNode, ok := prop.GetNode().(*LocalNode)
	if !ok || parentNode == nil {
		return fmt.Errorf(errorPropertyNoParentNode)
	}

	if !parentNode.IsRunning() {
		return nil
	}

	return parentNode.AnnounceProperty(prop)
}

// toProtocolProperty returns the new property of the property.
func (prop *Property) toProtocolProperty() *protocol.Property {
	newProp := protocol.NewProperty()
	newProp.SetCode(prop.Code())
	newProp.SetData(prop.GetData())
	return newProp
}

// Equals returns true if the specified property is same, otherwise false.
func (prop *Property) Equals(otherProp *Property) bool {
	if prop.Code() != otherProp.Code() {
		return false
	}
	if prop.GetAttribute() != otherProp.GetAttribute() {
		return false
	}
	if !bytes.Equal(prop.GetData(), otherProp.GetData()) {
		return false
	}
	return true
}

// Copy copies the property instance without the data.
func (prop *Property) Copy() *Property {
	return &Property{
		Name:         prop.Name,
		code:         prop.code,
		Attr:         prop.Attr,
		Data:         make([]byte, 0),
		ParentObject: nil,
	}
}
