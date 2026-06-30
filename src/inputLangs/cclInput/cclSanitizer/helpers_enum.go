package cclSanitizer

import (
	"math"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func sanitizeEnumDeclaration(
	ctx *cclValues.CCLCodeContext,
	definition *cclValues.SourceCodeDefinition,
	sourceFileId cclValues.SourceFileId,
	currentNamespace string,
	ownerModel *cclValues.ModelDefinition,
	enumAst *cclAst.EnumDecl,
) (*cclValues.EnumDefinition, *cclValues.CCLTypeDefinition, error) {
	enumNamespace := currentNamespace
	if ownerModel != nil {
		enumNamespace = ownerModel.GetFullName()
	}

	enumDef := &cclValues.EnumDefinition{
		SourceFileId:   sourceFileId,
		Name:           enumAst.Name,
		Namespace:      enumNamespace,
		OwnedBy:        ownerModel,
		SourcePosition: enumAst.SourcePosition,
	}

	baseType, err := resolveEnumBaseType(ctx, currentNamespace, ownerModel, enumAst)
	if err != nil {
		return nil, nil, err
	}
	enumDef.BaseType = baseType

	for _, attrAst := range enumAst.Attributes {
		attrUsage, err := ResolveAttributeUsage(ctx, attrAst)
		if err != nil {
			return nil, nil, err
		}
		attrUsage.SourceFileId = sourceFileId
		enumDef.Attributes = append(enumDef.Attributes, attrUsage)
	}

	if err := sanitizeEnumMembers(enumDef, enumAst); err != nil {
		return nil, nil, err
	}

	enumTypeDef, err := ctx.NewEnumTypeDefinition(&cclValues.SimpleTypeName{
		TypeName:  enumDef.Name,
		Namespace: enumNamespace,
	}, enumDef)
	if err != nil {
		return nil, nil, err
	}
	enumTypeDef.ChangeSourceFileId(sourceFileId)

	_ = definition
	return enumDef, enumTypeDef, nil
}

func resolveEnumBaseType(
	ctx *cclValues.CCLCodeContext,
	currentNamespace string,
	ownerModel *cclValues.ModelDefinition,
	enumAst *cclAst.EnumDecl,
) (*cclValues.CCLTypeUsage, error) {
	if enumAst.BaseType == nil {
		return ctx.NewBuiltinTypeUsage(cclValues.TypeNameInt32), nil
	}

	currentModelName := ""
	if ownerModel != nil {
		currentModelName = ownerModel.Name
	}

	baseType, err := ResolveTypeUsageForModel(
		ctx,
		currentNamespace,
		currentModelName,
		enumAst.BaseType,
	)
	if err != nil {
		return nil, err
	}

	name := baseType.GetName()
	switch name {
	case cclValues.TypeNameInt:
		return ctx.NewBuiltinTypeUsage(cclValues.TypeNameInt32), nil
	case cclValues.TypeNameUint:
		return ctx.NewBuiltinTypeUsage(cclValues.TypeNameUint32), nil
	case cclValues.TypeNameInt8,
		cclValues.TypeNameInt16,
		cclValues.TypeNameInt32,
		cclValues.TypeNameInt64,
		cclValues.TypeNameUint8,
		cclValues.TypeNameUint16,
		cclValues.TypeNameUint32,
		cclValues.TypeNameUint64:
		return baseType, nil
	default:
		return nil, &AstSanitizationError{
			Message:        "enum base type must be an integer type",
			SourcePosition: enumAst.SourcePosition,
		}
	}
}

func sanitizeEnumMembers(
	enumDef *cclValues.EnumDefinition,
	enumAst *cclAst.EnumDecl,
) error {
	minValue, maxValue := enumBaseTypeRange(enumDef.BaseType.GetName())
	nextValue := int64(0)
	memberNames := map[string]bool{}

	for _, memberAst := range enumAst.Members {
		if memberNames[memberAst.Name] {
			return &AstSanitizationError{
				Message:        "duplicate enum member: " + enumDef.Name + "." + memberAst.Name,
				SourcePosition: memberAst.SourcePosition,
			}
		}
		memberNames[memberAst.Name] = true

		value := nextValue
		if memberAst.Value != nil {
			value = *memberAst.Value
		}

		if value < minValue || value > maxValue {
			return &AstSanitizationError{
				Message:        "enum member value is outside the base type range",
				SourcePosition: memberAst.SourcePosition,
			}
		}

		enumDef.Members = append(enumDef.Members, &cclValues.EnumMemberDefinition{
			OwnedBy:        enumDef,
			Name:           memberAst.Name,
			Value:          value,
			SourcePosition: memberAst.SourcePosition,
		})

		if value == maxValue {
			nextValue = maxValue
			continue
		}
		nextValue = value + 1
	}

	return nil
}

func enumBaseTypeRange(typeName string) (int64, int64) {
	switch typeName {
	case cclValues.TypeNameInt8:
		return math.MinInt8, math.MaxInt8
	case cclValues.TypeNameInt16:
		return math.MinInt16, math.MaxInt16
	case cclValues.TypeNameInt32:
		return math.MinInt32, math.MaxInt32
	case cclValues.TypeNameInt64:
		return math.MinInt64, math.MaxInt64
	case cclValues.TypeNameUint8:
		return 0, math.MaxUint8
	case cclValues.TypeNameUint16:
		return 0, math.MaxUint16
	case cclValues.TypeNameUint32:
		return 0, math.MaxUint32
	case cclValues.TypeNameUint64:
		return 0, math.MaxInt64
	default:
		return math.MinInt32, math.MaxInt32
	}
}
