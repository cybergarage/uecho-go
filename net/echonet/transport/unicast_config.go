// Copyright 2018 The uecho-go Authors. All rights reserved.
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
	BindRetryCount    uint
	BindRetryWaitTime time.Duration
}

// NewDefaultUnicastConfig returns a default configuration.
func NewDefaultUnicastConfig() *UnicastConfig {
	conf := &UnicastConfig{
		TCPEnabled:        false,
		ConnectionTimeout: DefaultConnectimeTimeOut,
		RequestTimeout:    DefaultRequestTimeout,
		BindRetryCount:    DefaultBindRetryCount,
		BindRetryWaitTime: DefaultBindRetryWaitTime,
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

// SetBindRetryEnabled sets a flag for the bind retry function.
func (conf *UnicastConfig) SetBindRetryEnabled(flag bool) {
	if !flag {
		conf.BindRetryCount = 0
		return
	}
	conf.BindRetryCount = DefaultBindRetryCount
}

// IsBindRetryEnabled returns true whether the bind retry function is enabled, otherwise false.
func (conf *UnicastConfig) IsBindRetryEnabled() bool {
	if conf.BindRetryCount == 0 {
		return false
	}
	return true
}

// SetBindRetryCount sets a retry count when the binding is failed.
func (conf *Config) SetBindRetryCount(n uint) {
	conf.BindRetryCount = n
}

// GetBindRetryCount returns the retry count when the binding is failed.
func (conf *Config) GetBindRetryCount() uint {
	return conf.BindRetryCount
}

// SetBindRetryWaitTime sets a wait time when the binding is failed.
func (conf *Config) SetBindRetryWaitTime(d time.Duration) {
	conf.BindRetryWaitTime = d
}

// GetBindRetryWaitTime return the wait time when the binding is failed.
func (conf *Config) GetBindRetryWaitTime() time.Duration {
	return conf.BindRetryWaitTime
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *UnicastConfig) Equals(otherConf *UnicastConfig) bool {
	return reflect.DeepEqual(conf, otherConf)
}
