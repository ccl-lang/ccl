package pyGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *PythonGenerationContext) generateSerializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	builder.WriteLine("def serialize_binary(self) -> bytes:").
		Indent().
		WriteLine("buffer = bytearray()").
		NewLine()

	for i := range model.Fields {
		field := model.Fields[i]
		if field.IsArray() {
			c.generateArraySerializeBinary(field, builder)
			continue
		}

		c.generateFieldSerializeBinary(field, builder)
	}

	builder.WriteLine("return bytes(buffer)").
		UnindentLine().
		NewLine()

	return nil
}

func (c *PythonGenerationContext) generateFieldSerializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	fieldRawName := caseUtils.ToSnakeCase(field.Name)
	fieldName := "self." + fieldRawName
	builder.MapVarPairs(
		"field", fieldName,
		"rawField", fieldRawName,
		"rawFieldBytes", fieldRawName+"_bytes",
	)
	defer builder.UnmapVar(
		"field",
		"rawField",
		"rawFieldBytes",
	)

	switch pythonStorageTypeName(field.Type) {
	case cclValues.TypeNameString:
		builder.LineD(`$rawFieldBytes = $field.encode("utf-8")`).
			LineD("buffer.extend(struct.pack('<I', len($rawFieldBytes)))").
			LineD("buffer.extend($rawFieldBytes)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.LineD("buffer.extend(struct.pack('<i', $field))")
	case cclValues.TypeNameInt8:
		builder.LineD("buffer.extend(struct.pack('<b', $field))")
	case cclValues.TypeNameInt16:
		builder.LineD("buffer.extend(struct.pack('<h', $field))")
	case cclValues.TypeNameInt64:
		builder.LineD("buffer.extend(struct.pack('<q', $field))")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.LineD("buffer.extend(struct.pack('<I', $field))")
	case cclValues.TypeNameUint8:
		builder.LineD("buffer.extend(struct.pack('<B', $field))")
	case cclValues.TypeNameUint16:
		builder.LineD("buffer.extend(struct.pack('<H', $field))")
	case cclValues.TypeNameUint64:
		builder.LineD("buffer.extend(struct.pack('<Q', $field))")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.LineD("buffer.extend(struct.pack('<f', $field))")
	case cclValues.TypeNameFloat64:
		builder.LineD("buffer.extend(struct.pack('<d', $field))")
	case cclValues.TypeNameBool:
		builder.LineD("buffer.extend(struct.pack('<b', 1 if $field else 0))")
	case cclValues.TypeNameBytes:
		builder.LineD("buffer.extend(struct.pack('<I', len($field)))").
			LineD("buffer.extend($field)")
	case cclValues.TypeNameDateTime:
		builder.LineD("buffer.extend(struct.pack('<q', $field))")
	default:
		// Custom type handling
		if field.IsCustomTypeModel() {
			builder.LineD("if $field:").
				Indent().
				LineD("$rawFieldBytes = $field.serialize_binary()").
				Unindent().
				WriteLine("else:").
				Indent().
				LineD("$rawFieldBytes = b'\\x00'").
				Unindent().
				LineD("buffer.extend(struct.pack('<I', len($rawFieldBytes)))").
				LineD("buffer.extend($rawFieldBytes)")
		}
	}
	builder.NewLine()
}

