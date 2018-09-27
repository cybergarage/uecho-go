// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

// A MessageHandler represents a handler for message.
type MessageHandler interface {
	ProtocolMessageReceived(*Message) (*Message, error)
}
