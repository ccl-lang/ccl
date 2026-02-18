package gdGenerator

import (
	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclUtils"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func (c *GDScriptGenerationContext) generateSerializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	builder.WriteLine("func serialize_binary() -> PackedByteArray:").
		Indent().
		WriteLine("var buffer = StreamPeerBuffer.new()").
		NewLine()

	for i := range model.Fields {
		field := model.Fields[i]
		if field.IsArray() {
			err := c.generateArraySerializeBinary(field, builder)
			if err != nil {
				return err
			}

			continue
		}

		err := c.generateFieldSerializeBinary(field, builder)
		if err != nil {
			return err
		}
	}

	builder.WriteLine("return buffer.data_array").
		UnindentLine()

	return nil
}

func (c *GDScriptGenerationContext) generateFieldSerializeBinary(
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldRawName := cclUtils.ToSnakeCase(field.Name)
	resultField := "self." + fieldRawName
	modelName := field.OwnedBy.GetName()
	targetFieldTypeName := field.Type.GetName()

	switch targetFieldTypeName {
	case cclValues.TypeNameString:
		builder.WriteLine("var " + fieldRawName + "_bytes = " + resultField + ".to_utf8_buffer()").
			WriteLine("buffer.put_u32(" + fieldRawName + "_bytes.size())").
			WriteLine("buffer.put_data(" + fieldRawName + "_bytes)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine("buffer.put_32(" + resultField + ")")
	case cclValues.TypeNameInt8:
		builder.WriteLine("buffer.put_8(" + resultField + ")")
	case cclValues.TypeNameInt16:
		builder.WriteLine("buffer.put_16(" + resultField + ")")
	case cclValues.TypeNameInt64:
		builder.WriteLine("buffer.put_64(" + resultField + ")")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine("buffer.put_u32(" + resultField + ")")
	case cclValues.TypeNameUint8:
		builder.WriteLine("buffer.put_u8(" + resultField + ")")
	case cclValues.TypeNameUint16:
		builder.WriteLine("buffer.put_u16(" + resultField + ")")
	case cclValues.TypeNameUint64:
		builder.WriteLine("buffer.put_u64(" + resultField + ")")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.WriteLine("buffer.put_float(" + resultField + ")")
	case cclValues.TypeNameBool:
		builder.WriteLine("buffer.put_8(1 if " + resultField + " else 0)")
	case cclValues.TypeNameBytes:
		builder.WriteLine("buffer.put_u32(" + resultField + ".size())").
			WriteLine("buffer.put_data(" + resultField + ")")
	case cclValues.TypeNameDateTime:
		builder.WriteLine("buffer.put_64(" + resultField + ")")
	default:
		// Custom type handling
		if field.Type.IsCustomTypeModel() {
			builder.WriteLine("var " + fieldRawName + "_bytes = " +
				resultField + ".serialize_binary() if " + fieldRawName + " else PackedByteArray([0])").
				WriteLine("buffer.put_u32(" + fieldRawName + "_bytes.size())").
				WriteLine("buffer.put_data(" + fieldRawName + "_bytes)")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldTypeName,
				ModelName:      modelName,
				FieldName:      fieldRawName,
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}
	builder.NewLine()

	return nil
}

func (c *GDScriptGenerationContext) generateArraySerializeBinary(
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	targetFieldTypeName := targetFieldType.GetName()
	modelName := field.OwnedBy.GetName()
	fieldRawName := cclUtils.ToSnakeCase(field.GetName())
	resultField := "self." + fieldRawName

	builder.WriteLine("buffer.put_u32(" + resultField + ".size())").
		WriteLine("for item in " + resultField + ":").
		Indent()

	switch targetFieldTypeName {
	case cclValues.TypeNameString:
		builder.WriteLine("var item_bytes = item.to_utf8_buffer()").
			WriteLine("buffer.put_u32(item_bytes.size())").
			WriteLine("buffer.put_data(item_bytes)")
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
	case cclValues.TypeNameBytes:
		builder.WriteLine("buffer.put_u32(item.size())").
			WriteLine("buffer.put_data(item)")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("var item_bytes = item.serialize_binary() if item else PackedByteArray([0])").
				WriteLine("buffer.put_u32(item_bytes.size())").
				WriteLine("buffer.put_data(item_bytes)")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldTypeName,
				ModelName:      modelName,
				FieldName:      fieldRawName,
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}

	builder.UnindentLine()
	return nil
}

func (c *GDScriptGenerationContext) generateDeserializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	builder.WriteLine("static func deserialize_binary(data: PackedByteArray) -> " + model.Name + ":").
		Indent().
		// null-safety check
		WriteLine("if not data or data.is_empty() or (data.size() == 1 and data[0] == 0):").
		Indent().
		WriteLine("return null").
		UnindentLine().
		WriteLine("var buffer = StreamPeerBuffer.new()").
		WriteLine("buffer.data_array = data")

	// TODO: maybe dynamically generate this later?
	resultName := "model_result"
	builder.WriteLine("var " + resultName + " = " + model.Name + ".new()").
		NewLine()

	for _, field := range model.Fields {
		if field.IsArray() {
			err := c.generateArrayDeserializeBinary(resultName, field, builder)
			if err != nil {
				return err
			}

			continue
		}

		err := c.generateFieldDeserializeBinary(resultName, field, builder)
		if err != nil {
			return err
		}
	}

	builder.WriteLine("return " + resultName).
		UnindentLine()

	return nil
}

func (c *GDScriptGenerationContext) generateFieldDeserializeBinary(
	resultName string,
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldRawName := cclUtils.ToSnakeCase(field.GetName())
	resultField := resultName + "." + fieldRawName
	targetFieldTypeName := field.Type.GetName()
	modelName := field.OwnedBy.GetName()

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("var " + fieldRawName + "_len = buffer.get_u32()").
			WriteLine("if " + fieldRawName + "_len > buffer.get_size() - buffer.get_position():").
			Indent().
			WriteLine("return null").
			Unindent().
			WriteLine(resultField + " = buffer.get_data(" +
				fieldRawName + "_len)[1].get_string_from_utf8()")
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
		builder.WriteLine("var " + fieldRawName + "_len = buffer.get_u32()").
			WriteLine(resultField + " = buffer.get_data(" +
				fieldRawName + "_len)[1]")
	case cclValues.TypeNameDateTime:
		builder.WriteLine(resultField + " = buffer.get_64()")
	default:
		// Custom type handling
		if field.Type.IsCustomTypeModel() {
			builder.WriteLine("var " + fieldRawName + "_len = buffer.get_u32()").
				WriteLine("var " + fieldRawName + "_bytes = buffer.get_data(" +
					fieldRawName + "_len)[1]").
				WriteLine(resultField + " = " + field.Type.GetName() +
					".deserialize_binary(" + fieldRawName + "_bytes)")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldTypeName,
				ModelName:      modelName,
				FieldName:      fieldRawName,
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}
	builder.NewLine()

	return nil
}

func (c *GDScriptGenerationContext) generateArrayDeserializeBinary(
	resultName string,
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldRawName := cclUtils.ToSnakeCase(field.GetName())
	resultField := resultName + "." + fieldRawName
	targetFieldTypeName := targetFieldType.GetName()
	modelName := field.OwnedBy.GetName()

	builder.WriteLine("var " + fieldRawName + "_len = buffer.get_u32()").
		WriteLine(resultField + " = []" + " as " + c.getGDScriptType(field)).
		WriteLine("for i in range(" + fieldRawName + "_len):").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("var item_len = buffer.get_u32()").
			WriteLine("if item_len > buffer.get_size() - buffer.get_position():").
			Indent().
			WriteLine("return null").
			Unindent().
			WriteLine("var item = buffer.get_data(item_len)[1].get_string_from_utf8()").
			WriteLine(resultField + ".append(item)")
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
	case cclValues.TypeNameBytes:
		builder.WriteLine("var item_len = buffer.get_u32()").
			WriteLine("if item_len > buffer.get_size() - buffer.get_position():").
			Indent().
			WriteLine("return null").
			Unindent().
			WriteLine(resultField + ".append(buffer.get_data(item_len)[1])")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("var item_len = buffer.get_u32()").
				WriteLine("if item_len > buffer.get_size() - buffer.get_position():").
				Indent().
				WriteLine("return null").
				Unindent().
				WriteLine("var item_bytes = buffer.get_data(item_len)[1]").
				WriteLine(resultField + ".append(" +
					targetFieldType.GetName() + ".deserialize_binary(item_bytes))")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldTypeName,
				ModelName:      modelName,
				FieldName:      fieldRawName,
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}
	builder.UnindentLine()

	return nil
}
