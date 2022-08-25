// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import "fmt"

const (
	SuperObjectCode = ObjectCode(0x000000)
)

const (
	errorPropertyMapNotFound = "property map (%2X) not found"
)

const (
	ObjectOperatingStatus  = 0x80
	ObjectManufacturerCode = 0x8A
	ObjectAnnoPropertyMap  = 0x9D
	ObjectSetPropertyMap   = 0x9E
	ObjectGetPropertyMap   = 0x9F
)

const (
	ObjectOperatingStatusOn             = 0x30
	ObjectOperatingStatusOff            = 0x31
	ObjectOperatingStatusSize           = 1
	ObjectManufacturerEvaluationCodeMin = 0xFFFFF0
	ObjectManufacturerEvaluationCodeMax = 0xFFFFFF
	ObjectManufacturerCodeSize          = 3
	ObjectPropertyMapMaxSize            = 16
	ObjectAnnoPropertyMapMaxSize        = (ObjectPropertyMapMaxSize + 1)
	ObjectSetPropertyMapMaxSize         = (ObjectPropertyMapMaxSize + 1)
	ObjectGetPropertyMapMaxSize         = (ObjectPropertyMapMaxSize + 1)
)

const (
	ObjectManufacturerUnknown = ObjectManufacturerEvaluationCodeMin
)

// SuperObject represents a super object of Echonet device and profile objects.
type SuperObject struct {
	*Object
}

// NewSuperObject returns a new device Object.
func NewSuperObject() *SuperObject {
	obj := &SuperObject{
		Object: NewObject(),
	}
	obj.SetCode(SuperObjectCode)
	obj.updatePropertyMap()
	return obj
}

// SetCode sets a code to the object.
func (obj *SuperObject) SetCode(code ObjectCode) {
	obj.Object.SetCode(code)
	obj.addStandardProperties()
}

// SetCodes sets codes to the object.
func (obj *SuperObject) SetCodes(codes []byte) {
	obj.Object.SetCodes(codes)
	obj.addStandardProperties()
}

// addStandardProperties sets mandatory properties of the object code.
func (obj *SuperObject) addStandardProperties() {
	stdObj, ok := GetStandardDatabase().FindObjectByCode(obj.Code())
	if !ok {
		return
	}
	obj.SetClassName(stdObj.ClassName())
	for _, stdProp := range stdObj.Properties() {
		obj.AddProperty(stdProp.Copy())
	}
}

// AddProperty adds a new property into the property map.
func (obj *SuperObject) AddProperty(prop *Property) {
	obj.Object.AddProperty(prop)
	obj.updatePropertyMap()
}

// setPropertyMapProperty sets a specified property map to the object.
func (obj *SuperObject) setPropertyMapProperty(propMapCode PropertyCode, propCodes []PropertyCode) error {
	if !obj.HasProperty(propMapCode) {
		return fmt.Errorf(errorPropertyMapNotFound, propMapCode)
	}

	// Description Format 1

	if isPropertyMapDescriptionFormat2(len(propCodes)) {
		propMapData := make([]byte, len(propCodes)+1)
		propMapData[0] = byte(len(propCodes))
		for n, propCode := range propCodes {
			propMapData[n+1] = byte(propCode)
		}
		return obj.SetPropertyData(propMapCode, propMapData)
	}

	// Description Format 2

	propMapData := make([]byte, (PropertyMapFormat2MapSize + 1))
	propMapData[0] = byte(len(propCodes))

	for _, propCode := range propCodes {
		propCodeIdx, propCodeBit, ok := propertyMapCodeToFormat2(propCode)
		if !ok {
			continue
		}
		propMapData[propCodeIdx] |= byte((0x01 << propCodeBit) & 0x0F)
	}

	return obj.SetPropertyData(propMapCode, propMapData)
}

// updatePropertyMaps updates property maps  in the object.
func (obj *SuperObject) updatePropertyMap() error {
	propMaps := []struct {
		code  PropertyCode
		codes []PropertyCode
	}{
		{code: ObjectGetPropertyMap, codes: make([]PropertyCode, 0)},
		{code: ObjectSetPropertyMap, codes: make([]PropertyCode, 0)},
		{code: ObjectAnnoPropertyMap, codes: make([]PropertyCode, 0)},
	}

	for _, prop := range obj.Properties() {
		propCode := prop.Code()
		if prop.IsReadable() {
			propMaps[0].codes = append(propMaps[0].codes, propCode)
		}
		if prop.IsWritable() {
			propMaps[1].codes = append(propMaps[1].codes, propCode)
		}
		if prop.IsAnnounceable() {
			propMaps[2].codes = append(propMaps[2].codes, propCode)
		}
	}

	var lastErr error
	for _, propMap := range propMaps {
		if err := obj.setPropertyMapProperty(propMap.code, propMap.codes); err != nil {
			lastErr = err
		}
	}

	return lastErr
}

// SetOperatingStatus sets a operating status to the object.
func (obj *SuperObject) SetOperatingStatus(stats bool) error {
	statsByte := byte(ObjectOperatingStatusOff)
	if stats {
		statsByte = ObjectOperatingStatusOn
	}
	return obj.SetPropertyByteData(ObjectOperatingStatus, statsByte)
}

// OperatingStatus return the operating status of the object.
func (obj *SuperObject) OperatingStatus() (bool, error) {
	statsByte, err := obj.FindPropertyByteData(ObjectOperatingStatus)
	if err != nil {
		return false, err
	}
	if statsByte == ObjectOperatingStatusOff {
		return false, nil
	}
	return true, nil
}

// SetManufacturerCode sets a manufacture codes to the object.
func (obj *SuperObject) SetManufacturerCode(code uint) error {
	return obj.SetPropertyIntegerData(ObjectManufacturerCode, code, ObjectManufacturerCodeSize)
}

// ManufacturerCode return the manufacture codes of the object.
func (obj *SuperObject) ManufacturerCode() (uint, error) {
	return obj.FindPropertyIntegerData(ObjectManufacturerCode)
}
