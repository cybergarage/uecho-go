// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	logControllerListenerFormat = "Controller::NodeMessageReceived : %s"
)

func (ctrl *Controller) isSelfMessage(msg *protocol.Message) bool {
	msgNode := NewRemoteNodeWithRequestMessage(msg)
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

// NodeMessageReceived is a listener of the local node.
func (ctrl *Controller) NodeMessageReceived(msg *protocol.Message) error {
	if !ctrl.IsSelfMessageEnabled() {
		if ctrl.isSelfMessage(msg) {
			return nil
		}
	}

	// log.Trace(logControllerListenerFormat, msg.String())

	// NodeProfile message ?
	isNodeProfileMessage := func(msg *protocol.Message) bool {
		if !msg.IsNotification() && !msg.IsReadResponse() {
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
func (ctrl *Controller) parseNodeProfileMessage(msg *protocol.Message) {
	if !ctrl.IsSelfMessageEnabled() {
		if ctrl.isSelfMessage(msg) {
			return
		}
	}

	node, err := NewRemoteNodeWithInstanceListMessageAndConfig(msg, ctrl.TransportConfig)
	if err != nil {
		return
	}

	ctrl.addNode(node)
}

// addNode adds a specified node if the node is not added.
func (ctrl *Controller) addNode(notifyNode *RemoteNode) bool {
	for _, node := range ctrl.foundNodes {
		if notifyNode.Equals(node) {
			return false
		}
	}

	ctrl.foundNodes = append(ctrl.foundNodes, notifyNode)

	if ctrl.controllerListener != nil {
		ctrl.controllerListener.ControllerNewNodeFound(notifyNode)
	}

	return true
}
