// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package echonet provides core APIs for building ECHONET Lite nodes
// (device implementations) and controllers (discovery / management agents)
// in Go. It covers creation and management of local nodes and devices,
// standardized object / property metadata lookup, message composition
// and parsing (ESV frames), transport handling (UDP multicast/unicast,
// optional TCP), and higher-level controller operations.
//
// Overview
//
//   - Controller
//     A Controller simplifies discovery (Search), enumerating nodes,
//     and posting request messages (PostMessage) to obtain responses.
//   - Nodes / Devices
//     A LocalNode represents the local ECHONET Lite node hosting one or more
//     Device (object) instances. Each device has a class group code,
//     class code, an instance code (together forming EOJ) and a set of
//     Properties.
//   - Standard Database
//     SharedStandardDatabase() returns a singleton with:
//   - Manufacturer codes table
//   - Standard (MRA) object and property definitions
//     The data is generated from machine-readable ECHONET resources
//     by scripts under net/echonet/std.
//
// License
//
//   - Provided under a BSD-style license (see LICENSE).
package echonet
