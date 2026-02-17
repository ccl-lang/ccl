package goGenerator

import (
	"strings"

	"github.com/ccl-lang/ccl/src/core/cclValues"
)

//---------------------------------------------------------

func (c *GoGenerationContext) generateSerializeBinaryMethod(model *CCLModel) error {
	c.MethodsCode.NewLine().
		WriteLine("func (m *" + model.Name + ") SerializeBinary() ([]byte, error) {").
		Indent().
		// handle m is nil by returning []byte(0) and nil
		WriteLine("if m == nil {").
		Indent().
		WriteLine("return []byte{0}, nil").
		Unindent().
		WriteLine("}").
		NewLine()

	c.MethodsCode.WriteLine("buf := new(bytes.Buffer)").
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
	switch fieldTypeName {
	case cclValues.TypeNameString:
		c.MethodsCode.WriteLine("if err := binary.Write(buf, binary.LittleEndian, uint32(len(m." + field.Name + "))); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}").
			WriteLine("if _, err := buf.WriteString(m." + field.Name + "); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}")
	case cclValues.TypeNameBytes:
		c.MethodsCode.WriteLine("if err := binary.Write(buf, binary.LittleEndian, uint32(len(m." + field.Name + "))); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}").
			WriteLine("if _, err := buf.Write(m." + field.Name + "); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}")
	case cclValues.TypeNameDateTime:
		c.MethodsCode.WriteLine("if err := binary.Write(buf, binary.LittleEndian, m." + field.Name + ".UnixNano()); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}")
	default:
		if isCustomType {
			currentBytesName := "current_" + field.Name + "Bytes"
			c.MethodsCode.WriteLine(currentBytesName + ", err := m." + field.Name + ".SerializeBinary()").
				WriteLine("if err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}").
				WriteLine("if err := binary.Write(buf, binary.LittleEndian, uint32(len(" + currentBytesName + "))); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}").
				WriteLine("if err := binary.Write(buf, binary.LittleEndian, " + currentBytesName + "); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}")
		} else {
			toWriteStr := "m." + field.Name
			if fieldTypeName == cclValues.TypeNameInt {
				// binary.Write does not support int type directly, so we need to convert it to int32
				toWriteStr = "int32(m." + field.Name + ")"
			}
			c.MethodsCode.WriteLine("if err := binary.Write(buf, binary.LittleEndian, " + toWriteStr + "); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}")
		}
	}
	return nil
}

