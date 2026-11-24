# Read a default and update a live counter

This example uses a zero receiver to read a counter's default value, then
invokes `Increment` on a live counter and shows that its state changed. Finally,
it asks Intel to invoke `Increment` through an incompatible string-returning
interface and reports that Intel rejected the call boundary.

The policy is to check every `logic.Check` before using its result. A rejected
call is an expected outcome here, so the example reports the rejection rather
than treating it as a program failure. Run from `docs/examples`:

```sh
go run ./call-methods
```

Output:

```text
zero receiver: 0
invoked: 5, live hits: 5
incompatible call rejected: true
```
