// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"testing"
)

func newTestDefaultConfig() *Config {
	conf := NewDefaultConfig()
	conf.SetAutoPortBindingEnabled(true)
	return conf
}
func TestNewDefaultConfig(t *testing.T) {
	NewDefaultConfig()
}
