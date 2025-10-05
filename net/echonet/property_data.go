// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// PropertyData is an interface for Echonet property data.
type PropertyData interface {
	// Code returns the property code.
	Code() PropertyCode
	// Size return the property data size.
	Size() int
	// Data returns the property data.
	Data() []byte
}

type propData struct {
	code PropertyCode
	data []byte
}

// NewPropertyData returns a new property data.
func NewPropertyDataWith(code PropertyCode, data []byte) PropertyData {
	return &propData{
		code: code,
		data: data,
	}
}

// Code returns the property code.
func (pd *propData) Code() PropertyCode {
	return pd.code
}

// Size return the property data size.
func (pd *propData) Size() int {
	return len(pd.data)
}

// Data returns the property data.
func (pd *propData) Data() []byte {
	return pd.data
}

// newProtocolPropertyFrom returns a new protocol property from the specified property data.
func newProtocolPropertyFrom(prop PropertyData) *protocol.Property {
	newProp := protocol.NewProperty()
	newProp.SetCode(prop.Code())
	newProp.SetData(prop.Data())
	return newProp
}
