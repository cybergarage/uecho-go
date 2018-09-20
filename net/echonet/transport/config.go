// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// Config represents a cofiguration for transport.
type Config struct {
	*UnicastConfig
}

// NewDefaultConfig returns a default configuration.
func NewDefaultConfig() *Config {
	conf := &Config{
		UnicastConfig: NewDefaultUnicastConfig(),
	}
	return conf
}
