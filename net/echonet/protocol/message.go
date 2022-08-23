// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
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

// Message is an instance for Echonet message.
type Message struct {
	EHD1Echonet byte
	EHD2Format1 byte
	tid         []byte
	seoj        []byte
	deoj        []byte
	esv         ESV
	opc         byte
	ep          []*Property
	From        *Address
	pktType     int
	Interface   *net.Interface
}

// NewMessage returns a new message.
func NewMessage() *Message {
	msg := &Message{
		EHD1Echonet: EHD1Echonet,
		EHD2Format1: EHD2Format1,
		tid:         make([]byte, TIDSize),
		seoj:        make([]byte, EOJSize),
		deoj:        make([]byte, EOJSize),
		esv:         0,
		opc:         0,
		ep:          make([]*Property, 0),
		From:        NewAddress(),
		pktType:     UnknownPacket,
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

	copyMsg.pktType = msg.pktType
	copyMsg.Interface = msg.Interface

	return copyMsg, nil
}

// NewResponseMessageWithMessage returns a response message of the specified message withtout the properties.
func NewResponseMessageWithMessage(reqMsg *Message) *Message {
	msg := NewMessage()
	msg.SetTID(reqMsg.TID())
	msg.SetSEOJ(reqMsg.DEOJ())
	msg.SetDEOJ(reqMsg.SEOJ())

	switch reqMsg.ESV() {
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
	msg.SetTID(reqMsg.TID())
	msg.SetSEOJ(reqMsg.DEOJ())
	msg.SetDEOJ(reqMsg.SEOJ())

	switch reqMsg.ESV() {
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

	reqMsgOPC := reqMsg.OPC()
	msg.SetOPC(reqMsgOPC)
	for n := 0; n < reqMsgOPC; n++ {
		reqProp := msg.PropertyAt(n)
		msg.AddProperty(reqProp)
	}
	return msg
}

// SetTID sets the specified TID.
func (msg *Message) SetTID(value uint) error {
	if TIDMax < value {
		value %= TIDMax
	}
	msg.tid[0] = (byte)((value & 0xFF00) >> 8)
	msg.tid[1] = (byte)(value & 0x00FF)
	return nil
}

// TID returns the stored TID.
func (msg *Message) TID() uint {
	return (((uint)(msg.tid[0]) << 8) + (uint)(msg.tid[1]))
}

// IsTID returns true whether the specified value equals the message TID, otherwise false.
func (msg *Message) IsTID(tid uint) bool {
	return msg.TID() == tid
}

// SetSEOJ sets a source object code.
func (msg *Message) SetSEOJ(code ObjectCode) {
	encoding.IntegerToByte(uint(code), msg.seoj)
}

// SEOJ returns the source object code.
func (msg *Message) SEOJ() ObjectCode {
	return ObjectCode(encoding.ByteToInteger(msg.seoj))
}

// IsSEOJ returns true whether the specified value equals the message source object code, otherwise false.
func (msg *Message) IsSEOJ(code ObjectCode) bool {
	return msg.SEOJ() == code
}

// SetDEOJ sets a destination object code.
func (msg *Message) SetDEOJ(code ObjectCode) {
	encoding.IntegerToByte(uint(code), msg.deoj)
}

// DEOJ returns the destination object code.
func (msg *Message) DEOJ() ObjectCode {
	return ObjectCode(encoding.ByteToInteger(msg.deoj))
}

// IsDEOJ returns true whether the specified value equals the message destination object code, otherwise false.
func (msg *Message) IsDEOJ(code ObjectCode) bool {
	return msg.DEOJ() == code
}

// SetESV sets the specified ESV.
func (msg *Message) SetESV(value ESV) {
	msg.esv = value
}

// ESV returns the stored ESV.
func (msg *Message) ESV() ESV {
	return msg.esv
}

// IsESV returns true whether the specified code equals the message ESV, otherwise false.
func (msg *Message) IsESV(esv ESV) bool {
	return msg.esv == esv
}

// IsValidESV returns true whether the specified code is valid, otherwise false.
func (msg *Message) IsValidESV() bool {
	return IsValidESV(msg.esv)
}

// IsWriteRequest returns true whether the message is a write request type, otherwise false.
func (msg *Message) IsWriteRequest() bool {
	return IsWriteRequest(msg.esv)
}

// IsReadRequest returns true whether the message is a read request type, otherwise false.
func (msg *Message) IsReadRequest() bool {
	return IsReadRequest(msg.esv)
}

// IsNotificationRequest returns true whether the message is a notification request type, otherwise false.
func (msg *Message) IsNotificationRequest() bool {
	return IsNotificationRequest(msg.esv)
}

// IsWriteResponse returns true whether the message is a write response type, otherwise false.
func (msg *Message) IsWriteResponse() bool {
	return IsWriteResponse(msg.esv)
}

// IsReadResponse returns true whether the message is a read response type, otherwise false.
func (msg *Message) IsReadResponse() bool {
	return IsReadResponse(msg.esv)
}

// IsNotification returns true whether the message is a notification type, otherwise false.
func (msg *Message) IsNotification() bool {
	return IsNotification(msg.esv)
}

// IsNotificationResponse returns true whether the message is a notification response type, otherwise false.
func (msg *Message) IsNotificationResponse() bool {
	return IsNotificationResponse(msg.esv)
}

// IsResponseRequired returns true whether the ESV requires the response, otherwise false.
func (msg *Message) IsResponseRequired() bool {
	return IsResponseRequired(msg.esv)
}

// SetOPC sets the specified OPC.
func (msg *Message) SetOPC(value int) error {
	msg.opc = byte(value & 0xFF)
	msg.ep = make([]*Property, msg.opc)
	for n := 0; n < int(msg.opc); n++ {
		msg.ep[n] = NewProperty()
	}
	return nil
}

// OPC returns the stored OPC.
func (msg *Message) OPC() int {
	return int(msg.opc)
}

// AddProperty adds a property.
func (msg *Message) AddProperty(prop *Property) {
	msg.opc++
	msg.ep = append(msg.ep, prop)
}

// AddProperties adds a properties.
func (msg *Message) AddProperties(props []*Property) {
	for _, prop := range props {
		msg.AddProperty(prop)
	}
}

// PropertyAt returns the specified property.
func (msg *Message) PropertyAt(n int) *Property {
	if (len(msg.ep) - 1) < n {
		return nil
	}
	return msg.ep[n]
}

// Properties returns the all properties.
func (msg *Message) Properties() []*Property {
	return msg.ep
}

// HasProperty returns true when the message has the specified property, otherwise false.
func (msg *Message) HasProperty(propCode PropertyCode) bool {
	for _, prop := range msg.ep {
		if prop.Code() == propCode {
			return true
		}
	}
	return false
}

// SourceAddress returns the source address of the message.
func (msg *Message) SourceAddress() string {
	return msg.From.IP.String()
}

// SourcePort returns the source address of the message.
func (msg *Message) SourcePort() int {
	return msg.From.Port
}

// SetPacketType sets the specified package type to the message.
func (msg *Message) SetPacketType(packetType int) {
	msg.pktType = packetType
}

// PacketType returns the packet type of the message.
func (msg *Message) PacketType() int {
	return msg.pktType
}

// IsPacketType returns true when the specified type equals the message type, otherwise false.
func (msg *Message) IsPacketType(packetType int) bool {
	return (msg.pktType & packetType) != 0
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

	for n := 0; n < int(msg.opc); n++ {
		prop := msg.PropertyAt(n)
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
	msgBytes[2] = msg.tid[0]
	msgBytes[3] = msg.tid[1]
	msgBytes[4] = msg.seoj[0]
	msgBytes[5] = msg.seoj[1]
	msgBytes[6] = msg.seoj[2]
	msgBytes[7] = msg.deoj[0]
	msgBytes[8] = msg.deoj[1]
	msgBytes[9] = msg.deoj[2]
	msgBytes[10] = byte(msg.esv)
	msgBytes[11] = msg.opc

	offset := 12
	for n := 0; n < int(msg.opc); n++ {
		prop := msg.PropertyAt(n)
		if prop == nil {
			continue
		}
		msgBytes[offset] = byte(prop.Code())
		offset++

		propSize := prop.Size()
		msgBytes[offset] = byte(propSize)
		offset++
		if propSize == 0 {
			continue
		}

		propData := prop.Data()
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
