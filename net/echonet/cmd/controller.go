// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cybergarage/uecho-go/net/echonet"
	"github.com/cybergarage/uecho-go/net/echonet/encoding"
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

// Controller represents an Echonet Lite controller.
type Controller struct {
	echonet.Controller
}

// NewController returns a new controller.
func NewController() *Controller {
	c := &Controller{
		Controller: echonet.NewController(),
	}
	return c
}

// DiscoverNodes searches for Echonet Lite nodes and returns the discovered node table.
func (ctrl *Controller) DiscoverNodes(ctx context.Context, query *Query) (Table, error) {
	if err := ctrl.Controller.Search(ctx); err != nil {
		return nil, err
	}
	return ctrl.DiscoveredNodeTable(query)
}

// DiscoveredNodeTable returns the discovered node table.
func (ctrl *Controller) DiscoveredNodeTable(query *Query) (Table, error) {
	db := echonet.SharedStandardDatabase()

	cols := []string{
		"address",
		"port",
		"manufacture",
		"object_code",
		"object_name",
	}
	if query.Details {
		cols = append(cols,
			"property_code",
			"property_name",
			"property_attr",
			"property_data",
		)
	}

	rows := [][]string{}

	for _, node := range ctrl.Nodes() {

		// Filters by address

		if 0 < len(query.Address) && node.Address() != query.Address {
			continue
		}

		// Gets manufacture code.

		manufactureName := unknown
		req := echonet.NewMessage(
			echonet.WithMessageESV(echonet.ESVReadRequest),
			echonet.WithMessageDEOJ(0x0EF001),
			echonet.WithMessageProperties(
				echonet.NewProperty(
					echonet.WithPropertyCode(0x8A)),
			),
		)
		res, err := ctrl.PostMessage(context.Background(), node, req)
		if err == nil {
			if props := res.Properties(); len(props) == 1 {
				manufacture, ok := db.LookupManufacture(echonet.ManufactureCode(encoding.ByteToInteger(props[0].Data())))
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

			if !query.Details {
				row := []string{
					node.Address(),
					fmt.Sprintf("%d", node.Port()),
					manufactureName,
					fmt.Sprintf("%06X", obj.Code()),
					objName,
				}
				rows = append(rows, row)
				continue
			}

			for _, prop := range obj.Properties() {

				propName := prop.Name()
				if len(propName) == 0 {
					propName = "(" + unknown + ")"
				}

				propAttrString := func(attr echonet.PropertyAttribute, s string) string {
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
					req := echonet.NewMessage(
						echonet.WithMessageESV(echonet.ESVReadRequest),
						echonet.WithMessageDEOJ(obj.Code()),
						echonet.WithMessageProperties(
							echonet.NewProperty(
								echonet.WithPropertyCode(prop.Code()),
							),
						),
					)
					res, err := ctrl.PostMessage(context.Background(), node, req)
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
					fmt.Sprintf("%d", node.Port()),
					manufactureName,
					fmt.Sprintf("%06X", obj.Code()),
					objName,
					fmt.Sprintf("%02X", prop.Code()),
					propName,
					propAttr,
					propData,
				}

				rows = append(rows, row)
			}
		}
	}

	return NewTable(cols, rows), nil
}

// ControllerMessageReceived is called when a message is received.
func (ctrl *Controller) ControllerMessageReceived(msg *protocol.Message) {
	// log.Infof("%s : %s\n", msg.From.String(), hex.EncodeToString(msg.Bytes()))
}

// ControllerNewNodeFound is called when a new node is found.
func (ctrl *Controller) ControllerNewNodeFound(echonet.Node) {
}
