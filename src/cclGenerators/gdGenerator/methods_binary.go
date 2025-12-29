package gdGenerator

import (
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func (c *GDScriptGenerationContext) generateSerializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	builder.WriteLine("func serialize_binary() -> PackedByteArray:")
	builder.Indent()
	builder.WriteLine("var buffer = StreamPeerBuffer.new()")
	builder.NewLine()

	for i := range model.Fields {
		field := model.Fields[i]
		if field.IsArray() {
			c.generateArraySerializeBinary(field, builder)
			continue
		}

		c.generateFieldSerializeBinary(field, builder)
	}

	builder.WriteLine("return buffer.data_array")
	builder.Unindent()
	builder.NewLine()

	return nil
}

func (c *GDScriptGenerationContext) generateFieldSerializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	fieldRawName := ToSnakeCase(field.Name)
	fieldName := "self." + fieldRawName
	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("var " + fieldRawName + "_bytes = " + fieldName + ".to_utf8_buffer()")
		builder.WriteLine("buffer.put_u32(" + fieldRawName + "_bytes.size())")
		builder.WriteLine("buffer.put_data(" + fieldRawName + "_bytes)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine("buffer.put_32(" + fieldName + ")")
	case cclValues.TypeNameInt8:
		builder.WriteLine("buffer.put_8(" + fieldName + ")")
	case cclValues.TypeNameInt16:
		builder.WriteLine("buffer.put_16(" + fieldName + ")")
	case cclValues.TypeNameInt64:
		builder.WriteLine("buffer.put_64(" + fieldName + ")")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine("buffer.put_u32(" + fieldName + ")")
	case cclValues.TypeNameUint8:
		builder.WriteLine("buffer.put_u8(" + fieldName + ")")
	case cclValues.TypeNameUint16:
		builder.WriteLine("buffer.put_u16(" + fieldName + ")")
	case cclValues.TypeNameUint64:
		builder.WriteLine("buffer.put_u64(" + fieldName + ")")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteLine("buffer.put_float(" + fieldName + ")")
	case cclValues.TypeNameBool:
		builder.WriteLine("buffer.put_8(1 if " + fieldName + " else 0)")
	case cclValues.TypeNameBytes:
		builder.WriteLine("buffer.put_u32(" + fieldName + ".size())")
		builder.WriteLine("buffer.put_data(" + fieldName + ")")
	case cclValues.TypeNameDateTime:
		builder.WriteLine("buffer.put_64(" + fieldName + ")")
	default:
		// Custom type handling
		if c.Options.CCLDefinition.IsCustomType(field.Type.GetName()) {
			builder.WriteLine("var " + fieldRawName + "_bytes = " +
				fieldName + ".serialize_binary() if " + fieldRawName + " else PackedByteArray([0])")
			builder.WriteLine("buffer.put_u32(" + fieldRawName + "_bytes.size())")
			builder.WriteLine("buffer.put_data(" + fieldRawName + "_bytes)")
		}
	}
	builder.NewLine()
}

func (c *GDScriptGenerationContext) generateArraySerializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldName := "self." + ToSnakeCase(field.Name)
	builder.WriteLine("buffer.put_u32(" + fieldName + ".size())")
	builder.WriteLine("for item in " + fieldName + ":")
	builder.Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("var item_bytes = item.to_utf8_buffer()")
		builder.WriteLine("buffer.put_u32(item_bytes.size())")
		builder.WriteLine("buffer.put_data(item_bytes)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine("buffer.put_32(item)")
	case cclValues.TypeNameInt8:
		builder.WriteLine("buffer.put_8(item)")
	case cclValues.TypeNameInt16:
		builder.WriteLine("buffer.put_16(item)")
	case cclValues.TypeNameInt64:
		builder.WriteLine("buffer.put_64(item)")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine("buffer.put_u32(item)")
	case cclValues.TypeNameUint8:
		builder.WriteLine("buffer.put_u8(item)")
	case cclValues.TypeNameUint16:
		builder.WriteLine("buffer.put_u16(item)")
	case cclValues.TypeNameUint64:
		builder.WriteLine("buffer.put_u64(item)")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteLine("buffer.put_float(item)")
	case cclValues.TypeNameBool:
		builder.WriteLine("buffer.put_8(1 if item else 0)")
	default:
		if c.Options.CCLDefinition.IsCustomType(targetFieldType.GetName()) {
			builder.WriteLine("var item_bytes = item.serialize_binary() if item else PackedByteArray([0])")
			builder.WriteLine("buffer.put_u32(item_bytes.size())")
			builder.WriteLine("buffer.put_data(item_bytes)")
		}
	}
	builder.Unindent()
	builder.NewLine()
}

func (c *GDScriptGenerationContext) generateDeserializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	builder.WriteLine("static func deserialize_binary(data: PackedByteArray) -> " + model.Name + ":")
	builder.Indent()

	// null-safety check
	builder.WriteLine("if not data or data.is_empty() or (data.size() == 1 and data[0] == 0):")
	builder.Indent()
	builder.WriteLine("return null")
	builder.Unindent()
	builder.NewLine()

	builder.WriteLine("var buffer = StreamPeerBuffer.new()")
	builder.WriteLine("buffer.data_array = data")

	// TODO: maybe dynamically generate this later?
	resultName := "model_result"
	builder.WriteLine("var " + resultName + " = " + model.Name + ".new()")
	builder.NewLine()

	for _, field := range model.Fields {
		if field.IsArray() {
			c.generateArrayDeserializeBinary(resultName, field, builder)
			continue
		}

		c.generateFieldDeserializeBinary(resultName, field, builder)
	}

	builder.WriteLine("return " + resultName)
	builder.Unindent()
	builder.NewLine()
	return nil
}

