// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// A UnicastTCPSocket represents a socket.
type UnicastTCPSocket struct {
	*TCPSocket
}

// NewUnicastTCPSocket returns a new UnicastTCPSocket.
func NewUnicastTCPSocket() *UnicastTCPSocket {
	sock := &UnicastTCPSocket{
		TCPSocket: NewTCPSocket(),
	}
	return sock
}
