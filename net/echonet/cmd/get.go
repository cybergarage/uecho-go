// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"

	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{ // nolint:exhaustruct
	Use:     "get <node-address> <object-code> <property-code>",
	Short:   "Get property value from Echonet Lite device.",
	Long:    "Get property value from Echonet Lite device. Object and property codes must be specified in hexadecimal format.",
	Example: "  uechoctl get 192.168.1.100 013001 80\n  uechoctl get 10.0.0.50 028001 B0",
	Args:    cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose := viper.GetBool(VerboseParamStr)
		if verbose {
			enableStdoutVerbose(true)
		}

		// Parses arguments

		address := args[0]

		objCode, err := hexStringToInt(args[1])
		if err != nil {
			return err
		}

		propCode, err := hexStringToInt(args[2])
		if err != nil {
			return err
		}

		prop, err := echonet.NewPropertyWith(
			echonet.WithPropertyCode(echonet.PropertyCode(propCode)),
		)
		if err != nil {
			return err
		}

		reqMsg := echonet.NewMessageWith(echonet.ObjectCode(objCode), echonet.ESVReadRequest, prop)

		// Create a controller

		ctrl := NewController()

		if verbose {
			ctrl.SetListener(ctrl)
		}

		err = ctrl.Start()
		if err != nil {
			return err
		}

		// Lookup node

		ctx := context.Background()
		err = ctrl.Search(ctx)
		if err != nil {
			return err
		}

		node, ok := ctrl.LookupNode(address)
		if !ok {
			return fmt.Errorf("node not found: %s", address)
		}

		// Post the specified request message to the destination node

		resMsg, err := ctrl.PostMessage(context.Background(), node, reqMsg)
		if err != nil {
			return err
		}
		outputResponseMessage(resMsg)

		// Stops the controller

		err = ctrl.Stop()
		if err != nil {
			return err
		}

		return nil
	},
}
