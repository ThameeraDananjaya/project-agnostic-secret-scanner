# PSCAN-06 Task Specification

## Title

Reproducible cross-platform build, SBOM/licence bundle, keyless-signing release
workflow, immutable-release plan, and offline verifier guidance.

## State

Previously activated by the recorded owner command on 2026-09-02 and claimed
alone on 2026-09-03 in a fresh implementation session from exact clean
activation commit `058ffcd446c6431b2e1afeed769d02c7b1f307f8` after the complete
reading-map and current preflight. Candidate
`a13c28fe7273bc8dc6545f97966a02889524eb4c` is locally implemented,
reproduced, validated and skeptically reviewed. The owner approved the exact
remote/signing gate on 2026-09-03. Repository creation and control read-back
succeeded, but authorized workflow run `33709197614` failed closed during
pinned dependency acquisition because the Linux runner could not write the
bind-mounted Go module cache. PSCAN-06 is not accepted or closed; no signing,
attestation, draft release or publication occurred. Exact evidence is in
`evidence/PSCAN-06/REMOTE-GATE-FAILURE-001.md`. No successor is selected or
activated. On 2026-09-03 the owner approved the bounded contract-preserving
Correction C1 recovery recorded in
`evidence/PSCAN-06/CORRECTION-C1-AUTHORITY.md`. This approval record does not
claim, implement or accept the correction. Correction C1 was subsequently
claimed alone from exact clean authority commit
`3fb1b0a55dc4f48dd35464c63c768f497efbc89b` and locally implemented as a
candidate `f24b832ebe5f6749aa0ab910e1b9be279065ebb1`, but independent review
rejected it after reproducing carriage-return corruption in Docker POSIX shell
payloads from a Windows CRLF checkout. Iteration 003 candidate
`5ca77226ed3996a8267caf02da366d0beb915c8d` fixed shell transport but was
independently rejected after a complete CRLF-checkout build consumed
checkout-transformed integrity bytes and failed the pinned-rule binding test.
Correction C1 iteration 004 repaired that defect as candidate
`a22579fd5af14473e1d49b5027f21ab591bbb589`, but independent review rejected
its remaining status-only trust gate after assume-unchanged driver tampering
produced an empty status while the modified driver executed. Iteration 005
replaces that gate with complete hostile-state index/path/mode/raw-byte
verification and exact committed workflow/harness/build entrypoints. Exact
candidate `3fb7592889820fa2739a4a53588e073689621809` rejects all twelve isolated
adversarial source states before output or untrusted driver action, passes two
complete builds from separate actual CRLF checkouts, and reproduces all 30
files byte-for-byte. Exact author evidence is in
`evidence/PSCAN-06/CORRECTION-C1-AUTHOR-VALIDATION-005.md`. Actual-Linux host
positive/wrong-owner proof and independent skeptical acceptance remain open.
PSCAN-06 remains the sole open task and remains unaccepted.

On 2026-09-04 the owner approved the exact Correction C1 iteration-005
build-only Linux proof gate. Remote `main` fast-forwarded to evidence commit
`3523e4409ebc53cc1931e3dcaf7d1eea74bb15ab`; exact tooling tag
`release-tooling-v1.0.0-c1` was created at candidate
`3fb7592889820fa2739a4a53588e073689621809` and protected against deletion or
update with no bypass. Workflow run `33829598255` passed immutable invocation,
credential-free checkout, locked product identity, exact entrypoint and all
twelve hostile-source checks, then failed closed before acquisition because
the fresh Ubuntu runner lacked the pinned Docker image required by the earlier
CRLF regression. The Linux cache proof and all later build steps were skipped;
the signing job was skipped and no artifact, deployment or release was
created. Exact evidence is in
`evidence/PSCAN-06/REMOTE-GATE-FAILURE-002.md`. Both tags remain immutable.
PSCAN-06 remains open and unaccepted; any correction, new tooling identity or
rerun requires a separate exact owner decision.

On 2026-09-04 the owner approved Correction C2 pinned-image bootstrap recovery,
recorded in `evidence/PSCAN-06/CORRECTION-C2-AUTHORITY.md`. The approval is a
local authority bundle only: it does not claim or implement C2 and authorizes no
push, tag, setting, workflow, signing, attestation, draft or publication action.
Correction C2 must begin in a fresh session from the committed authority bundle.
PSCAN-06 remains the sole open task and remains unaccepted.

