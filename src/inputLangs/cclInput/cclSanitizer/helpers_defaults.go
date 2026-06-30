package cclSanitizer

import (
	"fmt"
	"strings"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func resolveFieldDefaultValue(
	ctx *cclValues.CCLCodeContext,
	fieldType *cclValues.CCLTypeUsage,
	currentNamespace string,
	currentModelName string,
	value cclAst.ValueExpression,
) (any, error) {
	if fieldType == nil {
		return nil, &AstSanitizationError{
			Message: "missing field type for default value",
		}
	}

	if fieldType.IsCustomTypeEnum() {
		return resolveEnumDefaultValue(fieldType, value)
	}

	switch expr := value.(type) {
	case *cclAst.LiteralValueExpression:
		return resolvePrimitiveDefaultValue(fieldType, expr)
	case *cclAst.IdentifierValueExpression:
		return resolveLegacyVariableDefault(ctx, expr)
	case *cclAst.QualifiedIdentifierValueExpression:
		return nil, &AstSanitizationError{
			Message:        "qualified default values are only supported for enum fields",
			SourcePosition: expr.SourcePosition,
		}
	default:
		_ = currentNamespace
		_ = currentModelName
		return nil, &AstSanitizationError{
			Message:        fmt.Sprintf("unsupported field default value: %T", value),
			SourcePosition: value.GetSourcePosition(),
		}
	}
}

func resolveEnumDefaultValue(
	fieldType *cclValues.CCLTypeUsage,
	value cclAst.ValueExpression,
) (any, error) {
	expr, ok := value.(*cclAst.QualifiedIdentifierValueExpression)
	if !ok || len(expr.Parts) < 2 {
		return nil, &AstSanitizationError{
			Message:        "enum default value must be qualified with its enum type",
			SourcePosition: value.GetSourcePosition(),
		}
	}

	enumDef := fieldType.GetDefinition().GetEnumDefinition()
	memberName := expr.Parts[len(expr.Parts)-1]
	enumRefName := strings.Join(expr.Parts[:len(expr.Parts)-1], ".")
	if !doesEnumReferenceMatchFieldType(enumDef, enumRefName) {
		return nil, &AstSanitizationError{
			Message:        "enum default value type does not match field type",
			SourcePosition: expr.SourcePosition,
		}
	}

	member := enumDef.GetMemberByName(memberName)
	if member == nil {
		return nil, &AstSanitizationError{
			Message:        "unknown enum member: " + memberName,
			SourcePosition: expr.SourcePosition,
		}
	}

	return &cclValues.EnumMemberReference{
		Enum:   enumDef,
		Member: member,
	}, nil
}

func doesEnumReferenceMatchFieldType(
	enumDef *cclValues.EnumDefinition,
	enumRefName string,
) bool {
	if enumDef == nil {
		return false
	}

	if enumRefName == enumDef.Name || enumRefName == enumDef.GetFullName() {
		return true
	}

	if enumDef.IsNested() && enumDef.OwnedBy != nil {
		return enumRefName == enumDef.OwnedBy.Name+"."+enumDef.Name ||
			enumRefName == enumDef.GetFullName()
	}

	return false
}

func resolvePrimitiveDefaultValue(
	fieldType *cclValues.CCLTypeUsage,
	expr *cclAst.LiteralValueExpression,
) (any, error) {
	if !fieldType.IsBuiltIn() {
		return nil, &AstSanitizationError{
			Message:        "default values for model fields must use enum or primitive types",
			SourcePosition: expr.SourcePosition,
		}
	}

	switch fieldType.GetName() {
	case cclValues.TypeNameString:
		if expr.LiteralKind == cclAst.AttributeLiteralKindString {
			return expr.Value, nil
		}
	case cclValues.TypeNameBool:
		if value, ok := expr.Value.(bool); ok {
			return value, nil
		}
	case cclValues.TypeNameFloat,
		cclValues.TypeNameFloat32,
		cclValues.TypeNameFloat64:
		switch value := expr.Value.(type) {
		case float64:
			return value, nil
		case int64:
			return float64(value), nil
		}
	case cclValues.TypeNameInt,
		cclValues.TypeNameInt8,
		cclValues.TypeNameInt16,
		cclValues.TypeNameInt32,
		cclValues.TypeNameInt64,
		cclValues.TypeNameUint,
		cclValues.TypeNameUint8,
		cclValues.TypeNameUint16,
		cclValues.TypeNameUint32,
		cclValues.TypeNameUint64:
		if value, ok := expr.Value.(int64); ok {
			return value, nil
		}
	}

	return nil, &AstSanitizationError{
		Message:        "default value does not match field type",
		SourcePosition: expr.SourcePosition,
	}
}

func resolveLegacyVariableDefault(
	ctx *cclValues.CCLCodeContext,
	expr *cclAst.IdentifierValueExpression,
) (any, error) {
	targetVariable := ctx.GetGlobalVariable(expr.Name)
	if targetVariable == nil {
		return nil, &AstSanitizationError{
			Message:        "undefined identifier '" + expr.Name + "'",
			SourcePosition: expr.SourcePosition,
		}
	}

	if targetVariable.IsAutomatic() {
		return &cclValues.VariableUsageInstance{
			Name:       expr.Name,
			Definition: targetVariable,
		}, nil
	}

	return targetVariable.GetValue(), nil
}
