package goGen_test

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/inputLangs/cclInput/cclParser"
)

//go:embed contents/binary_bounds_runner.txt
var binaryBoundsRunner string

func TestGoBinaryDecodingBounds(t *testing.T) {
	for _, endian := range []string{"little", "big"} {
		t.Run(endian, func(t *testing.T) {
			definitions, err := cclParser.ParseCCLSourceContent(&cclParser.CCLParseOptions{SourceContent: fmt.Sprintf(`
#[SerializationType("binary")]
#[BinarySerializationEndian("%s")]
model Child { Number: int64; }
model Empty {}
enum Choice: uint16 { First, Second }
model Envelope {
 IV: bytes;
 Ciphertext: bytes;
 Label: string;
 Nested: Child;
 Children: Child[];
 Labels: string[];
 Chunks: bytes[];
 Numbers: int64[];
 Flags: bool[];
 Choices: Choice[];
 EmptyValue: Empty;
}
[StrictBinaryParsing(false)]
model Loose { IV: bytes; Ciphertext: bytes; }
`, endian)})
			if err != nil {
				t.Fatal(err)
			}
			target := t.TempDir()
			cclLoader.LoadGenerators()
			if _, err = cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
				CodeContext: definitions.CodeContext, OutputPath: filepath.Join(target, "models"), TargetLanguage: "go",
			}); err != nil {
				t.Fatal(err)
			}
			if _, err = RunGoProject(&RunGoOptions{TargetPath: target, RunnerContent: binaryBoundsRunner}); err != nil {
				t.Fatal(err)
			}
		})
	}
}
