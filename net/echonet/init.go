// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

var sharedStandardDatabase *StandardDatabase

func init() {
	sharedStandardDatabase = NewStandardDatabase()
}

func GetStandardDatabase() *StandardDatabase {
	return sharedStandardDatabase
}
