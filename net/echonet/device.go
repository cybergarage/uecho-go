// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
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

// DeviceOption is a function that applies a configuration to a device.
type DeviceOption func(*device) error

// Device represents an instance for a device object of Echonet.
type Device interface {
	SuperObject
}

type device struct {
	SuperObject
}

// WithDeviceCode returns a DeviceOption that sets the object code for a device.
// It also adds standard properties using the StandardDatabase.
// Returns an error if the object code is not found in the database.
func WithDeviceCode(code ObjectCode) DeviceOption {
	return func(dev *device) error {
		_, ok := SharedStandardDatabase().LookupObjectByCode(code)
		if !ok {
			return fmt.Errorf("object code (%X) not found in standard database", code)
		}
		dev.SetCode(code)
		return nil
	}
}

// WithDeviceManufacturerCode sets a manufacture codes to the device.
func WithDeviceManufacturerCode(code uint) DeviceOption {
	return func(dev *device) error {
		return dev.SetManufacturerCode(code)
	}
}

// WithDeviceRequestHandler returns a DeviceOption that sets the request handler for a device.
func WithDeviceRequestHandler(handler ObjectRequestHandler) DeviceOption {
	return func(dev *device) error {
		dev.SetRequestHandler(handler)
		return nil
	}
}

// NewDevice returns a new device with the specified options.
func NewDevice(opts ...DeviceOption) (Device, error) {
	dev := newDevice()
	for _, opt := range opts {
		if err := opt(dev); err != nil {
			return nil, err
		}
	}
	return dev, nil
}

func newDevice() *device {
	dev := &device{
		SuperObject: NewSuperObject(),
	}
	return dev
}

// NewDeviceWithCode returns a new device of the specified object code.
// It also adds standard properties using the StandardDatabase.
// Returns an error if the object code is not found in the database.
func NewDeviceWithCode(code ObjectCode) (Device, error) {
	obj := newDevice()
	if err := WithDeviceCode(code)(obj); err != nil {
		return nil, err
	}
	return obj, nil
}

// SetInstallationLocation sets a installation location to the device.
func (dev *device) SetInstallationLocation(locByte byte) error {
	return dev.SetPropertyByte(DeviceInstallationLocation, locByte)
}

// InstallationLocation return the installation location of the device.
func (dev *device) InstallationLocation() (byte, error) {
	return dev.LookupPropertyByte(DeviceInstallationLocation)
}

// SetStandardVersion sets a standard version to the device.
func (dev *device) SetStandardVersion(ver byte) error {
	verBytes := []byte{0x00, 0x00, ver, 0x00}
	return dev.SetPropertyData(DeviceStandardVersion, verBytes)
}

// StandardVersion return the standard version of the device.
func (dev *device) StandardVersion() (byte, error) {
	verBytes, err := dev.LookupPropertyData(DeviceStandardVersion)
	if err != nil {
		return 0, err
	}
	if len(verBytes) <= DeviceStandardVersionSize {
		return 0, fmt.Errorf(errorDeviceInvalidDeviceStandardVersion, string(verBytes))
	}
	return verBytes[2], nil
}

// SetFaultStatus sets a fault status to the device.
func (dev *device) SetFaultStatus(stats bool) error {
	statsByte := byte(DeviceNoFaultOccurred)
	if stats {
		statsByte = DeviceFaultOccurred
	}
	return dev.SetPropertyByte(DeviceFaultStatus, statsByte)
}

// FaultStatus return the fault status of the device.
func (dev *device) FaultStatus() (bool, error) {
	statsByte, err := dev.LookupPropertyByte(DeviceFaultStatus)
	if err != nil {
		return false, err
	}
	if statsByte == DeviceFaultOccurred {
		return true, nil
	}
	return false, nil
}
