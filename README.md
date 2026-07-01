# ccl-gen

The Common Code Language generator implementation.

ccl is a code generation tool for converting a .ccl source file to certain programming languages. It is mostly used for defining models across servers and clients.

## Grammar Example

A basic example of the ccl grammar is as follow:

```ccl
// This is a global attribute, applied to the whole generation process
#[CCLVersion("1.0.0")]
#[SerializationType("binary")]

// by default, the binary serialization endianness is little,
// you can uncomment this line and use this attribute to change it.
// #[BinarySerializationEndian("big")]

// models can also have attribute on them.
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

// models support multiple attributes
[SerializationType("binary")]
[SerializationType("json")]
model GetUsersResult {
    Users: UserInfo[];
    OtherUsers: UserInfo[];
}
```

> NOTE: currently we do not have support for maps, unions, etc.

## Syntax features

  - [attributes](https://ccl-lang.github.io/docs/attributes/)
  - [imports](https://ccl-lang.github.io/docs/imports/)
  - [enums](https://ccl-lang.github.io/docs/enums/)

## Installation

```bash
go install github.com/ccl-lang/ccl@latest
```

## Usage

```bash
ccl generate --source definitions.ccl --output path/to/output --language Go
```

## Programming languages

A list of all programming languages that are either supported or we plan to support in the future are shown here:

- [Go](https://ccl-lang.github.io/docs/languages/go/)
- [GDScript](https://ccl-lang.github.io/docs/languages/gdscript/)
- [Python](https://ccl-lang.github.io/docs/languages/python/)
- [TypeScript](https://ccl-lang.github.io/docs/languages/typescript/)
- [JavaScript](https://ccl-lang.github.io/docs/languages/javascript/)
- [CSharp](https://ccl-lang.github.io/docs/languages/csharp/)
- [Rust](https://ccl-lang.github.io/docs/languages/rust/)

If you do not see your desired language in the list, please open an issue and we will consider adding it.

## Contributing

We will be glad to accept any contributions to the project!
