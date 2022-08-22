// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
	"time"
)

const (
	propertyWaitRetryCount = 20
)

const (
	errorPropertyNotFound = "Property (%02X) Not Found"
)

// PropertyMap represents a property map.
type PropertyMap struct {
	properties   map[PropertyCode]*Property
	parentObject *Object
}

// NewPropertyMap returns a new property map.
func NewPropertyMap() *PropertyMap {
	propMap := &PropertyMap{
		properties:   map[PropertyCode]*Property{},
		parentObject: nil,
	}

	return propMap
}

// SetParentObject sets a parent object.
func (propMap *PropertyMap) SetParentObject(obj *Object) {
	propMap.parentObject = obj
	for _, prop := range propMap.properties {
		prop.SetParentObject(obj)
	}
}

// ParentObject returns a parent object.
func (propMap *PropertyMap) ParentObject() *Object {
	return propMap.parentObject
}

// AddProperty adds a new property into the property map.
func (propMap *PropertyMap) AddProperty(prop *Property) {
	propMap.properties[prop.Code()] = prop
	prop.SetParentObject(propMap.parentObject)
}

// CreateProperty creates a new property to the property map.
func (propMap *PropertyMap) CreateProperty(propCode PropertyCode, propAttr PropertyAttr) {
	prop := NewProperty()
	prop.SetCode(propCode)
	prop.SetAttribute(propAttr)
	prop.SetParentObject(propMap.parentObject)
	propMap.AddProperty(prop)
}

// ClearAllProperties removes all properties in the property map.
func (propMap *PropertyMap) ClearAllProperties(prop *Property) {
	for code := range propMap.properties {
		delete(propMap.properties, code)
	}
}

// Properties returns the all properties in the property map.
func (propMap *PropertyMap) Properties() []*Property {
	props := []*Property{}
	for _, prop := range propMap.properties {
		props = append(props, prop)
	}
	return props
}

// FindProperty returns the specified property in the property map.
func (propMap *PropertyMap) FindProperty(code PropertyCode) (*Property, bool) {
	prop, ok := propMap.properties[code]
	return prop, ok
}

// FindPropertyWait returns the specified property in the property map with the specified waiting time.
func (propMap *PropertyMap) FindPropertyWait(code PropertyCode, waitTime time.Duration) (*Property, bool) {
	for n := 0; n < propertyWaitRetryCount; n++ {
		time.Sleep(waitTime / propertyWaitRetryCount)
		prop, ok := propMap.FindProperty(code)
		if ok {
			return prop, true
		}
	}
	return nil, false
}

// PropertyCount returns the property count in the property map.
func (propMap *PropertyMap) PropertyCount() int {
	return len(propMap.properties)
}

// SetPropertyAttribute sets an attribute to the existing property.
func (propMap *PropertyMap) SetPropertyAttribute(propCode PropertyCode, propAttr PropertyAttr) error {
	prop, ok := propMap.FindProperty(propCode)
	if !ok {
		return fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	prop.SetAttribute(propAttr)
	return nil
}

// FindPropertyAttribute returns the specified property attribute in the property map.
func (propMap *PropertyMap) FindPropertyAttribute(propCode PropertyCode) (PropertyAttr, error) {
	prop, ok := propMap.FindProperty(propCode)
	if !ok {
		return PropertyAttrNone, fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	return prop.Attribute(), nil
}

// SetPropertyData sets a data to the existing property.
func (propMap *PropertyMap) SetPropertyData(propCode PropertyCode, propData []byte) error {
	prop, ok := propMap.FindProperty(propCode)
	if !ok {
		return fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	prop.SetData(propData)
	return nil
}

// SetPropertyByteData sets a byte to the existing property.
func (propMap *PropertyMap) SetPropertyByteData(propCode PropertyCode, propData byte) error {
	prop, ok := propMap.FindProperty(propCode)
	if !ok {
		return fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	prop.SetData([]byte{propData})
	return nil
}

// SetPropertyIntegerData sets a integer to the existing property.
func (propMap *PropertyMap) SetPropertyIntegerData(propCode PropertyCode, propData uint, propSize uint) error {
	prop, ok := propMap.FindProperty(propCode)
	if !ok {
		return fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	prop.SetIntegerData(propData, propSize)
	return nil
}

// HasProperty return true when the specified property exists, otherwise false.
func (propMap *PropertyMap) HasProperty(propCode PropertyCode) bool {
	_, ok := propMap.FindProperty(propCode)
	return ok
}

// FindPropertyDataSize return the specified property data size in the property map.
func (propMap *PropertyMap) FindPropertyDataSize(propCode PropertyCode) (int, error) {
	prop, ok := propMap.FindProperty(propCode)
	if !ok {
		return -1, fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	return len(prop.Data()), nil
}

// FindPropertyData return the specified property data in the property map.
func (propMap *PropertyMap) FindPropertyData(propCode PropertyCode) ([]byte, error) {
	prop, ok := propMap.FindProperty(propCode)
	if !ok {
		return nil, fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	return prop.Data(), nil
}

// FindPropertyByteData return the specified property byte data in the property map.
func (propMap *PropertyMap) FindPropertyByteData(propCode PropertyCode) (byte, error) {
	prop, ok := propMap.FindProperty(propCode)
	if !ok {
		return 0, fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	return prop.ByteData()
}

// FindPropertyIntegerData return the specified property integer data in the property map.
func (propMap *PropertyMap) FindPropertyIntegerData(propCode PropertyCode) (uint, error) {
	prop, ok := propMap.FindProperty(propCode)
	if !ok {
		return 0, fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	return prop.IntegerData()
}
