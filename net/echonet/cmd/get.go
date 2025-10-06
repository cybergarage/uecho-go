// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{ // nolint:exhaustruct
	Use:   "get [address] [object-code] [property-code]",
	Short: "Get property value from Echonet Lite device.",
	Long:  "Get property value from Echonet Lite device.",
	Args:  cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose := viper.GetBool(VerboseParamStr)
		if verbose {
			enableStdoutVerbose(true)
		}

		ctrl := NewController()

		if verbose {
			ctrl.SetListener(ctrl)
		}

		err := ctrl.Start()
		if err != nil {
			return err
		}

		// Stops the controller

		err = ctrl.Stop()
		if err != nil {
			return err
		}

		return nil
	},
}
