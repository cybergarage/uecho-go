// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// NewLightNode returns a new light device.
func NewLightNode() echonet.LocalNode {
	onRequest := func(obj echonet.Object, esv protocol.ESV, reqProp protocol.Property) error {
		// Only handle write requests.
		if !esv.IsWriteRequest() {
			return nil
		}
		reqPropCode := reqProp.Code()
		reqPropData := reqProp.Data()
		// NOTE: Object::LookupProperty() will always find the property here because onRequest is only called if the object has the specified property.
		targetProp, _ := obj.LookupProperty(reqPropCode)
		switch reqPropCode {
		case 0x80: // Operation status
			// Accept only 0x30 (ON) or 0x31 (OFF) as valid data.
			reqData, err := reqProp.AsByte()
			if err != nil {
				return err
			}
			switch reqData {
			case 0x30, 0x31:
				targetProp.SetData(reqPropData)
				OutputMessage("0x%02X : 0x%s -> 0x%s", esv, hex.EncodeToString(targetProp.Data()), hex.EncodeToString(reqPropData))
				return nil
			default:
			}
		}
		err := fmt.Errorf("invalid request : %02X %s", reqPropCode, hex.EncodeToString(reqPropData))
		OutputError(err)
		return err
	}

	dev, _ := echonet.NewDevice(
		echonet.WithDeviceCode(LightObjectCode),
		echonet.WithDeviceManufacturerCode(echonet.DeviceManufacturerExperimental),
		echonet.WithDeviceRequestHandler(onRequest),
	)

	// Set initial property values
	dev.SetPropertyInteger(LightPropertyPowerCode, LightPropertyPowerOn, 1)

	node := echonet.NewLocalNode(
		echonet.WithLocalNodeDevices(dev),
	)
	return node
}
