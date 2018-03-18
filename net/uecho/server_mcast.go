// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

const ()

// multiCastServer is an instance for Echonet node.
type multiCastServer struct {
}

// NewNode returns a new object.
func newMultiCastServer() *multiCastServer {
	server := &multiCastServer{}
	return server
}

// Start starts the server.
func (server *multiCastServer) Start() error {
	return nil
}

// Stop stop the server.
func (server *multiCastServer) Stop() error {
	return nil
}
