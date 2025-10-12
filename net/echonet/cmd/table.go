// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

// Table represents a table for CLI output.
type Table interface {
	Columns() []string
	Rows() [][]string
}

type table struct {
	columns []string
	rows    [][]string
}

// NewTable returns a new table instance.
func NewTable(columns []string, rows [][]string) Table {
	return &table{
		columns: columns,
		rows:    rows,
	}
}

// Columns returns the column names of the table.
func (t *table) Columns() []string {
	return t.columns
}

// Rows returns the rows of the table.
func (t *table) Rows() [][]string {
	return t.rows
}
