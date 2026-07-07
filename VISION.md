# CCL Vision

The long-term vision for CCL is to let developers describe more of a project in a clear, human-readable language and generate clean, idiomatic code for multiple target runtimes.

CCL should make common cross-language work easier without hiding important behavior behind heavy tooling or target-specific configuration files.

## What CCL Is Today

Today, CCL is a schema-like model language and code generator.

It is focused on:

- defining shared data models
- attaching explicit attributes to source definitions
- generating target-language code
- keeping generated output inspectable
- supporting serialization behavior across targets
- helping servers, clients, tools, and game projects share model definitions

This current version of CCL is useful when the same data shapes need to exist in more than one runtime.

## What CCL Is Not Yet

CCL is currently not a general-purpose programming language.

It does not yet have full method bodies, general control flow, a standard library, a runtime, a package ecosystem, or a broad type system comparable to established programming languages.

Those ideas will become part of CCL over time, but they need to be introduced carefully.

## Where CCL Is Going

The intended direction is gradual:

1. Make model definitions and generated output reliable.
2. Improve attributes, imports, enums, serialization, and target-language quality.
3. Add missing data-shape features such as maps and generics.
4. Explore methods, reusable behavior, inheritance, composition, and richer project syntax.
5. Keep every new feature understandable in CCL and practical in generated target languages.

The larger dream is for CCL to describe not only data shapes, but also more of the code and behavior that projects repeat across languages.

## Design Values

CCL should stay guided by these values:

- **Simple source syntax**: CCL files should be easy to read and write.
- **Explicit behavior**: Important generation choices should be visible in source through attributes and clear declarations.
- **Clean generated code**: Output should feel natural in the target language.
- **Low dependency weight**: Generated code should avoid third-party dependencies unless there is a strong reason.
- **Clear errors**: Failing with a useful error is better than silently generating wrong code.
- **No reflection**: CCL core code and generated output should avoid reflection.
- **Target respect**: Every supported language has its own conventions, and generated code should respect them.

## Why Not Existing Tools

CCL is inspired by the need for an easier and more flexible way to describe shared code structures than many existing IDLs provide.

Tools like protobuf, OpenAPI, and other schema systems are useful, but CCL aims for a different balance:

- easier source syntax
- explicit source-level attributes
- clean generated code
- practical multi-target output
- room to grow beyond data schema definitions over time

CCL does not need to replace every existing tool. It should be good at the workflows where its syntax, attributes, and generated output make projects simpler.

## Community Direction

CCL should grow through real use, honest feedback, focused issues, small contributions, and public design discussion.

The best early community is not a crowd around a promise. It is a small group of people who can use CCL today, report what is confusing, and help shape what it becomes next.
