// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/cybergarage/uecho-go/net/echonet"
)

func NewLightDevice() *echonet.Device {

	dev := echonet.NewDevice()

	// TODO : Set your manufacture code
	dev.SetManufacturerCode(echonet.DeviceManufacturerUnknown)

	dev.SetCode(LightObjectCode)

	dev.CreateProperty(LightPropertyPowerCode, echonet.PropertyAttributeGetSet)
	dev.SetPropertyIntegerData(LightPropertyPowerCode, LightPropertyPowerOn, 1)

	return dev
}
