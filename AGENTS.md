# Agents guide for CCL

CCL stands for "Common Code Language". It is a language designed to be easily written by anyone and then used for generating code in any output programming language. CCL is meant to be a high-level, human-readable language that abstracts away the complexities of specific programming languages, allowing users to focus on the logic and structure of their code rather than syntax.

The reason CCL was made is because protobuf and other similar IDLs (Interface Definition Language) were either too hard to work with, or just weren't strong enough. The main goal with CCL is to make it "easy" for the developers to get all the hard works done without extra headaches.
The output code **should** also be "clean" and "idiomatic" in the output programming languages AND not using third-party libraries as much as possible.

This is the main repository for CCL, and it contains the core components of the language, including the parser, code generator, and various tools for working with CCL, written in Go programming language.

Organizing things is very very important. Both docs files and source code files should be organized in a clear and consistent manner (e.g. do NOT just put everything in the same folder. Make sure to create proper sub-folders in the proper folder when you want to create a new file).

## Enforcement script

Make sure to run `./scripts/EnsureGoContent.ps1` (or the .sh version for non-windows systems) and fix its reported issues/violations.

## Code Rules
While working on the code, please make sure to follow these rules:

### Rule 0: No Reflection
CCL does not allow the use of reflection in any of its agents. Reflection is a powerful feature in some programming languages that allows a program to inspect and modify its own structure and behavior at runtime. However, it can lead to code that is difficult to understand, maintain, and debug. It will also increase CPU usage and lower code performance. By disallowing reflection, CCL encourages developers to write clear and straightforward code that is easier to reason about.

Using reflection in both CCL's core code and ANY of its output programming languages is strictly prohibited. This rule applies to all agents working on CCL-related tasks, including those responsible for parsing, code generation, and any other aspect of the language's implementation.

Please note that reflections can have different types in different programming languages, e.g. in Python: `getattr(obj, 'attribute_name')`, in Go: `reflect.ValueOf(obj).FieldByName("FieldName")`. All types of reflection are prohibited in CCL.

### Rule 1: Failing with a clear error message is better than silently doing the wrong thing
Imagine you expect an array to have length of 10, if the array does not have length of 10, it is better to fail with a clear error message that says "Expected array length of 10, but got X" rather than silently doing the wrong thing and potentially causing more issues down the line (e.g. setting a variable to 0).

### Rule 2: File separations by category in Go
If you read my Go code, you might notice file names such as this:
- `helpers.go`: Contains ONLY helper functions.
- `types.go`: Contains ONLY type definitions.
- `constants.go`: Contains ONLY constant definitions.
- `methods.go`: Contains ONLY method implementations.
- `methods_binary.go`: Contains ONLY method implementations that are related to binary serialization/deserialization.
- `methods_json.go`: Contains ONLY method implementations that are related to JSON serialization/deserialization.
- `handlers.go`: Contains ONLY handler functions for handling specific tasks or operations.
- `vars.go`: Contains ONLY variable definitions.
- `errors.go`: Contains ONLY error variables (e.g. `var ErrInvalidInput = errors.New("invalid input")`).

My philosophy behind this file separation is to keep the code organized and easy to navigate. By categorizing the code into different files based on their purpose, it becomes easier for developers to find what they are looking for and understand the structure of the codebase. This also promotes better code readability and maintainability, as developers can quickly identify where specific functionality is implemented without having to sift through a large monolithic file.

The "main" category MUST be inside of the package name (folder name). E.g.:
- `src/cclGenerators/csGenerator/`: Which contains constants.go, helpers.go, etc...
- `src/inputLangs/cclInput/cclParser/`: Which contains constants.go, helpers.go, etc...

If the code you wrote does not have these styles, you will be asked to refactor it to follow this style; this is a strict rule.
The only place this rule doesn't apply is to the tests files and generated output languages other than Go.


### Rule 3: Long files are forbidden that is longer than 500 lines is forbidden in CCL-codegen
This is ok inside of the output files, since they are literally generated code, but in the core code of CCL, any file that is longer than 500 lines is forbidden. If you write a file that is longer than ~500 lines, you should refactor the content into multiple smaller files with category name in them (e.g. move helper functions to `helpers_xyz.go`, `helpers_abc.go`, types to `types_xyz.go`, etc...)


### Rule 4: Don't panic (except for internal misuse)
Panicking is WRONG for user input or runtime conditions in both CCL's core code and ANY of its output programming languages; it makes users confused and makes the code harder to maintain. Instead of panicking, return an error with a clear message that explains what went wrong, where and how to fix it (if possible, suggest a solution in the error message).

Panics are allowed ONLY for internal programmer misuse in utility code (e.g. broken invariants like misuse of a fluent builder). These cases must be explicitly agreed by core developers and documented in the relevant package docs.

### Rule 5: Tests
Tests are very very important. Whenever you want to implement a new feature, fix a bug, or refactor some code, you should write tests for it. (if the feature already has tests, you should make sure to run them and they are passing before and after your changes).

Please keep the tests organized in a clear and consistent manner. The rule #2 does not apply to the tests files, since they are usually shorter and more focused on specific functionality. However, you should still try to keep them organized and easy to navigate. You can use descriptive names for your test functions and group them by functionality or feature.

**IMPORTANT**: CCL's tests currently include runtime tests. because CCL's nature is **about** generating code that will be used in runtime, runtime tests are ALWAYS preferred (tests that generate code and then actually run the generated code while testing things INSIDE of that code).

### Rule 6: Meaningful names
Avoid using short and meaningless names for variables, functions, types, etc... instead, use descriptive and meaningful names that clearly indicate their purpose and functionality (both in code-gen, output and docs). Avoid using shortened names if they are not understandable.

### Rule 7: Avoid ambiguous/common type names
Do not introduce generic names like `File`, `Node`, `Context`, or `Result` as public types when a more specific name is possible. Use names that convey domain and avoid confusion with stdlib types (e.g. `CCLFileAST` instead of `File`).


### Rule 8: Git is not for LLM agents
Git is used by the code reviewers to review the codes generated by the LLM agents, but it is not for the LLM agents themselves. Do NOT use/rely on git commands such as `git commit` or `git push`.

### Rule 9: Prefer codeBuilder mapped variables in generators
When writing generator code that emits dynamic source lines, prefer `codeBuilder.MapVarPairs` with `LineD`/`AppendD` over manual string concatenation inside `WriteLine`/`AppendLine`. This keeps generator templates readable and consistent with the existing codebase.

Small local string composition is acceptable for computing generator-side names (for example deriving `fieldName + "Value"`), but emitted code lines should use mapped placeholders whenever practical.

### Rule 10: CCL language packages live under cclInput
All packages that implement the CCL input language itself must live under `src/inputLangs/cclInput`. This includes lexer, parser, AST, sanitizer, language IR, attributes, language errors, and language utilities.

Top-level `src/ccl*` packages outside `cclInput` are reserved for the CCL compiler/project infrastructure. For example, `src/cclCmd` is the command-line interface layer, `src/cclLoader` registers generators, and `src/cclGenerators` emits target-language code.

### Rule 11: Public docs belong in the public docs site
If information is useful to CCL users, contributors, or anyone reading the codebase publicly, document it under `ccl-lang.github.io/src/content/docs` instead.

Write the docs in a way that it's easy to understand for the users and easy to maintain for the developers as well.

