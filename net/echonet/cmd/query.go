// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

type Query struct {
	Details bool
	Address string
}

// NewQuery creates and returns a new Query.
func NewQuery() *Query {
	return &Query{
		Details: false,
		Address: "",
	}
}

func NewQueryWith(details bool, address string) *Query {
	return &Query{
		Details: details,
		Address: address,
	}
}
