# Intel — v1 guarantees

This document describes what a caller can rely on at each layer of the
architecture, and what a successful check does not imply.
Intel supplies checked operations; the caller supplies domain meaning.

## Inspection

- **Type identity and interface satisfaction are separate.** `Is` tests exact
  Go type identity. `Satisfies` tests interface satisfaction for that type's
  method set. Neither says whether a call on a particular value will succeed.
- **Visibility and traversal order are separate.** Field lookup follows Go's
  promotion, shadowing, and ambiguity rules. Choosing depth-first or
  breadth-first traversal changes the order of visible fields, not which
  fields are visible. Promotion does not shorten a field's index path.
- **A missing tag is not an empty tag.** Tag lookup can distinguish an absent
  key from one with an explicitly empty value. The caller decides what either
  means.

## Interrogation

- **Capture preserves the distinction between a value and its storage.**
  Capturing a non-nil pointer retains access to the supplied storage. Capturing
  a non-pointer makes a shallow copy; referenced data can still be shared.
- **A field path is resolved against a particular instance.** `ResolveField`
  takes a non-nil pointer to the exact expected struct type and checks the
  path. On success, the handle refers to writable storage in that instance.
  Resolution does not allocate nil intermediate pointers or prove where the
  path came from.
- **Slice element handles refer to actual elements.** `Elements` captures
  them in index order. `AppendZero` returns a handle to the zero element it
  appended. If an append reallocates the backing array, prior handles no
  longer refer to elements of the current slice; reacquire them when needed.
- **Invocation is distinct from discovery.** A getter calls its method on a
  fresh zero receiver. `Invoke` passes the captured instance to the caller's
  typed callback; the method name labels diagnostics, not dispatch. A panic
  from the getter call or the `Invoke` callback is reported as a failed
  check. Neither operation promises that the call has no side effects.

## Rules

- **No match is different from failed evaluation.** An optional registration
  can permit a valid non-match. It does not turn a failed evaluation into a
  successful one. An `ExpectMatch` callback can use its supplied
  `rules.ValidationBranch` to include nested rule diagnostics in the same trace.
- **Closed composition depends on the domain.** Closed property composition
  rejects each unrecognized addressable property. Closed contract and
  extension composition instead require a match in their respective domains;
  finding one is not a check that every candidate was recognized.

## Definitions and compilation

- **Bodies are shared within a session.** `Compile` reuses the body for the
  same definition instance and exact target type. A carrier associates a
  declaration with a particular field use, so sharing the body does not merge
  the fields or their metadata. A new top-level compilation starts a new
  session.
- **A recursive link can precede completion.** A declaration is registered
  before its rules run, so encounters can link to the body while it is still
  being built. Within one body, validation runs type rules, then contracts,
  extensions, and properties, stopping on failure. That order is local to the
  body, not a global order for callbacks across nested bodies.
- **A failed session's partial results must not be reused.** Compilation
  does not roll back the caller's callback effects. Successful compilation
  produces a declaration; it does not populate an instance or prove that
  later domain validation or population will succeed.

## Reading outcomes

Check `OK()` before using the result of a checked operation. For struct
compilation, use `IsOkay()` and check the outcome of `GetDeclaredValue()`
when retrieving the declaration. `Why()` and `Trace()` explain failures to
people; their text and layout are not formats for programs to parse.
