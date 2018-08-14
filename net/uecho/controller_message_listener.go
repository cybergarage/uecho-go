// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

// MessageReceived is an override message listener of LocalNode to get the announce messages.
func (ctrl *Controller) MessageReceived(msg *protocol.Message) {
	// NodeProfile message ?
	if msg.IsNotificationResponse() || msg.IsReadResponse() {
		msgDsgObj := msg.GetDestinationObjectCode()
		if isNodeProfileObjectCode(msgDsgObj) {
			ctrl.parseNodeProfileMessage(msg)
		}
	}

	ctrl.LocalNode.MessageReceived(msg)
}

// parseNodeProfileMessage parses the specified message to check new objects.
func (ctrl *Controller) parseNodeProfileMessage(msg *protocol.Message) {
	node, err := NewRemoteNodeWithNotificationMessage(msg)
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

	return true
}
