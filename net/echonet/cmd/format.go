// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"strings"
)

// Format represents the output format.
type Format int

const (
	FormatTable Format = iota
	FormatJSON
	FormatCSV
)

var (
	FormatParamStr = "format"
	FormatTableStr = "table"
	FormatJSONStr  = "json"
	FormatCSVStr   = "csv"
)

func allSupportedFormats() []string {
	return []string{
		FormatTableStr,
		FormatJSONStr,
		FormatCSVStr,
	}
}

var formatMap = map[string]Format{
	FormatTableStr: FormatTable,
	FormatJSONStr:  FormatJSON,
	FormatCSVStr:   FormatCSV,
}

// NewFormatFromString returns the format from the string.
func NewFormatFromString(s string) (Format, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if format, ok := formatMap[s]; ok {
		return format, nil
	}
	return FormatTable, fmt.Errorf("invalid format: %s", s)
}

// String returns the string representation of the format.
func (f Format) String() string {
	for k, v := range formatMap {
		if v == f {
			return k
		}
	}
	return "unknown"
}
