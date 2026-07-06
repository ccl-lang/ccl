package rsGenerator

import (
	"sort"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *RustGenerationContext) GenerateCode() ([]string, error) {
	if err := c.generateLibFile(); err != nil {
		return nil, err
	}
	if err := c.generateTypeFiles(); err != nil {
		return nil, err
	}

	c.finalizeCodeBuilders()
	return c.writeCodeBuilders()
}

func (c *RustGenerationContext) generateLibFile() error {
	builder := c.getCodeBuilder(LibFileName, "lib")
	builder.WriteLine("#![allow(unused_assignments)]").
		NewLine()
	declaredModules := map[string]bool{}
	for _, typeDef := range c.GetGenerationTypeDefinitions() {
		if typeDef.IsCustomEnum() {
			enumDef := typeDef.GetEnumDefinition()
			if enumDef.IsNested() {
				continue
			}
			c.writeRustModuleExport(builder, enumDef.Name, enumDef.Name, declaredModules)
			continue
		}

		if !typeDef.IsCustomModel() {
			return &cclErrors.UnsupportedTypeDefinitionError{
				TypeName:       typeDef.GetFullName(),
				TargetLanguage: CurrentLanguage.String(),
			}
		}

		model := typeDef.GetModelDefinition()
		c.writeRustModuleExport(builder, model.Name, model.Name, declaredModules)
		for _, enumDef := range model.Enums {
			enumTypeName, err := c.getRustEnumTypeName(enumDef)
			if err != nil {
				return err
			}
			c.writeRustModuleExport(builder, model.Name, enumTypeName, declaredModules)
		}
	}
	return nil
}

func (c *RustGenerationContext) writeRustModuleExport(
	builder *codeBuilder.CodeBuilder,
	sourceName string,
	typeName string,
	declaredModules map[string]bool,
) {
	moduleName := rustModuleName(sourceName)
	if !declaredModules[moduleName] {
		builder.WriteLine("pub mod " + moduleName + ";")
		declaredModules[moduleName] = true
	}
	builder.WriteLine("pub use " + moduleName + "::" + typeName + ";")
}

func (c *RustGenerationContext) generateTypeFiles() error {
	for _, typeDef := range c.GetGenerationTypeDefinitions() {
		if typeDef.IsCustomEnum() {
			enumDef := typeDef.GetEnumDefinition()
			if enumDef.IsNested() {
				continue
			}
			builder := c.getCodeBuilder("src/"+rustModuleFileName(enumDef.Name), enumDef.GetFullName())
			if err := c.generateEnum(builder, enumDef); err != nil {
				return err
			}
			continue
		}

		if !typeDef.IsCustomModel() {
			return &cclErrors.UnsupportedTypeDefinitionError{
				TypeName:       typeDef.GetFullName(),
				TargetLanguage: CurrentLanguage.String(),
			}
		}

		model := typeDef.GetModelDefinition()
		builder := c.getCodeBuilder("src/"+rustModuleFileName(model.Name), model.GetFullName())
		if err := c.generateModel(builder, model); err != nil {
			return err
		}
	}
	return nil
}

func (c *RustGenerationContext) generateModel(builder *codeBuilder.CodeBuilder, model *CCLModel) error {
	if err := c.writeModelImports(builder, model); err != nil {
		return err
	}
	for _, enumDef := range model.Enums {
		if err := c.generateEnum(builder, enumDef); err != nil {
			return err
		}
		builder.NewLine()
	}

	deriveLine := "#[derive(Debug, PartialEq, serde::Serialize, serde::Deserialize)]"
	if c.NeedsCloneMethods(CurrentLanguage, model) {
		deriveLine = "#[derive(Debug, Clone, PartialEq, serde::Serialize, serde::Deserialize)]"
	}
	builder.WriteLine(deriveLine).
		WriteLine("#[serde(default)]").
		WriteLine("pub struct " + model.Name + " {").
		Indent()
	for _, field := range model.Fields {
		if err := c.generateRustField(builder, model, field); err != nil {
			return err
		}
	}
	builder.Unindent().
		WriteLine("}").
		NewLine()

	if err := c.generateModelDefault(builder, model); err != nil {
		return err
	}
	if c.NeedsCloneMethods(CurrentLanguage, model) {
		if err := c.generateCloneMethods(builder, model); err != nil {
			return err
		}
	}
	if c.NeedsJsonSerialization(CurrentLanguage, model) {
		c.generateModelJsonMethods(builder, model)
		c.generateModelJsonAdapters(builder, model)
	}
	if c.NeedsBinarySerialization(CurrentLanguage, model) {
		return c.generateModelBinaryMethods(builder, model)
	}
	return nil
}

