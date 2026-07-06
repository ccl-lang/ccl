package csGenerator

import "github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclUtils/codeBuilder"

func (c *CSharpGenerationContext) generateCSharpCustomModelFieldDeserializeBinary(
	builder *codeBuilder.CodeBuilder,
) {
	builder.WriteLine("{").
		Indent()
	c.generateCSharpBinaryDeserializeBoundsCheck(builder, "1")
	builder.WriteLine("var present = reader.ReadByte();").
		WriteLine("if (present == 0)").
		WriteLine("{").
		Indent().
		LineD("$field = null;").
		Unindent().
		WriteLine("}").
		WriteLine("else if (present == 1)").
		WriteLine("{").
		Indent()
	c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
	builder.WriteLine("var len = reader.ReadUInt32();").
		LineD("if (len > ms.Length - ms.Position) return $binaryParseFallback;").
		WriteLine("var bytes = reader.ReadBytes((int)len);").
		LineD("$field = $type.DeserializeBinary(bytes);").
		Unindent().
		WriteLine("}").
		WriteLine("else").
		WriteLine("{").
		Indent().
		LineD("return $binaryParseFallback;").
		Unindent().
		WriteLine("}").
		Unindent().
		WriteLine("}")
}

func (c *CSharpGenerationContext) generateCSharpCustomModelArrayItemDeserializeBinary(
	builder *codeBuilder.CodeBuilder,
) {
	c.generateCSharpBinaryDeserializeBoundsCheck(builder, "1")
	builder.WriteLine("var itemPresent = reader.ReadByte();").
		WriteLine("if (itemPresent == 0)").
		WriteLine("{").
		Indent().
		LineD("$field.Add(null);").
		Unindent().
		WriteLine("}").
		WriteLine("else if (itemPresent == 1)").
		WriteLine("{").
		Indent()
	c.generateCSharpBinaryDeserializeBoundsCheck(builder, "4")
	builder.WriteLine("var itemLen = reader.ReadUInt32();").
		LineD("if (itemLen > ms.Length - ms.Position) return $binaryParseFallback;").
		WriteLine("var bytes = reader.ReadBytes((int)itemLen);").
		LineD("$field.Add($type.DeserializeBinary(bytes));").
		Unindent().
		WriteLine("}").
		WriteLine("else").
		WriteLine("{").
		Indent().
		LineD("return $binaryParseFallback;").
		Unindent().
		WriteLine("}")
}

func (c *CSharpGenerationContext) generateCSharpBinaryDeserializeBoundsCheck(
	builder *codeBuilder.CodeBuilder,
	requiredBytes string,
) {
	builder.MapVarPairs("requiredBytes", requiredBytes)
	builder.LineD("if (ms.Length - ms.Position < $requiredBytes) return $binaryParseFallback;")
	builder.UnmapVar("requiredBytes")
}
