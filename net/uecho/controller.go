// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

// ControllerListener is a listener for Echonet messages.
type ControllerListener interface {
	ControllerMessageReceived(msg *protocol.Message)
}

// Controller is an instance for Echonet controller.
type Controller struct {
	node     *Node
	Nodes    []*Node
	listener ControllerListener
}

// NewController returns a new contorller.
func NewController() *Controller {
	ctrl := &Controller{
		node:     NewNode(),
		Nodes:    make([]*Node, 0),
		listener: nil,
	}
	return ctrl
}

// SetListener sets a listener to receive the Echonet messages.
func (ctrl *Controller) SetListener(l ControllerListener) {
	ctrl.listener = l
}

// Start starts the controller.
func (ctrl *Controller) Start() error {
	err := ctrl.node.Start()
	if err != nil {
		return err
	}

	return nil
}

// Stop stop the controller.
func (ctrl *Controller) Stop() error {
	err := ctrl.node.Stop()
	if err != err {
		return nil
	}

	return nil
}
