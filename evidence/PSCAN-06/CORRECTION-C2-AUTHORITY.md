# PSCAN-06 Correction C2 Authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Title: Pinned-image bootstrap recovery
- Event: exact owner approval of a bounded pinned-image bootstrap recovery
- Owner command:
  `APPROVE PSCAN-06 CORRECTION C2 PINNED-IMAGE BOOTSTRAP RECOVERY`
- Date: `2026-09-04`
- Approval session: `01a06bc1-466a-7423-9689-edfd9a5a7672`
- Branch: `main`
- Required starting HEAD before this authority bundle:
  `a72a38d0d4a1322b064683520588155740ec512a`
- Correction authority commit: this authority-bundle commit
- Pre-change Git status: clean
- State after commit: owner-approved; unclaimed; not implemented; not accepted
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

This is a correction inside the already-activated PSCAN-06. It does not create,
select or activate another task. A fresh session must claim only Correction C2
from this committed authority bundle before implementation.

## Controlling and preserved authority

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Historical remote failure record:
  `evidence/PSCAN-06/REMOTE-GATE-FAILURE-002.md`
- Failure-record SHA-256:
  `62D2AA42674234F06F2A79AD59966089237A1A1C98A2EABE885B0B66369F7A89`
- Failed run: `33829598255`
- Locked product tag: `v1.0.0`
- Locked product-source commit:
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Locked Correction C1 tooling tag: `release-tooling-v1.0.0-c1`
- Locked Correction C1 tooling commit:
  `3fb7592889820fa2739a4a53588e073689621809`
- Pinned build-image identity:
  `docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`

The failed run, both locked tags, all Correction C1 records, the original
activation and implementation records, earlier task evidence and predecessor
contracts remain immutable facts. Correction C2 references them; it does not
rewrite or rehabilitate them.

## Approved objective

Resolve the fresh-runner ordering contradiction proved by run `33829598255`:
the CRLF regression required the digest-pinned Docker image before the only
step authorized to acquire it. Add an image-independent host cache canary and
CRLF-normalization proof before a narrowly separated bootstrap phase that may
acquire only the exact pinned image. Prove the image identity, then run the
container canary and shell-parser proof offline. Preserve every later sandbox,
dual-identity, reproducibility and fail-closed boundary.

## Allowed implementation paths

Correction C2 may change only:

```text
.github/workflows/release-recovery-v1.0.0.yml
build/release/**
contracts/release-manifest/schema-2.1.json
internal/verify/**
tests/integration/supply-chain/**
tests/acceptance/supply-chain/**
docs/decisions/DEC-003-DUAL-IDENTITY-RELEASE-RECOVERY.md
docs/release/**
docs/validation/**
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/**
README.md
```

This list supersedes the Correction C1 path list for Correction C2 only. Every
absent path is forbidden. Existing release-manifest schemas, including
`schema-2.0.json`, are preserved byte-for-byte.

## Mandatory correction design

### Exact image-bootstrap boundary

1. Use only the canonical exact image reference recorded above. No tag-only,
   short-name, mutable, `latest`, platform-ambiguous or substitute image may be
   admitted.
2. Before any networked image pull, run an image-independent host-side
   write/atomic-rename/read/delete canary on the exact cache paths using the
   invoking host identity and bounded modes. Wrong-owner, read-only or failed
   host state stops before the pull and emits no success ledger.
3. Before the first image-dependent check, inspect local image state. If the
   exact digest is absent, one explicit network-enabled bootstrap step may pull
   only that canonical digest reference.
4. After any pull, independently inspect the local image and prove that the
   admitted repository digest equals the approved canonical reference before
   the image may execute. Missing, ambiguous or mismatched identity fails
   closed.
5. The bootstrap step downloads no Go archive, Gitleaks source, module, Cosign
   binary or other input, runs no untrusted image command, and emits no
   acquisition-success ledger or build artifact.
6. Before the image pull, separately prove from an actual CRLF checkout that
   every covered raw POSIX shell payload is rejected and its normalized form
   contains no carriage return. This proof must not require Docker or execute a
   payload.
7. Pull or inspection failure stops before the container canary, offline shell
   parsing, dependency acquisition, build or artifact transfer. Partial state
   is never evidence of success.
8. After image admission, run the container cache canary and normalized shell-
   parser proof with the exact digest, `--pull=never` and `--network none`.
   Wrong-owner and read-only container cases fail before dependency download
   and emit no success ledger.
9. All Go, Gitleaks and module downloads remain after both successful canaries
   and inside the existing bounded acquisition phase.
