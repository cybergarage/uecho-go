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

// Class is an interface for Echonet class.
type Class interface {
	// GroupCode returns the group code of the class.
	GroupCode() byte
	// Code returns the code of the class.
	Code() byte
	// Equals returns true whether the specified other class is same, otherwise false.
	Equals(other Class) bool
	// Bytes returns the all codes of the class.
	Bytes() []byte
	// ClasstMutator is an interface to mutate the class.
	ClasstMutator
}

// ClasstMutator is an interface to mutate the class.
type ClasstMutator interface {
	// SetGroupCode sets a group code to the class.
	SetGroupCode(code byte)
	// SetCode sets a code to the class.
	SetCode(code byte)
}

// class is an instance for Echonet class.
type class struct {
	codes []byte
}

// NewClass returns a new class.
func NewClass() Class {
	cls := &class{
		codes: make([]byte, 2),
	}
	return cls
}

// NewClassWithBytes returns a new class with the specified byte codes.
func NewClassWithBytes(codes []byte) Class {
	cls := &class{
		codes: codes,
	}
	return cls
}

// SetGroupCode sets a group code to the class.
func (cls *class) SetGroupCode(code byte) {
	cls.codes[0] = code
}

// GroupCode returns the group code of the class.
func (cls *class) GroupCode() byte {
	return cls.codes[0]
}

// SetCode sets a code to the class.
func (cls *class) SetCode(code byte) {
	cls.codes[1] = code
}

// Code returns the group code of the class.
func (cls *class) Code() byte {
	return cls.codes[1]
}

// Bytes returns the all codes of the class.
func (cls *class) Bytes() []byte {
	return cls.codes
}

// Equals returns true whether the specified other class is same, otherwise false.
func (cls *class) Equals(other Class) bool {
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
