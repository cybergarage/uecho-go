// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// PropertyDataCode is a type for Echonet property data code.
type PropertyDataOption func(*propData) error

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

// WithPropertyDataCode sets a property code to the property data.
func WithPropertyDataCode(code PropertyCode) PropertyDataOption {
	return func(pd *propData) error {
		pd.code = code
		return nil
	}
}

// WithPropertyDataBytes sets data bytes to the property data.
func WithPropertyDataBytes(data []byte) PropertyDataOption {
	return func(pd *propData) error {
		pd.data = make([]byte, len(data))
		copy(pd.data, data)
		return nil
	}
}

// NewPropertyData returns a new property data.
func NewPropertyData() PropertyData {
	return newPropertyData()
}

func newPropertyData() *propData {
	return &propData{
		code: 0,
		data: make([]byte, 0),
	}
}

// NewPropertyDataWith returns a new property data with the specified options.
func NewPropertyDataWith(opts ...PropertyDataOption) (PropertyData, error) {
	pd := newPropertyData()
	for _, opt := range opts {
		if err := opt(pd); err != nil {
			return nil, err
		}
	}
	return pd, nil
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
