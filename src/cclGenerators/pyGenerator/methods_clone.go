package pyGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *PythonGenerationContext) generateCloneMethods(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.ExpectMappedVars("model")

	builder.LineD(`def clone_empty(self) -> "$model":`).
		Indent().
		LineD("return $model()").
		UnindentLine().
		NewLine()

	return c.generateDeepCloneMethod(model, builder)
}

func (c *PythonGenerationContext) generateDeepCloneMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.ExpectMappedVars("model")

	builder.LineD(`def deep_clone(self) -> "$model":`).
		Indent().
		LineD("clone = $model()")
	for _, field := range model.Fields {
		c.generateDeepCloneField(field, builder)
	}
	builder.WriteLine("return clone").
		UnindentLine().
		NewLine()
	return nil
}

func (c *PythonGenerationContext) generateDeepCloneField(
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
			builder.LineD("clone.$field = [item.deep_clone() if item is not None else None for item in self.$field]")
		case targetType.GetName() == cclValues.TypeNameBytes:
			builder.LineD("clone.$field = [bytes(item) for item in self.$field]")
		default:
			builder.LineD("clone.$field = list(self.$field)")
		}
		return
	}

	switch {
	case field.Type.IsCustomTypeModel():
		builder.LineD("clone.$field = self.$field.deep_clone() if self.$field is not None else None")
	case field.Type.GetName() == cclValues.TypeNameBytes:
		builder.LineD("clone.$field = bytes(self.$field)")
	default:
		builder.LineD("clone.$field = self.$field")
	}
}
