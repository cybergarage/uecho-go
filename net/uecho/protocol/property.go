// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

const (
	PropertyCodeMin = 0x80
	PropertyCodeMax = 0xFF

	PropertyAttrNone      = 0x00
	PropertyAttrRead      = 0x01
	PropertyAttrWrite     = 0x02
	PropertyAttrAnno      = 0x10
	PropertyAttrReadWrite = PropertyAttrRead | PropertyAttrWrite
	PropertyAttrReadAnno  = PropertyAttrRead | PropertyAttrAnno
)

// Property is an instance for Echonet property.
type Property struct {
	Code byte
	Attr int
	Data []byte
}

// NewProperty returns a new property.
func NewProperty() *Property {
	prop := &Property{
		Code: 0,
		Attr: PropertyAttrNone,
		Data: make([]byte, 0),
	}
	return prop
}

// SetCode sets a code
func (prop *Property) SetCode(code byte) {
	prop.Code = code
}

// GetCode returns the code.
func (prop *Property) GetCode() byte {
	return prop.Code
}

// SetData sets a code
func (prop *Property) SetData(data []byte) {
	prop.Data = data
}

// GetData returns the data.
func (prop *Property) GetData() []byte {
	return prop.Data
}

// Size return the data size.
func (prop *Property) Size() int {
	return len(prop.Data)
}
