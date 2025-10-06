// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

var sharedStandardDatabase *StandardDatabase

func init() {
	sharedStandardDatabase = NewStandardDatabase()
}

// SharedStandardDatabase returns the shared standard database.
func SharedStandardDatabase() *StandardDatabase {
	return sharedStandardDatabase
}
