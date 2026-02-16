package csGen_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type RunCSOptions struct {
	TargetPath    string
	RunnerContent string
	CSProjContent string
	ProjectName   string
}

// RunCSProject compiles and runs a C# project located at targetPath with the provided runnerContent script.
func RunCSProject(opts *RunCSOptions) (string, error) {
	// Write the runner script
	if err := os.WriteFile(
		filepath.Join(opts.TargetPath, "Program.cs"),
		[]byte(opts.RunnerContent), 0644,
	); err != nil {
		return "", fmt.Errorf("failed to write Program.cs: %v", err)
	}

	projFileName := opts.ProjectName + ".csproj"
	if err := os.WriteFile(
		filepath.Join(opts.TargetPath, projFileName),
		[]byte(opts.CSProjContent),
		0644,
	); err != nil {
		return "", fmt.Errorf("failed to write %s: %v", projFileName, err)
	}

	// Run the project using dotnet run
	cmd := exec.Command("dotnet", "run", "--project", projFileName)
	cmd.Dir = opts.TargetPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run C# project: %v\nOutput:\n%s", err, output)
	}

	return string(output), nil
}
