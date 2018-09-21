// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet/log"
)

const (
	logSocketTypeUDPMulticast = "UM"
	logSocketTypeUDPUnicast   = "UU"
	logSocketTypeTCPUnicast   = "TU"
	logSocketWriteFormat      = "W (%s) : %21s -> %21s (%d) : %s"
	logSocketReadFormat       = "R (%s) : %21s <- %21s (%d) : %s"
)

const (
	logSocketDirectionWrite = 0
	logSocketDirectionRead  = 1
)

func outputSocketLog(logLevel log.LogLevel, socketType string, socketDirection int, msgFrom string, msgTo string, msg string, msgSize int) {
	switch socketDirection {
	case logSocketDirectionWrite:
		{
			log.Output(logLevel, fmt.Sprintf(logSocketWriteFormat, socketType, msgFrom, msgTo, msgSize, msg))

		}
	case logSocketDirectionRead:
		{
			log.Output(logLevel, fmt.Sprintf(logSocketReadFormat, socketType, msgTo, msgFrom, msgSize, msg))

		}
	}
}
