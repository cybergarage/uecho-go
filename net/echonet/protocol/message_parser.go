// Copyright (C) 2018 Satoshi Konno. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"fmt"
)

// parseFrameHeaderBytes parses the specified frame header bytes.
func (msg *Message) parseFrameHeaderBytes(data []byte) error {
	if headerSize := len(data); headerSize < FrameHeaderSize {
		return fmt.Errorf(errorShortMessageSize, headerSize, FrameHeaderSize)
	}

	// Check Headers

	if data[0] != EHD1Echonet {
		return fmt.Errorf(errorInvalidMessageHeader, 0, data[0], EHD1Echonet)
	}

	if data[1] != EHD2Format1 {
		return fmt.Errorf(errorInvalidMessageHeader, 1, data[1], EHD2Format1)
	}

	// TID

	msg.tid[0] = data[2]
	msg.tid[1] = data[3]

	return nil
}

// parseFormat1HeaderBytes parses the specified header bytes.
func (msg *Message) parseFormat1HeaderBytes(data []byte) error {
	if headerSize := len(data); headerSize < Format1HeaderSize {
		return fmt.Errorf(errorShortMessageSize, (headerSize + FrameHeaderSize), (Format1HeaderSize + FrameHeaderSize))
	}

	// SEOJ

	msg.seoj[0] = data[0]
	msg.seoj[1] = data[1]
	msg.seoj[2] = data[2]

	// DEOJ

	msg.deoj[0] = data[3]
	msg.deoj[1] = data[4]
	msg.deoj[2] = data[5]

	// ESV

	msg.esv = ESV(data[6])

	// OPC

	err := msg.SetOPC(int(data[7]))
	if err != nil {
		return err
	}

	return nil
}

// parseFormat1PropertyBytes parses the specified property bytes.
func (msg *Message) parseFormat1PropertyBytes(data []byte) error {
	dataSize := len(data)

	offset := 0
	for n := 0; n < int(msg.opc); n++ {
		prop := msg.PropertyAt(n)
		if prop == nil {
			continue
		}

		// EPC

		if (dataSize - 1) < offset {
			continue
		}

		prop.code = PropertyCode(data[offset])
		offset++

		// PDC

		if (dataSize - 1) < offset {
			continue
		}

		propDataSize := int(data[offset])
		offset++

		// EDT

		if (dataSize - 1) < (offset + propDataSize - 1) {
			continue
		}

		prop.Data = data[offset:(offset + propDataSize)]

		offset += propDataSize
	}

	return nil
}

// ParseBytes parses the specified bytes.
func (msg *Message) ParseBytes(data []byte) error {
	// Frame header

	err := msg.parseFrameHeaderBytes(data)
	if err != nil {
		return err
	}

	// Echonet Format1 Header

	err = msg.parseFormat1HeaderBytes(data[FrameHeaderSize:])
	if err != nil {
		return err
	}

	// Propety data

	err = msg.parseFormat1PropertyBytes(data[Format1MinSize:])
	if err != nil {
		return err
	}

	return nil
}
