package tsGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclValues"
)

func (c *TypeScriptGenerationContext) generateSerializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	builder.WriteLine("public serializeBinary(): Uint8Array {").
		Indent().
		WriteLine("const buffer: number[] = [];").
		WriteLine("const dataView = new DataView(new ArrayBuffer(8));").
		NewLine()

	for i := range model.Fields {
		field := model.Fields[i]
		if field.IsArray() {
			c.generateArraySerializeBinary(field, builder)
			continue
		}

		c.generateFieldSerializeBinary(field, builder)
	}

	builder.WriteLine("return new Uint8Array(buffer);").
		UnindentLine().
		WriteLine("}")
	return nil
}

func (c *TypeScriptGenerationContext) generateFieldSerializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	fieldName := "this." + caseUtils.ToCamelCase(field.Name)
	builder.MapVarPairs(
		"field", fieldName,
	)
	defer builder.UnmapVar("field")

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("{").
			Indent().
			LineD("const strBytes = new TextEncoder().encode($field);").
			WriteLine("dataView.setUint32(0, strBytes.length, true);").
			WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));").
			WriteLine("strBytes.forEach(b => buffer.push(b));").
			UnindentLine().
			WriteLine("}")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.LineD("dataView.setInt32(0, $field, true);").
			WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameInt8:
		builder.LineD("dataView.setInt8(0, $field);").
			WriteLine("buffer.push(dataView.getUint8(0));")
	case cclValues.TypeNameInt16:
		builder.LineD("dataView.setInt16(0, $field, true);").
			WriteLine("for (let i = 0; i < 2; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameInt64:
		// JS doesn't support 64-bit integers well without BigInt.
		// Assuming modern JS environment (Node, modern browsers) supports BigInt.
		builder.LineD("dataView.setBigInt64(0, BigInt($field), true);").
			WriteLine("for (let i = 0; i < 8; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.LineD("dataView.setUint32(0, $field, true);").
			WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameUint8:
		builder.LineD("dataView.setUint8(0, $field);").
			WriteLine("buffer.push(dataView.getUint8(0));")
	case cclValues.TypeNameUint16:
		builder.LineD("dataView.setUint16(0, $field, true);").
			WriteLine("for (let i = 0; i < 2; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameUint64:
		builder.LineD("dataView.setBigUint64(0, BigInt($field), true);").
			WriteLine("for (let i = 0; i < 8; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.LineD("dataView.setFloat32(0, $field, true);").
			WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameFloat64:
		builder.LineD("dataView.setFloat64(0, $field, true);").
			WriteLine("for (let i = 0; i < 8; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameBool:
		builder.LineD("buffer.push($field ? 1 : 0);")
	case cclValues.TypeNameBytes:
		builder.LineD("dataView.setUint32(0, $field.length, true);").
			WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));").
			LineD("$field.forEach(b => buffer.push(b));")
	case cclValues.TypeNameDateTime:
		builder.LineD("dataView.setBigInt64(0, BigInt($field), true);").
			WriteLine("for (let i = 0; i < 8; i++) buffer.push(dataView.getUint8(i));")
	default:
		if field.IsCustomTypeModel() {
			builder.LineD("if ($field) {").
				Indent().
				LineD("const customBytes = $field.serializeBinary();").
				WriteLine("dataView.setUint32(0, customBytes.length, true);").
				WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));").
				WriteLine("customBytes.forEach(b => buffer.push(b));").
				Unindent().
				WriteLine("} else {").
				Indent().
				WriteLine("dataView.setUint32(0, 1, true);").
				WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));").
				WriteLine("buffer.push(0);").
				UnindentLine().
				WriteLine("}")
		}
	}
	builder.NewLine()
}

