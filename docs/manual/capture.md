# Capture and live values

Inspection describes types and members. Capture lets you narrow those descriptions
through checked operations and, separately, work with storage in a particular
value. Start with `intel.CaptureType[T]()` for a type, or `intel.CaptureValue(v)`
for a value; narrow with operations such as `AsStruct`, `AsSlice`, and
`AsInterface`. Check `OK()` before using each result. In the illustrations,
`must` stops on a failed check. `Why()` is for people, not a format to parse.

- [Struct member views](#struct-member-views)
- [Field metadata and live storage](#field-metadata-and-live-storage)
- [Slice type and value capabilities](#slice-type-and-value-capabilities)
- [Method capabilities and receiver choice](#method-capabilities-and-receiver-choice)

## Struct member views

`StructInterrogator` offers several views of the same Go struct. The narrower
views let callers work with a kind of member directly instead of checking every
field and deciding what role it has. The same field can belong to several views.

| View | What it gives you |
| --- | --- |
| **Fields** | The general structural view, including the embedding itself and the fields it promotes. It is the starting point when the kind of member is not yet known. |
| **Properties** | Non-embedded, non-blank fields, including promoted ones. This view separates ordinary field declarations from the composition that makes some of them visible. |
| **Embeddings** | Fields declared using Go embedding syntax. This view exposes the composition itself, rather than the ordinary properties an embedding may promote. |
| **Extensions** | Embeddings whose field type is a struct. They expose a nested struct shape. An embedded pointer to a struct is still an embedding, but it is not an extension. |
| **Contracts** | Embeddings whose field type is an interface. They identify an interface field supplied through embedding, not every interface the enclosing type satisfies. |

For example, these declarations differ by both field type and how that type
is used:

```go
type Counters struct{ Hits int }
type Limits struct{ Max int }
type Source interface{ Read() string }

type Report struct {
	Counters
	*Limits
	Source
	Backup Source
}
```

`Field("Counters")` finds the embedding itself; `Property("Hits")` finds a
property promoted through it. `Counters` is an extension; `Limits` is not.
`Source` is a contract, but `Backup` is a property even though both fields
have the same type. Finding `Source` says nothing about whether a particular
`Report` holds a non-nil `Source` value. For method-set compatibility of the
enclosing type, use `shape.SatisfiesType[Source]()` on the captured `Report`
type. `shape.Contracts.Satisfies[Source]()` instead searches the embedded
interface fields for a matching type.

The specialized views let you look up a contract or extension directly. If
you already have a general field handle, `AsEmbeddedInterface` and
`AsEmbeddedStruct` narrow it directly to the corresponding kind.

Classification, visibility, and storage are separate. `Properties.Declared`
is direct-only; visible lookup and traversal can include promoted fields.
`AllAddressable` is a type-level view, not a promise of writable storage on a
live instance. `WithoutDescendantsOf` changes the view used during discovery:
it can keep an embedding while omitting members reached through it.

For the individual contracts, see `go doc` for `capture.StructInterrogator`,
`capture.Properties`, `capture.Embeddings`, `capture.Extensions`, and
`capture.Contracts`. The member views also determine which members are
considered by [definition registrations](compilation.md#compositional-selection).

## Field metadata and live storage

Choose the operation based on what you have now and what you need later:

| Situation | Route |
| --- | --- |
| You only need a field's type, tags, or position. | Start with `intel.CaptureType[T]().AsStruct()`. Use `Field(name)` for a named lookup or `Properties.AllAddressable()` to discover properties. Stop at the metadata. |
| You have the instance and a field descriptor. | Use `intel.CaptureValue(&instance)` and narrow with `AsStruct` if you haven't captured it yet. Pass the descriptor to `StructValueInterrogator.Field`, then use `Address` if you need the field's storage. |
| You discover fields now but will receive instances later. | Save the captured root type and a copy of each field's `Index()`. Use `capture.ResolveField` on each supplied instance, then `Address` to read or write the field. Don't save one instance's address as part of the type-level plan. |

For example, a captured type can find the promoted `Hits int` on `Report`
and retain its index path for `current`, a `Report` supplied later:

```go
root := intel.CaptureType[Report]()
shape, check := root.AsStruct()
must(check)
hits, check := shape.Field("Hits")
must(check)
path := slices.Clone(hits.Index())

field, check := capture.ResolveField(&current, root, path)
must(check)
address, check := field.Address()
must(check)
*address.(*int) = 7 // The discovery step selected a field known to be int.
```

The root type and index path are a reusable recipe for locating a field, not
a reference to storage. Successful lookup during discovery does not mean a
particular instance can be traversed. Handle failures at the point where you
resolve the path, and do not use an address until that check succeeds.

Use `go doc` in the importing module for exact input, path, and storage
contracts:

```sh
go doc github.com/qlustra/intel.CaptureValue
go doc github.com/qlustra/intel/capture.ResolveField
go doc github.com/qlustra/intel/capture.StructValueInterrogator.Field
go doc github.com/qlustra/intel/capture.ValueInterrogator.Address
go doc github.com/qlustra/intel/inspect.FieldDescriptor.Index
```

## Slice type and value capabilities

Slice type metadata lets you reason about element shape even when there are no
elements yet. `AsSlice` on a captured type leads to `AsStructSlice` when the
immediate elements are structs, not pointers to structs. This can guide what
operations you prepare without borrowing any storage from an instance.

A captured slice value serves a different purpose: its `AsSlice` provides
`Elements` for visiting stored elements and `AppendZero` for growing the
captured slice. For example, after the checks succeed, the appended element's
address can be used:

```go
value, check := intel.CaptureValue(&items) // retain the slice header
must(check)
slice, check := value.AsSlice()
must(check)
added, check := slice.AppendZero()
must(check)
address, check := added.Address()
must(check)
address.(*Item).Count = 3
```

Element handles and storage changes have to be coordinated. Finish using an
element address before an append that might reallocate, or reacquire handles
after growth. See `go doc` for `capture.SliceInterrogator.AsStructSlice`,
`capture.SliceValueInterrogator.AppendZero`, and
`capture.SliceValueInterrogator.Elements`.

## Method capabilities and receiver choice

Intel separates a method's shape from the receiver a call uses:

| Capability | What it gives you |
| --- | --- |
| Interface method inspection | A declared signature without a receiver value. |
| `ZeroReceiver` → `Method` → `Getter` → `Read` | A checked zero-input, one-result call on a fresh zero receiver, without an existing instance. |
| `CaptureValue` → `AsStruct` → `Invoke` | A typed callback on the captured instance when its state matters. The callback, not the diagnostic method name, chooses the call. |

For example, a narrow interface keeps method selection in Go code while Intel
checks that the captured instance satisfies it:

```go
type Incrementer interface{ Increment() int }

// structure is a checked capture of the live struct value.
hits, check := structure.Invoke[Incrementer, int]("Increment",
    func(c Incrementer) (int, logic.Check) { return c.Increment(), logic.Ok() })
must(check)
```

Discovering a signature or checking a method set does not prove a call will
succeed. See `go doc` for `capture.InterfaceInterrogator.Method`,
`capture.GetterInterrogator.Read`, and `capture.StructValueInterrogator.Invoke`
for individual failure and call boundaries.

Optional runnable scenarios: [apply a tag-driven report update](../examples/resolve-field/readme.md),
[append a batch of records](../examples/slice-elements/readme.md), and
[read a default and update a live counter](../examples/call-methods/readme.md).

---

[Previous: Inspection](inspection.md) · [Next: Rules](rules.md)
