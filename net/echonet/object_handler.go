// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// ObjectRequestHandler is called when a property request is received.
// The node returns the standard responses of Echonet when the listener function returns no error.
// Otherwise, the node does not return any responses when the listener function returns an error.
type ObjectRequestHandler func(obj Object, esv protocol.ESV, prop protocol.Property) error
