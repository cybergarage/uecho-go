// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	errorNodeRequestInvalidDestinationNode = "Invalid Destination Node : %v"
)

type Message struct {
	*protocol.Message
}

// newMessageWithProtocolMessage returns a new message.
func newMessageWithProtocolMessage(protoMsg *protocol.Message) *Message {
	return &Message{
		Message: protoMsg,
	}
}

// NewMessage returns a new message.
func NewMessage() *Message {
	return newMessageWithProtocolMessage(protocol.NewMessage())
}

// NewMessageWithParameters returns a new message of the specified parameters.
func NewMessageWithParameters(objCode ObjectCode, esv ESV, props []*Property) *Message {
	return NewMessage().SetESV(esv).SetDEOJ(objCode).AddProperties(props)
}

// SetESV sets the specified ESV.
func (msg *Message) SetESV(esv ESV) *Message {
	msg.Message.SetESV(esv)
	return msg
}

// SetDEOJ sets a destination object code.
func (msg *Message) SetDEOJ(code ObjectCode) *Message {
	msg.Message.SetDEOJ(code)
	return msg
}

// AddProperty adds a property to the message.
func (msg *Message) AddProperty(prop *Property) *Message {
	msg.Message.AddProperty(prop.toProtocolProperty())
	return msg
}

// AddProperties adds all properties to the message.
func (msg *Message) AddProperties(props []*Property) *Message {
	for _, prop := range props {
		msg.Message.AddProperty(prop.toProtocolProperty())
	}
	return msg
}

// protocolMessage returns the protocol message.
func (msg *Message) protocolMessage() *protocol.Message {
	return msg.Message
}
