# Scoped Attributes

This document describes the intended CCL attribute syntax and the proposed
compiler architecture for resolving scoped attributes. It is a design contract;
the implementation may still lag behind parts of this document.

## Syntax

CCL supports normal local attributes on declarations:

```ccl
[SomeAttribute]
model User {
    Id: int64;
}
```

Scoped attributes start with `#` and are not attached to the immediately
following declaration:

```ccl
#[SomeAttribute]
#global:[SomeAttribute]
#file:[SomeAttribute]
#namespace:[SomeAttribute]
```

`#[SomeAttribute]` is an alias for `#global:[SomeAttribute]`. The short form is
the preferred spelling for global attributes.

Scope names are case-insensitive, so these are equivalent:

```ccl
#file:[SomeAttribute]
#File:[SomeAttribute]
#FILE:[SomeAttribute]
```

## Attribute Scopes

`global`
: Applies to the whole compilation context. This is the existing behavior of
  `#[SomeAttribute]`.

`file`
: Applies only to the source file where the attribute appears.

`namespace`
: Applies to the current namespace. Child namespaces inherit the attribute
  unless they define a matching namespace-scoped attribute of their own.

Current CCL code is considered to be in namespace `main` when no namespace is
declared.

```ccl
namespace main;
#namespace:[SomeAttribute("root")]

namespace main.users;
#namespace:[SomeAttribute("users")]
```

In this example, `main.users` overrides the matching namespace attribute from
`main`. A namespace under `main.users` inherits the `main.users` value unless it
overrides the attribute again.

## Language Selectors

Any scoped or local attribute can optionally be limited to one or more target
languages:

```ccl
#file:[$Go:SomeAttribute("hello from Go")]
#file:[$CSharp:SomeAttribute("hello from C#")]
```

Language names are case-insensitive and should resolve through the same alias
table used by code generation targets:

```ccl
#File:[$go:SomeAttribute("hello from Go")]
#File:[$csharp:SomeAttribute("hello from C#")]
```

Multiple languages can share one attribute declaration:

```ccl
#file:[$go,$js:SomeAttribute("hi from Go and JS")]
#file:[$go, $js:SomeAttribute("hi from Go and JS")]
```

Whitespace is insignificant around language-selector separators.

The normalized internal representation should store languages as
`globalValues.LanguageType` values. An attribute with no language selector
applies to all languages.

## Resolution Semantics

Attribute lookup should be explicit about the target declaration, target
language, attribute name, and whether the attribute is allowed to fall back to
scoped attributes.

For model-level lookup:

1. Check attributes directly on the model.
2. If the requested attribute can be inherited from scoped attributes, resolve
   scoped attributes in this order:
   1. file scope for the model's owning file
   2. current namespace
   3. parent namespace, continuing upward
   4. global scope

For field-level lookup:

1. Check attributes directly on the field.
2. Check attributes directly on the owning model when the attribute supports
   model-level fallback.
3. If the requested attribute can be inherited from scoped attributes, use the
   same scoped lookup order as model-level lookup.

For scoped lookup, each level should search only matching attributes:

1. attribute name matches
2. target language matches, or the attribute has no language selector
3. any future attribute-specific constraints match

Once a matching attribute is found at a lower scope, resolution stops and parent
scopes are not considered for that attribute name. This gives lower scopes a
clear override rule.

When multiple attributes with the same name exist in the same effective scope,
the resolver should preserve all matches for collection-style APIs and return
the first source-order match for single-attribute APIs. We may later add
attribute metadata to distinguish single-value override attributes from
repeatable attributes.

## Source Files And Context

`SourceCodeDefinition` should represent exactly one `.ccl` source file. It is
not thread-safe and should not be used as a merged compilation unit.

The shared compilation state belongs in `CCLCodeContext`. It is thread-safe and
is the right owner for cross-file indexes:

- all type definitions by namespace-qualified name
- incomplete type definitions
- global and automatic variables
- source-file definitions by absolute file path
- namespace-scoped attributes
- global-scoped attributes

Each parsed file should produce its own AST and sanitized
`SourceCodeDefinition`. Import resolution should build a dependency graph of
source files instead of merging imported AST nodes into the importer AST.

This keeps file-local state file-local, makes concurrent parsing and
sanitization possible, and prepares the compiler for future executable code
inside models or namespaces.

## Proposed Internal Model

Keep `SourceCodeDefinition` as the owner of file-local data:

- absolute source file path
- file namespace active for declarations without an explicit namespace
- imports
- custom type definitions declared in that file
- file-scoped attributes

Move shared indexes and cross-file resolution into `CCLCodeContext`:

- register one `SourceCodeDefinition` per file
- register each type into the context type cache
- register global-scoped attributes into a context-level global attribute index
- register namespace-scoped attributes into a context-level namespace index

Model and field definitions should retain enough source ownership information
to resolve file-scoped attributes later. A direct pointer to the owning
`SourceCodeDefinition` is simple, but an immutable source-file ID or absolute
path is also acceptable if we want to keep model values easier to serialize.

## Import And Concurrency Plan

The import resolver should become graph-oriented:

1. Resolve each source path to an absolute canonical path.
2. Parse each unique file once.
3. Detect cycles using the active path stack.
4. Store each file AST separately.
5. Sanitize each file into one `SourceCodeDefinition`.
6. Register definitions and scoped attributes into `CCLCodeContext`.

Parsing can run concurrently once imports are discovered. Sanitization can also
be parallelized per file if registration into `CCLCodeContext` remains
synchronized and validation that depends on all files runs after registration.

The current behavior of merging imported AST attributes and models into one AST
should be retired. It blurs file ownership, breaks file-scoped attributes, and
makes future real code harder to compile incrementally.

## Generator API Direction

Generators should stop treating `SourceCodeDefinition` as the entire program.
The preferred input should be:

- `CCLCodeContext` as the compilation-wide index and resolver
- a root `SourceCodeDefinition` or list of source files as the generation entry
  set
- target language and output options

Attribute helper methods should move toward resolver-style APIs, for example:

```go
ResolveAttribute(targetLang, attrName, subject, options)
ResolveAttributes(targetLang, attrName, subject, options)
```

The subject can be a field, model, source file, or namespace. Options can
describe whether model fallback, scoped fallback, or global fallback is allowed.

For backward compatibility, existing generator helpers can delegate to the new
resolver first, then gradually be replaced at call sites.

## Output File Groups

Generators may use `OutputFileGroup` to route generated artifacts into grouped
output files.

```ccl
#file:[$go:OutputFileGroup("users")]
```

The attribute is resolved per type definition, not per source file. File scope
is only one convenient way to apply the group to all models in a file. The same
attribute may also be applied to a model directly or inherited through namespace
or global scope:

```ccl
#[OutputFileGroup("shared")]

namespace main.users;
#namespace:[$go:OutputFileGroup("users")]

[$go:OutputFileGroup("auth")]
model LoginRequest {
    Id: int64;
}
```

For Go generation, a model whose resolved group is `users` writes type and
method artifacts to files such as `types_users.go` and `methods_users.go`.
Models without a resolved group keep the default files, such as `types.go` and
`methods.go`.

Constants are routed by artifact semantics, not by a blanket rule. Currently,
Go model ID constants remain in `constants.go`.
