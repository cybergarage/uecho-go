// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/cybergarage/go-logger/log"
)

func outputf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func errorf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func enableStdoutVerbose(flag bool) {
	if flag {
		log.SetSharedLogger(log.NewStdoutLogger(log.LevelInfo))
	} else {
		log.SetSharedLogger(nil)
	}
}
