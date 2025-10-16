// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"encoding/json"
	"os"
	"strings"
	"text/tabwriter"
)

// Table represents a table for CLI output.
type Table interface {
	Columns() []string
	Rows() [][]string
	Output()
	OutputJSON()
	OutputCSV()
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

// Output outputs the table in a tabular format.
func (t *table) Output() {
	columns, rows := t.Columns(), t.Rows()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	printRow := func(cols ...string) {
		if len(cols) == 0 {
			return
		}
		for i, col := range cols {
			if i == len(cols)-1 {
				_, _ = w.Write([]byte(col + "\n"))
			} else {
				_, _ = w.Write([]byte(col + "\t"))
			}
		}
	}
	printRow(columns...)
	for _, row := range rows {
		printRow(row...)
	}
	w.Flush()
}

// OutputJSON outputs the table in JSON format.
func (t *table) OutputJSON() {
	data := make([]map[string]string, len(t.rows))
	for i, row := range t.rows {
		data[i] = make(map[string]string)
		for j, col := range t.columns {
			data[i][col] = row[j]
		}
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		errorf("failed to marshal table to JSON: %v", err)
		return
	}
	outputf("%s\n", jsonData)
}

// OutputCSV outputs the table in CSV format.
func (t *table) OutputCSV() {
	columns, rows := t.Columns(), t.Rows()
	printRow := func(cols ...string) {
		if len(cols) == 0 {
			return
		}
		outputf("%s\n", strings.Join(cols, ","))
	}
	printRow(columns...)
	for _, row := range rows {
		printRow(row...)
	}
}
