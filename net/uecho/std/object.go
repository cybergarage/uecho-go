// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package std

const (
	ObjectCodeMin     = 0x000000
	ObjectCodeMax     = 0xFFFFFF
	ObjectCodeUnknown = ObjectCodeMin

	NodeProfileObject         = 0x0EF001
	NodeProfileObjectReadOnly = 0x0EF002
)

const (
	ObjectManufacturerCode = 0x8A
	ObjectAnnoPropertyMap  = 0x9D
	ObjectSetPropertyMap   = 0x9E
	ObjectGetPropertyMap   = 0x9F
)

const (
	ObjectManufacturerCodeLen   = 3
	ObjectPropertyMapMaxLen     = 16
	ObjectAnnoPropertyMapMaxLen = (ObjectPropertyMapMaxLen + 1)
	ObjectSetPropertyMapMaxLen  = (ObjectPropertyMapMaxLen + 1)
	ObjectGetPropertyMapMaxLen  = (ObjectPropertyMapMaxLen + 1)
)

const (
	ManufacturerCodeMin    = 0xFFFFF0
	ManufacturerCodeMax    = 0xFFFFFF
	ManufactureCodeDefault = ManufacturerCodeMin
)
