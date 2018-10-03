// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
	"time"
)

func TestNewUnicastDefaultConfig(t *testing.T) {
	NewDefaultUnicastConfig()
}

func TestUnicastEquals(t *testing.T) {
	conf01 := NewDefaultUnicastConfig()
	conf02 := NewDefaultUnicastConfig()

	// Testing Set*()

	if !conf01.Equals(conf02) {
		t.Errorf("%v != %v", conf01, conf02)
	}

	conf01.SetTCPEnabled(true)
	conf02.SetTCPEnabled(false)
	if conf01.Equals(conf02) {
		t.Errorf("%v == %v", conf01, conf02)
	}

	conf01.SetConnectionTimeout(time.Microsecond)
	conf02.SetConnectionTimeout(time.Second)
	if conf01.Equals(conf02) {
		t.Errorf("%v == %v", conf01, conf02)
	}

	// Testing SetConfig()

	conf03 := NewDefaultUnicastConfig()
	conf03.SetConfig(conf01)
	if !conf01.Equals(conf03) {
		t.Errorf("%v != %v", conf01, conf03)
	}
	if conf02.Equals(conf03) {
		t.Errorf("%v == %v", conf01, conf02)
	}

	conf03.SetConfig(conf02)
	if !conf02.Equals(conf03) {
		t.Errorf("%v != %v", conf01, conf03)
	}
}
