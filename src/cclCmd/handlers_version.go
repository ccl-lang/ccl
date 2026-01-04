package cclCmd

import (
	"fmt"
	"runtime"
	"runtime/debug"

	gValues "github.com/ccl-lang/ccl/src/core/globalValues"
)

func HandleVersionCommand() {
	fmt.Printf("ccl version %s %s/%s\n", gValues.CurrentCCLVersion, runtime.GOOS, runtime.GOARCH)

	// Read the build info embedded by the Go compiler
	if info, ok := debug.ReadBuildInfo(); ok {
		var revision string
		var time string
		var modified bool

		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.time":
				time = setting.Value
			case "vcs.modified":
				if setting.Value == "true" {
					modified = true
				}
			}
		}

		if revision != "" {
			// Use short hash (first 7 chars) for cleaner look
			if len(revision) > 7 {
				revision = revision[:7]
			}
			if modified {
				revision += " (dirty)"
			}
			fmt.Printf("Commit: %s\n", revision)
		}
		if time != "" {
			fmt.Printf("Commit Date: %s\n", time)
		}
	}
}
