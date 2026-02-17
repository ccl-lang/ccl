package csGenerator

import (
	"github.com/ccl-lang/ccl/src/core/cclUtils"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func (c *CSharpGenerationContext) generateSerializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	builder.WriteLine("public byte[] SerializeBinary()")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("using (var ms = new MemoryStream())")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("using (var writer = new BinaryWriter(ms, Encoding.UTF8, true))")
	builder.WriteLine("{")
	builder.Indent()

	for i := range model.Fields {
		field := model.Fields[i]
		if field.IsArray() {
			c.generateArraySerializeBinary(field, builder)
			continue
		}

		c.generateFieldSerializeBinary(field, builder)
	}

	builder.Unindent()
	builder.WriteLine("}")
	builder.WriteLine("return ms.ToArray();")
	builder.Unindent()
	builder.WriteLine("}")
	builder.Unindent()
	builder.WriteLine("}")
	builder.NewLine()
	return nil
}

func (c *CSharpGenerationContext) generateFieldSerializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	fieldName := "this." + cclUtils.ToPascalCase(field.Name)

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("{")
		builder.Indent()
		builder.WriteLine("var strBytes = Encoding.UTF8.GetBytes(" + fieldName + " ?? \"\");")
		builder.WriteLine("writer.Write((uint)strBytes.Length);")
		builder.WriteLine("writer.Write(strBytes);")
		builder.Unindent()
		builder.WriteLine("}")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine("writer.Write(" + fieldName + ");")
	case cclValues.TypeNameInt8:
		builder.WriteLine("writer.Write((sbyte)" + fieldName + ");")
	case cclValues.TypeNameInt16:
		builder.WriteLine("writer.Write(" + fieldName + ");")
	case cclValues.TypeNameInt64:
		builder.WriteLine("writer.Write(" + fieldName + ");")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine("writer.Write(" + fieldName + ");")
	case cclValues.TypeNameUint8:
		builder.WriteLine("writer.Write(" + fieldName + ");")
	case cclValues.TypeNameUint16:
		builder.WriteLine("writer.Write(" + fieldName + ");")
	case cclValues.TypeNameUint64:
		builder.WriteLine("writer.Write(" + fieldName + ");")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.WriteLine("writer.Write(" + fieldName + ");")
	case cclValues.TypeNameFloat64:
		builder.WriteLine("writer.Write(" + fieldName + ");")
	case cclValues.TypeNameBool:
		builder.WriteLine("writer.Write((byte)(" + fieldName + " ? 1 : 0));")
	case cclValues.TypeNameBytes:
		builder.WriteLine("writer.Write((uint)(" + fieldName + "?.Length ?? 0));")
		builder.WriteLine("if (" + fieldName + " != null) writer.Write(" + fieldName + ");")
	case cclValues.TypeNameDateTime:
		builder.WriteLine("writer.Write(" + fieldName + ");")
	default:
		if field.IsCustomTypeModel() {
			builder.WriteLine("if (" + fieldName + " != null)")
			builder.WriteLine("{")
			builder.Indent()
			builder.WriteLine("var customBytes = " + fieldName + ".SerializeBinary();")
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
}

func (c *CSharpGenerationContext) generateArraySerializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	fieldName := "this." + cclUtils.ToPascalCase(field.Name)
	targetFieldType := field.Type.GetUnderlyingType()

	builder.WriteLine("writer.Write((uint)(" + fieldName + "?.Count ?? 0));")
	builder.WriteLine("if (" + fieldName + " != null)")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("foreach (var item in " + fieldName + ")")
	builder.WriteLine("{")
	builder.Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("var strBytes = Encoding.UTF8.GetBytes(item ?? \"\");")
		builder.WriteLine("writer.Write((uint)strBytes.Length);")
		builder.WriteLine("writer.Write(strBytes);")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine("writer.Write(item);")
	case cclValues.TypeNameInt8:
		builder.WriteLine("writer.Write((sbyte)item);")
	case cclValues.TypeNameInt16:
		builder.WriteLine("writer.Write(item);")
	case cclValues.TypeNameInt64:
		builder.WriteLine("writer.Write(item);")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine("writer.Write(item);")
	case cclValues.TypeNameUint8:
		builder.WriteLine("writer.Write(item);")
	case cclValues.TypeNameUint16:
		builder.WriteLine("writer.Write(item);")
	case cclValues.TypeNameUint64:
		builder.WriteLine("writer.Write(item);")
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

	builder.Unindent()
	builder.WriteLine("}")
	builder.Unindent()
	builder.WriteLine("}")
}

func (c *CSharpGenerationContext) generateDeserializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	builder.WriteLine("public static " + model.Name + " DeserializeBinary(byte[] data)")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("if (data == null || data.Length == 0) return null;")
	builder.WriteLine("var result = new " + model.Name + "();")
	builder.WriteLine("using (var ms = new MemoryStream(data))")
	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("using (var reader = new BinaryReader(ms, Encoding.UTF8))")
	builder.WriteLine("{")
	builder.Indent()

	for _, field := range model.Fields {
		if field.IsArray() {
			c.generateArrayDeserializeBinary(field, builder)
			continue
		}

		c.generateFieldDeserializeBinary(field, builder)
	}

	builder.Unindent()
	builder.WriteLine("}")
	builder.Unindent()
	builder.WriteLine("}")
	builder.WriteLine("return result;")
	builder.Unindent()
	builder.WriteLine("}")
	builder.NewLine()
	return nil
}

