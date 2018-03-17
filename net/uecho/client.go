// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uecho

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
)

// Client is an instance for Graphite protocols.
type Client struct {
}

// NewClient returns a new Client.
func NewClient() *Client {
	client := &Client{
		}
	return client
}
