# Contributing to CCL

Thank you for being interested in CCL.

CCL is young, and small contributions matter a lot:

- fixing docs
- adding examples
- improving generated output
- reporting confusing syntax
- testing CCL on real projects
- proposing language features through focused issues

## Good First Contributions

If this is your first time here, good areas to start with are:

- documentation fixes
- small examples
- parser tests
- sanitizer tests
- generator tests
- clearer error messages
- generated output comparisons

## Before Large Changes

Please open an issue before starting work on:

- syntax changes
- new target languages
- generator architecture changes
- new attributes
- method or code syntax proposals
- breaking changes to generated output

CCL is a language and code generator, so small design choices can affect every target language. Discussing the direction first avoids wasted work.

## Local Setup

The core project is written in Go.

Install a Go version compatible with `go.mod`, then run:

```bash
go test -v ./...
```

Some runtime tests also use target-language tooling. The CI setup currently installs:

- Python 3.11
- Go stable
- Godot 4.6.0
- Node.js 22
- .NET 6.0 runtime

If you cannot run every runtime locally, run the tests you can and mention the missing toolchain in your pull request.

## Required Checks

Before opening a pull request, run:

```bash
go test -v ./...
```

On Windows, also run:

```powershell
./scripts/EnsureGoContent.ps1
```

On non-Windows systems, run:

```bash
./scripts/EnsureGoContent.sh
```

Fix reported issues before submitting the change.

## Code Rules

The detailed project rules live in `AGENTS.md`. The most important rules for contributors are:

- Do not use reflection in CCL core code or generated output.
- Return clear errors for user input and runtime failures.
- Do not panic for user input or normal runtime conditions.
- Keep Go files separated by category, such as `types.go`, `helpers.go`, `methods.go`, `constants.go`, and `errors.go`.
- Keep core Go source files under 500 lines.
- Use meaningful names instead of short or ambiguous names.
- Put CCL input-language packages under `src/inputLangs/cclInput`.
- Put public documentation under `ccl-lang.github.io/src/content/docs`.

Generated output should be clean and idiomatic for its target language, with as few third-party dependencies as practical.

## Tests

Tests are important for every change.

For parser, sanitizer, or language changes, add focused tests that describe the accepted and rejected syntax.

For generator changes, prefer runtime tests when practical: generate code, run the generated code, and verify behavior inside the generated target language.

For bug fixes, add a regression test that fails without the fix.

## Documentation

If a change affects users, update the public docs site under:

```txt
ccl-lang.github.io/src/content/docs
```

Root-level documentation should stay focused on project entry points and contributor process.

## Pull Requests

Keep pull requests focused. A good pull request usually changes one behavior, one target, or one documentation area.

Please include:

- what changed
- why it changed
- which tests were run
- any target-language behavior that reviewers should inspect

If generated output changes, explain whether the change is intentional and whether it is breaking.
