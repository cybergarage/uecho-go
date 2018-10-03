// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// ExtentionConfig represents a cofiguration for extended specifications.
type ExtentionConfig struct {
	AutoBindingEnabled bool
}

// NewDefaultExtentionConfig returns a default configuration.
func NewDefaultExtentionConfig() *ExtentionConfig {
	conf := &ExtentionConfig{
		AutoBindingEnabled: false,
	}
	return conf
}

// SetConfig sets all flags.
func (conf *ExtentionConfig) SetConfig(newConfig *ExtentionConfig) {
	conf.AutoBindingEnabled = newConfig.AutoBindingEnabled
}

// SetAutoBindingEnabled sets a flag for TCP functions.
func (conf *ExtentionConfig) SetAutoBindingEnabled(flag bool) {
	conf.AutoBindingEnabled = flag
}

// IsAutoBindingEnabled returns true whether the TCP function is enabled, otherwise false.
func (conf *ExtentionConfig) IsAutoBindingEnabled() bool {
	return conf.AutoBindingEnabled
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *ExtentionConfig) Equals(otherConf *ExtentionConfig) bool {
	if conf.IsAutoBindingEnabled() != otherConf.IsAutoBindingEnabled() {
		return false
	}
	return true
}
