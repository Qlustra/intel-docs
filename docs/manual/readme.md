# Intel usage manual

Intel is a Go reflection library for inspecting types and values through
checked operations. You can use the inspection primitives by themselves, or
build on them to recognize and compile a shape of your own.

## Why separate the questions?

An interface match says something about a type's method set, not whether a
call on an instance will work. A promoted field can be visible in a type even
if reaching it on a particular instance would require following a nil pointer.
And a compiled declaration is not a populated instance. Intel gives these
questions separate operations and uses Go's rules to answer them.

Consider a `Report` with an embedded `Counters` struct. First inspect the
promoted `Hits` field and its index path. Then use that path with a particular
`Report`, or write rules that recognize `Hits` and compile a declaration from
its type. The same shape connects the layers; what each operation establishes
is different.

## Read in layer order

| Layer | Capabilities |
| --- | --- |
| [Inspection](inspection.md) | Exact types and method sets, field visibility and order, tag presence and formats. |
| [Capture and live values](capture.md) | Struct member views, saved field paths and live storage, slice elements, method discovery and receiver choice. |
| [Rules](rules.md) | Recognition and expectations, optionality, evaluation validity, nested diagnostics. |
| [Definitions and compilation](compilation.md) | Typed declarations, member composition, shared bodies and per-use carriers. |
| [Reference](reference.md) | Package entry points, guarantees, and diagnostic output. |

Use the manual to choose and combine operations. Use `go doc` in your own
module for the contract of an individual operation in the Intel version you
are using. The example bundle is optional: open it to run a workflow, not to
fill in an API requirement.

## Version and examples

Use Go 1.27 or newer. To use the library in your own Go module once the v1 tag
is public:

```sh
go get github.com/qlustra/intel@v1.0.0
```

The [example bundle](../examples/readme.md) depends on the same version. Each
example starts with a task and runs independently; the bundle provides setup
and run commands.
