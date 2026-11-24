# Inspect a report's field catalog

Project the addressable fields of a report type into a deterministic catalog.
The example shows each field's complete index path, Go type, parsed JSON tag
value and options, and whether an empty `label` tag is present. It also reports
the optional outcome of looking up a name promoted equally from two embedded
fields; ambiguity is not resolved by traversal order.

The catalog uses breadth-first visible-field order without filters. It parses
present JSON tags as comma-separated text, not as JSON encoding rules; absent
JSON tags leave the value and options empty. Label presence is kept separately
so a present empty label is not mistaken for an absent one.

Run from `docs/examples`:

```sh
go run ./inspect-report
```

Expected output:

```text
Report fields:
  Counters path=[0] type=main.Counters json-key=json value= options=[] label="" present=false
  Title path=[1] type=string json-key=json value=title options=[] label="" present=false
  Hits path=[0 0] type=int json-key=json value=hits options=[omitempty] label="" present=true
AmbiguousReport Hits: ambiguous or absent
```
