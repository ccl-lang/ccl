package rsGenerator

import (
	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *RustGenerationContext) generateModelBinaryMethods(
	builder *codeBuilder.CodeBuilder,
	model *CCLModel,
) error {
	endian, err := c.GetBinarySerializationEndian(CurrentLanguage, model)
	if err != nil {
		return err
	}
	strictBinaryParsing, err := c.UsesStrictBinaryParsing(CurrentLanguage, model)
	if err != nil {
		return err
	}

	byteOrder := "le"
	if endian == gValues.EndianBig {
		byteOrder = "be"
	}

	builder.WriteLine("impl " + model.Name + " {").
		Indent()
	c.generateSerializeBinaryMethod(builder, model, byteOrder)
	c.generateDeserializeBinaryMethod(builder, model, byteOrder, strictBinaryParsing)
	builder.Unindent().
		WriteLine("}").
		NewLine()
	return nil
}

func (c *RustGenerationContext) generateSerializeBinaryMethod(
	builder *codeBuilder.CodeBuilder,
	model *CCLModel,
	byteOrder string,
) {
	builder.WriteLine("pub fn serialize_binary(&self) -> Vec<u8> {").
		Indent().
		WriteLine("let mut buffer = Vec::new();")
	for _, field := range model.Fields {
		if field.IsArray() {
			c.generateArraySerializeBinary(builder, field, byteOrder)
		} else {
			c.generateFieldSerializeBinary(builder, field, byteOrder)
		}
	}
	builder.WriteLine("buffer").
		Unindent().
		WriteLine("}").
		NewLine()
}

func (c *RustGenerationContext) generateDeserializeBinaryMethod(
	builder *codeBuilder.CodeBuilder,
	model *CCLModel,
	byteOrder string,
	strictBinaryParsing bool,
) {
	fallback := `return Ok(result);`
	if strictBinaryParsing {
		fallback = `return Err("not enough binary data".to_string());`
	}

	builder.WriteLine("pub fn deserialize_binary(data: &[u8]) -> Result<Self, String> {").
		Indent().
		WriteLine("if data.is_empty() || (data.len() == 1 && data[0] == 0) {").
		Indent().
		WriteLine("return Ok(Self::default());").
		Unindent().
		WriteLine("}").
		WriteLine("let mut result = Self::default();").
		WriteLine("let mut offset: usize = 0;")
	for _, field := range model.Fields {
		if field.IsArray() {
			c.generateArrayDeserializeBinary(builder, field, byteOrder, fallback)
		} else {
			c.generateFieldDeserializeBinary(builder, field, byteOrder, fallback)
		}
	}
	builder.WriteLine("Ok(result)").
		Unindent().
		WriteLine("}").
		NewLine()
}

func (c *RustGenerationContext) generateFieldSerializeBinary(
	builder *codeBuilder.CodeBuilder,
	field *CCLField,
	byteOrder string,
) {
	fieldName := "self." + rustFieldName(field.Name)
	targetType := field.Type
	switch rustStorageTypeName(targetType) {
	case cclValues.TypeNameString:
		builder.WriteLine("{").
			Indent().
			WriteLine("let field_bytes = " + fieldName + ".as_bytes();").
			WriteLine("buffer.extend_from_slice(&(field_bytes.len() as u32).to_" + byteOrder + "_bytes());").
			WriteLine("buffer.extend_from_slice(field_bytes);").
			Unindent().
			WriteLine("}")
	case cclValues.TypeNameBytes:
		builder.WriteLine("buffer.extend_from_slice(&(" + fieldName + ".len() as u32).to_" + byteOrder + "_bytes());").
			WriteLine("buffer.extend_from_slice(&" + fieldName + ");")
	case cclValues.TypeNameBool:
		builder.WriteLine("buffer.push(if " + fieldName + " { 1 } else { 0 });")
	default:
		if targetType.IsCustomTypeModel() {
			builder.WriteLine("if let Some(value) = &" + fieldName + " {").
				Indent().
				WriteLine("let custom_bytes = value.serialize_binary();").
				WriteLine("buffer.extend_from_slice(&(custom_bytes.len() as u32).to_" + byteOrder + "_bytes());").
				WriteLine("buffer.extend_from_slice(&custom_bytes);").
				Unindent().
				WriteLine("} else {").
				Indent().
				WriteLine("buffer.extend_from_slice(&1u32.to_" + byteOrder + "_bytes());").
				WriteLine("buffer.push(0);").
				Unindent().
				WriteLine("}")
			return
		}

		c.generateRustScalarWrite(builder, targetType, fieldName, byteOrder)
	}
}

