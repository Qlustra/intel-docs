# Keep nested field-check diagnostics

An intake caller applies a field rule from a property-level check and wants
the nested evaluation's diagnostics in the same validation trace. Run from
`docs/examples`:

```sh
go run ./nested-rules
```

```text
match: matched=true valid=true
non-match: matched=false valid=true
failed evaluation: matched=false valid=false
property validation
  !custom field check: nested_failed (nested evaluation failed)
    nested field checks
      !: rule expectation is not configured
```

The program delegates the nested field evaluation to Intel and places it under
a named branch of the caller's validation for each candidate. An unconfigured
nested rule shows what failed evaluation looks like, in contrast with a valid
non-match. The program prints that failure's trace for a reader, not as a
format to parse. Intel provides rule evaluation and diagnostic recording; the
policy connecting a property to its field check and interpreting the result
belongs to the application.

[Example bundle](../readme.md)
