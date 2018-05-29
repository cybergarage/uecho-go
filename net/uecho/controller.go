// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"github.com/cybergarage/uecho-go/net/uecho/protocol"
	"github.com/cybergarage/uecho-go/net/uecho/std"
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

// AnnounceMessage announces a message.
func (ctrl *Controller) AnnounceMessage(msg *protocol.Message) error {
	/*
	      nodeProfObj = uecho_node_getnodeprofileclassobject(ctrl->node);
	      if (!nodeProfObj)
	   	 return false;

	      uecho_message_settid(msg, uecho_controller_getnexttid(ctrl));

	      return uecho_object_announcemessage(nodeProfObj, msg);
	*/
}

// SearchAllObjectsWithESV searches all specified objects.
func (ctrl *Controller) SearchAllObjectsWithESV(esv protocol.ESV) error {
	msg := std.NewSearchMessage()
	msg.SetESV(esv)
	return ctrl.AnnounceMessage(msg)
}

// SearchAllObjects searches all objects.
func (ctrl *Controller) SearchAllObjects(esv protocol.ESV) error {
	return ctrl.SearchAllObjectsWithESV(protocol.ESVReadRequest)
}

// SearchObjectWithESV searches a specified object.
func (ctrl *Controller) SearchObjectWithESV(code byte, esv protocol.ESV) error {
	msg := std.NewSearchMessage()
	msg.SetESV(esv)
	msg.SetDestinationObjectCode(code)
	return ctrl.AnnounceMessage(msg)
}

// SearchObject searches a specified object.
func (ctrl *Controller) SearchObject(code byte) error {
	return ctrl.SearchObjectWithESV(code, protocol.ESVReadRequest)

}
