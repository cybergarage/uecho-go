// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

const ()

// Property is an instance for Echonet property.
type Property struct {
	Code byte
	Attr int
	Data []byte
	Size int
}

// NewProperty returns a new property.
func NewProperty() *Property {
	prop := &Property{}
	return prop
}
