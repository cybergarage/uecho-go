// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
)

func TestNewMulticastServerManager(t *testing.T) {
	mgr := NewMulticastServerManager()

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
