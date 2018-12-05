// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
)

const (
	testUnicastTCPSocketPort = 32001
)

func TestUnicastTCPSocketOpenClose(t *testing.T) {
	sock := NewUnicastTCPSocket()

	ifs, err := GetAvailableInterfaces()
	if err != nil {
		t.Error(err)
		return
	}

	err = sock.Bind(ifs[0], testUnicastTCPSocketPort)
	if err != nil {
		t.Error(err)
		return
	}

	err = sock.Close()
	if err != nil {
		t.Error(err)
	}
}
