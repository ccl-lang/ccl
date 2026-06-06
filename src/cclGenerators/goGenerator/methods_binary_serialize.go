package goGenerator

import (
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

//---------------------------------------------------------

func (c *GoGenerationContext) generateSerializeBinaryMethod(
	builder *codeBuilder.CodeBuilder,
	model *CCLModel,
) error {
	// Generated SerializeBinary methods build output with bytes.Buffer.
	registerGoImport(builder, "bytes")
	// Generated SerializeBinary methods write fields with encoding/binary.
	registerGoImport(builder, "encoding/binary")

	endian, err := c.GetBinarySerializationEndian(CurrentLanguage, model)
	if err != nil {
		return err
	}
	binaryEndianInit := "binary.LittleEndian"
	if endian == gValues.EndianBig {
		binaryEndianInit = "binary.BigEndian"
	}

	builder.ExpectMappedVars(
		"model",
	)
	builder.MapVarPairs(
		"binaryEndianInit", binaryEndianInit,
	)
	defer builder.UnmapVar(
		"binaryEndianInit",
	)

	builder.NewLine().
		LineD("func (m $model) SerializeBinary() ([]byte, error) {").
		Indent().
		// handle m is nil by returning []byte(0) and nil
		WriteLine("if m == nil {").
		Indent().
		WriteLine("return []byte{0}, nil").
		Unindent().
		WriteLine("}").
		NewLine()

	builder.WriteLine("buf := new(bytes.Buffer)").
		LineD("binaryEndian := $binaryEndianInit").
		NewLine()
	for _, field := range model.Fields {
		if field.IsArray() {
			err := c.generateArraySerializeBinaryMethod(builder, field)
			if err != nil {
				return err
			}

			continue
		}

		err := c.generateFieldSerializeBinaryMethod(builder, field)
		if err != nil {
			return err
		}
	}

	builder.WriteLine("return buf.Bytes(), nil").
		Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *GoGenerationContext) generateFieldSerializeBinaryMethod(
	builder *codeBuilder.CodeBuilder,
	field *CCLField,
) error {
	fieldTypeName := field.Type.GetName()
	isCustomType := field.Type.IsCustomTypeModel()
	fieldName := field.GetName()
	fieldVar := "m." + fieldName
	fieldBytesName := "current_" + fieldName + "Bytes"

	builder.MapVarPairs(
		"field", fieldVar,
		"fieldBytes", fieldBytesName,
	)
	defer builder.UnmapVar(
		"field",
		"fieldBytes",
	)

	switch fieldTypeName {
	case cclValues.TypeNameString:
		builder.LineD("if err := binary.Write(buf, binaryEndian, uint32(len($field))); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}").
			LineD("if _, err := buf.WriteString($field); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}")
	case cclValues.TypeNameBytes:
		builder.LineD("if err := binary.Write(buf, binaryEndian, uint32(len($field))); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}").
			LineD("if _, err := buf.Write($field); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}")
	case cclValues.TypeNameDateTime:
		builder.LineD("if err := binary.Write(buf, binaryEndian, $field.UnixNano()); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}")
	default:
		if isCustomType {
			builder.LineD("$fieldBytes, err := $field.SerializeBinary()").
				WriteLine("if err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}").
				LineD("if err := binary.Write(buf, binaryEndian, uint32(len($fieldBytes))); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}").
				LineD("if err := binary.Write(buf, binaryEndian, $fieldBytes); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}")
		} else {
			toWriteStr := fieldVar
			if fieldTypeName == cclValues.TypeNameInt {
				// binary.Write does not support int type directly, so we need to convert it to int32
				toWriteStr = "int32(" + fieldVar + ")"
			}
			builder.MapVarPairs("toWrite", toWriteStr)
			builder.LineD("if err := binary.Write(buf, binaryEndian, $toWrite); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}")
			builder.UnmapVar("toWrite")
		}
	}
	return nil
}

func (c *GoGenerationContext) generateArraySerializeBinaryMethod(
	builder *codeBuilder.CodeBuilder,
	field *CCLField,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	targetFieldTypeName := targetFieldType.GetName()
	isCustomType := targetFieldType.IsCustomTypeModel()
	fieldName := field.Name
	fieldVar := "m." + fieldName
	fieldBytesName := "current_" + fieldName + "Bytes"

	builder.MapVarPairs(
		"field", fieldVar,
		"fieldBytes", fieldBytesName,
	)
	defer builder.UnmapVar(
		"field",
		"fieldBytes",
	)

	builder.LineD("if err := binary.Write(buf, binaryEndian, uint32(len($field))); err != nil {").
		Indent().
		WriteLine("return nil, err").
		Unindent().
		WriteLine("}").
		LineD("for _, elem := range $field {").
		Indent()
	switch targetFieldTypeName {
	case cclValues.TypeNameString:
		builder.WriteLine("if err := binary.Write(buf, binaryEndian, uint32(len(elem))); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}").
			WriteLine("if _, err := buf.WriteString(elem); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}")
	case cclValues.TypeNameBytes:
		builder.WriteLine("if err := binary.Write(buf, binaryEndian, uint32(len(elem))); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}").
			WriteLine("if _, err := buf.Write(elem); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}")
	case cclValues.TypeNameDateTime:
		builder.WriteLine("if err := binary.Write(buf, binaryEndian, elem.UnixNano()); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}")
	default:
		if isCustomType {
			builder.LineD("$fieldBytes, err := elem.SerializeBinary()").
				WriteLine("if err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}").
				LineD("if err := binary.Write(buf, binaryEndian, uint32(len($fieldBytes))); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}").
				LineD("if err := binary.Write(buf, binaryEndian, $fieldBytes); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}")
		} else {
			toWriteStr := "elem"
			if targetFieldTypeName == cclValues.TypeNameInt {
				// binary.Write does not support int type directly, so we need to convert it to int32
				toWriteStr = "int32(elem)"
			}
			builder.MapVarPairs("toWrite", toWriteStr)
			builder.LineD("if err := binary.Write(buf, binaryEndian, $toWrite); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}")
			builder.UnmapVar("toWrite")
		}
	}

	builder.Unindent().
		WriteLine("}")
	return nil
}

//---------------------------------------------------------
