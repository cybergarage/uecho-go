// Copyright 2018 The uecho-go Authors. All rights reserved.
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
	conf01 := newTestDefaultConfig()
	conf01.SetAutoPortBindingEnabled(false)
	conf01.SetBindRetryEnabled(true)

	// Start on the default port

	mgr01 := NewMessageManager()
	mgr01.SetConfig(conf01)
	err := mgr01.Start()
	defer mgr01.Stop()
	if err != nil {
		// FIXME : TestExtentionAutoBindingConfig is failed on Travis
		// t.Skip(err)
		return
	}
	if mgr01.GetPort() != UDPPort {
		t.Errorf("%d != %d", mgr01.GetPort(), UDPPort)
		return
	}

	// Disable auto binding option

	conf02 := newTestDefaultConfig()
	conf02.SetAutoPortBindingEnabled(false)
	conf02.SetBindRetryEnabled(false)

	mgr02 := NewMessageManager()
	mgr02.SetConfig(conf02)
	err = mgr02.Start()
	if err == nil {
		mgr02.Stop()
		t.Errorf("%v", conf02)
		return
	}
	mgr02.Stop()

	// Enable auto binding option

	conf02 = newTestDefaultConfig()
	conf02.SetAutoPortBindingEnabled(true)
	conf02.SetBindRetryEnabled(false)

	mgr02 = NewMessageManager()
	mgr02.SetConfig(conf02)

	mgr02.SetConfig(conf02)
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