Correction C2 candidate `e7faf0f81b3853e2090c75378bfbca568b52efad`
subsequently received partial author validation only. A skeptical inspection
then found that any nonzero Docker image-inspection result was classified as
image absence and could authorize the networked pull. On 2026-09-04 the owner
approved the bounded fail-closed image-admission correction as Correction C2
iteration 002. Its authority is recorded in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-002-AUTHORITY.md`. This authority
bundle does not claim or implement the iteration. A fresh session must claim
only iteration 002 from its exact committed authority. PSCAN-06 remains open
and unaccepted; no remote or successor action is authorized.

Correction C2 iteration-002 candidate
`52f7ee22ab722d7590e2d2c8326e14a0e9670462` subsequently received local author
validation only. Inspection then found that the production admission file
retains mutable test callbacks capable of substituting Docker and host proof
results, and that its declared Docker output and timeout limits are not enforced
during the complete native-process lifecycle. On 2026-09-05 the owner approved
the bounded private-admission and bounded-process correction as Correction C2
iteration 003. Its authority is recorded in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-003-AUTHORITY.md`. This authority
bundle does not claim or implement the iteration. A fresh session must claim
only iteration 003 from its exact committed authority. PSCAN-06 remains open
and unaccepted; no Docker, network, remote or successor action is authorized.

## Objective

After current read-only preflight and exact action-time owner approval, establish
the trusted build/release boundary for exact Windows amd64 and Linux amd64 assets,
complete supply-chain evidence and offline consumer verification. PSCAN-06 does
not publish `v1.0.0`; final acceptance/publication belongs to PSCAN-07.

## Preconditions and owner gates

- PSCAN-05 accepted, closed and clean at activation parent
  `30d27856bdfd404fc190501be854b78bc5147f1f`.
- Exact PSCAN-06 owner selection and activation recorded in
  `evidence/PSCAN-06/ACTIVATION.md`, followed by a fresh implementation session
  from that activation commit.
- Revalidate official Go, Gitleaks, Cosign/Sigstore, GitHub Actions, immutable
  releases, artifact attestations, rules/settings, plan and cost facts.
- Read-only preflight must confirm exact account/repository/workflow identity,
  availability, proposed settings, permissions and zero-cost posture.
- Owner separately approves each remote creation, push/settings/workflow action
  and any signing dry run. Activation alone is not approval.

## Activation-only bundle

This session may change and commit only:

```text
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/ACTIVATION.md
```

No implementation, dependency/action intake, scanner/toolchain execution, test
run, workflow execution, remote preflight mutation or successor work is
authorized in the activation session.

## Allowed implementation paths

```text
.github/workflows/**
.github/CODEOWNERS
build/**
licenses/**
contracts/release-manifest/**
contracts/global-revocation/**
internal/verify/**
tests/integration/supply-chain/**
tests/acceptance/supply-chain/**
docs/release/**
docs/security/**
docs/validation/**
docs/governance/RUNBOOK-REQUIREMENTS.md
docs/tasks/PSCAN-06.md
docs/tasks/TRACKER.md
THIRD_PARTY_NOTICES.md
SECURITY.md
evidence/PSCAN-06/**
```

## Correction C1 approved boundary

### Title

Linux cache ownership and dual-identity `v1.0.0` recovery.

### Purpose

Correct the proved Linux bind-mount acquisition failure without moving,
deleting or replacing the locked product tag and without weakening the release
contract. The correction keeps `v1.0.0` fixed at product-source commit
`a13c28fe7273bc8dc6545f97966a02889524eb4c` and introduces a separately bound,
immutable correction-tooling identity. The intended final release remains
signed public `v1.0.0`.

### Correction implementation paths

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

This list narrows the original PSCAN-06 implementation boundary for Correction
C1. A path allowed by the original task but absent here is forbidden to this
correction.

### Correction obligations

1. Preserve `--cap-drop ALL`, the read-only container-root controls and all
   other accepted sandbox controls.
2. On Linux, determine the host numeric UID/GID, run the pinned acquisition
   container with that identity, and give only the required cache/tmpfs paths
   compatible ownership and bounded modes. Do not use privileged mode,
   `chmod 777`, added capabilities or `--userns=host`.
