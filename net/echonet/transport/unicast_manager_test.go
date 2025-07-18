// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
)

func testUnicastManagerBinding(t *testing.T, conf *Config) {
	t.Helper()

	mgr := NewUnicastManager()
	mgr.SetConfig(conf)

	err := mgr.Start()
	if err != nil {
		t.Error(err)
	}

	err = mgr.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestUnicastManagerWithDefaultConfig(t *testing.T) {
	// log.EnableStdoutDebug(true)
	conf := newTestDefaultConfig()
	testUnicastManagerBinding(t, conf)
}

func TestUnicastManagerWithOnlyUDPConfig(t *testing.T) {
	// log.EnableStdoutDebug(true)
	conf := newTestDefaultConfig()
	conf.SetTCPEnabled(false)
	testUnicastManagerBinding(t, conf)
}

func TestUnicastManagerWithTCPConfig(t *testing.T) {
	// log.EnableStdoutDebug(true)
	conf := newTestDefaultConfig()
	conf.SetTCPEnabled(true)
	testUnicastManagerBinding(t, conf)
}
