// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

// ControllerListener is a listener for Echonet messages.
type ControllerListener interface {
	ControllerMessageReceived(msg *Message)
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
