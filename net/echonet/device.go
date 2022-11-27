// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
)

const (
	DeviceOperatingStatus                       = ObjectOperatingStatus
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
	DeviceOperatingStatusOn           = ObjectOperatingStatusOn
	DeviceOperatingStatusOff          = ObjectOperatingStatusOff
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
	DeviceManufacturerExperimental    = ObjectManufacturerExperimental
)

const (
	errorDeviceInvalidDeviceStandardVersion = "invalid device standard version (%s)"
)

// Device represents an instance for a device object of Echonet.
type Device struct {
	*SuperObject
}

// NewDevice returns a new device Object.
func NewDevice() *Device {
	dev := &Device{
		SuperObject: NewSuperObject(),
	}
	return dev
}

// NewDeviceWithCodes returns a new device of the specified object codes.
func NewDeviceWithCodes(codes []byte) (*Device, error) {
	objCode, err := BytesToObjectCode(codes)
	if err != nil {
		return nil, err
	}
	obj := NewDevice()
	obj.SetCode(objCode)
	return obj, nil
}

// NewDeviceWithCode returns a new device of the specified object code.
func NewDeviceWithCode(code ObjectCode) *Device {
	obj := NewDevice()
	obj.SetCode(code)
	return obj
}

// SetInstallationLocation sets a installation location to the device.
func (dev *Device) SetInstallationLocation(locByte byte) error {
	return dev.SetPropertyByteData(DeviceInstallationLocation, locByte)
}

// InstallationLocation return the installation location of the device.
func (dev *Device) InstallationLocation() (byte, error) {
	return dev.FindPropertyByteData(DeviceInstallationLocation)
}

// SetStandardVersion sets a standard version to the device.
func (dev *Device) SetStandardVersion(ver byte) error {
	verBytes := []byte{0x00, 0x00, ver, 0x00}
	return dev.SetPropertyData(DeviceStandardVersion, verBytes)
}

// StandardVersion return the standard version of the device.
func (dev *Device) StandardVersion() (byte, error) {
	verBytes, err := dev.FindPropertyData(DeviceStandardVersion)
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

// FaultStatus return the fault status of the device.
func (dev *Device) FaultStatus() (bool, error) {
	statsByte, err := dev.FindPropertyByteData(DeviceFaultStatus)
	if err != nil {
		return false, err
	}
	if statsByte == DeviceFaultOccurred {
		return true, nil
	}
	return false, nil
}
