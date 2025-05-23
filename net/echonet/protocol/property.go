// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

// PropertyCode is a type for property code.
type PropertyCode byte

// Property is an instance for Echonet property.
type Property struct {
	code PropertyCode
	data []byte
}

// NewProperty returns a new property.
func NewProperty() *Property {
	return NewPropertyWithCode(0)
}

// NewPropertyWithCode returns a new property with the specified code.
func NewPropertyWithCode(code PropertyCode) *Property {
	prop := &Property{
		code: code,
		data: make([]byte, 0),
	}
	return prop
}

// NewPropertiesWithCodes returns a new properties with the specified codes.
func NewPropertiesWithCodes(codes []PropertyCode) []*Property {
	props := make([]*Property, len(codes))
	for n, code := range codes {
		props[n] = NewPropertyWithCode(code)
	}
	return props
}

// SetCode sets a code to the property.
func (prop *Property) SetCode(code PropertyCode) {
	prop.code = code
}

// Code returns the property code.
func (prop *Property) Code() PropertyCode {
	return prop.code
}

// SetData sets a code to the property.
func (prop *Property) SetData(data []byte) {
	prop.data = make([]byte, len(data))
	copy(prop.data, data)
}

// Data returns the property data.
func (prop *Property) Data() []byte {
	return prop.data
}

// StringData returns the property string data.
func (prop *Property) StringData() string {
	return string(prop.data)
}

// IntegerData returns the property integer data.
func (prop *Property) IntegerData() uint {
	return encoding.ByteToInteger(prop.data)
}

// Size return the property data size.
func (prop *Property) Size() int {
	return len(prop.data)
}

// Copy returns a copy property of the property.
func Copy(prop *Property) *Property {
	copyProp := &Property{
		code: prop.code,
		data: make([]byte, len(prop.data)),
	}
	copy(copyProp.data, prop.data)
	return copyProp
}
