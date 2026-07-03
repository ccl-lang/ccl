package goGenerator

import (
	"strings"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

//---------------------------------------------------------

func (c *GoGenerationContext) generateDeserializeBinaryMethod(
	builder *codeBuilder.CodeBuilder,
	model *CCLModel,
) error {
	// Generated DeserializeBinary methods read input with bytes.Reader.
	registerGoImport(builder, "bytes")
	// Generated DeserializeBinary methods read fields with encoding/binary.
	registerGoImport(builder, "encoding/binary")

	endian, err := c.GetBinarySerializationEndian(CurrentLanguage, model)
	if err != nil {
		return err
	}
	strictBinaryParsing, err := c.UsesStrictBinaryParsing(CurrentLanguage, model)
	if err != nil {
		return err
	}
	binaryEndianInit := "binary.LittleEndian"
	if endian == gValues.EndianBig {
		binaryEndianInit = "binary.BigEndian"
	}
	binaryParseErrorReturn := "nil"
	if strictBinaryParsing {
		binaryParseErrorReturn = "err"
	}

	builder.ExpectMappedVars(
		"model",
	)
	builder.MapVarPairs(
		"binaryEndianInit", binaryEndianInit,
		"binaryParseErrorReturn", binaryParseErrorReturn,
	)
	defer builder.UnmapVar(
		"binaryEndianInit",
		"binaryParseErrorReturn",
	)

	builder.LineD("func (m $model) DeserializeBinary(data []byte) error {").
		Indent().
		// add nil checker or when the len(data) is 0 or (len(data) == 1 and data[0] == 0)
		WriteLine("if m == nil || len(data) == 0 || (len(data) == 1 && data[0] == 0) {").
		Indent().
		WriteLine("return nil").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("buf := bytes.NewReader(data)").
		LineD("binaryEndian := $binaryEndianInit").
		NewLine()

	for _, field := range model.Fields {
		if field.IsArray() {
			err := c.generateArrayDeserializeBinaryMethod(builder, field)
			if err != nil {
				return err
			}

			continue
		}

		err := c.generateFieldDeserializeBinaryMethod(builder, field)
		if err != nil {
			return err
		}
	}

	builder.WriteLine("return nil").
		Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *GoGenerationContext) generateFieldDeserializeBinaryMethod(
	builder *codeBuilder.CodeBuilder,
	field *CCLField,
) error {
	fieldTypeName := field.Type.GetName()
	isCustomType := field.Type.IsCustomTypeModel()
	// isPointer := isCustomType //TODO: Find a way to specify this
	fieldName := field.Name
	fieldVar := "m." + fieldName
	fName := strings.ToLower(string(fieldName[0])) + fieldName[1:]
	fLenName := fName + "Len"
	fNameStrBytes := fName + "StrBytes"
	fNameUnix := fName + "Unix"

	builder.MapVarPairs(
		"field", fieldVar,
		"fieldLen", fLenName,
		"fieldStrBytes", fNameStrBytes,
		"fieldUnix", fNameUnix,
		"fieldName", fieldName,
	)
	defer builder.UnmapVar(
		"field",
		"fieldLen",
		"fieldStrBytes",
		"fieldUnix",
		"fieldName",
	)

	switch fieldTypeName {
	case cclValues.TypeNameString:
		builder.LineD("var $fieldLen uint32").
			LineD("if err := binary.Read(buf, binaryEndian, &$fieldLen); err != nil {").
			Indent().
			LineD("return $binaryParseErrorReturn").
			Unindent().
			WriteLine("}").
			LineD("$fieldStrBytes := make([]byte, $fieldLen)").
			LineD("if $fieldLen > 0 {").
			Indent().
			LineD("if _, err := buf.Read($fieldStrBytes); err != nil {").
			Indent().
			LineD("return $binaryParseErrorReturn").
			Unindent().
			WriteLine("}").
			Unindent().
			WriteLine("}").
			LineD("$field = string($fieldStrBytes)")
	case cclValues.TypeNameBytes:
		builder.LineD("var $fieldLen uint32").
			LineD("if err := binary.Read(buf, binaryEndian, &$fieldLen); err != nil {").
			Indent().
			LineD("return $binaryParseErrorReturn").
			Unindent().
			WriteLine("}").
			LineD("bytesData := make([]byte, $fieldLen)").
			LineD("if $fieldLen > 0 {").
			Indent().
			WriteLine("if _, err := buf.Read(bytesData); err != nil {").
			Indent().
			LineD("return $binaryParseErrorReturn").
			Unindent().
			WriteLine("}").
			Unindent().
			WriteLine("}").
			LineD("$field = bytesData")
	case cclValues.TypeNameDateTime:
		// Generated datetime deserialization rebuilds values with time.Unix.
		registerGoImport(builder, "time")
		builder.LineD("var $fieldUnix int64").
			LineD("if err := binary.Read(buf, binaryEndian, &$fieldUnix); err != nil {").
			Indent().
			LineD("return $binaryParseErrorReturn").
			Unindent().
			WriteLine("}").
			LineD("$field = time.Unix(0, $fieldUnix)")
	default:
		if isCustomType {
			lenVarName := fName + "BytesLen"
			bytesVarName := fName + "Bytes"
			fieldType := field.Type.GetName()

			builder.MapVarPairs(
				"fieldBytesLen", lenVarName,
				"fieldBytes", bytesVarName,
				"fieldType", fieldType,
			)

			builder.LineD("var $fieldBytesLen uint32").
				LineD("if err := binary.Read(buf, binaryEndian, &$fieldBytesLen); err != nil {").
				Indent().
				LineD("return $binaryParseErrorReturn").
				Unindent().
				WriteLine("}")

			builder.LineD("$fieldBytes := make([]byte, $fieldBytesLen)").
				LineD("if $fieldBytesLen > 0 {").
				Indent().
				LineD("if _, err := buf.Read($fieldBytes); err != nil {").
				Indent().
				LineD("return $binaryParseErrorReturn").
				Unindent().
				WriteLine("}").
				Unindent().
				WriteLine("}")

			// make sure m.field is not nil ONLY when len(bytesVarName) != 0 and !(len(bytesVarName) == 1 and bytesVarName[0] == 0)
			builder.LineD("if m.$fieldName == nil && len($fieldBytes) != 0 && !(len($fieldBytes) == 1 && $fieldBytes[0] == 0) {").
				Indent().
				LineD("m.$fieldName = new($fieldType)").
				Unindent().
				WriteLine("}")

			builder.LineD("if err := m.$fieldName.DeserializeBinary($fieldBytes); err != nil {").
				Indent().
				LineD("return $binaryParseErrorReturn").
				Unindent().
				WriteLine("}")

			builder.UnmapVar(
				"fieldBytesLen",
				"fieldBytes",
				"fieldType",
			)
		} else {
			if fieldTypeName == cclValues.TypeNameInt || field.Type.IsCustomTypeEnum() {
				toReadName := "tmp" + fieldName
				readType := goBaseIntegerTypeForRead(field.Type)
				assignExpr := readType + "(" + toReadName + ")"
				if field.Type.IsCustomTypeEnum() {
					enumTypeName, err := c.getGoEnumTypeName(field.Type.GetDefinition().GetEnumDefinition())
					if err != nil {
						return err
					}
					assignExpr = enumTypeName + "(" + toReadName + ")"
				} else {
					assignExpr = "int(" + toReadName + ")"
				}
				builder.MapVarPairs(
					"toRead", toReadName,
					"readType", readType,
					"assignExpr", assignExpr,
				)
				builder.LineD("var $toRead $readType").
					LineD("if err := binary.Read(buf, binaryEndian, &$toRead); err != nil {").
					Indent().
					LineD("return $binaryParseErrorReturn").
					Unindent().
					WriteLine("}").
					LineD("$field = $assignExpr")
				builder.UnmapVar("toRead", "readType", "assignExpr")
				return nil
			}

			// all other types can be read directly
			builder.LineD("if err := binary.Read(buf, binaryEndian, &$field); err != nil {").
				Indent().
				LineD("return $binaryParseErrorReturn").
				Unindent().
				WriteLine("}")
		}
	}
	return nil
}

func (c *GoGenerationContext) generateArrayDeserializeBinaryMethod(
	builder *codeBuilder.CodeBuilder,
	field *CCLField,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	targetFieldTypeName := targetFieldType.GetName()
	isCustomType := targetFieldType.IsCustomTypeModel()
	isPointer := isCustomType //TODO: Find a way to specify this
	fieldName := field.Name
	fieldVar := "m." + fieldName
	fName := strings.ToLower(string(fieldName[0])) + fieldName[1:]
	fLenName := fName + "Len"
	fieldRealType := ""
	if isCustomType {
		fieldRealType = "*" + targetFieldTypeName
	} else {
		resolvedType, err := c.getGoTypeForUsage(targetFieldType, field)
		if err != nil {
			return err
		}
		fieldRealType = resolvedType
	}

	builder.MapVarPairs(
		"field", fieldVar,
		"fieldLen", fLenName,
		"fieldRealType", fieldRealType,
	)
	defer builder.UnmapVar(
		"field",
		"fieldLen",
		"fieldRealType",
	)

	builder.LineD("var $fieldLen uint32").
		LineD("if err := binary.Read(buf, binaryEndian, &$fieldLen); err != nil {").
		Indent().
		LineD("return $binaryParseErrorReturn").
		Unindent().
		WriteLine("}").
		LineD("$field = make([]$fieldRealType, $fieldLen)").
		LineD("for i := uint32(0); i < $fieldLen; i++ {").
		Indent()
	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("var elemLen uint32").
			WriteLine("if err := binary.Read(buf, binaryEndian, &elemLen); err != nil {").
			Indent().
			LineD("return $binaryParseErrorReturn").
			Unindent().
			WriteLine("}").
			WriteLine("elemBytes := make([]byte, elemLen)").
			WriteLine("if elemLen > 0 {").
			Indent().
			WriteLine("if _, err := buf.Read(elemBytes); err != nil {").
			Indent().
			LineD("return $binaryParseErrorReturn").
			Unindent().
			WriteLine("}").
			Unindent().
			WriteLine("}").
			LineD("$field[i] = string(elemBytes)")
	case cclValues.TypeNameBytes:
		builder.WriteLine("var elemLen uint32").
			WriteLine("if err := binary.Read(buf, binaryEndian, &elemLen); err != nil {").
			Indent().
			LineD("return $binaryParseErrorReturn").
			Unindent().
			WriteLine("}").
			WriteLine("elemBytes := make([]byte, elemLen)").
			WriteLine("if elemLen > 0 {").
			Indent().
			WriteLine("if _, err := buf.Read(elemBytes); err != nil {").
			Indent().
			LineD("return $binaryParseErrorReturn").
			Unindent().
			WriteLine("}").
			Unindent().
			WriteLine("}").
			LineD("$field[i] = elemBytes")
	case cclValues.TypeNameDateTime:
		// Generated datetime array deserialization rebuilds elements with time.Unix.
		registerGoImport(builder, "time")
		builder.WriteLine("var elemUnix int64").
			WriteLine("if err := binary.Read(buf, binaryEndian, &elemUnix); err != nil {").
			Indent().
			LineD("return $binaryParseErrorReturn").
			Unindent().
			WriteLine("}").
			LineD("$field[i] = time.Unix(0, elemUnix)")
	default:
		if isCustomType {
			builder.MapVarPairs("fieldType", targetFieldType.GetName())
			if isPointer {
				builder.LineD("var elem $fieldRealType = new($fieldType)")
			} else {
				builder.LineD("var elem $fieldRealType")
			}
			builder.WriteLine("var elemLen uint32").
				WriteLine("if err := binary.Read(buf, binaryEndian, &elemLen); err != nil {").
				Indent().
				LineD("return $binaryParseErrorReturn").
				Unindent().
				WriteLine("}").
				WriteLine("elemBytes := make([]byte, elemLen)").
				WriteLine("if elemLen > 0 {").
				Indent().
				WriteLine("if _, err := buf.Read(elemBytes); err != nil {").
				Indent().
				LineD("return $binaryParseErrorReturn").
				Unindent().
				WriteLine("}").
				Unindent().
				WriteLine("}").
				WriteLine("if err := elem.DeserializeBinary(elemBytes); err != nil {").
				Indent().
				LineD("return $binaryParseErrorReturn").
				Unindent().
				WriteLine("}").
				LineD("$field[i] = elem")
			builder.UnmapVar("fieldType")
		} else {
			if targetFieldTypeName == cclValues.TypeNameInt || targetFieldType.IsCustomTypeEnum() {
				readType := goBaseIntegerTypeForRead(targetFieldType)
				assignExpr := "int(elem)"
				if targetFieldType.IsCustomTypeEnum() {
					enumTypeName, err := c.getGoEnumTypeName(targetFieldType.GetDefinition().GetEnumDefinition())
					if err != nil {
						return err
					}
					assignExpr = enumTypeName + "(elem)"
				}
				builder.MapVarPairs(
					"readType", readType,
					"assignExpr", assignExpr,
				)
				builder.LineD("var elem $readType").
					WriteLine("if err := binary.Read(buf, binaryEndian, &elem); err != nil {").
					Indent().
					LineD("return $binaryParseErrorReturn").
					Unindent().
					WriteLine("}").
					LineD("$field[i] = $assignExpr")
				builder.UnmapVar("readType", "assignExpr")
				break
			}

			builder.LineD("if err := binary.Read(buf, binaryEndian, &$field[i]); err != nil {").
				Indent().
				LineD("return $binaryParseErrorReturn").
				Unindent().
				WriteLine("}")
		}
	}

	builder.Unindent().
		WriteLine("}")
	return nil
}

//---------------------------------------------------------
