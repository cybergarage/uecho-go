// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	TIDMin = 0
	TIDMax = 65535
)

const (
	errorNodeNotFound = "Node (%s) not found"
)

// Controller is an instance for Echonet controller.
type Controller struct {
	*LocalNode
	foundNodes []*RemoteNode

	controllerListener ControllerListener
}

// NewController returns a new contorller.
func NewController() *Controller {
	ctrl := &Controller{
		LocalNode:          NewLocalNode(),
		foundNodes:         make([]*RemoteNode, 0),
		controllerListener: nil,
	}

	ctrl.LocalNode.SetListener(ctrl)

	return ctrl
}

// SetListener sets a listener to receive the Echonet messages.
func (ctrl *Controller) SetListener(l ControllerListener) {
	ctrl.controllerListener = l
}

// GetNodes returns found nodes
func (ctrl *Controller) GetNodes() []*RemoteNode {
	return ctrl.foundNodes
}

// GetNode returns a node which has the specified address
func (ctrl *Controller) GetNode(addr string) (*RemoteNode, error) {
	for _, node := range ctrl.GetNodes() {
		if node.GetAddress() == addr {
			return node, nil
		}
	}
	return nil, fmt.Errorf(errorNodeNotFound, addr)
}

// GetObject returns a object which has the specified object code.
func (ctrl *Controller) GetObject(code ObjectCode) (*Object, error) {
	for _, node := range ctrl.GetNodes() {
		obj, err := node.GetObject(code)
		if err == nil {
			return obj, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// GetDevice returns a device object which has the specified object code.
func (ctrl *Controller) GetDevice(code ObjectCode) (*Device, error) {
	for _, node := range ctrl.GetNodes() {
		dev, err := node.GetDevice(code)
		if err == nil {
			return dev, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// GetProfile returns a profile object which has the specified object code.
func (ctrl *Controller) GetProfile(code ObjectCode) (*Profile, error) {
	for _, node := range ctrl.GetNodes() {
		prof, err := node.GetProfile(code)
		if err == nil {
			return prof, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
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
func (ctrl *Controller) SearchObjectWithESV(code ObjectCode, esv protocol.ESV) error {
	msg := NewSearchMessage()
	msg.SetESV(esv)
	msg.SetDestinationObjectCode(code)
	return ctrl.AnnounceMessage(msg)
}

// SearchObject searches a specified object.
func (ctrl *Controller) SearchObject(code ObjectCode) error {
	return ctrl.SearchObjectWithESV(code, protocol.ESVReadRequest)
}

// Clear clears all found nodes.
func (ctrl *Controller) Clear() error {
	ctrl.foundNodes = make([]*RemoteNode, 0)
	return nil
}

// Start starts the controller.
func (ctrl *Controller) Start() error {
	err := ctrl.Clear()
	if err != nil {
		return err
	}

	err = ctrl.LocalNode.Start()
	if err != nil {
		return err
	}

	return nil
}

// Stop stop the controller.
func (ctrl *Controller) Stop() error {
	err := ctrl.LocalNode.Stop()
	if err != err {
		return nil
	}

	return nil
}
