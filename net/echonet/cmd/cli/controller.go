// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// Controller represents an Echonet Lite controller.
type Controller struct {
	*echonet.Controller
}

// NewController returns a new controller.
func NewController() *Controller {
	c := &Controller{
		Controller: echonet.NewController(),
	}
	return c
}

// ControllerMessageReceived is called when a message is received.
func (ctrl *Controller) ControllerMessageReceived(msg *protocol.Message) {
	// log.Infof("%s : %s\n", msg.From.String(), hex.EncodeToString(msg.Bytes()))
}

// ControllerNewNodeFound is called when a new node is found.
func (ctrl *Controller) ControllerNewNodeFound(*echonet.RemoteNode) {
}
