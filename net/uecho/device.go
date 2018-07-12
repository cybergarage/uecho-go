// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

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
	DeviceManufacturerCodeSize                      = ObjectManufacturerCodeLen
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
	DeviceAnnoPropertyMapSize                       = ObjectAnnoPropertyMapMaxLen
	DeviceSetPropertyMapSize                        = ObjectSetPropertyMapMaxLen
	DeviceGetPropertyMapSize                        = ObjectGetPropertyMapMaxLen
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
)

// Device is an instance for Echonet device object.
type Device = Object

// NewDevice returns a new device object.
func NewDevice() *Device {
	dev := NewObject()
	return dev
}

// AddMandatoryProperties sets mandatory properties for Echonet specification
/*
 func (dev *Device) AddMandatoryProperties(dev *Device) error {
   // Operation Status

   uecho_object_setproperty(obj, uEchoDeviceOperatingStatus, uEchoPropertyAttrReadAnno);
   (dev *Device) setoperatingstatus(obj, true);

   // Installation Location

   uecho_object_setproperty(obj, uEchoDeviceInstallationLocation, uEchoPropertyAttrReadAnno);
   (dev *Device) setinstallationlocation(obj, uEchoDeviceInstallationLocationUnknown);

   // Standard Version Infomation

   uecho_object_setproperty(obj, uEchoDeviceStandardVersion, uEchoPropertyAttrRead);
   (dev *Device) setstandardversion(obj, uEchoDeviceDefaultVersionAppendix);

   // Fault Status

   uecho_object_setproperty(obj, uEchoDeviceFaultStatus, uEchoPropertyAttrReadAnno);
   (dev *Device) setfaultstatus(obj, false);

   return true;
 }

 bool (dev *Device) setoperatingstatus(dev *Device, bool stats)
 {
   byte statsByte;

   statsByte = stats ? uEchoDeviceOperatingStatusOn : uEchoDeviceOperatingStatusOff;
   return uecho_object_setpropertydata(obj, uEchoDeviceOperatingStatus, &statsByte, uEchoDeviceOperatingStatusSize);
 }

 bool (dev *Device) isoperatingstatus(dev *Device)
 {
   byte statsByte;

   if (!uecho_object_getpropertybytedata(obj, uEchoDeviceOperatingStatus, &statsByte))
	 return false;

   return (statsByte == uEchoDeviceOperatingStatusOn) ? true : false;
 }

 bool (dev *Device) setinstallationlocation(dev *Device, byte locByte)
 {
   return uecho_object_setpropertydata(obj, uEchoDeviceInstallationLocation, &locByte, uEchoDeviceInstallationLocationSize);
 }

 byte (dev *Device) getinstallationlocation(dev *Device)
 {
   byte locByte;

   if (!uecho_object_getpropertybytedata(obj, uEchoDeviceInstallationLocation, &locByte))
	 return uEchoDeviceInstallationLocationUnknown;

   return locByte;
 }

 bool (dev *Device) setstandardversion(dev *Device, char ver)
 {
   byte verBytes[uEchoDeviceStandardVersionSize];

   verBytes[0] = 0x00;
   verBytes[1] = 0x00;
   verBytes[2] = ver;
   verBytes[3] = 0x00;
   return uecho_object_setpropertydata(obj, uEchoDeviceStandardVersion, verBytes, uEchoDeviceStandardVersionSize);
 }

 char (dev *Device) getstandardversion(dev *Device)
 {
   uEchoProperty *prop;
   byte *verBytes;

   prop = uecho_object_getproperty(obj, uEchoDeviceStandardVersion);
   if (!prop)
	 return uEchoDeviceVersionUnknown;

   if (uecho_property_getdatasize(prop) != uEchoDeviceStandardVersionSize)
	 return uEchoDeviceVersionUnknown;

   verBytes = uecho_property_getdata(prop);
   if (!verBytes)
	 return uEchoDeviceVersionUnknown;

   return verBytes[2];
 }

 bool (dev *Device) setfaultstatus(dev *Device, bool stats)
 {
   byte faultByte;

   faultByte = stats ? uEchoDeviceFaultOccurred : uEchoDeviceNoFaultOccurred;
   return uecho_object_setpropertydata(obj, uEchoDeviceFaultStatus, &faultByte, uEchoDeviceFaultStatusSize);
 }

 bool (dev *Device) isfaultstatus(dev *Device)
 {
   byte statsByte;

   if (!uecho_object_getpropertybytedata(obj, uEchoDeviceFaultStatus, &statsByte))
	 return false;

   return (statsByte == uEchoDeviceFaultOccurred) ? true : false;
 }

 bool (dev *Device) setmanufacturercode(dev *Device, uEchoManufacturerCode code)
 {
   return uecho_object_setmanufacturercode(obj, code);
 }

 uEchoManufacturerCode (dev *Device) getmanufacturercode(dev *Device)
 {
   return uecho_object_getmanufacturercode(obj);
 }
*/
