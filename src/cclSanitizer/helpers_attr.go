package cclSanitizer

import (
	"github.com/ccl-lang/ccl/src/core/cclAst"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

// ResolveAttributeUsage resolves an attribute AST node into a usage info.
func ResolveAttributeUsage(
	ctx *cclValues.CCLCodeContext,
	node cclAst.AttributeNodeBase,
) (*cclValues.AttributeUsageInfo, error) {
	if node == nil {
		return nil, &AttributeResolutionError{
			Message: "missing attribute node",
		}
	}

	if ctx == nil {
		return nil, &AttributeResolutionError{
			Message: "missing code context for attribute resolution",
		}
	}

	params, err := resolveAttributeParams(ctx, node.GetAttributeParams())
	if err != nil {
		return nil, err
	}

	return &cclValues.AttributeUsageInfo{
		Name:           node.GetAttributeName(),
		Parameters:     params,
		SourcePosition: node.GetSourcePosition(),
	}, nil
}

func resolveAttributeParams(
	ctx *cclValues.CCLCodeContext,
	params []*cclAst.AttributeParamNode,
) ([]*cclValues.ParameterInstance, error) {
	if len(params) == 0 {
		return nil, nil
	}

	result := make([]*cclValues.ParameterInstance, 0, len(params))
	for _, param := range params {
		if param == nil {
			return nil, &AttributeResolutionError{
				Message: "nil attribute parameter",
			}
		}

		if param.Value == nil {
			return nil, &AttributeResolutionError{
				Message:        "attribute parameter has no value",
				SourcePosition: param.SourcePosition,
			}
		}

		resolvedParam, err := resolveAttributeParamValue(ctx, param)
		if err != nil {
			return nil, err
		}

		result = append(result, resolvedParam)
	}

	return result, nil
}

func resolveAttributeParamValue(
	ctx *cclValues.CCLCodeContext,
	param *cclAst.AttributeParamNode,
) (*cclValues.ParameterInstance, error) {
	paramInstance := &cclValues.ParameterInstance{
		Name:           param.Name,
		SourcePosition: param.SourcePosition,
	}

	switch value := param.Value.(type) {
	case *cclAst.LiteralValueExpression:
		typeUsage := resolveLiteralTypeUsage(ctx, value)
		paramInstance.ChangeValueType(typeUsage)
		paramInstance.ChangeValue(value.Value)
		return paramInstance, nil
	case *cclAst.IdentifierValueExpression:
		if value == nil {
			return nil, &AttributeResolutionError{
				Message:        "invalid identifier value expression",
				SourcePosition: param.SourcePosition,
			}
		}

		targetVariable := ctx.GetGlobalVariable(value.Name)
		if targetVariable == nil {
			return nil, &AttributeResolutionError{
				Message:        "undefined identifier '" + value.Name + "'",
				SourcePosition: value.SourcePosition,
			}
		}

		paramInstance.ChangeValueType(targetVariable.Type)
		paramInstance.ChangeValue(&cclValues.VariableUsageInstance{
			Name:       value.Name,
			Definition: targetVariable,
		})
		return paramInstance, nil
	default:
		return nil, &AttributeResolutionError{
			Message:        "unsupported attribute value expression",
			SourcePosition: param.SourcePosition,
		}
	}
}

func resolveLiteralTypeUsage(
	ctx *cclValues.CCLCodeContext,
	value *cclAst.LiteralValueExpression,
) *cclValues.CCLTypeUsage {
	if value == nil {
		return nil
	}

	switch value.LiteralKind {
	case cclAst.AttributeLiteralKindString:
		return ctx.NewBuiltinTypeUsage(cclValues.TypeNameString)
	case cclAst.AttributeLiteralKindInt:
		return ctx.NewBuiltinTypeUsage(cclValues.TypeNameInt)
	case cclAst.AttributeLiteralKindFloat:
		return ctx.NewBuiltinTypeUsage(cclValues.TypeNameFloat)
	case cclAst.AttributeLiteralKindReserved:
		switch value.ReservedLiteral {
		case cclValues.ReservedLiteralTrue, cclValues.ReservedLiteralFalse:
			return ctx.NewBuiltinTypeUsage(cclValues.TypeNameBool)
		default:
			return nil
		}
	default:
		return nil
	}
}
