package cclSanitizer

import (
	"strings"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAttr"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func isEnumMappingAttribute(name cclAttr.CCLAttributeName) bool {
	return name == cclAttr.AttrEnumMapSameMembers || name == cclAttr.AttrEnumMapUnknownMember
}

func resolveEnumAttributeUsage(
	ctx *cclValues.CCLCodeContext,
	node *cclAst.AttributeNode,
) (*cclValues.AttributeUsageInfo, error) {
	if !isEnumMappingAttribute(node.Name) {
		return ResolveAttributeUsage(ctx, node)
	}
	languages, err := resolveAttributeLanguages(node.Languages, node.SourcePosition)
	if err != nil {
		return nil, err
	}
	attr := &cclValues.AttributeUsageInfo{
		Name:           node.Name,
		Languages:      languages,
		SourcePosition: node.SourcePosition,
	}
	if len(node.Params) != 1 || node.Params[0].Name != "" {
		return nil, enumMappingError(attr, "requires exactly one positional enum symbol reference")
	}
	var parts []string
	switch value := node.Params[0].Value.(type) {
	case *cclAst.IdentifierValueExpression:
		parts = []string{value.Name}
	case *cclAst.QualifiedIdentifierValueExpression:
		parts = value.Parts
	default:
		return nil, enumMappingError(attr, "requires an enum symbol reference, not a literal or variable value")
	}
	if attr.Name == cclAttr.AttrEnumMapUnknownMember && len(parts) < 2 {
		return nil, enumMappingError(attr, "requires a qualified member reference: CurrentEnum.Member")
	}
	param := &cclValues.ParameterInstance{SourcePosition: node.Params[0].SourcePosition}
	param.ChangeValue(&enumMappingSymbol{parts: parts})
	attr.Parameters = []*cclValues.ParameterInstance{param}
	return attr, nil
}

func resolveEnumMappingSymbol(
	ctx *cclValues.CCLCodeContext,
	source *cclValues.EnumDefinition,
	attr *cclValues.AttributeUsageInfo,
) error {
	param := attr.GetParamAt(0)
	symbol, pending := param.GetValue().(*enumMappingSymbol)
	if !pending {
		return nil
	}
	parts := symbol.parts
	memberName := ""
	if attr.Name == cclAttr.AttrEnumMapUnknownMember {
		memberName = parts[len(parts)-1]
		parts = parts[:len(parts)-1]
	}
	namespace := source.Namespace
	modelName := ""
	if source.OwnedBy != nil {
		namespace = source.OwnedBy.Namespace
		modelName = source.OwnedBy.Name
	}
	usage, err := ResolveTypeUsageForModel(ctx, namespace, modelName, &cclAst.SimpleTypeExpression{
		TypeName: cclAst.SimpleTypeName{
			Name:      parts[len(parts)-1],
			Namespace: strings.Join(parts[:len(parts)-1], "."),
		},
		SourcePosition: param.SourcePosition,
	})
	if err != nil {
		return err
	}
	if !usage.IsCustomTypeEnum() {
		return enumMappingError(attr, "reference '"+strings.Join(parts, ".")+"' must resolve to an enum")
	}
	enumDef := usage.GetDefinition().GetEnumDefinition()
	param.ChangeValueType(usage)
	if memberName == "" {
		param.ChangeValue(enumDef)
		return nil
	}
	if enumDef != source {
		return enumMappingError(attr, "fallback must belong to the annotated enum '"+source.GetFullName()+"'")
	}
	member := enumDef.GetMemberByName(memberName)
	if member == nil {
		return enumMappingError(attr, "unknown enum member: "+enumDef.GetFullName()+"."+memberName)
	}
	param.ChangeValue(&cclValues.EnumMemberReference{Enum: enumDef, Member: member})
	return nil
}

func enumMappingError(attr *cclValues.AttributeUsageInfo, message string) error {
	return &cclErrors.InvalidAttributeUsageError{
		AttrName:       attr.Name,
		Message:        message,
		SourcePosition: attr.SourcePosition,
	}
}
