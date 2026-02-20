package goGenerator

import (
	"strings"

	"github.com/ccl-lang/ccl/src/core/cclErrors"
	"github.com/ccl-lang/ccl/src/core/cclValues"
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

//---------------------------------------------------------

func (c *GoGenerationContext) generateSerializeBinaryMethod(model *CCLModel) error {
	endian, err := c.GetBinarySerializationEndian(CurrentLanguage, model)
	if err != nil {
		return err
	}
	binaryEndianInit := "binary.LittleEndian"
	if endian == gValues.EndianBig {
		binaryEndianInit = "binary.BigEndian"
	}

	c.MethodsCode.ExpectMappedVars(
		"model",
	)
	c.MethodsCode.MapVarPairs(
		"binaryEndianInit", binaryEndianInit,
	)
	defer c.MethodsCode.UnmapVar(
		"binaryEndianInit",
	)

	c.MethodsCode.NewLine().
		LineD("func (m $model) SerializeBinary() ([]byte, error) {").
		Indent().
		// handle m is nil by returning []byte(0) and nil
		WriteLine("if m == nil {").
		Indent().
		WriteLine("return []byte{0}, nil").
		Unindent().
		WriteLine("}").
		NewLine()

	c.MethodsCode.WriteLine("buf := new(bytes.Buffer)").
		LineD("binaryEndian := $binaryEndianInit").
		NewLine()
	for _, field := range model.Fields {
		if field.IsArray() {
			err := c.generateArraySerializeBinaryMethod(field)
			if err != nil {
				return err
			}

			continue
		}

		err := c.generateFieldSerializeBinaryMethod(field)
		if err != nil {
			return err
		}
	}

	c.MethodsCode.WriteLine("return buf.Bytes(), nil").
		Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *GoGenerationContext) generateFieldSerializeBinaryMethod(field *CCLField) error {
	fieldTypeName := field.Type.GetName()
	isCustomType := field.Type.IsCustomTypeModel()
	fieldName := field.GetName()
	fieldVar := "m." + fieldName
	fieldBytesName := "current_" + fieldName + "Bytes"

	c.MethodsCode.MapVarPairs(
		"field", fieldVar,
		"fieldBytes", fieldBytesName,
	)
	defer c.MethodsCode.UnmapVar(
		"field",
		"fieldBytes",
	)

	switch fieldTypeName {
	case cclValues.TypeNameString:
		c.MethodsCode.LineD("if err := binary.Write(buf, binaryEndian, uint32(len($field))); err != nil {").
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
		c.MethodsCode.LineD("if err := binary.Write(buf, binaryEndian, uint32(len($field))); err != nil {").
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
		c.MethodsCode.LineD("if err := binary.Write(buf, binaryEndian, $field.UnixNano()); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}")
	default:
		if isCustomType {
			c.MethodsCode.LineD("$fieldBytes, err := $field.SerializeBinary()").
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
			c.MethodsCode.MapVarPairs("toWrite", toWriteStr)
			c.MethodsCode.LineD("if err := binary.Write(buf, binaryEndian, $toWrite); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}")
			c.MethodsCode.UnmapVar("toWrite")
		}
	}
	return nil
}

func (c *GoGenerationContext) generateArraySerializeBinaryMethod(field *CCLField) error {
	targetFieldType := field.Type.GetUnderlyingType()
	isCustomType := targetFieldType.IsCustomTypeModel()
	fieldName := field.Name
	fieldVar := "m." + fieldName
	fieldBytesName := "current_" + fieldName + "Bytes"

	c.MethodsCode.MapVarPairs(
		"field", fieldVar,
		"fieldBytes", fieldBytesName,
	)
	defer c.MethodsCode.UnmapVar(
		"field",
		"fieldBytes",
	)

	c.MethodsCode.LineD("if err := binary.Write(buf, binaryEndian, uint32(len($field))); err != nil {").
		Indent().
		WriteLine("return nil, err").
		Unindent().
		WriteLine("}").
		LineD("for _, elem := range $field {").
		Indent()
	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		c.MethodsCode.WriteLine("if err := binary.Write(buf, binaryEndian, uint32(len(elem))); err != nil {").
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
		c.MethodsCode.WriteLine("if err := binary.Write(buf, binaryEndian, uint32(len(elem))); err != nil {").
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
		c.MethodsCode.WriteLine("if err := binary.Write(buf, binaryEndian, elem.UnixNano()); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}")
	default:
		if isCustomType {
			c.MethodsCode.LineD("$fieldBytes, err := elem.SerializeBinary()").
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
			c.MethodsCode.WriteLine("if err := binary.Write(buf, binaryEndian, elem); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}")
		}
	}

	c.MethodsCode.Unindent().
		WriteLine("}")
	return nil
}

