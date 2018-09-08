// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type DumpController struct {
	*echonet.Controller
}

func NewDumpController() *DumpController {
	c := &DumpController{
		Controller: echonet.NewController(),
	}
	return c
}

func (ctrl *DumpController) RequestMessageReceived(msg *protocol.Message) {
	fmt.Printf("%s : %s\n", msg.From.String(), hex.EncodeToString(msg.Bytes()))
}

func (ctrl *DumpController) NewNodeFound(*echonet.RemoteNode) {

}
