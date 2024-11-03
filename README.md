# ccl-gen
The Common Code Language generator implementation.

ccl is a code generation tool for converting a .ccl source file to certain programming languages. It is mostly used for defining models across servers and clients.

## Installation
```bash
go install github.com/ALiwoto/ccl@latest
```

## Usage
```bash
ccl generate --source definitions.ccl -o path/to/output -language Go
```
