// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"net"
)

const (
	MessageHeaderLen = (1 + 1 + 2)
	MessageMinLen    = (MessageHeaderLen + 3 + 3 + 1 + 1)
	EHD1             = 0x10
	EHD2             = 0x81
	TIDSize          = 2
	TIDMax           = 65535
	EOJSize          = 3
)

const (
	EsvWriteRequest                      = 0x60
	EsvWriteRequestResponseRequired      = 0x61
	EsvReadRequest                       = 0x62
	EsvNotificationRequest               = 0x63
	EsvWriteReadRequest                  = 0x6E
	EsvWriteResponse                     = 0x71
	EsvReadResponse                      = 0x72
	EsvNotification                      = 0x73
	EsvNotificationResponseRequired      = 0x74
	EsvNotificationResponse              = 0x7A
	EsvWriteReadResponse                 = 0x7E
	EsvWriteRequestError                 = 0x50
	EsvWriteRequestResponseRequiredError = 0x51
	EsvReadRequestError                  = 0x52
	EsvNotificationRequestError          = 0x53
	EsvWriteReadRequestError             = 0x5E
)

// Message is an instance for Echonet message.
type Message struct {
	EHD1      byte
	EHD2      byte
	TID       [TIDSize]byte
	SEOJ      [EOJSize]byte
	DEOJ      [EOJSize]byte
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
	msg := &Message{
		EHD1: EHD1,
		EHD2: EHD2,
	}
	return msg
}

// NewMessageWithBytes returns a new message of the specified bytes.
func NewMessageWithBytes(data []byte) (*Message, error) {
	msg := NewMessage()
	err := msg.Parse(data)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// SetID sets the specified ID.
func (msg *Message) SetID(value int) error {
	if TIDMax < value {
		value %= TIDMax
	}
	msg.TID[0] = (byte)((value & 0xFF00) >> 8)
	msg.TID[1] = (byte)(value & 0x00FF)

	return nil
}

// GetID returns the stored ID.
func (msg *Message) GetID(value int) int {
	return (((int)(msg.TID[0]) << 8) + (int)(msg.TID[1]))
}

// Parse parses the specified bytes.
func (msg *Message) Parse(data []byte) error {
	return nil
}
