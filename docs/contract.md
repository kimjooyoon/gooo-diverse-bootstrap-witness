# Contract summary

The machine-readable v2 authority is
[`.gooo/diverse-bootstrap.gooo`](../.gooo/diverse-bootstrap.gooo). This page is
an orientation aid; if it differs from the metacode, the metacode wins.

The v2 contract is append-only: the released v1 denominator of six cases is
retained and eight cases are added, for fourteen fixed cases. No released tag,
release, or asset is rewritten.

The graph is:

`source -> path-a` and `source -> path-b` -> `compare` -> `case-judgment`, with
both generated artifacts also passing through `generated-go-build-run`. The
source node is the `.gooo` semantic kernel; Go supplies only bounded runtime
operations.

The closure gate requires path availability, exact kernel/input identity, two
independent witness identities, canonical IR/decision/provenance identity,
generated-artifact byte identity, terminal reason/effect identity, and runtime
identity. Same-lineage replay, forged digest, witness disagreement,
self-approval cycle, and frozen-bootstrap mismatch are `REFUTED`; missing
witness/lineage/toolchain/input identity is `UNKNOWN`. Every unknown record
carries `stage`, `step`, `reason`, `unknown_class`, `next_operation`, and
`blocked_by`.

Each machine evidence record includes FOUNDATION, COHERENCE, and REGRESSION
proof receipts with their independence basis. The contract explicitly records
the Münchhausen trilemma: external source/toolchain trust and correlated
attacks are not erased by witness agreement.

The root README is excluded only from inventory counts. All other Go/Gooo files,
physical lines, directories, regular files, generated files/bytes, and test
totals are reported in the CI evidence artifact. Each required witness has
actual CI-observed `wall_ms`, `peak_rss_kib`, `build_ms`, and `test_ms` on the
same canonical input. Improvement is `null`/`UNKNOWN` unless an exact matched
before/after pair has the same scenario, input, contract, toolchain, witness,
and trial identity; scores and aggregate performance claims are forbidden.
