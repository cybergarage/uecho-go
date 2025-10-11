// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

const (
	ObjectCodeMin     = 0x000000
	ObjectCodeMax     = 0xFFFFFF
	ObjectCodeSize    = 3
	ObjectCodeUnknown = ObjectCodeMin
)

// ObjectCode is a type for object code.
type ObjectCode uint

// NewObjectCodeFromBytes returns a new ObjectCode instance from the specified byte array.
func NewObjectCodeFromBytes(codes []byte) (ObjectCode, error) {
	if len(codes) != ObjectCodeSize {
		return 0, fmt.Errorf(errorInvalidObjectCodes, ErrInvalid, string(codes))
	}
	return ObjectCode(encoding.ByteToInteger(codes)), nil
}

// Bytes returns a byte array of the object code.
func (code ObjectCode) Bytes() []byte {
	codes := make([]byte, 3)
	encoding.IntegerToByte(uint(code), codes)
	return codes
}
