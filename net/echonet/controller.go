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

// NewController returns a new controller.
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

// GetLocalNode returns a local contoroller node.
func (ctrl *Controller) GetLocalNode() *LocalNode {
	return ctrl.LocalNode
}

// GetNodes returns found nodes.
func (ctrl *Controller) GetNodes() []*RemoteNode {
	return ctrl.foundNodes
}

// GetNode returns a node which has the specified address.
func (ctrl *Controller) GetNode(addr string) (*RemoteNode, error) {
	for _, node := range ctrl.GetNodes() {
		if node.GetAddress() == addr {
			return node, nil
		}
	}
	return nil, fmt.Errorf(errorNodeNotFound, addr)
}

// FindObject returns a object which has the specified object code.
func (ctrl *Controller) FindObject(code ObjectCode) (*Object, error) {
	for _, node := range ctrl.GetNodes() {
		obj, err := node.FindObject(code)
		if err == nil {
			return obj, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// FindDevice returns a device object which has the specified object code.
func (ctrl *Controller) FindDevice(code ObjectCode) (*Device, error) {
	for _, node := range ctrl.GetNodes() {
		dev, err := node.FindDevice(code)
		if err == nil {
			return dev, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// FindProfile returns a profile object which has the specified object code.
func (ctrl *Controller) FindProfile(code ObjectCode) (*Profile, error) {
	for _, node := range ctrl.GetNodes() {
		prof, err := node.FindProfile(code)
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
	if err := ctrl.Clear(); err != nil {
		return err
	}
	if err := ctrl.LocalNode.Start(); err != nil {
		return err
	}
	return nil
}

// Stop stop the controller.
func (ctrl *Controller) Stop() error {
	if err := ctrl.LocalNode.Stop(); err != nil {
		return err
	}
	return nil
}
