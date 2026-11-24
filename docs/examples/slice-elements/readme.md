# Append source records to a slice

## Task

Append each record from a source batch to an existing slice, then print the
complete result in order.

## Policy

Capture a pointer to the slice variable so appends that grow the slice update
the original header. Append one zero element at a time and immediately copy the
source record into it. Do not retain element addresses across appends: an
append may reallocate the backing array. Once all appends are complete, acquire
fresh element handles with `Elements` to visit the current slice.

## Run

From `docs/examples`:

```sh
go run ./slice-elements
```

## Output

```text
1: Linus (1)
2: Ada (2)
3: Grace (3)
records: 3
```