func (c *PythonGenerationContext) generateArraySerializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldName := "self." + caseUtils.ToSnakeCase(field.Name)
	builder.MapVarPairs(
		"field", fieldName,
	)
	defer builder.UnmapVar("field")

	builder.LineD("buffer.extend(struct.pack('<I', len($field)))").
		LineD("for item in $field:").
		Indent()

	switch pythonStorageTypeName(targetFieldType) {
	case cclValues.TypeNameString:
		builder.WriteLine("item_bytes = item.encode(\"utf-8\")").
			WriteLine("buffer.extend(struct.pack('<I', len(item_bytes)))").
			WriteLine("buffer.extend(item_bytes)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine("buffer.extend(struct.pack('<i', item))")
	case cclValues.TypeNameInt8:
		builder.WriteLine("buffer.extend(struct.pack('<b', item))")
	case cclValues.TypeNameInt16:
		builder.WriteLine("buffer.extend(struct.pack('<h', item))")
	case cclValues.TypeNameInt64:
		builder.WriteLine("buffer.extend(struct.pack('<q', item))")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine("buffer.extend(struct.pack('<I', item))")
	case cclValues.TypeNameUint8:
		builder.WriteLine("buffer.extend(struct.pack('<B', item))")
	case cclValues.TypeNameUint16:
		builder.WriteLine("buffer.extend(struct.pack('<H', item))")
	case cclValues.TypeNameUint64:
		builder.WriteLine("buffer.extend(struct.pack('<Q', item))")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.WriteLine("buffer.extend(struct.pack('<f', item))")
	case cclValues.TypeNameFloat64:
		builder.WriteLine("buffer.extend(struct.pack('<d', item))")
	case cclValues.TypeNameBool:
		builder.WriteLine("buffer.extend(struct.pack('<b', 1 if item else 0))")
	case cclValues.TypeNameBytes:
		builder.WriteLine("buffer.extend(struct.pack('<I', len(item)))").
			WriteLine("buffer.extend(item)")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("if item:").
				Indent().
				WriteLine("item_bytes = item.serialize_binary()").
				Unindent().
				WriteLine("else:").
				Indent().
				WriteLine("item_bytes = b'\\x00'").
				Unindent().
				WriteLine("buffer.extend(struct.pack('<I', len(item_bytes)))").
				WriteLine("buffer.extend(item_bytes)")
		}
	}
	builder.UnindentLine()
}

func (c *PythonGenerationContext) generateDeserializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	resultName := "model_result"
	strictBinaryParsing, err := c.UsesStrictBinaryParsing(CurrentLanguage, model)
	if err != nil {
		return err
	}
	binaryParseFallback := resultName
	if strictBinaryParsing {
		binaryParseFallback = "None"
	}

	builder.MapVarPairs(
		"model", model.Name,
		"result", resultName,
		"binaryParseFallback", binaryParseFallback,
	)
	defer builder.UnmapVar(
		"model",
		"result",
		"binaryParseFallback",
	)

	builder.WriteLine("@staticmethod").
		LineD("def deserialize_binary(data: bytes) -> $model | None:").
		Indent().
		// null-safety check
		WriteLine("if not data or len(data) == 0 or (len(data) == 1 and data[0] == 0):").
		Indent().
		WriteLine("return None").
		UnindentLine().
		WriteLine("buffer = memoryview(data)").
		WriteLine("offset = 0").
		NewLine()

	builder.LineD("$result = $model()").
		WriteLine("try:").
		Indent()

	for _, field := range model.Fields {
		if field.IsArray() {
			c.generateArrayDeserializeBinary(resultName, field, builder)
			continue
		}

		c.generateFieldDeserializeBinary(resultName, field, builder)
	}

	builder.Unindent().
		WriteLine("except (struct.error, ValueError):").
		Indent().
		LineD("return $binaryParseFallback").
		Unindent().
		LineD("return $result").
		UnindentLine().
		NewLine()
	return nil
}

