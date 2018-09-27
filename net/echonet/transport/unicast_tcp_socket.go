// Copyright 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

import (
	"time"

	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

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

// PostMessage send a message to the destination address.
func (sock *UnicastTCPSocket) PostMessage(addr string, port int, msg *protocol.Message, timeout time.Duration) (int, error) {
	return sock.Write(addr, port, msg.Bytes(), timeout)
}
