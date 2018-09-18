// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/transport"
)

type transportConfig = transport.Config

// Config represents a cofiguration for transport.
type Config struct {
	*transportConfig
}

// NewDefaultConfig returns a default configuration.
func NewDefaultConfig() *Config {
	conf := &Config{
		transportConfig: transport.NewDefaultConfig(),
	}
	return conf
}