func (c *PythonGenerationContext) generateFieldDeserializeBinary(
	resultName string,
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) {
	fieldName := caseUtils.ToSnakeCase(field.Name)
	resultField := resultName + "." + fieldName
	builder.MapVarPairs(
		"field", resultField,
		"len", fieldName+"_len",
		"rawField", fieldName,
		"rawFieldBytes", fieldName+"_bytes",
		"type", field.Type.GetName(),
	)
	defer builder.UnmapVar(
		"field",
		"len",
		"rawField",
		"rawFieldBytes",
		"type",
	)

	switch pythonStorageTypeName(field.Type) {
	case cclValues.TypeNameString:
		builder.LineD("$len = struct.unpack_from('<I', buffer, offset)[0]").
			WriteLine("offset += 4").
			LineD("if $len > len(buffer) - offset:").
			Indent().
			LineD("return $binaryParseFallback").
			Unindent().
			LineD(`$field = bytes(buffer[offset:offset+$len]).decode("utf-8")`).
			LineD("offset += $len")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.LineD("$field = " + c.pythonEnumCastExpression(field.Type, "struct.unpack_from('<i', buffer, offset)[0]")).
			WriteLine("offset += 4")
	case cclValues.TypeNameInt8:
		builder.LineD("$field = " + c.pythonEnumCastExpression(field.Type, "struct.unpack_from('<b', buffer, offset)[0]")).
			WriteLine("offset += 1")
	case cclValues.TypeNameInt16:
		builder.LineD("$field = " + c.pythonEnumCastExpression(field.Type, "struct.unpack_from('<h', buffer, offset)[0]")).
			WriteLine("offset += 2")
	case cclValues.TypeNameInt64:
		builder.LineD("$field = " + c.pythonEnumCastExpression(field.Type, "struct.unpack_from('<q', buffer, offset)[0]")).
			WriteLine("offset += 8")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.LineD("$field = " + c.pythonEnumCastExpression(field.Type, "struct.unpack_from('<I', buffer, offset)[0]")).
			WriteLine("offset += 4")
	case cclValues.TypeNameUint8:
		builder.LineD("$field = " + c.pythonEnumCastExpression(field.Type, "struct.unpack_from('<B', buffer, offset)[0]")).
			WriteLine("offset += 1")
	case cclValues.TypeNameUint16:
		builder.LineD("$field = " + c.pythonEnumCastExpression(field.Type, "struct.unpack_from('<H', buffer, offset)[0]")).
			WriteLine("offset += 2")
	case cclValues.TypeNameUint64:
		builder.LineD("$field = " + c.pythonEnumCastExpression(field.Type, "struct.unpack_from('<Q', buffer, offset)[0]")).
			WriteLine("offset += 8")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.LineD("$field = struct.unpack_from('<f', buffer, offset)[0]").
			WriteLine("offset += 4")
	case cclValues.TypeNameFloat64:
		builder.LineD("$field = struct.unpack_from('<d', buffer, offset)[0]").
			WriteLine("offset += 8")
	case cclValues.TypeNameBool:
		builder.LineD("$field = struct.unpack_from('<b', buffer, offset)[0] != 0").
			WriteLine("offset += 1")
	case cclValues.TypeNameBytes:
		builder.LineD("$len = struct.unpack_from('<I', buffer, offset)[0]").
			WriteLine("offset += 4").
			LineD("if $len > len(buffer) - offset:").
			Indent().
			LineD("return $binaryParseFallback").
			Unindent().
			LineD("$field = bytes(buffer[offset:offset+$len])").
			LineD("offset += $len")
	case cclValues.TypeNameDateTime:
		builder.LineD("$field = struct.unpack_from('<q', buffer, offset)[0]").
			WriteLine("offset += 8")
	default:
		// Custom type handling
		if field.IsCustomTypeModel() {
			builder.LineD("$len = struct.unpack_from('<I', buffer, offset)[0]").
				WriteLine("offset += 4").
				LineD("if $len > len(buffer) - offset:").
				Indent().
				LineD("return $binaryParseFallback").
				Unindent().
				LineD("$rawFieldBytes = bytes(buffer[offset:offset+$len])").
				LineD("offset += $len").
				LineD("$field = $type.deserialize_binary($rawFieldBytes)")
		}
	}
	builder.NewLine()
}

