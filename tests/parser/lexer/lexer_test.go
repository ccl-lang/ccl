package lexer_test

import (
	"fmt"
	"testing"

	"github.com/ccl-lang/ccl/src/cclParser/cclLexer"
)

//---------------------------------------------------------

// expectedToken is a struct that holds the expected token type and value.
type expectedToken struct {
	token cclLexer.CCLTokenType
	value string
}

//---------------------------------------------------------

const input1 = `
// This is a global attribute, applied to the whole generation process
#[CCLVersion("1.0.0")]
#[SerializationType("binary")]

// This is a comment

// models can also have attribute on them.
[SerializationType("binary")]
model UserInfo {
    Id: int64;
    Username: string;
    Email: string;
    ProfileImage: bytes;
    CreatedAt: datetime;
    UpdatedAt: datetime;
}

// models support multiple attributes
[SerializationType("binary")]
[SerializationType("json")]
model GetUsersResult {
    Users: UserInfo[];
    OtherUsers: UserInfo[];
}
`

func TestLexer1(t *testing.T) {
	lexResults, err := cclLexer.Lex(input1)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(lexResults) == 0 {
		t.Fatalf("No tokens found")
	}

	fmt.Println("Tokens:")
	for _, token := range lexResults {
		fmt.Println(token)
	}
}

const input2 = `#[System.Text.SerializationType("C#", "binary")]`

func TestLexer2(t *testing.T) {
	lexResults, err := cclLexer.Lex(input2)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(lexResults) == 0 {
		t.Fatalf("No tokens found")
	}

	fmt.Println("Tokens:")
	for _, token := range lexResults {
		fmt.Println(token)
	}
}

const input3 = `model MyModelName { field1: string = "default value"; field2: int32; }`

var input3Tokens = []expectedToken{
	{cclLexer.TokenTypeKeywordModel, "model"},          // model
	{cclLexer.TokenTypeIdentifier, "MyModelName"},      // MyModelName
	{cclLexer.TokenTypeLeftBrace, "{"},                 // {
	{cclLexer.TokenTypeIdentifier, "field1"},           // field1
	{cclLexer.TokenTypeColon, ":"},                     // :
	{cclLexer.TokenTypeDataType, "string"},             // string
	{cclLexer.TokenTypeAssignment, "="},                // =
	{cclLexer.TokenTypeStringLiteral, "default value"}, // "default value"
	{cclLexer.TokenTypeSemicolon, ";"},                 // ;
	{cclLexer.TokenTypeIdentifier, "field2"},           // field2
	{cclLexer.TokenTypeColon, ":"},                     // :
	{cclLexer.TokenTypeDataType, "int32"},              // int32
	{cclLexer.TokenTypeSemicolon, ";"},                 // ;
	{cclLexer.TokenTypeRightBrace, "}"},                // }
}

func TestLexer3(t *testing.T) {
	lexResults, err := cclLexer.Lex(input3)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(lexResults) == 0 {
		t.Fatalf("No tokens found")
	}

	fmt.Println("Tokens:")
	for i, token := range lexResults {
		if i >= len(input3Tokens) {
			t.Fatalf("Unexpected token: %v", token)
			return
		}
		if token.Type != input3Tokens[i].token {
			t.Fatalf("Expected token type %v, got %v", input3Tokens[i], token.Type)
			return
		} else if formatted := token.FormatValueAsString(); formatted != input3Tokens[i].value {
			t.Fatalf("Expected token value %v, got %v", input3Tokens[i].value, formatted)
			return
		}

		fmt.Println(token)
	}
}

const input4 = `
model MyModelName {
	field1: string = "default value";
	field2: int32;
}

something *= 1234;
anotherValue /= 1234;
anotherValue = something ** 1234;
`

