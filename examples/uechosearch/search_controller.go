// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type SearchController struct {
	*echonet.Controller
}

func NewSearchController() *SearchController {
	c := &SearchController{
		Controller: echonet.NewController(),
	}
	return c
}

func (ctrl *SearchController) ControllerMessageReceived(msg *protocol.Message) {
	// log.Info("%s : %s\n", msg.From.String(), hex.EncodeToString(msg.Bytes()))
}

func (ctrl *SearchController) ControllerNewNodeFound(*echonet.RemoteNode) {
}
