package pyGenerator

import (
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
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
	fieldRawName := toSnakeCase(field.Name)
	fieldName := "self." + fieldRawName
	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine(fieldRawName + "_bytes = " + fieldName + ".encode('utf-8')").
			WriteLine("buffer.extend(struct.pack('<I', len(" + fieldRawName + "_bytes)))").
			WriteLine("buffer.extend(" + fieldRawName + "_bytes)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine("buffer.extend(struct.pack('<i', " + fieldName + "))")
	case cclValues.TypeNameInt8:
		builder.WriteLine("buffer.extend(struct.pack('<b', " + fieldName + "))")
	case cclValues.TypeNameInt16:
		builder.WriteLine("buffer.extend(struct.pack('<h', " + fieldName + "))")
	case cclValues.TypeNameInt64:
		builder.WriteLine("buffer.extend(struct.pack('<q', " + fieldName + "))")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine("buffer.extend(struct.pack('<I', " + fieldName + "))")
	case cclValues.TypeNameUint8:
		builder.WriteLine("buffer.extend(struct.pack('<B', " + fieldName + "))")
	case cclValues.TypeNameUint16:
		builder.WriteLine("buffer.extend(struct.pack('<H', " + fieldName + "))")
	case cclValues.TypeNameUint64:
		builder.WriteLine("buffer.extend(struct.pack('<Q', " + fieldName + "))")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.WriteLine("buffer.extend(struct.pack('<f', " + fieldName + "))")
	case cclValues.TypeNameFloat64:
		builder.WriteLine("buffer.extend(struct.pack('<d', " + fieldName + "))")
	case cclValues.TypeNameBool:
		builder.WriteLine("buffer.extend(struct.pack('<b', 1 if " + fieldName + " else 0))")
	case cclValues.TypeNameBytes:
		builder.WriteLine("buffer.extend(struct.pack('<I', len(" + fieldName + ")))").
			WriteLine("buffer.extend(" + fieldName + ")")
	case cclValues.TypeNameDateTime:
		builder.WriteLine("buffer.extend(struct.pack('<q', " + fieldName + "))")
	default:
		// Custom type handling
		if c.Options.CCLDefinition.IsCustomType(field.Type.GetName()) {
			builder.WriteLine("if " + fieldName + ":").
				Indent().
				WriteLine(fieldRawName + "_bytes = " + fieldName + ".serialize_binary()").
				Unindent().
				WriteLine("else:").
				Indent().
				WriteLine(fieldRawName + "_bytes = b'\\x00'").
				Unindent().
				WriteLine("buffer.extend(struct.pack('<I', len(" + fieldRawName + "_bytes)))").
				WriteLine("buffer.extend(" + fieldRawName + "_bytes)")
		}
	}
	builder.NewLine()
}

func (c *PythonGenerationContext) generateArraySerializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldName := "self." + toSnakeCase(field.Name)
	builder.WriteLine("buffer.extend(struct.pack('<I', len(" + fieldName + ")))").
		WriteLine("for item in " + fieldName + ":").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("item_bytes = item.encode('utf-8')").
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
	default:
		if c.Options.CCLDefinition.IsCustomType(targetFieldType.GetName()) {
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
	builder.WriteLine("@staticmethod").
		WriteLine("def deserialize_binary(data: bytes) -> " + model.Name + " | None:").
		Indent().
		// null-safety check
		WriteLine("if not data or len(data) == 0 or (len(data) == 1 and data[0] == 0):").
		Indent().
		WriteLine("return None").
		UnindentLine().
		WriteLine("buffer = memoryview(data)").
		WriteLine("offset = 0").
		NewLine()

	resultName := "model_result"
	builder.WriteLine(resultName + " = " + model.Name + "()").
		NewLine()

	for _, field := range model.Fields {
		if field.IsArray() {
			c.generateArrayDeserializeBinary(resultName, field, builder)
			continue
		}

		c.generateFieldDeserializeBinary(resultName, field, builder)
	}

	builder.WriteLine("return " + resultName).
		UnindentLine().
		NewLine()
	return nil
}

