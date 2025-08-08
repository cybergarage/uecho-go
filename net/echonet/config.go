// Copyright 2018 The uecho-go Authors. All rights reserved.
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

	SelfMessageEnabled bool
}

// NewDefaultConfig returns a default configuration.
func NewDefaultConfig() *Config {
	conf := &Config{
		SelfMessageEnabled: true,
		TransportConfig:    transport.NewDefaultConfig(),
	}
	return conf
}

// SetSelfMessageEnabled sets a flag for self messages.
func (conf *Config) SetSelfMessageEnabled(flag bool) {
	conf.SelfMessageEnabled = flag
}

// IsSelfMessageEnabled returns true whether the self messages are enabled, otherwise false.
func (conf *Config) IsSelfMessageEnabled() bool {
	return conf.SelfMessageEnabled
}

// SetConfig sets all configuration flags.
func (conf *Config) SetConfig(newConfig *Config) {
	conf.TransportConfig.SetConfig(newConfig.TransportConfig)
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *Config) Equals(other *Config) bool {
	return conf.TransportConfig.Equals(other.TransportConfig)
}
