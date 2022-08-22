// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

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
	codes []byte
}

// NewClass returns a new class.
func NewClass() *Class {
	cls := &Class{
		codes: make([]byte, 2),
	}
	return cls
}

// NewClassWithCodes returns a new class with the specified codes.
func NewClassWithCodes(codes []byte) *Class {
	cls := &Class{
		codes: codes,
	}
	return cls
}

// SetGroupCode sets a group code to the class.
func (cls *Class) SetGroupCode(code byte) {
	cls.codes[0] = code
}

// GroupCode returns the group code of the class.
func (cls *Class) GroupCode() byte {
	return cls.codes[0]
}

// SetCode sets a code to the class.
func (cls *Class) SetCode(code byte) {
	cls.codes[1] = code
}

// Code returns the group code of the class.
func (cls *Class) Code() byte {
	return cls.codes[1]
}

// Codes returns the all codes of the class.
func (cls *Class) Codes() []byte {
	return cls.codes
}

// Equals returns true whether the specified other class is same, otherwise false.
func (cls *Class) Equals(other *Class) bool {
	if len(cls.codes) != len(other.codes) {
		return false
	}
	if len(cls.codes) != 2 {
		return false
	}
	if cls.GroupCode() != other.GroupCode() {
		return false
	}
	if cls.Code() != other.Code() {
		return false
	}
	return true
}
