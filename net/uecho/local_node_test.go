// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"testing"
)

func TestNewLocalNode(t *testing.T) {
	node := NewLocalNode()

	_, err := node.GetNodeProfile()
	if err != nil {
		t.Error(err)
	}
}