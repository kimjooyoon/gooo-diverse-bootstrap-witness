# Gooo diverse bootstrap witness

`gooo-diverse-bootstrap-witness` compares two independently implemented
Gooo lowerer/executor paths. The same `.gooo` semantic source is parsed by
`internal/patha` and `internal/pathb`; each produces a canonical semantic IR,
a runnable generated Go artifact, and a terminal reason/effect trace.

The `.gooo` file at
[`.gooo/diverse-bootstrap.gooo`](.gooo/diverse-bootstrap.gooo) is the
authority for the grammar, normalization, generation graph, independence
predicate, comparison indicators, resolution precedence, trace policy, fixed
cases, and the optional digest-pinned `gooo-two-generation-bootstrap v0.1.1`
input. Go implements the parser, lowerers, executors, and verifier; it does not
replace the metacode rules.

The conformance corpus intentionally contains three `CLOSED` convergences, one
detected semantic injection (`REFUTED`), one detected terminal-reason drift
(`REFUTED`), and one unavailable diverse path (`UNKNOWN`). `REFUTED > UNKNOWN >
CLOSED` is used for every case. The suite is `CLOSED` only when every fixed
case receives its declared judgment and all required CI evidence is present.

Artifact byte identity and semantic identity are separate indicators. Neither
one can close a case by itself. A `CLOSED` case also needs equal terminal
traces and equal runtime output from the built generated artifacts.

## Run through CI

GitHub Actions uses Go 1.27 and is the verification authority. It checks the
independence import intersection, runs Go tests, generates the six cases into a
caller-owned temporary directory, builds and runs every available generated Go
artifact, compares runtime output, and uploads an evidence JSON artifact.

The local development contract intentionally does not require local
test/build/vet/conformance runs. Generated output never enters the repository.
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
counterevidence boundary.
