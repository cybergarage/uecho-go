// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

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

// NewDevice returns a new device Object.
func NewSuperObject() *SuperObject {
	obj := &SuperObject{
		Object: NewObject(),
	}
	return obj
}

// CreateProperty creates a new property to the property map. (Override function for PropertyMap)
func (obj *SuperObject) CreateProperty(propCode PropertyCode, propAttr PropertyAttribute) {
	obj.PropertyMap.CreateProperty(propCode, propAttr)
	obj.updatePropertyMap()
}

// setPropertyMapProperty sets a specified property map to the object.
func (obj *SuperObject) setPropertyMapProperty(propMapCode PropertyCode, propCodes []PropertyCode) error {
	if !obj.HasProperty(propMapCode) {
		obj.PropertyMap.CreateProperty(propMapCode, PropertyAttributeRead)
	}

	// Description Format 1

	if len(propCodes) <= PropertyMapFormat1MaxSize {
		propMapData := make([]byte, len(propCodes)+1)
		propMapData[0] = byte(len(propCodes))
		for n, propCode := range propCodes {
			propMapData[n+1] = byte(propCode)
		}
		obj.SetPropertyData(propMapCode, propMapData)
		return nil
	}

	// Description Format 2

	propMapData := make([]byte, PropertyMapFormat2Size)
	propMapData[0] = byte(len(propCodes))

	for _, propCode := range propCodes {
		if (propCode < PropertyCodeMin) || (PropertyCodeMax < propCode) {
			continue
		}
		propByteIdx := ((propCode - PropertyCodeMin) & 0x0F) + 1
		propMapData[propByteIdx] |= byte(((int(propCode-PropertyCodeMin) & 0xF0) >> 8) & 0x0F)
	}

	return nil
}

// updatePropertyMaps updates property maps  in the object.
func (obj *SuperObject) updatePropertyMap() error {
	getPropMapCodes := make([]PropertyCode, 0)
	setPropMapCodes := make([]PropertyCode, 0)
	annoPropMapCodes := make([]PropertyCode, 0)

	for _, prop := range obj.properties {
		if prop.IsReadable() {
			getPropMapCodes = append(getPropMapCodes, prop.GetCode())
		}
		if prop.IsWritable() {
			setPropMapCodes = append(setPropMapCodes, prop.GetCode())
		}
		if prop.IsAnnouncement() {
			annoPropMapCodes = append(annoPropMapCodes, prop.GetCode())
		}
	}

	obj.setPropertyMapProperty(ObjectGetPropertyMap, getPropMapCodes)
	obj.setPropertyMapProperty(ObjectSetPropertyMap, setPropMapCodes)
	obj.setPropertyMapProperty(ObjectAnnoPropertyMap, annoPropMapCodes)

	return nil
}

// SetOperatingStatus sets a operating status to the object.
func (obj *SuperObject) SetOperatingStatus(stats bool) error {
	statsByte := byte(ObjectOperatingStatusOff)
	if stats {
		statsByte = ObjectOperatingStatusOn
	}
	return obj.SetPropertyByteData(ObjectOperatingStatus, statsByte)
}

// GetOperatingStatus return the operating status of the object.
func (obj *SuperObject) GetOperatingStatus() (bool, error) {
	statsByte, err := obj.GetPropertyByteData(ObjectOperatingStatus)
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

// GetManufacturerCode return the manufacture codes of the object.
func (obj *SuperObject) GetManufacturerCode() (uint, error) {
	return obj.GetPropertyIntegerData(ObjectManufacturerCode)
}
