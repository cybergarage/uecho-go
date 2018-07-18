// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"testing"
)

const (
	errorMandatoryPropertyNotFound = "Mandatory Property (%0X) Not Found"
)

func TestNewDevice(t *testing.T) {
	dev := NewDevice()

	mandatoryPropertyCodes := []PropertyCode{
		DeviceOperatingStatus,
		DeviceInstallationLocation,
		DeviceStandardVersion,
		DeviceFaultStatus,
		DeviceManufacturerCode,
		ProfileGetPropertyMap,
		ProfileSetPropertyMap,
		ProfileAnnoPropertyMap,
	}

	for _, propCode := range mandatoryPropertyCodes {
		if !dev.HasProperty(propCode) {
			t.Errorf(errorMandatoryPropertyNotFound, propCode)
		}
	}
}
