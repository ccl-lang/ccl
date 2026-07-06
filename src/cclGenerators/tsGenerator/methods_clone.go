package tsGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *TypeScriptGenerationContext) generateDeepCloneMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.ExpectMappedVars("model")

	builder.LineD("public deepClone(): $model {").
		Indent().
		LineD("const clone = new $model();")
	for _, field := range model.Fields {
		c.generateDeepCloneField(field, builder)
	}
	builder.WriteLine("return clone;").
		Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *TypeScriptGenerationContext) generateDeepCloneField(
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) {
	fieldName := caseUtils.ToCamelCase(field.Name)
	builder.MapVarPairs("field", fieldName)
	defer builder.UnmapVar("field")

	if field.IsArray() {
		targetType := field.Type.GetUnderlyingType()
		switch {
		case targetType.IsCustomTypeModel():
			builder.LineD("clone.$field = this.$field.map(item => item ? item.deepClone() : null);")
		case targetType.GetName() == cclValues.TypeNameBytes:
			builder.LineD("clone.$field = this.$field.map(item => new Uint8Array(item));")
		default:
			builder.LineD("clone.$field = [...this.$field];")
		}
		return
	}

	switch {
	case field.Type.IsCustomTypeModel():
		builder.LineD("clone.$field = this.$field ? this.$field.deepClone() : null as any;")
	case field.Type.GetName() == cclValues.TypeNameBytes:
		builder.LineD("clone.$field = new Uint8Array(this.$field);")
	default:
		builder.LineD("clone.$field = this.$field;")
	}
}
