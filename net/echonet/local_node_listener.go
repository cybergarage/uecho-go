// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// LocalNodeListener is an instance of the listner.
type LocalNodeListener interface {
	// NodeMessageReceived is first called when a new message is received.
	// The node returns the standard responses of Echonet when the listener function returns no error.
	// Otherwise, the node does not return any responses when the listener function returns an error.
	NodeMessageReceived(*protocol.Message) error
}
