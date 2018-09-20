// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
)

import (
	"github.com/cybergarage/uecho-go/net/echonet/log"
)

func testUnicastManagerBinding(t *testing.T, conf *UnicastConfig) {
	mgr := NewUnicastManager()
	mgr.SetConfig(conf)
	err := mgr.Start()
	if err != nil {
		t.Error(err)
		return
	}

	err = mgr.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestUnicastManagerWithDefaultConfig(t *testing.T) {
	conf := NewDefaultUnicastConfig()
	testUnicastManagerBinding(t, conf)
}

func TestUnicastManagerWithOnlyUDPConfig(t *testing.T) {
	conf := NewDefaultUnicastConfig()
	conf.SetTCPEnabled(false)
	conf.SetUDPEnabled(true)
	testUnicastManagerBinding(t, conf)
}

func TestUnicastManagerWithTCPConfig(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	conf := NewDefaultUnicastConfig()
	conf.SetTCPEnabled(true)
	conf.SetUDPEnabled(true)
	testUnicastManagerBinding(t, conf)
}
