package goGenerator

import (
	"strconv"

	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

// Check the wire size before allocating. uint64 comparison avoids narrowing an
// untrusted uint32 length to int on 32-bit targets, and division avoids overflow.
func generateBinaryLengthCheck(
	builder *codeBuilder.CodeBuilder,
	lengthName string,
	minimumBytes int,
) {
	builder.MapVarPairs(
		"checkedLength", lengthName,
		"minimumBytes", strconv.Itoa(minimumBytes),
	)
	defer builder.UnmapVar(
		"checkedLength",
		"minimumBytes",
	)

	builder.LineD("if uint64($checkedLength) > uint64(buf.Len()) / $minimumBytes {").
		Indent().
		LineD("return $binaryLengthErrorReturn").
		Unindent().
		WriteLine("}")
}

func generateBinaryBytesRead(builder *codeBuilder.CodeBuilder, lengthName, dataName string) {
	registerGoImport(builder, "io")
	builder.MapVarPairs(
		"readLength", lengthName,
		"readData", dataName,
	)
	defer builder.UnmapVar("readLength", "readData")

	builder.LineD("var $readLength uint32").
		LineD("if err := binary.Read(buf, binaryEndian, &$readLength); err != nil {").
		Indent().
		LineD("return $binaryParseErrorReturn").
		Unindent().
		WriteLine("}")

	generateBinaryLengthCheck(builder, lengthName, 1)

	builder.LineD("$readData := make([]byte, $readLength)").
		LineD("if _, err := io.ReadFull(buf, $readData); err != nil {").
		Indent().
		LineD("return $binaryParseErrorReturn").
		Unindent().
		WriteLine("}")
}

func goBinaryMinimumElementSize(element *cclValues.CCLTypeUsage) int {
	typeName := element.GetName()
	if element.IsCustomTypeEnum() {
		typeName = element.GetEnumBaseTypeName()
	}
	switch typeName {
	case cclValues.TypeNameInt16, cclValues.TypeNameUint16:
		return 2
	case cclValues.TypeNameString, cclValues.TypeNameBytes,
		cclValues.TypeNameInt, cclValues.TypeNameUint,
		cclValues.TypeNameInt32, cclValues.TypeNameUint32,
		cclValues.TypeNameFloat32:
		return 4
	case cclValues.TypeNameInt64, cclValues.TypeNameUint64,
		cclValues.TypeNameFloat, cclValues.TypeNameFloat64,
		cclValues.TypeNameDateTime:
		return 8
	default:
		// A nullable model has at least its presence byte; bool/int8/uint8 use one byte.
		return 1
	}
}
