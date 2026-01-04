package cclCmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
)

func HandleInfoCommand() {
	fs := flag.NewFlagSet("info", flag.ExitOnError)
	file := fs.String("file", "", "Path to the generated file")
	line := fs.Int("line", 0, "Line number in the generated file")

	fs.Parse(os.Args[2:])

	if *file == "" || *line == 0 {
		fmt.Println("Error: --file and --line flags are required")
		fs.Usage()
		os.Exit(1)
	}

	debugInfoPath := *file + ".cclinfo"
	data, err := os.ReadFile(debugInfoPath)
	if err != nil {
		fmt.Printf("Error: failed to read debug info file %s: %v\n", debugInfoPath, err)
		os.Exit(1)
	}

	var debugInfos []*codeBuilder.DebugInfo
	err = json.Unmarshal(data, &debugInfos)
	if err != nil {
		fmt.Printf("Error: failed to parse debug info file: %v\n", err)
		os.Exit(1)
	}

	found := false
	for _, info := range debugInfos {
		if info.GeneratedLine == *line {
			fmt.Printf("%s:%d\n", info.SourceFile, info.SourceLine)
			found = true
		}
	}

	if !found {
		var bestMatch *codeBuilder.DebugInfo
		for _, info := range debugInfos {
			if info.GeneratedLine <= *line {
				if bestMatch == nil || info.GeneratedLine > bestMatch.GeneratedLine {
					bestMatch = info
				}
			}
		}
		if bestMatch != nil {
			fmt.Printf("%s:%d\n", bestMatch.SourceFile, bestMatch.SourceLine)
			found = true
		}
	}

	if !found {
		fmt.Println("No debug info found for this line.")
	}
}
