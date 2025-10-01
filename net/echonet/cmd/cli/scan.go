// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
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

		_, err := NewFormatFromString(viper.GetString(FormatParamStr))
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

		// Outputs all found nodes

		db := echonet.GetStandardDatabase()

		for i, node := range ctrl.Nodes() {

			// Gets manufacture code.

			manufactureName := unknown
			req := echonet.NewMessage()
			req.SetESV(echonet.ESVReadRequest)
			req.SetDEOJ(0x0EF001)
			req.AddProperty(echonet.NewProperty().SetCode(0x8A))
			res, err := ctrl.PostMessage(node, req)
			if err == nil {
				if props := res.Properties(); len(props) == 1 {
					manufacture, ok := db.FindManufacture(echonet.ManufactureCode(encoding.ByteToInteger(props[0].Data())))
					if ok {
						manufactureName = manufacture.Name()
					}
				}
			}

			// Prints node data.

			fmt.Printf("[%d] %-15s:%d (%s)\n", i, node.Address(), node.Port(), manufactureName)

			for j, obj := range node.Objects() {
				// Prints object data.

				objName := obj.ClassName()
				if len(objName) == 0 {
					objName = unknown
				}
				fmt.Printf("    [%d] %06X (%s)\n", j, obj.Code(), objName)

				// Prints only read required properties with the current property data.

				for _, prop := range obj.Properties() {
					if !prop.IsReadRequired() {
						continue
					}
					propName := prop.Name()
					if len(propName) == 0 {
						propName = "(" + unknown + ")"
					}
					propData := "--"
					req := echonet.NewMessage()
					req.SetESV(echonet.ESVReadRequest)
					req.SetDEOJ(obj.Code())
					req.AddProperty(echonet.NewProperty().SetCode(prop.Code()))
					res, err := ctrl.PostMessage(node, req)
					if err == nil {
						if props := res.Properties(); len(props) == 1 {
							propData = hex.EncodeToString(props[0].Data())
						}
					} else {
						propData = err.Error()
					}
					fmt.Printf("        [%02X] %s: %s\n", prop.Code(), propName, propData)
				}
			}
		}

		/*
			scanner := SharedCommissioner().Scannar()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = scanner.Scan(ctx)
			if err != nil {
				return err
			}
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

		// Stops the controller

		err = ctrl.Stop()
		if err != nil {
			return err
		}

		return nil
	},
}
