// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// Controller represents an Echonet Lite controller.
type Controller struct {
	*echonet.Controller
}

// NewController returns a new controller.
func NewController() *Controller {
	c := &Controller{
		Controller: echonet.NewController(),
	}
	return c
}

// DiscoveredNodeTable returns the discovered node table.
func (ctrl *Controller) DiscoveredNodeTable() ([]string, [][]string, error) {
	db := echonet.GetStandardDatabase()

	cols := []string{
		"address",
		"port",
		"manufacture",
		"object_code",
		"object_name",
		"property_code",
		"property_name",
		"property_attribute",
		"property_data",
	}
	rows := [][]string{}

	for n, node := range ctrl.Nodes() {

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

		for _, obj := range node.Objects() {
			objName := obj.ClassName()
			if len(objName) == 0 {
				objName = unknown
			}

			for _, prop := range obj.Properties() {

				propName := prop.Name()
				if len(propName) == 0 {
					propName = "(" + unknown + ")"
				}

				propAttrString := func(attr echonet.PropertyAttr, s string) string {
					switch attr {
					case echonet.Required:
						return strings.ToUpper(s)
					case echonet.Optional:
						return s
					default:
						return "-"
					}
				}

				propAttr := propAttrString(prop.ReadAttribute(), "r")
				propAttr += propAttrString(prop.WriteAttribute(), "w")
				propAttr += propAttrString(prop.AnnoAttribute(), "a")

				propData := ""
				if prop.IsReadRequired() {
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
				}

				row := []string{
					node.Address(),
					strconv.Itoa(node.Port()),
					manufactureName,
					fmt.Sprintf("%06X", obj.Code()),
					objName,
					fmt.Sprintf("%02X", prop.Code()),
					propName,
					propAttr,
					propData,
				}

				// Skip duplicate values for better readability.
				if 0 < n {
					prevRows := rows[len(rows)-1]
					skipColumnIndexes := []int{0, 3}
					for _, skipColumnIndex := range skipColumnIndexes {
						if prevRows[skipColumnIndex] != row[skipColumnIndex] {
							break
						}
						switch skipColumnIndex {
						case 0:
							row[0] = ""
							row[1] = ""
							row[2] = ""
						case 3:
							row[3] = ""
							row[4] = ""
						}
					}

				}

				rows = append(rows, row)
			}
		}
	}

	return cols, rows, nil
}

// ControllerMessageReceived is called when a message is received.
func (ctrl *Controller) ControllerMessageReceived(msg *protocol.Message) {
	// log.Infof("%s : %s\n", msg.From.String(), hex.EncodeToString(msg.Bytes()))
}

// ControllerNewNodeFound is called when a new node is found.
func (ctrl *Controller) ControllerNewNodeFound(*echonet.RemoteNode) {
}
