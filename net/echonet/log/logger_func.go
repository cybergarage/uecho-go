// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

// SetDebugEnabled sets the logging level of the current shared logger to trace.
func SetDebugEnabled(flag bool) {
	logger := GetSharedLogger()
	if logger == nil {
		return
	}
	logger.SetLevel(LoggerLevelTrace)
}

// SetStdoutEnbled sets the shared logger to stdout
func SetStdoutEnbled(flag bool) {
	if flag {
		logger := GetSharedLogger()
		if logger == nil {
			return
		}
		SetSharedLogger(NewStdoutLogger(logger.GetLevel()))
	}
}

// SetStdoutDebugEnbled sets a trace stdout logger for debug.
func SetStdoutDebugEnbled(flag bool) {
	if flag {
		SetSharedLogger(NewStdoutLogger(LoggerLevelTrace))
	}
}
