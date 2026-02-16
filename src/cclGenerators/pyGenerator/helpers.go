package pyGenerator

import (
	"os"
	"strings"

	"github.com/ALiwoto/ssg/ssg"
	gen "github.com/ccl-lang/ccl/src/cclGenerators"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

// GenerateCode generates Python code from the provided CCL source file.
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

	pyCtx := &PythonGenerationContext{
		CodeGenerationBase: &gen.CodeGenerationBase{
			Options: options,
		},
	}
	err := pyCtx.GenerateCode()
	if err != nil {
		return nil, err
	}

	return &gen.CodeGenerationResult{
		SourceLanguage: gValues.LanguageCCL.String(),
		TargetLanguage: CurrentLanguage.String(),
		OutputFiles:    pyCtx.OutputFiles,
	}, nil
}

func toSnakeCase(s string) string {
	var b strings.Builder
	for i, v := range s {
		if v >= 'A' && v <= 'Z' {
			if i > 0 && s[i-1] != '_' && (s[i-1] < 'A' || s[i-1] > 'Z') {
				b.WriteRune('_')
			}
			b.WriteRune(v + ('a' - 'A'))
		} else {
			b.WriteRune(v)
		}
	}
	return b.String()
}
