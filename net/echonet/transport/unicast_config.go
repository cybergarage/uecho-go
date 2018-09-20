// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// UnicastConfig represents a cofiguration for unicast server.
type UnicastConfig struct {
	TCPEnabled bool
	UDPEnabled bool
}

// NewDefaultUnicastConfig returns a default configuration.
func NewDefaultUnicastConfig() *UnicastConfig {
	conf := &UnicastConfig{
		TCPEnabled: true,
		UDPEnabled: true,
	}
	return conf
}

// SetConfig sets all flags.
func (conf *UnicastConfig) SetConfig(newConfig *UnicastConfig) {
	conf.TCPEnabled = newConfig.TCPEnabled
	conf.UDPEnabled = newConfig.UDPEnabled
}

// SetTCPEnabled sets a flag for TCP functions.
func (conf *UnicastConfig) SetTCPEnabled(flag bool) {
	conf.TCPEnabled = flag
}

// IsTCPEnabled returns true whether the TCP function is enabled, otherwise false.
func (conf *UnicastConfig) IsTCPEnabled() bool {
	return conf.TCPEnabled
}

// SetUDPEnabled sets a flag for UDP functions.
func (conf *UnicastConfig) SetUDPEnabled(flag bool) {
	conf.UDPEnabled = flag
}

// IsUDPEnabled returns true whether the UDP function is enabled, otherwise false.
func (conf *UnicastConfig) IsUDPEnabled() bool {
	return conf.UDPEnabled
}
