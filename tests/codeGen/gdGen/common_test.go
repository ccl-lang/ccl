package gdGen_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type RunGodotOptions struct {
	TargetPath    string
	RunnerContent string
}

// Helper function to create the command based on OS and environment
func getGodotCmd(args ...string) *exec.Cmd {
	// 1. Try to find 'godot' in the system PATH directly.
	// This handles Linux/Mac and correctly configured Windows machines
	path, err := exec.LookPath("godot")
	if err == nil {
		return exec.Command(path, args...)
	}

	// 2. If LookPath failed, we might be in a complex Windows environment.
	if runtime.GOOS == "windows" {
		// Check if 'bash' is available (GitHub Actions and Git Bash users have this).
		// This solves the "godot not recognized" error on GH Actions.
		if _, err := exec.LookPath("bash"); err == nil {
			// We use bash to run godot.
			// "-c", "godot \"$@\"", "--" allows us to pass arguments safely.
			bashArgs := append([]string{"-c", "godot \"$@\"", "--"}, args...)
			return exec.Command("bash", bashArgs...)
		}

		// 3. Fallback to 'cmd' if bash is missing (Original fallback).
		cmdArgs := append([]string{"/C", "godot"}, args...)
		return exec.Command("cmd", cmdArgs...)
	}

	// Default for Linux/Mac if not found in PATH (will likely fail execution later)
	return exec.Command("godot", args...)
}

// RunGodotProject runs a Godot project located at targetPath with the provided runnerContent script.
func RunGodotProject(opts *RunGodotOptions) (string, error) {
	// Create project.godot to make it a valid Godot project
	projectGodotContent := `config_version=5

[application]

config/name="CCL GDScript Test"
config/features=PackedStringArray("4.6", "Forward Plus")
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

	cmd := getGodotCmd("--headless", "--import", "--path", opts.TargetPath)
	cmd.Dir = opts.TargetPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("godot import step failed: %v\nOutput:\n%s", err, out)
	}
	// Run the generated code using Godot headless
	// Assuming 'godot' is in PATH.
	cmd = getGodotCmd("--headless", "--script", "runner.gd")
	cmd.Dir = opts.TargetPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run Godot: %v\nOutput:\n%s", err, output)
	}
	return string(output), nil
}
