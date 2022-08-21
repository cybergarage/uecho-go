// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
)

// NewStandardDeviceWithCodes returns a new device of the specified object codes.
func NewStandardDeviceWithCodes(codes []byte) (*Device, error) {
	if len(codes) != ObjectCodeSize {
		return nil, fmt.Errorf(errorInvalidObjectCodes, string(codes))
	}
	obj := NewDevice()
	obj.SetCodes(codes)
	addStandardProperties(obj.Object)
	return obj, nil
}