func (c *CSharpGenerationContext) generateFieldDeserializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	fieldName := cclUtils.ToPascalCase(field.Name)
	resultField := "result." + fieldName

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("{")
		builder.Indent()
		builder.WriteLine("var len = reader.ReadUInt32();")
		builder.WriteLine("var bytes = reader.ReadBytes((int)len);")
		builder.WriteLine(resultField + " = Encoding.UTF8.GetString(bytes);")
		builder.Unindent()
		builder.WriteLine("}")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine(resultField + " = reader.ReadInt32();")
	case cclValues.TypeNameInt8:
		builder.WriteLine(resultField + " = reader.ReadSByte();")
	case cclValues.TypeNameInt16:
		builder.WriteLine(resultField + " = reader.ReadInt16();")
	case cclValues.TypeNameInt64:
		builder.WriteLine(resultField + " = reader.ReadInt64();")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine(resultField + " = reader.ReadUInt32();")
	case cclValues.TypeNameUint8:
		builder.WriteLine(resultField + " = reader.ReadByte();")
	case cclValues.TypeNameUint16:
		builder.WriteLine(resultField + " = reader.ReadUInt16();")
	case cclValues.TypeNameUint64:
		builder.WriteLine(resultField + " = reader.ReadUInt64();")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.WriteLine(resultField + " = reader.ReadSingle();")
	case cclValues.TypeNameFloat64:
		builder.WriteLine(resultField + " = reader.ReadDouble();")
	case cclValues.TypeNameBool:
		builder.WriteLine(resultField + " = reader.ReadByte() != 0;")
	case cclValues.TypeNameBytes:
		builder.WriteLine("{")
		builder.Indent()
		builder.WriteLine("var len = reader.ReadUInt32();")
		builder.WriteLine(resultField + " = reader.ReadBytes((int)len);")
		builder.Unindent()
		builder.WriteLine("}")
	case cclValues.TypeNameDateTime:
		builder.WriteLine(resultField + " = reader.ReadInt64();")
	default:
		if field.IsCustomTypeModel() {
			builder.WriteLine("{")
			builder.Indent()
			builder.WriteLine("var len = reader.ReadUInt32();")
			builder.WriteLine("if (len > 0)")
			builder.WriteLine("{")
			builder.Indent()
			builder.WriteLine("var bytes = reader.ReadBytes((int)len);")
			builder.WriteLine(resultField + " = " + field.Type.GetName() + ".DeserializeBinary(bytes);")
			builder.Unindent()
			builder.WriteLine("}")
			builder.Unindent()
			builder.WriteLine("}")
		}
	}
}

func (c *CSharpGenerationContext) generateArrayDeserializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	fieldName := cclUtils.ToPascalCase(field.Name)
	resultField := "result." + fieldName
	targetFieldType := field.Type.GetUnderlyingType()

	builder.WriteLine("{")
	builder.Indent()
	builder.WriteLine("var len = reader.ReadUInt32();")
	builder.WriteLine(resultField + " = new " + c.getCSharpType(field) + "();")
	builder.WriteLine("for (var i = 0; i < len; i++)")
	builder.WriteLine("{")
	builder.Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("var itemLen = reader.ReadUInt32();")
		builder.WriteLine("var bytes = reader.ReadBytes((int)itemLen);")
		builder.WriteLine(resultField + ".Add(Encoding.UTF8.GetString(bytes));")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine(resultField + ".Add(reader.ReadInt32());")
	case cclValues.TypeNameInt8:
		builder.WriteLine(resultField + ".Add(reader.ReadSByte());")
	case cclValues.TypeNameInt16:
		builder.WriteLine(resultField + ".Add(reader.ReadInt16());")
	case cclValues.TypeNameInt64:
		builder.WriteLine(resultField + ".Add(reader.ReadInt64());")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine(resultField + ".Add(reader.ReadUInt32());")
	case cclValues.TypeNameUint8:
		builder.WriteLine(resultField + ".Add(reader.ReadByte());")
	case cclValues.TypeNameUint16:
		builder.WriteLine(resultField + ".Add(reader.ReadUInt16());")
	case cclValues.TypeNameUint64:
		builder.WriteLine(resultField + ".Add(reader.ReadUInt64());")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.WriteLine(resultField + ".Add(reader.ReadSingle());")
	case cclValues.TypeNameFloat64:
		builder.WriteLine(resultField + ".Add(reader.ReadDouble());")
	case cclValues.TypeNameBool:
		builder.WriteLine(resultField + ".Add(reader.ReadByte() != 0);")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("var itemLen = reader.ReadUInt32();")
			builder.WriteLine("if (itemLen > 0)")
			builder.WriteLine("{")
			builder.Indent()
			builder.WriteLine("var bytes = reader.ReadBytes((int)itemLen);")
			builder.WriteLine(resultField + ".Add(" + targetFieldType.GetName() + ".DeserializeBinary(bytes));")
			builder.Unindent()
			builder.WriteLine("}")
			builder.WriteLine("else")
			builder.WriteLine("{")
			builder.Indent()
			builder.WriteLine(resultField + ".Add(null);")
			builder.Unindent()
			builder.WriteLine("}")
		}
	}

	builder.Unindent()
	builder.WriteLine("}")
	builder.Unindent()
	builder.WriteLine("}")
}
