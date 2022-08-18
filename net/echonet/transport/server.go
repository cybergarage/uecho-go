// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// A Server represents a server.
type Server struct {
}

// NewServer returns a new UnicastServer.
func NewServer() *Server {
	server := &Server{}
	return server
}
