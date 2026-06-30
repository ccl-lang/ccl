package gdGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclErrors"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *GDScriptGenerationContext) generateSerializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	endian, err := c.GetBinarySerializationEndian(CurrentLanguage, model)
	if err != nil {
		return err
	}
	bigEndian := "false"
	if endian == gValues.EndianBig {
		bigEndian = "true"
	}

	builder.WriteLine("func serialize_binary() -> PackedByteArray:").
		Indent().
		WriteLine("var buffer := StreamPeerBuffer.new()").
		WriteLine("buffer.big_endian = " + bigEndian).
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
	fieldRawName := caseUtils.ToSnakeCase(field.GetName())
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

func (c *GDScriptGenerationContext) generateDeserializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	endian, err := c.GetBinarySerializationEndian(CurrentLanguage, model)
	if err != nil {
		return err
	}
	strictBinaryParsing, err := c.UsesStrictBinaryParsing(CurrentLanguage, model)
	if err != nil {
		return err
	}
	useWGodot, err := c.UsesWGodot(CurrentLanguage, model)
	if err != nil {
		return err
	}
	bigEndian := "false"
	if endian == gValues.EndianBig {
		bigEndian = "true"
	}
	binaryParseFallback := modelResultName
	if strictBinaryParsing {
		binaryParseFallback = "null"
	}

	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"modelResult", modelResultName,
		"binaryParseFallback", binaryParseFallback,
	)
	defer builder.UnmapVar(
		"modelResult",
		"binaryParseFallback",
	)

	builder.LineD("static func deserialize_binary(data: PackedByteArray) -> $model:").
		Indent().
		// null-safety check
		WriteLine("if not data or data.is_empty() or (data.size() == 1 and data[0] == 0):").
		Indent().
		WriteLine("return null").
		UnindentLine().
		WriteLine("var buffer := StreamPeerBuffer.new()").
		WriteLine("buffer.big_endian = " + bigEndian).
		WriteLine("buffer.data_array = data")

	builder.LineD("var $modelResult := $model.new()").
		NewLine()

	for _, field := range model.Fields {
		if field.IsArray() {
			err := c.generateArrayDeserializeBinary(
				field,
				builder,
				useWGodot,
			)
			if err != nil {
				return err
			}

			continue
		}

		err := c.generateFieldDeserializeBinary(
			field,
			builder,
			useWGodot,
		)
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
	useWGodot bool,
) error {
	fieldRawName := caseUtils.ToSnakeCase(field.GetName())
	resultField := modelResultName + "." + fieldRawName
	targetFieldTypeName := field.Type.GetName()
	fieldLenName := fieldRawName + "_len"
	getDataCall := ""
	if useWGodot {
		getDataCall = "get_data_bytes(" + fieldLenName + ")"
	} else {
		getDataCall = "get_data(" + fieldLenName + ")[1]"
	}

	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"field", resultField,
		"fieldT", targetFieldTypeName,
		"fieldLen", fieldLenName,
		"fieldBytes", fieldRawName+"_bytes",
		"getDataCall", getDataCall,
	)
	defer builder.UnmapVar(
		"field",
		"fieldT",
		"fieldLen",
		"fieldBytes",
		"getDataCall",
	)

	switch targetFieldTypeName {
	case cclValues.TypeNameString:
		c.generateBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("var $fieldLen := buffer.get_u32()").
			LineD("if $fieldLen > buffer.get_size() - buffer.get_position():").
			Indent().
			LineD("return $binaryParseFallback").
			Unindent().
			LineD("$field = buffer.$getDataCall.get_string_from_utf8()")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		c.generateBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("$field = buffer.get_32()")
	case cclValues.TypeNameInt8:
		c.generateBinaryDeserializeBoundsCheck(builder, "1")
		builder.LineD("$field = buffer.get_8()")
	case cclValues.TypeNameInt16:
		c.generateBinaryDeserializeBoundsCheck(builder, "2")
		builder.LineD("$field = buffer.get_16()")
	case cclValues.TypeNameInt64:
		c.generateBinaryDeserializeBoundsCheck(builder, "8")
		builder.LineD("$field = buffer.get_64()")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		c.generateBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("$field = buffer.get_u32()")
	case cclValues.TypeNameUint8:
		c.generateBinaryDeserializeBoundsCheck(builder, "1")
		builder.LineD("$field = buffer.get_u8()")
	case cclValues.TypeNameUint16:
		c.generateBinaryDeserializeBoundsCheck(builder, "2")
		builder.LineD("$field = buffer.get_u16()")
	case cclValues.TypeNameUint64:
		c.generateBinaryDeserializeBoundsCheck(builder, "8")
		builder.LineD("$field = buffer.get_u64()")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32, cclValues.TypeNameFloat64:
		c.generateBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("$field = buffer.get_float()")
	case cclValues.TypeNameBool:
		c.generateBinaryDeserializeBoundsCheck(builder, "1")
		builder.LineD("$field = buffer.get_8() != 0")
	case cclValues.TypeNameBytes:
		c.generateBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("var $fieldLen := buffer.get_u32()").
			LineD("if $fieldLen > buffer.get_size() - buffer.get_position():").
			Indent().
			LineD("return $binaryParseFallback").
			Unindent().
			LineD("$field = buffer.$getDataCall")
	case cclValues.TypeNameDateTime:
		c.generateBinaryDeserializeBoundsCheck(builder, "8")
		builder.LineD("$field = buffer.get_64()")
	default:
		// Custom type handling
		if field.IsCustomTypeModel() {
			c.generateBinaryDeserializeBoundsCheck(builder, "4")
			builder.LineD("var $fieldLen := buffer.get_u32()").
				LineD("if $fieldLen > buffer.get_size() - buffer.get_position():").
				Indent().
				LineD("return $binaryParseFallback").
				Unindent().
				LineD("var $fieldBytes: PackedByteArray = buffer.$getDataCall").
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

func (c *GDScriptGenerationContext) generateBinaryDeserializeBoundsCheck(
	builder *codeBuilder.CodeBuilder,
	requiredBytes string,
) {
	builder.MapVarPairs("requiredBytes", requiredBytes)
	builder.LineD("if buffer.get_size() - buffer.get_position() < $requiredBytes:").
		Indent().
		LineD("return $binaryParseFallback").
		Unindent()
	builder.UnmapVar("requiredBytes")
}
