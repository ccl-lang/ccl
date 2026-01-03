package tsGenerator

import (
	gen "github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

// CCLModel is a type alias for the model definition type from the CCL's
// internal packages.
type CCLModel = cclValues.ModelDefinition

// CCLField is a type alias for the field definition type from the CCL's
// internal packages.
type CCLField = cclValues.ModelFieldDefinition

type TypeScriptGenerationContext struct {
	*gen.CodeGenerationBase

	// One builder per model class file (or single builder for single file mode)
	ModelClasses map[string]*codeBuilder.CodeBuilder

	// SingleFileBuilder is the builder used in single file mode
	SingleFileBuilder *codeBuilder.CodeBuilder

	// IsSingleFile indicates if we are generating a single file
	IsSingleFile bool

	// SingleFileName is the name of the single file to generate
	SingleFileName string
}
