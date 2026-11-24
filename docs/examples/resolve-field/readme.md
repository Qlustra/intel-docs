# Apply a tag-driven report update

A report holds metric counters and unrelated local state. The application
wants to put incoming readings into fields marked with a `metric` tag, while
leaving untagged fields alone. It builds a plan once from the report type and
reuses it for each report it receives.

From `docs/examples`, run:

```sh
go run ./resolve-field
```

```text
plan: hits at [0 0]
plan: misses at [0 1]
first: hits=12 misses=3 untouched=99
second: hits=4 misses=1 untouched=88
incomplete update rejected; first unchanged
different root type rejected
```

[The program](main.go) separates discovery from updating:

1. `buildPlan` scans visible properties, including promoted ones. The program
   requires each `metric` tag to name a unique, non-empty key on an exported
   `int` field. It saves independent copies of the paths, not addresses into
   any particular report.
2. `apply` looks up each reading, resolves its path on the supplied report,
   and collects its destination address. Only after every lookup and
   resolution succeeds does it write the batch. Missing data leaves the report
   unchanged; a newly defined type with the same structure is still rejected.

Intel handles inspection and path resolution. The tag meaning and the choice
to check the entire batch before writing are application policy, not behavior
provided by Intel.

[Field-path guidance in the manual](../../manual/capture.md#field-metadata-and-live-storage)
· [Example bundle](../readme.md)
