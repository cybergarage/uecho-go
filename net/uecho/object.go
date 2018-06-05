// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/uecho/encoding"
	"github.com/cybergarage/uecho-go/net/uecho/protocol"
)

const (
	errorParentNodeNotFound = "Parent node not found"
)

// Object is an instance for Echonet object.
type Object struct {
	Code             []byte
	Properties       []*Property
	annoPropMapSize  int
	annoPropMapBytes []byte
	setPropMapSize   int
	setPropMapBytes  []byte
	getPropMapSize   int
	getPropMapBytes  []byte

	parentNode *Node
}

// NewObject returns a new object.
func NewObject() *Object {
	obj := &Object{
		Code:       make([]byte, 0),
		Properties: make([]*Property, 0),
		parentNode: nil,
	}
	return obj
}

// GetCode returns the code.
func (obj *Object) GetCode() uint {
	return encoding.ByteToInteger(obj.Code)
}

// AnnounceMessage announces a message.
func (obj *Object) AnnounceMessage(msg *protocol.Message) error {
	if obj.parentNode == nil {
		return fmt.Errorf(errorParentNodeNotFound)
	}
	msg.SetSourceObjectCode(obj.GetCode())
	return obj.parentNode.AnnounceMessage(msg)
}
