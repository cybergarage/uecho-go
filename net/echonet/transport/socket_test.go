// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
)

func TestNewSocket(t *testing.T) {
	sock := NewUDPSocket()
	sock.IsBound()
}
