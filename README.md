# ccl-gen

The Common Code Language generator implementation.

ccl is a code generation tool for converting a .ccl source file to certain programming languages. It is mostly used for defining models across servers and clients.

## Grammar Example

A basic example of the ccl grammar is as follow:

```ccl
// NOTE: Support for attributes is not added as of yet
// This is a global attribute, applied to the whole generation process
#[CCLVersion("1.0.0")]
#[SerializationType("binary")]

// This is a comment

// models can also have attribute on them.
[SerializationType("binary")]
model UserInfo {
    Id: int64;
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

> NOTE: ccl is a very minimal language, currently we do not have support for complex types such as enums, maps, unions, etc.

## Installation

```bash
go install github.com/ALiwoto/ccl@latest
```

## Usage

```bash
ccl generate --source definitions.ccl --output path/to/output --language Go
```

## Programming languages

A list of all programming languages that are either supported or we plan to support in the future are shown here:

- [Golang](https://github.com/ALiwoto/ccl/wiki/Programming-Languages#golang)
- [GDScript](https://github.com/ALiwoto/ccl/wiki/Programming-Languages#gdscript)
- [Python](https://github.com/ALiwoto/ccl/wiki/Programming-Languages#python)
- [CSharp](https://github.com/ALiwoto/ccl/wiki/Programming-Languages#csharp)
- [Rust](https://github.com/ALiwoto/ccl/wiki/Programming-Languages#rust)

If you do not see your desired language in the list, please open an issue and we will consider adding it.

## Contributing

We will be glad to accept any contributions to the project!
