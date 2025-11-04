package goGenerator

import (
	"strings"

	"github.com/ccl-lang/ccl/src/core/cclValues"
)

//---------------------------------------------------------

func (c *GoGenerationContext) generateSerializeBinaryMethod(model *CCLModel) error {
	// TODO:
	// This SerializeBinary method generation HAS TO BE MOVED to an attribute handler
	// so users can specify different types of attributes for serialize method generation
	// like JSON, XML, etc.
	c.MethodsCode.WriteString("\nfunc (m *" + model.Name + ") SerializeBinary() ([]byte, error) {\n")

	// handle m is nil by returning []byte(0) and nil
	c.MethodsCode.WriteString("\tif m == nil {\n")
	c.MethodsCode.WriteString("\t\treturn []byte{0}, nil\n")
	c.MethodsCode.WriteString("\t}\n\n")

	c.MethodsCode.WriteString("\tbuf := new(bytes.Buffer)\n\n")
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

	c.MethodsCode.WriteString("\treturn buf.Bytes(), nil\n")
	c.MethodsCode.WriteString("}\n\n")
	return nil
}

func (c *GoGenerationContext) generateFieldSerializeBinaryMethod(field *CCLField) error {
	fieldTypeName := field.Type.GetName()
	isCustomType := c.Options.CCLDefinition.IsCustomType(fieldTypeName)
	switch fieldTypeName {
	case cclValues.TypeNameString:
		c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, uint32(len(m." + field.Name + "))); err != nil {\n")
		c.MethodsCode.WriteString("\t\treturn nil, err\n")
		c.MethodsCode.WriteString("\t}\n")
		c.MethodsCode.WriteString("\tif _, err := buf.WriteString(m." + field.Name + "); err != nil {\n")
		c.MethodsCode.WriteString("\t\treturn nil, err\n")
		c.MethodsCode.WriteString("\t}\n")
	case cclValues.TypeNameBytes:
		c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, uint32(len(m." + field.Name + "))); err != nil {\n")
		c.MethodsCode.WriteString("\t\treturn nil, err\n")
		c.MethodsCode.WriteString("\t}\n")
		c.MethodsCode.WriteString("\tif _, err := buf.Write(m." + field.Name + "); err != nil {\n")
		c.MethodsCode.WriteString("\t\treturn nil, err\n")
		c.MethodsCode.WriteString("\t}\n")
	case cclValues.TypeNameDateTime:
		c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, m." + field.Name + ".UnixNano()); err != nil {\n")
		c.MethodsCode.WriteString("\t\treturn nil, err\n")
		c.MethodsCode.WriteString("\t}\n")
	default:
		if isCustomType {
			currentBytesName := "current_" + field.Name + "Bytes"
			c.MethodsCode.WriteString("\t" + currentBytesName + ", err := m." + field.Name + ".SerializeBinary()\n")
			c.MethodsCode.WriteString("\tif err != nil {\n")
			c.MethodsCode.WriteString("\t\treturn nil, err\n")
			c.MethodsCode.WriteString("\t}\n")
			c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, uint32(len(" + currentBytesName + "))); err != nil {\n")
			c.MethodsCode.WriteString("\t\treturn nil, err\n")
			c.MethodsCode.WriteString("\t}\n")
			c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, " + currentBytesName + "); err != nil {\n")
			c.MethodsCode.WriteString("\t\treturn nil, err\n")
			c.MethodsCode.WriteString("\t}\n")
		} else {
			toWriteStr := "m." + field.Name
			if fieldTypeName == cclValues.TypeNameInt {
				// binary.Write does not support int type directly, so we need to convert it to int32
				toWriteStr = "int32(m." + field.Name + ")"
			}
			c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, " + toWriteStr + "); err != nil {\n")
			c.MethodsCode.WriteString("\t\treturn nil, err\n")
			c.MethodsCode.WriteString("\t}\n")
		}
	}
	return nil
}

