package cclSanitizer

import (
	"fmt"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclAst"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

// ResolveTypeUsage resolves a type expression into a CCL type usage.
func ResolveTypeUsage(
	ctx *cclValues.CCLCodeContext,
	currentNamespace string,
	expr cclAst.TypeExpression,
) (*cclValues.CCLTypeUsage, error) {
	return ResolveTypeUsageForModel(ctx, currentNamespace, "", expr)
}

// ResolveTypeUsageForModel resolves a type expression with nested model scope.
func ResolveTypeUsageForModel(
	ctx *cclValues.CCLCodeContext,
	currentNamespace string,
	currentModelName string,
	expr cclAst.TypeExpression,
) (*cclValues.CCLTypeUsage, error) {
	if expr == nil {
		return nil, &TypeUsageResolutionError{
			Message: "missing type expression",
		}
	}

	if ctx == nil {
		return nil, &TypeUsageResolutionError{
			Message: "missing code context for type resolution",
		}
	}

	if currentNamespace == "" {
		currentNamespace = gValues.DefaultMainNamespace
	}

	switch node := expr.(type) {
	case *cclAst.SimpleTypeExpression:
		if node == nil {
			return nil, &TypeUsageResolutionError{
				Message: "invalid simple type expression",
			}
		}

		if node.IsBuiltinToken {
			usage := ctx.NewBuiltinTypeUsage(node.TypeName.Name)
			if usage == nil {
				return nil, &TypeUsageResolutionError{
					Message:        "unknown built-in type: " + node.TypeName.Name,
					SourcePosition: node.SourcePosition,
				}
			}
			return usage, nil
		}

		return resolveCustomTypeUsage(ctx, currentNamespace, currentModelName, node), nil
	case *cclAst.ArrayTypeExpression:
		if node == nil {
			return nil, &TypeUsageResolutionError{
				Message: "invalid array type expression",
			}
		}

		elementType, err := ResolveTypeUsageForModel(
			ctx,
			currentNamespace,
			currentModelName,
			node.ElementType,
		)
		if err != nil {
			return nil, err
		}

		if elementType == nil {
			return nil, &TypeUsageResolutionError{
				Message:        "array element type resolved to nil",
				SourcePosition: node.SourcePosition,
			}
		}

		return ctx.NewArrayTypeUsage(elementType, node.Length), nil
	default:
		return nil, &TypeUsageResolutionError{
			Message: fmt.Sprintf("unsupported type expression kind: %T", node),
		}
	}
}

func resolveCustomTypeUsage(
	ctx *cclValues.CCLCodeContext,
	currentNamespace string,
	currentModelName string,
	node *cclAst.SimpleTypeExpression,
) *cclValues.CCLTypeUsage {
	typeName := node.TypeName.Name
	namespace := node.TypeName.Namespace
	if namespace != "" {
		exactName := &cclValues.SimpleTypeName{
			TypeName:  typeName,
			Namespace: namespace,
		}
		if ctx.GetTypeDefinition(exactName) != nil {
			return ctx.NewCustomTypeUsage(exactName)
		}

		if currentNamespace != "" {
			nestedName := &cclValues.SimpleTypeName{
				TypeName:  typeName,
				Namespace: currentNamespace + "." + namespace,
			}
			if ctx.GetTypeDefinition(nestedName) != nil {
				return ctx.NewCustomTypeUsage(nestedName)
			}
		}

		return ctx.NewCustomTypeUsage(exactName)
	}

	if currentModelName != "" {
		nestedName := &cclValues.SimpleTypeName{
			TypeName:  typeName,
			Namespace: currentNamespace + "." + currentModelName,
		}
		if ctx.GetTypeDefinition(nestedName) != nil {
			return ctx.NewCustomTypeUsage(nestedName)
		}
	}

	return ctx.NewCustomTypeUsage(&cclValues.SimpleTypeName{
		TypeName:  typeName,
		Namespace: currentNamespace,
	})
}

func validateFieldTypeUsages(typeUsages []*fieldTypeUsageCheck) error {
	for _, usage := range typeUsages {
		if err := usage.validateTypeUsageCompleteness(usage.typeUsage); err != nil {
			return err
		}
	}

	return nil
}
