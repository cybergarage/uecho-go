// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

// TableFormatter provides an interface for formatting and transforming tables.
type TableFormatter interface {
	// HideDuplicateColumns removes duplicate columns from the table.
	HideDuplicateColumns(table Table, columnIdxes ...int) Table
	// FilterEmptyRows removes empty rows from the table.
	FilterEmptyRows(table Table) Table
}

// defaultTableFormatter provides a default implementation of TableFormatter.
type defaultTableFormatter struct{}

// NewTableFormatter returns a new default table formatter.
func NewTableFormatter() TableFormatter {
	return &defaultTableFormatter{}
}

// HideDuplicateColumns removes duplicate columns from the table.
func (f *defaultTableFormatter) HideDuplicateColumns(table Table, columnIdxes ...int) Table {
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
					switch columnIdx {
					case 0:
						uniqRow[columnIdx] = ""
					default:
						if uniqRow[columnIdx-1] == "" {
							uniqRow[columnIdx] = ""
						}
					}
				} else {
					lastRowContext = row[columnIdx]
				}
			}
			uniqRows = append(uniqRows, uniqRow)
		}
		return uniqRows
	}

	uniqRows := table.Rows()
	for _, columnIdx := range columnIdxes {
		uniqRows = stripDuplicateRowColumns(uniqRows, columnIdx)
	}
	return f.FilterEmptyRows(NewTable(table.Columns(), uniqRows))
}

// FilterEmptyRows removes empty rows from the table.
func (f *defaultTableFormatter) FilterEmptyRows(table Table) Table {
	uniqRows := [][]string{}
	for _, row := range table.Rows() {
		isBlankRow := true
		for _, cell := range row {
			if cell != "" {
				isBlankRow = false
				break
			}
		}
		if !isBlankRow {
			uniqRows = append(uniqRows, row)
		}
	}
	return NewTable(table.Columns(), uniqRows)
}
