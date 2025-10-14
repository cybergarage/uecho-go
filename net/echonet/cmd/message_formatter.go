// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet"
)

// MessageFormatter is an interface for formatting Echonet messages.
type MessageFormatter interface {
	// Columns returns the column names for the message.
	Columns() []string
	// HexStrings returns the hex string representation of the message.
	HexStrings() []string
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

// Columns returns the column names for the message.
func (f *defaultMessageFormatter) Columns() []string {
	columns := []string{
		"SEOJ",
		"DEOJ",
		"ESV",
		"OPC",
	}
	for n := range f.msg.OPC() {
		columns = append(
			columns,
			fmt.Sprintf("EPC%d", n),
			fmt.Sprintf("PDC%d", n),
			fmt.Sprintf("EDT%d", n),
		)
	}
	return columns
}

// HexStrings returns the hex string representation of the message.
func (f *defaultMessageFormatter) HexStrings() []string {
	strs := []string{
		f.msg.SEOJ().String(),
		f.msg.DEOJ().String(),
		f.msg.ESV().String(),
		fmt.Sprintf("%02X", f.msg.OPC()),
	}
	for _, prop := range f.msg.Properties() {
		strs = append(
			strs,
			fmt.Sprintf("%02X", prop.Code()),
			fmt.Sprintf("%02X", prop.Size()),
			fmt.Sprintf("%X", prop.Data()),
		)
	}
	return strs
}
