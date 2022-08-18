// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"time"
)

const (
	Port                 = 3610
	UDPPort              = Port
	TCPPort              = Port
	MulticastIPv4Address = "224.0.23.0"
	MulticastIPv6Address = "ff02::1"
	MaxPacketSize        = 1024
)

// Extension for Echonet Lite.
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
