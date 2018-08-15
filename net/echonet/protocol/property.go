// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

// Property is an instance for Echonet property.
type Property struct {
	Code byte
	Data []byte
}

// NewProperty returns a new property.
func NewProperty() *Property {
	prop := &Property{
		Code: 0,
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
		Data: make([]byte, len(prop.Data)),
	}
	copy(copyProp.Data, prop.Data)
	return copyProp
}