10. Preserve the read-only container root, dropped capabilities,
   `no-new-privileges`, bounded resources, exact host UID/GID mapping and all
   other accepted sandbox controls. Both reproducibility builds remain
   network-disabled and use the completed cache read-only.

### Immutable recovery identity

1. Product source remains locked tag `v1.0.0`, commit
   `a13c28fe7273bc8dc6545f97966a02889524eb4c` and its recorded tree.
2. The locked `release-tooling-v1.0.0-c1` tag remains historical and must not be
   moved, deleted, recreated or used to identify changed C2 tooling.
3. Proposed C2 tooling tag `release-tooling-v1.0.0-c2` may be created only under
   a later exact remote gate at an independently accepted C2 tooling commit.
4. Preserve release-manifest schema `2.0`. Add schema `2.1` to bind the C2
   tooling tag, exact commit/tree, recovery workflow path, workflow ref,
   workflow SHA and trigger while retaining separate, non-interchangeable
   product-source and release-tooling roles.
5. Version omission, old-schema masquerade, unknown versions, role swaps and
   mutation of either identity fail closed. Historical valid `2.0` evidence
   remains parseable only under its exact C1 identity policy.
6. Cosign verification remains bound to the exact approved repository, numeric
   owner, workflow, C2 ref, workflow SHA, trigger and issuer. Signing remains
   keyless, workflow-scoped and separately owner-gated.

## Required local checks and success criteria

1. Positive host-side cache state passes the image-independent canary before
   any pull; wrong-owner and read-only host cases stop before network access and
   emit no ledger or artifact.
2. A fresh-runner image-absent case admits only the canonical digest-pinned
   image and then passes exact post-pull identity inspection.
3. Mutable tag, short-name, wrong digest, ambiguous identity, pull failure and
   post-pull inspection mismatch cases reject before image execution and emit
   no ledger or build artifact.
4. A pre-existing exact image follows the same identity proof without a
   network pull; an unproved cached image cannot pass.
5. The image-independent CRLF proof rejects raw payloads and proves LF
   normalization before the pull; the admitted image then parses normalized
   payloads with `--pull=never --network none`.
6. The container canary passes only after exact image admission; wrong-owner
   and read-only-cache cases stop before dependency acquisition and emit no
   success ledger.
7. The completed cache is mounted read-only into two network-disabled builds
   whose complete outputs are byte-identical.
8. Product-source and C2 release-tooling commit, tree, ref, tag, workflow and
   role mutations reject; the identities cannot be swapped.
9. Manifest `2.1` validates exactly; schemas `1.0`, `1.1` and `2.0` cannot
   masquerade as `2.1`, and unknown versions reject.
10. Existing Go tests, `go vet`, Windows compilation, Linux execution, hostile
   source-state, licensing, SBOM and offline release-verification checks pass
   without weakening an accepted boundary.
11. Independent skeptical review examines the bounded diff and proves the
    fresh-runner bootstrap behavior without relying on author summaries or a
    passing subset.

Local success does not close PSCAN-06. It establishes eligibility to propose a
new exact build-only Linux proof gate for the accepted C2 tooling identity.

## Forbidden scope and reserved gates

- No modification of PASS-OUTCOME-SPEC-001, PASS-SPEC-001, their source record,
  TRACEABILITY, DEC-002, existing release-manifest schemas or predecessor
  historical evidence.
- No movement, deletion, recreation or reuse of locked tag `v1.0.0` or
  `release-tooling-v1.0.0-c1`; no replacement product version.
- No scanner detection/rule, policy, allowlist, receipt, revocation, project
  integration, consuming-project behavior or real-project-data change.
- No broad registry login, mutable image admission, unpinned network input,
  privileged container, added capability, `chmod 777` or `--userns=host`.
- No TruffleHog assessment, legal work, download, integration, distribution or
  enablement.
- No credential, long-lived signing key, paid capability or spending.
- No push, C2 tooling-tag creation, ruleset or setting change, workflow run,
  signing, attestation, draft/public release, publication or other remote
  mutation.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

After complete local acceptance, a separate exact owner decision is required
before pushing a C2 correction commit, creating and protecting the C2 tooling
tag, changing any remote setting, or running even a build-only workflow. Any
signing, attestation or draft creation requires a later separately bounded
gate. Final publication remains separately owner-gated.

## Approval-session exclusions

This approval session introduces no product code, schema, verifier, build,
workflow or test implementation. It performs no image pull, dependency
download, scanner/toolchain execution, build, remote read/write, tag operation,
workflow run, signing, attestation, draft or release action. It creates no
credential, key, remote resource or paid commitment. Its only durable change
is this bounded authority bundle and synchronized living task state.
