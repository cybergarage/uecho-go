// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"github.com/cybergarage/go-logger/log"
)

const (
	logSocketTypeUDPMulticast = "UM"
	logSocketTypeUDPUnicast   = "UU"
	logSocketTypeTCPUnicast   = "TU"
	logSocketWriteFormat      = "S (%s) : %21s -> %21s (%d) : %s"
	logSocketReadFormat       = "R (%s) : %21s <- %21s (%d) : %s"
)

const (
	logSocketDirectionWrite = 0
	logSocketDirectionRead  = 1
)

func outputSocketLog(logLevel log.Level, socketType string, socketDirection int, msgFrom string, msgTo string, msg string, msgSize int) {
	switch socketDirection {
	case logSocketDirectionWrite:
		{
			log.Outputf(logLevel, logSocketWriteFormat, socketType, msgFrom, msgTo, msgSize, msg)
		}
	case logSocketDirectionRead:
		{
			log.Outputf(logLevel, logSocketReadFormat, socketType, msgTo, msgFrom, msgSize, msg)
		}
	}
}
