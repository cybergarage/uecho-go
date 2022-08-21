// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

// StandardDatabase represents a standard database of Echonet.
type StandardDatabase struct {
}

// NewStandardDatabase returns a standard database instance.
func NewStandardDatabase() *StandardDatabase {
	db := &StandardDatabase{}
	return db
}
