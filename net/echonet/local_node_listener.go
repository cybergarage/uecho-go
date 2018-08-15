// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/echonet-go/net/echonet/protocol"
)

// LocalNodeListener is an instance of the listner.
type LocalNodeListener interface {
	protocol.MessageListener
}
