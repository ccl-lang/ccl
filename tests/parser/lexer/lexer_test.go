package lexer_test

import (
	"fmt"
	"testing"

	"github.com/ALiwoto/ccl/src/cclParser/cclLexer"
)

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

const input3 = `model MyModelName { field1: string; field2: int32; }`

func TestLexer3(t *testing.T) {
	lexResults, err := cclLexer.Lex(input3)
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
