// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
)

import (
	"github.com/cybergarage/uecho-go/net/echonet/log"
)

func testUnicastManagerBinding(t *testing.T, conf *Config) {
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
	conf := NewDefaultConfig()
	testUnicastManagerBinding(t, conf)
}

func TestUnicastManagerWithOnlyUDPConfig(t *testing.T) {
	conf := NewDefaultConfig()
	conf.SetTCPEnabled(false)
	testUnicastManagerBinding(t, conf)
}

func TestUnicastManagerWithTCPConfig(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	conf := NewDefaultConfig()
	conf.SetTCPEnabled(true)
	testUnicastManagerBinding(t, conf)
}
