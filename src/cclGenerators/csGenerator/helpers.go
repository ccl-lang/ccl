package csGenerator

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	gen "github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/core/cclUtils"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

var (
	namespaceRegex = regexp.MustCompile(`namespace\s+([\w\.]+)`)
)

// GenerateCode generates C# code from the provided CCL source file.
func GenerateCode(options *gen.CodeGenerationOptions) (*gen.CodeGenerationResult, error) {
	context := &CSharpGenerationContext{
		CodeGenerationBase: &gen.CodeGenerationBase{
			Options: options,
		},
		ModelClasses: make(map[string]*codeBuilder.CodeBuilder),
	}

	return context.generateCode()
}

func (c *CSharpGenerationContext) getNamespace() string {
	// 1. Check for attribute on the first model (assuming all models in the same file share namespace preference if defined there)
	// Or we can check global attributes if CCL supports them.
	// The user said: "try to get #[CSharpNamespace("MyNamespace")] attribute"
	// We should check if any model has this attribute or if it's a file-level attribute (if supported).
	// Since we are generating for a set of models, let's check the first one or look for a common one.
	// Actually, the user might have meant an attribute on the model definition.
	// Let's check the first model's attributes.

	if len(c.Options.CCLDefinition.TypeDefinitions) > 0 {
		firstModel := c.Options.CCLDefinition.TypeDefinitions[0]
		if firstModel.IsCustomModel() {
			model := firstModel.GetModelDefinition()
			attr := c.GetGlobalOrModelAttributes(gValues.LanguageCS, AttributeNamespace, model)
			if !attr.IsEmpty() {
				params := attr.GetParamsAtAsStrings(0)
				if len(params) > 0 {
					return params[0]
				}
			}
		}
	}

	// 2. Peek at parent directory
	parentDir := filepath.Dir(c.Options.OutputPath)
	baseNamespace := ""
	foundCsFile := false

	entries, err := os.ReadDir(parentDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".cs") {
				content, err := os.ReadFile(filepath.Join(parentDir, entry.Name()))
				if err == nil {
					matches := namespaceRegex.FindStringSubmatch(string(content))
					if len(matches) > 1 {
						baseNamespace = matches[1]
						foundCsFile = true
						break
					}
				}
			}
		}
	}

	if foundCsFile {
		targetDirName := filepath.Base(c.Options.OutputPath)
		return baseNamespace + "." + targetDirName
	}

	// 3. Default
	return DefaultNamespace
}

func (c *CSharpGenerationContext) getCSharpType(field *CCLField) string {
	targetType := field.Type
	if targetType.IsArray() {
		targetType = targetType.GetUnderlyingType()
	}

	csType := ""
	switch targetType.GetName() {
	case cclValues.TypeNameString:
		csType = "string"
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		csType = "int"
	case cclValues.TypeNameInt8:
		csType = "sbyte"
	case cclValues.TypeNameInt16:
		csType = "short"
	case cclValues.TypeNameInt64:
		csType = "long"
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		csType = "uint"
	case cclValues.TypeNameUint8:
		csType = "byte"
	case cclValues.TypeNameUint16:
		csType = "ushort"
	case cclValues.TypeNameUint64:
		csType = "ulong"
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		csType = "float"
	case cclValues.TypeNameFloat64:
		csType = "double"
	case cclValues.TypeNameBool:
		csType = "bool"
	case cclValues.TypeNameBytes:
		csType = "byte[]"
	case cclValues.TypeNameDateTime:
		csType = "long" // Using timestamp
	default:
		if targetType.IsCustomTypeModel() {
			csType = targetType.GetName()
		} else {
			csType = "object"
		}
	}

	if field.IsArray() {
		if targetType.IsCustomTypeModel() {
			csType = "List<" + csType + ">"
		} else {
			csType = "List<" + csType + ">"
		}
	}

	return csType
}

func (c *CSharpGenerationContext) toPascalCase(s string) string {
	return cclUtils.ToPascalCase(s)
}
