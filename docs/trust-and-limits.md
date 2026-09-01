# Trusting-trust context, evidence, and limits

## What the literature establishes

Ken Thompson's original *Reflections on Trusting Trust* explains that a
compiler can recognize and subvert the source of a compiler or a login
program, while also reproducing the subversion in future compiler binaries.
The important observation for this repository is that agreement between two
successive generations on one compiler lineage is not an independent witness.
The primary publication is [Thompson's 1984 Communications of the ACM
article](https://doi.org/10.1145/358198.358210).

David Wheeler's *Countering Trusting Trust through Diverse Double-Compiling*
(DDC) proposes compiling with a second compiler and then using that result to
compile the compiler source again, comparing the resulting binaries. The
[author-published ACSAC paper](https://dwheeler.com/trusting-trust/wheelerd-trust.pdf)
and the [full dissertation](https://dwheeler.com/trusting-trust/) describe the
assumptions and the practical technique. DDC is stronger than a same-lineage
two-generation replay because it changes the compiler path, but it still
depends on the trusted compiler, source correspondence, build inputs, and
comparison discipline.

The [reproducible-builds definition](https://reproducible-builds.org/docs/definition/)
requires that the same source, build environment, and build instructions let
any party recreate bit-for-bit identical specified artifacts. The
[bootstrappable builds project](https://bootstrappable.org/) makes the related
bootstrap concern explicit: opaque binary seeds make the source-to-binary path
hard to audit, so the project aims to minimize those seeds.

## What this witness changes

The authoritative `.gooo` contract declares one source node and two distinct
parser/lowerer/emitter/executor nodes. `path-a` and `path-b` may share only the
wire schema and the fixture bytes. CI computes their dependency intersection
and fails if another internal package is shared. Each path independently
produces:

1. canonical semantic IR;
2. generated Go source that CI builds and runs; and
3. a terminal reason plus ordered effect trace.

The verifier reports independent indicators for semantic identity, generated
artifact byte identity, terminal-trace identity, and runtime identity. A case
cannot become `CLOSED` because only one of those indicators matches. The
precedence is `REFUTED > UNKNOWN > CLOSED`; unavailable evidence is not silently
treated as convergence.

## Counterevidence and closure conditions

The fixed corpus is a small executable falsification suite:

- semantic IR divergence is `REFUTED`;
- generated artifact byte divergence is `REFUTED`;
- terminal-reason/effect drift is `REFUTED`;
- an absent diverse path is `UNKNOWN` with six required next-operation fields;
- exact replay and comment-insensitive canonical convergence can be `CLOSED`
  only when all required identities and generated runtime outputs match.

The injected path-B case is a controlled mutation, not a claim that every
possible malicious mutation has been modeled. It proves that this contract
does not accept a matching first-generation semantic answer when the second
path produces a different answer. The terminal-drift case separately proves
that equal semantic IR does not authorize ignoring an execution-trace drift.

## What `CLOSED` does not prove

`CLOSED` is scoped to the declared Gooo slice and the evidence recorded by one
GitHub Actions run. It does not prove that:

- the Go 1.27 compiler, linker, standard library, kernel, runner, or CI control
  plane is honest;
- both path implementations were not designed by a correlated adversary;
- a payload outside the fixed grammar or fixture corpus cannot be triggered;
- the generated artifacts correspond to an entire production compiler; or
- byte identity alone proves source-level semantic correctness.

The two paths also share the Go toolchain and the operating environment. That
shared trust is intentionally visible rather than hidden behind the word
"bootstrap". A stronger claim would require more diverse toolchains,
independent environments, source review, signed provenance, and a broader
corpus. This repository makes no global-language self-improvement claim.

## Optional predecessor input and improvements

`gooo-two-generation-bootstrap v0.1.1` is retained only as a digest-pinned
optional input with `required_gate=0`. Its presence cannot close any case.

The repository release policy is recorded in
[docs/release-policy.md](release-policy.md). The endpoint is an
administrator-read API, so a regular Actions GITHUB_TOKEN 403 is recorded as
insufficient observation rather than inferred as enabled=false. The release
workflow instead verifies the public release API after publication and fails
closed unless immutable=true, the exact annotated tag object target, one asset,
and its expected digest all match. GitHub documents that immutable releases
lock the associated tag and assets after publication and recommends creating a
draft, attaching all assets, then publishing it; this repository preserves the
historical mutable v0.1.0 and uses the enabled policy for later releases.

An improvement claim is accepted only for an exact before/after pair with the
same scenario, source digest, contract digest, and toolchain digest, and with
integer measurements on both sides. Scores, weighted averages, and estimated
rates are not evidence of improvement and are forbidden by the `.gooo`
contract.
