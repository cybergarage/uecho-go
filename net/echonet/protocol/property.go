// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

// PropertyCode is a type for property code.
type PropertyCode byte

type Property interface {
	SetCode(code PropertyCode)
	Code() PropertyCode
	SetData(data []byte)
	Data() []byte
	StringData() string
	IntegerData() uint
	Size() int
}

// property is an instance for Echonet property.
type property struct {
	code PropertyCode
	data []byte
}

func newProperty() *property {
	prop := &property{
		code: 0,
		data: make([]byte, 0),
	}
	return prop
}

// NewProperty returns a new property.
func NewProperty() Property {
	return newProperty()
}

// NewPropertyWithCode returns a new property with the specified code.
func NewPropertyWithCode(code PropertyCode) Property {
	prop := newProperty()
	prop.SetCode(code)
	return prop
}

// NewPropertiesWithCodes returns a new properties with the specified codes.
func NewPropertiesWithCodes(codes []PropertyCode) []Property {
	props := make([]Property, len(codes))
	for n, code := range codes {
		props[n] = NewPropertyWithCode(code)
	}
	return props
}

// SetCode sets a code to the property.
func (prop *property) SetCode(code PropertyCode) {
	prop.code = code
}

// Code returns the property code.
func (prop *property) Code() PropertyCode {
	return prop.code
}

// SetData sets a code to the property.
func (prop *property) SetData(data []byte) {
	prop.data = make([]byte, len(data))
	copy(prop.data, data)
}

// Data returns the property data.
func (prop *property) Data() []byte {
	return prop.data
}

// StringData returns the property string data.
func (prop *property) StringData() string {
	return string(prop.data)
}

// IntegerData returns the property integer data.
func (prop *property) IntegerData() uint {
	return encoding.ByteToInteger(prop.data)
}

// Size return the property data size.
func (prop *property) Size() int {
	return len(prop.data)
}
