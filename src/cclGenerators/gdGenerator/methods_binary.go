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

	for _, field := range model.Fields {
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
	fieldRawName := cclUtils.ToSnakeCase(field.GetName())
	resultField := "self." + fieldRawName
	modelName := field.OwnedBy.GetName()
	targetFieldTypeName := field.Type.GetName()

	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"fieldBytes", fieldRawName+"_bytes",
		"field", resultField,
		"fieldT", targetFieldTypeName,
	)
	defer builder.UnmapVar(
		"fieldBytes",
		"field",
		"fieldT",
	)

	switch targetFieldTypeName {
	case cclValues.TypeNameString:
		builder.LineD("var $fieldBytes = $field.to_utf8_buffer()").
			LineD("buffer.put_u32($fieldBytes.size())").
			LineD("buffer.put_data($fieldBytes)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.LineD("buffer.put_32($field)")
	case cclValues.TypeNameInt8:
		builder.LineD("buffer.put_8($field)")
	case cclValues.TypeNameInt16:
		builder.LineD("buffer.put_16($field)")
	case cclValues.TypeNameInt64:
		builder.LineD("buffer.put_64($field)")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.LineD("buffer.put_u32($field)")
	case cclValues.TypeNameUint8:
		builder.LineD("buffer.put_u8($field)")
	case cclValues.TypeNameUint16:
		builder.LineD("buffer.put_u16($field)")
	case cclValues.TypeNameUint64:
		builder.LineD("buffer.put_u64($field)")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.LineD("buffer.put_float($field)")
	case cclValues.TypeNameBool:
		builder.LineD("buffer.put_8(1 if $field else 0)")
	case cclValues.TypeNameBytes:
		builder.LineD("buffer.put_u32($field.size())").
			LineD("buffer.put_data($field)")
	case cclValues.TypeNameDateTime:
		builder.LineD("buffer.put_64($field)")
	default:
		// Custom type handling
		if field.Type.IsCustomTypeModel() {
			builder.LineD("var $fieldBytes = $field.serialize_binary() if $field else PackedByteArray([0])").
				LineD("buffer.put_u32($fieldBytes.size())").
				LineD("buffer.put_data($fieldBytes)")
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

	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"fieldBytes", fieldRawName+"_bytes",
		"field", resultField,
		"fieldT", targetFieldTypeName,
	)
	defer builder.UnmapVar(
		"fieldBytes",
		"field",
		"fieldT",
	)

	builder.LineD("buffer.put_u32($field.size())").
		LineD("for item in $field:").
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
	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"modelResult", modelResultName,
	)
	defer builder.UnmapVar(
		"modelResult",
	)

	builder.LineD("static func deserialize_binary(data: PackedByteArray) -> $model:").
		Indent().
		// null-safety check
		WriteLine("if not data or data.is_empty() or (data.size() == 1 and data[0] == 0):").
		Indent().
		WriteLine("return null").
		UnindentLine().
		WriteLine("var buffer = StreamPeerBuffer.new()").
		WriteLine("buffer.data_array = data")

	builder.LineD("var $modelResult = $model.new()").
		NewLine()

	for _, field := range model.Fields {
		if field.IsArray() {
			err := c.generateArrayDeserializeBinary(field, builder)
			if err != nil {
				return err
			}

			continue
		}

		err := c.generateFieldDeserializeBinary(field, builder)
		if err != nil {
			return err
		}
	}

	builder.LineD("return $modelResult").
		UnindentLine()

	return nil
}

