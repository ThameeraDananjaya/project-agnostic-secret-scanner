# PSCAN-06 Correction C2 iteration 002 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Iteration: `002`
- Title: Fail-closed image admission correction
- Event: exact owner approval of the bounded correction described below
- Owner command:
  `APPROVE PSCAN-06 C2 ITERATION-002 FAIL-CLOSED IMAGE ADMISSION CORRECTION`
- Date: `2026-09-04`
- Approval session: `01a06dfe-53c1-7c12-9540-d6aa086b88d9`
- Branch: `main`
- Required starting HEAD before this authority bundle:
  `636c2c037686b807967afa1c80b3b9ffd8134950`
- Correction authority commit: this authority-bundle commit
- Pre-change Git status: clean
- State after commit: owner-approved; unclaimed; not implemented; not accepted
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

This is a bounded iteration inside the already-activated PSCAN-06 Correction
C2. It does not create, select or activate another task. A fresh session must
claim only Correction C2 iteration 002 from this committed authority bundle
before implementation.

## Controlling and preserved authority

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Original Correction C2 authority:
  `evidence/PSCAN-06/CORRECTION-C2-AUTHORITY.md`
- Original authority SHA-256:
  `3CE8E20151993A8FC77890E51C5BF1557F68040CAFF2C3DB1A46237792A75D63`
- Unaccepted Correction C2 candidate:
  `e7faf0f81b3853e2090c75378bfbca568b52efad`
- Candidate author-validation record:
  `evidence/PSCAN-06/CORRECTION-C2-AUTHOR-VALIDATION.md`
- Author-validation SHA-256:
  `FA4D3C9E7F767575916C885449399943B7F809070F5A1E58268FFA85693AF637`
- Locked product tag: `v1.0.0`
- Locked product-source commit:
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Locked Correction C1 tooling tag: `release-tooling-v1.0.0-c1`
- Locked Correction C1 tooling commit:
  `3fb7592889820fa2739a4a53588e073689621809`
- Proposed Correction C2 tooling tag: `release-tooling-v1.0.0-c2`
- Pinned build-image identity:
  `docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`

The original C2 authority, candidate, partial author-validation record, both
locked tags, existing schemas and all predecessor evidence remain immutable
facts. The partial author validation is not acceptance and is superseded as a
readiness claim by this newly identified fail-closed gap.

## Demonstrated gap

Candidate `e7faf0f` treats every nonzero
`docker image inspect` result as the single condition `Pinned build image is
absent`. The admission orchestrator then permits the networked pull whenever
that message is observed and `AllowImagePull` is enabled. A daemon failure,
permission denial, timeout, protocol failure, invalid invocation or other
unclassified inspection failure can therefore enter the pull branch without
positive proof that the exact image is merely absent.

The existing image-admission unit record exercises canonical repository-digest
classification directly. It does not exercise the Docker-command failure
boundary or prove that all non-absence inspection failures produce zero pull
attempts. Missing proof at this network admission boundary fails closed.

## Approved objective

Make the image-admission state machine explicit and fail closed. A networked
pull may occur only after trustworthy, structured evidence proves that the
exact canonical digest-pinned image is absent while the Docker engine is
otherwise reachable and responsive. Every unavailable, malformed, ambiguous,
unsupported or unclassified inspection condition must stop with no pull, no
image execution, no success ledger, no dependency acquisition and no build
artifact.

## Allowed implementation paths

Correction C2 iteration 002 may change only:

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

This list supersedes the broader Correction C2 path list for iteration 002
only. Every absent path is forbidden. A need to change another path stops the
iteration and returns to the owner with exact evidence.

## Mandatory correction design

1. Model at least three distinct pre-pull outcomes: exact image present and
   proved; exact image conclusively absent; and inspection untrusted or failed.
   Only the conclusively absent outcome may reach the pull branch.
2. Do not classify absence from a generic nonzero exit code, exception text,
   free-form or localized stderr, ambient `$LASTEXITCODE`, missing output, or
   malformed output. Engine reachability and the absence conclusion must be
   established with deterministic bounded command results and structured data.
