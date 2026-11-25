# Intel Docs

Learn to use Intel, a Go reflection library, to inspect struct fields and tags, work with live values, and validate struct shapes. Start with the usage manual or run a standalone example.

## What it is

- A usage manual covering inspection, capture, rules, and compilation.
- Standalone Go programs you can run independently.
- A draft of Intel’s v1 guarantees. The library source lives in a separate repository.

## Why it exists

A promoted field can be visible on a Go type even when a nil embedded pointer makes it unreachable on a particular value. The manual and examples show what each Intel check establishes before you rely on it in your own code.

## Core idea

For Go developers building tools or frameworks that use reflection:

- Inspect a type to learn its structure; capture a value to work with a particular instance.
- Write rules for the shapes your application recognizes.
- Compile declarations from types, not populated instances.

## Setup

No setup is needed to read the docs. To run the examples against the public release, use Go 1.27 or newer and wait until Intel `v1.0.0` is available:

```sh
git clone https://github.com/Qlustra/intel-docs.git
cd intel-docs/docs/examples
go mod download github.com/qlustra/intel
```

## Quick start

Open the [field catalog example](docs/examples/inspect-report/readme.md) to see how Intel finds promoted fields and reads their tags. From `docs/examples`, run:

```sh
go run ./inspect-report
```

Look for `Hits path=[0 0]` in the output: `Hits` is promoted, but its index path still goes through the embedded struct.

## Project shape

- `docs/manual/`
- `docs/examples/`
- `docs/guarantees.md`

## Status

The manual and examples target Intel `v1.0.0`; the v1 guarantees are still a draft.

## Links

- [Usage manual](docs/manual/readme.md)
- [Example bundle](docs/examples/readme.md)
- [Draft v1 guarantees](docs/guarantees.md)
- [Intel library source](https://github.com/Qlustra/intel)

## License

The documentation and examples are licensed under the [Apache License 2.0](LICENSE).
