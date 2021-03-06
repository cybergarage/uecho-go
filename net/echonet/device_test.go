// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"testing"
)

func TestNewDevice(t *testing.T) {
	dev := NewDevice()

	mandatoryPropertyCodes := []PropertyCode{
		DeviceOperatingStatus,
		DeviceInstallationLocation,
		DeviceStandardVersion,
		DeviceFaultStatus,
		DeviceManufacturerCode,
		DeviceGetPropertyMap,
		DeviceSetPropertyMap,
		DeviceAnnoPropertyMap,
	}

	for _, propCode := range mandatoryPropertyCodes {
		if !dev.HasProperty(propCode) {
			t.Errorf(errorMandatoryPropertyNotFound, propCode)
		}
	}
}
