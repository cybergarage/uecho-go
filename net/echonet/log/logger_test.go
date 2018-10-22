// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log

import (
	"errors"
	"testing"
)

const (
	testLogMessage         = "hello"
	nullOutputErrorMessage = "Shared Logger is null, but message is output"
	outputErrorMessage     = "Message can't be output"
)

func TestNullLogger(t *testing.T) {
	SetSharedLogger(nil)

	nOutput := Trace(testLogMessage)
	if 0 < nOutput {
		t.Error(errors.New(nullOutputErrorMessage))
	}

	nOutput = Info(testLogMessage)
	if 0 < nOutput {
		t.Error(errors.New(nullOutputErrorMessage))
	}

	nOutput = Error(testLogMessage)
	if 0 < nOutput {
		t.Error(errors.New(nullOutputErrorMessage))
	}

	nOutput = Warn(testLogMessage)
	if 0 < nOutput {
		t.Error(errors.New(nullOutputErrorMessage))
	}

	nOutput = Fatal(testLogMessage)
	if 0 < nOutput {
		t.Error(errors.New(nullOutputErrorMessage))
	}
}

func TestStdoutLogger(t *testing.T) {
	SetSharedLogger(NewStdoutLogger(LevelTrace))
	defer SetSharedLogger(nil)

	nOutput := Trace(testLogMessage)
	if nOutput <= 0 {
		t.Error(errors.New(outputErrorMessage))
	}

	nOutput = Info(testLogMessage)
	if nOutput <= 0 {
		t.Error(errors.New(outputErrorMessage))
	}

	nOutput = Error(testLogMessage)
	if nOutput <= 0 {
		t.Error(errors.New(outputErrorMessage))
	}

	nOutput = Warn(testLogMessage)
	if nOutput <= 0 {
		t.Error(errors.New(outputErrorMessage))
	}

	nOutput = Fatal(testLogMessage)
	if nOutput <= 0 {
		t.Error(errors.New(outputErrorMessage))
	}
}
