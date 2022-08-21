// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/cybergarage/uecho-go/net/echonet"
)

const (
	ProgramName            = "uecholight"
	LightObjectCode        = echonet.ObjectCode(0x029101)
	LightPropertyPowerCode = 0x80
	LightPropertyPowerOn   = 0x30
	LightPropertyPowerOff  = 0x31
	EXIT_SUCCESS           = 0
	EXIT_FAILURE           = 1
)
