package cclCmd

import "fmt"

func HandleHelpCommand() {
	fmt.Println("ccl is a tool for managing ccl source code.\n\n" +
		"Usage:\n" +
		"\t ccl <command> [arguments]\n\n" +
		"The commands are:\n\n" +
		"\tgenerate    Generate code from a ccl source file\n" +
		"\tinfo        Show debug info for a generated file\n" +
		"\thelp        Show this help message",
	)
}
