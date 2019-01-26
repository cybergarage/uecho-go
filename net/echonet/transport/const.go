// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"time"
)

const (
	Port             = 3610
	UDPPort          = Port
	TCPPort          = Port
	MulticastAddress = "224.0.23.0"
	MaxPacketSize    = 1024
)

//  Extention for Echonet Lite
const (
	UDPPortRange = 100
)

const (
	MessageFormatSpecified = 0x01
	MessageFormatArbitrary = 0x02
)

const (
	DefaultConnectimeTimeOut = (time.Millisecond * 5000)
	DefaultRequestTimeout    = (time.Millisecond * 5000)
	DefaultBindRetryCount    = 5
	DefaultBindRetryWaitTime = (time.Millisecond * 500)
)
