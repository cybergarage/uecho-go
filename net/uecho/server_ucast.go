// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

const ()

// uniCastServer is an instance for Echonet node.
type uniCastServer struct {
}

// newUniCastServer returns a new object.
func newUniCastServer() *uniCastServer {
	server := &uniCastServer{}
	return server
}

// Start starts the server.
func (server *uniCastServer) Start() error {
	return nil
}

// Stop stop the server.
func (server *uniCastServer) Stop() error {
	return nil
}
