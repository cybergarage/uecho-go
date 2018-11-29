// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"reflect"
	"time"
)

// UnicastConfig represents a cofiguration for unicast server.
type UnicastConfig struct {
	TCPEnabled        bool
	ConnectionTimeout time.Duration
	RequestTimeout    time.Duration
}

// NewDefaultUnicastConfig returns a default configuration.
func NewDefaultUnicastConfig() *UnicastConfig {
	conf := &UnicastConfig{
		TCPEnabled:        false,
		ConnectionTimeout: DefaultConnectimeTimeOut,
		RequestTimeout:    DefaultRequestTimeout,
	}
	return conf
}

// SetConfig sets all flags.
func (conf *UnicastConfig) SetConfig(newConfig *UnicastConfig) {
	conf.TCPEnabled = newConfig.TCPEnabled
	conf.ConnectionTimeout = newConfig.ConnectionTimeout
}

// SetTCPEnabled sets a flag for TCP functions.
func (conf *UnicastConfig) SetTCPEnabled(flag bool) {
	conf.TCPEnabled = flag
}

// IsTCPEnabled returns true whether the TCP function is enabled, otherwise false.
func (conf *UnicastConfig) IsTCPEnabled() bool {
	return conf.TCPEnabled
}

// SetConnectionTimeout sets a connection timeout setting.
func (conf *UnicastConfig) SetConnectionTimeout(timeout time.Duration) {
	conf.ConnectionTimeout = timeout
}

// GetConnectionTimeout returns the current connection timeout setting.
func (conf *UnicastConfig) GetConnectionTimeout() time.Duration {
	return conf.ConnectionTimeout
}

// SetRequestTimeout sets a request timeout.
func (conf *Config) SetRequestTimeout(d time.Duration) {
	conf.RequestTimeout = d
}

// GetRequestTimeout return the request timeout.
func (conf *Config) GetRequestTimeout() time.Duration {
	return conf.RequestTimeout
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *UnicastConfig) Equals(otherConf *UnicastConfig) bool {
	return reflect.DeepEqual(conf, otherConf)
}
