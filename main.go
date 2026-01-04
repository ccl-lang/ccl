package main

import (
	"fmt"
	"os"

	"github.com/ccl-lang/ccl/src/cclCmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No command provided. Use 'generate' to generate code.")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "generate":
		cclCmd.HandleGenerateCommand()
	case "info":
		cclCmd.HandleInfoCommand()
	case "version", "--version":
		cclCmd.HandleVersionCommand()
	case "help":
		cclCmd.HandleHelpCommand()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Use 'ccl help' for usage.")
		os.Exit(1)
	}
}
