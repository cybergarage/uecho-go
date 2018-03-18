// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

const ()

// server is an instance for Echonet node.
type server struct {
	*multiCastServer
	*uniCastServer
}

// NewNode returns a new object.
func newServer() *server {
	server := &server{
		multiCastServer: newMultiCastServer(),
		uniCastServer:   newUniCastServer(),
	}
	return server
}

// Start starts the server.
func (server *server) Start() error {
	err := server.multiCastServer.Start()
	if err != nil {
		return nil
	}

	err = server.uniCastServer.Start()
	if err != nil {
		return nil
	}

	return nil
}

// Stop stop the server.
func (server *server) Stop() error {
	err := server.multiCastServer.Stop()
	if err != nil {
		return nil
	}

	err = server.uniCastServer.Stop()
	if err != nil {
		return nil
	}

	return nil
}
