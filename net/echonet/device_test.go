// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
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
	}

	for _, propCode := range mandatoryPropertyCodes {
		t.Run(fmt.Sprintf("%02X", propCode), func(t *testing.T) {
			if !dev.HasProperty(propCode) {
				t.Errorf(errorMandatoryPropertyNotFound, propCode)
			}
		})
	}

	testObjectPropertyMaps(t, dev)
}
