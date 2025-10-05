// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package echonet

import (
	"github.com/cybergarage/uecho-go/net/echonet/protocol"
)

type ESV = protocol.ESV

const (
	ESVWriteRequest                      = protocol.ESVWriteRequest
	ESVWriteRequestResponseRequired      = protocol.ESVWriteRequestResponseRequired
	ESVReadRequest                       = protocol.ESVReadRequest
	ESVNotificationRequest               = protocol.ESVNotificationRequest
	ESVWriteReadRequest                  = protocol.ESVWriteReadRequest
	ESVWriteResponse                     = protocol.ESVWriteResponse
	ESVReadResponse                      = protocol.ESVReadResponse
	ESVNotification                      = protocol.ESVNotification
	ESVNotificationResponseRequired      = protocol.ESVNotificationResponseRequired
	ESVNotificationResponse              = protocol.ESVNotificationResponse
	ESVWriteReadResponse                 = protocol.ESVWriteReadResponse
	ESVWriteRequestError                 = protocol.ESVWriteRequestError
	ESVWriteRequestResponseRequiredError = protocol.ESVWriteRequestResponseRequiredError
	ESVReadRequestError                  = protocol.ESVReadRequestError
	ESVNotificationRequestError          = protocol.ESVNotificationRequestError
	ESVWriteReadRequestError             = protocol.ESVWriteReadRequestError
)
