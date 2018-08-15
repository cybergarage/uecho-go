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
	Code []byte
}

// NewClass returns a new class.
func NewClass() *Class {
	cls := &Class{
		Code: make([]byte, 2),
	}
	return cls
}

// NewClassWithCodes returns a new class with the specified codes.
func NewClassWithCodes(codes []byte) *Class {
	cls := &Class{
		Code: codes,
	}
	return cls
}

// SetGroupCode sets a group code to the class.
func (cls *Class) SetGroupCode(code byte) {
	cls.Code[0] = code
}

// GetGroupCode returns the group code of the class.
func (cls *Class) GetGroupCode() byte {
	return cls.Code[0]
}

// SetCode sets a code to the class.
func (cls *Class) SetCode(code byte) {
	cls.Code[1] = code
}

// GetCode returns the group code of the class.
func (cls *Class) GetCode() byte {
	return cls.Code[1]
}

// GetCodes returns the all codes of the class.
func (cls *Class) GetCodes() []byte {
	return cls.Code
}

// Equals returns true whether the specified other class is same, otherwise false.
func (cls *Class) Equals(other *Class) bool {
	if len(cls.Code) != len(other.Code) {
		return false
	}
	if len(cls.Code) != 2 {
		return false
	}
	if cls.GetGroupCode() != other.GetGroupCode() {
		return false
	}
	if cls.GetCode() != other.GetCode() {
		return false
	}
	return true
}
