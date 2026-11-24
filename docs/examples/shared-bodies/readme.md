# Share a totals layout across report fields

An intake report has two totals fields with the same `Counters` layout but
different labels. The application wants one reusable body definition, a
distinct carrier for each field use, and the shared body declaration to drive
the output while each carrier retains its own field policy metadata.

From `docs/examples`, run:

```sh
go run ./shared-bodies
```

```text
First.Hits label=first
Second.Hits label=second
```

The public `define` API selects `Counters` as the target for both scope uses.
Each use receives a separate carrier, and both carriers point to the same
body declaration. The body declaration supplies the shared `Hits` output;
the original property on each carrier supplies its own field name and `label`
tag. Thus the shared body defines layout policy without collapsing per-use
metadata.

[Example bundle](../readme.md)