3. Before any networked dependency acquisition, prove cache write, atomic
   rename, read and delete with a synthetic canary. A wrong owner, read-only
   cache or failed canary must stop before download and must create no success
   ledger.
4. Preserve the accepted Windows Docker Desktop behavior and require the same
   semantic canary. Prove the Linux path on an actual Linux host.
5. After acquisition, mount the completed cache read-only into two independent
   network-disabled builds and require byte-identical outputs.
6. Keep the product source identity fixed at locked tag `v1.0.0` and commit
   `a13c28fe7273bc8dc6545f97966a02889524eb4c`. Do not introduce a replacement
   product version or reinterpret a tooling commit as product source.
7. Use the proposed immutable correction-tooling tag
   `release-tooling-v1.0.0-c1` only after a separate remote gate. Bind that tag,
   its exact commit/tree, the recovery workflow path, workflow ref and workflow
   SHA as release-tooling identity distinct from product-source identity.
8. Add release-manifest schema major `2.0` so the two identity roles are
   explicit and non-interchangeable. Preserve all existing schemas. Unknown
   majors, role swaps and mutation of either identity must fail closed.
9. Bind Cosign verification to the exact approved repository, workflow, ref,
   SHA, trigger and issuer. Keep keyless signing inside the separately approved
   GitHub workflow; create no key.

### Correction acceptance

- Positive Linux cache ownership and canary evidence passes; wrong-owner and
  read-only-cache tests fail before acquisition and emit no ledger.
- Product-source and release-tooling commit, tree, ref, tag, workflow and role
  mutations reject; the two tag roles cannot be swapped.
- Manifest `2.0` validates exactly; old schemas cannot masquerade as `2.0`, and
  unknown majors reject.
- Two network-disabled builds from the completed read-only cache are
  byte-identical, and runner, Gitleaks engine and rules match the accepted
  product-source digests.
- Existing Go tests, `go vet`, Windows compilation, Linux execution, licensing,
  SBOM and offline release-verification checks pass without weakening any
  accepted boundary.
- No remote action occurs during local Correction C1 implementation or local
  acceptance. PSCAN-06 remains open after local success until a separate exact
  remote/signing correction gate is approved and successfully evidenced.

### Correction exclusions

- No modification of PASS-OUTCOME-SPEC-001, PASS-SPEC-001, their source record,
  TRACEABILITY, DEC-002 or any predecessor historical evidence.
- No movement, deletion, recreation or reuse of locked product tag `v1.0.0` for
  another commit; no `v1.0.1` substitution for the contracted first release.
- No scanner detection, rules, policy, allowlist, receipt, revocation,
  project-integration or consuming-project behavior change.
- No TruffleHog work, paid capability, credential, long-lived signing key,
  push, tag creation, workflow run, signing, attestation, draft/public release,
  GitHub settings change or other remote mutation.
- No selection, activation, claim or implementation of PSCAN-07, PSCAN-08 or
  any other successor.

## Correction C2 approved boundary

### Title

Pinned-image bootstrap recovery.

### Purpose

Resolve the fresh-runner ordering contradiction proved by workflow run
`33829598255`: the CRLF regression required the digest-pinned Docker image
before the only authorized acquisition step could obtain it. Admit only the
exact pinned image through a narrow bootstrap phase, prove its repository
digest before execution, and preserve every later offline, cache-canary,
sandbox, reproducibility and dual-identity boundary. An image-independent host
cache canary and CRLF-normalization proof remain before the networked pull.

### Correction implementation paths

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
`schema-2.0.json`, remain immutable.

### Correction obligations

1. Bootstrap only
   `docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`.
   Reject every mutable, shortened, ambiguous or substituted reference.
2. Before any networked pull, run an image-independent host-side cache canary
   and prove raw CRLF payload rejection plus LF normalization without Docker.
   Wrong-owner, read-only or failed host state stops before network access.
3. If the exact image is absent, allow one explicit networked pull before any
   image-dependent check, then independently prove the admitted repository
   digest before execution. Pull or inspection failure stops without a ledger
   or artifact.
4. The image bootstrap acquires no Go, Gitleaks, module, Cosign or other input
   and executes no image command. Every later Docker invocation uses the exact
   digest with `--pull=never`.
5. After image admission, run the container cache canary and normalized shell-
   parser proof with `--pull=never --network none`. Keep all dependency
   downloads after both canaries and keep both reproducibility builds network-
   disabled with the completed cache mounted read-only.
