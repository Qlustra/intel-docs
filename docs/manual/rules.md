# Rules: recognize and check

Rules let you ask two separate questions about a type or member: does this
subject match, and what do its expectations establish? A match is not itself a
claim that every expectation passed; choose the recognizer and expectations to
express the question your application needs.

## Rule families and their subjects

The family determines what the predicate receives. Choose a family for the
subject, rather than converting every member to a generic field yourself.

| Subject | Rule family |
| --- | --- |
| A string identifier, such as a type name | `rules/identifier` |
| A Go type, such as the type of a property | `rules/typerule` |
| A field descriptor, including its name, type, and tags | `rules/field` |
| An ordinary property or a particular kind of embedding | `rules/property`, `rules/extension`, or `rules/contract`, according to the [member view](capture.md#struct-member-views) |

Some families provide adapters. For example, `property.Field` passes the
property's field descriptor to a field rule without treating every field as
a property.

## Recognition and expectations

For example, recognize a property named `Title`, then check its type and JSON
tag:

```go
namedTitle := property.Property()
namedTitle.Expect("named Title", func(e property.Encounter) bool {
	return e.Property.FieldName() == "Title"
})
stringType := property.Property()
stringType.ExpectCheck("string type", func(e property.Encounter) logic.Check {
	return e.Property.Is[string]()
})

var checks property.Rules
checks.Match(options.Optional, namedTitle)
checks.Expect(stringType, property.Field(field.TaggedWith("json")))
```

`Expect` adapts a predicate; `ExpectCheck` adapts a `logic.Check` and can
retain its diagnostic. The built-in `property.Field` adapter applies a field
rule to a property's field; prefer it when it expresses the check. A field
rule and a property rule see different subjects ([member views](capture.md#struct-member-views)).

## Match versus evaluation validity

`MatchSubject` reports recognition only. `ValidateRules` applies optionality
and reports whether the **set is satisfied** and whether evaluation was
**valid**. Those are different questions:

| Recognition outcome | What `ValidateRules` does |
| --- | --- |
| Valid non-match, optional selector | Reports satisfaction without running expectations. |
| Valid non-match, required selector | Reports non-satisfaction, but evaluation was valid. |
| Failed recognition | Reports failed evaluation, even for an optional selector. |
| Match | Evaluates expectations. Use `ValidateExpectations` instead if recognition already ran and you do not want to repeat it. |

For multiple checks, compose rules in a set: a rule's `Expect`, `ExpectCheck`,
and `ExpectMatch` replace its prior expectation, while a set's `Expect`
replaces its list and `Also` appends to it.

## Nested rules and diagnostics

Use `ExpectMatch` when a rule needs to report evaluation validity separately
from a match, or compose a nested rule that needs the current validation
branch. For example, a custom property rule can pass a child branch to a field
rule and propagate both of its results.
Propagate both results: a failed evaluation is not a valid non-match, and
optionality cannot repair it. For a single `Rule.ValidateRule` call, a failed
encounter after a match can return `matched=true, valid=false`. A validation
branch's status does not replace the returned pair. `Trace()` and `Why()` are
for human diagnostics, not stable formats to parse. Composition may
short-circuit, so a trace need not include every branch.

Use `go doc` for the individual contracts: `rules.Rule.ExpectCheck`,
`rules.Rule.ExpectMatch`, `rules.ValidationBranch`, and the matching and
validation methods in your installed Intel version.

Optional runnable scenarios: [check a report's title schema](../examples/match-rules/readme.md)
and [keep nested field-check diagnostics](../examples/nested-rules/readme.md).

[Previous: Capture](capture.md) · [Next: Compilation](compilation.md)
