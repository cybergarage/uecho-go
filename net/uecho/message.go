// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"net"
)

const (
	uEchoMessageHeaderLen = (1 + 1 + 2)
	uEchoMessageMinLen    = (uEchoMessageHeaderLen + 3 + 3 + 1 + 1)
	uEchoEhd1             = 0x10
	uEchoEhd2             = 0x81
	uEchoTIDSize          = 2
	uEchoEOJSize          = 3
)

// Message is an instance for Echonet message.
type Message struct {
	EHD1      byte
	EHD2      byte
	TID       [uEchoTIDSize]byte
	SEOJ      [uEchoEOJSize]byte
	DEOJ      [uEchoEOJSize]byte
	ESV       int
	OPC       byte
	Property  *Property
	bytes     []byte
	srcAddr   string
	From      net.UDPAddr
	Interface net.Interface
}

// NewMessage returns a new message.
func NewMessage() *Message {
	msg := &Message{}
	return msg
}

// Parse parses the specified bytes.
func (msg *Message) Parse(data []byte) error {
	return nil
}
