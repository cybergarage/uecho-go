// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		_, err := NewFormatFromString(viper.GetString(FormatParamStr))
		if err != nil {
			return err
		}
		return nil
		/*
			scanner := SharedCommissioner().Scannar()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = scanner.Scan(ctx)
			if err != nil {
				return err
			}
			columns := []string{"Name", "Addr", "VendorID", "ProductID", "Discriminator"}
			deviceColumns := func(dev ble.Device) ([]string, error) {
				service, err := dev.Service()
				if err != nil {
					return nil, err
				}
				return []string{
					dev.LocalName(),
					dev.Address().String(),
					strconv.Itoa(int(service.VendorID())),
					strconv.Itoa(int(service.ProductID())),
					strconv.Itoa(int(service.Discriminator())),
				}, nil
			}

			printDevicesTable := func(devs []ble.Device) error {
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
				for _, dev := range devs {
					devColumns, err := deviceColumns(dev)
					if err != nil {
						return err
					}
					printRow(devColumns...)
				}
				w.Flush()
				return nil
			}

			printDevicesCSV := func(devs []ble.Device) error {
				printRow := func(cols ...string) {
					if len(cols) == 0 {
						return
					}
					outputf("%s\n", strings.Join(cols, ","))
				}
				printRow(columns...)
				for _, dev := range devs {
					devColumns, err := deviceColumns(dev)
					if err != nil {
						return err
					}
					printRow(devColumns...)
				}
				return nil
			}

			printDevicesJSON := func(devs []ble.Device) error {
				devObjs := make([]any, 0)
				for _, dev := range devs {
					devObjs = append(devObjs, dev.MarshalObject())
				}
				b, err := json.MarshalIndent(devObjs, "", "  ")
				if err != nil {
					return err
				}
				outputf("%s\n", string(b))
				return nil
			}

			devs := scanner.DiscoveredDevices()
			if len(devs) == 0 {
				return nil
			}

			switch format {
			case FormatJSON:
				return printDevicesJSON(devs)
			case FormatCSV:
				return printDevicesCSV(devs)
			default:
				return printDevicesTable(devs)
			}
		*/

	},
}
