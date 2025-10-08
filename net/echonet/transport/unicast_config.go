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
	tcpEnabled        bool
	connectionTimeout time.Duration
	requestTimeout    time.Duration
	bindRetryCount    uint
	bindRetryWaitTime time.Duration
}

// NewDefaultUnicastConfig returns a default configuration.
func NewDefaultUnicastConfig() *UnicastConfig {
	conf := &UnicastConfig{
		tcpEnabled:        false,
		connectionTimeout: DefaultConnectimeTimeOut,
		requestTimeout:    DefaultRequestTimeout,
		bindRetryCount:    DefaultBindRetryCount,
		bindRetryWaitTime: 0,
	}
	return conf
}

// SetConfig sets all flags.
func (conf *UnicastConfig) SetConfig(newConfig *UnicastConfig) {
	conf.tcpEnabled = newConfig.tcpEnabled
	conf.connectionTimeout = newConfig.connectionTimeout
}

// SetTCPEnabled sets a flag for TCP functions.
func (conf *UnicastConfig) SetTCPEnabled(flag bool) {
	conf.tcpEnabled = flag
}

// TCPEnabled returns true whether the TCP function is enabled, otherwise false.
func (conf *UnicastConfig) TCPEnabled() bool {
	return conf.tcpEnabled
}

// SetConnectionTimeout sets a connection timeout setting.
func (conf *UnicastConfig) SetConnectionTimeout(timeout time.Duration) {
	conf.connectionTimeout = timeout
}

// ConnectionTimeout returns the current connection timeout setting.
func (conf *UnicastConfig) ConnectionTimeout() time.Duration {
	return conf.connectionTimeout
}

// SetRequestTimeout sets a request timeout.
func (conf *Config) SetRequestTimeout(d time.Duration) {
	conf.requestTimeout = d
}

// RequestTimeout return the request timeout.
func (conf *Config) RequestTimeout() time.Duration {
	return conf.requestTimeout
}

// SetBindRetryEnabled sets a flag for the bind retry function.
func (conf *UnicastConfig) SetBindRetryEnabled(flag bool) {
	if !flag {
		conf.bindRetryCount = 0
		return
	}
	conf.bindRetryCount = DefaultBindRetryCount
}

// BindRetryEnabled returns true whether the bind retry function is enabled, otherwise false.
func (conf *UnicastConfig) BindRetryEnabled() bool {
	return conf.bindRetryCount != 0
}

// SetBindRetryCount sets a retry count when the binding is failed.
func (conf *Config) SetBindRetryCount(n uint) {
	conf.bindRetryCount = n
}

// BindRetryCount returns the retry count when the binding is failed.
func (conf *Config) BindRetryCount() uint {
	return conf.bindRetryCount
}

// SetBindRetryWaitTime sets a wait time when the binding is failed.
func (conf *Config) SetBindRetryWaitTime(d time.Duration) {
	conf.bindRetryWaitTime = d
}

// BindRetryWaitTime return the wait time when the binding is failed.
func (conf *Config) BindRetryWaitTime() time.Duration {
	return conf.bindRetryWaitTime
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *UnicastConfig) Equals(otherConf *UnicastConfig) bool {
	return reflect.DeepEqual(conf, otherConf)
}
