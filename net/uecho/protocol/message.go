// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"fmt"
	"net"

	"github.com/cybergarage/uecho-go/net/uecho/encoding"
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

type ESV int

const (
	ESVWriteRequest                      = 0x60
	ESVWriteRequestResponseRequired      = 0x61
	ESVReadRequest                       = 0x62
	ESVNotificationRequest               = 0x63
	ESVWriteReadRequest                  = 0x6E
	ESVWriteResponse                     = 0x71
	ESVReadResponse                      = 0x72
	ESVNotification                      = 0x73
	ESVNotificationResponseRequired      = 0x74
	ESVNotificationResponse              = 0x7A
	ESVWriteReadResponse                 = 0x7E
	ESVWriteRequestError                 = 0x50
	ESVWriteRequestResponseRequiredError = 0x51
	ESVReadRequestError                  = 0x52
	ESVNotificationRequestError          = 0x53
	ESVWriteReadRequestError             = 0x5E
)

const (
	errorShortMessageLength   = "Short message length : %d < %d"
	errorInvalidMessageHeader = "Invalid Message header [%d] : %02X != %02X"
)

// Message is an instance for Echonet message.
type Message struct {
	EHD1      byte
	EHD2      byte
	TID       []byte
	SEOJ      []byte
	DEOJ      []byte
	ESV       ESV
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
		TID:  make([]byte, TIDSize),
		SEOJ: make([]byte, EOJSize),
		DEOJ: make([]byte, EOJSize),
		ESV:  0,
		OPC:  0,
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

// SetTID sets the specified TID.
func (msg *Message) SetTID(value int) error {
	if TIDMax < value {
		value %= TIDMax
	}
	msg.TID[0] = (byte)((value & 0xFF00) >> 8)
	msg.TID[1] = (byte)(value & 0x00FF)
	return nil
}

// GetTID returns the stored TID.
func (msg *Message) GetTID() int {
	return (((int)(msg.TID[0]) << 8) + (int)(msg.TID[1]))
}

// SetSourceObjectCode sets a source object code.
func (msg *Message) SetSourceObjectCode(code uint) {
	encoding.IntegerToByte(code, msg.SEOJ)
}

// GetSourceObjectCode returns the source object code.
func (msg *Message) GetSourceObjectCode() uint {
	return encoding.ByteToInteger(msg.SEOJ)
}

// SetDestinationObjectCode sets a source object code.
func (msg *Message) SetDestinationObjectCode(code uint) {
	encoding.IntegerToByte(code, msg.DEOJ)
}

// GetDestinationObjectCode returns the source object code.
func (msg *Message) GetDestinationObjectCode() uint {
	return encoding.ByteToInteger(msg.DEOJ)
}

// SetESV sets the specified ESV.
func (msg *Message) SetESV(value ESV) {
	msg.ESV = value
}

// GetESV returns the stored ESV.
func (msg *Message) GetESV() ESV {
	return msg.ESV
}

// SetOPC sets the specified OPC.
func (msg *Message) SetOPC(value byte) {
	msg.OPC = value
}

// GetOPC returns the stored OPC.
func (msg *Message) GetOPC() byte {
	return msg.OPC
}

// Parse parses the specified bytes.
func (msg *Message) Parse(data []byte) error {
	msgLen := len(data)
	if msgLen < MessageMinLen {
		return fmt.Errorf(errorShortMessageLength, msgLen, MessageMinLen)
	}

	// Check Headers

	if data[0] != EHD1 {
		return fmt.Errorf(errorInvalidMessageHeader, 0, data[0], EHD1)
	}

	if data[1] != EHD2 {
		return fmt.Errorf(errorInvalidMessageHeader, 1, data[1], EHD2)
	}

	// TID

	msg.TID[0] = data[2]
	msg.TID[1] = data[3]

	// SEOJ

	msg.SEOJ[0] = data[4]
	msg.SEOJ[1] = data[5]
	msg.SEOJ[2] = data[6]

	// DEOJ

	msg.DEOJ[0] = data[7]
	msg.DEOJ[1] = data[8]
	msg.DEOJ[2] = data[9]

	// ESV

	msg.ESV = ESV(data[10])

	// OPC

	msg.OPC = data[11]

	/*

		// EP

		offset = 12;
		for (n = 0; n<(int)(msg.OPC); n++) {
		  prop = uecho_message_getproperty(msg, n);
		  if (!prop)
			return false;

		  // EPC

		  if ((dataLen - 1) < offset)
			return false;
		  uecho_property_setcode(prop, data[offset++]);

		  // PDC

		  if ((dataLen - 1) < offset)
			return false;
		  count = data[offset++];

		  // EDT

		  if ((dataLen - 1) < (offset + count - 1))
			return false;
		  if (!uecho_property_setdata(prop, (data + offset), count))
			return false;
		  offset += count;
		}
	*/

	return nil
}
