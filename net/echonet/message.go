// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// Message represents an Echonet message.
type Message interface {
	// SourceAddress returns the source address of the message.
	SourceAddress() string
	// SourcePort returns the source port of the message.
	SourcePort() int
	// SEOJ returns the source object code of the message.
	SEOJ() ObjectCode
	// DEOJ returns the destination object code of the message.
	DEOJ() ObjectCode
	// ESV returns the ESV of the message.
	ESV() ESV
	// OPC returns the OPC of the message.
	OPC() int
	// Properties returns the all properties of the message.
	Properties() []PropertyData
	// Property returns the n-th property of the message.
	Property(n int) (PropertyData, bool)
	// MessageMutator is an interface to mutate a message.
	MessageMutator
	// messageInternal is an interface to represent a message internal.
	messageInternal
}

type MessageMutator interface {
	// SetESV sets the specified ESV.
	SetESV(esv ESV) Message
	// SetDEOJ sets a destination object code.
	SetDEOJ(code ObjectCode) Message
	// AddProperty adds a property to the message.
	AddProperty(prop PropertyData) Message
	// AddProperties adds all properties to the message.
	AddProperties(props ...PropertyData) Message
}

// messageInternal is an interface to represent a message internal.
type messageInternal interface {
	// ToProtocol returns the protocol message.
	ToProtocol() *protocol.Message
}

type message struct {
	*protocol.Message
}

// newMessageWithProtocolMessage returns a new message.
func newMessageWithProtocolMessage(protoMsg *protocol.Message) *message {
	return &message{
		Message: protoMsg,
	}
}

// NewMessage returns a new message.
func NewMessage() Message {
	return newMessageWithProtocolMessage(protocol.NewMessage())
}

// NewMessageWith returns a new message of the specified parameters.
func NewMessageWith(objCode ObjectCode, esv ESV, props ...PropertyData) Message {
	return NewMessage().SetESV(esv).SetDEOJ(objCode).AddProperties(props...)
}

// SetESV sets the specified ESV.
func (msg *message) SetESV(esv ESV) Message {
	msg.Message.SetESV(esv)
	return msg
}

// SetDEOJ sets a destination object code.
func (msg *message) SetDEOJ(code ObjectCode) Message {
	msg.Message.SetDEOJ(code)
	return msg
}

// AddProperty adds a property to the message.
func (msg *message) AddProperty(prop PropertyData) Message {
	msg.Message.AddProperty(newProtocolPropertyFrom(prop))
	return msg
}

// AddProperties adds all properties to the message.
func (msg *message) AddProperties(props ...PropertyData) Message {
	for _, prop := range props {
		msg.Message.AddProperty(newProtocolPropertyFrom(prop))
	}
	return msg
}

// Properties returns the all properties of the message.
func (msg *message) Properties() []PropertyData {
	protoProps := msg.Message.Properties()
	props := make([]PropertyData, len(protoProps))
	for n, protoProp := range protoProps {
		props[n] = protoProp
	}
	return props
}

// Property returns the n-th property of the message.
func (msg *message) Property(n int) (PropertyData, bool) {
	props := msg.Properties()
	if n < 0 || n >= len(props) {
		return nil, false
	}
	return props[n], true
}

// ToProtocol returns the protocol message.
func (msg *message) ToProtocol() *protocol.Message {
	return msg.Message
}
