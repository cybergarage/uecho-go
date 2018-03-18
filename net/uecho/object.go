// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

const ()

// Object is an instance for Echonet object.
type Object struct {
	Code             [3]byte
	Properties       []*Property
	annoPropMapSize  int
	annoPropMapBytes []byte
	setPropMapSize   int
	setPropMapBytes  []byte
	getPropMapSize   int
	getPropMapBytes  []byte
}

// NewObject returns a new object.
func NewObject() *Object {
	obj := &Object{}
	return obj
}
