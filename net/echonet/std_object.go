// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

// NewStandardObjectWithCodes returns a new object of the specified object codes.
func NewStandardObjectWithCodes(codes []byte) (interface{}, error) {
	objCode, err := BytesToObjectCode(codes)
	if err != nil {
		return nil, err
	}
	if isProfileObjectCode(codes[0]) {
		obj := NewProfile()
		obj.SetCode(objCode)
		return obj, nil
	}
	return NewStandardDeviceWithCodes(codes)
}
