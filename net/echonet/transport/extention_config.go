// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import "reflect"

// ExtentionConfig represents a cofiguration for extended specifications.
type ExtentionConfig struct {
	AutoPortBindingEnabled      bool
	EachInterfaceBindingEnabled bool
	AutoInterfaceBindingEnabled bool
}

// NewDefaultExtentionConfig returns a default configuration.
func NewDefaultExtentionConfig() *ExtentionConfig {
	conf := &ExtentionConfig{
		AutoPortBindingEnabled:      false,
		EachInterfaceBindingEnabled: true,
		AutoInterfaceBindingEnabled: true,
	}
	return conf
}

// SetConfig sets all flags.
func (conf *ExtentionConfig) SetConfig(newConfig *ExtentionConfig) {
	conf.AutoPortBindingEnabled = newConfig.AutoPortBindingEnabled
	conf.EachInterfaceBindingEnabled = newConfig.EachInterfaceBindingEnabled
	conf.AutoInterfaceBindingEnabled = newConfig.AutoInterfaceBindingEnabled
}

// SetAutoPortBindingEnabled sets a flag for TCP functions.
func (conf *ExtentionConfig) SetAutoPortBindingEnabled(flag bool) {
	conf.AutoPortBindingEnabled = flag
}

// IsAutoPortBindingEnabled returns true whether the TCP function is enabled, otherwise false.
func (conf *ExtentionConfig) IsAutoPortBindingEnabled() bool {
	return conf.AutoPortBindingEnabled
}

// SetEachInterfaceBindingEnabled sets a flag for binding functions.
func (conf *ExtentionConfig) SetEachInterfaceBindingEnabled(flag bool) {
	conf.EachInterfaceBindingEnabled = flag
}

// IsEachInterfaceBindingEnabled returns true whether the binding functions is enabled, otherwise false.
func (conf *ExtentionConfig) IsEachInterfaceBindingEnabled() bool {
	return conf.EachInterfaceBindingEnabled
}

// SetAutoInterfaceBindingEnabled sets a flag for the auto interface binding.
func (conf *ExtentionConfig) SetAutoInterfaceBindingEnabled(flag bool) {
	conf.AutoInterfaceBindingEnabled = flag
}

// IsAutoInterfaceBindingEnabled returns true whether the the auto interface binding is enabled, otherwise false.
func (conf *ExtentionConfig) IsAutoInterfaceBindingEnabled() bool {
	return conf.AutoInterfaceBindingEnabled
}

// Equals returns true whether the specified other class is same, otherwise false.
func (conf *ExtentionConfig) Equals(otherConf *ExtentionConfig) bool {
	return reflect.DeepEqual(conf, otherConf)
}
