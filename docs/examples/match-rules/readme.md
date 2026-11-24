# Check a report's Title schema

An intake caller allows reports to omit `Title`. If the property exists, it
must be a string with a `json` tag. Unrelated fields are not subject to those
checks. Run from `docs/examples`:

```sh
go run ./match-rules
```

```text
tagged title and unrelated fields: accepted=true valid=true
no title: accepted=true valid=true
non-string title: accepted=false valid=true
untagged title: accepted=false valid=true
```

The program visits each candidate's visible properties. An optional selector
recognizes `Title`; the type and tag rules run only on a match. The caller
accepts a candidate if all property evaluations satisfy the rule set. A wrong
type or missing tag rejects the candidate without making evaluation invalid.
Intel supplies recognition and validation; the caller supplies the title
policy and chooses which outcome permits processing.

[Example bundle](../readme.md)
