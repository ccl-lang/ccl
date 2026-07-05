package rsGenerator

import (
	"os"

	"github.com/ALiwoto/ssg/ssg"
	gen "github.com/ccl-lang/ccl/src/cclGenerators"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

func GenerateCode(options *gen.CodeGenerationOptions) (*gen.CodeGenerationResult, error) {
	if _, err := os.Stat(options.OutputPath); os.IsNotExist(err) {
		if err = os.MkdirAll(options.OutputPath, os.ModePerm); err != nil {
			return nil, err
		}
	}
	if options.PackageName == "" {
		pathParts := ssg.Split(options.OutputPath, "/", "\\")
		options.PackageName = pathParts[len(pathParts)-1]
	}

	context := &RustGenerationContext{
		CodeGenerationBase: &gen.CodeGenerationBase{
			Options: options,
		},
	}

	outputFiles, err := context.GenerateCode()
	if err != nil {
		return nil, err
	}

	return &gen.CodeGenerationResult{
		SourceLanguage: gValues.LanguageCCL.String(),
		TargetLanguage: CurrentLanguage.String(),
		OutputFiles:    outputFiles,
	}, nil
}
