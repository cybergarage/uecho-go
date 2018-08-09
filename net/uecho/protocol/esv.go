// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

type ESV byte

const (
	ESVWriteRequest                      = 0x60
	ESVWriteRequestResponseRequired      = 0x61
	ESVReadRequest                       = 0x62
	ESVNotificationRequest               = 0x63
	ESVWriteReadRequest                  = 0x6E
	ESVWriteResponse                     = 0x71
	ESVReadResponse                      = 0x72
	ESVNotification                      = 0x73
	ESVNotificationResponseRequired      = 0x74
	ESVNotificationResponse              = 0x7A
	ESVWriteReadResponse                 = 0x7E
	ESVWriteRequestError                 = 0x50
	ESVWriteRequestResponseRequiredError = 0x51
	ESVReadRequestError                  = 0x52
	ESVNotificationRequestError          = 0x53
	ESVWriteReadRequestError             = 0x5E
)

// IsWriteRequest returns true whether the message is a write request type, otherwise false.
func IsWriteRequest(esv ESV) bool {
	switch esv {
	case ESVWriteRequest:
		return true
	case ESVWriteReadRequest:
		return true
	}
	return false
}

// IsReadRequest returns true whether the message is a read request type, otherwise false.
func IsReadRequest(esv ESV) bool {
	switch esv {
	case ESVReadResponse:
		return true
	case ESVWriteReadRequest:
		return true
	}
	return false
}

// IsNotificationRequest returns true whether the message is a notification request type, otherwise false.
func IsNotificationRequest(esv ESV) bool {
	switch esv {
	case ESVNotificationRequest:
		return true
	case ESVNotification:
		return true
	}
	return false
}

// IsResponseRequired returns true whether the ESV requires the response, otherwise false.
func IsResponseRequired(esv ESV) bool {
	switch esv {
	case ESVReadRequest:
		return true
	case ESVWriteRequestResponseRequired:
		return true
	case ESVNotificationResponseRequired:
		return true
	case ESVReadResponse:
		return true
	case ESVWriteReadResponse:
		return true
	}
	return false
}
