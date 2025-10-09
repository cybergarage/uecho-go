// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// MessageOptions is a function type to set options in Message.
type MessageOptions func(*message)

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
	Properties() []Property
	// Property returns the n-th property of the message.
	Property(n int) (Property, bool)
	// messageInternal is an interface to represent a message internal.
	messageInternal
}

// messageInternal is an interface to represent a message internal.
type messageInternal interface {
	// ToProtocol returns the protocol message.
	ToProtocol() *protocol.Message
}

type message struct {
	*protocol.Message
}

// WithMessageDEOJ sets a destination object code.
func WithMessageDEOJ(code ObjectCode) MessageOptions {
	return func(msg *message) {
		msg.SetDEOJ(code)
	}
}

// WithMessageSEOJ sets a source object code.
func WithMessageESV(esv ESV) MessageOptions {
	return func(msg *message) {
		msg.SetESV(esv)
	}
}

// WithMessageProperties sets properties to the message.
func WithMessageProperties(props ...Property) MessageOptions {
	return func(msg *message) {
		msg.AddProperties(props...)
	}
}

// NewMessage returns a new message.
func NewMessage(opts ...MessageOptions) Message {
	msg := newMessage()
	for _, opt := range opts {
		opt(msg)
	}
	return msg
}

// newMessageWithProtocolMessage returns a new message.
func newMessageWithProtocolMessage(protoMsg *protocol.Message) *message {
	return &message{
		Message: protoMsg,
	}
}

func newMessage() *message {
	return newMessageWithProtocolMessage(protocol.NewMessage())
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
func (msg *message) AddProperty(prop Property) Message {
	msg.Message.AddProperty(newProtocolPropertyFrom(prop))
	return msg
}

// AddProperties adds all properties to the message.
func (msg *message) AddProperties(props ...Property) Message {
	for _, prop := range props {
		msg.Message.AddProperty(newProtocolPropertyFrom(prop))
	}
	return msg
}

// Properties returns the all properties of the message.
func (msg *message) Properties() []Property {
	protoProps := msg.Message.Properties()
	props := make([]Property, len(protoProps))
	for n, protoProp := range protoProps {
		props[n] = NewProperty(
			WithPropertyCode(PropertyCode(protoProp.Code())),
			WithPropertyData(protoProp.Data()),
		)
	}
	return props
}

// Property returns the n-th property of the message.
func (msg *message) Property(n int) (Property, bool) {
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
