// Copyright 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transport

// A MessageManager represents a multicast server list.
type MessageManager struct {
}

// NewMessageManager returns a new message manager.
func NewMessageManager() *MessageManager {
	mgr := &MessageManager{}
	return mgr
}

// Start starts this server.
func (mgr *MessageManager) Start() error {
	return nil
}

// Stop stops this server.
func (mgr *MessageManager) Stop() error {
	return nil
}
