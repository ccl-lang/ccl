# ccl-gen
The Common Code Language generator implementation.

ccl is a code generation tool for converting a .ccl source file to certain programming languages. It is mostly used for defining models across servers and clients.

## Grammar Example

A basic example of the ccl grammar is as follows:

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


## Installation

```bash
go install github.com/ALiwoto/ccl@latest
```

## Usage

```bash
ccl generate --source definitions.ccl -o path/to/output -language Go
```