func (c *RustGenerationContext) generateArraySerializeBinary(
	builder *codeBuilder.CodeBuilder,
	field *CCLField,
	byteOrder string,
) {
	fieldName := "self." + rustFieldName(field.Name)
	targetType := field.Type.GetUnderlyingType()
	builder.WriteLine("buffer.extend_from_slice(&(" + fieldName + ".len() as u32).to_" + byteOrder + "_bytes());").
		WriteLine("for item in &" + fieldName + " {").
		Indent()

	switch rustStorageTypeName(targetType) {
	case cclValues.TypeNameString:
		builder.WriteLine("let item_bytes = item.as_bytes();").
			WriteLine("buffer.extend_from_slice(&(item_bytes.len() as u32).to_" + byteOrder + "_bytes());").
			WriteLine("buffer.extend_from_slice(item_bytes);")
	case cclValues.TypeNameBytes:
		builder.WriteLine("buffer.extend_from_slice(&(item.len() as u32).to_" + byteOrder + "_bytes());").
			WriteLine("buffer.extend_from_slice(item);")
	case cclValues.TypeNameBool:
		builder.WriteLine("buffer.push(if *item { 1 } else { 0 });")
	default:
		if targetType.IsCustomTypeModel() {
			builder.WriteLine("if let Some(value) = item {").
				Indent().
				WriteLine("let custom_bytes = value.serialize_binary();").
				WriteLine("buffer.extend_from_slice(&(custom_bytes.len() as u32).to_" + byteOrder + "_bytes());").
				WriteLine("buffer.extend_from_slice(&custom_bytes);").
				Unindent().
				WriteLine("} else {").
				Indent().
				WriteLine("buffer.extend_from_slice(&1u32.to_" + byteOrder + "_bytes());").
				WriteLine("buffer.push(0);").
				Unindent().
				WriteLine("}")
		} else {
			c.generateRustScalarWrite(builder, targetType, "*item", byteOrder)
		}
	}

	builder.Unindent().
		WriteLine("}")
}

func (c *RustGenerationContext) generateRustScalarWrite(
	builder *codeBuilder.CodeBuilder,
	targetType *cclValues.CCLTypeUsage,
	expression string,
	byteOrder string,
) {
	if targetType.IsCustomTypeEnum() {
		enumDef := targetType.GetDefinition().GetEnumDefinition()
		builder.WriteLine("buffer.extend_from_slice(&(" + expression + " as " +
			c.getRustEnumBaseType(enumDef) + ").to_" + byteOrder + "_bytes());")
		return
	}

	rustType := CCLTypesToRustTypes[targetType.GetName()]
	builder.WriteLine("buffer.extend_from_slice(&(" + expression + " as " +
		rustType + ").to_" + byteOrder + "_bytes());")
}

func (c *RustGenerationContext) generateFieldDeserializeBinary(
	builder *codeBuilder.CodeBuilder,
	field *CCLField,
	byteOrder string,
	fallback string,
) {
	fieldName := "result." + rustFieldName(field.Name)
	targetType := field.Type
	switch rustStorageTypeName(targetType) {
	case cclValues.TypeNameString:
		c.generateRustLengthRead(builder, "field_len", byteOrder, fallback)
		c.generateRustBoundsCheck(builder, "field_len", fallback)
		builder.WriteLine(fieldName + " = String::from_utf8(data[offset..offset + field_len].to_vec())").
			Indent().
			WriteLine(`.map_err(|err| err.to_string())?;`).
			Unindent().
			WriteLine("offset += field_len;")
	case cclValues.TypeNameBytes:
		c.generateRustLengthRead(builder, "field_len", byteOrder, fallback)
		c.generateRustBoundsCheck(builder, "field_len", fallback)
		builder.WriteLine(fieldName + " = data[offset..offset + field_len].to_vec();").
			WriteLine("offset += field_len;")
	case cclValues.TypeNameBool:
		c.generateRustBoundsCheck(builder, "1", fallback)
		builder.WriteLine(fieldName + " = data[offset] != 0;").
			WriteLine("offset += 1;")
	default:
		if targetType.IsCustomTypeModel() {
			c.generateRustCustomModelRead(builder, fieldName, targetType.GetName(), byteOrder, fallback)
		} else {
			valueExpr := c.generateRustScalarRead(builder, targetType, "field_value", byteOrder, fallback)
			builder.WriteLine(fieldName + " = " + valueExpr + ";")
		}
	}
}

