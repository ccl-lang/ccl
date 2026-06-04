package gdGenerator

import (
	"github.com/ALiwoto/ssg/ssg/caseUtils"
	"github.com/ccl-lang/ccl/src/core/cclUtils/codeBuilder"
	"github.com/ccl-lang/ccl/src/core/cclValues"
)

func (c *GDScriptGenerationContext) generateSerializeJsonMethods(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	if err := c.generateSerializeJsonDictMethod(model, builder); err != nil {
		return err
	}
	if err := c.generateSerializeJsonMethod(builder); err != nil {
		return err
	}
	if err := c.generateDeserializeJsonDictMethod(model, builder); err != nil {
		return err
	}
	if err := c.generateDeserializeJsonMethod(builder); err != nil {
		return err
	}

	return nil
}

func (c *GDScriptGenerationContext) generateSerializeJsonDictMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.WriteLine("func serialize_json_dict() -> Dictionary:").
		Indent().
		WriteLine("var data: Dictionary = {}").
		NewLine()

	for _, field := range model.Fields {
		jsonName, err := c.GetJsonFieldName(CurrentLanguage, model, field)
		if err != nil {
			return err
		}

		if field.IsArray() {
			err = c.generateArraySerializeJson(field, jsonName, builder)
			if err != nil {
				return err
			}
			continue
		}

		c.generateFieldSerializeJson(field, jsonName, builder)
	}

	builder.WriteLine("return data").
		UnindentLine()

	return nil
}

func (c *GDScriptGenerationContext) generateSerializeJsonMethod(
	builder *codeBuilder.CodeBuilder,
) error {
	builder.WriteLine("func serialize_json() -> String:").
		Indent().
		WriteLine("return JSON.stringify(self.serialize_json_dict())").
		UnindentLine()

	return nil
}

func (c *GDScriptGenerationContext) generateDeserializeJsonDictMethod(
	model *CCLModel,
	builder *codeBuilder.CodeBuilder,
) error {
	builder.ExpectMappedVars(
		"model",
	).MapVarPairs(
		"modelResult", modelResultName,
	)
	defer builder.UnmapVar(
		"modelResult",
	)

	builder.LineD("static func deserialize_json_dict(data: Dictionary) -> $model:").
		Indent().
		WriteLine("if not data:").
		Indent().
		WriteLine("return null").
		UnindentLine().
		LineD("var $modelResult = $model.new()").
		NewLine()

	for _, field := range model.Fields {
		jsonName, err := c.GetJsonFieldName(CurrentLanguage, model, field)
		if err != nil {
			return err
		}

		if field.IsArray() {
			err = c.generateArrayDeserializeJson(field, jsonName, builder)
			if err != nil {
				return err
			}
			continue
		}

		err = c.generateFieldDeserializeJson(field, jsonName, builder)
		if err != nil {
			return err
		}
	}

	builder.LineD("return $modelResult").
		UnindentLine()

	return nil
}

// generateDeserializeJsonMethod generates code for deserializing json method.
// The method expects `model` mapped-var to be present in the builder.
func (c *GDScriptGenerationContext) generateDeserializeJsonMethod(
	builder *codeBuilder.CodeBuilder,
) error {
	builder.ExpectMappedVars(
		"model",
	)

	builder.LineD("static func deserialize_json(json_text: String) -> $model:").
		Indent().
		WriteLine(`if not json_text:`).
		Indent().
		WriteLine("return null").
		UnindentLine().
		WriteLine("var json = JSON.new()").
		WriteLine("var err = json.parse(json_text)").
		WriteLine("if err != OK:").
		Indent().
		LineD(`push_error("Failed to parse JSON for $model: " + json.get_error_message())`).
		WriteLine("return null").
		UnindentLine().
		WriteLine("if not (json.data is Dictionary):").
		Indent().
		LineD(`push_error("Expected JSON object for $model, got ", json.data)`).
		WriteLine("return null").
		UnindentLine().
		LineD(`return $model.deserialize_json_dict(json.data)`).
		UnindentLine()

	return nil
}

func (c *GDScriptGenerationContext) generateFieldSerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) {
	fieldRawName := caseUtils.ToSnakeCase(field.Name)
	resultField := "self." + fieldRawName

	builder.MapVarPairs(
		"field", resultField,
		"jsonName", jsonName,
	)
	defer builder.UnmapVar(
		"field",
		"jsonName",
	)

	switch field.Type.GetName() {
	case cclValues.TypeNameBytes:
		builder.LineD(`data["$jsonName"] = Marshalls.raw_to_base64($field)`)
	default:
		if field.IsCustomTypeModel() {
			builder.WriteLine(`@warning_ignore("incompatible_ternary")`).
				LineD(`data["$jsonName"] = $field.serialize_json_dict() if $field else null`)
		} else {
			builder.LineD(`data["$jsonName"] = $field`)
		}
	}

	builder.NewLine()
}

func (c *GDScriptGenerationContext) generateArraySerializeJson(
	field *CCLField,
	jsonName string,
	builder *codeBuilder.CodeBuilder,
) error {
	targetFieldType := field.Type.GetUnderlyingType()
	fieldRawName := caseUtils.ToSnakeCase(field.Name)
	resultField := "self." + fieldRawName
	listName := fieldRawName + "_list"

	builder.MapVarPairs(
		"list", listName,
		"field", resultField,
		"jsonName", jsonName,
	)
	defer builder.UnmapVar(
		"list",
		"field",
		"jsonName",
	)

	builder.LineD("var $list = []").
		LineD("if $field != null:").
		Indent().
		LineD("for item in $field:").
		Indent()

	switch targetFieldType.GetName() {
	case cclValues.TypeNameBytes:
		builder.LineD("$list.append(Marshalls.raw_to_base64(item))")
	default:
		if targetFieldType.IsCustomTypeModel() {
			// this will definitely break for < godot 4.0 (e.g. 3.x)
			// but then we are using lots of other features which are not available
			// for 3.x and lower anyway...
			builder.WriteLine(`@warning_ignore("incompatible_ternary")`).
				LineD("$list.append(item.serialize_json_dict() if item else null)")
		} else {
			// are we sure about this plain else? are all "non-custom types" out there
			// supported by json schema (or at least Godot's json serializer)?
			// well, if that's not the case, this is the place where it should be added
			// in future.
			builder.LineD("$list.append(item)")
		}
	}

	builder.Unindent().
		LineD(`data["$jsonName"] = $list`).
		Unindent().
		WriteLine("else:").
		Indent().
		LineD(`data["$jsonName"] = null`).
		UnindentLine()

	return nil
}
