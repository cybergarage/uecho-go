// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protocol

import (
	"fmt"
	"io"
)

// parseFormat1PropertyReader parses the specified property reader.
func (msg *Message) parseFormat1PropertyReader(reader io.Reader) error {
	propertyHeader := make([]byte, Format1PropertyHeaderSize)

	for n := 0; n < int(msg.opc); n++ {
		prop := msg.PropertyAt(n)
		if prop == nil {
			continue
		}

		nRead, err := reader.Read(propertyHeader)
		if err != nil {
			return err
		}
		if nRead < Format1PropertyHeaderSize {
			return fmt.Errorf(errorShortMessageSize, n, Format1PropertyHeaderSize)
		}

		prop.code = PropertyCode(propertyHeader[0])

		propDataSize := int(propertyHeader[1])
		prop.data = make([]byte, propDataSize)
		nRead, err = reader.Read(prop.data)
		if err != nil {
			return err
		}
		if nRead < propDataSize {
			return fmt.Errorf(errorShortMessageSize, n, propDataSize)
		}
	}

	return nil
}

// ParseReader parses the specified bytes.
func (msg *Message) ParseReader(reader io.Reader) error {
	// Frame header

	frameHeader := make([]byte, FrameHeaderSize)
	n, err := reader.Read(frameHeader)
	if err != nil {
		return err
	}
	if n < FrameHeaderSize {
		return fmt.Errorf(errorShortMessageSize, n, FrameHeaderSize)
	}
	err = msg.parseFrameHeaderBytes(frameHeader)
	if err != nil {
		return err
	}

	// Echonet Format1 Header

	format1Header := make([]byte, Format1HeaderSize)
	n, err = reader.Read(format1Header)
	if err != nil {
		return err
	}
	if n < Format1HeaderSize {
		return fmt.Errorf(errorShortMessageSize, n, Format1HeaderSize)
	}
	err = msg.parseFormat1HeaderBytes(format1Header)
	if err != nil {
		return err
	}

	// Propety data

	err = msg.parseFormat1PropertyReader(reader)
	if err != nil {
		return err
	}

	return nil
}
