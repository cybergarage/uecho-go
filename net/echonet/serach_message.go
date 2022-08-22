// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// NewSearchMessage returns a new search message.
func NewSearchMessage() *protocol.Message {
	msg := protocol.NewMessage()

	msg.SetESV(protocol.ESVReadRequest)

	msg.SetSEOJ(NodeProfileObject)
	msg.SetDestinationObjectCode(NodeProfileObject)

	prop := protocol.NewProperty()
	prop.SetCode(NodeProfileClassSelfNodeInstanceListS)
	prop.SetData(make([]byte, 0))
	msg.AddProperty(prop)

	return msg
}
