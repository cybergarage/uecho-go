// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// Config represents a cofiguration for transport.
type Config struct {
	TCPEnabled bool
	UDPEnabled bool
}

// NewDefaultConfig returns a default configuration.
func NewDefaultConfig() *Config {
	conf := &Config{
		TCPEnabled: true,
		UDPEnabled: true,
	}
	return conf
}

// SetConfig sets all flags.
func (conf *Config) SetConfig(flags *Config) {
	conf.TCPEnabled = flags.TCPEnabled
	conf.UDPEnabled = flags.UDPEnabled
}

// SetTCPEnabled sets a flag for TCP functions.
func (conf *Config) SetTCPEnabled(flag bool) {
	conf.TCPEnabled = flag
}

// IsTCPEnabled returns true whether the TCP function is enabled, otherwise false.
func (conf *Config) IsTCPEnabled() bool {
	return conf.TCPEnabled
}

// SetUDPEnabled sets a flag for UDP functions.
func (conf *Config) SetUDPEnabled(flag bool) {
	conf.UDPEnabled = flag
}

// IsUDPEnabled returns true whether the UDP function is enabled, otherwise false.
func (conf *Config) IsUDPEnabled() bool {
	return conf.UDPEnabled
}
