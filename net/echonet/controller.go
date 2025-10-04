// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"context"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	// TIDMin is the minimum value of transaction ID.
	TIDMin = 0
	// TIDMax is the maximum value of transaction ID.
	TIDMax = 65535
	// DefaultResponseTimeout is the default timeout value to wait for a response message.
	DefaultResponseTimeout = time.Duration(3) * time.Second
)

// Controller represents the Echonet controller.
type Controller interface {
	// SetConfig sets a configuration.
	SetConfig(*Config)
	// SetListener sets a listener to receive the Echonet messages.
	SetListener(ControllerListener)
	// Addresses returns the local addresses that this controller is bound to.
	Addresses() []string
	// Search searches echonet nodes until the context is done.
	Search(ctx context.Context) error
	// Nodes returns discovered nodes.
	Nodes() []*RemoteNode
	// LookupNode returns a node which has the specified address.
	LookupNode(addr string) (*RemoteNode, bool)
	// SendMessage sends a message to the node.
	SendMessage(ctx context.Context, dstNode Node, msg *Message) error
	// PostMessage posts a message to the node, and wait the response message.
	PostMessage(ctx context.Context, dstNode Node, msg *Message) (*Message, error)
	// Start starts the controller.
	Start() error
	// Stop stops the controller.
	Stop() error
}

type controller struct {
	*LocalNode

	foundNodes []*RemoteNode

	controllerListener ControllerListener
}

// NewController returns a new controller.
func NewController() *controller {
	return newController()
}

func newController() *controller {
	ctrl := &controller{
		LocalNode:          NewLocalNode(),
		foundNodes:         make([]*RemoteNode, 0),
		controllerListener: nil,
	}

	ctrl.LocalNode.SetListener(ctrl)

	return ctrl
}

// SetListener sets a listener to receive the Echonet messages.
func (ctrl *controller) SetListener(l ControllerListener) {
	ctrl.controllerListener = l
}

// Nodes returns found nodes.
func (ctrl *controller) Nodes() []*RemoteNode {
	return ctrl.foundNodes
}

// LookupNode returns a node which has the specified address.
func (ctrl *controller) LookupNode(addr string) (*RemoteNode, bool) {
	for _, node := range ctrl.Nodes() {
		if node.Address() == addr {
			return node, true
		}
	}
	return nil, false
}

// SearchAllObjectsWithESV searches all specified objects.
func (ctrl *controller) SearchAllObjectsWithESV(esv protocol.ESV) error {
	msg := NewSearchMessage()
	msg.SetESV(esv)
	return ctrl.AnnounceMessage(msg)
}

// SearchAllObjects searches all objects.
func (ctrl *controller) SearchAllObjects() error {
	return ctrl.SearchAllObjectsWithESV(protocol.ESVReadRequest)
}

// SearchObjectWithESV searches a specified object.
func (ctrl *controller) SearchObjectWithESV(code ObjectCode, esv protocol.ESV) error {
	msg := NewSearchMessage()
	msg.SetESV(esv)
	msg.SetDEOJ(code)
	return ctrl.AnnounceMessage(msg)
}

// SearchObject searches a specified object.
func (ctrl *controller) SearchObject(code ObjectCode) error {
	return ctrl.SearchObjectWithESV(code, protocol.ESVReadRequest)
}

// Search searches echonet nodes until the context is done.
func (ctrl *controller) Search(ctx context.Context) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultResponseTimeout)
		defer cancel()
	}
	err := ctrl.SearchAllObjects()
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(time.Millisecond * 200)
		}
	}
}

// Clear clears all found nodes.
func (ctrl *controller) Clear() error {
	ctrl.foundNodes = make([]*RemoteNode, 0)
	return nil
}

// Start starts the controller.
func (ctrl *controller) Start() error {
	if err := ctrl.Clear(); err != nil {
		return err
	}
	if err := ctrl.LocalNode.Start(); err != nil {
		return err
	}
	return nil
}

// Stop stop the controller.
func (ctrl *controller) Stop() error {
	if err := ctrl.LocalNode.Stop(); err != nil {
		return err
	}
	return nil
}