func (c *RustGenerationContext) generateArrayDeserializeBinary(
	builder *codeBuilder.CodeBuilder,
	field *CCLField,
	byteOrder string,
	fallback string,
) {
	fieldName := "result." + rustFieldName(field.Name)
	targetType := field.Type.GetUnderlyingType()
	c.generateRustLengthRead(builder, "array_len", byteOrder, fallback)
	builder.WriteLine(fieldName + " = Vec::with_capacity(array_len);").
		WriteLine("for _ in 0..array_len {").
		Indent()

	switch rustStorageTypeName(targetType) {
	case cclValues.TypeNameString:
		c.generateRustLengthRead(builder, "item_len", byteOrder, fallback)
		c.generateRustBoundsCheck(builder, "item_len", fallback)
		builder.WriteLine("let item = String::from_utf8(data[offset..offset + item_len].to_vec())").
			Indent().
			WriteLine(`.map_err(|err| err.to_string())?;`).
			Unindent().
			WriteLine("offset += item_len;").
			WriteLine(fieldName + ".push(item);")
	case cclValues.TypeNameBytes:
		c.generateRustLengthRead(builder, "item_len", byteOrder, fallback)
		c.generateRustBoundsCheck(builder, "item_len", fallback)
		builder.WriteLine("let item = data[offset..offset + item_len].to_vec();").
			WriteLine("offset += item_len;").
			WriteLine(fieldName + ".push(item);")
	case cclValues.TypeNameBool:
		c.generateRustBoundsCheck(builder, "1", fallback)
		builder.WriteLine(fieldName + ".push(data[offset] != 0);").
			WriteLine("offset += 1;")
	default:
		if targetType.IsCustomTypeModel() {
			c.generateRustCustomModelArrayRead(builder, fieldName, targetType.GetName(), byteOrder, fallback)
		} else {
			valueExpr := c.generateRustScalarRead(builder, targetType, "item_value", byteOrder, fallback)
			builder.WriteLine(fieldName + ".push(" + valueExpr + ");")
		}
	}

	builder.Unindent().
		WriteLine("}")
}

func (c *RustGenerationContext) generateRustLengthRead(
	builder *codeBuilder.CodeBuilder,
	name string,
	byteOrder string,
	fallback string,
) {
	c.generateRustBoundsCheck(builder, "4", fallback)
	builder.WriteLine("let " + name + " = u32::from_" + byteOrder + "_bytes(data[offset..offset + 4].try_into().map_err(|_| \"invalid binary data\".to_string())?) as usize;").
		WriteLine("offset += 4;")
}

func (c *RustGenerationContext) generateRustBoundsCheck(
	builder *codeBuilder.CodeBuilder,
	requiredBytes string,
	fallback string,
) {
	builder.WriteLine("if data.len().saturating_sub(offset) < " + requiredBytes + " {").
		Indent().
		WriteLine(fallback).
		Unindent().
		WriteLine("}")
}

func (c *RustGenerationContext) generateRustScalarRead(
	builder *codeBuilder.CodeBuilder,
	targetType *cclValues.CCLTypeUsage,
	name string,
	byteOrder string,
	fallback string,
) string {
	rustType := CCLTypesToRustTypes[targetType.GetName()]
	if targetType.IsCustomTypeEnum() {
		rustType = c.getRustEnumBaseType(targetType.GetDefinition().GetEnumDefinition())
	}

	size := rustScalarByteSize(rustType)
	c.generateRustBoundsCheck(builder, size, fallback)
	builder.WriteLine("let " + name + " = " + rustType + "::from_" + byteOrder + "_bytes(data[offset..offset + " + size + "].try_into().map_err(|_| \"invalid binary data\".to_string())?);").
		WriteLine("offset += " + size + ";")

	if targetType.IsCustomTypeEnum() {
		enumType, _ := c.getRustEnumTypeName(targetType.GetDefinition().GetEnumDefinition())
		return enumType + "::from_raw(" + name + ").ok_or_else(|| \"invalid enum value\".to_string())?"
	}
	return name
}

func (c *RustGenerationContext) generateRustCustomModelRead(
	builder *codeBuilder.CodeBuilder,
	fieldName string,
	typeName string,
	byteOrder string,
	fallback string,
) {
	c.generateRustLengthRead(builder, "field_len", byteOrder, fallback)
	c.generateRustBoundsCheck(builder, "field_len", fallback)
	builder.WriteLine("let field_bytes = &data[offset..offset + field_len];").
		WriteLine("offset += field_len;").
		WriteLine("if field_bytes.is_empty() || (field_bytes.len() == 1 && field_bytes[0] == 0) {").
		Indent().
		WriteLine(fieldName + " = None;").
		Unindent().
		WriteLine("} else {").
		Indent().
		WriteLine(fieldName + " = Some(Box::new(" + typeName + "::deserialize_binary(field_bytes)?));").
		Unindent().
		WriteLine("}")
}

func (c *RustGenerationContext) generateRustCustomModelArrayRead(
	builder *codeBuilder.CodeBuilder,
	fieldName string,
	typeName string,
	byteOrder string,
	fallback string,
) {
	c.generateRustLengthRead(builder, "item_len", byteOrder, fallback)
	c.generateRustBoundsCheck(builder, "item_len", fallback)
	builder.WriteLine("let item_bytes = &data[offset..offset + item_len];").
		WriteLine("offset += item_len;").
		WriteLine("if item_bytes.is_empty() || (item_bytes.len() == 1 && item_bytes[0] == 0) {").
		Indent().
		WriteLine(fieldName + ".push(None);").
		Unindent().
		WriteLine("} else {").
		Indent().
		WriteLine(fieldName + ".push(Some(Box::new(" + typeName + "::deserialize_binary(item_bytes)?)));").
		Unindent().
		WriteLine("}")
}
