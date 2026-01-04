package gdGen_test

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
	gdSource1      = "definitions1.ccl"
	mainGdContent1 = `
extends SceneTree

func _init():
	print("Starting GDScript runtime verification...")
	
	var user_script = load("res://models/user.gd")
	if user_script == null:
		print("Error: Could not load user.gd")
		quit(1)
		return

	var user = user_script.new()
	user.id = "user123"
	user.username = "test_user"
	
	# Test SkinInfo
	var skin_info_script = load("res://models/skin_info.gd")
	var skin_info = skin_info_script.new()
	skin_info.type = 0

	var basic_skin_script = load("res://models/basic_skin.gd")
	var basic_skin = basic_skin_script.new()
	basic_skin.r = 255
	basic_skin.g = 0
	basic_skin.b = 0

	skin_info.basic = basic_skin
	user.skin = skin_info

	if user.username != "test_user":
		print("Error: Username mismatch")
		quit(1)
		return

	# Test Position (direct access)
	var pos = Position.new()
	pos.x = 10
	pos.y = 20
	pos.z = 30

	if pos.x != 10:
		print("Error: Position X mismatch")
		quit(1)
		return

	print("Runtime verification successful!")
	quit(0)

func _ready():
	quit(1) # Should not reach here
`
)

func TestGdGenerator1(t *testing.T) {
	// Setup output directory
	tmpDir, err := filepath.Abs("ccl_gd_test_1")
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

	fmt.Printf("Generating GDScript code to: %s\n", tmpDir)

	// Parse CCL
	parsedDefinitions, parseErr := cclParser.ParseCCLSourceFile(&cclParser.CCLParseOptions{
		SourceFilePath: gdSource1,
	})
	if parseErr != nil {
		t.Fatalf("Error: failed to parse source file %s: %v\n", gdSource1, parseErr)
		return
	}

	// Generate Code
	cclLoader.LoadGenerators()
	result, err := cclGenerators.DoGenerateCode(&cclGenerators.CodeGenerationOptions{
		CCLDefinition:     parsedDefinitions,
		OutputPath:        filepath.Join(tmpDir, "models"),
		TargetLanguage:    "gd",
		GenerateDebugInfo: true,
	})
	if err != nil {
		t.Fatalf("Error: failed to generate code: %v\n", err)
		return
	} else if result == nil {
		t.Fatalf("Unknown error: failed to generate code")
		return
	}

	fmt.Printf("Running GDScript code from: %s\n", tmpDir)

	output, err := RunGodotProject(&RunGodotOptions{
		TargetPath:    tmpDir,
		RunnerContent: mainGdContent1,
	})
	if err != nil {
		t.Fatalf("Failed to run Godot: %v\nOutput:\n%s", err, output)
		return
	}

	fmt.Printf("Output:\n%s\n", output)
}
