package gdGenerator

import (
	"strings"

	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func (c *GDScriptGenerationContext) generateSerializeBinaryMethod(model *CCLModel, builder *strings.Builder) error {
	builder.WriteString("func serialize_binary() -> PackedByteArray:\n")
	builder.WriteString("\tvar buffer = StreamPeerBuffer.new()\n\n")

	for i := range model.Fields {
		field := model.Fields[i]
		if field.IsArray() {
			c.generateArraySerializeBinary(field, builder)
			continue
		}

		c.generateFieldSerializeBinary(field, builder)
	}

	builder.WriteString("\treturn buffer.data_array\n\n")

	return nil
}

func (c *GDScriptGenerationContext) generateFieldSerializeBinary(field *CCLField, builder *strings.Builder) {
	fieldRawName := ToSnakeCase(field.Name)
	fieldName := "self." + fieldRawName
	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteString("\tvar " + fieldRawName + "_bytes = " + fieldName + ".to_utf8_buffer()\n")
		builder.WriteString("\tbuffer.put_u32(" + fieldRawName + "_bytes.size())\n")
		builder.WriteString("\tbuffer.put_data(" + fieldRawName + "_bytes)\n\n")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteString("\tbuffer.put_32(" + fieldName + ")\n\n")
	case cclValues.TypeNameInt8:
		builder.WriteString("\tbuffer.put_8(" + fieldName + ")\n\n")
	case cclValues.TypeNameInt16:
		builder.WriteString("\tbuffer.put_16(" + fieldName + ")\n\n")
	case cclValues.TypeNameInt64:
		builder.WriteString("\tbuffer.put_64(" + fieldName + ")\n\n")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteString("\tbuffer.put_u32(" + fieldName + ")\n\n")
	case cclValues.TypeNameUint8:
		builder.WriteString("\tbuffer.put_u8(" + fieldName + ")\n\n")
	case cclValues.TypeNameUint16:
		builder.WriteString("\tbuffer.put_u16(" + fieldName + ")\n\n")
	case cclValues.TypeNameUint64:
		builder.WriteString("\tbuffer.put_u64(" + fieldName + ")\n\n")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteString("\tbuffer.put_float(" + fieldName + ")\n\n")
	case cclValues.TypeNameBool:
		builder.WriteString("\tbuffer.put_8(1 if " + fieldName + " else 0)\n\n")
	case cclValues.TypeNameBytes:
		builder.WriteString("\tbuffer.put_u32(" + fieldName + ".size())\n")
		builder.WriteString("\tbuffer.put_data(" + fieldName + ")\n\n")
	case cclValues.TypeNameDateTime:
		builder.WriteString("\tbuffer.put_64(" + fieldName + ")\n\n")
	default:
		// Custom type handling
		if c.Options.CCLDefinition.IsCustomType(field.Type.GetName()) {
			builder.WriteString("\tvar " + fieldRawName + "_bytes = " +
				fieldName + ".serialize_binary() if " + fieldRawName + " else PackedByteArray([0])\n")
			builder.WriteString("\tbuffer.put_u32(" + fieldRawName + "_bytes.size())\n")
			builder.WriteString("\tbuffer.put_data(" + fieldRawName + "_bytes)\n\n")
		}
	}
}

