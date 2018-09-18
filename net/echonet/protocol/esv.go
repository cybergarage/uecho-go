// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
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

// IsValidESV returns true whether the specified code is valid, otherwise false.
func IsValidESV(esv ESV) bool {
	validCodes := []ESV{
		ESVWriteRequest,
		ESVWriteRequestResponseRequired,
		ESVReadRequest,
		ESVNotificationRequest,
		ESVWriteReadRequest,
		ESVWriteResponse,
		ESVReadResponse,
		ESVNotification,
		ESVNotificationResponseRequired,
		ESVNotificationResponse,
		ESVWriteReadResponse,
		ESVWriteRequestError,
		ESVWriteRequestResponseRequiredError,
		ESVReadRequestError,
		ESVNotificationRequestError,
		ESVWriteReadRequestError,
	}

	for _, code := range validCodes {
		if esv == code {
			return true
		}
	}

	return false
}

// IsWriteRequest returns true whether the specified code is a write request type, otherwise false.
func IsWriteRequest(esv ESV) bool {
	switch esv {
	case ESVWriteRequest:
		return true
	case ESVWriteRequestResponseRequired:
		return true
	case ESVWriteReadRequest:
		return true
	}
	return false
}

// IsReadRequest returns true whether the specified code is a read request type, otherwise false.
func IsReadRequest(esv ESV) bool {
	switch esv {
	case ESVReadRequest:
		return true
	case ESVWriteReadRequest:
		return true
	}
	return false
}

// IsNotificationRequest returns true whether the specified code is a notification request type, otherwise false.
func IsNotificationRequest(esv ESV) bool {
	switch esv {
	case ESVNotificationRequest:
		return true
	}
	return false
}

// IsWriteResponse returns true whether the specified code is a write response type, otherwise false.
func IsWriteResponse(esv ESV) bool {
	switch esv {
	case ESVWriteResponse:
		return true
	case ESVWriteReadResponse:
		return true
	}
	return false
}

// IsReadResponse returns true whether the specified code is a read response type, otherwise false.
func IsReadResponse(esv ESV) bool {
	switch esv {
	case ESVReadResponse:
		return true
	case ESVWriteReadResponse:
		return true
	}
	return false
}

// IsNotification returns true whether the specified code is a notification type, otherwise false.
func IsNotification(esv ESV) bool {
	switch esv {
	case ESVNotification:
		return true
	case ESVNotificationResponseRequired:
		return true
	}
	return false
}

// IsNotificationResponse returns true whether the specified code is a notification response type, otherwise false.
func IsNotificationResponse(esv ESV) bool {
	switch esv {
	case ESVNotificationResponse:
		return true
	}
	return false
}

// IsResponseRequired returns true whether the ESV requires the response, otherwise false.
func IsResponseRequired(esv ESV) bool {
	switch esv {
	case ESVReadRequest:
		return true
	case ESVWriteReadRequest:
		return true
	case ESVWriteRequestResponseRequired:
		return true
	case ESVNotificationResponseRequired:
		return true
	}
	return false
}