func (c *TypeScriptGenerationContext) generateArraySerializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	fieldName := "this." + caseUtils.ToCamelCase(field.Name)
	targetFieldType := field.Type.GetUnderlyingType()
	builder.MapVarPairs(
		"field", fieldName,
	)
	defer builder.UnmapVar("field")

	builder.LineD("dataView.setUint32(0, $field.length, true);").
		WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));").
		LineD("for (const item of $field) {").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("const strBytes = new TextEncoder().encode(item);").
			WriteLine("dataView.setUint32(0, strBytes.length, true);").
			WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));").
			WriteLine("strBytes.forEach(b => buffer.push(b));")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.WriteLine("dataView.setInt32(0, item, true);").
			WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameInt8:
		builder.WriteLine("dataView.setInt8(0, item);").
			WriteLine("buffer.push(dataView.getUint8(0));")
	case cclValues.TypeNameInt16:
		builder.WriteLine("dataView.setInt16(0, item, true);").
			WriteLine("for (let i = 0; i < 2; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameInt64:
		builder.WriteLine("dataView.setBigInt64(0, BigInt(item), true);").
			WriteLine("for (let i = 0; i < 8; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.WriteLine("dataView.setUint32(0, item, true);").
			WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameUint8:
		builder.WriteLine("dataView.setUint8(0, item);").
			WriteLine("buffer.push(dataView.getUint8(0));")
	case cclValues.TypeNameUint16:
		builder.WriteLine("dataView.setUint16(0, item, true);").
			WriteLine("for (let i = 0; i < 2; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameUint64:
		builder.WriteLine("dataView.setBigUint64(0, BigInt(item), true);").
			WriteLine("for (let i = 0; i < 8; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.WriteLine("dataView.setFloat32(0, item, true);").
			WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameFloat64:
		builder.WriteLine("dataView.setFloat64(0, item, true);").
			WriteLine("for (let i = 0; i < 8; i++) buffer.push(dataView.getUint8(i));")
	case cclValues.TypeNameBool:
		builder.WriteLine("buffer.push(item ? 1 : 0);")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("if (item) {").
				Indent().
				WriteLine("const customBytes = item.serializeBinary();").
				WriteLine("dataView.setUint32(0, customBytes.length, true);").
				WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));").
				WriteLine("customBytes.forEach(b => buffer.push(b));").
				Unindent().
				WriteLine("} else {").
				Indent().
				WriteLine("dataView.setUint32(0, 1, true);").
				WriteLine("for (let i = 0; i < 4; i++) buffer.push(dataView.getUint8(i));").
				WriteLine("buffer.push(0);").
				UnindentLine().
				WriteLine("}")
		}
	}
	builder.UnindentLine().
		WriteLine("}")
}

func (c *TypeScriptGenerationContext) generateDeserializeBinaryMethod(model *CCLModel, builder *codeBuilder.CodeBuilder) error {
	builder.MapVarPairs(
		"model", model.Name,
	)
	defer builder.UnmapVar("model")

	builder.LineD("public static deserializeBinary(data: Uint8Array): $model | null {").
		Indent().
		WriteLine("if (!data || data.length === 0 || (data.length === 1 && data[0] === 0)) return null;").
		WriteLine("const view = new DataView(data.buffer, data.byteOffset, data.byteLength);").
		WriteLine("let offset = 0;").
		LineD("const result = new $model();").
		NewLine()

	for _, field := range model.Fields {
		if field.IsArray() {
			c.generateArrayDeserializeBinary(field, builder)
			continue
		}

		c.generateFieldDeserializeBinary(field, builder)
	}

	builder.WriteLine("return result;").
		UnindentLine().
		WriteLine("}")
	return nil
}

func (c *TypeScriptGenerationContext) generateFieldDeserializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	fieldName := caseUtils.ToCamelCase(field.Name)
	resultField := "result." + fieldName
	builder.MapVarPairs(
		"field", resultField,
		"type", field.Type.GetName(),
	)
	defer builder.UnmapVar(
		"field",
		"type",
	)

	switch field.Type.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("{").
			Indent().
			WriteLine("const len = view.getUint32(offset, true);").
			WriteLine("offset += 4;").
			WriteLine("if (len > data.length - offset) return null;").
			WriteLine("const bytes = data.subarray(offset, offset + len);").
			LineD("$field = new TextDecoder().decode(bytes);").
			WriteLine("offset += len;").
			UnindentLine().
			WriteLine("}")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.LineD("$field = view.getInt32(offset, true);").
			WriteLine("offset += 4;")
	case cclValues.TypeNameInt8:
		builder.LineD("$field = view.getInt8(offset);").
			WriteLine("offset += 1;")
	case cclValues.TypeNameInt16:
		builder.LineD("$field = view.getInt16(offset, true);").
			WriteLine("offset += 2;")
	case cclValues.TypeNameInt64:
		builder.LineD("$field = Number(view.getBigInt64(offset, true));").
			WriteLine("offset += 8;")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.LineD("$field = view.getUint32(offset, true);").
			WriteLine("offset += 4;")
	case cclValues.TypeNameUint8:
		builder.LineD("$field = view.getUint8(offset);").
			WriteLine("offset += 1;")
	case cclValues.TypeNameUint16:
		builder.LineD("$field = view.getUint16(offset, true);").
			WriteLine("offset += 2;")
	case cclValues.TypeNameUint64:
		builder.LineD("$field = Number(view.getBigUint64(offset, true));").
			WriteLine("offset += 8;")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.LineD("$field = view.getFloat32(offset, true);").
			WriteLine("offset += 4;")
	case cclValues.TypeNameFloat64:
		builder.LineD("$field = view.getFloat64(offset, true);").
			WriteLine("offset += 8;")
	case cclValues.TypeNameBool:
		builder.LineD("$field = view.getInt8(offset) !== 0;").
			WriteLine("offset += 1;")
	case cclValues.TypeNameBytes:
		builder.WriteLine("{").
			Indent().
			WriteLine("const len = view.getUint32(offset, true);").
			WriteLine("offset += 4;").
			LineD("$field = data.slice(offset, offset + len);").
			WriteLine("offset += len;").
			UnindentLine().
			WriteLine("}")
	case cclValues.TypeNameDateTime:
		builder.LineD("$field = Number(view.getBigInt64(offset, true));").
			WriteLine("offset += 8;")
	default:
		if field.IsCustomTypeModel() {
			builder.WriteLine("{").
				Indent().
				WriteLine("const len = view.getUint32(offset, true);").
				WriteLine("offset += 4;").
				WriteLine("if (len > data.length - offset) return null;").
				WriteLine("const bytes = data.subarray(offset, offset + len);").
				LineD("$field = $type.deserializeBinary(bytes);").
				WriteLine("offset += len;").
				UnindentLine().
				WriteLine("}")
		}
	}
	builder.NewLine()
}