func (c *GoGenerationContext) generateArraySerializeBinaryMethod(field *CCLField) error {
	targetFieldType := field.Type.GetUnderlyingType()
	isCustomType := c.Options.CCLDefinition.IsCustomType(targetFieldType.GetName())
	c.MethodsCode.WriteString("\tif err := binary.Write(buf, binary.LittleEndian, uint32(len(m." + field.Name + "))); err != nil {\n")
	c.MethodsCode.WriteString("\t\treturn nil, err\n")
	c.MethodsCode.WriteString("\t}\n")
	c.MethodsCode.WriteString("\tfor _, elem := range m." + field.Name + " {\n")
	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		c.MethodsCode.WriteString("\t\tif err := binary.Write(buf, binary.LittleEndian, uint32(len(elem))); err != nil {\n")
		c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
		c.MethodsCode.WriteString("\t\t}\n")
		c.MethodsCode.WriteString("\t\tif _, err := buf.WriteString(elem); err != nil {\n")
		c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
		c.MethodsCode.WriteString("\t\t}\n")
	case cclValues.TypeNameBytes:
		c.MethodsCode.WriteString("\t\tif err := binary.Write(buf, binary.LittleEndian, uint32(len(elem))); err != nil {\n")
		c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
		c.MethodsCode.WriteString("\t\t}\n")
		c.MethodsCode.WriteString("\t\tif _, err := buf.Write(elem); err != nil {\n")
		c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
		c.MethodsCode.WriteString("\t\t}\n")
	case cclValues.TypeNameDateTime:
		c.MethodsCode.WriteString("\t\tif err := binary.Write(buf, binary.LittleEndian, elem.UnixNano()); err != nil {\n")
		c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
		c.MethodsCode.WriteString("\t\t}\n")
	default:
		if isCustomType {
			currentBytesName := "current_" + field.Name + "Bytes"
			c.MethodsCode.WriteString("\t\t" + currentBytesName + ", err := elem.SerializeBinary()\n")
			c.MethodsCode.WriteString("\t\tif err != nil {\n")
			c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
			c.MethodsCode.WriteString("\t\t}\n")
			c.MethodsCode.WriteString("\t\tif err := binary.Write(buf, binary.LittleEndian, uint32(len(" + currentBytesName + "))); err != nil {\n")
			c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
			c.MethodsCode.WriteString("\t\t}\n")
			c.MethodsCode.WriteString("\t\tif err := binary.Write(buf, binary.LittleEndian, " + currentBytesName + "); err != nil {\n")
			c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
			c.MethodsCode.WriteString("\t\t}\n")
		} else {
			c.MethodsCode.WriteString("\t\tif err := binary.Write(buf, binary.LittleEndian, elem); err != nil {\n")
			c.MethodsCode.WriteString("\t\t\treturn nil, err\n")
			c.MethodsCode.WriteString("\t\t}\n")
		}
	}

	c.MethodsCode.WriteString("\t}\n")
	return nil
}

func (c *GoGenerationContext) generateDeserializeBinaryMethod(model *CCLModel) error {
	c.MethodsCode.WriteString("func (m *" + model.Name + ") DeserializeBinary(data []byte) error {\n")
	// add nil checker or when the len(data) is 0 or (len(data) == 1 and data[0] == 0)
	c.MethodsCode.WriteString("\tif m == nil || len(data) == 0 || (len(data) == 1 && data[0] == 0) {\n")
	c.MethodsCode.WriteString("\t\treturn nil\n")
	c.MethodsCode.WriteString("\t}\n\n")
	c.MethodsCode.WriteString("\tbuf := bytes.NewReader(data)\n\n")

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

	c.MethodsCode.WriteString("\treturn nil\n")
	c.MethodsCode.WriteString("}\n\n")
	return nil
}

