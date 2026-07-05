package rsGen_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RunRustOptions struct {
	TargetPath    string
	RunnerContent string
}

func RunRustProject(opts *RunRustOptions) (string, error) {
	if err := writeRustCargoToml(opts.TargetPath); err != nil {
		return "", err
	}

	testsPath := filepath.Join(opts.TargetPath, "tests")
	if err := os.MkdirAll(testsPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create Rust tests dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(testsPath, "runtime_test.rs"), []byte(opts.RunnerContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write Rust runtime test: %v", err)
	}

	cmd := exec.Command("cargo", "test", "--quiet")
	cmd.Dir = opts.TargetPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run Rust generated code: %v\nOutput:\n%s", err, output)
	}

	return string(output), nil
}

func writeRustCargoToml(targetPath string) error {
	packageName := strings.ReplaceAll(filepath.Base(targetPath), "_", "-")
	cargoTomlContent := `[package]
name = "` + packageName + `"
version = "0.1.0"
edition = "2021"

[dependencies]
base64 = "0.22"
serde = { version = "1", features = ["derive"] }
serde_json = "1"
`
	if err := os.WriteFile(filepath.Join(targetPath, "Cargo.toml"), []byte(cargoTomlContent), 0644); err != nil {
		return fmt.Errorf("failed to write Cargo.toml: %v", err)
	}

	return nil
}
