// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

// ControllerListener is a listener for Echonet messages.
type ControllerListener interface {
	addedNewNode(*RemoteNode)
}
