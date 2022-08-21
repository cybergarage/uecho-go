// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

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
	return db
}

func (db *StandardDatabase) addManufacture(man *Manufacture) {
	db.Manufactures[man.Code] = man
}

// GetManufacture returns the registered manuracture by the specified manuracture code.
func (db *StandardDatabase) GetManufacture(code ManufactureCode) (*Manufacture, bool) {
	m, ok := db.Manufactures[code]
	return m, ok
}

func (db *StandardDatabase) addObject(obj *Object) {
	db.Objects[obj.GetCode()] = obj
}

// GetObject returns the registered object by the specified object code.
func (db *StandardDatabase) GetObject(code ObjectCode) (*Object, bool) {
	obj, ok := db.Objects[code]
	return obj, ok
}
