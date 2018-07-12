// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

const (
	ClassCodeMin = 0x0000
	ClassCodeMax = 0xFFFF
)

const (
	ClassGroupDeviceMin = 0x00
	ClassGroupDeviceMax = 0x06
	ClassGroupProfile   = 0x0E
)

// Class is an instance for Echonet class.
type Class struct {
	Code [2]byte
}

// NewClass returns a new class.
func NewClass() *Class {
	cls := &Class{}
	return cls
}