func (c *GoGenerationContext) generateFieldDeserializeBinaryMethod(field *CCLField) error {
	fieldTypeName := field.Type.GetName()
	isCustomType := c.Options.CCLDefinition.IsCustomType(fieldTypeName)
	// isPointer := isCustomType //TODO: Find a way to specify this
	fName := strings.ToLower(string(field.Name[0])) + field.Name[1:]
	fLenName := fName + "Len"
	fNameStrBytes := fName + "StrBytes"
	fNameUnix := fName + "Unix"

	switch fieldTypeName {
	case cclValues.TypeNameString:
		c.MethodsCode.WriteString("\tvar " + fLenName + " uint32\n")
		c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &" + fLenName + "); err != nil {\n")
		c.MethodsCode.WriteString("\t\treturn err\n")
		c.MethodsCode.WriteString("\t}\n")
		c.MethodsCode.WriteString("\t" + fNameStrBytes + " := make([]byte, " + fLenName + ")\n")
		c.MethodsCode.WriteString("\tif _, err := buf.Read(" + fNameStrBytes + "); err != nil {\n")
		c.MethodsCode.WriteString("\t\treturn err\n")
		c.MethodsCode.WriteString("\t}\n")
		c.MethodsCode.WriteString("\tm." + field.Name + " = string(" + fNameStrBytes + ")\n")
	case cclValues.TypeNameBytes:
		c.MethodsCode.WriteString("\tvar " + fLenName + " uint32\n")
		c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &" + fLenName + "); err != nil {\n")
		c.MethodsCode.WriteString("\t\treturn err\n")
		c.MethodsCode.WriteString("\t}\n")
		c.MethodsCode.WriteString("\tbytesData := make([]byte, " + fLenName + ")\n")
		c.MethodsCode.WriteString("\tif _, err := buf.Read(bytesData); err != nil {\n")
		c.MethodsCode.WriteString("\t\treturn err\n")
		c.MethodsCode.WriteString("\t}\n")
		c.MethodsCode.WriteString("\tm." + field.Name + " = bytesData\n")
	case cclValues.TypeNameDateTime:
		c.MethodsCode.WriteString("\tvar " + fNameUnix + " int64\n")
		c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &" + fNameUnix + "); err != nil {\n")
		c.MethodsCode.WriteString("\t\treturn err\n")
		c.MethodsCode.WriteString("\t}\n")
		c.MethodsCode.WriteString("\tm." + field.Name + " = time.Unix(0, " + fNameUnix + ")\n")
	default:
		if isCustomType {
			// read the length of the next buffer that we need
			// var basicLen uint32
			// if err := binary.Read(buf, binary.LittleEndian, &basicLen); err != nil {
			// 	return err
			// }
			lenVarName := fName + "BytesLen"
			c.MethodsCode.WriteString("\tvar " + lenVarName + " uint32\n")
			c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &" + lenVarName + "); err != nil {\n")
			c.MethodsCode.WriteString("\t\treturn err\n")
			c.MethodsCode.WriteString("\t}\n")

			bytesVarName := fName + "Bytes"
			c.MethodsCode.WriteString("\t" + bytesVarName + " := make([]byte, " + lenVarName + ")\n")
			c.MethodsCode.WriteString("\tif _, err := buf.Read(" + bytesVarName + "); err != nil {\n")
			c.MethodsCode.WriteString("\t\treturn err\n")
			c.MethodsCode.WriteString("\t}\n")

			// make sure m.field is not nil ONLY when len(bytesVarName) != 0 and !(len(bytesVarName) == 1 and bytesVarName[0] == 0)
			c.MethodsCode.WriteString("\tif m." + field.Name + " == nil && len(" + bytesVarName + ") != 0 && !(len(" + bytesVarName + ") == 1 && " + bytesVarName + "[0] == 0) {\n")
			c.MethodsCode.WriteString("\t\tm." + field.Name + " = new(" + field.Type.GetName() + ")\n")
			c.MethodsCode.WriteString("\t}\n")

			c.MethodsCode.WriteString("\tif err := m." + field.Name + ".DeserializeBinary(" + bytesVarName + "); err != nil {\n")
			c.MethodsCode.WriteString("\t\treturn err\n")
			c.MethodsCode.WriteString("\t}\n")
		} else {
			if fieldTypeName == cclValues.TypeNameInt {
				// binary.Read does not support int type directly, so we need to read it into an int32 first
				toReadName := "tmp" + field.Name
				c.MethodsCode.WriteString("\tvar " + toReadName + " int32\n")
				c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &" + toReadName + "); err != nil {\n")
				c.MethodsCode.WriteString("\t\treturn err\n")
				c.MethodsCode.WriteString("\t}\n")
				c.MethodsCode.WriteString("\tm." + field.Name + " = int(" + toReadName + ")\n")
				return nil
			}

			// all other types can be read directly
			c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &m." + field.Name + "); err != nil {\n")
			c.MethodsCode.WriteString("\t\treturn err\n")
			c.MethodsCode.WriteString("\t}\n")
		}
	}
	return nil
}

