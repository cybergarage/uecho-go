// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"fmt"
)

/****************************************
 * Device Object Super Class
 ****************************************/

const (
	DeviceOperatingStatus                       = 0x80
	DeviceInstallationLocation                  = 0x81
	DeviceStandardVersion                       = 0x82
	DeviceIdentificationNumber                  = 0x83
	DeviceMeasuredInstantaneousPowerConsumption = 0x84
	DeviceMeasuredCumulativePowerConsumption    = 0x85
	DeviceManufacturerFaultCode                 = 0x86
	DeviceCurrentLimitSetting                   = 0x87
	DeviceFaultStatus                           = 0x88
	DeviceFaultDescription                      = 0x89
	DeviceManufacturerCode                      = ObjectManufacturerCode
	DeviceBusinessFacilityCode                  = 0x8B
	DeviceProductCode                           = 0x8C
	DeviceProductionNumber                      = 0x8D
	DeviceProductionDate                        = 0x8E
	DevicePowerSavingOperationSetting           = 0x8F
	DeviceRemoteControlSetting                  = 0x93
	DeviceCurrentTimeSetting                    = 0x97
	DeviceCurrentDateSetting                    = 0x98
	DevicePowerLimitSetting                     = 0x99
	DeviceCumulativeOperatingTime               = 0x9A
	DeviceAnnoPropertyMap                       = ObjectAnnoPropertyMap
	DeviceSetPropertyMap                        = ObjectSetPropertyMap
	DeviceGetPropertyMap                        = ObjectGetPropertyMap
)

const (
	DeviceOperatingStatusSize                       = 1
	DeviceInstallationLocationSize                  = 1
	DeviceStandardVersionSize                       = 4
	DeviceIdentificationNumberSize                  = 9
	DeviceMeasuredInstantaneousPowerConsumptionSize = 2
	DeviceMeasuredCumulativePowerConsumptionSize    = 4
	DeviceManufacturerFaultCodeSize                 = 255
	DeviceCurrentLimitSettingSize                   = 1
	DeviceFaultStatusSize                           = 1
	DeviceFaultDescriptionSize                      = 2
	DeviceManufacturerCodeSize                      = ObjectManufacturerCodeSize
	DeviceBusinessFacilityCodeSize                  = 3
	DeviceProductCodeSize                           = 12
	DeviceProductionNumberSize                      = 12
	DeviceProductionDateSize                        = 4
	DevicePowerSavingOperationSettingSize           = 1
	DeviceRemoteControlSettingSize                  = 1
	DeviceCurrentTimeSettingSize                    = 2
	DeviceCurrentDateSettingSize                    = 4
	DevicePowerLimitSettingSize                     = 2
	DeviceCumulativeOperatingTimeSize               = 5
	DeviceAnnoPropertyMapSize                       = ObjectAnnoPropertyMapMaxSize
	DeviceSetPropertyMapSize                        = ObjectSetPropertyMapMaxSize
	DeviceGetPropertyMapSize                        = ObjectGetPropertyMapMaxSize
)

const (
	DeviceOperatingStatusOn           = 0x30
	DeviceOperatingStatusOff          = 0x31
	DeviceVersionAppendixA            = 'A'
	DeviceVersionAppendixB            = 'B'
	DeviceVersionAppendixC            = 'C'
	DeviceVersionAppendixD            = 'D'
	DeviceVersionAppendixE            = 'E'
	DeviceVersionAppendixF            = 'F'
	DeviceVersionAppendixG            = 'G'
	DeviceVersionUnknown              = 0
	DeviceDefaultVersionAppendix      = DeviceVersionAppendixG
	DeviceFaultOccurred               = 0x41
	DeviceNoFaultOccurred             = 0x42
	DeviceInstallationLocationUnknown = 0x00
	DeviceManufacturerUnknown         = 0xFFFFFF
)

const (
	errorDeviceInvalidDeviceStandardVersion = "Invalid Device Standard Version (%s)"
)

// Device is an instance for Echonet device Object.
type Device = Object

// NewDevice returns a new device Object.
func NewDevice() *Device {
	dev := NewObject()
	dev.addMandatoryProperties()
	return dev
}

// CreateProperty creates a new property to the property map. (Override)
func (dev *Device) CreateProperty(propCode PropertyCode, propAttr PropertyAttribute) {
	dev.PropertyMap.CreateProperty(propCode, propAttr)
	dev.updatePropertyMap()
}

