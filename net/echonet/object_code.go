// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

const (
	errorInvalidObjectCodes = "invalid object code : %s"
)

const (
	ObjectCodeMin     = 0x000000
	ObjectCodeMax     = 0xFFFFFF
	ObjectCodeSize    = 3
	ObjectCodeUnknown = ObjectCodeMin
)

// ObjectCode is a type for object code.
type ObjectCode = protocol.ObjectCode

// BytesToObjectCode converts the specified object code bytes to the object code.
func BytesToObjectCode(codes []byte) (ObjectCode, error) {
	if len(codes) != ObjectCodeSize {
		return 0, fmt.Errorf(errorInvalidObjectCodes, string(codes))
	}
	return ObjectCode(encoding.ByteToInteger(codes)), nil
}