func (c *GoGenerationContext) generateArrayDeserializeBinaryMethod(field *CCLField) error {
	targetFieldType := field.Type.GetUnderlyingType()
	isCustomType := c.Options.CCLDefinition.IsCustomType(targetFieldType.GetName())
	isPointer := isCustomType //TODO: Find a way to specify this
	fName := strings.ToLower(string(field.Name[0])) + field.Name[1:]
	fLenName := fName + "Len"
	// fNameStrBytes := fName + "StrBytes"
	// fNameUnix := fName + "Unix"
	fieldRealType := targetFieldType.GetName()
	if isPointer {
		fieldRealType = "*" + fieldRealType
	}

	c.MethodsCode.WriteString("\tvar " + fLenName + " uint32\n")
	c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &" + fLenName + "); err != nil {\n")
	c.MethodsCode.WriteString("\t\treturn err\n")
	c.MethodsCode.WriteString("\t}\n")
	c.MethodsCode.WriteString("\tm." + field.Name + " = make([]" + fieldRealType + ", " + fLenName + ")\n")
	c.MethodsCode.WriteString("\tfor i := uint32(0); i < " + fLenName + "; i++ {\n")
	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		c.MethodsCode.WriteString("\t\tvar elemLen uint32\n")
		c.MethodsCode.WriteString("\t\tif err := binary.Read(buf, binary.LittleEndian, &elemLen); err != nil {\n")
		c.MethodsCode.WriteString("\t\t\treturn err\n")
		c.MethodsCode.WriteString("\t\t}\n")
		c.MethodsCode.WriteString("\t\telemBytes := make([]byte, elemLen)\n")
		c.MethodsCode.WriteString("\t\tif _, err := buf.Read(elemBytes); err != nil {\n")
		c.MethodsCode.WriteString("\t\t\treturn err\n")
		c.MethodsCode.WriteString("\t\t}\n")
		c.MethodsCode.WriteString("\t\tm." + field.Name + "[i] = string(elemBytes)\n")
	case cclValues.TypeNameBytes:
		c.MethodsCode.WriteString("\t\tvar elemLen uint32\n")
		c.MethodsCode.WriteString("\t\tif err := binary.Read(buf, binary.LittleEndian, &elemLen); err != nil {\n")
		c.MethodsCode.WriteString("\t\t\treturn err\n")
		c.MethodsCode.WriteString("\t\t}\n")
		c.MethodsCode.WriteString("\t\telemBytes := make([]byte, elemLen)\n")
		c.MethodsCode.WriteString("\t\tif _, err := buf.Read(elemBytes); err != nil {\n")
		c.MethodsCode.WriteString("\t\t\treturn err\n")
		c.MethodsCode.WriteString("\t\t}\n")
		c.MethodsCode.WriteString("\t\tm." + field.Name + "[i] = elemBytes\n")
	case cclValues.TypeNameDateTime:
		c.MethodsCode.WriteString("\t\tvar elemUnix int64\n")
		c.MethodsCode.WriteString("\tif err := binary.Read(buf, binary.LittleEndian, &elemUnix); err != nil {\n")
		c.MethodsCode.WriteString("\t\treturn err\n")
		c.MethodsCode.WriteString("\t}\n")
		c.MethodsCode.WriteString("\t\tm." + field.Name + "[i] = time.Unix(0, elemUnix)\n")
	default:
		if isCustomType {
			c.MethodsCode.WriteString("\t\tvar elem " + fieldRealType)
			if isPointer {
				c.MethodsCode.WriteString(" = new(" + targetFieldType.GetName() + ")\n")
			} else {
				c.MethodsCode.WriteString("\n")
			}
			// we need to read the bytes from the buffer
			// and then deserialize the element
			c.MethodsCode.WriteString("\t\tvar elemLen uint32\n")
			c.MethodsCode.WriteString("\t\tif err := binary.Read(buf, binary.LittleEndian, &elemLen); err != nil {\n")
			c.MethodsCode.WriteString("\t\t\treturn err\n")
			c.MethodsCode.WriteString("\t\t}\n")
			c.MethodsCode.WriteString("\t\telemBytes := make([]byte, elemLen)\n")
			c.MethodsCode.WriteString("\t\tif _, err := buf.Read(elemBytes); err != nil {\n")
			c.MethodsCode.WriteString("\t\t\treturn err\n")
			c.MethodsCode.WriteString("\t\t}\n")
			c.MethodsCode.WriteString("\t\tif err := elem.DeserializeBinary(elemBytes); err != nil {\n")
			c.MethodsCode.WriteString("\t\t\treturn err\n")
			c.MethodsCode.WriteString("\t\t}\n")
			c.MethodsCode.WriteString("\t\tm." + field.Name + "[i] = elem\n")
		} else {
			c.MethodsCode.WriteString("\t\tif err := binary.Read(buf, binary.LittleEndian, &m." + field.Name + "[i]); err != nil {\n")
			c.MethodsCode.WriteString("\t\t\treturn err\n")
			c.MethodsCode.WriteString("\t\t}\n")
		}
	}

	c.MethodsCode.WriteString("\t}\n")
	return nil
}

//---------------------------------------------------------
