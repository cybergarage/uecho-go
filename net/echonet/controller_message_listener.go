// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"slices"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	logControllerListenerFormat = "controller::OnMessage : %s"
)

func (ctrl *controller) isSelfMessage(msg *protocol.Message) bool {
	msgNode := newRemoteNodeWithRequestMessage(msg)
	for _, server := range ctrl.MulticastManager().Servers {
		port, err := server.Port()
		if err != nil {
			continue
		}
		if msgNode.Port() != port {
			continue
		}
		addr, err := server.Address()
		if err != nil {
			continue
		}
		if msgNode.Address() != addr {
			continue
		}
		return true
	}
	return false
}

// OnMessage is a listener of the local node.
func (ctrl *controller) OnMessage(msg *protocol.Message) error {
	if !ctrl.SelfMessageEnabled() {
		if ctrl.isSelfMessage(msg) {
			return nil
		}
	}

	// log.Trace(logControllerListenerFormat, msg.String())

	// NodeProfile message ?
	isNodeProfileMessage := func(msg *protocol.Message) bool {
		if !msg.ESV().IsNotification() && !msg.ESV().IsReadResponse() {
			return false
		}
		if !isNodeProfileObjectCode(msg.DEOJ()) {
			return false
		}
		return true
	}

	if isNodeProfileMessage(msg) {
		ctrl.parseNodeProfileMessage(msg)
	}

	if ctrl.controllerListener != nil {
		ctrl.controllerListener.ControllerMessageReceived(msg)
	}

	return nil
}

// parseNodeProfileMessage parses the specified message to check new objects.
func (ctrl *controller) parseNodeProfileMessage(msg *protocol.Message) {
	if !ctrl.SelfMessageEnabled() {
		if ctrl.isSelfMessage(msg) {
			return
		}
	}

	node, err := newRemoteNodeWithInstanceListMessageAndConfig(msg, ctrl.TransportConfig())
	if err != nil {
		return
	}

	ctrl.addNode(node)
}

// addNode adds a specified node if the node is not added.
func (ctrl *controller) addNode(notifyNode Node) bool {
	if slices.ContainsFunc(ctrl.foundNodes, notifyNode.Equals) {
		return false
	}

	ctrl.foundNodes = append(ctrl.foundNodes, notifyNode)

	if ctrl.controllerListener != nil {
		ctrl.controllerListener.ControllerNewNodeFound(notifyNode)
	}

	return true
}
