package cclSanitizer

import (
	"fmt"

	"github.com/ccl-lang/ccl/src/core/cclAst"
	"github.com/ccl-lang/ccl/src/core/cclValues"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

// ResolveTypeUsage resolves a type expression into a CCL type usage.
func ResolveTypeUsage(
	ctx *cclValues.CCLCodeContext,
	currentNamespace string,
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

		namespace := node.TypeName.Namespace
		if namespace == "" {
			namespace = currentNamespace
		}

		return ctx.NewCustomTypeUsage(&cclValues.SimpleTypeName{
			TypeName:  node.TypeName.Name,
			Namespace: namespace,
		}), nil
	case *cclAst.ArrayTypeExpression:
		if node == nil {
			return nil, &TypeUsageResolutionError{
				Message: "invalid array type expression",
			}
		}

		elementType, err := ResolveTypeUsage(ctx, currentNamespace, node.ElementType)
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
