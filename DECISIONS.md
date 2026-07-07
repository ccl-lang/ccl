# Project Decisions

Design decisions need to stay clear, practical, and consistent because small syntax or generator changes can affect many target languages.

## Maintainer-Led Decisions

CCL currently uses a maintainer-led decision process.

Discussion, issues, experiments, and pull requests are welcome, but final decisions rest with the project maintainers.

This is not meant to discourage contribution. It exists so the language can keep a coherent direction while it is still small.

## Changes That Need Discussion First

Please open an issue before starting work on:

- syntax changes
- new built-in types
- new attributes
- new target languages
- generator architecture changes
- serialization format changes
- method or code syntax proposals
- inheritance, composition, generics, maps, or other language-level features
- breaking changes to generated output

Small bug fixes, tests, docs fixes, and clearly scoped generator improvements usually do not need a design discussion first.

## How Decisions Are Made

Maintainers evaluate proposals using these questions:

- Does this solve a real user problem?
- Can the syntax stay simple?
- Can errors remain clear?
- Can this be generated cleanly across supported target languages?
- Does this avoid reflection in CCL core code and generated output?
- Does this keep generated code idiomatic and inspectable?
- Does this create long-term maintenance cost?
- Is this the right time for the feature?

Not every good idea should be accepted immediately. Some ideas may be correct for CCL later but too early for the current language.

## Language Design Proposals

Language-level proposals should include:

- the problem being solved
- example CCL syntax
- expected generated behavior
- impact on existing CCL files
- impact on target languages
- alternatives considered
- unresolved questions

Use the language design issue template for early proposals. A more formal RFC process may be added later when the project needs it.

## Compatibility

CCL should avoid breaking existing source files and generated output without a strong reason.

Breaking changes may still happen while CCL is young, but they should be explained clearly and, when practical, documented with migration notes.

## Experiments

Experimental ideas are welcome, especially when they come with tests or generated-output examples.

An experiment can be useful even if it is not accepted. It may reveal target-language problems, syntax issues, or better alternatives.

Experimental features should not be treated as stable until maintainers explicitly document them as supported.

## Reconsidering Decisions

Decisions can be revisited when new information appears:

- real user feedback
- target-language problems
- implementation complexity
- performance problems
- unclear errors
- better syntax or generator designs

Changing direction is acceptable when the reason is clear and the project becomes better because of it.
