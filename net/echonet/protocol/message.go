// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"bytes"
	"encoding/hex"
	"io"
	"net"

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

const (
	FrameHeaderSize           = (1 + 1 + 2)
	Format1HeaderSize         = (3 + 3 + 1 + 1)
	Format1MinSize            = (FrameHeaderSize + Format1HeaderSize)
	Format1PropertyHeaderSize = 2
	EHD1Echonet               = 0x10
	EHD2Format1               = 0x81
	TIDSize                   = 2
	TIDMax                    = 65535
	EOJSize                   = 3
)

const (
	UnknownPacket    = 0x00
	MulticastPacket  = 0x01
	UDPUnicastPacket = 0x10
	TCPUnicastPacket = 0x20
	UnicastPacket    = (UDPUnicastPacket | TCPUnicastPacket)
)

const (
	errorShortMessageSize     = "Short message length : %d < %d"
	errorInvalidMessageHeader = "Invalid Message header [%d] : %02X != %02X"
)

// Message is an instance for Echonet message.
type Message struct {
	EHD1Echonet byte
	EHD2Format1 byte
	TID         []byte
	SEOJ        []byte
	DEOJ        []byte
	ESV         ESV
	OPC         byte
	EP          []*Property
	From        *Address
	PacketType  int
	Interface   *net.Interface
}

// NewMessage returns a new message.
func NewMessage() *Message {
	msg := &Message{
		EHD1Echonet: EHD1Echonet,
		EHD2Format1: EHD2Format1,
		TID:         make([]byte, TIDSize),
		SEOJ:        make([]byte, EOJSize),
		DEOJ:        make([]byte, EOJSize),
		ESV:         0,
		OPC:         0,
		EP:          make([]*Property, 0),
		From:        NewAddress(),
		PacketType:  UnknownPacket,
		Interface:   nil,
	}
	return msg
}

// NewMessageWithReader returns a new message with the specified reader.
func NewMessageWithReader(reader io.Reader) (*Message, error) {
	msg := NewMessage()
	if err := msg.ParseReader(reader); err != nil {
		return nil, err
	}
	return msg, nil
}

// NewMessageWithBytes returns a new message of the specified bytes.
func NewMessageWithBytes(data []byte) (*Message, error) {
	msg := NewMessage()
	if err := msg.ParseBytes(data); err != nil {
		return nil, err
	}
	return msg, nil
}

// NewMessageWithMessage copies the specified message.
func NewMessageWithMessage(msg *Message) (*Message, error) {
	copyMsg, err := NewMessageWithBytes(msg.Bytes())
	if err != nil {
		return nil, err
	}

	from := *msg.From
	copyMsg.From = &from

	copyMsg.PacketType = msg.PacketType
	copyMsg.Interface = msg.Interface

	return copyMsg, nil
}

