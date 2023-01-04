// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

// StandardDatabase represents a standard database of Echonet.
type StandardDatabase struct {
	Manufactures map[ManufactureCode]*Manufacture
	Objects      map[ObjectCode]*Object
}

// NewStandardDatabase returns a standard database instance.
func NewStandardDatabase() *StandardDatabase {
	db := &StandardDatabase{
		Manufactures: map[ManufactureCode]*Manufacture{},
		Objects:      map[ObjectCode]*Object{},
	}
	db.initManufactures()
	db.initObjects()
	return db
}

func (db *StandardDatabase) addManufacture(man *Manufacture) {
	db.Manufactures[man.code] = man
}

// FindManufacture returns the registered manuracture by the specified manuracture code.
func (db *StandardDatabase) FindManufacture(code ManufactureCode) (*Manufacture, bool) {
	m, ok := db.Manufactures[code]
	return m, ok
}

func (db *StandardDatabase) addObject(obj *Object) {
	db.Objects[obj.Code()] = obj
}

// FindObjectByCode returns the registered object by the specified object code.
func (db *StandardDatabase) FindObjectByCode(code ObjectCode) (*Object, bool) {
	obj, ok := db.Objects[(code & 0xFFFF00)]
	return obj, ok
}

// FindObjectByCodes returns the registered object by the specified object code.
func (db *StandardDatabase) FindObjectByCodes(codes []byte) (*Object, bool) {
	if len(codes) != ObjectCodeSize {
		return nil, false
	}
	return db.FindObjectByCode(ObjectCode(encoding.ByteToInteger([]byte{codes[0], codes[1], 0x00})))
}

// SuperObject returns the super object.
func (db *StandardDatabase) SuperObject() (*Object, bool) {
	return db.FindObjectByCode(0x000000)
}

// NodeProfile returns the node profile object.
func (db *StandardDatabase) NodeProfile() (*Object, bool) {
	return db.FindObjectByCode(0x0EF000)
}
