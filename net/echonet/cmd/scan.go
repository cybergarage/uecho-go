// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	unknown = "unknown"
)

func init() {
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "scan [address]",
	Short: "Scan for Echonet Lite devices.",
	Long:  "Scan for Echonet Lite devices.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose := viper.GetBool(VerboseParamStr)
		if verbose {
			enableStdoutVerbose(true)
		}

		format, err := NewFormatFromString(viper.GetString(FormatParamStr))
		if err != nil {
			return err
		}

		address := ""
		if 0 < len(args) {
			address = args[0]
		}

		// Creates a query

		query := &Query{
			Details: false,
			Address: address,
		}

		if verbose || 0 < len(address) {
			query.Details = true
		}

		// Starts a controller for Echonet Lite node

		ctrl := NewController()

		if verbose {
			ctrl.SetListener(ctrl)
		}

		err = ctrl.Start()
		if err != nil {
			return err
		}

		table, err := ctrl.DiscoverNodes(context.Background(), query)
		if err != nil {
			return err
		}

		printDevicesTable := func(tbl Table) {
			formatter := NewTableFormatter()
			tbl = formatter.HideDuplicateColumns(tbl, 0, 1, 2, 3, 4)
			columns, rows := tbl.Columns(), tbl.Rows()
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

		printDevicesCSV := func(tbl Table) {
			columns, rows := tbl.Columns(), tbl.Rows()
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

		printDevicesJSON := func(tbl Table) error {
			columns, rows := tbl.Columns(), tbl.Rows()
			devObjs := make([]map[string]string, 0, len(rows))
			for _, row := range rows {
				obj := make(map[string]string)
				for i, col := range columns {
					obj[col] = row[i]
				}
				devObjs = append(devObjs, obj)
			}
			b, err := json.MarshalIndent(devObjs, "", "  ")
			if err != nil {
				return err
			}
			outputf("%s\n", string(b))
			return nil
		}

		switch format {
		case FormatJSON:
			printDevicesJSON(table)
		case FormatCSV:
			printDevicesCSV(table)
		default:
			printDevicesTable(table)
		}

		// Stops the controller

		err = ctrl.Stop()
		if err != nil {
			return err
		}

		return nil
	},
}
