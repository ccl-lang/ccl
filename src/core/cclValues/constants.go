package cclValues

// general ccl type names
const (
	TypeNameString   = "string"
	TypeNameBytes    = "bytes"
	TypeNameInt      = "int"
	TypeNameInt8     = "int8"
	TypeNameInt16    = "int16"
	TypeNameInt32    = "int32"
	TypeNameInt64    = "int64"
	TypeNameUint     = "uint"
	TypeNameUint8    = "uint8"
	TypeNameUint16   = "uint16"
	TypeNameUint32   = "uint32"
	TypeNameUint64   = "uint64"
	TypeNameFloat    = "float"
	TypeNameFloat32  = "float32"
	TypeNameFloat64  = "float64"
	TypeNameBool     = "bool"
	TypeNameDateTime = "datetime"
)

// general ccl keyword names
const (
	KeywordNameModel = "model"
)

const (
	TypeFlagBuiltIn cclTypeFlag = 1 << iota
	TypeFlagArray
	TypeFlagMap
)
