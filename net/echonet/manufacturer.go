// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

// ManufactureCode represent a manufacture code.
type ManufactureCode int

// Manufacture represents a manufacture object.
type Manufacture struct {
	code ManufactureCode
	name string
}

// NewManufacture returns a manufacture instance.
func NewManufacture(code ManufactureCode, name string) *Manufacture {
	return &Manufacture{
		code: code,
		name: name,
	}
}

// Name returns the manufacture name.
func (man *Manufacture) Name() string {
	return man.name
}

// Code returns the manufacture code.
func (man *Manufacture) Code() ManufactureCode {
	return man.code
}
