package pyGen_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ccl-lang/ccl/src/cclGenerators"
	"github.com/ccl-lang/ccl/src/cclLoader"
	"github.com/ccl-lang/ccl/src/cclParser"
)

const (
	pySource1 = "definitions1.ccl"
)

func TestPythonGenerator1(t *testing.T) {
	// Setup output directory
	tmpDir, err := filepath.Abs("ccl_py_test_1")
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Clean up previous run
	if err := os.RemoveAll(tmpDir); err != nil {
		t.Fatalf("Failed to remove existing dir: %v", err)
	}

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}

	fmt.Printf("Generating Python code to: %s\n", tmpDir)

	// Parse CCL
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: pySource1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", pySource1, parseErr)
		return
	}

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CCLDefinition:  parsedDefinitions,
		OutputPath:     filepath.Join(tmpDir, "models"),
		TargetLanguage: "python",
	})
	if err != nil {
		t.Fatalf("Error: failed to generate code: %v\n", err)
		return
	} else if result == nil {
		t.Fatalf("Unknown error: failed to generate code")
		return
	}

	// Create verify.py to test generated code
	verifyPyContent := `
import sys
import os
import struct

# Add the current directory to sys.path to ensure we can import the generated models
sys.path.append(os.getcwd())

from models import User, SkinInfo, BasicSkin, Position

def test_user():
	print("Testing User...")
	user = User()
	user.id = "user123"
	user.username = "test_user"
	
	skin = SkinInfo()
	skin.type = 0
	
	basic_skin = BasicSkin()
	basic_skin.r = 255
	basic_skin.g = 0
	basic_skin.b = 0
	
	skin.basic = basic_skin
	user.skin = skin

	# Test serialization
	user_binary = user.serialize_binary()
	if not user_binary:
		raise Exception("Failed to serialize user")

	# Test deserialization
	deserialized_user = User.deserialize_binary(user_binary)
	if not deserialized_user:
		raise Exception("Failed to deserialize user")

	if user.username != "test_user":
		raise Exception("Username mismatch")
	if deserialized_user.username != "test_user":
		raise Exception("Deserialized username mismatch")
	
	if deserialized_user.skin.basic.r != 255:
		raise Exception("Deserialized skin color mismatch")

	print("User test passed!")

def test_position():
	print("Testing Position...")
	pos = Position()
	pos.x = 10
	pos.y = 20
	pos.z = 30

	if pos.x != 10:
		raise Exception("Position X mismatch")

	# Test serialization
	pos_binary = pos.serialize_binary()
	if not pos_binary:
		raise Exception("Failed to serialize position")

	# Test deserialization
	deserialized_pos = Position.deserialize_binary(pos_binary)
	if not deserialized_pos:
		raise Exception("Failed to deserialize position")

	if deserialized_pos.x != 10:
		raise Exception("Deserialized Position X mismatch")

	print("Position test passed!")

if __name__ == "__main__":
	try:
		test_user()
		test_position()
		print("Runtime verification successful!")
	except Exception as e:
		print(f"Verification failed: {e}")
		sys.exit(1)
`
	output, err := RunPythonProject(&RunPythonOptions{
		TargetPath:    tmpDir,
		RunnerContent: verifyPyContent,
	})
	if err != nil {
		t.Fatalf("Failed to run generated code: %v\nOutput:\n%s", err, output)
		return
	}
	fmt.Printf("Output:\n%s\n", output)
}
