// Package main is the entry point for the NBP Exchange Rate CLI tool.
package main

import "github.com/MarcusSlover/nbp-cli/internal/cli"

// main is the entry point of the application.
// It delegates the execution to the cli package.
func main() {
	cli.Run()
}
