package gdGen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// RunGodotProject runs a Godot project located at targetPath with the provided runnerContent script.
func RunGodotProject(opts *RunGodotOptions) (string, error) {
	// Create project.godot to make it a valid Godot project
	projectGodotContent := `config_version=5

[application]

config/name="CCL GDScript Test"
config/features=PackedStringArray("4.5", "Forward Plus")
`
	if err := os.WriteFile(
		filepath.Join(opts.TargetPath, "project.godot"),
		[]byte(projectGodotContent),
		0644,
	); err != nil {
		return "", fmt.Errorf("failed to write project.godot: %v", err)
	}

	if err := os.WriteFile(filepath.Join(opts.TargetPath, "runner.gd"), []byte(opts.RunnerContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write runner.gd: %v", err)
	}

	cmd := exec.Command("godot", "--headless", "--import", "--path", opts.TargetPath)
	cmd.Dir = opts.TargetPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("godot import step failed: %v\nOutput:\n%s", err, out)
	}
	// Run the generated code using Godot headless
	// Assuming 'godot' is in PATH.
	cmd = exec.Command("godot", "--headless", "--script", "runner.gd")
	cmd.Dir = opts.TargetPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run Godot: %v\nOutput:\n%s", err, output)
	}
	return string(output), nil
}
