// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

// PropertyAttribute is a type for property attribute.
type PropertyAttribute uint

const (
	Prohibited = PropertyAttribute(0x00)
	Required   = PropertyAttribute(0x01)
	Optional   = PropertyAttribute(0x02)
)

// IsRequired returns true when the attribute is Required.
func (attr PropertyAttribute) IsRequired() bool {
	return attr == Required
}

// IsOptional returns true when the attribute is Optional.
func (attr PropertyAttribute) IsOptional() bool {
	return attr == Optional
}

// IsProhibited returns true when the attribute is not Prohibited.
func (attr PropertyAttribute) IsProhibited() bool {
	return attr == Prohibited
}
