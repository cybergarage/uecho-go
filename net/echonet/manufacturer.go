// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

// ManufactureCode represent a manufacture code.
type ManufactureCode int

// Manufacture represents a manufacture interface.
type Manufacture interface {
	// Name returns the manufacture name.
	Name() string

	// Code returns the manufacture code.
	Code() ManufactureCode
}

// manufacture represents a manufacture object.
type manufacture struct {
	code ManufactureCode
	name string
}

// newManufacture returns a manufacture instance.
func newManufacture(code ManufactureCode, name string) Manufacture {
	return &manufacture{
		code: code,
		name: name,
	}
}

// Name returns the manufacture name.
func (man *manufacture) Name() string {
	return man.name
}

// Code returns the manufacture code.
func (man *manufacture) Code() ManufactureCode {
	return man.code
}