func (c *PythonGenerationContext) generateArrayDeserializeBinary(
	resultName string,
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldName := caseUtils.ToSnakeCase(field.Name)
	resultField := resultName + "." + fieldName
	builder.MapVarPairs(
		"field", resultField,
		"len", fieldName+"_len",
		"type", targetFieldType.GetName(),
	)
	defer builder.UnmapVar(
		"field",
		"len",
		"type",
	)

	builder.LineD("$len = struct.unpack_from('<I', buffer, offset)[0]").
		WriteLine("offset += 4").
		LineD("$field = []").
		LineD("for _ in range($len):").
		Indent()

	switch pythonStorageTypeName(targetFieldType) {
	case cclValues.TypeNameString:
		builder.WriteLine("item_len = struct.unpack_from('<I', buffer, offset)[0]").
			WriteLine("offset += 4").
			WriteLine("if item_len > len(buffer) - offset:").
			Indent().
			LineD("return $binaryParseFallback").
			Unindent().
			WriteLine("item = bytes(buffer[offset:offset+item_len]).decode(\"utf-8\")").
			WriteLine("offset += item_len").
			LineD("$field.append(item)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.LineD("$field.append(" + c.pythonEnumCastExpression(targetFieldType, "struct.unpack_from('<i', buffer, offset)[0]") + ")").
			WriteLine("offset += 4")
	case cclValues.TypeNameInt8:
		builder.LineD("$field.append(" + c.pythonEnumCastExpression(targetFieldType, "struct.unpack_from('<b', buffer, offset)[0]") + ")").
			WriteLine("offset += 1")
	case cclValues.TypeNameInt16:
		builder.LineD("$field.append(" + c.pythonEnumCastExpression(targetFieldType, "struct.unpack_from('<h', buffer, offset)[0]") + ")").
			WriteLine("offset += 2")
	case cclValues.TypeNameInt64:
		builder.LineD("$field.append(" + c.pythonEnumCastExpression(targetFieldType, "struct.unpack_from('<q', buffer, offset)[0]") + ")").
			WriteLine("offset += 8")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.LineD("$field.append(" + c.pythonEnumCastExpression(targetFieldType, "struct.unpack_from('<I', buffer, offset)[0]") + ")").
			WriteLine("offset += 4")
	case cclValues.TypeNameUint8:
		builder.LineD("$field.append(" + c.pythonEnumCastExpression(targetFieldType, "struct.unpack_from('<B', buffer, offset)[0]") + ")").
			WriteLine("offset += 1")
	case cclValues.TypeNameUint16:
		builder.LineD("$field.append(" + c.pythonEnumCastExpression(targetFieldType, "struct.unpack_from('<H', buffer, offset)[0]") + ")").
			WriteLine("offset += 2")
	case cclValues.TypeNameUint64:
		builder.LineD("$field.append(" + c.pythonEnumCastExpression(targetFieldType, "struct.unpack_from('<Q', buffer, offset)[0]") + ")").
			WriteLine("offset += 8")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.LineD("$field.append(struct.unpack_from('<f', buffer, offset)[0])").
			WriteLine("offset += 4")
	case cclValues.TypeNameFloat64:
		builder.LineD("$field.append(struct.unpack_from('<d', buffer, offset)[0])").
			WriteLine("offset += 8")
	case cclValues.TypeNameBool:
		builder.LineD("$field.append(struct.unpack_from('<b', buffer, offset)[0] != 0)").
			WriteLine("offset += 1")
	case cclValues.TypeNameBytes:
		builder.WriteLine("item_len = struct.unpack_from('<I', buffer, offset)[0]").
			WriteLine("offset += 4").
			WriteLine("if item_len > len(buffer) - offset:").
			Indent().
			LineD("return $binaryParseFallback").
			Unindent().
			LineD("$field.append(bytes(buffer[offset:offset+item_len]))").
			WriteLine("offset += item_len")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("item_len = struct.unpack_from('<I', buffer, offset)[0]").
				WriteLine("offset += 4").
				WriteLine("if item_len > len(buffer) - offset:").
				Indent().
				LineD("return $binaryParseFallback").
				Unindent().
				WriteLine("item_bytes = bytes(buffer[offset:offset+item_len])").
				WriteLine("offset += item_len").
				LineD("$field.append($type.deserialize_binary(item_bytes))")
		}
	}
	builder.UnindentLine()
}
