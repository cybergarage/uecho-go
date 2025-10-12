// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/transport"
)

type transportConfig = transport.Config

// Config is an interface for Echonet configuration.
type Config interface {
	configInternal
}

type configInternal interface {
	TransportConfig() *transportConfig
	SelfMessageEnabled() bool
	TCPEnabled() bool
	RequestTimeout() time.Duration
}

// ConfigOption is a function that configures a configuration.
type ConfigOption func(*config)

// config represents a cofiguration for transport.
type config struct {
	*transportConfig
	selfMsgEnabled bool
}

// WithSelfMessageEnabled sets a flag for self messages.
func WithConfigTCPEnabled(flag bool) ConfigOption {
	return func(conf *config) {
		conf.SetTCPEnabled(flag)
	}
}

// NewDefaultConfig returns a new default configuration.
func NewDefaultConfig(opts ...ConfigOption) Config {
	return newDefaultConfig(opts...)
}

func newDefaultConfig(opts ...ConfigOption) *config {
	conf := &config{
		selfMsgEnabled:  true,
		transportConfig: transport.NewDefaultConfig(),
	}
	for _, opt := range opts {
		opt(conf)
	}
	return conf
}

// TransportConfig returns the transport configuration.
func (conf *config) TransportConfig() *transportConfig {
	return conf.transportConfig
}

// SetTCPEnabled sets a flag for TCP functions.
func (conf *config) SetTCPEnabled(flag bool) Config {
	conf.transportConfig.SetTCPEnabled(flag)
	return conf
}

// SetSelfMessageEnabled sets a flag for self messages.
func (conf *config) SetSelfMessageEnabled(flag bool) Config {
	conf.selfMsgEnabled = flag
	return conf
}

// SelfMessageEnabled returns true whether the self messages are enabled, otherwise false.
func (conf *config) SelfMessageEnabled() bool {
	return conf.selfMsgEnabled
}