// NewResponseMessageWithMessage returns a response message of the specified message withtout the properties.
func NewResponseMessageWithMessage(reqMsg *Message) *Message {
	msg := NewMessage()
	msg.SetTID(reqMsg.GetTID())
	msg.SetSourceObjectCode(reqMsg.GetDestinationObjectCode())
	msg.SetDestinationObjectCode(reqMsg.GetSourceObjectCode())

	switch reqMsg.GetESV() {
	case ESVWriteRequestResponseRequired:
		msg.SetESV(ESVWriteResponse)
	case ESVReadRequest:
		msg.SetESV(ESVReadResponse)
	case ESVNotificationRequest:
		msg.SetESV(ESVNotification)
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

// IsTID returns true whether the specified value equals the message TID, otherwise false.
func (msg *Message) IsTID(tid uint) bool {
	return msg.GetTID() == tid
}

// SetSourceObjectCode sets a source object code.
func (msg *Message) SetSourceObjectCode(code ObjectCode) {
	encoding.IntegerToByte(uint(code), msg.SEOJ)
}

// GetSourceObjectCode returns the source object code.
func (msg *Message) GetSourceObjectCode() ObjectCode {
	return ObjectCode(encoding.ByteToInteger(msg.SEOJ))
}

// IsSourceObjectCode returns true whether the specified value equals the message source object code, otherwise false.
func (msg *Message) IsSourceObjectCode(code ObjectCode) bool {
	return msg.GetSourceObjectCode() == code
}

// SetDestinationObjectCode sets a source object code.
func (msg *Message) SetDestinationObjectCode(code ObjectCode) {
	encoding.IntegerToByte(uint(code), msg.DEOJ)
}

// GetDestinationObjectCode returns the source object code.
func (msg *Message) GetDestinationObjectCode() ObjectCode {
	return ObjectCode(encoding.ByteToInteger(msg.DEOJ))
}

// IsDestinationObjectCode returns true whether the specified value equals the message destination object code, otherwise false.
func (msg *Message) IsDestinationObjectCode(code ObjectCode) bool {
	return msg.GetDestinationObjectCode() == code
}

// SetESV sets the specified ESV.
func (msg *Message) SetESV(value ESV) {
	msg.ESV = value
}

// GetESV returns the stored ESV.
func (msg *Message) GetESV() ESV {
	return msg.ESV
}

// IsESV returns true whether the specified code equals the message ESV, otherwise false.
func (msg *Message) IsESV(esv ESV) bool {
	return msg.ESV == esv
}

// IsValidESV returns true whether the specified code is valid, otherwise false.
func (msg *Message) IsValidESV() bool {
	return IsValidESV(msg.ESV)
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

// IsNotification returns true whether the message is a notification type, otherwise false.
func (msg *Message) IsNotification() bool {
	return IsNotification(msg.ESV)
}

// IsNotificationResponse returns true whether the message is a notification response type, otherwise false.
func (msg *Message) IsNotificationResponse() bool {
	return IsNotificationResponse(msg.ESV)
}

// IsResponseRequired returns true whether the ESV requires the response, otherwise false.
func (msg *Message) IsResponseRequired() bool {
	return IsResponseRequired(msg.ESV)
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

// AddProperties adds a properties.
func (msg *Message) AddProperties(props []*Property) error {
	for _, prop := range props {
		err := msg.AddProperty(prop)
		if err != nil {
			return err
		}
	}
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

// HasProperty returns true when the message has the specified property, otherwise false.
func (msg *Message) HasProperty(propCode PropertyCode) bool {
	for _, prop := range msg.EP {
		if prop.GetCode() == propCode {
			return true
		}
	}
	return false
}

// GetSourceAddress returns the source address of the message.
func (msg *Message) GetSourceAddress() string {
	return msg.From.IP.String()
}

// GetSourcePort returns the source address of the message.
func (msg *Message) GetSourcePort() int {
	return msg.From.Port
}

// SetPacketType sets the specified package type to the message.
func (msg *Message) SetPacketType(packetType int) {
	msg.PacketType = packetType
}

// GetPacketType returns the packet type of the message.
func (msg *Message) GetPacketType() int {
	return msg.PacketType
}

// IsPacketType returns true when the specified type equals the message type, otherwise false.
func (msg *Message) IsPacketType(packetType int) bool {
	return (msg.PacketType & packetType) != 0
}

// IsMulticastPacket returns true when the message was sent by multicast, otherwise false.
func (msg *Message) IsMulticastPacket() bool {
	return msg.IsPacketType(MulticastPacket)
}

// IsUnicastPacket returns true when the message was sent by TCP or UDP unicast, otherwise false.
func (msg *Message) IsUnicastPacket() bool {
	return msg.IsPacketType(UnicastPacket)
}

// IsTCPUnicastPacket returns true when the message was sent by TCP unicast, otherwise false.
func (msg *Message) IsTCPUnicastPacket() bool {
	return msg.IsPacketType(TCPUnicastPacket)
}

// IsUDPUnicastPacket returns true when the message was sent by UDP unicast, otherwise false.
func (msg *Message) IsUDPUnicastPacket() bool {
	return msg.IsPacketType(UDPUnicastPacket)
}

// Size return the byte size.
func (msg *Message) Size() int {
	msgSize := Format1MinSize

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
	if msg == nil {
		return make([]byte, 0)
	}

	msgBytes := make([]byte, msg.Size())

	msgBytes[0] = msg.EHD1Echonet
	msgBytes[1] = msg.EHD2Format1
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
		msgBytes[offset] = byte(prop.GetCode())
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

// Equals returns true whether the specified other message is same, otherwise false.
func (msg *Message) Equals(other *Message) bool {
	if msg == nil || other == nil {
		return false
	}
	if !bytes.Equal(msg.Bytes(), other.Bytes()) {
		return false
	}
	return true
}

// String return the string .
func (msg *Message) String() string {
	if msg == nil {
		return ""
	}
	return hex.EncodeToString(msg.Bytes())
}
