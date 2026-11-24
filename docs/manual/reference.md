# Reference: find the right operation

| If you need to… | Start with |
| --- | --- |
| Compare a Go type or read field and tag metadata | [`inspect`](inspection.md): `CaptureType`, `InspectFields`, `CaptureFieldTag` |
| Work with a specific category of struct member | [Struct member views](capture.md#struct-member-views): fields, properties, embeddings, extensions, and contracts |
| Narrow a type or work with an instance | [`intel` and `capture`](capture.md): `CaptureType`, `CaptureValue`, `ResolveField` |
| Select a subject and apply expectations | [`rules` and its rule families](rules.md): start with the family for the subject, such as `rules/property` or `rules/field` |
| Describe and compile a struct shape | [`define` and `compile`](compilation.md): `define.Struct` to build the definition, `compile.Struct` to use it |
| Choose traversal order or optionality | `options`: `DepthFirst`, `BreadthFirst`, `Required`, `Optional` |

Use `go doc` for signatures from a module that imports Intel. For example:

```sh
go doc github.com/qlustra/intel/capture.ResolveField
go doc github.com/qlustra/intel/rules.ValidationBranch
```

## Read outcomes at the right layer

- For a result paired with a `logic.Check`, use `check.OK()` before using
  the result. Show `check.Why()` to explain failure.
- When evaluating a rule, check evaluation validity separately from whether
  the rule matched. A valid non-match need not be an error.
- For struct compilation, check `IsOkay()` and use `Trace()` to explain a
  failure. Check the result from `GetDeclaredValue()` before using the
  declaration.

Diagnostic text and trace layout are for readers, not formats to parse. See
[v1 guarantees](../guarantees.md) for storage, identity, ordering,
and session boundaries.

[Previous: Definitions and compilation](compilation.md) · [Contents](readme.md)
· [Example bundle](../examples/readme.md)