3. Capture each Docker invocation's exit status, stdout and stderr at its own
   call boundary. Reject mixed, truncated, extra, null, scalar, duplicate,
   malformed or otherwise unexpected identity evidence.
4. Preserve exact case-sensitive caller input
   `docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`.
   No tag, short caller reference, substitute repository or digest may be
   admitted.
5. After conclusive absence, permit at most one explicit pull of that exact
   reference. Pull failure or any post-pull inspection problem stops without a
   retry, execution, ledger, dependency acquisition or artifact.
6. A pre-existing exact image remains pull-free and must pass the same exact
   repository-digest proof required after a pull.
7. Add a deterministic fake-engine or equivalent command-boundary harness that
   records invocations without network or Docker and proves pull count and
   ordering for positive, absent and adversarial outcomes.
8. Preserve the host cache and CRLF proofs before admission; preserve all
   post-admission Docker commands at the exact digest with `--pull=never` and
   `--network none`; preserve acquisition ordering and both network-disabled
   reproducibility builds.
9. Preserve manifest schema `2.1`, all earlier schemas, product/tooling role
   separation, signature identity policy, locked tags and every original C2
   sandbox and owner gate. This iteration does not change release identity or
   scanner behavior.

## Required local checks and success criteria

1. Exact pre-existing image plus exactly one canonical repository digest
   admits with zero pull calls.
2. Conclusive absence on a responsive engine permits exactly one pull of the
   exact canonical digest, followed by a fresh independent identity inspection.
3. Generic inspect failure, daemon unavailable, permission denial, timeout,
   invalid invocation, nonzero exit with deceptive absence text, and any
   unclassified failure reject with zero pull calls.
4. Empty, null, scalar, malformed, mixed, duplicate, wrong-repository,
   wrong-digest and ambiguous structured results reject with zero image
   execution and no success ledger.
5. Pull failure and every post-pull inspection or identity failure reject with
   no retry, image execution, dependency acquisition, ledger or artifact.
6. The harness proves exact Docker command count, arguments and order and
   rejects any unexpected command. Its negative cases require no Docker daemon,
   network or mutable external state.
7. The original C2 host-cache, CRLF, container-cache, source-trust, acquisition,
   reproducibility, schema/verifier, licensing, SBOM and offline verification
   checks still pass without weakening any boundary.
8. Actual Linux execution and independent skeptical review confirm the bounded
   correction before any remote-gate proposal. Author validation alone cannot
   accept this iteration or PSCAN-06.

## Forbidden scope and reserved gates

- No implementation in this approval session.
- No modification of PASS-OUTCOME-SPEC-001, PASS-SPEC-001, their source record,
  TRACEABILITY, DEC-002, existing manifest schemas, verifier/product behavior,
  scanner detection, rules, policy, allowlist, receipt, revocation or
  consuming-project behavior.
- No movement, deletion, recreation or reuse of either locked tag; no C2
  tooling-tag creation and no replacement product version.
- No image pull, image execution, dependency download, scanner/toolchain run,
  build or external network action in this approval session.
- No weakened digest pin, mutable or short image reference, registry login,
  retry-on-unknown behavior, privileged container, host networking, added
  capability, `chmod 777` or `--userns=host`.
- No credential, signing key, paid capability, TruffleHog work, push, setting or
  ruleset change, workflow run, signing, attestation, draft, publication or
  other remote mutation.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

After bounded local implementation and independent acceptance, a separate
exact owner decision is still required before a push, C2 tooling tag, remote
setting change or build-only workflow run. Signing, attestation, draft creation
and publication remain separately gated.

## Approval-session exclusions

This approval session introduces no workflow, image-admission, build, test,
schema, verifier or product implementation. It performs no Docker command,
image pull, dependency download, scanner/toolchain execution, build, remote
read/write, tag operation, workflow run, signing, attestation, draft or release
action. It creates no credential, key, remote resource or paid commitment. Its
only durable change is this bounded authority bundle and synchronized living
task state.