func (c *GDScriptGenerationContext) generateFieldDeserializeBinary(
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) error {
	fieldRawName := cclUtils.ToSnakeCase(field.GetName())
	resultField := modelResultName + "." + fieldRawName
	targetFieldTypeName := field.Type.GetName()

	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"field", resultField,
		"fieldT", targetFieldTypeName,
		"fieldLen", fieldRawName+"_len",
		"fieldBytes", fieldRawName+"_bytes",
	)
	defer builder.UnmapVar(
		"field",
		"fieldT",
		"fieldLen",
		"fieldBytes",
	)

	switch targetFieldTypeName {
	case cclValues.TypeNameString:
		builder.LineD("var $fieldLen = buffer.get_u32()").
			LineD("if $fieldLen > buffer.get_size() - buffer.get_position():").
			Indent().
			WriteLine("return null").
			Unindent().
			LineD("$field = buffer.get_data($fieldLen)[1].get_string_from_utf8()")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.LineD("$field = buffer.get_32()")
	case cclValues.TypeNameInt8:
		builder.LineD("$field = buffer.get_8()")
	case cclValues.TypeNameInt16:
		builder.LineD("$field = buffer.get_16()")
	case cclValues.TypeNameInt64:
		builder.LineD("$field = buffer.get_64()")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.LineD("$field = buffer.get_u32()")
	case cclValues.TypeNameUint8:
		builder.LineD("$field = buffer.get_u8()")
	case cclValues.TypeNameUint16:
		builder.LineD("$field = buffer.get_u16()")
	case cclValues.TypeNameUint64:
		builder.LineD("$field = buffer.get_u64()")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.LineD("$field = buffer.get_float()")
	case cclValues.TypeNameBool:
		builder.LineD("$field = buffer.get_8() != 0")
	case cclValues.TypeNameBytes:
		builder.LineD("var $fieldLen = buffer.get_u32()").
			LineD("$field = buffer.get_data($fieldLen)[1]")
	case cclValues.TypeNameDateTime:
		builder.LineD("$field = buffer.get_64()")
	default:
		// Custom type handling
		if field.IsCustomTypeModel() {
			builder.LineD("var $fieldLen = buffer.get_u32()").
				LineD("var $fieldBytes = buffer.get_data($fieldLen)[1]").
				LineD("$field = $fieldT.deserialize_binary($fieldBytes)")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldTypeName,
				ModelName:      field.GetModelFullName(),
				FieldName:      fieldRawName,
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}
	builder.NewLine()

	return nil
}

func (c *GDScriptGenerationContext) generateArrayDeserializeBinary(
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldRawName := cclUtils.ToSnakeCase(field.GetName())
	resultField := modelResultName + "." + fieldRawName
	targetFieldTypeName := targetFieldType.GetName()
	fieldTargetLangType := c.getGDScriptType(field)

	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"fieldLen", fieldRawName+"_len",
		"field", resultField,
		"fieldT", targetFieldTypeName,
		"fieldTargetT", fieldTargetLangType,
	)
	defer builder.UnmapVar(
		"fieldLen",
		"field",
		"fieldT",
		"fieldTargetT",
	)

	builder.LineD("var $fieldLen = buffer.get_u32()").
		LineD("$field = [] as $fieldTargetT").
		LineD("for i in range($fieldLen):").
		Indent()

	switch targetFieldTypeName {
	case cclValues.TypeNameString:
		builder.WriteLine("var item_len = buffer.get_u32()").
			WriteLine("if item_len > buffer.get_size() - buffer.get_position():").
			Indent().
			WriteLine("return null").
			Unindent().
			WriteLine("var item = buffer.get_data(item_len)[1].get_string_from_utf8()").
			LineD("$field.append(item)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.LineD("$field.append(buffer.get_32())")
	case cclValues.TypeNameInt8:
		builder.LineD("$field.append(buffer.get_8())")
	case cclValues.TypeNameInt16:
		builder.LineD("$field.append(buffer.get_16())")
	case cclValues.TypeNameInt64:
		builder.LineD("$field.append(buffer.get_64())")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.LineD("$field.append(buffer.get_u32())")
	case cclValues.TypeNameUint8:
		builder.LineD("$field.append(buffer.get_u8())")
	case cclValues.TypeNameUint16:
		builder.LineD("$field.append(buffer.get_u16())")
	case cclValues.TypeNameUint64:
		builder.LineD("$field.append(buffer.get_u64())")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		builder.LineD("$field.append(buffer.get_float())")
	case cclValues.TypeNameBool:
		builder.LineD("$field.append(buffer.get_8() != 0)")
	case cclValues.TypeNameBytes:
		builder.WriteLine("var item_len = buffer.get_u32()").
			WriteLine("if item_len > buffer.get_size() - buffer.get_position():").
			Indent().
			WriteLine("return null").
			Unindent().
			LineD("$field.append(buffer.get_data(item_len)[1])")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("var item_len = buffer.get_u32()").
				WriteLine("if item_len > buffer.get_size() - buffer.get_position():").
				Indent().
				WriteLine("return null").
				Unindent().
				WriteLine("var item_bytes = buffer.get_data(item_len)[1]").
				LineD("$field.append($fieldT.deserialize_binary(item_bytes))")
		} else {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldTypeName,
				ModelName:      field.GetModelFullName(),
				FieldName:      fieldRawName,
				TargetLanguage: CurrentLanguage.String(),
			}
		}
	}
	builder.UnindentLine()

	return nil
}
