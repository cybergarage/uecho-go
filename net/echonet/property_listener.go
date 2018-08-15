// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/echonet-go/net/echonet/protocol"
)

// PropertyRequestListener is an instance for Echonet requests.
type PropertyRequestListener interface {
	PropertyRequestReceived(obj *Object, esv protocol.ESV, prop *protocol.Property) error
}