func (c *PythonGenerationContext) generateFieldDeserializeBinary(
	resultName string,
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) {
	fieldName := toSnakeCase(field.Name)
	resultField := resultName + "." + fieldName

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine(fieldName + "_len = struct.unpack_from('<I', buffer, offset)[0]").
			WriteLine("offset += 4").
			WriteLine("if " + fieldName + "_len > len(buffer) - offset:").
			Indent().
			WriteLine("return None").
			Unindent().
			WriteLine(resultField + " = bytes(buffer[offset:offset+" + fieldName + "_len]).decode('utf-8')").
			WriteLine("offset += " + fieldName + "_len")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine(resultField + " = struct.unpack_from('<i', buffer, offset)[0]").
			WriteLine("offset += 4")
	case cclValues.TypeNameInt8:
		builder.WriteLine(resultField + " = struct.unpack_from('<b', buffer, offset)[0]").
			WriteLine("offset += 1")
	case cclValues.TypeNameInt16:
		builder.WriteLine(resultField + " = struct.unpack_from('<h', buffer, offset)[0]").
			WriteLine("offset += 2")
	case cclValues.TypeNameInt64:
		builder.WriteLine(resultField + " = struct.unpack_from('<q', buffer, offset)[0]").
			WriteLine("offset += 8")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine(resultField + " = struct.unpack_from('<I', buffer, offset)[0]").
			WriteLine("offset += 4")
	case cclValues.TypeNameUint8:
		builder.WriteLine(resultField + " = struct.unpack_from('<B', buffer, offset)[0]").
			WriteLine("offset += 1")
	case cclValues.TypeNameUint16:
		builder.WriteLine(resultField + " = struct.unpack_from('<H', buffer, offset)[0]").
			WriteLine("offset += 2")
	case cclValues.TypeNameUint64:
		builder.WriteLine(resultField + " = struct.unpack_from('<Q', buffer, offset)[0]").
			WriteLine("offset += 8")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.WriteLine(resultField + " = struct.unpack_from('<f', buffer, offset)[0]").
			WriteLine("offset += 4")
	case cclValues.TypeNameFloat64:
		builder.WriteLine(resultField + " = struct.unpack_from('<d', buffer, offset)[0]").
			WriteLine("offset += 8")
	case cclValues.TypeNameBool:
		builder.WriteLine(resultField + " = struct.unpack_from('<b', buffer, offset)[0] != 0").
			WriteLine("offset += 1")
	case cclValues.TypeNameBytes:
		builder.WriteLine(fieldName + "_len = struct.unpack_from('<I', buffer, offset)[0]").
			WriteLine("offset += 4").
			WriteLine(resultField + " = bytes(buffer[offset:offset+" + fieldName + "_len])").
			WriteLine("offset += " + fieldName + "_len")
	case cclValues.TypeNameDateTime:
		builder.WriteLine(resultField + " = struct.unpack_from('<q', buffer, offset)[0]").
			WriteLine("offset += 8")
	default:
		// Custom type handling
		if c.Options.CCLDefinition.IsCustomType(field.Type.GetName()) {
			builder.WriteLine(fieldName + "_len = struct.unpack_from('<I', buffer, offset)[0]").
				WriteLine("offset += 4").
				WriteLine(fieldName + "_bytes = bytes(buffer[offset:offset+" + fieldName + "_len])").
				WriteLine("offset += " + fieldName + "_len").
				WriteLine(resultField + " = " + field.Type.GetName() + ".deserialize_binary(" + fieldName + "_bytes)")
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
	fieldName := toSnakeCase(field.Name)
	resultField := resultName + "." + fieldName

	builder.WriteLine(fieldName + "_len = struct.unpack_from('<I', buffer, offset)[0]").
		WriteLine("offset += 4").
		WriteLine(resultField + " = []").
		WriteLine("for _ in range(" + fieldName + "_len):").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("item_len = struct.unpack_from('<I', buffer, offset)[0]").
			WriteLine("offset += 4").
			WriteLine("if item_len > len(buffer) - offset:").
			Indent().
			WriteLine("return None").
			Unindent().
			WriteLine("item = bytes(buffer[offset:offset+item_len]).decode('utf-8')").
			WriteLine("offset += item_len").
			WriteLine(resultField + ".append(item)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine(resultField + ".append(struct.unpack_from('<i', buffer, offset)[0])").
			WriteLine("offset += 4")
	case cclValues.TypeNameInt8:
		builder.WriteLine(resultField + ".append(struct.unpack_from('<b', buffer, offset)[0])").
			WriteLine("offset += 1")
	case cclValues.TypeNameInt16:
		builder.WriteLine(resultField + ".append(struct.unpack_from('<h', buffer, offset)[0])").
			WriteLine("offset += 2")
	case cclValues.TypeNameInt64:
		builder.WriteLine(resultField + ".append(struct.unpack_from('<q', buffer, offset)[0])").
			WriteLine("offset += 8")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine(resultField + ".append(struct.unpack_from('<I', buffer, offset)[0])").
			WriteLine("offset += 4")
	case cclValues.TypeNameUint8:
		builder.WriteLine(resultField + ".append(struct.unpack_from('<B', buffer, offset)[0])").
			WriteLine("offset += 1")
	case cclValues.TypeNameUint16:
		builder.WriteLine(resultField + ".append(struct.unpack_from('<H', buffer, offset)[0])").
			WriteLine("offset += 2")
	case cclValues.TypeNameUint64:
		builder.WriteLine(resultField + ".append(struct.unpack_from('<Q', buffer, offset)[0])").
			WriteLine("offset += 8")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.WriteLine(resultField + ".append(struct.unpack_from('<f', buffer, offset)[0])").
			WriteLine("offset += 4")
	case cclValues.TypeNameFloat64:
		builder.WriteLine(resultField + ".append(struct.unpack_from('<d', buffer, offset)[0])").
			WriteLine("offset += 8")
	case cclValues.TypeNameBool:
		builder.WriteLine(resultField + ".append(struct.unpack_from('<b', buffer, offset)[0] != 0)").
			WriteLine("offset += 1")
	default:
		if c.Options.CCLDefinition.IsCustomType(targetFieldType.GetName()) {
			builder.WriteLine("item_len = struct.unpack_from('<I', buffer, offset)[0]").
				WriteLine("offset += 4").
				WriteLine("if item_len > len(buffer) - offset:").
				Indent().
				WriteLine("return None").
				Unindent().
				WriteLine("item_bytes = bytes(buffer[offset:offset+item_len])").
				WriteLine("offset += item_len").
				WriteLine(resultField + ".append(" + targetFieldType.GetName() + ".deserialize_binary(item_bytes))")
		}
	}
	builder.UnindentLine()
}
