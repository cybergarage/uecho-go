// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

// Table represents a table for CLI output.
type Table interface {
	Columns() []string
	Rows() [][]string
	TableHelper
}

type TableHelper interface {
	StripDuplicateRowColumns(indexes ...int) Table
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

func (t *table) StripDuplicateRowColumns(columnIdxes ...int) Table {
	stripDuplicateRowColumns := func(rows [][]string, columnIdx int) [][]string {
		uniqRows := [][]string{}
		lastRowContext := ""
		for rowIdx, row := range rows {
			uniqRow := make([]string, len(row))
			copy(uniqRow, row)
			switch {
			case rowIdx == 0:
				lastRowContext = row[columnIdx]
			default:
				if row[columnIdx] == lastRowContext {
					if 0 < columnIdx && row[columnIdx-1] == "" {
						uniqRow[columnIdx-1] = ""
					}
				} else {
					lastRowContext = row[columnIdx]
				}
			}
			uniqRows = append(uniqRows, uniqRow)
		}
		return uniqRows
	}

	uniqRows := t.rows
	for _, columnIdx := range columnIdxes {
		uniqRows = stripDuplicateRowColumns(uniqRows, columnIdx)
	}
	return NewTable(t.columns, uniqRows)
}
