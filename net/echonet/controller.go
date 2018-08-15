// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"

	"github.com/cybergarage/echonet-go/net/echonet/protocol"
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
	*LocalNode
	foundNodes []*RemoteNode

	listener ControllerListener
}

// NewController returns a new contorller.
func NewController() *Controller {
	ctrl := &Controller{
		LocalNode:  NewLocalNode(),
		foundNodes: make([]*RemoteNode, 0),
		listener:   nil,
	}

	ctrl.SetMessageListener(ctrl)

	return ctrl
}

// SetListener sets a listener to receive the Echonet messages.
func (ctrl *Controller) SetListener(l ControllerListener) {
	ctrl.listener = l
}

// GetNodes returns found nodes
func (ctrl *Controller) GetNodes() []*RemoteNode {
	return ctrl.foundNodes
}

// GetObject returns the specified object.
func (ctrl *Controller) GetObject(code uint) (*Object, error) {
	for _, node := range ctrl.GetNodes() {
		obj, err := node.GetObject(code)
		if err == nil {
			return obj, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// GetDevice returns the specified device object.
func (ctrl *Controller) GetDevice(code uint) (*Device, error) {
	for _, node := range ctrl.GetNodes() {
		dev, err := node.GetDevice(code)
		if err == nil {
			return dev, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// GetProfile returns the specified profile object.
func (ctrl *Controller) GetProfile(code uint) (*Profile, error) {
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
