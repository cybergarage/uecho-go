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

type Message = protocol.Message

// NewMessage returns a new message.
func NewMessage() *Message {
	return protocol.NewMessage()
}

// NewMessageWithParameters returns a new message of the specified parameters.
func NewMessageWithParameters(objCode ObjectCode, esv ESV, props []*Property) *Message {
	msg := NewMessage()
	msg.SetESV(esv)
	msg.SetDestinationObjectCode(objCode)
	for _, prop := range props {
		msg.AddProperty(prop.toProtocolProperty())
	}
	return msg
}
