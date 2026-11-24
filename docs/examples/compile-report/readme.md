# Validate an intake report shape

An intake service accepts reports with a field marked by an `intake` tag. It
builds the reusable declaration from `Report` once, reads the tag-derived
declaration independently of compilation, and rejects a shape that lacks the
required property with a human-readable trace.

From `docs/examples`, run:

```sh
go run ./compile-report
```

```text
Intake inventory: hits
Missing compiled: false
Missing trace:
!type Missing struct {}: required properties not matched
  type name
  type
  !required properties: definition 1 did not match any property
Missing declaration available: false
```

The program uses only Intel's public `define`, `compile`, `options`,
`logic`, and `rules/property` APIs. The property rule accepts the `intake`
tag; the encounter stores its value in a declaration. The service uses the
compiled declaration to print an intake inventory. The `Missing` compilation
shows the rejection and its human-readable trace; a rejected compilation has
no declared value to read.

[Example bundle](../readme.md)
