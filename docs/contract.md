# Contract summary

The machine-readable authority is
[`.gooo/diverse-bootstrap.gooo`](../.gooo/diverse-bootstrap.gooo). This page is
an orientation aid; if it differs from the metacode, the metacode wins.

The graph is:

`source -> path-a` and `source -> path-b` -> `compare` -> `case-judgment`, with
both generated artifacts also passing through `generated-go-build-run`.

The closure gate requires path availability, canonical semantic identity,
generated-artifact byte identity, terminal reason/effect identity, and runtime
identity. `REFUTED` is emitted for observed divergence, while missing path or
trace evidence is `UNKNOWN`. Every unknown record carries `stage`, `step`,
`reason`, `unknown_class`, `next_operation`, and `blocked_by`.

The root README is excluded only from inventory counts. All other Go/Gooo files,
physical lines, directories, regular files, generated files/bytes, wall-clock
milliseconds, peak RSS KiB, and test totals are reported as integer fields in
the CI evidence artifact.
