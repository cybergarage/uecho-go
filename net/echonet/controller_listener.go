// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// ControllerListener is a listener for Echonet messages.
type ControllerListener interface {
	ControllerMessageReceived(*protocol.Message)
	ControllerNewNodeFound(*RemoteNode)
}
