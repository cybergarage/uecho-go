// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"bytes"
	"fmt"

	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

const (
	PropertyAttrNone      = 0x00
	PropertyAttrRead      = 0x01
	PropertyAttrWrite     = 0x02
	PropertyAttrAnno      = 0x10
	PropertyAttrReadWrite = PropertyAttrRead | PropertyAttrWrite
	PropertyAttrReadAnno  = PropertyAttrRead | PropertyAttrAnno
)

const (
	errorPropertyNoParentNode = "Property has no parent node"
)

// Property is an instance for Echonet property.
type Property struct {
	Code         byte
	Attr         int
	Data         []byte
	ParentObject *Object
}

// NewProperty returns a new property.
func NewProperty() *Property {
	prop := &Property{
		Code:         0,
		Attr:         PropertyAttrNone,
		Data:         make([]byte, 0),
		ParentObject: nil,
	}
	return prop
}

// SetParentObject sets a parent object into the property
func (prop *Property) SetParentObject(obj *Object) {
	prop.ParentObject = obj
}

// GetParentObject returns the parent object
func (prop *Property) GetParentObject() *Object {
	return prop.ParentObject
}

// SetCode sets a specified code to the property
func (prop *Property) SetCode(code byte) {
	prop.Code = code
}

// GetCode returns the property code
func (prop *Property) GetCode() byte {
	return prop.Code
}

// SetData sets a specified data to the property
func (prop *Property) SetData(data []byte) {
	prop.ClearData()
	prop.AddData(data)
}

// AddData adds a specified data to the property
func (prop *Property) AddData(data []byte) {
	if len(data) <= 0 {
		return
	}

	prop.Data = append(prop.Data, data...)

	// (D) Basic sequence for autonomous notification

	if prop.IsAnnouncement() {
		prop.Announce()
	}
}

// GetData returns the property data
func (prop *Property) GetData() []byte {
	return prop.Data
}

// ClearData clears the property data
func (prop *Property) ClearData() {
	prop.Data = make([]byte, 0)
}

// Size return the property data size.
func (prop *Property) Size() int {
	return len(prop.Data)
}

// SetAttribute sets an attribute to the property
func (prop *Property) SetAttribute(attr int) {
	prop.Attr = attr
}

// GetAttribute returns the property attribute
func (prop *Property) GetAttribute() int {
	return prop.Attr
}

// IsReadable returns true when the property attribute is readable, otherwise false
func (prop *Property) IsReadable() bool {
	if (prop.Attr & PropertyAttrRead) == 0 {
		return false
	}
	return true
}

// IsWritable returns true when the property attribute is writable, otherwise false
func (prop *Property) IsWritable() bool {
	if (prop.Attr & PropertyAttrWrite) == 0 {
		return false
	}
	return true
}

// IsReadOnly returns true when the property attribute is read only, otherwise false
func (prop *Property) IsReadOnly() bool {
	if (prop.Attr & PropertyAttrRead) == 0 {
		return false
	}

	if (prop.Attr & PropertyAttrWrite) != 0 {
		return false
	}

	return true
}

// IsWriteOnly returns true when the property attribute is write only, otherwise false
func (prop *Property) IsWriteOnly() bool {
	if (prop.Attr & PropertyAttrWrite) == 0 {
		return false
	}

	if (prop.Attr & PropertyAttrRead) != 0 {
		return false
	}

	return true
}

// IsAnnouncement returns true when the property attribute is announcement, otherwise false
func (prop *Property) IsAnnouncement() bool {
	if (prop.Attr & PropertyAttrAnno) == 0 {
		return false
	}
	return true
}

// toProtocolProperty returns the new property of the property
func (prop *Property) toProtocolProperty() *protocol.Property {
	newProp := protocol.NewProperty()
	newProp.SetCode(prop.GetCode())
	newProp.SetAttribute(prop.GetAttribute())
	newProp.SetData(prop.GetData())
	return newProp
}

// Equals returns true if the specified property is same, otherwise false
func (prop *Property) Equals(otherProp *Property) bool {
	if prop.GetCode() != otherProp.GetCode() {
		return false
	}
	if prop.GetAttribute() != otherProp.GetAttribute() {
		return false
	}
	if bytes.Compare(prop.GetData(), otherProp.GetData()) != 0 {
		return false
	}
	return true
}

// GetNode returns a parent node
func (prop *Property) GetNode() *Node {
	parentObj := prop.GetParentObject()
	if parentObj == nil {
		return nil
	}
	return parentObj.GetParentNode()
}

// Announce announces the property
func (prop *Property) Announce() error {
	parentNode := prop.GetNode()
	if parentNode == nil {
		return fmt.Errorf(errorPropertyNoParentNode)
	}
	return parentNode.AnnounceProperty(prop)
}