func (c *GDScriptGenerationContext) generateFieldDeserializeBinary(
	resultName string,
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) {
	fieldName := ToSnakeCase(field.Name)
	resultField := resultName + "." + fieldName

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("var " + fieldName + "_len = buffer.get_u32()")
		builder.WriteLine("if " + fieldName + "_len > buffer.get_size() - buffer.get_position():")
		builder.Indent()
		builder.WriteLine("return null")
		builder.Unindent()
		builder.WriteLine(resultField + " = buffer.get_data(" +
			fieldName + "_len)[1].get_string_from_utf8()")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine(resultField + " = buffer.get_32()")
	case cclValues.TypeNameInt8:
		builder.WriteLine(resultField + " = buffer.get_8()")
	case cclValues.TypeNameInt16:
		builder.WriteLine(resultField + " = buffer.get_16()")
	case cclValues.TypeNameInt64:
		builder.WriteLine(resultField + " = buffer.get_64()")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine(resultField + " = buffer.get_u32()")
	case cclValues.TypeNameUint8:
		builder.WriteLine(resultField + " = buffer.get_u8()")
	case cclValues.TypeNameUint16:
		builder.WriteLine(resultField + " = buffer.get_u16()")
	case cclValues.TypeNameUint64:
		builder.WriteLine(resultField + " = buffer.get_u64()")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteLine(resultField + " = buffer.get_float()")
	case cclValues.TypeNameBool:
		builder.WriteLine(resultField + " = buffer.get_8() != 0")
	case cclValues.TypeNameBytes:
		builder.WriteLine("var " + fieldName + "_len = buffer.get_u32()")
		builder.WriteLine(resultField + " = buffer.get_data(" +
			fieldName + "_len)[1]")
	case cclValues.TypeNameDateTime:
		builder.WriteLine(resultField + " = buffer.get_64()")
	default:
		// Custom type handling
		if c.Options.CCLDefinition.IsCustomType(field.Type.GetName()) {
			builder.WriteLine("var " + fieldName + "_len = buffer.get_u32()")
			builder.WriteLine("var " + fieldName + "_bytes = buffer.get_data(" +
				fieldName + "_len)[1]")
			builder.WriteLine(resultField + " = " + field.Type.GetName() +
				".deserialize_binary(" + fieldName + "_bytes)")
		}
	}
	builder.NewLine()
}

func (c *GDScriptGenerationContext) generateArrayDeserializeBinary(
	resultName string,
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldName := ToSnakeCase(field.Name)
	resultField := resultName + "." + fieldName

	builder.WriteLine("var " + fieldName + "_len = buffer.get_u32()")
	builder.WriteLine(resultField + " = []" + " as " + c.getGDScriptType(field))
	builder.WriteLine("for i in range(" + fieldName + "_len):")
	builder.Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("var item_len = buffer.get_u32()")
		builder.WriteLine("if item_len > buffer.get_size() - buffer.get_position():")
		builder.Indent()
		builder.WriteLine("return null")
		builder.Unindent()
		builder.WriteLine("var item = buffer.get_data(item_len)[1].get_string_from_utf8()")
		builder.WriteLine(resultField + ".append(item)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine(resultField + ".append(buffer.get_32())")
	case cclValues.TypeNameInt8:
		builder.WriteLine(resultField + ".append(buffer.get_8())")
	case cclValues.TypeNameInt16:
		builder.WriteLine(resultField + ".append(buffer.get_16())")
	case cclValues.TypeNameInt64:
		builder.WriteLine(resultField + ".append(buffer.get_64())")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine(resultField + ".append(buffer.get_u32())")
	case cclValues.TypeNameUint8:
		builder.WriteLine(resultField + ".append(buffer.get_u8())")
	case cclValues.TypeNameUint16:
		builder.WriteLine(resultField + ".append(buffer.get_u16())")
	case cclValues.TypeNameUint64:
		builder.WriteLine(resultField + ".append(buffer.get_u64())")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteLine(resultField + ".append(buffer.get_float())")
	case cclValues.TypeNameBool:
		builder.WriteLine(resultField + ".append(buffer.get_8() != 0)")
	default:
		if c.Options.CCLDefinition.IsCustomType(targetFieldType.GetName()) {
			builder.WriteLine("var item_len = buffer.get_u32()")
			builder.WriteLine("if item_len > buffer.get_size() - buffer.get_position():")
			builder.Indent()
			builder.WriteLine("return null")
			builder.Unindent()
			builder.WriteLine("var item_bytes = buffer.get_data(item_len)[1]")
			builder.WriteLine(resultField + ".append(" +
				targetFieldType.GetName() + ".deserialize_binary(item_bytes))")
		}
	}
	builder.Unindent()
	builder.NewLine()
}
