package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/cclParser"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No command provided. Use 'generate' to generate code.")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "generate":
		// Define flags
		fs := flag.NewFlagSet("generate", flag.ExitOnError)
		source := fs.String("source", "", "Path to the source file")
		language := fs.String("language", "", "Programming language for the generated code")
		output := fs.String("output", "", "Output path for the generated package")
		generateDebugInfo := fs.Bool("generate-debug-info", false, "Generate debug info file")
		o := fs.String("o", "", "Output path for the generated package")
		if o != nil && *o != "" {
			output = o
		}

		// Parse flags
		fs.Parse(os.Args[2:])

		// Validate flags
		if *source == "" || *language == "" || *output == "" {
			fmt.Println("Error: --source, --language, and --output flags are required when using 'generate'")
			fs.Usage()
			os.Exit(1)
		}

		// Handle generate command
		fmt.Printf("Generating code from %s in %s language to %s\n", *source, *language, *output)

		// check if the source file exists
		if _, err := os.Stat(*source); err != nil {
			switch {
			case os.IsNotExist(err):
				fmt.Printf("Error: source file %s does not exist\n", *source)
			case os.IsPermission(err):
				fmt.Printf("Error: permission denied for source file %s\n", *source)
			default:
				fmt.Printf("Error: failed to check source file %s: %v\n", *source, err)
			}
			os.Exit(1)
		}

		parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
			SourceFilePath: *source,
		})
		if parseErr != nil {
			fmt.Printf("Error: failed to parse source file %s: %v\n", *source, parseErr)
			os.Exit(1)
		}

		cclLoader.LoadGenerators()
		result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
			CCLDefinition:     parsedDefinitions,
			OutputPath:        *output,
			TargetLanguage:    *language,
			GenerateDebugInfo: *generateDebugInfo,
		})
		if err != nil {
			fmt.Printf("Error: failed to generate code: %v\n", err)
			os.Exit(1)
		} else if result == nil {
			fmt.Println("Unknown error: failed to generate code")
			os.Exit(1)
		}

		fmt.Println("\nCode generation completed successfully")
	case "info":
		fs := flag.NewFlagSet("info", flag.ExitOnError)
		file := fs.String("file", "", "Path to the generated file")
		line := fs.Int("line", 0, "Line number in the generated file")
		// col := fs.Int("col", 0, "Column number in the generated file") // Not used yet as per user request for line mapping mainly

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

	case "version":
		fmt.Printf("ccl version %s %s/%s\n", gValues.CurrentCCLVersion, runtime.GOOS, runtime.GOARCH)
	case "help":
		fmt.Println("ccl is a tool for managing ccl source code.\n\n" +
			"Usage:\n" +
			"\t ccl <command> [arguments]\n\n" +
			"The commands are:\n\n" +
			"\tgenerate    Generate code from a ccl source file\n" +
			"\tinfo        Show debug info for a generated file\n" +
			"\thelp        Show this help message",
		)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Use 'ccl help' for usage.")
		os.Exit(1)
	}
}
