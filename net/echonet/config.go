// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/transport"
)

// TransportConfig represents a cofiguration for transport.
type TransportConfig = transport.Config

// Config represents a cofiguration for transport.
type Config struct {
	*TransportConfig
}

// NewDefaultConfig returns a default configuration.
func NewDefaultConfig() *Config {
	conf := &Config{
		TransportConfig: transport.NewDefaultConfig(),
	}
	return conf
}

// SetConfig sets all configuration flags.
func (conf *Config) SetConfig(newConfig *Config) {
	conf.TransportConfig.SetConfig(newConfig.TransportConfig)
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *Config) Equals(other *Config) bool {
	if !conf.TransportConfig.Equals(other.TransportConfig) {
		return false
	}
	return true
}
