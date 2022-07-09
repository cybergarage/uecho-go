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

// NodeMessageReceived is a listener of the local node.
func (ctrl *Controller) NodeMessageReceived(msg *protocol.Message) error {
	// Ignore own messages
	msgNode := NewRemoteNodeWithRequestMessage(msg)
	if msgNode.Equals(ctrl) {
		return nil
	}

	// log.Trace(logControllerListenerFormat, msg.String()))

	// NodeProfile message ?
	isNodeProfileMessage := func(msg *protocol.Message) bool {
		if !msg.IsNotification() && !msg.IsReadResponse() {
			return false
		}
		if !isNodeProfileObjectCode(msg.GetDestinationObjectCode()) {
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
	node, err := NewRemoteNodeWithInstanceListMessageAndConfig(msg, ctrl.GetConfig())
	if err != nil {
		return
	}

	if node.Equals(ctrl) {
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
