// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	PropertyCodeMin           = 0x80
	PropertyCodeMax           = 0xFF
	PropertyMapFormat1MaxSize = 15
	PropertyMapFormat2MapSize = 16
	PropertyMapFormatMaxSize  = PropertyMapFormat2MapSize + 1
)

const (
	errorPropertyNoParentNode   = "property has no parent node"
	errorPropertyNoData         = "property has no data"
	errorInvalidPropertyCode    = "invalid property code: %02X"
	errorInvalidPropertyMapData = "invalid property map data: %0s"
)

// PropertyCode is a type for property code.
type PropertyCode = protocol.PropertyCode

type Property interface {
	// SetParentObject sets a parent object into the property.
	SetParentObject(obj Object)
	// ParentObject returns the parent object.
	ParentObject() Object
	// Node returns a parent node of the parent object.
	Node() Node
	// SetName sets a name to the property.
	SetName(name string) *property
	// Name returns the property name.
	Name() string
	// Code returns the property code.
	Code() PropertyCode
	// ClearData clears the property data.
	ClearData()
	// Size return the property data size.
	Size() int
	// SetCode sets a specified code to the property.
	SetCode(code PropertyCode) Property
	// SetReadAttribute sets an attribute to the read property.
	SetReadAttribute(attr PropertyAttribute) Property
	// SetWriteAttribute sets an attribute to the write property.
	SetWriteAttribute(attr PropertyAttribute) Property
	// SetAnnoAttribute sets an attribute to the announce property.
	SetAnnoAttribute(attr PropertyAttribute) Property
	// GetAttribute returns the get attribute.
	ReadAttribute() PropertyAttribute
	// SetAttribute returns the set attribute.
	WriteAttribute() PropertyAttribute
	// AnnoAttribute returns the announce attribute.
	AnnoAttribute() PropertyAttribute
	// IsReadable returns true when the get attribute is readable, otherwise false.
	IsReadable() bool
	// IsWritable returns true when the set attribute is writable, otherwise false.
	IsWritable() bool
	// IsAnnounceable returns true when the anno attribute is announcement, otherwise false.
	IsAnnounceable() bool
	// IsReadRequired returns true when the get attribute is required, otherwise false.
	IsReadRequired() bool
	// IsWriteRequired returns true when the set attribute is required, otherwise false.
	IsWriteRequired() bool
	// IsAnnounceRequired returns true when the announce attribute is required, otherwise false.
	IsAnnounceRequired() bool
	// IsReadOnly returns true when the property attribute is read only, otherwise false.
	IsReadOnly() bool
	// IsWriteOnly returns true when the property attribute is write only, otherwise false.
	IsWriteOnly() bool
	// IsAvailableService returns true whether the specified service can execute, otherwise false.
	IsAvailableService(esv protocol.ESV) bool
	// SetData sets a specified data to the property.
	SetData(data []byte) Property
	// SetByteData is an alias of SetData.
	SetByteData(data []byte) Property
	// SetIntegerData sets a specified integer data to the property.
	SetIntegerData(data uint, size uint) Property
	// Data returns the property data.
	Data() []byte
	// ByteData returns a byte value of the property data.
	ByteData() (byte, error)
	// StringData returns a byte value of the property string data.
	StringData() (string, error)
	// IntegerData returns a integer value of the property integer data.
	IntegerData() (uint, error)
	// PropertyMapData returns a property map.
	PropertyMapData() ([]PropertyCode, error)
	// announce announces the property.
	announce() error
	// toProtocolProperty returns the new property of the property.
	toProtocolProperty() *protocol.Property
	// Equals returns true if the specified property is same, otherwise false.
	Equals(otherProp *property) bool
	// Copy copies the property instance without the data.
	Copy() Property
}

// property is an instance for Echonet property.
type property struct {
	name         string
	code         PropertyCode
	data         []byte
	parentObject Object
	getAttr      PropertyAttribute
	setAttr      PropertyAttribute
	annoAttr     PropertyAttribute
}

// NewProperty returns a new property.
func NewProperty() Property {
	return &property{
		name:         "",
		code:         0,
		data:         make([]byte, 0),
		parentObject: nil,
		getAttr:      Prohibited,
		setAttr:      Prohibited,
		annoAttr:     Prohibited,
	}
}

