// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

/****************************************
 * Object Super Class
 ****************************************/

const (
	ObjectOperatingStatus  = 0x80
	ObjectManufacturerCode = 0x8A
	ObjectAnnoPropertyMap  = 0x9D
	ObjectSetPropertyMap   = 0x9E
	ObjectGetPropertyMap   = 0x9F
)

const (
	ObjectOperatingStatusOn      = 0x30
	ObjectOperatingStatusOff     = 0x31
	ObjectOperatingStatusSize    = 1
	ObjectManufacturerCodeSize   = 3
	ObjectPropertyMapMaxSize     = 16
	ObjectAnnoPropertyMapMaxSize = (ObjectPropertyMapMaxSize + 1)
	ObjectSetPropertyMapMaxSize  = (ObjectPropertyMapMaxSize + 1)
	ObjectGetPropertyMapMaxSize  = (ObjectPropertyMapMaxSize + 1)
)

// SetOperatingStatus sets a operating status to the Object.
func (obj *Object) SetOperatingStatus(stats bool) error {
	statsByte := byte(ObjectOperatingStatusOff)
	if stats {
		statsByte = ObjectOperatingStatusOn
	}
	return obj.SetPropertyByteData(ObjectOperatingStatus, statsByte)
}

// GetOperatingStatus return the operating status of the Object.
func (obj *Object) GetOperatingStatus() (bool, error) {
	statsByte, err := obj.GetPropertyByteData(ObjectOperatingStatus)
	if err != nil {
		return false, err
	}
	if statsByte == ObjectOperatingStatusOff {
		return false, nil
	}
	return true, nil
}
