// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func TestPostController(t *testing.T) {
	ctrl := NewPostController()
	err := ctrl.Start()
	if err != nil {
		t.Error(err)
		return
	}

	err = ctrl.SearchAllObjects()
	if err != nil {
		t.Error(err)
	}

	err = ctrl.Stop()
	if err != nil {
		t.Error(err)
	}

}
