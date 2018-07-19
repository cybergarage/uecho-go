// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

const (
	errorObjectNotFound              = "Object (%d) not found"
	errorObjectProfileObjectNotFound = "Object profile object not found"
)

// Node is an instance for Echonet node.
type Node struct {
	devices  []*Device
	profiles []*Profile
	Address  string
	server   *server
}

// NewNode returns a new object.
func NewNode() *Node {
	node := &Node{
		devices: make([]*Device, 0),
		server:  newServer(),
	}

	node.AddProfile(NewNodeProfile())

	return node
}

// AddDevice adds a new device into the node.
func (node *Node) AddDevice(dev *Device) error {
	node.devices = append(node.devices, dev)
	return node.updateNodeProfile()
}

// GetDevices returns all device objects.
func (node *Node) GetDevices() []*Device {
	return node.devices
}

// GetDeviceByCode returns a specified device object.
func (node *Node) GetDeviceByCode(code uint) (*Device, error) {
	for _, obj := range node.devices {
		if obj.GetCode() == code {
			return obj, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// AddProfile adds a new profile object into the node.
func (node *Node) AddProfile(prof *Profile) error {
	node.profiles = append(node.profiles, prof)
	return node.updateNodeProfile()
}

// GetProfiles returns all profile objects.
func (node *Node) GetProfiles() []*Profile {
	return node.profiles
}

// GetProfilesByCode returns a specified profile object.
func (node *Node) GetProfilesByCode(code uint) (*Profile, error) {
	for _, prof := range node.profiles {
		if prof.GetCode() == code {
			return prof, nil
		}
	}
	return nil, fmt.Errorf(errorObjectNotFound, code)
}

// GetNodeProfileObject returns a specified object.
func (node *Node) GetNodeProfileObject() (*Profile, error) {
	return node.GetProfilesByCode(NodeProfileObject)
}

// AnnounceMessage announces a message.
func (node *Node) AnnounceMessage(msg *protocol.Message) error {
	return node.server.SendMessageAll(msg)
}

// GetAddress returns a IP address of the node.
func (node *Node) GetAddress() string {
	return node.Address
}

// AnnounceProperty announces a specified property.
func (node *Node) AnnounceProperty(prop *Property) error {
	msg := protocol.NewMessage()
	msg.SetESV(protocol.ESVNotification)
	msg.SetSourceObjectCode(NodeProfileObject)
	msg.SetDestinationObjectCode(NodeProfileObject)
	msg.AddProperty(prop.toProtocolProperty())

	return node.AnnounceMessage(msg)
}

// Announce announces the node
func (node *Node) Announce() error {
	nodePropObj, err := node.GetNodeProfileObject()
	if err != nil {
		return err
	}

	nodeProp, ok := nodePropObj.GetProperty(NodeProfileClassInstanceListNotification)
	if !ok {
		return fmt.Errorf(errorObjectProfileObjectNotFound)
	}

	return node.AnnounceProperty(nodeProp)
}

// SendMessage send a message to the node
func (node *Node) SendMessage(dstNode *Node, msg *protocol.Message) error {
	return node.server.SendMessage(dstNode.GetAddress(), msg)
}

// Start starts the node.
func (node *Node) Start() error {
	err := node.server.Start()
	if err != nil {
		return err
	}

	return nil
}

// Stop stop the node.
func (node *Node) Stop() error {
	err := node.server.Stop()
	if err != nil {
		return err
	}

	return nil
}

// updateNodeProfile updates the node profile in the node.
func (node *Node) updateNodeProfile() error {
	/*
		   Node *node;
		   Class *nodeCls;
		   Object *nodeObj;
		   byte *nodeClassList, *nodeInstanceList;
		   int nodeClassListCnt, nodeInstanceListCnt;
		   int nodeClassCnt, nodeInstanceCnt;
		   int idx;

		   if(!obj)
			 return false;

		   node = _object_getparentnode(obj);
		   if (!node)
			 return false;

		   // Class Properties

		   nodeClassList = (byte *)realloc(NULL, 1);
		   nodeClassListCnt = 0;
		   nodeClassCnt = 0;

		   for (nodeCls = _node_getclasses(node); nodeCls; nodeCls = _class_next(nodeCls)) {
			 nodeClassCnt++;

			 if (_class_isprofile(nodeCls))
			   continue;

			 nodeClassListCnt++;
			 nodeClassList = (byte *)realloc(nodeClassList, ((2 * nodeClassListCnt) + 1));
			 idx = (2 * (nodeClassListCnt - 1)) + 1;
			 nodeClassList[idx + 0] = _class_getclassgroupcode(nodeCls);
			 nodeClassList[idx + 1] = _class_getclasscode(nodeCls);
		   }

		   _nodeprofileclass_setclasscount(obj, nodeClassCnt);
		   _nodeprofileclass_setclasslist(obj, nodeClassListCnt, nodeClassList);

		   free(nodeClassList);

		   // Instance Properties

		   nodeInstanceList = (byte *)realloc(NULL, 1);
		   nodeInstanceListCnt = 0;
		   nodeInstanceCnt = 0;

		   for (nodeObj = _node_getobjects(node); nodeObj; nodeObj = _object_next(nodeObj)) {
			 if (_object_isprofile(nodeObj))
			   continue;

			 nodeInstanceCnt++;

			 nodeInstanceListCnt++;
			 nodeInstanceList = (byte *)realloc(nodeInstanceList, ((3 * nodeInstanceListCnt) + 1));
			 idx = (3 * (nodeInstanceListCnt - 1)) + 1;
			 nodeInstanceList[idx + 0] = _object_getclassgroupcode(nodeObj);
			 nodeInstanceList[idx + 1] = _object_getclasscode(nodeObj);
			 nodeInstanceList[idx + 2] = _object_getinstancecode(nodeObj);
		   }

		   _nodeprofileclass_setinstancecount(obj, nodeInstanceCnt);
		   _nodeprofileclass_setinstancelist(obj, nodeInstanceListCnt, nodeInstanceList);

		   free(nodeInstanceList);

		   return true;
	*/
	return nil
}
