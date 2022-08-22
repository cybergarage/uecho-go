// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

func addStandardProperties(obj *Object) {
	db := GetStandardDatabase()
	stdObj, ok := db.FindObjectByCode(obj.Code())
	if !ok {
		return
	}
	obj.SetClassName(stdObj.ClassName())
	for _, prop := range stdObj.Properties() {
		obj.AddProperty(prop.Copy())
	}
}

// NewStandardObjectWithCodes returns a new object of the specified object codes.
func NewStandardObjectWithCodes(codes []byte) (interface{}, error) {
	objCode, err := BytesToObjectCode(codes)
	if err != nil {
		return nil, err
	}
	if isProfileObjectCode(codes[0]) {
		obj := NewProfile()
		obj.SetCode(objCode)
		addStandardProperties(obj.Object)
		return obj, nil
	}
	return NewStandardDeviceWithCodes(codes)
}
