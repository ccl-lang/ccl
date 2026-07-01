package gdGen_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type RunGodotOptions struct {
	TargetPath    string
	RunnerContent string
}

const godotCommandTimeout = 30 * time.Second

// Helper function to create the command based on OS and environment
func getGodotCmd(args ...string) *exec.Cmd {
	// 1. Try to find 'godot' in the system PATH directly.
	// This handles Linux/Mac and correctly configured Windows machines
	path, err := exec.LookPath("godot")
	if err == nil {
		return newGodotCommand(path, args...)
	}

	// 2. If LookPath failed, we might be in a complex Windows environment.
	if runtime.GOOS == "windows" {
		if path, err := lookPathWithPowerShellPathext("godot"); err == nil {
			return newGodotCommand(path, args...)
		}

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

func newGodotCommand(path string, args ...string) *exec.Cmd {
	if runtime.GOOS == "windows" && strings.EqualFold(filepath.Ext(path), ".ps1") {
		powerShellArgs := append(
			[]string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-File", path},
			args...,
		)
		return exec.Command("powershell.exe", powerShellArgs...)
	}

	return exec.Command(path, args...)
}

func lookPathWithPowerShellPathext(fileName string) (string, error) {
	currentPathext := os.Getenv("PATHEXT")
	if hasPathext(currentPathext, ".PS1") {
		return exec.LookPath(fileName)
	}

	if currentPathext == "" {
		os.Setenv("PATHEXT", ".PS1")
	} else {
		os.Setenv("PATHEXT", currentPathext+";.PS1")
	}
	defer os.Setenv("PATHEXT", currentPathext)

	return exec.LookPath(fileName)
}

func hasPathext(pathext string, targetExt string) bool {
	for _, currentExt := range strings.Split(pathext, ";") {
		if strings.EqualFold(strings.TrimSpace(currentExt), targetExt) {
			return true
		}
	}

	return false
}

func runGodotCommand(targetPath string, stepName string, timeout time.Duration, args ...string) (string, error) {
	cmd := getGodotCmd(args...)
	cmd.Dir = targetPath

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	if err := cmd.Start(); err != nil {
		return output.String(), err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case err := <-done:
		return output.String(), err
	case <-timer.C:
		_ = killGodotProcessesForProject(targetPath)
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		waitErr := <-done
		_ = killGodotProcessesForProject(targetPath)
		if waitErr != nil {
			return output.String(), fmt.Errorf("%s timed out after %s: %v", stepName, timeout, waitErr)
		}

		return output.String(), fmt.Errorf("%s timed out after %s", stepName, timeout)
	}
}

func killGodotProcessesForProject(targetPath string) error {
	if runtime.GOOS != "windows" {
		return nil
	}

	projectPath := filepath.Clean(targetPath)
	powerShellScript := fmt.Sprintf(
		`$target = '%s'; `+
			`Get-CimInstance Win32_Process | `+
			`Where-Object { $_.Name -like '*godot*' -and $_.CommandLine -like "*$target*" } | `+
			`ForEach-Object { Stop-Process -Id $_.ProcessId -Force -ErrorAction SilentlyContinue }`,
		strings.ReplaceAll(projectPath, `'`, `''`),
	)

	return exec.Command("powershell.exe", "-NoProfile", "-Command", powerShellScript).Run()
}

// RunGodotProject runs a Godot project located at targetPath with the provided runnerContent script.
func RunGodotProject(opts *RunGodotOptions) (string, error) {
	// Create project.godot to make it a valid Godot project
	projectGodotContent := `config_version=5

[application]

config/name="CCL GDScript Test"
config/features=PackedStringArray("4.6", "Forward Plus")

[wgodot]

gdscript/strict_type_checking=false
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

	out, err := runGodotCommand(
		opts.TargetPath,
		"godot import step",
		godotCommandTimeout,
		"--headless",
		"--import",
		"--path",
		opts.TargetPath,
	)
	if err != nil {
		return "", fmt.Errorf("godot import step failed: %v\nOutput:\n%s", err, out)
	}
	// Run the generated code using Godot headless
	// Assuming 'godot' is in PATH.
	output, err := runGodotCommand(
		opts.TargetPath,
		"godot script step",
		godotCommandTimeout,
		"--headless",
		"--path",
		opts.TargetPath,
		"--script",
		"runner.gd",
	)
	if err != nil {
		return "", fmt.Errorf("failed to run Godot: %v\nOutput:\n%s", err, output)
	}
	return string(output), nil
}
