// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

const (
	PropertyAttributeNone          = 0x00
	PropertyAttributeRead          = 0x01
	PropertyAttributeWrite         = 0x02
	PropertyAttributeAnno          = 0x10
	PropertyAttributeReadWrite     = PropertyAttributeRead | PropertyAttributeWrite
	PropertyAttributeReadAnno      = PropertyAttributeRead | PropertyAttributeAnno
	PropertyAttributeReadWriteAnno = PropertyAttributeRead | PropertyAttributeWrite | PropertyAttributeAnno
)

// PropertyCode is a type for property code.
type PropertyCode byte

// PropertyAttribute is a type for property attribute.
type PropertyAttribute uint

// Property is an instance for Echonet property.
type Property struct {
	Code byte
	Attr uint
	Data []byte
}

// NewProperty returns a new property.
func NewProperty() *Property {
	prop := &Property{
		Code: 0,
		Attr: PropertyAttributeNone,
		Data: make([]byte, 0),
	}
	return prop
}

// SetCode sets a code to the property
func (prop *Property) SetCode(code byte) {
	prop.Code = code
}

// GetCode returns the property code.
func (prop *Property) GetCode() byte {
	return prop.Code
}

// SetAttribute sets an attribute to the property
func (prop *Property) SetAttribute(attr uint) {
	prop.Attr = attr
}

// GetAttribute returns the property attribute
func (prop *Property) GetAttribute() uint {
	return prop.Attr
}

// SetData sets a code to the property
func (prop *Property) SetData(data []byte) {
	prop.Data = make([]byte, len(data))
	copy(prop.Data, data)
}

// GetData returns the property data.
func (prop *Property) GetData() []byte {
	return prop.Data
}

// Size return the property data size.
func (prop *Property) Size() int {
	return len(prop.Data)
}

// Copy returns a copy property of the property.
func Copy(prop *Property) *Property {
	copyProp := &Property{
		Code: prop.Code,
		Attr: prop.Attr,
		Data: make([]byte, len(prop.Data)),
	}
	copy(copyProp.Data, prop.Data)
	return copyProp
}