var input4Tokens = []expectedToken{
	{cclLexer.TokenTypeKeywordModel, "model"},          // model
	{cclLexer.TokenTypeIdentifier, "MyModelName"},      // MyModelName
	{cclLexer.TokenTypeLeftBrace, "{"},                 // {
	{cclLexer.TokenTypeIdentifier, "field1"},           // field1
	{cclLexer.TokenTypeColon, ":"},                     // :
	{cclLexer.TokenTypeDataType, "string"},             // string
	{cclLexer.TokenTypeAssignment, "="},                // =
	{cclLexer.TokenTypeStringLiteral, "default value"}, // "default value"
	{cclLexer.TokenTypeSemicolon, ";"},                 // ;
	{cclLexer.TokenTypeIdentifier, "field2"},           // field2
	{cclLexer.TokenTypeColon, ":"},                     // :
	{cclLexer.TokenTypeDataType, "int32"},              // int32
	{cclLexer.TokenTypeSemicolon, ";"},                 // ;
	{cclLexer.TokenTypeRightBrace, "}"},                // }
	{cclLexer.TokenTypeIdentifier, "something"},        // something
	{cclLexer.TokenTypeMultiplyAssignment, "*="},       // *=
	{cclLexer.TokenTypeIntLiteral, "1234"},             // 1234
	{cclLexer.TokenTypeSemicolon, ";"},                 // ;
	{cclLexer.TokenTypeIdentifier, "anotherValue"},     // anotherValue
	{cclLexer.TokenTypeDivideAssignment, "/="},         // /=
	{cclLexer.TokenTypeIntLiteral, "1234"},             // 1234
	{cclLexer.TokenTypeSemicolon, ";"},                 // ;
	{cclLexer.TokenTypeIdentifier, "anotherValue"},     // anotherValue
	{cclLexer.TokenTypeAssignment, "="},                // =
	{cclLexer.TokenTypeIdentifier, "something"},        // something
	{cclLexer.TokenTypePower, "**"},                    // **
	{cclLexer.TokenTypeIntLiteral, "1234"},             // 1234
	{cclLexer.TokenTypeSemicolon, ";"},                 // ;
}

func TestLexer4(t *testing.T) {
	lexResults, err := cclLexer.Lex(input4)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(lexResults) == 0 {
		t.Fatalf("No tokens found")
	}

	fmt.Println("Tokens:")
	for i, token := range lexResults {
		if i >= len(input4Tokens) {
			t.Fatalf("Unexpected token: %v", token)
			return
		}
		if token.Type != input4Tokens[i].token {
			t.Fatalf("Expected token type %v, got %v", input4Tokens[i], token.Type)
			return
		} else if formatted := token.FormatValueAsString(); formatted != input4Tokens[i].value {
			t.Fatalf("Expected token value %v, got %v", input4Tokens[i].value, formatted)
			return
		}

		fmt.Println(token)
	}
}

//---------------------------------------------------------

const inputAttributeWithInteger = `
#[System.Text.SerializationType(
	"C#",
	"binary",
	123,
	4.56, 
	0x123,
	0b1010,
	-19,
)]
`

func TestAttributeWithInteger(t *testing.T) {
	lexResults, err := cclLexer.Lex(inputAttributeWithInteger)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(lexResults) == 0 {
		t.Fatalf("No tokens found")
	}

	stringLiteral := lexResults[10].GetStringLiteral()
	if stringLiteral != "binary" {
		t.Fatalf("Expected string literal 'binary', got '%s'", stringLiteral)
	}

	fmt.Println("Tokens:")
	for _, token := range lexResults {
		fmt.Println(token)
	}
}

//---------------------------------------------------------

const wrongInput1 = `modelZ MyModelName { field1: WrongTypeHere; field2: int32; }`

func TestLexerWrongInput1(t *testing.T) {
	lexResults, err := cclLexer.Lex(wrongInput1)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if len(lexResults) == 0 {
		t.Fatalf("No tokens found")
		return
	} else if len(lexResults) < 12 {
		t.Fatalf("Expected at least 12 tokens, got %d", len(lexResults))
		return
	}

	// The lexer's job is just to recognize tokens based on their patterns,
	// not their semantic meaning or context.
	// So having two IDENTIFIERs in a row (modelZ and MyModelName) is perfectly
	// valid at the lexical level.
	if lexResults[0].Type != cclLexer.TokenTypeIdentifier ||
		lexResults[1].Type != cclLexer.TokenTypeIdentifier {
		t.Fatalf("Expected 2 consecutive IDENTIFIER tokens")
		return
	}

	// For the "WrongTypeHere" issue, this is typically caught by a semantic analyzer
	// (sometimes called a type checker) that runs after parsing. This phase would:
	// 1. Successfully parse the syntax (if we fixed the modelZ issue)
	// 2. Build an AST
	// 3. Check that all referenced types exist
	// 4. Report: "Unknown type 'WrongTypeHere' for field 'field1'"
	// But at the lexical level, "WrongTypeHere" is just an IDENTIFIER token.

	// To sum things up:
	// 	The Three Levels of Errors:
	// 1. Lexical errors: Invalid tokens (e.g., @#$% if those aren't valid in our language)
	// 2. Syntax errors: Valid tokens but invalid syntax (e.g., modelZ instead of model)
	// 3. Semantic errors: Syntactically valid but meaningless constructs (e.g., undefined types)
}

//---------------------------------------------------------
