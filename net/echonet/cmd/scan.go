// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"

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

		switch format {
		case FormatJSON:
			table.OutputJSON()
		case FormatCSV:
			table.OutputCSV()
		default:
			formatter := NewTableFormatter(table)
			table = formatter.HideDuplicateColumns(0, 1, 2, 3, 4)
			table.Output()
		}

		// Stops the controller

		err = ctrl.Stop()
		if err != nil {
			return err
		}

		return nil
	},
}
