// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/cybergarage/uecho-go/net/uecho"
	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

const (
	LightObjectCode        = 0x029101
	LightPropertyPowerCode = 0x80
	LightPropertyPowerOn   = 0x30
	LightPropertyPowerOff  = 0x31
)

type Light struct {
	*uecho.Node
}

// NewLight returns a new light device.
func NewLight() (*Light, error) {

	light := &Light{
		Node: uecho.NewNode(),
	}

	dev := uecho.NewDevice()

	// TODO : Set your manufacture code
	dev.SetManufacturerCode(uecho.DeviceManufacturerUnknown)

	dev.SetCode(LightObjectCode)

	dev.CreateProperty(LightPropertyPowerCode, protocol.PropertyAttributeReadWriteAnno)
	dev.SetPropertyIntegerData(LightPropertyPowerCode, LightPropertyPowerOn, 1)

	light.AddDevice(dev)
	dev.AddListener(light)

	return light, nil
}

// PropertyRequestReceived is a listener for Echonet requests.
func (light *Light) PropertyRequestReceived(obj *Object, esv protocol.ESV, prop *protocol.Property) error {}
	return nil
}
