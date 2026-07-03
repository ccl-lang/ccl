package csGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *CSharpGenerationContext) generateSerializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	builder.WriteLine("public byte[] SerializeBinary()").
		WriteLine("{").
		Indent().
		WriteLine("using (var ms = new MemoryStream())").
		WriteLine("{").
		Indent().
		WriteLine("using (var writer = new BinaryWriter(ms, Encoding.UTF8, true))").
		WriteLine("{").
		Indent()

	for i := range model.Fields {
		field := model.Fields[i]
		if field.IsArray() {
			if err := c.generateArraySerializeBinary(field, builder); err != nil {
				return err
			}
			continue
		}

		if err := c.generateFieldSerializeBinary(field, builder); err != nil {
			return err
		}
	}

	builder.Unindent().
		WriteLine("}").
		WriteLine("return ms.ToArray();").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *CSharpGenerationContext) generateFieldSerializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) error {
	fieldName := "this." + caseUtils.ToPascalCase(field.Name)
	fieldWrite, err := c.csharpBinaryWriteExpression(field.Type, fieldName)
	if err != nil {
		return err
	}
	builder.MapVarPairs(
		"field", fieldName,
		"fieldWrite", fieldWrite,
	)
	defer builder.UnmapVar(
		"field",
		"fieldWrite",
	)

	switch csharpStorageTypeName(field.Type) {
	case cclValues.TypeNameString:
		builder.WriteLine("{").
			Indent().
			LineD(`var strBytes = Encoding.UTF8.GetBytes($field ?? "");`).
			WriteLine("writer.Write((uint)strBytes.Length);").
			WriteLine("writer.Write(strBytes);").
			Unindent().
			WriteLine("}")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.LineD("writer.Write($fieldWrite);")
	case cclValues.TypeNameInt8:
		builder.LineD("writer.Write((sbyte)$fieldWrite);")
	case cclValues.TypeNameInt16:
		builder.LineD("writer.Write($fieldWrite);")
	case cclValues.TypeNameInt64:
		builder.LineD("writer.Write($fieldWrite);")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.LineD("writer.Write($fieldWrite);")
	case cclValues.TypeNameUint8:
		builder.LineD("writer.Write($fieldWrite);")
	case cclValues.TypeNameUint16:
		builder.LineD("writer.Write($fieldWrite);")
	case cclValues.TypeNameUint64:
		builder.LineD("writer.Write($fieldWrite);")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.LineD("writer.Write($field);")
	case cclValues.TypeNameFloat64:
		builder.LineD("writer.Write($field);")
	case cclValues.TypeNameBool:
		builder.LineD("writer.Write((byte)($field ? 1 : 0));")
	case cclValues.TypeNameBytes:
		builder.LineD("writer.Write((uint)($field?.Length ?? 0));").
			LineD("if ($field != null) writer.Write($field);")
	case cclValues.TypeNameDateTime:
		builder.LineD("writer.Write($field);")
	default:
		if field.IsCustomTypeModel() {
			builder.LineD("if ($field != null)").
				WriteLine("{").
				Indent().
				LineD("var customBytes = $field.SerializeBinary();").
				WriteLine("writer.Write((uint)customBytes.Length);").
				WriteLine("writer.Write(customBytes);").
				Unindent().
				WriteLine("}").
				WriteLine("else").
				WriteLine("{").
				Indent().
				WriteLine("writer.Write((uint)0);").
				Unindent().
				WriteLine("}")
		}
	}
	return nil
}

func (c *CSharpGenerationContext) generateArraySerializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) error {
	fieldName := "this." + caseUtils.ToPascalCase(field.Name)
	targetFieldType := field.Type.GetUnderlyingType()
	itemWrite, err := c.csharpBinaryWriteExpression(targetFieldType, "item")
	if err != nil {
		return err
	}
	builder.MapVarPairs(
		"field", fieldName,
		"itemWrite", itemWrite,
	)
	defer builder.UnmapVar(
		"field",
		"itemWrite",
	)

	builder.LineD("writer.Write((uint)($field?.Count ?? 0));").
		LineD("if ($field != null)").
		WriteLine("{").
		Indent().
		LineD("foreach (var item in $field)").
		WriteLine("{").
		Indent()

	switch csharpStorageTypeName(targetFieldType) {
	case cclValues.TypeNameString:
		builder.WriteLine("var strBytes = Encoding.UTF8.GetBytes(item ?? \"\");")
		builder.WriteLine("writer.Write((uint)strBytes.Length);")
		builder.WriteLine("writer.Write(strBytes);")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.LineD("writer.Write($itemWrite);")
	case cclValues.TypeNameInt8:
		builder.LineD("writer.Write((sbyte)$itemWrite);")
	case cclValues.TypeNameInt16:
		builder.LineD("writer.Write($itemWrite);")
	case cclValues.TypeNameInt64:
		builder.LineD("writer.Write($itemWrite);")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.LineD("writer.Write($itemWrite);")
	case cclValues.TypeNameUint8:
		builder.LineD("writer.Write($itemWrite);")
	case cclValues.TypeNameUint16:
		builder.LineD("writer.Write($itemWrite);")
	case cclValues.TypeNameUint64:
		builder.LineD("writer.Write($itemWrite);")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.WriteLine("writer.Write(item);")
	case cclValues.TypeNameFloat64:
		builder.WriteLine("writer.Write(item);")
	case cclValues.TypeNameBool:
		builder.WriteLine("writer.Write((byte)(item ? 1 : 0));")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("if (item != null)")
			builder.WriteLine("{")
			builder.Indent()
			builder.WriteLine("var customBytes = item.SerializeBinary();")
			builder.WriteLine("writer.Write((uint)customBytes.Length);")
			builder.WriteLine("writer.Write(customBytes);")
			builder.Unindent()
			builder.WriteLine("}")
			builder.WriteLine("else")
			builder.WriteLine("{")
			builder.Indent()
			builder.WriteLine("writer.Write((uint)0);")
			builder.Unindent()
			builder.WriteLine("}")
		}
	}

	builder.Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}")
	return nil
}

