// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

const (
	Port             = 3610
	UDPPort          = Port
	TCPPort          = Port
	MulticastAddress = "224.0.23.0"
	MaxPacketSize    = 1024
)

//  Extention for Echonet Lite
const (
	UDPPortMin = UDPPort
	UDPPortMax = UDPPort + 80
)

const (
	MessageFormatSpecified = 0x01
	MessageFormatArbitrary = 0x02
)
