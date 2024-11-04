package goGenerator

import (
	"os"

	gen "github.com/ALiwoto/ccl/src/cclGenerators"
	"github.com/ALiwoto/ssg/ssg"
)

// GenerateCode generates Go code from the provided CCL source file.
func GenerateCode(options *gen.CodeGenerationOptions) (*gen.CodeGenerationResult, error) {
	// If the output directory does not exist, create it recursively.
	if _, err := os.Stat(options.OutputPath); os.IsNotExist(err) {
		err := os.MkdirAll(options.OutputPath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	if options.PackageName == "" {
		// if there is no package name provided, use the last folder name
		// in the output path as the package name.
		pathParts := ssg.Split(options.OutputPath, "/", "\\")
		options.PackageName = pathParts[len(pathParts)-1]
	}

	goCtx := &GoGenerationContext{
		Options: options,
	}
	err := goCtx.GenerateCode()
	if err != nil {
		return nil, err
	}

	return &gen.CodeGenerationResult{}, nil
}
