// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// ExtentionConfig represents a cofiguration for extended specifications.
type ExtentionConfig struct {
	AutoBindingEnabled          bool
	EachInterfaceBindingEnabled bool
}

// NewDefaultExtentionConfig returns a default configuration.
func NewDefaultExtentionConfig() *ExtentionConfig {
	conf := &ExtentionConfig{
		AutoBindingEnabled:          true,
		EachInterfaceBindingEnabled: true,
	}
	return conf
}

// SetConfig sets all flags.
func (conf *ExtentionConfig) SetConfig(newConfig *ExtentionConfig) {
	conf.AutoBindingEnabled = newConfig.AutoBindingEnabled
	conf.EachInterfaceBindingEnabled = newConfig.EachInterfaceBindingEnabled
}

// SetAutoBindingEnabled sets a flag for TCP functions.
func (conf *ExtentionConfig) SetAutoBindingEnabled(flag bool) {
	conf.AutoBindingEnabled = flag
}

// IsAutoBindingEnabled returns true whether the TCP function is enabled, otherwise false.
func (conf *ExtentionConfig) IsAutoBindingEnabled() bool {
	return conf.AutoBindingEnabled
}

// SetEachInterfaceBindingEnabled sets a flag for binding functions.
func (conf *ExtentionConfig) SetEachInterfaceBindingEnabled(flag bool) {
	conf.EachInterfaceBindingEnabled = flag
}

// IsEachInterfaceBindingEnabled returns true whether the binding functions is enabled, otherwise false.
func (conf *ExtentionConfig) IsEachInterfaceBindingEnabled() bool {
	return conf.EachInterfaceBindingEnabled
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *ExtentionConfig) Equals(otherConf *ExtentionConfig) bool {
	if conf.IsAutoBindingEnabled() != otherConf.IsAutoBindingEnabled() {
		return false
	}
	if conf.IsEachInterfaceBindingEnabled() != otherConf.IsEachInterfaceBindingEnabled() {
		return false
	}
	return true
}
