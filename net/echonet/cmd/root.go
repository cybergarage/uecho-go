// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"strings"

	"github.com/cybergarage/uecho-go/net/echonet"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	VerboseParamStr = "verbose"
)

var rootCmd = &cobra.Command{ // nolint:exhaustruct
	Use:               "uechoctl",
	Version:           echonet.Version,
	Short:             "",
	Long:              "",
	DisableAutoGenTag: true,
}

func GetRootCommand() *cobra.Command {
	return rootCmd
}

func Execute() error {
	err := rootCmd.Execute()
	return err
}

func init() {
	viper.SetEnvPrefix("UECHO")

	viper.SetDefault(FormatParamStr, FormatTableStr)
	rootCmd.PersistentFlags().String(FormatParamStr, FormatTableStr, fmt.Sprintf("output format: %s", strings.Join(allSupportedFormats(), "|")))
	viper.BindPFlag(FormatParamStr, rootCmd.PersistentFlags().Lookup(FormatParamStr))
	viper.BindEnv(FormatParamStr) // UECHO_FORMAT

	viper.SetDefault(VerboseParamStr, false)
	rootCmd.PersistentFlags().Bool((VerboseParamStr), false, "enable verbose output")
	viper.BindPFlag(VerboseParamStr, rootCmd.PersistentFlags().Lookup(VerboseParamStr))
	viper.BindEnv(VerboseParamStr) // UECHO_VERBOSE
}
