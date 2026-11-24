# Inspection: type and field facts

`inspect` reports type and field facts without capturing a live value. Go
package documentation describes each operation contract. Start with
`go doc github.com/qlustra/intel/inspect`, then look up individual operations,
for example `go doc github.com/qlustra/intel/inspect.FieldsInspector.LookupNamed`.

## Exact type identity and method sets

`type HitCount int` defines a distinct type, even though Go allows conversion
to `int`. With an `AddHit` method declared on `*Report`, a pointer to
`Report` satisfies `hitAdder`, while the value type does not:

```go
count := inspect.CaptureType(reflect.TypeFor[HitCount]())
same := count.Is[int]().OK()
assignable := count.AssignableToType[int]()
convertible := count.ConvertibleToType[int]()

valueMatches := inspect.Satisfies[hitAdder](reflect.TypeFor[Report]()).OK()
pointerMatches := inspect.Satisfies[hitAdder](reflect.TypeFor[*Report]()).OK()
```

`Is` compares exact Go types; it doesn't compare underlying kinds.
`Satisfies` checks the method set of the type you give it, without substituting
a pointer for a value. This step neither runs a conversion nor calls a method.

## Field visibility and complete paths

Create a field inspector for a known struct type. A fresh descriptor registry
and no field filters let you compare declared and visible views. For the
high-level member categories, see [struct member views](capture.md#struct-member-views).

```go
reportType := inspect.CaptureType(reflect.TypeFor[Report]())
fields := inspect.InspectFields(reportType, inspect.NewDescriptorRegistry(), nil, nil)
hits, check := fields.LookupNamed("Hits")
must(check)
```

```go
type Counters struct{ Hits int }
type Report struct {
	Counters
	Title string
}
```

`AllDeclared` reports `Counters` and `Title` in declaration order.
`AllAddressable` also includes `Hits`; depth-first puts it after `Counters`,
while breadth-first puts it after `Title`. `LookupNamed` follows Go's selection
rules: two equally promoted fields with the same name are ambiguous, and
changing traversal order does not resolve the ambiguity. `Hits` keeps its
complete path `[0 0]`, which can be used in a
[saved field path](capture.md#field-metadata-and-live-storage).

## Tag values, parsing, and presence

For a `Hits` field tagged `json:"hits,omitempty" label:""`, `Lookup("label")`
returns `("", true)`; an absent key returns `("", false)`.
Keep the presence bit when an empty string and a missing key mean different
things. Use `Parse` when you want a structured value rather than raw text:

```go
tags := hits.Tags()
tag, check := tags.Parse[formats.CSV]("json")
must(check)
label, present := tags.Lookup("label")
missing, found := tags.Lookup("missing")
```

`formats.CSV` splits the value on commas; it does not apply JSON encoding
semantics. For a custom format, supply `Parse(string)` and
`Validate() logic.Check` methods. Check the result of parsing before using the
format. `Raw` alone cannot distinguish missing from present-empty.

[Contents](readme.md) · [Next: Capture and live values](capture.md)

Optional example: [inspect a report's shape and tags](../examples/inspect-report/readme.md).
