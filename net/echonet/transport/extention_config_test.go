// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
)

func TestNewDefaultExtentionConfigConfig(t *testing.T) {
	NewDefaultExtentionConfig()
}

func TestExtentionConfigEquals(t *testing.T) {
	conf01 := NewDefaultExtentionConfig()
	conf02 := NewDefaultExtentionConfig()

	// Testing Set*()

	if !conf01.Equals(conf02) {
		t.Errorf("%v != %v", conf01, conf02)
	}

	conf01.SetAutoPortBindingEnabled(true)
	conf02.SetAutoPortBindingEnabled(false)
	if conf01.Equals(conf02) {
		t.Errorf("%v == %v", conf01, conf02)
	}

	// Testing SetConfig()

	conf03 := NewDefaultExtentionConfig()
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

func TestExtentionAutoBindingConfig(t *testing.T) {
	conf := newTestDefaultConfig()
	conf.SetAutoPortBindingEnabled(false)

	// Start on the default port

	mgr01 := NewMessageManager()
	mgr01.SetConfig(conf)
	err := mgr01.Start()
	defer mgr01.Stop()
	if err != nil {
		t.Error(err)
		return
	}
	if mgr01.GetPort() != UDPPort {
		t.Errorf("%d != %d", mgr01.GetPort(), UDPPort)
		return
	}

	// Disable auto binding option

	mgr02 := NewMessageManager()
	mgr02.SetConfig(conf)
	err = mgr02.Start()
	if err == nil {
		mgr02.Stop()
		t.Errorf("%v", conf)
		return
	}
	mgr02.Stop()

	// Enable auto binding option

	mgr02 = NewMessageManager()
	conf.SetAutoPortBindingEnabled(true)
	mgr02.SetConfig(conf)

	mgr02.SetConfig(conf)
	err = mgr02.Start()
	if err != nil {
		t.Error(err)
	}

	if mgr02.GetPort() == UDPPort {
		t.Errorf("%d == %d", mgr02.GetPort(), UDPPort)
		return
	}

	mgr02.Stop()
}
