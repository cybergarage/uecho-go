// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"encoding/hex"
	"fmt"
	"net"

	"github.com/cybergarage/uecho-go/net/uecho/encoding"
)

const (
	MessageHeaderSize = (1 + 1 + 2)
	MessageMinSize    = (MessageHeaderSize + 3 + 3 + 1 + 1)
	EHD1              = 0x10
	EHD2              = 0x81
	TIDSize           = 2
	TIDMax            = 65535
	EOJSize           = 3
)

const (
	errorShortMessageSize     = "Short message length : %d < %d"
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
	EP        []*Property
	rawBytes  []byte
	From      net.UDPAddr
	Interface net.Interface
}

// NewMessage returns a new message.
func NewMessage() *Message {
	msg := &Message{
		EHD1:     EHD1,
		EHD2:     EHD2,
		TID:      make([]byte, TIDSize),
		SEOJ:     make([]byte, EOJSize),
		DEOJ:     make([]byte, EOJSize),
		ESV:      0,
		OPC:      0,
		EP:       make([]*Property, 0),
		rawBytes: make([]byte, 0),
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
	msg.rawBytes = data
	return msg, nil
}

// NewResponseMessageWithMessage returns a response message of the specified message withtout the properties.
func NewResponseMessageWithMessage(reqMsg *Message) *Message {
	msg := NewMessage()
	msg.SetTID(reqMsg.GetTID())
	msg.SetSourceObjectCode(reqMsg.GetDestinationObjectCode())
	msg.SetDestinationObjectCode(reqMsg.GetSourceObjectCode())

	switch reqMsg.GetESV() {
	case ESVWriteRequest:
		msg.SetESV(ESVWriteResponse)
	case ESVWriteRequestResponseRequired:
		msg.SetESV(ESVWriteResponse)
	case ESVReadRequest:
		msg.SetESV(ESVReadResponse)
	case ESVNotificationRequest:
		msg.SetESV(ESVNotificationResponse)
	case ESVWriteReadRequest:
		msg.SetESV(ESVWriteReadResponse)
	case ESVNotificationResponseRequired:
		msg.SetESV(ESVNotificationResponse)
	default:
		msg.SetESV(0)
	}

	return msg
}

// NewImpossibleMessageWithMessage returns a impossible message of the specified message.
func NewImpossibleMessageWithMessage(reqMsg *Message) *Message {
	msg := NewMessage()
	msg.SetTID(reqMsg.GetTID())
	msg.SetSourceObjectCode(reqMsg.GetDestinationObjectCode())
	msg.SetDestinationObjectCode(reqMsg.GetSourceObjectCode())

	switch reqMsg.GetESV() {
	case ESVWriteRequest:
		msg.SetESV(ESVWriteRequestError)
	case ESVWriteRequestResponseRequired:
		msg.SetESV(ESVWriteRequestResponseRequiredError)
	case ESVReadRequest:
		msg.SetESV(ESVReadRequestError)
	case ESVNotificationRequest:
		msg.SetESV(ESVNotificationRequestError)
	case ESVWriteReadRequest:
		msg.SetESV(ESVWriteReadRequestError)
	case ESVNotificationResponseRequired:
		msg.SetESV(ESVNotificationRequestError)
	default:
		msg.SetESV(0)
	}

	reqMsgOPC := reqMsg.GetOPC()
	msg.SetOPC(reqMsgOPC)
	for n := 0; n < reqMsgOPC; n++ {
		reqProp := msg.GetProperty(n)
		msg.AddProperty(reqProp)
	}
	return msg
}

// SetTID sets the specified TID.
func (msg *Message) SetTID(value uint) error {
	if TIDMax < value {
		value %= TIDMax
	}
	msg.TID[0] = (byte)((value & 0xFF00) >> 8)
	msg.TID[1] = (byte)(value & 0x00FF)
	return nil
}

// GetTID returns the stored TID.
func (msg *Message) GetTID() uint {
	return (((uint)(msg.TID[0]) << 8) + (uint)(msg.TID[1]))
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

// IsWriteRequest returns true whether the message is a write request type, otherwise false.
func (msg *Message) IsWriteRequest() bool {
	return IsWriteRequest(msg.ESV)
}

// IsReadRequest returns true whether the message is a read request type, otherwise false.
func (msg *Message) IsReadRequest() bool {
	return IsReadRequest(msg.ESV)
}

// IsNotificationRequest returns true whether the message is a notification request type, otherwise false.
func (msg *Message) IsNotificationRequest() bool {
	return IsNotificationRequest(msg.ESV)
}

// IsWriteResponse returns true whether the message is a write response type, otherwise false.
func (msg *Message) IsWriteResponse() bool {
	return IsWriteResponse(msg.ESV)
}

// IsReadResponse returns true whether the message is a read response type, otherwise false.
func (msg *Message) IsReadResponse() bool {
	return IsReadResponse(msg.ESV)
}

// IsNotificationResponse returns true whether the message is a notification response type, otherwise false.
func (msg *Message) IsNotificationResponse() bool {
	return IsNotificationResponse(msg.ESV)
}

// IsResponseRequired returns true whether the ESV requires the response, otherwise false.
func (msg *Message) IsResponseRequired() bool {
	return IsResponseRequired(msg.ESV)
}

// GetSourceAddress returns the source address of the message.
func (msg *Message) GetSourceAddress() string {
	return msg.From.IP.String()
}

// GetSourcePort returns the source address of the message.
func (msg *Message) GetSourcePort() int {
	return msg.From.Port
}

// SetOPC sets the specified OPC.
func (msg *Message) SetOPC(value int) error {
	msg.OPC = byte(value & 0xFF)
	msg.EP = make([]*Property, msg.OPC)
	for n := 0; n < int(msg.OPC); n++ {
		msg.EP[n] = NewProperty()
	}
	return nil
}

// GetOPC returns the stored OPC.
func (msg *Message) GetOPC() int {
	return int(msg.OPC)
}

// AddProperty adds a property.
func (msg *Message) AddProperty(prop *Property) error {
	msg.OPC++
	msg.EP = append(msg.EP, prop)
	return nil
}

// GetProperty returns the specified property.
func (msg *Message) GetProperty(n int) *Property {
	if (len(msg.EP) - 1) < n {
		return nil
	}
	return msg.EP[n]
}

// GetProperties returns the all properties.
func (msg *Message) GetProperties() []*Property {
	return msg.EP
}

// Parse parses the specified bytes.
func (msg *Message) Parse(data []byte) error {
	dataSize := len(data)
	if dataSize < MessageMinSize {
		return fmt.Errorf(errorShortMessageSize, dataSize, MessageMinSize)
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

	err := msg.SetOPC(int(data[11]))
	if err != nil {
		return err
	}

	// EP

	offset := 12
	for n := 0; n < int(msg.OPC); n++ {
		prop := msg.GetProperty(n)
		if prop == nil {
			continue
		}

		// EPC

		if (dataSize - 1) < offset {
			continue
		}

		prop.Code = data[offset]
		offset++

		// PDC

		if (dataSize - 1) < offset {
			continue
		}

		propSize := int(data[offset])
		offset++

		// EDT

		if (dataSize - 1) < (offset + propSize - 1) {
			continue
		}

		prop.Data = data[offset:(offset + propSize)]

		offset += propSize
	}

	return nil
}

// Size return the byte size.
func (msg *Message) Size() int {

	msgSize := MessageMinSize

	for n := 0; n < int(msg.OPC); n++ {
		prop := msg.GetProperty(n)
		if prop == nil {
			continue
		}
		msgSize += 2
		msgSize += prop.Size()
	}

	return msgSize
}

// Bytes return the message bytes.
func (msg *Message) Bytes() []byte {

	msgBytes := make([]byte, msg.Size())

	msgBytes[0] = msg.EHD1
	msgBytes[1] = msg.EHD2
	msgBytes[2] = msg.TID[0]
	msgBytes[3] = msg.TID[1]
	msgBytes[4] = msg.SEOJ[0]
	msgBytes[5] = msg.SEOJ[1]
	msgBytes[6] = msg.SEOJ[2]
	msgBytes[7] = msg.DEOJ[0]
	msgBytes[8] = msg.DEOJ[1]
	msgBytes[9] = msg.DEOJ[2]
	msgBytes[10] = byte(msg.ESV)
	msgBytes[11] = msg.OPC

	offset := 12
	for n := 0; n < int(msg.OPC); n++ {
		prop := msg.GetProperty(n)
		if prop == nil {
			continue
		}
		msgBytes[offset] = prop.GetCode()
		offset++

		propSize := int(prop.Size())
		msgBytes[offset] = byte(propSize)
		offset++
		if propSize == 0 {
			continue
		}

		propData := prop.GetData()
		for i := 0; i < propSize; i++ {
			msgBytes[offset+i] = propData[i]
		}

		offset += propSize
	}

	return msgBytes
}

// String return the string .
func (msg *Message) String() string {
	return hex.EncodeToString(msg.Bytes())
}
