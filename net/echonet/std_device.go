// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

// NewStandardDeviceWithCodes returns a new device of the specified object codes.
func NewStandardDeviceWithCodes(codes []byte) (*Device, error) {
	objCode, err := BytesToObjectCode(codes)
	if err != nil {
		return nil, err
	}
	obj := NewDevice()
	obj.SetCode(objCode)
	return obj, nil
}

// NewStandardDeviceWithCode returns a new device of the specified object code.
func NewStandardDeviceWithCode(code ObjectCode) *Device {
	obj := NewDevice()
	obj.SetCode(code)
	return obj
}