func (c *CSharpGenerationContext) generateDeserializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	strictBinaryParsing, err := c.UsesStrictBinaryParsing(CurrentLanguage, model)
	if err != nil {
		return err
	}
	binaryParseFallback := "result"
	if strictBinaryParsing {
		binaryParseFallback = "null"
	}

	builder.MapVarPairs(
		"model", model.Name,
		"binaryParseFallback", binaryParseFallback,
	)
	defer builder.UnmapVar(
		"model",
		"binaryParseFallback",
	)

	builder.LineD("public static $model DeserializeBinary(byte[] data)").
		WriteLine("{").
		Indent().
		WriteLine("if (data == null || data.Length == 0) return null;").
		LineD("var result = new $model();").
		WriteLine("try").
		WriteLine("{").
		Indent().
		WriteLine("using (var ms = new MemoryStream(data))").
		WriteLine("{").
		Indent().
		WriteLine("using (var reader = new BinaryReader(ms, Encoding.UTF8))").
		WriteLine("{").
		Indent()

	for _, field := range model.Fields {
		if field.IsArray() {
			if err := c.generateArrayDeserializeBinary(field, builder); err != nil {
				return err
			}
			continue
		}

		if err := c.generateFieldDeserializeBinary(field, builder); err != nil {
			return err
		}
	}

	builder.Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}").
		WriteLine("catch (EndOfStreamException)").
		WriteLine("{").
		Indent().
		LineD("return $binaryParseFallback;").
		Unindent().
		WriteLine("}").
		WriteLine("catch (ArgumentException)").
		WriteLine("{").
		Indent().
		LineD("return $binaryParseFallback;").
		Unindent().
		WriteLine("}").
		WriteLine("return result;").
		Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *CSharpGenerationContext) generateFieldDeserializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) error {
	fieldName := caseUtils.ToPascalCase(field.Name)
	resultField := "result." + fieldName
	builder.MapVarPairs(
		"field", resultField,
		"type", field.Type.GetName(),
	)
	defer builder.UnmapVar(
		"field",
		"type",
	)

	if field.Type.IsCustomTypeEnum() {
		baseType := csharpStorageTypeName(field.Type)
		readerMethod := c.csharpReaderMethod(baseType)
		enumName, err := c.getCSharpEnumTypeName(field.Type.GetDefinition().GetEnumDefinition())
		if err != nil {
			return err
		}
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, enumBinarySize(baseType))
		builder.LineD("$field = (" + enumName + ")reader." + readerMethod + "();")
		builder.NewLine()
		return nil
	}

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("{").
			Indent()
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
		builder.WriteLine("var len = reader.ReadUInt32();").
			LineD("if (len > ms.Length - ms.Position) return $binaryParseFallback;").
			WriteLine("var bytes = reader.ReadBytes((int)len);").
			LineD("$field = Encoding.UTF8.GetString(bytes);").
			Unindent().
			WriteLine("}")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("$field = reader.ReadInt32();")
	case cclValues.TypeNameInt8:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "1")
		builder.LineD("$field = reader.ReadSByte();")
	case cclValues.TypeNameInt16:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "2")
		builder.LineD("$field = reader.ReadInt16();")
	case cclValues.TypeNameInt64:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "8")
		builder.LineD("$field = reader.ReadInt64();")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("$field = reader.ReadUInt32();")
	case cclValues.TypeNameUint8:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "1")
		builder.LineD("$field = reader.ReadByte();")
	case cclValues.TypeNameUint16:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "2")
		builder.LineD("$field = reader.ReadUInt16();")
	case cclValues.TypeNameUint64:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "8")
		builder.LineD("$field = reader.ReadUInt64();")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("$field = reader.ReadSingle();")
	case cclValues.TypeNameFloat64:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "8")
		builder.LineD("$field = reader.ReadDouble();")
	case cclValues.TypeNameBool:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "1")
		builder.LineD("$field = reader.ReadByte() != 0;")
	case cclValues.TypeNameBytes:
		builder.WriteLine("{").
			Indent()
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
		builder.WriteLine("var len = reader.ReadUInt32();").
			LineD("if (len > ms.Length - ms.Position) return $binaryParseFallback;").
			LineD("$field = reader.ReadBytes((int)len);").
			Unindent().
			WriteLine("}")
	case cclValues.TypeNameDateTime:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "8")
		builder.LineD("$field = reader.ReadInt64();")
	default:
		if field.IsCustomTypeModel() {
			builder.WriteLine("{").
				Indent()
			c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
			builder.WriteLine("var len = reader.ReadUInt32();").
				WriteLine("if (len > 0)").
				WriteLine("{").
				Indent().
				LineD("if (len > ms.Length - ms.Position) return $binaryParseFallback;").
				WriteLine("var bytes = reader.ReadBytes((int)len);").
				LineD("$field = $type.DeserializeBinary(bytes);").
				Unindent().
				WriteLine("}").
				Unindent().
				WriteLine("}")
		}
	}
	return nil
}