func (c *GoGenerationContext) generateArraySerializeBinaryMethod(field *CCLField) error {
	targetFieldType := field.Type.GetUnderlyingType()
	isCustomType := targetFieldType.IsCustomTypeModel()
	c.MethodsCode.WriteLine("if err := binary.Write(buf, binary.LittleEndian, uint32(len(m." + field.Name + "))); err != nil {").
		Indent().
		WriteLine("return nil, err").
		Unindent().
		WriteLine("}").
		WriteLine("for _, elem := range m." + field.Name + " {").
		Indent()
	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		c.MethodsCode.WriteLine("if err := binary.Write(buf, binary.LittleEndian, uint32(len(elem))); err != nil {").
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
		c.MethodsCode.WriteLine("if err := binary.Write(buf, binary.LittleEndian, uint32(len(elem))); err != nil {").
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
		c.MethodsCode.WriteLine("if err := binary.Write(buf, binary.LittleEndian, elem.UnixNano()); err != nil {").
			Indent().
			WriteLine("return nil, err").
			Unindent().
			WriteLine("}")
	default:
		if isCustomType {
			currentBytesName := "current_" + field.Name + "Bytes"
			c.MethodsCode.WriteLine(currentBytesName + ", err := elem.SerializeBinary()").
				WriteLine("if err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}").
				WriteLine("if err := binary.Write(buf, binary.LittleEndian, uint32(len(" + currentBytesName + "))); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}").
				WriteLine("if err := binary.Write(buf, binary.LittleEndian, " + currentBytesName + "); err != nil {").
				Indent().
				WriteLine("return nil, err").
				Unindent().
				WriteLine("}")
		} else {
			c.MethodsCode.WriteLine("if err := binary.Write(buf, binary.LittleEndian, elem); err != nil {").
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
	c.MethodsCode.WriteLine("func (m *" + model.Name + ") DeserializeBinary(data []byte) error {").
		Indent().
		// add nil checker or when the len(data) is 0 or (len(data) == 1 and data[0] == 0)
		WriteLine("if m == nil || len(data) == 0 || (len(data) == 1 && data[0] == 0) {").
		Indent().
		WriteLine("return nil").
		Unindent().
		WriteLine("}").
		NewLine().
		WriteLine("buf := bytes.NewReader(data)").
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
	fName := strings.ToLower(string(field.Name[0])) + field.Name[1:]
	fLenName := fName + "Len"
	fNameStrBytes := fName + "StrBytes"
	fNameUnix := fName + "Unix"

	switch fieldTypeName {
	case cclValues.TypeNameString:
		c.MethodsCode.WriteLine("var " + fLenName + " uint32").
			WriteLine("if err := binary.Read(buf, binary.LittleEndian, &" + fLenName + "); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			WriteLine(fNameStrBytes + " := make([]byte, " + fLenName + ")").
			WriteLine("if _, err := buf.Read(" + fNameStrBytes + "); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			WriteLine("m." + field.Name + " = string(" + fNameStrBytes + ")")
	case cclValues.TypeNameBytes:
		c.MethodsCode.WriteLine("var " + fLenName + " uint32").
			WriteLine("if err := binary.Read(buf, binary.LittleEndian, &" + fLenName + "); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			WriteLine("bytesData := make([]byte, " + fLenName + ")").
			WriteLine("if _, err := buf.Read(bytesData); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			WriteLine("m." + field.Name + " = bytesData")
	case cclValues.TypeNameDateTime:
		c.MethodsCode.WriteLine("var " + fNameUnix + " int64").
			WriteLine("if err := binary.Read(buf, binary.LittleEndian, &" + fNameUnix + "); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			WriteLine("m." + field.Name + " = time.Unix(0, " + fNameUnix + ")")
	default:
		if isCustomType {
			// read the length of the next buffer that we need
			// var basicLen uint32
			// if err := binary.Read(buf, binary.LittleEndian, &basicLen); err != nil {
			// 	return err
			// }
			lenVarName := fName + "BytesLen"
			c.MethodsCode.WriteLine("var " + lenVarName + " uint32").
				WriteLine("if err := binary.Read(buf, binary.LittleEndian, &" + lenVarName + "); err != nil {").
				Indent().
				WriteLine("return err").
				Unindent().
				WriteLine("}")

			bytesVarName := fName + "Bytes"
			c.MethodsCode.WriteLine(bytesVarName + " := make([]byte, " + lenVarName + ")").
				WriteLine("if _, err := buf.Read(" + bytesVarName + "); err != nil {").
				Indent().
				WriteLine("return err").
				Unindent().
				WriteLine("}")

			// make sure m.field is not nil ONLY when len(bytesVarName) != 0 and !(len(bytesVarName) == 1 and bytesVarName[0] == 0)
			c.MethodsCode.WriteLine("if m." + field.Name + " == nil && len(" + bytesVarName + ") != 0 && !(len(" + bytesVarName + ") == 1 && " + bytesVarName + "[0] == 0) {").
				Indent().
				WriteLine("m." + field.Name + " = new(" + field.Type.GetName() + ")").
				Unindent().
				WriteLine("}")

			c.MethodsCode.WriteLine("if err := m." + field.Name + ".DeserializeBinary(" + bytesVarName + "); err != nil {").
				Indent().
				WriteLine("return err").
				Unindent().
				WriteLine("}")
		} else {
			if fieldTypeName == cclValues.TypeNameInt {
				// binary.Read does not support int type directly, so we need to read it into an int32 first
				toReadName := "tmp" + field.Name
				c.MethodsCode.WriteLine("var " + toReadName + " int32").
					WriteLine("if err := binary.Read(buf, binary.LittleEndian, &" + toReadName + "); err != nil {").
					Indent().
					WriteLine("return err").
					Unindent().
					WriteLine("}").
					WriteLine("m." + field.Name + " = int(" + toReadName + ")")
				return nil
			}

			// all other types can be read directly
			c.MethodsCode.WriteLine("if err := binary.Read(buf, binary.LittleEndian, &m." + field.Name + "); err != nil {").
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
	isCustomType := targetFieldType.IsCustomTypeModel()
	isPointer := isCustomType //TODO: Find a way to specify this
	fName := strings.ToLower(string(field.Name[0])) + field.Name[1:]
	fLenName := fName + "Len"
	// fNameStrBytes := fName + "StrBytes"
	// fNameUnix := fName + "Unix"
	fieldRealType := targetFieldType.GetName()
	if isPointer {
		fieldRealType = "*" + fieldRealType
	}

	c.MethodsCode.WriteLine("var " + fLenName + " uint32").
		WriteLine("if err := binary.Read(buf, binary.LittleEndian, &" + fLenName + "); err != nil {").
		Indent().
		WriteLine("return err").
		Unindent().
		WriteLine("}").
		WriteLine("m." + field.Name + " = make([]" + fieldRealType + ", " + fLenName + ")").
		WriteLine("for i := uint32(0); i < " + fLenName + "; i++ {").
		Indent()
	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		c.MethodsCode.WriteLine("var elemLen uint32").
			WriteLine("if err := binary.Read(buf, binary.LittleEndian, &elemLen); err != nil {").
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
			WriteLine("m." + field.Name + "[i] = string(elemBytes)")
	case cclValues.TypeNameBytes:
		c.MethodsCode.WriteLine("var elemLen uint32").
			WriteLine("if err := binary.Read(buf, binary.LittleEndian, &elemLen); err != nil {").
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
			WriteLine("m." + field.Name + "[i] = elemBytes")
	case cclValues.TypeNameDateTime:
		c.MethodsCode.WriteLine("var elemUnix int64").
			WriteLine("if err := binary.Read(buf, binary.LittleEndian, &elemUnix); err != nil {").
			Indent().
			WriteLine("return err").
			Unindent().
			WriteLine("}").
			WriteLine("m." + field.Name + "[i] = time.Unix(0, elemUnix)")
	default:
		if isCustomType {
			c.MethodsCode.WriteStr("var elem " + fieldRealType)
			if isPointer {
				c.MethodsCode.AppendLine(" = new(" + targetFieldType.GetName() + ")")
			} else {
				c.MethodsCode.NewLine()
			}
			// we need to read the bytes from the buffer
			// and then deserialize the element
			c.MethodsCode.WriteLine("var elemLen uint32").
				WriteLine("if err := binary.Read(buf, binary.LittleEndian, &elemLen); err != nil {").
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
				WriteLine("m." + field.Name + "[i] = elem")
		} else {
			c.MethodsCode.WriteLine("if err := binary.Read(buf, binary.LittleEndian, &m." + field.Name + "[i]); err != nil {").
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
