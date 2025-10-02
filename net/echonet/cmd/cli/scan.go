// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"encoding/json"
	"os"
	"strings"
	"text/tabwriter"
	"time"

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
	Use:   "scan",
	Short: "Scan for Echonet Lite devices.",
	Long:  "Scan for Echonet Lite devices.",
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose := viper.GetBool(VerboseParamStr)
		if verbose {
			enableStdoutVerbose(true)
		}

		format, err := NewFormatFromString(viper.GetString(FormatParamStr))
		if err != nil {
			return err
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

		err = ctrl.SearchAllObjects()
		if err != nil {
			return err
		}

		// Waits node responses in the local network

		time.Sleep(time.Second * 1)

		printDevicesTable := func(columns []string, rows [][]string) error {
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
			return nil
		}

		printDevicesCSV := func(columns []string, rows [][]string) error {
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
			return nil
		}

		printDevicesJSON := func(columns []string, rows [][]string) error {
			devObjs := make([]any, 0)
			b, err := json.MarshalIndent(devObjs, "", "  ")
			if err != nil {
				return err
			}
			outputf("%s\n", string(b))
			return nil
		}

		columns, rows, err := ctrl.DiscoveredNodeTable()
		if err != nil {
			return err
		}

		switch format {
		case FormatJSON:
			return printDevicesJSON(columns, rows)
		case FormatCSV:
			return printDevicesCSV(columns, rows)
		default:
			return printDevicesTable(columns, rows)
		}

		// Stops the controller

		err = ctrl.Stop()
		if err != nil {
			return err
		}

		return nil
	},
}
