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

// IsValidESV returns true whether the specified code is valid, otherwise false.
func IsValidESV(esv ESV) bool {
	return protocol.IsValidESV(esv)
}

// IsWriteRequest returns true whether the specified code is a write request type, otherwise false.
func IsWriteRequest(esv ESV) bool {
	return protocol.IsWriteRequest(esv)
}

// IsReadRequest returns true whether the specified code is a read request type, otherwise false.
func IsReadRequest(esv ESV) bool {
	return protocol.IsReadRequest(esv)
}

// IsNotificationRequest returns true whether the specified code is a notification request type, otherwise false.
func IsNotificationRequest(esv ESV) bool {
	return protocol.IsNotificationRequest(esv)
}

// IsWriteResponse returns true whether the specified code is a write response type, otherwise false.
func IsWriteResponse(esv ESV) bool {
	return protocol.IsWriteResponse(esv)
}

// IsReadResponse returns true whether the specified code is a read response type, otherwise false.
func IsReadResponse(esv ESV) bool {
	return protocol.IsReadResponse(esv)
}

// IsNotification returns true whether the specified code is a notification type, otherwise false.
func IsNotification(esv ESV) bool {
	return protocol.IsNotification(esv)
}

// IsNotificationResponse returns true whether the specified code is a notification response type, otherwise false.
func IsNotificationResponse(esv ESV) bool {
	return protocol.IsNotificationResponse(esv)
}

// IsResponseRequired returns true whether the ESV requires the response, otherwise false.
func IsResponseRequired(esv ESV) bool {
	return protocol.IsResponseRequired(esv)
}