func (c *CSharpGenerationContext) generateArrayDeserializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) error {
	fieldName := caseUtils.ToPascalCase(field.Name)
	resultField := "result." + fieldName
	targetFieldType := field.Type.GetUnderlyingType()
	fieldType, err := c.getCSharpType(field)
	if err != nil {
		return err
	}
	builder.MapVarPairs(
		"field", resultField,
		"fieldType", fieldType,
		"type", targetFieldType.GetName(),
	)
	defer builder.UnmapVar(
		"field",
		"fieldType",
		"type",
	)

	builder.WriteLine("{").
		Indent()
	c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
	builder.WriteLine("var len = reader.ReadUInt32();").
		LineD("$field = new $fieldType();").
		WriteLine("for (var i = 0; i < len; i++)").
		WriteLine("{").
		Indent()

	if targetFieldType.IsCustomTypeEnum() {
		baseType := csharpStorageTypeName(targetFieldType)
		readerMethod := c.csharpReaderMethod(baseType)
		enumName, err := c.getCSharpEnumTypeName(targetFieldType.GetDefinition().GetEnumDefinition())
		if err != nil {
			return err
		}
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, enumBinarySize(baseType))
		builder.LineD("$field.Add((" + enumName + ")reader." + readerMethod + "());")
		builder.Unindent().
			WriteLine("}").
			Unindent().
			WriteLine("}")
		return nil
	}

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
		builder.WriteLine("var itemLen = reader.ReadUInt32();")
		builder.LineD("if (itemLen > ms.Length - ms.Position) return $binaryParseFallback;")
		builder.WriteLine("var bytes = reader.ReadBytes((int)itemLen);")
		builder.LineD("$field.Add(Encoding.UTF8.GetString(bytes));")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("$field.Add(reader.ReadInt32());")
	case cclValues.TypeNameInt8:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "1")
		builder.LineD("$field.Add(reader.ReadSByte());")
	case cclValues.TypeNameInt16:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "2")
		builder.LineD("$field.Add(reader.ReadInt16());")
	case cclValues.TypeNameInt64:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "8")
		builder.LineD("$field.Add(reader.ReadInt64());")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("$field.Add(reader.ReadUInt32());")
	case cclValues.TypeNameUint8:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "1")
		builder.LineD("$field.Add(reader.ReadByte());")
	case cclValues.TypeNameUint16:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "2")
		builder.LineD("$field.Add(reader.ReadUInt16());")
	case cclValues.TypeNameUint64:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "8")
		builder.LineD("$field.Add(reader.ReadUInt64());")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
		builder.LineD("$field.Add(reader.ReadSingle());")
	case cclValues.TypeNameFloat64:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "8")
		builder.LineD("$field.Add(reader.ReadDouble());")
	case cclValues.TypeNameBool:
		c.generateCSharpBinaryDeserializeBoundsCheck(builder, "1")
		builder.LineD("$field.Add(reader.ReadByte() != 0);")
	default:
		if targetFieldType.IsCustomTypeModel() {
			c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
			builder.WriteLine("var itemLen = reader.ReadUInt32();").
				WriteLine("if (itemLen > 0)").
				WriteLine("{").
				Indent().
				LineD("if (itemLen > ms.Length - ms.Position) return $binaryParseFallback;").
				WriteLine("var bytes = reader.ReadBytes((int)itemLen);").
				LineD("$field.Add($type.DeserializeBinary(bytes));").
				Unindent().
				WriteLine("}").
				WriteLine("else").
				WriteLine("{").
				Indent().
				LineD("$field.Add(null);").
				Unindent().
				WriteLine("}")
		}
	}

	builder.Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}")
	return nil
}

func (c *CSharpGenerationContext) generateCSharpBinaryDeserializeBoundsCheck(
	builder *codeBuilder.CodeBuilder,
	requiredBytes string,
) {
	builder.MapVarPairs("requiredBytes", requiredBytes)
	builder.LineD("if (ms.Length - ms.Position < $requiredBytes) return $binaryParseFallback;")
	builder.UnmapVar("requiredBytes")
}
