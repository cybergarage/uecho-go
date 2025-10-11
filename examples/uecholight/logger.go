// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
)

func OutputError(err error) {
	fmt.Printf("ERROR : %s\n", err.Error())
}

func OutputMessage(format string, args ...any) {
	fmt.Printf(format, args...)
	fmt.Println()
}
