// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
uechodump is a dump utility for Echonet Lite.

	NAME
	uechopost

	SYNOPSIS
	uechopost [options] <address> <obj> <esv> <property (epc, pdc, edt) ...

	DESCRIPTION
	uechopost is a controller utility to send any messages to Echonet Lite nodes.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cybergarage/uecho-go/net/echonet"
)

const (
	EXIT_SUCCESS = 0
	EXIT_FAIL    = 1
)

func outputTransportMessage(prefix string, addr string, obj echonet.ObjectCode, msg *echonet.Message) {
	fmt.Printf("%s %-15s : %06X %02X ",
		prefix,
		addr,
		obj,
		msg.GetESV())
	for _, prop := range msg.GetProperties() {
		fmt.Printf("%2X%s ",
			prop.GetCode(),
			hex.EncodeToString(prop.GetData()))
	}
	fmt.Printf("\n")
}

func outputRequestMessage(ctrl *PostController, msg *echonet.Message) {
	sourceAddr := ""
	boundAddrs := ctrl.GetBoundAddresses()
	if 0 < len(boundAddrs) {
		sourceAddr = boundAddrs[0]
	}
	outputTransportMessage("->", sourceAddr, msg.GetDestinationObjectCode(), msg)
}

func outputResponseMessage(msg *echonet.Message) {
	outputTransportMessage("<-", msg.GetSourceAddress(), msg.GetSourceObjectCode(), msg)
}

func outputUsage() {
	program := filepath.Base(os.Args[0])
	fmt.Printf("Usage : %s [options] <address> <obj> <esv> <property (code, data) ...>\n", program)
}

func exitWithErrorMessage(errMsg string) {
	fmt.Printf("ERROR : %s\n", errMsg)
	os.Exit(EXIT_FAIL)
}

func exitWithError(err error) {
	exitWithErrorMessage(err.Error())
}

func main() {
	ctrl := NewPostController()

	argc := len(os.Args)
	if argc < 5 {
		outputUsage()
		os.Exit(EXIT_FAIL)
		return
	}

	err := ctrl.Start()
	if err != nil {
		exitWithError(err)
		return
	}

	err = ctrl.SearchAllObjects()
	if err != nil {
		exitWithError(err)
		return
	}

	// Wait node responses in the local network

	time.Sleep(time.Second * 1)

	// Find the specified destination node

	dstNodeAddr := os.Args[1]
	dstNode, err := ctrl.GetNode(dstNodeAddr)
	if err != nil {
		exitWithErrorMessage(fmt.Sprintf("The destination node (%s) is not found", dstNodeAddr))
	}

	// Find the specified destination object in the found node

	dstObjStr := os.Args[2]
	dstObjVal, err := strconv.ParseUint(dstObjStr, 16, 32)
	if err != nil {
		exitWithErrorMessage(fmt.Sprintf("The destination object (%s) is invalid", dstObjStr))
	}
	dstObjCode := echonet.ObjectCode(dstObjVal)
	_, err = dstNode.FindObject(dstObjCode)
	if err != nil {
		exitWithErrorMessage(fmt.Sprintf("The destination object (%06X) is not found", dstObjCode))
	}

	// Create a request message of the specified ESV and properties

	esvStr := os.Args[3]
	esvVal, err := strconv.ParseUint(esvStr, 16, 8)
	if (err != nil) || !echonet.IsValidESV(echonet.ESV(esvVal)) {
		exitWithErrorMessage(fmt.Sprintf("The ESV (%s) is invalid", esvStr))
	}
	esv := echonet.ESV(esvVal)

	props := make([]*echonet.Property, 0)
	for n := 4; n < len(os.Args); n++ {
		propHexBytes := []byte(os.Args[n])
		propBytes := make([]byte, hex.DecodedLen(len(propHexBytes)))
		_, err := hex.Decode(propBytes, propHexBytes)
		if err != nil {
			exitWithErrorMessage(fmt.Sprintf("The property code (%s) is invalid", string(propHexBytes)))
		}
		if len(propBytes) <= 1 {
			exitWithErrorMessage(fmt.Sprintf("The property code (%s) is short", string(propHexBytes)))
		}

		prop := echonet.NewProperty()
		prop.SetCode(echonet.PropertyCode(propBytes[0]))
		if 2 <= len(propBytes) {
			prop.SetData(propBytes[1:])
		}
		props = append(props, prop)
	}

	reqMsg := echonet.NewMessageWithParameters(dstObjCode, esv, props)

	// Send the specified request message to the destination node

	outputRequestMessage(ctrl, reqMsg)

	if reqMsg.IsResponseRequired() {
		resMsg, err := ctrl.PostMessage(dstNode, reqMsg)
		if err != nil {
			exitWithError(err)
		}
		outputResponseMessage(resMsg)
	} else {
		err := ctrl.SendMessage(dstNode, reqMsg)
		if err != nil {
			exitWithError(err)
		}
	}

	// Stop the controller

	err = ctrl.Stop()
	if err != nil {
		return
	}

	os.Exit(EXIT_SUCCESS)
}
