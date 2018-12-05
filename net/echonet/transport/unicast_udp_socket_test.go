// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
)

const (
	testUnicastUDPSocketPort = testUnicastTCPSocketPort + 1
)

func TestUnicastUDPSocketOpenClose(t *testing.T) {
	sock := NewUnicastUDPSocket()

	ifs, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
		return
	}

	err = sock.Bind(ifs[0], testUnicastUDPSocketPort)
	if err != nil {
		t.Error(err)
		return
	}

	err = sock.Close()
	if err != nil {
		t.Error(err)
	}
}
