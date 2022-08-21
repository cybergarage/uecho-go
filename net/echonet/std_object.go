// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"fmt"
)

func addStandardProperties(obj *Object) {
	db := GetStandardDatabase()
	stdObj, ok := db.GetObjectByCode(obj.GetCode())
	if !ok {
		return
	}
	for _, prop := range stdObj.properties {
		obj.AddProperty(prop.Copy())
	}
}

// NewStandardObjectWithCodes returns a new object of the specified object codes.
func NewStandardObjectWithCodes(codes []byte) (interface{}, error) {
	if len(codes) != ObjectCodeSize {
		return nil, fmt.Errorf(errorInvalidObjectCodes, string(codes))
	}
	if isProfileObjectCode(codes[0]) {
		obj := NewProfile()
		obj.SetCodes(codes)
		addStandardProperties(obj.Object)
		return obj, nil
	}
	return NewStandardDeviceWithCodes(codes)
}
