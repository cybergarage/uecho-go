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
	msg := NewMessage()
	msg.SetESV(esv)
	msg.SetDEOJ(objCode)
	for _, prop := range props {
		msg.AddProperty(prop)
	}
	return msg
}

// AddProperty adds a property to the message.
func (msg *Message) AddProperty(prop *Property) {
	msg.Message.AddProperty(prop.toProtocolProperty())
}

// AddProperties adds all properties to the message.
func (msg *Message) AddProperties(props []*Property) {
	for _, prop := range props {
		msg.Message.AddProperty(prop.toProtocolProperty())
	}
}

// protocolMessage returns the protocol message.
func (msg *Message) protocolMessage() *protocol.Message {
	return msg.Message
}
