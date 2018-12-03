// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
)

func TestMulticastSocketBindWithInterface(t *testing.T) {
	sock := NewMulticastSocket()

	ifs, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
	}

	err = sock.Bind(ifs[0])
	if err != nil {
		t.Error(err)
	}

	err = sock.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestMulticastSocketBindWithNoInterface(t *testing.T) {
	sock := NewMulticastSocket()

	err := sock.Bind(nil)
	if err != nil {
		t.Error(err)
	}

	err = sock.Close()
	if err != nil {
		t.Error(err)
	}
}
