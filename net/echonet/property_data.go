// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// propertyData is an interface for Echonet property data.
type propertyData interface {
	// Code returns the property code.
	Code() PropertyCode
	// Size return the property data size.
	Size() int
	// Data returns the property data.
	Data() []byte
}

// newProtocolPropertyFrom returns a new protocol property from the specified property data.
func newProtocolPropertyFrom(prop propertyData) *protocol.Property {
	newProp := protocol.NewProperty()
	newProp.SetCode(prop.Code())
	newProp.SetData(prop.Data())
	return newProp
}