6. Preserve all accepted sandbox, UID/GID, hostile-source and exact-entrypoint
   controls.
7. Preserve locked product tag `v1.0.0` and locked C1 tooling tag
   `release-tooling-v1.0.0-c1`. Changed tooling must use proposed immutable tag
   `release-tooling-v1.0.0-c2` only after a later exact remote gate.
8. Preserve manifest schema `2.0`; add schema `2.1` for the C2 tooling identity.
   Product-source and release-tooling roles remain mandatory and
   non-interchangeable; omission, mutation, masquerade and unknown versions
   fail closed.
9. Keep keyless signing bound to the exact repository, owner, workflow, C2 ref,
   workflow SHA, trigger and issuer and behind a separate owner gate.

### Correction acceptance

- Positive host-side cache state passes before any pull; wrong-owner and read-
  only host cases stop before network access. Image-absent and pre-existing-
  image cases prove only the canonical digest;
  mutable tag, short-name, wrong-digest, ambiguous, pull-failure and inspection-
  mismatch cases reject before execution and emit no ledger or artifact.
- Image-independent CRLF rejection and normalization precede the pull; offline
  container parsing and cache-canary cases follow exact image admission and
  retain their fail-closed outcomes.
- Two network-disabled builds from the read-only completed cache reproduce all
  output bytes.
- Product-source and C2 tooling identity mutations and role swaps reject;
  manifest `2.1` validates exactly without changing historical schemas.
- Existing tests, hostile-source cases, `go vet`, Windows compilation, Linux
  execution, licensing, SBOM and offline verification all pass.
- Independent skeptical review proves the bounded fresh-runner behavior. Local
  success authorizes only a proposal for a separate build-only Linux proof gate.

### Correction exclusions

- No change to controlling or historical contracts, TRACEABILITY, DEC-002,
  existing manifest schemas or predecessor evidence.
- No movement, deletion, recreation or reuse of either locked tag; no
  replacement product version.
- No scanner, policy, allowlist, receipt, revocation, consuming-project,
  credential, paid capability or TruffleHog work.
- No push, C2 tooling tag, setting, workflow run, signing, attestation, draft,
  publication or other remote mutation.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

## Forbidden scope

- No modification of PASS-OUTCOME-SPEC-001, PASS-SPEC-001, their source record,
  TRACEABILITY, DEC-002 or PSCAN-01 through PSCAN-05 and PSCAN-09/10 historical
  evidence.
- No secrets, credentials, long-lived signing keys, project receipt keys,
  consuming-project integration, source, policy, allowlist, receipt,
  revocation, identity, finding, customer data, statistics, project-specific
  behavior, deployment gates or real candidate data.
- No `latest`, mutable refs, auto-update, pipe-to-shell, unpinned actions,
  incomplete asset sets, unverified identity, unresolved licence material or
  pass from missing, stale, conflicting or unsupported evidence.
- No TruffleHog assessment, legal work, download, material, integration,
  distribution or enablement.
- No paid capability, provider commitment, final release publication,
  deployment, production, legal/compliance or go-live action.
- No remote creation, push, settings/workflow enablement, signing dry run or
  other remote action without separate exact recorded owner approval after the
  required current read-only preflight.
- No selection, activation, claim or implementation of PSCAN-07, PSCAN-08 or
  any other successor.

## Deliverables

- Full-SHA-pinned minimal-permission trusted workflows and reproducible builds.
- Exact runner/engine/rules/schema assets for Windows amd64 and Linux amd64.
- Signed release-manifest design binding every digest, checksum, licence/notice,
  SBOM, compatibility, test summary, limitation and revocation location.
- Exact keyless identity policy for repository, workflow, release ref/tag and
  GitHub OIDC issuer; no private signing key.
- Draft/all-assets-together and immutable-release control plan based on current
  verified capability.
- Offline verifier and complete plain-English runbook covering acquisition,
  identity, bundle, digest, licence, SBOM, revocation, rollback and retirement.
- Current account/settings/cost evidence and explicit remaining publication gate.

## Acceptance

Two platform assets reproduce and verify; one-byte, identity, issuer, tag,
bundle, licence, SBOM and revocation mutations reject before scan; actions are
fully pinned and permissions minimal; no credential enters scan execution; no
unapproved remote/spend/publication occurred.