// NewPropertyWithCode returns a new property with the specified property code.
func NewPropertyWithCode(code PropertyCode) Property {
	prop := NewProperty()
	prop.SetCode(code)
	return prop
}

// NewPropertiesWithCodes returns a new properties with the specified property codes.
func NewPropertiesWithCodes(codes []PropertyCode) []Property {
	props := make([]Property, len(codes))
	for n, code := range codes {
		props[n] = NewPropertyWithCode(code)
	}
	return props
}

// SetParentObject sets a parent object into the property.
func (prop *property) SetParentObject(obj Object) {
	prop.parentObject = obj
}

// ParentObject returns the parent object.
func (prop *property) ParentObject() Object {
	return prop.parentObject
}

// Node returns a parent node of the parent object.
func (prop *property) Node() Node {
	parentObj := prop.ParentObject()
	if parentObj == nil {
		return nil
	}
	return parentObj.ParentNode()
}

// SetName sets a name to the property.
func (prop *property) SetName(name string) *property {
	prop.name = name
	return prop
}

// Name returns the property name.
func (prop *property) Name() string {
	return prop.name
}

// SetCode sets a specified code to the property.
func (prop *property) SetCode(code PropertyCode) Property {
	prop.code = code
	return prop
}

// Code returns the property code.
func (prop *property) Code() PropertyCode {
	return prop.code
}

// ClearData clears the property data.
func (prop *property) ClearData() {
	prop.data = make([]byte, 0)
}

// Size return the property data size.
func (prop *property) Size() int {
	return len(prop.data)
}

// SetReadAttribute sets an attribute to the read property.
func (prop *property) SetReadAttribute(attr PropertyAttribute) Property {
	prop.getAttr = attr
	return prop
}

// SetWriteAttribute sets an attribute to the write property.
func (prop *property) SetWriteAttribute(attr PropertyAttribute) Property {
	prop.setAttr = attr
	return prop
}

// SetAnnoAttribute sets an attribute to the announce property.
func (prop *property) SetAnnoAttribute(attr PropertyAttribute) Property {
	prop.annoAttr = attr
	return prop
}

// GetAttribute returns the get attribute.
func (prop *property) ReadAttribute() PropertyAttribute {
	return prop.getAttr
}

// SetAttribute returns the set attribute.
func (prop *property) WriteAttribute() PropertyAttribute {
	return prop.setAttr
}

// AnnoAttribute returns the announce attribute.
func (prop *property) AnnoAttribute() PropertyAttribute {
	return prop.annoAttr
}

// IsReadable returns true when the get attribute is readable, otherwise false.
func (prop *property) IsReadable() bool {
	return !prop.getAttr.IsProhibited()
}

// IsWritable returns true when the set attribute is writable, otherwise false.
func (prop *property) IsWritable() bool {
	return !prop.setAttr.IsProhibited()
}

// IsAnnounceable returns true when the anno attribute is announcement, otherwise false.
func (prop *property) IsAnnounceable() bool {
	return !prop.annoAttr.IsProhibited()
}

// IsReadRequired returns true when the get attribute is required, otherwise false.
func (prop *property) IsReadRequired() bool {
	return (prop.getAttr.IsRequired())
}

// IsWriteRequired returns true when the set attribute is required, otherwise false.
func (prop *property) IsWriteRequired() bool {
	return (prop.setAttr.IsRequired())
}

// IsAnnounceRequired returns true when the announce attribute is required, otherwise false.
func (prop *property) IsAnnounceRequired() bool {
	return (prop.annoAttr.IsRequired())
}

// IsReadOnly returns true when the property attribute is read only, otherwise false.
func (prop *property) IsReadOnly() bool {
	if prop.IsWritable() {
		return false
	}
	if !prop.IsReadable() {
		return false
	}
	return true
}

// IsWriteOnly returns true when the property attribute is write only, otherwise false.
func (prop *property) IsWriteOnly() bool {
	if prop.IsReadable() {
		return false
	}
	if !prop.IsWritable() {
		return false
	}
	return true
}

// IsAvailableService returns true whether the specified service can execute, otherwise false.
func (prop *property) IsAvailableService(esv protocol.ESV) bool {
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
		if prop.IsAnnounceable() {
			return true
		}
		return false
	case protocol.ESVWriteReadRequest:
		if prop.IsWritable() && prop.IsReadable() {
			return true
		}
		return false
	case protocol.ESVNotificationResponseRequired:
		if prop.IsAnnounceable() {
			return true
		}
		return false
	}
	return false
}

