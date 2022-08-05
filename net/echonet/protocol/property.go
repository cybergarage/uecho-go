// Copyright (C) 2018 Satoshi Konno. All rights reserved.
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
	Code PropertyCode
	Data []byte
}

// NewProperty returns a new property.
func NewProperty() *Property {
	return NewPropertyWithCode(0)
}

// NewPropertyWithCode returns a new property with the specified code.
func NewPropertyWithCode(code PropertyCode) *Property {
	prop := &Property{
		Code: code,
		Data: make([]byte, 0),
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
	prop.Code = code
}

// GetCode returns the property code.
func (prop *Property) GetCode() PropertyCode {
	return prop.Code
}

// SetData sets a code to the property.
func (prop *Property) SetData(data []byte) {
	prop.Data = make([]byte, len(data))
	copy(prop.Data, data)
}

// GetData returns the property data.
func (prop *Property) GetData() []byte {
	return prop.Data
}

// GetStringData returns the property string data.
func (prop *Property) GetStringData() string {
	return string(prop.Data)
}

// GetIntegerData returns the property integer data.
func (prop *Property) GetIntegerData() uint {
	return encoding.ByteToInteger(prop.Data)
}

// Size return the property data size.
func (prop *Property) Size() int {
	return len(prop.Data)
}

// Copy returns a copy property of the property.
func Copy(prop *Property) *Property {
	copyProp := &Property{
		Code: prop.Code,
		Data: make([]byte, len(prop.Data)),
	}
	copy(copyProp.Data, prop.Data)
	return copyProp
}
