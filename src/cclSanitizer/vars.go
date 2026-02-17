package cclSanitizer

import (
	"strings"

	"github.com/ccl-lang/ccl/src/core/cclValues"
)

var builtinTypeNamesLower = map[string]string{
	strings.ToLower(cclValues.TypeNameString):   cclValues.TypeNameString,
	strings.ToLower(cclValues.TypeNameBytes):    cclValues.TypeNameBytes,
	strings.ToLower(cclValues.TypeNameInt):      cclValues.TypeNameInt,
	strings.ToLower(cclValues.TypeNameInt8):     cclValues.TypeNameInt8,
	strings.ToLower(cclValues.TypeNameInt16):    cclValues.TypeNameInt16,
	strings.ToLower(cclValues.TypeNameInt32):    cclValues.TypeNameInt32,
	strings.ToLower(cclValues.TypeNameInt64):    cclValues.TypeNameInt64,
	strings.ToLower(cclValues.TypeNameUint):     cclValues.TypeNameUint,
	strings.ToLower(cclValues.TypeNameUint8):    cclValues.TypeNameUint8,
	strings.ToLower(cclValues.TypeNameUint16):   cclValues.TypeNameUint16,
	strings.ToLower(cclValues.TypeNameUint32):   cclValues.TypeNameUint32,
	strings.ToLower(cclValues.TypeNameUint64):   cclValues.TypeNameUint64,
	strings.ToLower(cclValues.TypeNameFloat):    cclValues.TypeNameFloat,
	strings.ToLower(cclValues.TypeNameFloat32):  cclValues.TypeNameFloat32,
	strings.ToLower(cclValues.TypeNameFloat64):  cclValues.TypeNameFloat64,
	strings.ToLower(cclValues.TypeNameBool):     cclValues.TypeNameBool,
	strings.ToLower(cclValues.TypeNameDateTime): cclValues.TypeNameDateTime,
	strings.ToLower(cclValues.TypeNamePointer):  cclValues.TypeNamePointer,
	strings.ToLower(cclValues.TypeNameArray):    cclValues.TypeNameArray,
}