func (c *TypeScriptGenerationContext) generateArrayDeserializeBinary(field *CCLField, builder *codeBuilder.CodeBuilder) {
	fieldName := caseUtils.ToCamelCase(field.Name)
	resultField := "result." + fieldName
	targetFieldType := field.Type.GetUnderlyingType()
	builder.MapVarPairs(
		"field", resultField,
		"type", targetFieldType.GetName(),
	)
	defer builder.UnmapVar(
		"field",
		"type",
	)

	builder.WriteLine("{").
		Indent().
		WriteLine("const len = view.getUint32(offset, true);").
		WriteLine("offset += 4;").
		LineD("$field = [];").
		WriteLine("for (let i = 0; i < len; i++) {").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameString:
		builder.WriteLine("const itemLen = view.getUint32(offset, true);").
			WriteLine("offset += 4;").
			WriteLine("if (itemLen > data.length - offset) return null;").
			WriteLine("const bytes = data.subarray(offset, offset + itemLen);").
			LineD("$field.push(new TextDecoder().decode(bytes));").
			WriteLine("offset += itemLen;")
	case cclValues.TypeNameInt, cclValues.TypeNameInt32:
		builder.LineD("$field.push(view.getInt32(offset, true));").
			WriteLine("offset += 4;")
	case cclValues.TypeNameInt8:
		builder.LineD("$field.push(view.getInt8(offset));").
			WriteLine("offset += 1;")
	case cclValues.TypeNameInt16:
		builder.LineD("$field.push(view.getInt16(offset, true));").
			WriteLine("offset += 2;")
	case cclValues.TypeNameInt64:
		builder.LineD("$field.push(Number(view.getBigInt64(offset, true)));").
			WriteLine("offset += 8;")
	case cclValues.TypeNameUint, cclValues.TypeNameUint32:
		builder.LineD("$field.push(view.getUint32(offset, true));").
			WriteLine("offset += 4;")
	case cclValues.TypeNameUint8:
		builder.LineD("$field.push(view.getUint8(offset));").
			WriteLine("offset += 1;")
	case cclValues.TypeNameUint16:
		builder.LineD("$field.push(view.getUint16(offset, true));").
			WriteLine("offset += 2;")
	case cclValues.TypeNameUint64:
		builder.LineD("$field.push(Number(view.getBigUint64(offset, true)));").
			WriteLine("offset += 8;")
	case cclValues.TypeNameFloat, cclValues.TypeNameFloat32:
		builder.LineD("$field.push(view.getFloat32(offset, true));").
			WriteLine("offset += 4;")
	case cclValues.TypeNameFloat64:
		builder.LineD("$field.push(view.getFloat64(offset, true));").
			WriteLine("offset += 8;")
	case cclValues.TypeNameBool:
		builder.LineD("$field.push(view.getInt8(offset) !== 0);").
			WriteLine("offset += 1;")
	default:
		if targetFieldType.IsCustomTypeModel() {
			builder.WriteLine("const itemLen = view.getUint32(offset, true);").
				WriteLine("offset += 4;").
				WriteLine("if (itemLen > data.length - offset) return null;").
				WriteLine("const bytes = data.subarray(offset, offset + itemLen);").
				LineD("$field.push($type.deserializeBinary(bytes));").
				WriteLine("offset += itemLen;")
		}
	}
	builder.UnindentLine().
		WriteLine("}").
		UnindentLine().
		WriteLine("}")
}