Stop after PSCAN-06 closeout with PSCAN-07 unselected.

## Correction C2 implementation status

Correction C2 was claimed alone in a fresh session from exact clean authority
commit `d4eca19e04862d660ca6ac0e9b64eec4fb06b61c`. Bounded local implementation
candidate `e7faf0f81b3853e2090c75378bfbca568b52efad`
separates host-only pre-image proofs, exact canonical digest image admission,
offline container proofs, networked dependency acquisition, and offline
reproducible builds. It adds release-manifest schema `2.1` and C2 verifier
policy without modifying schema `2.0` or either locked tag.

This is author implementation only. Independent skeptical acceptance, actual
fresh-runner execution and every pull/tag/push/workflow/signing/attestation/
draft/publication gate remain open. PSCAN-07 and PSCAN-08 remain unselected.

## Correction C2 iteration 002 approved boundary

The unaccepted candidate collapses every nonzero `docker image inspect` result
into image absence, so failures other than proved absence can reach the pull
branch when pull authority is enabled. Iteration 002 must distinguish exact
present, conclusively absent and untrusted/failed inspection states. Only
conclusive absence on an otherwise responsive engine may permit one exact
digest-pinned pull. Every daemon, permission, timeout, protocol, malformed,
ambiguous or unclassified result must reject with no pull or later action.

The iteration is limited to the recovery workflow; the three image-admission
scripts; the named release, validation and living-task documents; and
`evidence/PSCAN-06/**`. It must add deterministic command-boundary tests that
prove exact Docker arguments, call order and zero pulls for every non-absence
failure without requiring Docker or network access. It preserves all original
C2 ordering, sandbox, identity, schema, signature, locked-tag, remote and
successor gates. Exact scope, acceptance cases and exclusions are controlling
in `evidence/PSCAN-06/CORRECTION-C2-ITERATION-002-AUTHORITY.md`.

This approval supersedes candidate `e7faf0f` only as a readiness claim. The
candidate and its partial author evidence remain immutable history. Iteration
002 was claimed alone in fresh session
`01a06e1c-5c06-7763-958a-fdb623bae6aa` from exact clean authority commit
`281bea031bb6bbaf1be3059977074df2a88ecdf4` and locally implemented as a bounded
author candidate. Author validation is not acceptance. Actual Linux execution,
fresh-runner proof and independent skeptical review remain open; PSCAN-06 is
unaccepted and no successor or remote action is authorized.

## Correction C2 iteration 003 approved boundary

The unaccepted iteration-002 candidate keeps mutable script-scope test
callbacks inside the production admission file, so its production decision is
not structurally private from synthetic Docker, cache-canary and CRLF proof
substitution. Its native Docker runner also reads both redirected streams to
completion before checking their declared size, then can wait without a bound
after timeout. The recorded fake-result cases do not prove live byte limits,
bounded process-tree termination or pipe closure.

Iteration 003 must remove test substitution from the workflow-facing
production boundary and enforce the existing stdout, stderr and per-command
limits while the native process is running. Timeout, overflow, read/termination
failure, incomplete stream closure or surviving descendants must reject with
no pull or later action. Deterministic process fixtures must prove the native
Windows and actual-Linux lifecycle, including concurrent floods, hangs and
child/grandchild cleanup.

The iteration is limited to the recovery workflow; the three existing image-
admission scripts; the named release, validation and living-task documents;
and `evidence/PSCAN-06/**`. It preserves the iteration-002 three-state
admission model and every original C2 ordering, image, sandbox, identity,
schema, signature, locked-tag, remote and successor gate. Exact scope,
acceptance cases and exclusions are controlling in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-003-AUTHORITY.md`.

This approval supersedes candidate `52f7ee2` only as a readiness claim. The
candidate and its author evidence remain immutable history. Iteration 003 was
claimed alone in fresh session
`01a06e8c-6a86-74a2-819f-edcbc15fe6f8` from exact authority commit
`5762d1ef6a7708629bad5f5bb33eba21ca1cdcd4` in an isolated checkout with clean
tracked, untracked and ignored state. The bounded local implementation removes
the production callback seams, closes the workflow entrypoint, and adds live
byte, UTF-8, monotonic-time, pipe and process-tree tests. This is author work,
not acceptance; actual-Linux execution and independent skeptical review remain
open.
