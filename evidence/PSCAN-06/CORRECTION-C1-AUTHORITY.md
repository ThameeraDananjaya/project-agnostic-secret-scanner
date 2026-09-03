# PSCAN-06 Correction C1 Authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C1`
- Title: Linux cache ownership and dual-identity `v1.0.0` recovery
- Event: exact owner approval of a bounded contract-preserving recovery
- Owner command:
  `APPROVE PSCAN-06 CORRECTION C1 CONTRACT-PRESERVING RECOVERY`
- Date: `2026-09-03`
- Approval session: `01a053a9-298b-75e1-9d49-7d2261bbb951`
- Branch: `main`
- Required starting HEAD before this authority bundle:
  `a5f5b34693d6ad83f841dabfab175998dc6a7be3`
- Correction authority commit: this authority-bundle commit
- Pre-change Git status: clean
- State after commit: owner-approved; unclaimed; not implemented; not accepted
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

This is a correction iteration inside the already-activated PSCAN-06. It does
not create, select or activate another task. A fresh session must claim only
Correction C1 from this committed authority bundle before implementation.

## Controlling and preserved authority

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Historical remote failure record:
  `evidence/PSCAN-06/REMOTE-GATE-FAILURE-001.md`
- Failure-record SHA-256:
  `52D45468BEF600012DBBD0D7A94D648D182A589A8FC63333D9AF85FED5C8F145`
- Failed run: `33709197614`
- Locked product tag: `v1.0.0`
- Locked product-source commit:
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`

The failed run, its evidence, the original activation and implementation
records, all earlier task evidence and all predecessor contracts remain
immutable facts. Correction C1 references them; it does not rewrite them.

## Approved objective

Correct the Linux bind-mounted Go module-cache ownership failure and add the
regression proof that would have caught it. Recover the contracted signed
public `v1.0.0` release without moving, deleting or replacing its locked product
tag by separating immutable product-source identity from immutable correction-
tooling/workflow identity.

## Allowed implementation paths

Correction C1 may change only:

```text
.github/workflows/release-recovery-v1.0.0.yml
build/release/**
contracts/release-manifest/schema-2.0.json
internal/verify/**
tests/integration/supply-chain/**
tests/acceptance/supply-chain/**
docs/decisions/DEC-003-DUAL-IDENTITY-RELEASE-RECOVERY.md
docs/release/**
docs/validation/**
docs/tasks/PSCAN-06.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/**
README.md
```

This narrower list supersedes the broader original PSCAN-06 path list for this
correction only. Every absent path is forbidden.

## Mandatory correction design

### Linux acquisition boundary

- Preserve the accepted container sandbox, including `--cap-drop ALL` and a
  read-only container root.
- On Linux, map the host numeric UID/GID into the pinned acquisition container
  and provide compatible, bounded ownership/mode only for required cache and
  tmpfs paths.
- Perform a synthetic write, atomic rename, read and delete cache canary before
  networked dependency acquisition. Any failure stops before download and
  produces no success ledger.
- Preserve Windows Docker Desktop behavior with the same semantic canary.
- Prove actual Linux positive, wrong-owner and read-only-cache cases.
- Do not use privileged mode, `chmod 777`, added capabilities or
  `--userns=host`.
- Mount the completed cache read-only into two independent network-disabled
  builds and require byte-identical results.

### Dual release identity

- Product source remains locked tag `v1.0.0` at commit
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`.
- Proposed correction-tooling tag: `release-tooling-v1.0.0-c1`. It may be
  created only under a later exact remote gate at the accepted correction
  tooling commit.
- Release-manifest schema major `2.0` must bind product-source tag, commit and
  tree separately from release-tooling tag, commit, tree, recovery-workflow
  path, workflow ref and workflow SHA.
- The two identities are mandatory and non-interchangeable. Mutation, omission,
  ambiguity or role swapping is rejection.
- Preserve existing manifest schemas. An older schema cannot masquerade as
  `2.0`; every unknown major fails closed.
- Cosign verification must bind the exact repository, workflow, ref, SHA,
  trigger and issuer. Signing remains keyless and workflow-scoped.
- The intended final release remains signed public `v1.0.0`; no replacement
  version is authorized.

## Required local checks and success criteria

1. Positive Linux cache ownership and canary succeeds.
2. Wrong-owner and read-only-cache cases stop before acquisition and emit no
   success ledger.
3. Product-source and release-tooling commit, tree, ref, tag, workflow and role
   mutation tests reject; tag roles cannot be swapped.
4. Schema `2.0` exact-validation, old-schema masquerade and unknown-major tests
   behave fail closed.
5. Two builds use the completed read-only cache with network disabled and
   produce byte-identical Windows amd64 and Linux amd64 assets.
6. Runner, pinned Gitleaks engine and rules match the accepted product-source
   digests.
7. Existing Go tests, `go vet`, Windows compilation, Linux execution, licence,
   SBOM and offline release-verification checks pass.
8. Required evidence records exact commands, pins, expected/actual outcomes,
   limitations, hashes, changed paths and clean post-commit state.
9. Independent skeptical review confirms the correction without relying on an
   implementation summary or a passing subset.

Local success does not close PSCAN-06. It establishes eligibility to propose a
new exact remote/signing correction gate.

## Forbidden scope and reserved gates

- No modification of PASS-OUTCOME-SPEC-001, PASS-SPEC-001, their source record,
  TRACEABILITY, DEC-002, existing release-manifest schemas or predecessor
  historical evidence.
- No movement, deletion, recreation or reuse of locked product tag `v1.0.0` for
  any other commit; no `v1.0.1` substitution.
- No scanner detection/rule, policy, allowlist, receipt, revocation, project
  integration, consuming-project behavior or real-project-data change.
- No TruffleHog assessment, legal work, download, integration, distribution or
  enablement.
- No credential, long-lived signing key, paid capability or spending.
- No push, tooling-tag creation, GitHub setting change, workflow run, signing,
  attestation, draft/public release, publication or other remote mutation.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

After complete local acceptance, a separate exact owner decision is required
before pushing the correction commit, creating or locking the tooling tag,
changing settings, running the recovery workflow, signing, attesting or
creating a draft release. Final publication remains separately owner-gated.

## Approval-session exclusions

This approval session introduces no product code, schema, verifier, build,
workflow, test, signing, release or remote change. It executes no scanner,
toolchain, build, test or workflow. It creates no credential, key, remote
resource or paid commitment. Its only durable change is this bounded authority
bundle and synchronized living task state.
