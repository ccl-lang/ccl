package cclValues

import "github.com/ccl-lang/ccl/src/core/cclUtils"

// cclTypeFlag represents flags for a specific ccl type definition.
type cclTypeFlag int

// CCLReservedLiteral is a type alias for reserved literals in CCL.
// These include: null, true, false, nil, self, super, this, etc...
type CCLReservedLiteral = string

// ParameterInstance is a struct that represents a passed parameter instance.
// This parameter is for when the user is passing a parameter to a function or any
// other place in the source code, such as attributes.
// The difference between this struct and VariableUsageInstance is that
// this struct must be used when passing a value directly, such as in
// functionName(1, 2, 3) or [AttrName(1, 2, 3)], while VariableUsageInstance
// is used when passing a variable by its name, such as in
// functionName(var1, var2, var3) or [AttrName(var1, var2, var3)].
type ParameterInstance struct {
	// Name is the name of the parameter.
	// Please note that this field might be empty, if the programmer is
	// passing a parameter without specifying its name; such as in
	// functionName(1, 2, 3) or [AttrName(1, 2, 3)]
	Name string

	// value is the value of the parameter, specified in the source code.
	// Please note that this field is not exported, you should use the
	// methods to get or set this field.
	value any

	// ValueType is the type of the parameter.
	ValueType *CCLTypeUsage

	// SourcePosition is the position of the attribute in the source code.
	SourcePosition *cclUtils.SourceCodePosition
}

// SimpleTypeName is a struct that represents a simple type name
// with its namespace.
type SimpleTypeName struct {
	// TypeName is the normalized type name.
	TypeName string

	// Namespace is the namespace of the type.
	Namespace string
}
