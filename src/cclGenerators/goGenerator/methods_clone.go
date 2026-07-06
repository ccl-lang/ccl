package goGenerator

import (
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

//---------------------------------------------------------

func (c *GoGenerationContext) generateCloneMethods(
	builder *codeBuilder.CodeBuilder,
	model *CCLModel,
) error {
	builder.ExpectMappedVars(
		"model",
		"modelName",
	)

	builder.NewLine().
		LineD("func (m $model) CloneEmpty() $model {").
		Indent().
		WriteLine("if m == nil {").
		Indent().
		WriteLine("return nil").
		Unindent().
		WriteLine("}").
		LineD("return &$modelName{}").
		Unindent().
		WriteLine("}")

	builder.NewLine().
		LineD("func (m $model) DeepClone() $model {").
		Indent().
		WriteLine("if m == nil {").
		Indent().
		WriteLine("return nil").
		Unindent().
		WriteLine("}").
		LineD("clone := &$modelName{}")
	for _, field := range model.Fields {
		if err := c.generateDeepCloneField(builder, field); err != nil {
			return err
		}
	}
	builder.WriteLine("return clone").
		Unindent().
		WriteLine("}")

	if c.GetGlobalAttribute(CurrentLanguage, "AddSerializableInterface") != nil {
		builder.NewLine().
			LineD("func (m $model) CloneEmptySerializable() Serializable {").
			Indent().
			WriteLine("if m == nil {").
			Indent().
			WriteLine("return nil").
			Unindent().
			WriteLine("}").
			LineD("return &$modelName{}").
			Unindent().
			WriteLine("}")
	}
	return nil
}

func (c *GoGenerationContext) generateDeepCloneField(
	builder *codeBuilder.CodeBuilder,
	field *CCLField,
) error {
	fieldName := field.Name
	builder.MapVarPairs(
		"field", fieldName,
	)
	defer builder.UnmapVar("field")

	if field.IsArray() {
		targetType := field.Type.GetUnderlyingType()
		fieldType, err := c.getGoTypeForField(field)
		if err != nil {
			return err
		}

		builder.MapVarPairs("fieldType", fieldType)
		defer builder.UnmapVar("fieldType")

		builder.LineD("if m.$field != nil {").
			Indent().
			LineD("clone.$field = make($fieldType, len(m.$field))")
		switch {
		case targetType.IsCustomTypeModel():
			builder.LineD("for i, item := range m.$field {").
				Indent().
				WriteLine("if item != nil {").
				Indent().
				LineD("clone.$field[i] = item.DeepClone()").
				Unindent().
				WriteLine("}").
				Unindent().
				WriteLine("}")
		case targetType.GetName() == cclValues.TypeNameBytes:
			builder.LineD("for i, item := range m.$field {").
				Indent().
				WriteLine("if item != nil {").
				Indent().
				LineD("clone.$field[i] = append([]byte(nil), item...)").
				Unindent().
				WriteLine("}").
				Unindent().
				WriteLine("}")
		default:
			builder.LineD("copy(clone.$field, m.$field)")
		}
		builder.Unindent().
			WriteLine("}")
		return nil
	}

	switch {
	case field.Type.IsCustomTypeModel():
		builder.LineD("if m.$field != nil {").
			Indent().
			LineD("clone.$field = m.$field.DeepClone()").
			Unindent().
			WriteLine("}")
	case field.Type.GetName() == cclValues.TypeNameBytes:
		builder.LineD("if m.$field != nil {").
			Indent().
			LineD("clone.$field = append([]byte(nil), m.$field...)").
			Unindent().
			WriteLine("}")
	default:
		builder.LineD("clone.$field = m.$field")
	}
	return nil
}

//---------------------------------------------------------
