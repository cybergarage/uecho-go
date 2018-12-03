// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// A UnicastHandler represents a listener for UnicastServer.
type UnicastHandler interface {
	protocol.MessageHandler
}
