// uechotui is the compatibility entry point for uechoctl tui.
package main

import (
	"fmt"
	"github.com/cybergarage/uecho-go/internal/tui"
	"os"
)

func main() {
	command := tui.NewCommand()
	command.Use = "uechotui"
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
