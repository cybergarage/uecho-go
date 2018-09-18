// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet/log"
)

const (
	logSocketTypeUDP     = "U"
	logSocketTypeTCP     = "T"
	logSocketWriteFormat = "W (%s) : %s -> %s (%d) : %s"
	logSocketReadFormat  = "R (%s) : %s <- %s (%d) : %s"
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
