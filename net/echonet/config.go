// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/transport"
)

type transportConfig = transport.Config

// config represents a cofiguration for transport.
type Config struct {
	*transportConfig
	selfMsgEnabled bool
}

// NewDefaultConfig returns a default configuration.
func NewDefaultConfig() *Config {
	conf := &Config{
		selfMsgEnabled:  true,
		transportConfig: transport.NewDefaultConfig(),
	}
	return conf
}

// SetTCPEnabled sets a flag for TCP functions.
func (conf *Config) SetTCPEnabled(flag bool) *Config {
	conf.transportConfig.SetTCPEnabled(flag)
	return conf
}

// setSelfMessageEnabled sets a flag for self messages.
func (conf *Config) setSelfMessageEnabled(flag bool) *Config {
	conf.selfMsgEnabled = flag
	return conf
}

// selfMessageEnabled returns true whether the self messages are enabled, otherwise false.
func (conf *Config) selfMessageEnabled() bool {
	return conf.selfMsgEnabled
}
