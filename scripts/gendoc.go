// Copyright (C) 2018 The uecho-go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//go:build ignore

package main

import (
	"log"

	"github.com/cybergarage/uecho-go/net/echonet/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	rootCmd := cmd.GetRootCommand()

	log.Println("Generating CLI documentation...")

	// Generate Markdown documentation
	err := doc.GenMarkdownTree(rootCmd, "./doc/cmd")
	if err != nil {
		log.Fatal("Failed to generate markdown docs:", err)
	}

	log.Println("Documentation generated successfully!")
}
