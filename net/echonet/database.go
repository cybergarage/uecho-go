// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
)

// StandardDatabase represents a standard database of Echonet. The database has standard manufacture and device and profile objects from Echonet Consortium. The object database is based on Machine Readable Appendix (MRA) which  is a data file that describes the contents of “APPENDIX Detailed Requirements for ECHONET Device objects”.
type StandardDatabase interface {
	// Manufactures returns the all registered manufactures.
	Manufactures() []*Manufacture

	// Objects returns the all registered objects.
	Objects() []Object

	// LookupManufacture returns the registered manuracture by the specified manuracture code.
	LookupManufacture(code ManufactureCode) (*Manufacture, bool)

	// LookupObjectByCode returns the registered object by the specified object code.
	LookupObjectByCode(code ObjectCode) (Object, bool)

	// LookupObjectByCodes returns the registered object by the specified object code.
	LookupObjectByCodes(codes []byte) (Object, bool)

	// SuperObject returns the super object.
	SuperObject() Object

	// NodeProfile returns the node profile object.
	NodeProfile() Object
}

// stdDatabase represents a standard database of Echonet.
type stdDatabase struct {
	manufactures map[ManufactureCode]*Manufacture
	objects      map[ObjectCode]Object
}

// newStandardDatabase returns a standard database instance.
func newStandardDatabase() StandardDatabase {
	db := &stdDatabase{
		manufactures: map[ManufactureCode]*Manufacture{},
		objects:      map[ObjectCode]Object{},
	}
	db.initManufactures()
	db.initObjects()
	return db
}

func (db *stdDatabase) addManufacture(man *Manufacture) {
	db.manufactures[man.code] = man
}

// Manufactures returns the all registered manufactures.
func (db *stdDatabase) Manufactures() []*Manufacture {
	mans := []*Manufacture{}
	for _, man := range db.manufactures {
		mans = append(mans, man)
	}
	return mans
}

// Objects returns the all registered objects.
func (db *stdDatabase) Objects() []Object {
	objs := []Object{}
	for _, obj := range db.objects {
		objs = append(objs, obj)
	}
	return objs
}

// LookupManufacture returns the registered manuracture by the specified manuracture code.
func (db *stdDatabase) LookupManufacture(code ManufactureCode) (*Manufacture, bool) {
	m, ok := db.manufactures[code]
	return m, ok
}

func (db *stdDatabase) addObject(obj Object) {
	db.objects[obj.Code()] = obj
}

// LookupObjectByCode returns the registered object by the specified object code.
func (db *stdDatabase) LookupObjectByCode(code ObjectCode) (Object, bool) {
	obj, ok := db.objects[(code & 0xFFFF00)]
	return obj, ok
}

// LookupObjectByCodes returns the registered object by the specified object code.
func (db *stdDatabase) LookupObjectByCodes(codes []byte) (Object, bool) {
	if len(codes) != ObjectCodeSize {
		return nil, false
	}
	return db.LookupObjectByCode(ObjectCode(encoding.ByteToInteger([]byte{codes[0], codes[1], 0x00})))
}

// SuperObject returns the super object.
func (db *stdDatabase) SuperObject() Object {
	obj, _ := db.LookupObjectByCode(0x000000)
	return obj
}

// NodeProfile returns the node profile object.
func (db *stdDatabase) NodeProfile() Object {
	obj, _ := db.LookupObjectByCode(0x0EF000)
	return obj
}
