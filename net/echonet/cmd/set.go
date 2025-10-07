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
	rootCmd.AddCommand(setCmd)
}

var setCmd = &cobra.Command{ // nolint:exhaustruct
	Use:     "set <node-address> <object-code> <property-code> <property-value>",
	Short:   "Set property value to Echonet Lite device.",
	Long:    "Set property value to Echonet Lite device. Object code, property code and property value must be specified in hexadecimal format.",
	Example: "set 192.168.1.100 013001 80 30",
	Args:    cobra.MinimumNArgs(4),
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

		propData, err := hexStringToByte(args[3])
		if err != nil {
			return err
		}

		prop := echonet.NewProperty(
			echonet.WithPropertyCode(echonet.PropertyCode(propCode)),
			echonet.WithPropertyData(propData),
		)

		reqMsg := echonet.NewMessage(
			echonet.WithMessageDEOJ(echonet.ObjectCode(objCode)),
			echonet.WithMessageESV(echonet.ESVWriteRequest),
			echonet.WithMessageProperties(prop),
		)

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

		// Send the specified request message to the destination node

		err = ctrl.SendMessage(context.Background(), node, reqMsg)
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
