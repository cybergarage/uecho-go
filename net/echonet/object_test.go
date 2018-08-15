// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"testing"

	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

const (
	errorMandatoryPropertyNotFound = "Mandatory Property (%0X) Not Found"
	errorInvalidGroupClassCode     = "Invalid Group Class Code (%0X)"
)

func TestNewObject(t *testing.T) {
	NewObject()
}

func TestObjectCodes(t *testing.T) {
	objCodes := []uint{
		NodeProfileObject,
		NodeProfileObjectReadOnly,
	}

	objCodeBytes := make([]byte, 3)
	for n, objCode := range objCodes {
		encoding.IntegerToByte(objCode, objCodeBytes)
		if objCode != encoding.ByteToInteger(objCodeBytes) {
			t.Errorf("[%d] : %X != %X", n, encoding.ByteToInteger(objCodeBytes), objCode)
		}
	}
}
