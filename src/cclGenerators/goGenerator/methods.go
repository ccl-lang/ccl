package goGenerator

import (
	"errors"
	"path/filepath"
	"sort"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

//---------------------------------------------------------

func (c *GoGenerationContext) GenerateCode() ([]string, error) {
	if c.Options.PackageName == "" {
		return nil, errors.New("package name is required for Go code generation")
	}

	if err := c.GenerateConstants(); err != nil {
		return nil, err
	}
	if err := c.GenerateVars(); err != nil {
		return nil, err
	}
	if err := c.GenerateTypes(); err != nil {
		return nil, err
	}
	if err := c.GenerateHelpers(); err != nil {
		return nil, err
	}
	if err := c.GenerateMethods(); err != nil {
		return nil, err
	}

	c.finalizeCodeBuilders()
	return c.writeCodeBuilders()
}

//---------------------------------------------------------

func (c *GoGenerationContext) GenerateVars() error {
	c.getCodeBuilder(VarsFileName, "vars")
	return nil
}

//---------------------------------------------------------

func (c *GoGenerationContext) GenerateTypes() error {
	for _, currentTypeDef := range c.GetGenerationTypeDefinitions() {
		if currentTypeDef.IsCustomEnum() {
			enumDef := currentTypeDef.GetEnumDefinition()
			builder, err := c.getEnumCodeBuilder("types", enumDef)
			if err != nil {
				return err
			}
			if err = c.generateTypesForEnum(builder, enumDef); err != nil {
				return err
			}
			continue
		}

		if !currentTypeDef.IsCustomModel() {
			return &cclErrors.UnsupportedTypeDefinitionError{
				TypeName:       currentTypeDef.GetFullName(),
				TargetLanguage: CurrentLanguage.String(),
			}
		}

		currentModel := currentTypeDef.GetModelDefinition()
		builder, err := c.getTypeDefCodeBuilder("types", currentModel)
		if err != nil {
			return err
		}

		for _, currentField := range currentModel.Fields {
			targetType := currentField.Type
			if targetType.IsArray() {
				targetType = targetType.GetUnderlyingType()
			}
			if targetType.GetName() == cclValues.TypeNameDateTime {
				// Generated datetime fields use time.Time in the types file.
				registerGoImport(builder, "time")
			}
		}
	}

	for _, currentTypeDef := range c.GetGenerationTypeDefinitions() {
		if currentTypeDef.IsCustomEnum() {
			continue
		}

		currentModel := currentTypeDef.GetModelDefinition()
		builder, err := c.getTypeDefCodeBuilder("types", currentModel)
		if err != nil {
			return err
		}
		if err = c.generateTypesForModel(builder, currentModel); err != nil {
			return err
		}
	}

	return nil
}

func (c *GoGenerationContext) generateTypesForModel(
	builder *codeBuilder.CodeBuilder,
	currentModel *CCLModel,
) error {
	builder.WriteLine("type " + currentModel.Name + " struct {").
		Indent()
	for _, currentField := range currentModel.Fields {
		theGoType, err := c.getGoTypeForField(currentField)
		if err != nil {
			return err
		}
		builder.WriteLine(currentField.Name + " " + theGoType)
	}
	builder.Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *GoGenerationContext) generateTypesForEnum(
	builder *codeBuilder.CodeBuilder,
	enumDef *CCLEnum,
) error {
	enumTypeName, err := c.getGoEnumTypeName(enumDef)
	if err != nil {
		return err
	}

	builder.WriteLine("type " + enumTypeName + " " +
		c.getGoEnumBaseType(enumDef)).
		NewLine()
	return nil
}

//---------------------------------------------------------

func (c *GoGenerationContext) GenerateHelpers() error {
	c.getCodeBuilder(HelpersFileName, "helpers")
	for _, currentTypeDef := range c.GetGenerationTypeDefinitions() {
		if !currentTypeDef.IsCustomModel() {
			continue
		}

		currentModel := currentTypeDef.GetModelDefinition()
		if !goModelHasDefaults(currentModel) {
			continue
		}

		builder, err := c.getTypeDefCodeBuilder("helpers", currentModel)
		if err != nil {
			return err
		}

		if err = c.generateHelpersForModel(builder, currentModel); err != nil {
			return err
		}
	}
	return nil
}

func (c *GoGenerationContext) generateHelpersForModel(
	builder *codeBuilder.CodeBuilder,
	currentModel *CCLModel,
) error {
	builder.WriteLine("func New" + currentModel.Name + "() *" + currentModel.Name + " {").
		Indent().
		WriteLine("return &" + currentModel.Name + "{").
		Indent()
	for _, field := range currentModel.Fields {
		if !field.HasDefaultValue() {
			continue
		}

		defaultValue := goDefaultLiteral(field.GetDefaultValue())
		if enumRef := c.GetEnumDefaultReference(field); enumRef != nil {
			memberName, err := c.getGoEnumMemberName(enumRef.Enum, enumRef.Member)
			if err != nil {
				return err
			}
			defaultValue = memberName
		}
		builder.WriteLine(field.Name + ": " + defaultValue + ",")
	}
	builder.Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

//---------------------------------------------------------

func (c *GoGenerationContext) GenerateMethods() error {
	for _, currentTypeDef := range c.GetGenerationTypeDefinitions() {
		if currentTypeDef.IsCustomEnum() {
			continue
		}

		if !currentTypeDef.IsCustomModel() {
			return &cclErrors.UnsupportedTypeDefinitionError{
				TypeName:       currentTypeDef.GetFullName(),
				TargetLanguage: CurrentLanguage.String(),
			}
		}

		currentModel := currentTypeDef.GetModelDefinition()
		builder, err := c.getTypeDefCodeBuilder("methods", currentModel)
		if err != nil {
			return err
		}
		if err = c.generateMethodsForModel(builder, currentModel); err != nil {
			return err
		}
	}
	return nil
}

func (c *GoGenerationContext) generateMethodsForModel(
	builder *codeBuilder.CodeBuilder,
	currentModel *CCLModel,
) error {
	modelName := currentModel.GetName()
	builder.MapVarPairs(
		"model", "*"+modelName,
		"modelName", modelName,
		"modelIdConst", "ModelId"+modelName,
	)
	defer builder.UnmapVar(
		"model",
		"modelName",
		"modelIdConst",
	)

	builder.WriteLine("//------------------------------------------------------------").
		NewLine().
		LineD("func (m $model) GetModelId() int {").
		Indent().
		LineD("return $modelIdConst").
		Unindent().
		WriteLine("}")

	// generate CloneEmpty() method
	if c.NeedsCloneMethods(CurrentLanguage, currentModel) {
		if err := c.generateCloneMethods(builder); err != nil {
			return err
		}
	}

	if c.NeedsBinarySerialization(CurrentLanguage, currentModel) {
		if err := c.generateSerializeBinaryMethod(builder, currentModel); err != nil {
			return err
		}
		if err := c.generateDeserializeBinaryMethod(builder, currentModel); err != nil {
			return err
		}
	}
	if c.NeedsJsonSerialization(CurrentLanguage, currentModel) {
		if err := c.generateSerializeJsonMethods(builder, currentModel); err != nil {
			return err
		}
	}
	return nil
}

//---------------------------------------------------------

func (c *GoGenerationContext) getModelOutputFileGroup(currentModel *CCLModel) (string, error) {
	return c.GetOutputFileGroup(CurrentLanguage, currentModel)
}

func (c *GoGenerationContext) getTypeDefCodeBuilder(
	category string,
	currentModel *CCLModel,
) (*codeBuilder.CodeBuilder, error) {
	group, err := c.getModelOutputFileGroup(currentModel)
	if err != nil {
		return nil, err
	}

	return c.getCodeBuilder(getGoCategoryFileName(category, group), category), nil
}

func (c *GoGenerationContext) getEnumCodeBuilder(
	category string,
	enumDef *CCLEnum,
) (*codeBuilder.CodeBuilder, error) {
	group, err := c.getGoEnumOutputFileGroup(enumDef)
	if err != nil {
		return nil, err
	}

	return c.getCodeBuilder(getGoCategoryFileName(category, group), category), nil
}

func (c *GoGenerationContext) getCodeBuilder(
	relativePath string,
	section string,
) *codeBuilder.CodeBuilder {
	if c.CodeByPath == nil {
		c.CodeByPath = map[string]*codeBuilder.CodeBuilder{}
	}

	builder := c.CodeByPath[relativePath]
	if builder != nil {
		return builder
	}

	builder = codeBuilder.NewCodeBuilderWithOptions(&codeBuilder.CodeBuilderOptions{
		IndentationStr:  "\t",
		NewLineStr:      "\n",
		EnableDebugInfo: c.Options.GenerateDebugInfo,
	})
	builder.BeginSection(section)
	builder.AddCommentHeader("// THIS FILE IS AUTOGENERATED BY A CCL TOOL. DO NOT EDIT.").
		AddHeader("package " + c.Options.PackageName)

	c.CodeByPath[relativePath] = builder
	return builder
}

func (c *GoGenerationContext) finalizeCodeBuilders() {
	for _, builder := range c.CodeByPath {
		closeGoImportGroup(builder)
		builder.EndSection()
	}
}

func (c *GoGenerationContext) writeCodeBuilders() ([]string, error) {
	relativePaths := make([]string, 0, len(c.CodeByPath))
	for relativePath := range c.CodeByPath {
		relativePaths = append(relativePaths, relativePath)
	}
	sort.Strings(relativePaths)

	outputFiles := make([]string, 0, len(relativePaths))
	for _, relativePath := range relativePaths {
		builder := c.CodeByPath[relativePath]
		fullPath := filepath.Join(c.Options.OutputPath, relativePath)
		if err := c.WriteCodeFile(fullPath, builder.Build(nil)); err != nil {
			return nil, err
		}
		outputFiles = append(outputFiles, fullPath)
	}

	return outputFiles, nil
}