func (c *GoGenerationContext) generateDeserializeBinaryMethod(model *CCLModel) error {
	endian, err := c.GetBinarySerializationEndian(CurrentLanguage, model)
	if err != nil {
		return err
	}
	binaryEndianInit := "binary.LittleEndian"
	if endian == gValues.EndianBig {
		binaryEndianInit = "binary.BigEndian"
	}

	c.MethodsCode.ExpectMappedVars(
		"model",
	)
	c.MethodsCode.MapVarPairs(
		"binaryEndianInit", binaryEndianInit,
	)
	defer c.MethodsCode.UnmapVar(
		"binaryEndianInit",
	)

	c.MethodsCode.LineD("func (m $model) DeserializeBinary(data []byte) error {").
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
			err := c.generateArrayDeserializeBinaryMethod(field)
			if err != nil {
				return err
			}

			continue
		}

		err := c.generateFieldDeserializeBinaryMethod(field)
		if err != nil {
			return err
		}
	}

	c.MethodsCode.WriteLine("return nil").
		Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *GoGenerationContext) generateFieldDeserializeBinaryMethod(field *CCLField) error {
	fieldTypeName := field.Type.GetName()
	isCustomType := field.Type.IsCustomTypeModel()
	// isPointer := isCustomType //TODO: Find a way to specify this
	fieldName := field.Name
	fieldVar := "m." + fieldName
	fName := strings.ToLower(string(fieldName[0])) + fieldName[1:]
	fLenName := fName + "Len"
	fNameStrBytes := fName + "StrBytes"
	fNameUnix := fName + "Unix"

	c.MethodsCode.MapVarPairs(
		"field", fieldVar,
		"fieldLen", fLenName,
		"fieldStrBytes", fNameStrBytes,
		"fieldUnix", fNameUnix,
		"fieldName", fieldName,
	)
	defer c.MethodsCode.UnmapVar(
		"field",
		"fieldLen",
		"fieldStrBytes",
		"fieldUnix",
		"fieldName",
	)

	switch fieldTypeName {
	case cclValues.TypeNameString:
		c.MethodsCode.LineD("var $fieldLen uint32").
			LineD("if err := binary.Read(buf, binaryEndian, &$fieldLen); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			LineD("$fieldStrBytes := make([]byte, $fieldLen)").
			LineD("if _, err := buf.Read($fieldStrBytes); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			LineD("$field = string($fieldStrBytes)")
	case cclValues.TypeNameBytes:
		c.MethodsCode.LineD("var $fieldLen uint32").
			LineD("if err := binary.Read(buf, binaryEndian, &$fieldLen); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			LineD("bytesData := make([]byte, $fieldLen)").
			WriteLine("if _, err := buf.Read(bytesData); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			LineD("$field = bytesData")
	case cclValues.TypeNameDateTime:
		c.MethodsCode.LineD("var $fieldUnix int64").
			LineD("if err := binary.Read(buf, binaryEndian, &$fieldUnix); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			LineD("$field = time.Unix(0, $fieldUnix)")
	default:
		if isCustomType {
			// read the length of the next buffer that we need
			// var basicLen uint32
			// if err := binary.Read(buf, binaryEndian, &basicLen); err != nil {
			// \treturn err
			// }
			lenVarName := fName + "BytesLen"
			bytesVarName := fName + "Bytes"
			fieldType := field.Type.GetName()

			c.MethodsCode.MapVarPairs(
				"fieldBytesLen", lenVarName,
				"fieldBytes", bytesVarName,
				"fieldType", fieldType,
			)

			c.MethodsCode.LineD("var $fieldBytesLen uint32").
				LineD("if err := binary.Read(buf, binaryEndian, &$fieldBytesLen); err != nil {").
				Indent().
				WriteLine("return err").
				Unindent().
				WriteLine("}")

			c.MethodsCode.LineD("$fieldBytes := make([]byte, $fieldBytesLen)").
				LineD("if _, err := buf.Read($fieldBytes); err != nil {").
				Indent().
				WriteLine("return err").
				Unindent().
				WriteLine("}")

			// make sure m.field is not nil ONLY when len(bytesVarName) != 0 and !(len(bytesVarName) == 1 and bytesVarName[0] == 0)
			c.MethodsCode.LineD("if m.$fieldName == nil && len($fieldBytes) != 0 && !(len($fieldBytes) == 1 && $fieldBytes[0] == 0) {").
				Indent().
				LineD("m.$fieldName = new($fieldType)").
				Unindent().
				WriteLine("}")

			c.MethodsCode.LineD("if err := m.$fieldName.DeserializeBinary($fieldBytes); err != nil {").
				Indent().
				WriteLine("return err").
				Unindent().
				WriteLine("}")

			c.MethodsCode.UnmapVar(
				"fieldBytesLen",
				"fieldBytes",
				"fieldType",
			)
		} else {
			if fieldTypeName == cclValues.TypeNameInt {
				// binary.Read does not support int type directly, so we need to read it into an int32 first
				toReadName := "tmp" + fieldName
				c.MethodsCode.MapVarPairs("toRead", toReadName)
				c.MethodsCode.LineD("var $toRead int32").
					LineD("if err := binary.Read(buf, binaryEndian, &$toRead); err != nil {").
					Indent().
					WriteLine("return err").
					Unindent().
					WriteLine("}").
					LineD("$field = int($toRead)")
				c.MethodsCode.UnmapVar("toRead")
				return nil
			}

			// all other types can be read directly
			c.MethodsCode.LineD("if err := binary.Read(buf, binaryEndian, &$field); err != nil {").
				Indent().
				WriteLine("return err").
				Unindent().
				WriteLine("}")
		}
	}
	return nil
}

func (c *GoGenerationContext) generateArrayDeserializeBinaryMethod(field *CCLField) error {
	targetFieldType := field.Type.GetUnderlyingType()
	targetFieldTypeName := targetFieldType.GetName()
	isCustomType := targetFieldType.IsCustomTypeModel()
	isPointer := isCustomType //TODO: Find a way to specify this
	fieldName := field.Name
	fieldVar := "m." + fieldName
	fName := strings.ToLower(string(fieldName[0])) + fieldName[1:]
	fLenName := fName + "Len"
	// fNameStrBytes := fName + "StrBytes"
	// fNameUnix := fName + "Unix"
	fieldRealType := ""
	if isCustomType {
		fieldRealType = "*" + targetFieldTypeName
	} else {
		mappedType, ok := CCLTypesToGoTypes[targetFieldTypeName]
		if !ok {
			return &cclErrors.UnsupportedFieldTypeError{
				TypeName:       targetFieldTypeName,
				FieldName:      fieldName,
				ModelName:      field.GetModelFullName(),
				TargetLanguage: CurrentLanguage.String(),
			}
		}
		fieldRealType = mappedType
	}

	c.MethodsCode.MapVarPairs(
		"field", fieldVar,
		"fieldLen", fLenName,
		"fieldRealType", fieldRealType,
	)
	defer c.MethodsCode.UnmapVar(
		"field",
		"fieldLen",
		"fieldRealType",
	)

	c.MethodsCode.LineD("var $fieldLen uint32").
		LineD("if err := binary.Read(buf, binaryEndian, &$fieldLen); err != nil {").
		Indent().
		WriteLine("return err").
		Unindent().
		WriteLine("}").
		LineD("$field = make([]$fieldRealType, $fieldLen)").
		LineD("for i := uint32(0); i < $fieldLen; i++ {").
		Indent()
	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		c.MethodsCode.WriteLine("var elemLen uint32").
			WriteLine("if err := binary.Read(buf, binaryEndian, &elemLen); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			WriteLine("elemBytes := make([]byte, elemLen)").
			WriteLine("if _, err := buf.Read(elemBytes); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			LineD("$field[i] = string(elemBytes)")
	case cclValues.TypeNameBytes:
		c.MethodsCode.WriteLine("var elemLen uint32").
			WriteLine("if err := binary.Read(buf, binaryEndian, &elemLen); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			WriteLine("elemBytes := make([]byte, elemLen)").
			WriteLine("if _, err := buf.Read(elemBytes); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			LineD("$field[i] = elemBytes")
	case cclValues.TypeNameDateTime:
		c.MethodsCode.WriteLine("var elemUnix int64").
			WriteLine("if err := binary.Read(buf, binaryEndian, &elemUnix); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			LineD("$field[i] = time.Unix(0, elemUnix)")
	default:
		if isCustomType {
			c.MethodsCode.MapVarPairs("fieldType", targetFieldType.GetName())
			if isPointer {
				c.MethodsCode.LineD("var elem $fieldRealType = new($fieldType)")
			} else {
				c.MethodsCode.LineD("var elem $fieldRealType")
			}
			// we need to read the bytes from the buffer
			// and then deserialize the element
			c.MethodsCode.WriteLine("var elemLen uint32").
				WriteLine("if err := binary.Read(buf, binaryEndian, &elemLen); err != nil {").
				Indent().
				WriteLine("return err").
				Unindent().
				WriteLine("}").
				WriteLine("elemBytes := make([]byte, elemLen)").
				WriteLine("if _, err := buf.Read(elemBytes); err != nil {").
				Indent().
				WriteLine("return err").
				Unindent().
				WriteLine("}").
				WriteLine("if err := elem.DeserializeBinary(elemBytes); err != nil {").
				Indent().
				WriteLine("return err").
				Unindent().
				WriteLine("}").
				LineD("$field[i] = elem")
			c.MethodsCode.UnmapVar("fieldType")
		} else {
			c.MethodsCode.LineD("if err := binary.Read(buf, binaryEndian, &$field[i]); err != nil {").
				Indent().
				WriteLine("return err").
				Unindent().
				WriteLine("}")
		}
	}

	c.MethodsCode.Unindent().
		WriteLine("}")
	return nil
}

//---------------------------------------------------------
