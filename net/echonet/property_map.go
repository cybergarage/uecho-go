// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
	"sort"
	"time"
)

const (
	propertyWaitRetryCount = 20
)

const (
	errorPropertyNotFound = "Property (%02X) Not Found"
)

// propertyMap represents a property map.
type propertyMap struct {
	properties   map[PropertyCode]Property
	parentObject Object
}

// newPropertyMap returns a new property map.
func newPropertyMap() *propertyMap {
	propMap := &propertyMap{
		properties:   map[PropertyCode]Property{},
		parentObject: nil,
	}

	return propMap
}

// SetObject sets a parent object.
func (propMap *propertyMap) SetObject(obj Object) {
	propMap.parentObject = obj
	for _, prop := range propMap.properties {
		prop.SetObject(obj)
	}
}

// Object returns a parent object.
func (propMap *propertyMap) Object() Object {
	return propMap.parentObject
}

// AddProperty adds a new property into the property map.
func (propMap *propertyMap) AddProperty(prop Property) {
	propMap.properties[prop.Code()] = prop
	prop.SetObject(propMap.parentObject)
}

// ClearAllProperties removes all properties in the property map.
func (propMap *propertyMap) ClearAllProperties(prop Property) {
	for code := range propMap.properties {
		delete(propMap.properties, code)
	}
}

// Properties returns the all properties in the property map.
func (propMap *propertyMap) Properties() []Property {
	codes := make([]PropertyCode, len(propMap.properties))
	n := 0
	for code := range propMap.properties {
		codes[n] = code
		n++
	}

	sort.Slice(codes, func(i, j int) bool { return codes[i] < codes[j] })

	props := []Property{}
	for _, code := range codes {
		prop, ok := propMap.properties[code]
		if !ok {
			continue
		}
		props = append(props, prop)
	}
	return props
}

// LookupProperty returns the specified property in the property map.
func (propMap *propertyMap) LookupProperty(code PropertyCode) (Property, bool) {
	prop, ok := propMap.properties[code]
	return prop, ok
}

// FindPropertyWait returns the specified property in the property map with the specified waiting time.
func (propMap *propertyMap) FindPropertyWait(code PropertyCode, waitTime time.Duration) (Property, bool) {
	for range propertyWaitRetryCount {
		time.Sleep(waitTime / propertyWaitRetryCount)
		prop, ok := propMap.LookupProperty(code)
		if ok {
			return prop, true
		}
	}
	return nil, false
}

// PropertyCount returns the property count in the property map.
func (propMap *propertyMap) PropertyCount() int {
	return len(propMap.properties)
}

// SetPropertyData sets a data to the existing property.
func (propMap *propertyMap) SetPropertyData(propCode PropertyCode, propData []byte) error {
	prop, ok := propMap.LookupProperty(propCode)
	if !ok {
		return fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	prop.SetData(propData)
	return nil
}

// SetPropertyByte sets a byte to the existing property.
func (propMap *propertyMap) SetPropertyByte(propCode PropertyCode, propData byte) error {
	prop, ok := propMap.LookupProperty(propCode)
	if !ok {
		return fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	prop.SetData([]byte{propData})
	return nil
}

// SetPropertyInteger sets a integer to the existing property.
func (propMap *propertyMap) SetPropertyInteger(propCode PropertyCode, propData uint, propSize uint) error {
	prop, ok := propMap.LookupProperty(propCode)
	if !ok {
		return fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	prop.SetInteger(propData, propSize)
	return nil
}

// HasProperty return true when the specified property exists, otherwise false.
func (propMap *propertyMap) HasProperty(propCode PropertyCode) bool {
	_, ok := propMap.LookupProperty(propCode)
	return ok
}

// LookupPropertyDataSize return the specified property data size in the property map.
func (propMap *propertyMap) LookupPropertyDataSize(propCode PropertyCode) (int, error) {
	prop, ok := propMap.LookupProperty(propCode)
	if !ok {
		return -1, fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	return len(prop.Data()), nil
}

// LookupPropertyData return the specified property data in the property map.
func (propMap *propertyMap) LookupPropertyData(propCode PropertyCode) ([]byte, error) {
	prop, ok := propMap.LookupProperty(propCode)
	if !ok {
		return nil, fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	return prop.Data(), nil
}

// LookupPropertyByte return the specified property byte data in the property map.
func (propMap *propertyMap) LookupPropertyByte(propCode PropertyCode) (byte, error) {
	prop, ok := propMap.LookupProperty(propCode)
	if !ok {
		return 0, fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	return prop.ByteData()
}

// LookupPropertyInteger return the specified property integer data in the property map.
func (propMap *propertyMap) LookupPropertyInteger(propCode PropertyCode) (uint, error) {
	prop, ok := propMap.LookupProperty(propCode)
	if !ok {
		return 0, fmt.Errorf(errorPropertyNotFound, uint(propCode))
	}
	return prop.IntegerData()
}