func (c *RustGenerationContext) generateRustField(
	builder *codeBuilder.CodeBuilder,
	model *CCLModel,
	field *CCLField,
) error {
	fieldType, err := c.getRustTypeForField(field)
	if err != nil {
		return err
	}
	jsonName, err := c.GetJsonFieldName(CurrentLanguage, model, field)
	if err != nil {
		return err
	}

	fieldName := rustFieldName(field.Name)
	if jsonName != fieldName {
		builder.WriteLine(`#[serde(rename = "` + jsonName + `")]`)
	}
	if isRustBytesType(field.Type) {
		builder.WriteLine(`#[serde(with = "ccl_bytes_json")]`)
	}
	if isRustBytesArrayType(field.Type) {
		builder.WriteLine(`#[serde(with = "ccl_bytes_array_json")]`)
	}
	builder.WriteLine("pub " + fieldName + ": " + fieldType + ",")
	return nil
}

func (c *RustGenerationContext) writeModelImports(builder *codeBuilder.CodeBuilder, model *CCLModel) error {
	imports := map[string]string{}
	for _, field := range model.Fields {
		targetType := field.Type
		if targetType.IsArray() {
			targetType = targetType.GetUnderlyingType()
		}

		if targetType.IsCustomTypeModel() {
			targetModel := targetType.GetDefinition().GetModelDefinition()
			if targetModel != nil && targetModel != model {
				imports[targetModel.Name] = "use crate::" + targetModel.Name + ";"
			}
		} else if targetType.IsCustomTypeEnum() {
			enumDef := targetType.GetDefinition().GetEnumDefinition()
			enumTypeName, err := c.getRustEnumTypeName(enumDef)
			if err != nil {
				return err
			}
			if enumDef.IsNested() && enumDef.OwnedBy == model {
				continue
			}
			imports[enumTypeName] = "use crate::" + enumTypeName + ";"
		}
	}

	keys := make([]string, 0, len(imports))
	for key := range imports {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		builder.WriteLine(imports[key])
	}
	if len(keys) != 0 {
		builder.NewLine()
	}
	return nil
}

func (c *RustGenerationContext) getRustTypeForField(field *CCLField) (string, error) {
	targetType := field.Type
	if targetType.IsArray() {
		targetType = targetType.GetUnderlyingType()
	}

	rustType, err := c.getRustTypeForUsage(targetType, field.OwnedBy)
	if err != nil {
		return "", err
	}

	if field.IsArray() {
		return "Vec<" + rustType + ">", nil
	}
	return rustType, nil
}

func (c *RustGenerationContext) getRustTypeForUsage(
	targetType *cclValues.CCLTypeUsage,
	currentModel *CCLModel,
) (string, error) {
	if mappedType, ok := CCLTypesToRustTypes[targetType.GetName()]; ok {
		return mappedType, nil
	}
	if targetType.IsCustomTypeEnum() {
		return c.getRustEnumTypeReference(targetType.GetDefinition().GetEnumDefinition(), currentModel)
	}
	if targetType.IsCustomTypeModel() {
		return "Option<Box<" + targetType.GetName() + ">>", nil
	}
	return "", &cclErrors.UnsupportedFieldTypeError{
		TypeName:       targetType.GetName(),
		TargetLanguage: CurrentLanguage.String(),
	}
}

func (c *RustGenerationContext) generateModelDefault(builder *codeBuilder.CodeBuilder, model *CCLModel) error {
	builder.WriteLine("impl Default for " + model.Name + " {").
		Indent().
		WriteLine("fn default() -> Self {").
		Indent().
		WriteLine("Self {").
		Indent()
	for _, field := range model.Fields {
		defaultValue, err := c.rustDefaultValue(field)
		if err != nil {
			return err
		}
		builder.WriteLine(rustFieldName(field.Name) + ": " + defaultValue + ",")
	}
	builder.Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *RustGenerationContext) rustDefaultValue(field *CCLField) (string, error) {
	if field.IsArray() {
		return "Vec::new()", nil
	}
	if enumRef := c.GetEnumDefaultReference(field); enumRef != nil {
		enumType, err := c.getRustEnumTypeReference(enumRef.Enum, field.OwnedBy)
		if err != nil {
			return "", err
		}
		memberName, err := c.getRustEnumMemberName(enumRef.Enum, enumRef.Member)
		if err != nil {
			return "", err
		}
		return enumType + "::" + memberName, nil
	}
	if field.HasDefaultValue() {
		return rustPrimitiveDefaultLiteral(field.GetDefaultValue()), nil
	}

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		return "String::new()", nil
	case cclValues.TypeNameBytes:
		return "Vec::new()", nil
	case cclValues.TypeNameBool:
		return "false", nil
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		return "0.0", nil
	default:
		if field.Type.IsCustomTypeEnum() {
			enumDef := field.Type.GetDefinition().GetEnumDefinition()
			enumType, err := c.getRustEnumTypeReference(enumDef, field.OwnedBy)
			if err != nil {
				return "", err
			}
			memberName, err := c.getRustEnumMemberName(enumDef, enumDef.Members[0])
			if err != nil {
				return "", err
			}
			return enumType + "::" + memberName, nil
		}
		if field.Type.IsCustomTypeModel() {
			return "None", nil
		}
		return "0", nil
	}
}
