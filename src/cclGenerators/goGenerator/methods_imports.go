package goGenerator

import (
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func (c *GoGenerationContext) registerGoMethodImports() {
	imports := map[string]bool{}

	for _, currentTypeDef := range c.GetGenerationTypeDefinitions() {
		if !currentTypeDef.IsCustomModel() {
			continue
		}

		model := currentTypeDef.GetModelDefinition()
		needsBinary := c.NeedsBinarySerialization(CurrentLanguage, model)
		needsJson := c.NeedsJsonSerialization(CurrentLanguage, model)

		if needsBinary {
			imports["bytes"] = true
			imports["encoding/binary"] = true
		}
		if needsJson {
			imports["bytes"] = true
			if len(model.Fields) > 0 {
				imports["strconv"] = true
				imports["strings"] = true
			}
		}

		for _, field := range model.Fields {
			targetType := field.Type
			if targetType.IsArray() {
				targetType = targetType.GetUnderlyingType()
			}

			switch targetType.GetName() {
			case cclValues.TypeNameDateTime:
				if needsBinary || needsJson {
					imports["time"] = true
				}
			case cclValues.TypeNameBytes:
				if needsJson {
					imports["encoding/base64"] = true
				}
			}
		}
	}

	registerGoImports(c.MethodsCode, imports)
}
