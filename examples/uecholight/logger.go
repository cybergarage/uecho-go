// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
)

func OutputError(err error) {
	fmt.Println(fmt.Sprintf("ERROR : %s", err.Error()))
}

func OutputMessage(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}
