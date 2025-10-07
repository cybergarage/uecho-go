// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet"
)

func outputf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func errorf(format string, args ...any) {
	fmt.Printf(format, args...)
}

func enableStdoutVerbose(flag bool) {
	if flag {
		log.SetSharedLogger(log.NewStdoutLogger(log.LevelInfo))
	} else {
		log.SetSharedLogger(nil)
	}
}

func hexStringToByte(hexStr string) ([]byte, error) {
	// Remove any whitespace and convert to lowercase for consistency
	cleanHex := strings.ReplaceAll(strings.ToLower(hexStr), " ", "")

	// Check if the string has an odd number of characters
	if len(cleanHex)%2 != 0 {
		return nil, fmt.Errorf("invalid hex string: odd number of characters")
	}

	// Decode the hex string to bytes first
	bytes, err := hex.DecodeString(cleanHex)
	if err != nil {
		return nil, fmt.Errorf("invalid hex string: %w", err)
	}

	return bytes, nil
}

func hexStringToInt(hexStr string) (int, error) {
	// Decode the hex string to bytes first
	bytes, err := hexStringToByte(hexStr)
	if err != nil {
		return 0, fmt.Errorf("invalid hex string: %w", err)
	}

	// Convert bytes to integer (big-endian)
	result := 0
	for _, b := range bytes {
		result = result<<8 + int(b)
	}

	return result, nil
}

func outputTransportMessage(prefix string, addr string, obj echonet.ObjectCode, msg echonet.Message) {
	fmt.Printf("%s %-15s : %06X %02X ",
		prefix,
		addr,
		obj,
		msg.ESV())
	for _, prop := range msg.Properties() {
		fmt.Printf("%2X%s ",
			prop.Code(),
			hex.EncodeToString(prop.Data()))
	}
	fmt.Printf("\n")
}

func outputResponseMessage(msg echonet.Message) {
	outputTransportMessage("<-", msg.SourceAddress(), msg.SEOJ(), msg)
}
