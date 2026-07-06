package gdGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *GDScriptGenerationContext) generateArraySerializeBinary(
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	targetFieldTypeName := gdStorageTypeName(targetFieldType)
	modelName := field.OwnedBy.GetName()
	fieldRawName := caseUtils.ToSnakeCase(field.GetName())
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
			builder.WriteLine("if item:").
				Indent().
				WriteLine("buffer.put_u8(1)").
				WriteLine("var item_bytes = item.serialize_binary()").
				WriteLine("buffer.put_u32(item_bytes.size())").
				WriteLine("buffer.put_data(item_bytes)").
				Unindent().
				WriteLine("else:").
				Indent().
				WriteLine("buffer.put_u8(0)").
				Unindent()
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

func (c *GDScriptGenerationContext) generateArrayDeserializeBinary(
	field *CCLField,
	builder *codeBuilder.CodeBuilder,
	useWGodot bool,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldRawName := caseUtils.ToSnakeCase(field.GetName())
	resultField := modelResultName + "." + fieldRawName
	targetFieldTypeName := gdStorageTypeName(targetFieldType)
	fieldTargetLangType, err := c.getGDScriptType(field)
	if err != nil {
		return err
	}
	enumCastSuffix, err := c.getGDScriptEnumCastSuffix(targetFieldType, field.OwnedBy)
	if err != nil {
		return err
	}
	getDataCall := ""
	if useWGodot {
		getDataCall = "get_data_bytes(item_len)"
	} else {
		getDataCall = "get_data(item_len)[1]"
	}

	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"fieldLen", fieldRawName+"_len",
		"fieldItemPresent", fieldRawName+"_item_present",
		"field", resultField,
		"fieldT", targetFieldTypeName,
		"fieldTargetT", fieldTargetLangType,
		"getDataCall", getDataCall,
		"enumCast", enumCastSuffix,
	)
	defer builder.UnmapVar(
		"fieldLen",
		"fieldItemPresent",
		"field",
		"fieldT",
		"fieldTargetT",
		"getDataCall",
		"enumCast",
	)

	c.generateBinaryDeserializeBoundsCheck(builder, "4")
	builder.LineD("var $fieldLen := buffer.get_u32()").
		LineD("$field = [] as $fieldTargetT").
		LineD("for i in range($fieldLen):").
		Indent()

	switch targetFieldTypeName {
	case cclValues.TypeNameString:
		c.generateBinaryDeserializeBoundsCheck(builder, "4")
		builder.WriteLine("var item_len := buffer.get_u32()").
			WriteLine("if item_len > buffer.get_size() - buffer.get_position():").
			Indent().
			LineD("return $binaryParseFallback").
			Unindent().
			LineD("var item: String = buffer.$getDataCall.get_string_from_utf8()").
			LineD("$field.append(item)")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		c.generateBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("$field.append(buffer.get_32()$enumCast)")
	case cclValues.TypeNameInt8:
		c.generateBinaryDeserializeBoundsCheck(builder, "1")
		builder.LineD("$field.append(buffer.get_8()$enumCast)")
	case cclValues.TypeNameInt16:
		c.generateBinaryDeserializeBoundsCheck(builder, "2")
		builder.LineD("$field.append(buffer.get_16()$enumCast)")
	case cclValues.TypeNameInt64:
		c.generateBinaryDeserializeBoundsCheck(builder, "8")
		builder.LineD("$field.append(buffer.get_64()$enumCast)")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		c.generateBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("$field.append(buffer.get_u32()$enumCast)")
	case cclValues.TypeNameUint8:
		c.generateBinaryDeserializeBoundsCheck(builder, "1")
		builder.LineD("$field.append(buffer.get_u8()$enumCast)")
	case cclValues.TypeNameUint16:
		c.generateBinaryDeserializeBoundsCheck(builder, "2")
		builder.LineD("$field.append(buffer.get_u16()$enumCast)")
	case cclValues.TypeNameUint64:
		c.generateBinaryDeserializeBoundsCheck(builder, "8")
		builder.LineD("$field.append(buffer.get_u64()$enumCast)")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		c.generateBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("$field.append(buffer.get_float())")
	case cclValues.TypeNameBool:
		c.generateBinaryDeserializeBoundsCheck(builder, "1")
		builder.LineD("$field.append(buffer.get_8() != 0)")
	case cclValues.TypeNameBytes:
		c.generateBinaryDeserializeBoundsCheck(builder, "4")
		builder.WriteLine("var item_len := buffer.get_u32()").
			WriteLine("if item_len > buffer.get_size() - buffer.get_position():").
			Indent().
			LineD("return $binaryParseFallback").
			Unindent().
			LineD("$field.append(buffer.$getDataCall)")
	default:
		if targetFieldType.IsCustomTypeModel() {
			c.generateBinaryDeserializeBoundsCheck(builder, "1")
			builder.LineD("var $fieldItemPresent := buffer.get_u8()").
				LineD("if $fieldItemPresent == 0:").
				Indent().
				LineD("$field.append(null)").
				Unindent().
				LineD("elif $fieldItemPresent == 1:").
				Indent()
			c.generateBinaryDeserializeBoundsCheck(builder, "4")
			builder.WriteLine("var item_len := buffer.get_u32()").
				WriteLine("if item_len > buffer.get_size() - buffer.get_position():").
				Indent().
				LineD("return $binaryParseFallback").
				Unindent().
				LineD("var item_bytes: PackedByteArray = buffer.$getDataCall").
				LineD("$field.append($fieldT.deserialize_binary(item_bytes))").
				Unindent().
				WriteLine("else:").
				Indent().
				LineD("return $binaryParseFallback").
				Unindent()
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
