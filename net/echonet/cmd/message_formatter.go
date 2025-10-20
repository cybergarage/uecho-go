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
	// Rows returns the data rows for the m
	Rows() [][]string
}

// defaultMessageFormatter provides a default implementation of MessageFormatter.
type defaultMessageFormatter struct {
	msgs []echonet.Message
}

// NewMessageFormatter returns a new default message formatter.
func NewMessageFormatter(msgs ...echonet.Message) MessageFormatter {
	return &defaultMessageFormatter{
		msgs: msgs,
	}
}

// Columns returns the column names for the message.
func (f *defaultMessageFormatter) Columns() []string {
	columns := []string{
		"EHD",
		"TID",
		"SEOJ",
		"DEOJ",
		"ESV",
		"OPC",
	}
	opc := 0
	for _, msg := range f.msgs {
		if opc < msg.OPC() {
			opc = msg.OPC()
		}
	}
	for n := range opc {
		columns = append(
			columns,
			fmt.Sprintf("EPC%d", n),
			fmt.Sprintf("PDC%d", n),
			fmt.Sprintf("EDT%d", n),
		)
	}
	return columns
}

func (f *defaultMessageFormatter) Rows() [][]string {
	rows := [][]string{}
	for _, msg := range f.msgs {
		ehd := msg.EHD()
		strs := []string{
			fmt.Sprintf("%02X%02X", ehd[0], ehd[1]),
			fmt.Sprintf("%04X", msg.TID()),
			msg.SEOJ().String(),
			msg.DEOJ().String(),
			msg.ESV().String(),
			fmt.Sprintf("%02X", msg.OPC()),
		}
		for _, prop := range msg.Properties() {
			strs = append(
				strs,
				fmt.Sprintf("%02X", prop.Code()),
				fmt.Sprintf("%02X", prop.Size()),
				fmt.Sprintf("%X", prop.Data()),
			)
		}
		rows = append(rows, strs)
	}
	return rows
}