func (c *GDScriptGenerationContext) generateArraySerializeBinary(field *CCLField, builder *strings.Builder) {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldName := "self." + ToSnakeCase(field.Name)
	builder.WriteString("\tbuffer.put_u32(" + fieldName + ".size())\n")
	builder.WriteString("\tfor item in " + fieldName + ":\n")

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteString("\t\tvar item_bytes = item.to_utf8_buffer()\n")
		builder.WriteString("\t\tbuffer.put_u32(item_bytes.size())\n")
		builder.WriteString("\t\tbuffer.put_data(item_bytes)\n")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteString("\t\tbuffer.put_32(item)\n")
	case cclValues.TypeNameInt8:
		builder.WriteString("\t\tbuffer.put_8(item)\n")
	case cclValues.TypeNameInt16:
		builder.WriteString("\t\tbuffer.put_16(item)\n")
	case cclValues.TypeNameInt64:
		builder.WriteString("\t\tbuffer.put_64(item)\n")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteString("\t\tbuffer.put_u32(item)\n")
	case cclValues.TypeNameUint8:
		builder.WriteString("\t\tbuffer.put_u8(item)\n")
	case cclValues.TypeNameUint16:
		builder.WriteString("\t\tbuffer.put_u16(item)\n")
	case cclValues.TypeNameUint64:
		builder.WriteString("\t\tbuffer.put_u64(item)\n")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteString("\t\tbuffer.put_float(item)\n")
	case cclValues.TypeNameBool:
		builder.WriteString("\t\tbuffer.put_8(1 if item else 0)\n")
	default:
		if c.Options.CCLDefinition.IsCustomType(targetFieldType.GetName()) {
			builder.WriteString("\t\tvar item_bytes = item.serialize_binary() if item else PackedByteArray([0])\n")
			builder.WriteString("\t\tbuffer.put_u32(item_bytes.size())\n")
			builder.WriteString("\t\tbuffer.put_data(item_bytes)\n")
		}
	}
	builder.WriteString("\n")
}

func (c *GDScriptGenerationContext) generateDeserializeBinaryMethod(model *CCLModel, builder *strings.Builder) error {
	builder.WriteString("static func deserialize_binary(data: PackedByteArray) -> " + model.Name + ":\n")

	// null-safety check
	builder.WriteString("\tif not data or data.is_empty() or (data.size() == 1 and data[0] == 0):\n")
	builder.WriteString("\t\treturn null\n\n")

	builder.WriteString("\tvar buffer = StreamPeerBuffer.new()\n")
	builder.WriteString("\tbuffer.data_array = data\n")

	// TODO: maybe dynamically generate this later?
	resultName := "model_result"
	builder.WriteString("\tvar " + resultName + " = " + model.Name + ".new()\n\n")

	for _, field := range model.Fields {
		if field.IsArray() {
			c.generateArrayDeserializeBinary(resultName, field, builder)
			continue
		}

		c.generateFieldDeserializeBinary(resultName, field, builder)
	}

	builder.WriteString("\treturn " + resultName + "\n")
	return nil
}

func (c *GDScriptGenerationContext) generateFieldDeserializeBinary(
	resultName string,
	field *CCLField,
	builder *strings.Builder,
) {
	fieldName := ToSnakeCase(field.Name)
	resultField := resultName + "." + fieldName

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteString("\tvar " + fieldName + "_len = buffer.get_u32()\n")
		builder.WriteString("\tif " + fieldName + "_len > buffer.get_size() - buffer.get_position():\n")
		builder.WriteString("\t\treturn null\n")
		builder.WriteString("\t" + resultField + " = buffer.get_data(" +
			fieldName + "_len)[1].get_string_from_utf8()\n\n")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteString("\t" + resultField + " = buffer.get_32()\n\n")
	case cclValues.TypeNameInt8:
		builder.WriteString("\t" + resultField + " = buffer.get_8()\n\n")
	case cclValues.TypeNameInt16:
		builder.WriteString("\t" + resultField + " = buffer.get_16()\n\n")
	case cclValues.TypeNameInt64:
		builder.WriteString("\t" + resultField + " = buffer.get_64()\n\n")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteString("\t" + resultField + " = buffer.get_u32()\n\n")
	case cclValues.TypeNameUint8:
		builder.WriteString("\t" + resultField + " = buffer.get_u8()\n\n")
	case cclValues.TypeNameUint16:
		builder.WriteString("\t" + resultField + " = buffer.get_u16()\n\n")
	case cclValues.TypeNameUint64:
		builder.WriteString("\t" + resultField + " = buffer.get_u64()\n\n")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteString("\t" + resultField + " = buffer.get_float()\n\n")
	case cclValues.TypeNameBool:
		builder.WriteString("\t" + resultField + " = buffer.get_8() != 0\n\n")
	case cclValues.TypeNameBytes:
		builder.WriteString("\tvar " + fieldName + "_len = buffer.get_u32()\n")
		builder.WriteString("\t" + resultField + " = buffer.get_data(" +
			fieldName + "_len)[1]\n\n")
	case cclValues.TypeNameDateTime:
		builder.WriteString("\t" + resultField + " = buffer.get_64()\n\n")
	default:
		// Custom type handling
		if c.Options.CCLDefinition.IsCustomType(field.Type.GetName()) {
			builder.WriteString("\tvar " + fieldName + "_len = buffer.get_u32()\n")
			builder.WriteString("\tvar " + fieldName + "_bytes = buffer.get_data(" +
				fieldName + "_len)[1]\n")
			builder.WriteString("\t" + resultField + " = " + field.Type.GetName() +
				".deserialize_binary(" + fieldName + "_bytes)\n\n")
		}
	}
}

