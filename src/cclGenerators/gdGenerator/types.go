package gdGenerator

import (
	"strings"

	gen "github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

// CCLModel is a type alias for the model definition type from the CCL library.
type CCLModel = cclValues.ModelDefinition

// CCLField is a type alias for the field definition type from the CCL library.
type CCLField = cclValues.FieldDefinition

type GDScriptGenerationContext struct {
	Options *gen.CodeGenerationOptions
	// One builder per model class file
	ModelClasses map[string]*strings.Builder
}
