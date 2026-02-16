package gdGenerator

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

type GDScriptGenerationContext struct {
	*gen.CodeGenerationBase

	// One builder per model class file
	ModelClasses  map[string]*codeBuilder.CodeBuilder
	ModelSections map[string][]string

	// OutputFiles is the list of generated files.
	OutputFiles []string
}
