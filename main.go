package main

import (
	"fmt"
	"os"

	"github.com/TheCoolRobot/asana-cli/cmd"
)

var (
	Version = "0.1.0"
	Commit  = "unknown"
	Date    = "Feb 17, 2026"
)

func init() {
	// Version info is injected at build time
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}