func (c *GDScriptGenerationContext) generateArrayDeserializeBinary(
	resultName string,
	field *CCLField,
	builder *strings.Builder,
) {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldName := ToSnakeCase(field.Name)
	resultField := resultName + "." + fieldName

	builder.WriteString("\tvar " + fieldName + "_len = buffer.get_u32()\n")
	builder.WriteString("\t" + resultField + " = []\n")
	builder.WriteString("\tfor i in range(" + fieldName + "_len):\n")

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteString("\t\tvar item_len = buffer.get_u32()\n")
		builder.WriteString("\t\tif item_len > buffer.get_size() - buffer.get_position():\n")
		builder.WriteString("\t\t\treturn null\n")
		builder.WriteString("\t\tvar item = buffer.get_data(item_len)[1].get_string_from_utf8()\n")
		builder.WriteString("\t\t" + resultField + ".append(item)\n")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteString("\t\t" + resultField + ".append(buffer.get_32())\n")
	case cclValues.TypeNameInt8:
		builder.WriteString("\t\t" + resultField + ".append(buffer.get_8())\n")
	case cclValues.TypeNameInt16:
		builder.WriteString("\t\t" + resultField + ".append(buffer.get_16())\n")
	case cclValues.TypeNameInt64:
		builder.WriteString("\t\t" + resultField + ".append(buffer.get_64())\n")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteString("\t\t" + resultField + ".append(buffer.get_u32())\n")
	case cclValues.TypeNameUint8:
		builder.WriteString("\t\t" + resultField + ".append(buffer.get_u8())\n")
	case cclValues.TypeNameUint16:
		builder.WriteString("\t\t" + resultField + ".append(buffer.get_u16())\n")
	case cclValues.TypeNameUint64:
		builder.WriteString("\t\t" + resultField + ".append(buffer.get_u64())\n")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteString("\t\t" + resultField + ".append(buffer.get_float())\n")
	case cclValues.TypeNameBool:
		builder.WriteString("\t\t" + resultField + ".append(buffer.get_8() != 0)\n")
	default:
		if c.Options.CCLDefinition.IsCustomType(targetFieldType.GetName()) {
			builder.WriteString("\t\tvar item_len = buffer.get_u32()\n")
			builder.WriteString("\t\tif item_len > buffer.get_size() - buffer.get_position():\n")
			builder.WriteString("\t\t\treturn null\n")
			builder.WriteString("\t\tvar item_bytes = buffer.get_data(item_len)[1]\n")
			builder.WriteString("\t\t" + resultField + ".append(" +
				targetFieldType.GetName() + ".deserialize_binary(item_bytes))\n")
		}
	}
	builder.WriteString("\n")
}