// addMandatoryProperties sets mandatory properties for Echonet specification
func (dev *Device) addMandatoryProperties() error {
	// Operation Status
	dev.CreateProperty(DeviceOperatingStatus, PropertyAttributeReadAnno)
	dev.SetOperatingStatus(true)

	// Installation Location
	dev.CreateProperty(DeviceInstallationLocation, PropertyAttributeReadAnno)
	dev.SetInstallationLocation(DeviceInstallationLocationUnknown)

	// Standard Version Infomation
	dev.CreateProperty(DeviceStandardVersion, PropertyAttributeRead)
	dev.SetStandardVersion(DeviceDefaultVersionAppendix)

	// Fault Status
	dev.CreateProperty(DeviceFaultStatus, PropertyAttributeReadAnno)
	dev.SetFaultStatus(false)

	// Manufacture Code
	dev.CreateProperty(DeviceManufacturerCode, PropertyAttributeRead)
	dev.SetManufacturerCode(DeviceManufacturerUnknown)

	return nil
}

// SetOperatingStatus sets a operating status to the device.
func (dev *Device) SetOperatingStatus(stats bool) error {
	statsByte := byte(DeviceOperatingStatusOff)
	if stats {
		statsByte = DeviceOperatingStatusOn
	}
	return dev.SetPropertyByteData(DeviceOperatingStatus, statsByte)
}

// GetOperatingStatus return the operating status of the device.
func (dev *Device) GetOperatingStatus() (bool, error) {
	statsByte, err := dev.GetPropertyByteData(DeviceOperatingStatus)
	if err != nil {
		return false, err
	}
	if statsByte == DeviceOperatingStatusOff {
		return false, nil
	}
	return true, nil
}

// SetInstallationLocation sets a installation location to the device.
func (dev *Device) SetInstallationLocation(locByte byte) error {
	return dev.SetPropertyByteData(DeviceInstallationLocation, locByte)
}

// GetInstallationLocation return the installation location of the device.
func (dev *Device) GetInstallationLocation() (byte, error) {
	return dev.GetPropertyByteData(DeviceInstallationLocation)
}

// SetStandardVersion sets a standard version to the device.
func (dev *Device) SetStandardVersion(ver byte) error {
	verBytes := []byte{0x00, 0x00, ver, 0x00}
	return dev.SetPropertyData(DeviceStandardVersion, verBytes)
}

// GetStandardVersion return the standard version of the device.
func (dev *Device) GetStandardVersion() (byte, error) {
	verBytes, err := dev.GetPropertyData(DeviceStandardVersion)
	if err != nil {
		return 0, err
	}
	if len(verBytes) <= DeviceStandardVersionSize {
		return 0, fmt.Errorf(errorDeviceInvalidDeviceStandardVersion, string(verBytes))
	}
	return verBytes[2], nil
}

// SetFaultStatus sets a fault status to the device.
func (dev *Device) SetFaultStatus(stats bool) error {
	statsByte := byte(DeviceNoFaultOccurred)
	if stats {
		statsByte = DeviceFaultOccurred
	}
	return dev.SetPropertyByteData(DeviceFaultStatus, statsByte)
}

// GetFaultStatus return the fault status of the device.
func (dev *Device) GetFaultStatus() (bool, error) {
	statsByte, err := dev.GetPropertyByteData(DeviceFaultStatus)
	if err != nil {
		return false, err
	}
	if statsByte == DeviceFaultOccurred {
		return true, nil
	}
	return false, nil
}

// SetManufacturerCode sets a manufacture codes to the device.
func (dev *Device) SetManufacturerCode(code uint) error {
	return dev.SetPropertyIntegerData(DeviceManufacturerCode, code, DeviceManufacturerCodeSize)
}

// GetManufacturerCode return the manufacture codes of the device.
func (dev *Device) GetManufacturerCode() (uint, error) {
	return dev.GetPropertyIntegerData(DeviceManufacturerCode)
}

// setPropertyMapProperty sets a specified property map to the device.
func (dev *Device) setPropertyMapProperty(propMapCode PropertyCode, propCodes []PropertyCode) error {
	if !dev.HasProperty(propMapCode) {
		dev.PropertyMap.CreateProperty(propMapCode, PropertyAttributeRead)
	}

	// Description Format 1

	if len(propCodes) <= PropertyMapFormat1MaxSize {
		propMapData := make([]byte, len(propCodes)+1)
		propMapData[0] = byte(len(propCodes))
		for n, propCode := range propCodes {
			propMapData[n+1] = byte(propCode)
		}
		dev.SetPropertyData(propMapCode, propMapData)
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

// updatePropertyMaps updates property maps  in the device.
func (dev *Device) updatePropertyMap() error {
	getPropMapCodes := make([]PropertyCode, 0)
	setPropMapCodes := make([]PropertyCode, 0)
	annoPropMapCodes := make([]PropertyCode, 0)

	for _, prop := range dev.properties {
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

	dev.setPropertyMapProperty(ProfileGetPropertyMap, getPropMapCodes)
	dev.setPropertyMapProperty(ProfileSetPropertyMap, setPropMapCodes)
	dev.setPropertyMapProperty(ProfileAnnoPropertyMap, annoPropMapCodes)

	return nil
}
