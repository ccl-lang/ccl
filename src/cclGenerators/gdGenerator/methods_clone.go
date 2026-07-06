package gdGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *GDScriptGenerationContext) generateCloneMethods(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.ExpectMappedVars("model")

	builder.WriteLine("func clone_empty() -> " + model.Name + ":").
		Indent().
		WriteLine("return " + model.Name + ".new()").
		UnindentLine()

	return c.generateDeepCloneMethod(model, builder)
}

func (c *GDScriptGenerationContext) generateDeepCloneMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.ExpectMappedVars("model")

	builder.NewLine().
		LineD("func deep_clone() -> $model:").
		Indent().
		LineD("var clone := $model.new()")
	for _, field := range model.Fields {
		c.generateDeepCloneField(field, builder)
	}
	builder.WriteLine("return clone").
		UnindentLine()
	return nil
}

func (c *GDScriptGenerationContext) generateDeepCloneField(
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) {
	fieldName := caseUtils.ToSnakeCase(field.Name)
	builder.MapVarPairs("field", fieldName)
	defer builder.UnmapVar("field")

	if field.IsArray() {
		targetType := field.Type.GetUnderlyingType()
		switch {
		case targetType.IsCustomTypeModel():
			builder.LineD("clone.$field = []").
				LineD("for item in $field:").
				Indent().
				LineD("clone.$field.append(item.deep_clone() if item != null else null)").
				Unindent()
		case targetType.GetName() == cclValues.TypeNameBytes:
			builder.LineD("clone.$field = []").
				LineD("for item in $field:").
				Indent().
				LineD("clone.$field.append(item.duplicate())").
				Unindent()
		default:
			builder.LineD("clone.$field = $field.duplicate()")
		}
		return
	}

	switch {
	case field.Type.IsCustomTypeModel():
		builder.LineD("clone.$field = $field.deep_clone() if $field != null else null")
	case field.Type.GetName() == cclValues.TypeNameBytes:
		builder.LineD("clone.$field = $field.duplicate()")
	default:
		builder.LineD("clone.$field = $field")
	}
}
