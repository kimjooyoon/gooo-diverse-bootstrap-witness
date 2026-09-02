# Gooo diverse bootstrap witness

`gooo-diverse-bootstrap-witness` compares two independently implemented
Gooo lowerer/executor paths. The same `.gooo` semantic kernel declares the
kernel source, canonical input, stage0/reference witness, generated/current
witness, optional diverse input, expected observation, and bounded evaluation
rules. `internal/patha` and `internal/pathb` each produce a canonical semantic
IR, a runnable generated Go artifact, and a terminal reason/effect trace.

The `.gooo` file at
[`.gooo/diverse-bootstrap.gooo`](.gooo/diverse-bootstrap.gooo) is the
authority for the grammar, normalization, generation graph, independence
predicate, comparison indicators, resolution precedence, trace policy, fixed
cases, and the optional digest-pinned `gooo-two-generation-bootstrap v0.1.1`
input. Go implements the parser, lowerers, executors, and verifier; it does not
replace the metacode rules.

The v3 conformance contract is append-only from the immutable v0.1.2 v2
denominator of fourteen fixed cases. Each case carries one authoritative
`proof_choice` (`FOUNDATION`, `COHERENCE`, or `REGRESSION`) and one
`indicator_class` (`DRIVER`, `OUTCOME`, or `GUARDRAIL`). Those labels are
copied into the case IR, verifier observation, proof receipts, and release
evidence. The fourteen-case corpus contains normal convergence,
missing witness/lineage/toolchain/input identity (`UNKNOWN`), semantic and
trace disagreement, forged digest, same-lineage replay, self-approval cycle,
and frozen-bootstrap mismatch (`REFUTED`). `REFUTED > UNKNOWN > CLOSED` is
used for every case. The suite is `CLOSED` only when every fixed case receives
its declared judgment and all required CI evidence is present.

Artifact byte identity and semantic identity are separate indicators. Neither
one can close a case by itself. A `CLOSED` case also needs equal decision and
provenance digests, two independently identified witnesses, equal terminal
traces, and equal runtime output from the built generated artifacts. Identical
bytes from one lineage are not two witnesses.

## Run through CI

GitHub Actions uses Go 1.27 and is the verification authority. It checks the
independence import intersection, records per-witness and overall integer
wall/build/test/RSS measurements, records cache hits/misses as integers or
`null` with a six-field `UNKNOWN`, runs Go tests, generates the fourteen cases
into a caller-owned temporary directory, builds and runs every available
generated Go artifact, compares runtime output, and uploads an evidence JSON
artifact. The release workflow selects the successful CI evidence for the
exact merge commit, creates a draft with both source and evidence-dossier
assets, publishes it, then fails closed unless the public release API reports
immutable=true, the exact annotated tag object target, both expected assets,
and both expected digests.

The local development contract intentionally does not require local
test/build/vet/gofmt/actionlint/bash validation or conformance runs. Generated output never enters the repository.
The input repository is measured as `repository_writes=0`; commit, push, and
merge are human-authorized repository operations, not generated-program
authority.

## Scope and limits

This is a witness for a small language slice, not a proof that a host, kernel,
Go 1.27 toolchain, standard library, CI runner, or both implementations are
honest. The two paths share only the declared wire schema and fixture bytes;
the predicate does not remove trust in their source review or in the build
toolchain. The controlled path-B mutations demonstrate that semantic and trace
drift are observable; they do not model every possible correlated attacker.

See [`docs/trust-and-limits.md`](docs/trust-and-limits.md) for the Thompson,
DDC, reproducible-build, and bootstrappable-build references and the precise
counterevidence boundary. Improvement remains `null`/`UNKNOWN` without an
exact before/after pair sharing scenario, input, contract, toolchain, witness,
and trial identity. The machine evidence records FOUNDATION, COHERENCE, and
REGRESSION proof choices and independence grounds while acknowledging the
Münchhausen trilemma.
