// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
uechodump is a dump utility for Echonet Lite.

	NAME
	uechodump

	SYNOPSIS
	uechodump [OPTIONS]

	DESCRIPTION
	uechodump is a dump utility for Echonet Lite.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
*/
package main

import (
	"os"
	"os/signal"
)

func main() {
	ctrl := NewPostController()

	err := ctrl.Start()
	if err != nil {
		return
	}

	err = ctrl.SearchAllObjects()
	if err != nil {
		return
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-c:
			ctrl.Stop()
			os.Exit(0)
		}
	}()

	err = ctrl.Stop()
	if err != nil {
		return
	}
}
