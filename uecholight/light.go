// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/cybergarage/uecho-go/net/uecho"
)

const (
	LightObjectCode        = 0x029101
	LightPropertyPowerCode = 0x80
	LightPropertyPowerOn   = 0x30
	LightPropertyPowerOff  = 0x31
)

// NewLight returns a new light device.
func NewLight() *uecho.Device {
	dev := uecho.NewDevice()

	// TODO : Set your manufacture code
	dev.SetManufacturerCode(uecho.DeviceManufacturerUnknown)

	dev.SetCode(LightObjectCode)

	dev.CreateProperty(LightPropertyPowerCode, uecho.PropertyAttributeReadWriteAnno)
	dev.SetPropertyIntegerData(LightPropertyPowerCode, LightPropertyPowerOn, 1)

	return dev
}
