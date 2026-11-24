# Intel examples

These standalone programs accompany the [usage manual](../manual/readme.md).
Use Go 1.27 or newer. Run each command from `docs/examples`; no program depends
on output or state from another program. This module targets Intel `v1.0.0`,
which must be available to run against the public release.

On a fresh checkout, run this once from `docs/examples`:

```sh
go mod download github.com/qlustra/intel
```

| Task | Command |
| --- | --- |
| [Inspect a report's field catalog](inspect-report/readme.md). | `go run ./inspect-report` |
| [Apply a tag-driven update plan](resolve-field/readme.md) to multiple reports. | `go run ./resolve-field` |
| [Append a batch of source records](slice-elements/readme.md) to an existing slice. | `go run ./slice-elements` |
| [Read a counter's default and update a live counter](call-methods/readme.md). | `go run ./call-methods` |
| [Validate an optional title](match-rules/readme.md) without rejecting unrelated fields. | `go run ./match-rules` |
| [Keep nested field-check diagnostics](nested-rules/readme.md) in one trace. | `go run ./nested-rules` |
| [Validate and inventory an intake report shape](compile-report/readme.md). | `go run ./compile-report` |
| [Share a totals layout across report fields](shared-bodies/readme.md) without merging their labels. | `go run ./shared-bodies` |
