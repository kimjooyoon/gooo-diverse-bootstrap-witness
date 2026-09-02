# Contract summary

The machine-readable v3 authority is
[`.gooo/diverse-bootstrap.gooo`](../.gooo/diverse-bootstrap.gooo). This page is
an orientation aid; if it differs from the metacode, the metacode wins.

The v3 contract is append-only: the immutable v0.1.2 v2 denominator of
fourteen cases is retained and enriched with per-case evidence labels. No
released tag, release, or asset is rewritten.

Every fixed case declares exactly one `proof_choice` from
`FOUNDATION`/`COHERENCE`/`REGRESSION` and one `indicator_class` from
`DRIVER`/`OUTCOME`/`GUARDRAIL`. The labels are required in the authoritative
case record, copied into both generated IRs, checked by the verifier, included
in each case proof receipt, and required in release evidence.

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
proof receipts with their independence basis, plus the case-level indicator
class. The contract explicitly records
the Münchhausen trilemma: external source/toolchain trust and correlated
attacks are not erased by witness agreement.

The root README is excluded only from inventory counts. All other Go/Gooo files,
physical lines, directories, regular files, generated files/bytes, and test
totals are reported in the CI evidence artifact. Each required witness and
the overall CI run has actual integer `wall_ms`, `peak_rss_kib`, `build_ms`,
and `test_ms` when observed. Cache `hits`/`misses` are actual integers when
the runner exposes them; otherwise they are explicit `null` values with a
six-field `UNKNOWN` record. Improvement is `null`/`UNKNOWN` unless an exact
matched before/after pair has the same scenario, input, contract, toolchain,
witness, and trial identity; scores and aggregate performance claims are
forbidden.
