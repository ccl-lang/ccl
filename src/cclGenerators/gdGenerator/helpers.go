package gdGenerator

import (
	"os"
	"strings"

	gen "github.com/ALiwoto/ccl/src/cclGenerators"
	"github.com/ALiwoto/ssg/ssg"
)

// GenerateCode generates GDScript code from the provided CCL source file.
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

	goCtx := &GDScriptGenerationContext{
		Options: options,
	}
	err := goCtx.GenerateCode()
	if err != nil {
		return nil, err
	}

	return &gen.CodeGenerationResult{}, nil
}

func snakeToTitle(s string) string {
	bd := strings.Builder{}

	for _, split := range strings.Split(s, "_") {
		bd.WriteString(strings.Title(split))
	}

	return bd.String()
}

func ToCamelCase(s string) string {
	title := snakeToTitle(s)

	return strings.ToLower(title[:1]) + title[1:]
}

func ToPascalCase(str string) string {
	title := snakeToTitle(str)

	return strings.ToUpper(title[:1]) + title[1:]
}

func ToSnakeCase(camel string) (snake string) {
	var b strings.Builder
	diff := 'a' - 'A'
	l := len(camel)
	for i, v := range camel {
		// A is 65, a is 97
		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		// v is capital letter here
		// irregard first letter
		// add underscore if last letter is capital letter
		// add underscore when previous letter is lowercase
		// add underscore when next letter is lowercase
		if (i != 0 || i == l-1) && (          // head and tail
		(i > 0 && rune(camel[i-1]) >= 'a') || // pre
			(i < l-1 && rune(camel[i+1]) >= 'a')) { //next
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}
