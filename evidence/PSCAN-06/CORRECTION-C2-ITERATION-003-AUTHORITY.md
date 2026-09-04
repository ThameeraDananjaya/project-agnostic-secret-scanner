# PSCAN-06 Correction C2 iteration 003 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Iteration: `003`
- Title: Private admission boundary and bounded process correction
- Event: exact owner approval of the bounded correction described below
- Owner command:
  `APPROVE PSCAN-06 C2 ITERATION-003 PRIVATE ADMISSION BOUNDARY AND BOUNDED PROCESS CORRECTION`
- Date: `2026-09-05`
- Approval session: `01a06e82-c28c-7571-b961-3781c255b374`
- Branch: `main`
- Required starting HEAD before this authority bundle:
  `ee5280f00d26b59be532ca6cbb8c8f489285ad36`
- Required starting tree:
  `ad2787e5bacf942f5b8fc2cdeac750a3b9419073`
- Correction authority commit: this authority-bundle commit
- Pre-change tracked and untracked Git status: clean
- Preserved ignored workspace material: `graphify-out/**`; not read as product
  authority, modified, staged or included in this bundle
- State after commit: owner-approved; unclaimed; not implemented; not accepted
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

This is a bounded iteration inside the already-activated PSCAN-06 Correction
C2. It does not create, select or activate another task. A fresh session must
claim only Correction C2 iteration 003 from this committed authority bundle
before implementation. Because the current main checkout contains preserved
ignored workspace material, that fresh session must use an exact clean
authority checkout or another source-trust-compliant isolated worktree; a
normal clean `git status` alone is not sufficient.

## Controlling and preserved authority

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Original Correction C2 authority:
  `evidence/PSCAN-06/CORRECTION-C2-AUTHORITY.md`
- Original authority SHA-256:
  `3CE8E20151993A8FC77890E51C5BF1557F68040CAFF2C3DB1A46237792A75D63`
- Iteration-002 authority:
  `evidence/PSCAN-06/CORRECTION-C2-ITERATION-002-AUTHORITY.md`
- Iteration-002 authority SHA-256:
  `2FAD7986A1E4BF11D0F469CE2FA0334B6B0E7BE0A0C1372D50483F5846115173`
- Unaccepted iteration-002 implementation candidate:
  `52f7ee22ab722d7590e2d2c8326e14a0e9670462`
- Candidate tree:
  `475f186da2a2685740430afef905c633e6123ac7`
- Iteration-002 author-validation record:
  `evidence/PSCAN-06/CORRECTION-C2-ITERATION-002-AUTHOR-VALIDATION.md`
- Author-validation SHA-256:
  `1CAA3512BCB82ACB57C92DC257DA2A5AA10D95F64E1CFF59FDFBE97BD0CE3DBE`
- Locked product tag: `v1.0.0`
- Locked product-source commit:
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Locked Correction C1 tooling tag: `release-tooling-v1.0.0-c1`
- Locked Correction C1 tooling commit:
  `3fb7592889820fa2739a4a53588e073689621809`
- Proposed but absent Correction C2 tooling tag:
  `release-tooling-v1.0.0-c2`
- Pinned build-image identity:
  `docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`

The original C2 authority, iteration-002 authority, candidate and author
evidence, both locked tags, existing schemas and all predecessor evidence
remain immutable facts. Iteration-002 author validation is not acceptance and
is superseded as a readiness claim by the missing trust and resource-boundary
proof recorded here.

## Demonstrated gaps

### Production admission is not private from its test seam

Candidate `52f7ee2` places mutable script-scope callbacks named
`ReleaseDockerInvoker`, `ReleaseHostCacheCanaryInvoker` and
`ReleaseHostOnlyCrlfProofInvoker` in the same production file that makes the
admission decision. The production path conditionally trusts callback results
instead of necessarily executing the real Docker, host-cache and CRLF proof
boundaries. The deterministic harness changes those callbacks directly.

The committed workflow currently invokes a fixed exact entrypoint, so this
record does not claim that a remote bypass was executed. The gap is that the
production admission implementation itself contains a mutable proof-substitution
surface. A fake-engine pass through that surface does not independently prove
that the immutable production entrypoint cannot be switched to synthetic
evidence. Missing isolation at this admission boundary fails closed.

### Declared Docker process limits are enforced only after capture

