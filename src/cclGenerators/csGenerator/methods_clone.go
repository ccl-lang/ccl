package csGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *CSharpGenerationContext) generateCloneMethods(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.ExpectMappedVars("model")

	builder.LineD("public $model CloneEmpty()").
		WriteLine("{").
		Indent().
		LineD("return new $model();").
		Unindent().
		WriteLine("}").
		NewLine()

	return c.generateDeepCloneMethod(model, builder)
}

func (c *CSharpGenerationContext) generateDeepCloneMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.ExpectMappedVars("model")

	builder.LineD("public $model DeepClone()").
		WriteLine("{").
		Indent().
		LineD("var clone = new $model();")
	for _, field := range model.Fields {
		if err := c.generateDeepCloneField(field, builder); err != nil {
			return err
		}
	}
	builder.WriteLine("return clone;").
		Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *CSharpGenerationContext) generateDeepCloneField(
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldName := caseUtils.ToPascalCase(field.Name)
	builder.MapVarPairs("field", fieldName)
	defer builder.UnmapVar("field")

	if field.IsArray() {
		fieldType, err := c.getCSharpType(field)
		if err != nil {
			return err
		}
		builder.MapVarPairs("fieldType", fieldType)
		defer builder.UnmapVar("fieldType")

		targetType := field.Type.GetUnderlyingType()
		switch {
		case targetType.IsCustomTypeModel():
			builder.LineD("clone.$field = new $fieldType();").
				LineD("foreach (var item in $field)").
				WriteLine("{").
				Indent().
				LineD("clone.$field.Add(item == null ? null : item.DeepClone());").
				Unindent().
				WriteLine("}")
		case targetType.GetName() == cclValues.TypeNameBytes:
			builder.LineD("clone.$field = new $fieldType();").
				LineD("foreach (var item in $field)").
				WriteLine("{").
				Indent().
				LineD("clone.$field.Add(item == null ? null : (byte[])item.Clone());").
				Unindent().
				WriteLine("}")
		default:
			builder.LineD("clone.$field = new $fieldType($field);")
		}
		return nil
	}

	switch {
	case field.Type.IsCustomTypeModel():
		builder.LineD("clone.$field = $field == null ? null : $field.DeepClone();")
	case field.Type.GetName() == cclValues.TypeNameBytes:
		builder.LineD("clone.$field = $field == null ? null : (byte[])$field.Clone();")
	default:
		builder.LineD("clone.$field = $field;")
	}
	return nil
}
