// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

// ManufactureCode represent a manufacture code.
type ManufactureCode int

// Manufacture represents a manufacture object.
type Manufacture struct {
	Code ManufactureCode
	Name string
}

// NewManufacture returns a manufacture instance.
func NewManufacture(code ManufactureCode, name string) *Manufacture {
	return &Manufacture{
		Code: code,
		Name: name,
	}
}
