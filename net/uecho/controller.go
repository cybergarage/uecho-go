// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

const (
	TIDMin = 0
	TIDMax = 65535
)

// ControllerListener is a listener for Echonet messages.
type ControllerListener interface {
	ControllerMessageReceived(msg *protocol.Message)
}

// Controller is an instance for Echonet controller.
type Controller struct {
	node       *Node
	foundNodes []*Node

	lastTID  uint
	listener ControllerListener
}

// NewController returns a new contorller.
func NewController() *Controller {
	ctrl := &Controller{
		node:       NewNode(),
		foundNodes: make([]*Node, 0),
		lastTID:    TIDMin,
		listener:   nil,
	}
	return ctrl
}

// SetListener sets a listener to receive the Echonet messages.
func (ctrl *Controller) SetListener(l ControllerListener) {
	ctrl.listener = l
}

// GetNodes returns found nodes
func (ctrl *Controller) GetNodes() []*Node {
	return ctrl.foundNodes
}

// getNextTID returns a next TID.
func (ctrl *Controller) getNextTID() uint {
	if TIDMax <= ctrl.lastTID {
		ctrl.lastTID = TIDMin
	} else {
		ctrl.lastTID++
	}
	return ctrl.lastTID
}

// AnnounceMessage announces a message.
func (ctrl *Controller) AnnounceMessage(msg *protocol.Message) error {
	nodeProf, err := ctrl.node.GetNodeProfile()
	if err != nil {
		return err
	}
	msg.SetTID(ctrl.getNextTID())
	return nodeProf.AnnounceMessage(msg)
}

// SearchAllObjectsWithESV searches all specified objects.
func (ctrl *Controller) SearchAllObjectsWithESV(esv protocol.ESV) error {
	msg := NewSearchMessage()
	msg.SetESV(esv)
	return ctrl.AnnounceMessage(msg)
}

// SearchAllObjects searches all objects.
func (ctrl *Controller) SearchAllObjects() error {
	return ctrl.SearchAllObjectsWithESV(protocol.ESVReadRequest)
}

// SearchObjectWithESV searches a specified object.
func (ctrl *Controller) SearchObjectWithESV(code uint, esv protocol.ESV) error {
	msg := NewSearchMessage()
	msg.SetESV(esv)
	msg.SetDestinationObjectCode(code)
	return ctrl.AnnounceMessage(msg)
}

// SearchObject searches a specified object.
func (ctrl *Controller) SearchObject(code uint) error {
	return ctrl.SearchObjectWithESV(code, protocol.ESVReadRequest)
}

// Clear clears all found nodes.
func (ctrl *Controller) Clear() error {
	ctrl.foundNodes = make([]*Node, 0)
	return nil
}

// Start starts the controller.
func (ctrl *Controller) Start() error {
	err := ctrl.Clear()
	if err != nil {
		return err
	}

	err = ctrl.node.Start()
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