// SetData sets a specified data to the property.
func (prop *property) SetData(data []byte) Property {
	prop.data = make([]byte, len(data))
	copy(prop.data, data)

	// (D) Basic sequence for autonomous notification.

	if prop.IsAnnounceable() {
		prop.announce()
	}

	return prop
}

// SetByteData is an alias of SetData.
func (prop *property) SetByteData(data []byte) Property {
	return prop.SetData(data)
}

// SetIntegerData sets a specified integer data to the property.
func (prop *property) SetIntegerData(data uint, size uint) Property {
	binData := make([]byte, size)
	encoding.IntegerToByte(data, binData)
	return prop.SetData(binData)
}

// Data returns the property data.
func (prop *property) Data() []byte {
	return prop.data
}

// ByteData returns a byte value of the property data.
func (prop *property) ByteData() (byte, error) {
	if len(prop.data) == 0 {
		return 0, fmt.Errorf(errorPropertyNoData)
	}
	return prop.data[0], nil
}

// StringData returns a byte value of the property string data.
func (prop *property) StringData() (string, error) {
	return string(prop.data), nil
}

// IntegerData returns a integer value of the property integer data.
func (prop *property) IntegerData() (uint, error) {
	if len(prop.data) == 0 {
		return 0, fmt.Errorf(errorPropertyNoData)
	}
	return encoding.ByteToInteger(prop.Data()), nil
}

// PropertyMapData returns a property map.
func (prop *property) PropertyMapData() ([]PropertyCode, error) {
	switch prop.code {
	case ObjectGetPropertyMap, ObjectSetPropertyMap, ObjectAnnoPropertyMap:
		if len(prop.data) == 0 {
			return nil, fmt.Errorf(errorInvalidPropertyMapData, "")
		}
		propMapCount := int(prop.data[0])
		switch {
		case isPropertyMapDescriptionFormat1(propMapCount):
			if len(prop.data) != (propMapCount + 1) {
				return nil, fmt.Errorf(errorInvalidPropertyMapData, hex.EncodeToString(prop.data))
			}
			codes := make([]PropertyCode, 0)
			for n := range propMapCount {
				codes = append(codes, PropertyCode(prop.data[n+1]))
			}
			return codes, nil
		case isPropertyMapDescriptionFormat2(propMapCount):
			if len(prop.data) != (PropertyMapFormat2MapSize + 1) {
				return nil, fmt.Errorf(errorInvalidPropertyMapData, hex.EncodeToString(prop.data))
			}
			codes := make([]PropertyCode, 0)
			for n := range PropertyMapFormat2MapSize {
				codes = append(codes, propertyMapFormat2ByteToCodes(n, prop.data[n+1])...)
			}
			return codes, nil
		}
	}
	return nil, fmt.Errorf(errorInvalidPropertyCode, prop.code)
}

// announce announces the property.
func (prop *property) announce() error {
	parentNode, ok := prop.Node().(localNodeHelper)
	if !ok || parentNode == nil {
		return fmt.Errorf(errorPropertyNoParentNode)
	}

	if !parentNode.IsRunning() {
		return nil
	}

	return parentNode.AnnounceProperty(prop)
}

// toProtocolProperty returns the new property of the property.
func (prop *property) toProtocolProperty() *protocol.Property {
	newProp := protocol.NewProperty()
	newProp.SetCode(prop.Code())
	newProp.SetData(prop.Data())
	return newProp
}

// Equals returns true if the specified property is same, otherwise false.
func (prop *property) Equals(otherProp *property) bool {
	if prop.Code() != otherProp.Code() {
		return false
	}
	if !bytes.Equal(prop.Data(), otherProp.Data()) {
		return false
	}
	return true
}

// Copy copies the property instance without the data.
func (prop *property) Copy() Property {
	return &property{
		name:         prop.name,
		code:         prop.code,
		getAttr:      prop.getAttr,
		setAttr:      prop.setAttr,
		annoAttr:     prop.annoAttr,
		data:         make([]byte, 0),
		parentObject: nil,
	}
}
