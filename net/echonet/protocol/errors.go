// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

const (
	errorShortMessageSize     = "short message length : %d < %d"
	errorInvalidMessageHeader = "invalid message header [%d] : %02X != %02X"
	errorInvalidAddress       = "invalid address string : %s"
)