The production Docker runner starts unbounded `ReadToEndAsync` operations for
stdout and stderr. It compares the completed strings with the declared
`131072` limit only after the process exits and both streams have been fully
materialized. A process can therefore consume memory beyond the stated output
limit before rejection.

The nominal `15000` millisecond wait is also not a complete wall-clock or
process-tree bound. After timeout, tree termination errors are ignored, the
runner performs an unbounded `WaitForExit()`, and stream completion can remain
dependent on inherited pipe handles. The existing fake-result timeout and
oversized-string cases do not execute this real process lifecycle. Therefore
they do not prove bounded output, bounded termination or absence of surviving
descendants.

## Approved objective

Make the release-image admission decision a private production boundary and
make every process resource limit real at the point of consumption. The exact
production entrypoint must always run the real committed host prerequisites and
real Docker process boundary; callers and test code must not replace those
proofs through mutable variables, callbacks, exported functions, dot-sourced
state or ambient session state. Tests may exercise classification through a
separate non-production harness, but test doubles must be structurally unable
to reach the production admission decision used by the workflow.

For each Docker command, enforce fixed argument, executable, environment,
standard-stream byte, wall-clock and descendant-process bounds while the
process is running. Overflow, timeout, start failure, read failure, termination
failure, incomplete stream closure or uncertain descendant cleanup must reject
with no pull or later action. No rejected command may be retried or reclassified
as image absence.

## Allowed implementation paths

Correction C2 iteration 003 may change only:

```text
.github/workflows/release-recovery-v1.0.0.yml
build/release/image-admission.ps1
build/release/admit-image.ps1
build/release/test-image-admission.ps1
docs/release/OFFLINE-VERIFICATION-RUNBOOK.md
docs/release/PSCAN-06-RELEASE-PLAN.md
docs/validation/PSCAN-06-VALIDATION.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/**
```

This list supersedes the broader Correction C2 and iteration-002 path lists for
iteration 003 only. Every absent path is forbidden. A need for a new production
module, helper, binary, dependency or another path stops this iteration and
returns to the owner with exact evidence; authority to broaden scope is not
implied.

## Mandatory correction design

1. The workflow-facing `admit-image.ps1` entrypoint must construct and invoke
   only the real committed production boundary. It must not accept or consult a
   caller-provided command runner, callback, script block, module object,
   environment-selected implementation or mutable ambient override.
2. Remove production-path use of all three mutable iteration-002 callback
   seams. Tests must not dot-source the production boundary and then replace
   its command or prerequisite behavior. Classification tests may use a
   separate test-only model, but parity with production parsing and state
   transitions must be established byte-for-byte or by exact shared immutable
   logic without reopening injection.
3. Keep the executable fixed to `docker` and pass every Docker argument through
   an argument list. Preserve the exact command allowlist and exact canonical
   digest. Shell execution, command strings, aliases, functions and PATH-based
   PowerShell command resolution must not replace the intended native process
   boundary. If executable identity cannot be established sufficiently for the
   accepted environment, reject before pull.
4. Preserve the existing maximum of `131072` bytes independently for stdout
   and stderr, but enforce each limit during capture. Do not fully materialize,
   truncate, silently discard or continue consuming over-limit output before
   rejection. Define encoding and byte accounting deterministically; a
   multi-byte boundary split must not bypass or ambiguously change the limit.
5. Preserve the existing `15000` millisecond per-command ceiling as a complete
   monotonic wall-clock budget covering start, concurrent stream capture,
   process exit and termination/cleanup. A separate strictly bounded cleanup
   allowance may be used only if it is fixed, documented and never converts
   timeout or cleanup uncertainty into success.
6. On timeout, output overflow, read failure or cancellation, terminate the
   complete process tree. Prove bounded return and no live child/grandchild from
   the adversarial fixture. If process-tree termination or pipe closure cannot
   be proved, fail terminally and permit no later admission action.
7. Capture stdout and stderr concurrently without deadlock. Return a command
   result only after exit, both bounded streams and cleanup state are complete.
   Reject partial, missing, extra, mixed, malformed or uncertain lifecycle
   evidence before parsing it as engine, inventory, identity or pull evidence.
8. Preserve the iteration-002 three-state admission model. Only conclusive
   structured absence from a responsive engine may permit one exact pull. The
   correction adds no retry, fallback, mutable reference, registry login or
   alternative absence inference.
