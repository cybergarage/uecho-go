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

// ClassOption is a function type to set options for class.
type ClassOption func(*class) error

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
}

// class is an instance for Echonet class.
type class struct {
	codes []byte
}

// WithClassBytes returns a ClassOption to set class codes with the specified byte slice.
func WithClassBytes(codes []byte) ClassOption {
	return func(cls *class) error {
		cls.codes[0] = codes[0]
		cls.codes[1] = codes[1]
		return nil
	}
}

// WithClassCode returns a ClassOption to set class code.
func WithClassCode(code byte) ClassOption {
	return func(cls *class) error {
		cls.codes[1] = byte(code & 0xFF)
		return nil
	}
}

// WithClassGroupCode returns a ClassOption to set class group code.
func WithClassGroupCode(code byte) ClassOption {
	return func(cls *class) error {
		cls.codes[0] = byte(code & 0xFF)
		return nil
	}
}

// NewClass returns a new class.
func NewClass(opts ...ClassOption) (Class, error) {
	cls := newClass()
	for _, opt := range opts {
		if err := opt(cls); err != nil {
			return nil, err
		}
	}
	return cls, nil
}

func newClass() *class {
	return &class{
		codes: make([]byte, 2),
	}
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
