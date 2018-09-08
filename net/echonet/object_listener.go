// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// ObjectListener is an interface for Echonet requests.
type ObjectListener interface {
	MessageRequestReceived(*protocol.Message)
	PropertyRequestReceived(obj *Object, esv protocol.ESV, prop *protocol.Property) error
}
