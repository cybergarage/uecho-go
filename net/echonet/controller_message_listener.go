// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// MessageReceived is an override message listener of LocalNode to get the announce messages.
func (ctrl *Controller) MessageReceived(msg *protocol.Message) {
	// Ignore own messages
	msgNode := NewRemoteNodeWithRequestMessage(msg)
	if msgNode.Equals(ctrl) {
		return
	}

	// NodeProfile message ?
	if msg.IsNotificationResponse() || msg.IsReadResponse() {
		msgDsgObj := msg.GetDestinationObjectCode()
		if isNodeProfileObjectCode(msgDsgObj) {
			ctrl.parseNodeProfileMessage(msg)
		}
	}

	ctrl.LocalNode.MessageReceived(msg)

	if ctrl.controllerListener != nil {
		ctrl.controllerListener.MessageRequestReceived(msg)
	}
}

// parseNodeProfileMessage parses the specified message to check new objects.
func (ctrl *Controller) parseNodeProfileMessage(msg *protocol.Message) {
	node, err := NewRemoteNodeWithInstanceListMessage(msg)
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
		ctrl.controllerListener.NewNodeFound(notifyNode)
	}

	return true
}
