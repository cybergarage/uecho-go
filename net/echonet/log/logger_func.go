// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

// SetStdoutDebugEnbled sets a trace stdout logger for debug.
func SetStdoutDebugEnbled(flag bool) {
	if flag {
		SetSharedLogger(NewStdoutLogger(LoggerLevelTrace))
	} else {
		SetSharedLogger(nil)
	}
}
