// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

const (
	PropertyCodeMin = 0x80
	PropertyCodeMax = 0xFF

	PropertyAttributeNone      = 0x00
	PropertyAttributeRead      = 0x01
	PropertyAttributeWrite     = 0x02
	PropertyAttributeAnno      = 0x10
	PropertyAttributeReadWrite = PropertyAttributeRead | PropertyAttributeWrite
	PropertyAttributeReadAnno  = PropertyAttributeRead | PropertyAttributeAnno
)

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
	prop.Data = data
}

// GetData returns the property data.
func (prop *Property) GetData() []byte {
	return prop.Data
}

// Size return the property data size.
func (prop *Property) Size() int {
	return len(prop.Data)
}
