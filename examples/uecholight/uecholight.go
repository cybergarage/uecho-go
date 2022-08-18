// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
uechosearch is a search utility for Echonet Lite.

	NAME
	uecholight

	SYNOPSIS
	uecholight [OPTIONS]

	DESCRIPTION
	uecholight is a sample implementation of mono functional lighting device

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/uecho-go/net/echonet"
)

// See : APPENDIX Detailed Requirements for ECHONET Device objects
//       3.3.29 Requirements for mono functional lighting class

func main() {
	verbose := flag.Bool("v", false, "Enable verbose output")
	manufacturerCode := flag.Int("m", echonet.ObjectManufacturerUnknown, "Set manufacturer code")
	flag.Parse()

	// Setup logger

	if *verbose {
		log.SetSharedLogger(log.NewStdoutLogger(log.LevelTrace))
	}

	// Start a light node for Echonet Lite

	node := NewLightNode()
	node.SetManufacturerCode(uint(*manufacturerCode))

	err := node.Start()
	if err != nil {
		OutputError(err)
		os.Exit(EXIT_FAILURE)
	}

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM)

	exitCh := make(chan int)

	go func() {
		for {
			s := <-sigCh
			switch s {
			case syscall.SIGHUP:
				err = node.Restart()
				if err != nil {
					OutputError(err)
					os.Exit(EXIT_FAILURE)
				}
			case syscall.SIGINT, syscall.SIGTERM:
				err = node.Stop()
				if err != nil {
					OutputError(err)
					os.Exit(EXIT_FAILURE)
				}
				exitCh <- EXIT_SUCCESS
			}
		}
	}()

	code := <-exitCh

	os.Exit(code)

}
