# CCL

CCL stands for Common Code Language.

CCL is a small source language for defining shared models once and generating clean target-language code from them. It is useful when the same data shapes need to exist across servers, clients, tools, or game projects without manually rewriting those models in every runtime.

Today, CCL is focused on schema-like model definitions, attributes, enums, imports, and code generation. It is not a general-purpose programming language yet. The long-term direction is to grow carefully toward richer project descriptions, including explicit attributes, reusable methods, inheritance/composition, and more expressive code syntax without losing the simplicity that makes CCL useful now.

## Syntax Example

A basic CCL file looks like this:

```ccl
// This is a global attribute, applied to the whole generation process.
#[CCLVersion("1.0.0")]
#[SerializationType("binary")]

// By default, binary serialization uses little-endian byte order.
// Uncomment this attribute when big-endian output is needed.
// #[BinarySerializationEndian("big")]

// Models can also have attributes.
[SerializationType("binary")]
model UserInfo {
    enum UserType: uint32 {
        Unknown,
        NormalUser,
        Tester,
        Admin,
        Owner,
    }

    Id: int64;
    Type: UserType;
    Username: string;
    Email: string;
    ProfileImage: bytes;
    CreatedAt: datetime;
    UpdatedAt: datetime;
}

// Models can use multiple attributes.
[SerializationType("binary")]
[SerializationType("json")]
model GetUsersResult {
    Users: UserInfo[];
    OtherUsers: UserInfo[];
}
```

CCL currently supports model fields, arrays, nested models, enums, attributes, imports, binary serialization, and JSON serialization. Maps, unions, generic models, and method/code syntax are planned but should not be assumed to exist yet.

## Installation

```bash
go install github.com/ccl-lang/ccl@latest
```

## Generate Code

```bash
ccl generate --source examples/definitions.ccl --output ./generated --language Go
```

The general workflow is:

1. Write model definitions in `.ccl`.
2. Choose a target language.
3. Generate target-language source files.
4. Keep the `.ccl` file as the source of truth.

## Target Languages

CCL can generate code for multiple runtimes. See the language docs for target-specific behavior:

- [Go](https://ccl-lang.github.io/docs/languages/go/)
- [GDScript](https://ccl-lang.github.io/docs/languages/gdscript/)
- [Python](https://ccl-lang.github.io/docs/languages/python/)
- [TypeScript](https://ccl-lang.github.io/docs/languages/typescript/)
- [JavaScript](https://ccl-lang.github.io/docs/languages/javascript/)
- [CSharp](https://ccl-lang.github.io/docs/languages/csharp/)
- [Rust](https://ccl-lang.github.io/docs/languages/rust/)

If you do not see your desired language in the list, please open an issue with the use case and the target runtime requirements.

## Documentation

- [Getting Started](https://ccl-lang.github.io/docs/get-started/)
- [Language Overview](https://ccl-lang.github.io/docs/language-overview/)
- [Attributes](https://ccl-lang.github.io/docs/attributes/)
- [Imports](https://ccl-lang.github.io/docs/imports/)
- [Enums](https://ccl-lang.github.io/docs/enums/)
- [Serialization](https://ccl-lang.github.io/docs/serialization/)
- [Roadmap](https://ccl-lang.github.io/docs/roadmap/)
- [CLI Reference](https://ccl-lang.github.io/reference/cli/)

## Project Direction

CCL's short-term goal is practical, inspectable code generation for shared models. Generated code should feel natural in each target language and should avoid unnecessary third-party dependencies as much as possible.

The larger vision is to let developers describe more of a project in clear CCL syntax: attributes, data models, methods, reusable behavior, and eventually richer code structure. That work needs to happen slowly because every new language feature has to map cleanly into each supported target.

## Contributing

Contributions are welcome, especially:

- documentation fixes
- small examples
- parser and sanitizer tests
- generator tests
- clearer error messages
- generated output improvements
- language design proposals

For large changes, please open an issue first, especially for syntax changes, new target languages, generator architecture changes, new attributes, or method/code syntax proposals.
