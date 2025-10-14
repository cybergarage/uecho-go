// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/cybergarage/uecho-go/net/echonet"
)

type MessageFormatter interface {
}

// defaultMessageFormatter provides a default implementation of MessageFormatter.
type defaultMessageFormatter struct {
	msg echonet.Message
}

// NewMessageFormatter returns a new default message formatter.
func NewMessageFormatter(msg echonet.Message) MessageFormatter {
	return &defaultMessageFormatter{
		msg: msg,
	}
}
