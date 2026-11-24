# Definitions and compilation

A compiled declaration turns recognized facts about a Go struct shape into a
typed result that an application can reuse. [Rules](rules.md) decide which
types and members are recognized; the definition decides what to record. A
declaration describes a type, not the state of a populated instance.

## Reusable typed declarations

**Definitions** hold configuration. `define.Struct` constructs one and calls
`Define` to register its rules. **Compilation** applies that configuration to a
target type. Use `compile.Struct[T]` when `T` is known in Go code, or
`compile.StructReflection` when the type comes from an earlier discovery step.
**Declarations** hold the data recorded by successful encounters.

```go
type declaration struct {
    define.ScopeDeclaration
    Field string
}

type definition struct{ define.StructDefinition[declaration] }

func (d *definition) Define() {
    hits := property.Property()
    hits.Expect("field named Hits", func(e property.Encounter) bool {
        return e.Property.FieldName() == "Hits"
    })
    d.Property(options.Required, hits).Encounter(func(e property.Encounter) logic.Check {
        decl, check := e.Access[declaration](define.DeclarationKey)
        if !check.OK() { return check }
        decl.Field = e.Property.FieldName()
        return logic.Ok()
    })
}

d, err := define.Struct[definition, declaration](nil)
if err != nil { panic(err) }
result := compile.Struct[Report](d)
if !result.IsOkay() { panic(result.Trace()) }
decl, check := result.GetDeclaredValue()
if !check.OK() { panic(check.Why()) }
```

The registration's recognizer selects properties; additional rule arguments
to `Property` are expectations on a recognized property. The encounter runs
after they succeed. Required means this registration must match at least one
property in the body, not exactly one. All matching property registrations are
evaluated for each property.

Check `IsOkay()` before using a compilation, and check `GetDeclaredValue()`
before using its result. A failed compilation's partial results must not be
used. Intel does not roll back side effects that caller callbacks already made.
`Trace()` is for people; neither its wording nor layout is a parsing contract.

## Compositional selection

Inside `Define`, `d.Name` and `d.Type` hold type-name and type rules. For
example, `d.Type.Expect(typerule.Type().Is[Report]())` requires exact Go type
identity; `Expect` replaces the expectation list and `Also` appends to it.
Use the public `rules/typerule` package. Likewise, register member rules with
`d.Contract(options.Required, contractRule, validators...)`,
`d.Extension(options.Optional, extensionRule, validators...)`, and
`d.Property(options.Required, propertyRule, validators...)`. Construct these
rules using the respective `rules/contract`, `rules/extension`, and
`rules/property` packages. Each registration takes a recognizer followed by
zero or more expectations, and can receive an `Encounter` callback. Optional
absence is different from a failed rule evaluation.

`define.Struct` starts with breadth-first traversal, open contracts and
extensions, and closed properties. Change composition during `Define` using
`SetContractComposition`, `SetExtensionComposition`, or
`SetPropertyComposition` with `define.Open` or `define.Closed`. Closed
properties reject each unrecognized addressable property; closed contracts
and extensions instead require a match in their respective domains, not
recognition of every candidate. Composition does not normalize fields, tags,
or the declaration into a domain model. Within each body, type rules run
first, then contracts, extensions, and properties, stopping on failure; nested
bodies do not have a single global callback order.

## Per-use carriers for shared bodies

A parent definition can select the same target for multiple recognized fields
while retaining each use's carrier on its host declaration:

```go
scope := d.Scope[body, countersDefinition](fields)
scope.Target(func(property.Encounter) (*capture.StructInterrogator, logic.Check) {
    return capture.Type[Counters]().AsStruct()
})
scope.Encounter(func(_ property.Encounter, host *reportDeclaration, use *define.Carrier[body]) logic.Check {
    host.Uses = append(host.Uses, use)
    return logic.Ok()
})
```

`Scope` and `Leaf` register optional property handlers with a recognizer and
optional validators. By default `Scope` narrows the property to a struct and
calls its `ScopeIndicator.Enter` on a zero value; embedding `define.Scope`
provides an `Enter` that selects the current struct. `Leaf` instead selects the
property itself as a struct without calling `Enter`. Either can use a `Target`
callback to select a different struct body, as above. The carrier still
records the original property; install `Encounter` to receive it and decide
how to attach it to the host.

Within a compilation session the same definition instance and exact target
type share a body, not field metadata. A new top-level compilation starts a
new session. Recursion can link to a registered declaration before its body
has finished being built: linking is safe, treating it as complete is not.
Discard partial results after failure and start a fresh session rather than
reusing them; Intel cannot undo caller callback effects. Successful compilation
does not populate live instances or establish what labels, field order, or the
resulting graph mean. Those decisions and any further validation belong to the
caller.

For member views such as contracts, extensions, and properties, see
[Struct member views](capture.md#struct-member-views).

For exact operation contracts, see Go documentation for `define.Struct`,
`compile.Struct`, `compile.StructReflection`,
`compile.StructCompiler.GetDeclaredValue`, and `define.Carrier`.

Optional complete scenarios: [compile an intake report](../examples/compile-report/readme.md)
and [share a totals layout](../examples/shared-bodies/readme.md).

---

[Previous: Rules](rules.md) · [Next: Reference](reference.md)
