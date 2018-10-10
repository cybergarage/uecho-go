// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type PostController struct {
	*echonet.Controller
}

func NewPostController() *PostController {
	c := &PostController{
		Controller: echonet.NewController(),
	}
	return c
}

func (ctrl *PostController) ControllerMessageReceived(msg *protocol.Message) {
	fmt.Printf("%s : %s\n", msg.From.String(), hex.EncodeToString(msg.Bytes()))
}

func (ctrl *PostController) ControllerNewNodeFound(*echonet.RemoteNode) {
}
