// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

// PropertyCode is a type for property code.
type PropertyCode byte

// Property is an interface for Echonet property.
type Property interface {
	Code() PropertyCode
	Data() []byte
	Size() int
	// PropertyHelper is an interface for Echonet property helper functions.
	PropertyHelper
	// PropertyMutator is an interface for Echonet property mutator functions.
	PropertyMutator
}

// PropertyHelper is an interface for Echonet property helper functions.
type PropertyHelper interface {
	// AsByte returns a byte value of the property byte data.
	AsByte() (byte, error)
	// AsString returns the property string data.
	AsString() (string, error)
	// AsInteger returns the property integer data.
	AsInteger() (uint, error)
}

// PropertyMutator is an interface for Echonet property mutator functions.
type PropertyMutator interface {
	// SetCode sets a code to the property.
	SetCode(code PropertyCode)
	// SetData sets a code to the property.
	SetData(data []byte)
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

// AsByte returns a byte value of the property byte data.
func (prop *property) AsByte() (byte, error) {
	switch len(prop.data) {
	case 0:
		return 0, ErrNoData
	case 1:
		return prop.data[0], nil
		// ok
	default:
		return 0, fmt.Errorf("%w data size (%d)", ErrInvalid, len(prop.data))
	}
}

// AsString returns the property string data.
func (prop *property) AsString() (string, error) {
	if len(prop.data) == 0 {
		return "", ErrNoData
	}
	return string(prop.data), nil
}

// AsInteger returns the property integer data.
func (prop *property) AsInteger() (uint, error) {
	if len(prop.data) == 0 {
		return 0, ErrNoData
	}
	return encoding.ByteToInteger(prop.data), nil
}

// Size return the property data size.
func (prop *property) Size() int {
	return len(prop.data)
}
