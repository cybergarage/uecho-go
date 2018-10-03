// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// Config represents a cofiguration for transport.
type Config struct {
	*UnicastConfig
	*ExtentionConfig
}

// NewDefaultConfig returns a default configuration.
func NewDefaultConfig() *Config {
	conf := &Config{
		UnicastConfig:   NewDefaultUnicastConfig(),
		ExtentionConfig: NewDefaultExtentionConfig(),
	}
	return conf
}

// SetConfig sets all configuration flags.
func (conf *Config) SetConfig(newConfig *Config) {
	conf.UnicastConfig.SetConfig(newConfig.UnicastConfig)
	conf.ExtentionConfig.SetConfig(newConfig.ExtentionConfig)
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *Config) Equals(other *Config) bool {
	if !conf.UnicastConfig.Equals(other.UnicastConfig) {
		return false
	}
	if !conf.ExtentionConfig.Equals(other.ExtentionConfig) {
		return false
	}
	return true
}
