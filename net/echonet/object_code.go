// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	ObjectCodeMin     = protocol.ObjectCodeMin
	ObjectCodeMax     = protocol.ObjectCodeMax
	ObjectCodeSize    = protocol.ObjectCodeSize
	ObjectCodeUnknown = ObjectCodeMin
)

// ObjectCode is a type for object code.
type ObjectCode = protocol.ObjectCode

// NewObjectCodeFromBytes converts the specified object code bytes to the object code.
func NewObjectCodeFromBytes(codes []byte) (ObjectCode, error) {
	return protocol.NewObjectCodeFromBytes(codes)
}