9. Preserve host cache and host-only CRLF proofs before admission; every
   post-admission Docker command remains exact-digest, `--pull=never` and
   `--network none`; acquisition and both reproducibility builds retain their
   existing order and bounds.
10. Preserve manifest schema `2.1`, every earlier schema, product/tooling role
    separation, signature identity policy, locked tags and every C2 sandbox,
    credential, remote, publication and owner gate. This iteration changes no
    release identity, scanner behavior or consuming-project boundary.

## Required local checks and success criteria

1. Static and runtime proof shows that the workflow-facing production
   entrypoint has no mutable callback, test-double or caller-provided proof
   substitution path. Attempts to predefine variables/functions, dot-source
   helpers, pass script blocks or set ambient override state cannot replace the
   real prerequisites or Docker process runner.
2. The existing fake-engine matrix still proves exact present/no-pull,
   conclusive-absence/one-pull and all untrusted/zero-pull outcomes through a
   test-only boundary whose state-machine parity with production is proved.
3. A deterministic no-Docker process fixture emits exactly below, at and above
   the stdout byte cap; repeats those cases for stderr; and emits both streams
   concurrently. Below and exact-bound cases complete without truncation; the
   first over-bound condition terminates and rejects within the fixed budget.
4. UTF-8 multi-byte sequences split across read boundaries cannot bypass byte
   accounting or produce accepted malformed evidence. Invalid encoding,
   incomplete final sequences and stream read errors reject deterministically.
5. Deterministic fixtures cover immediate exit, nonzero exit, start failure,
   hang, output flood, simultaneous stdout/stderr flood, child-held pipe and a
   child/grandchild that outlives the parent. Every adverse case returns within
   the fixed total bound, leaves no live fixture descendant and records zero
   pull or later action.
6. Failure to terminate a process tree or complete both streams is explicitly
   represented as terminal untrusted evidence; it is never reported as a
   normal timeout result that admission code could reinterpret.
7. The real process-boundary suite runs on Windows and actual Linux. A
   fake-result unit case alone is insufficient proof of native process, pipe,
   timeout and descendant semantics.
8. The original C2 and iteration-002 host-cache, CRLF, container-cache,
   source-trust, admission, acquisition, reproducibility, schema/verifier,
   licensing, SBOM and offline-verification checks still pass without
   weakening any boundary.
9. Independent skeptical review confirms both private-boundary isolation and
   native process bounds before any remote-gate proposal. Author validation
   alone cannot accept iteration 003 or PSCAN-06.

## Forbidden scope and reserved gates

- No implementation in this approval session.
- No modification of PASS-OUTCOME-SPEC-001, PASS-SPEC-001, their source record,
  TRACEABILITY, DEC-002, existing manifest schemas, verifier/product behavior,
  scanner detection, rules, policy, allowlist, receipt, revocation or
  consuming-project behavior.
- No new runtime dependency, package, module, executable download or trust
  root. A need for one returns to the owner.
- No movement, deletion, recreation or reuse of either locked tag; no C2
  tooling-tag creation and no replacement product version.
- No Docker command, image pull, image execution, dependency download,
  scanner/toolchain run, build or external network action in this approval
  session.
- No weakened digest pin, mutable or short image reference, registry login,
  retry-on-unknown behavior, privileged container, host networking, added
  capability, `chmod 777` or `--userns=host`.
- No credential, signing key, paid capability, TruffleHog work, push, setting or
  ruleset change, workflow run, signing, attestation, draft, publication or
  other remote mutation.
- No deletion or movement of the preserved ignored `graphify-out/**` user
  material.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

After bounded local implementation and independent acceptance, a separate
exact owner decision is still required before a push, C2 tooling tag, remote
setting change or build-only workflow run. Signing, attestation, draft creation
and publication remain separately gated.

## Approval-session exclusions

This approval session introduces no workflow, image-admission, process-runner,
build, test, schema, verifier or product implementation. It performs no Docker
command, image pull, dependency download, scanner/toolchain execution, build,
remote read/write, tag operation, workflow run, signing, attestation, draft or
release action. It creates no credential, key, remote resource or paid
commitment. Its only durable change is this bounded authority bundle and
synchronized living task state.
