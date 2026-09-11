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

Correction C2 iteration-003 candidate
`dfbe897e9e47632ee5ca9437650bd62eafbaf341` subsequently received exact-commit
Windows author validation only. Inspection then found that workflow-reachable
acquisition, cache, CRLF and build scripts still use ambient `docker` execution
outside the closed admission runner, and that root `WaitForExit` plus redirected
pipe closure does not prove complete descendant termination. On 2026-09-05 the
owner approved the bounded closed-Docker-execution and process-tree correction
as Correction C2 iteration 004. Its authority is recorded in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-004-AUTHORITY.md`. This authority
bundle does not claim or implement the iteration. A genuinely fresh session
must claim only iteration 004 from its exact committed authority. PSCAN-06
remains open and unaccepted; no Docker, network, remote or successor action is
authorized.

Correction C2 iteration-004 was claimed alone in a fresh session and locally
implemented/refined through candidate
`9583aa3d18310c2e9275c665f69eb7e5b4fb82a4`; governance and author-evidence
HEAD `fad4c3e12c655944cf8dfc3c46d622666e777fb7` retained only partial author
validation. Independent review then found that a preloaded or stale CLR type
named `PscanNativeBoundary` can replace the exact committed Windows process
boundary and fabricate accepted lifecycle evidence. It also found the tracker
summary incorrectly describing iteration 004 as unclaimed and unimplemented.
On 2026-09-05 the owner approved the bounded ambient-type-isolation and
tracker-reconciliation correction as Correction C2 iteration 005. Its authority
is recorded in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-005-AUTHORITY.md`. This authority
bundle does not claim or implement the iteration. A genuinely fresh session
must claim only iteration 005 from its exact committed authority. PSCAN-06
remains open and unaccepted; no Docker, network, remote or successor action is
authorized.

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
byte, UTF-8, monotonic-time, pipe and process-tree tests. Exact-commit Windows
author validation passed source trust, hostile-source rejection, the host-only
CRLF proof and the native no-Docker process matrix. This remains author work,
not acceptance; actual-Linux execution, genuine Docker proof and independent
skeptical review remain open.

## Correction C2 iteration 004 approved boundary

The unaccepted iteration-003 candidate closes the bootstrap admission
entrypoint but does not close all later Docker execution. `acquire.ps1`,
`cache-canary.ps1`, `test-crlf-shell-payloads.ps1` and `build.ps1` still use
ambient PowerShell `docker` command resolution, and `image-admission.ps1`
retains a duplicated native runner. Those paths can bypass the fixed executable
identity, environment, stream and lifecycle boundary proved for admission.

The iteration-003 cleanup also treats bounded root exit and pipe closure as
process-tree completion. It retains no operating-system containment object or
complete membership proof. Child or grandchild creation can race with a tree
snapshot or root exit; post-hoc fixture PID checks do not prove complete
production containment.

Iteration 004 must route every workflow-reachable Docker operation through one
committed closed operation boundary with stable executable binding, internally
constructed operation-specific arguments, a cleared/private environment and
the existing output and monotonic-time limits. It must establish native process
containment before child execution and prove empty membership after normal exit
or bounded termination. Any identity, containment, termination or membership
uncertainty is terminal and permits no pull or later action.

The iteration is limited to the recovery workflow; the named admission,
acquisition, cache, CRLF, build and validation scripts; at most one new closed
Docker entrypoint and one no-Docker adversarial harness; the named release,
validation and living-task documents; and `evidence/PSCAN-06/**`. It preserves
the iteration-002 three-state admission model and every original C2 ordering,
image, sandbox, identity, schema, signature, locked-tag, remote and successor
gate. Exact scope, acceptance cases and exclusions are controlling in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-004-AUTHORITY.md`.

This approval supersedes candidate `dfbe897` only as a readiness claim. The
candidate and its author evidence remain immutable history. Iteration 004 was
claimed alone in fresh session `01a07160-ef63-7e80-a441-b2ce360cc23d` from
exact authority commit `30d849c11a15a7ce35f182140d07dbe66350c135` in a
separate clean LF checkout. The bounded implementation routes every release
Docker operation through one typed entrypoint, binds later operations to the
admitted executable digest, and introduces OS-backed containment plus empty-
membership proof. Candidate `6f791646413bfde52a7f034f6219d92f6fb44c03`
passed Windows author validation but a post-integration audit found a Linux
between-sample detached-member race. The append-only refinement record rejects
that design as final readiness evidence. Refined candidate
`9583aa3d18310c2e9275c665f69eb7e5b4fb82a4` uses a stopped, recorded PID-
namespace init before Docker execution so detached descendants cannot leave the
kernel containment boundary. Exact-commit Windows, host-only CRLF, admission
and hostile-source author checks pass. Actual Linux, genuine Docker, complete
builds and independent review remain open. This is author evidence, not
acceptance. No Docker, network, remote, signing, publication or successor
action occurred.

## Correction C2 iteration 005 approved boundary

The refined iteration-004 candidate closes literal ambient Docker command
resolution and adds OS-backed process containment, but its Windows native
boundary is compiled only when the PowerShell process does not already contain
a type named `PscanNativeBoundary`. Production then invokes and trusts that
type. A compatible hostile preload or stale type from an earlier script version
can therefore replace the committed implementation and fabricate successful
exit, stream and containment evidence. The current harness repeats the same
conditional reuse, so green tests do not close this substitution path.

Iteration 005 must treat any pre-existing expected native type as terminal
ambient state and fail before invoking it or Docker. Clean and hostile cases
must run in separate non-profile processes; a compatible fake and a stale-type
case must prove zero fake calls, zero Docker calls and no accepted evidence.
The exact committed source must be the only implementation behind a successful
clean result. The iteration must preserve the complete iteration-004 Windows
job-object and Linux PID-namespace designs, the nine typed Docker operations,
all executable/image/environment/stream/time limits and every earlier C2 gate.

The iteration also reconciles the living tracker: iteration 004 was claimed and
locally implemented/refined, but independent review rejected its readiness;
actual Linux, genuine Docker, complete builds and independent acceptance remain
open. Exact scope, checks and exclusions are controlling in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-005-AUTHORITY.md`.

This approval supersedes candidate `9583aa3` and governance/evidence HEAD
`fad4c3e` only as readiness claims. Their commits and evidence remain immutable
history. Iteration 005 is owner-approved but unclaimed, unimplemented and
unaccepted. It must begin in a genuinely fresh session from the exact committed
authority bundle and an exact clean source-trust-compliant checkout. No Docker,
network, remote, signing, publication or successor action is authorized.

Fresh implementation session `01a071ef-5eeb-7a91-8d82-1feb22a9c2f6` claimed
only iteration 005 from exact authority commit
`13496a35eab71482f2e908bd0aac4b563941682d` in a new local, no-network,
LF-preserving isolated clone after exact every-byte source trust passed. The
bounded implementation rejects an existing native type before compilation or
execution, validates the exact newly compiled type identities, and moves the
clean, compatible-hostile and stale-hostile proofs into separate non-profile
processes. Exact candidate
`050f5862841a1fcf764ffc75d98a742968ac301e` passed every required local
no-Docker author check on Windows. This is implementation and author evidence
only, not independent acceptance. All actual-Linux, Docker, build, remote,
signing, publication and successor gates remain open.

Independent reviewer task `01a07210-d00a-7e21-8816-483c53027c21`, completed
turn `01a07210-d2af-7ca1-b40f-0a9ce6bde7fa`, then reviewed exact evidence-
bearing candidate `faef8435322c9096df09b56969662411f17356ea`, tree
`3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`, in a separate local no-network
clone. It reported no blocking findings and terminal verdict
`LOCAL_ACCEPTANCE_PASS_REMOTE_PROOF_OPEN`. The complete local parse, exact-
source, 12-case hostile-source, actual-CRLF host-only, ambient/stale preload,
Windows native-boundary, nine-operation inventory, image-admission and process-
cleanup evidence passed. The review's cleanup was blocked before execution by
the host guard, no alternate destructive method was used, and the three literal
remaining paths are recorded in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-005-ACCEPTANCE.md`.

Iteration 005 is therefore independently accepted locally for its bounded
correction objective. PSCAN-06 remains open and unaccepted overall: actual
Linux, genuine Docker/image/container execution, dependency acquisition,
complete builds, byte comparison and remote proof remain open and separately
owner-gated. This result does not authorize the remote proof gate. PSCAN-07
remains proposed and unselected; PSCAN-08 remains inactive and ineligible.

On 2026-09-05 the owner approved the exact Correction C2 build-only actual-
Linux and genuine-Docker proof gate for independently accepted tooling
candidate `faef8435322c9096df09b56969662411f17356ea`, tree
`3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`, from evidence HEAD
`b6d341284cf63baa0502fcba319ea8cca3c7eb3a`. Exact authority is recorded in
`evidence/PSCAN-06/CORRECTION-C2-BUILD-ONLY-LINUX-DOCKER-PROOF-GATE-AUTHORITY.md`.
The gate is owner-approved but unexecuted. This approval session performs no
live remote preflight, push, C2 tooling tag, setting change, workflow run,
Linux/Docker execution, network acquisition or artifact action. A genuinely
fresh session must fail closed through the complete current preflight before
the exact fast-forward evidence push, immutable C2 tooling tag and single
build-only workflow dispatch. Maximum spend is USD 0; the sole permitted remote
output is one uncompressed unsigned workflow artifact retained for one day.
The C2 signing variable must remain absent, the signing job must remain skipped,
and signing, attestation, draft creation, publication, rerun and successor work
are forbidden. PSCAN-06 remains open and unaccepted overall.

## Correction C2 build-only proof-gate pre-dispatch failure and Recovery R1

Fresh Daybreak Blue `xhigh` execution session
`01a072e4-50f8-7842-8c7d-28bb034a773f` began from exact gate-authority commit
`268e8f339318a141d059115c9da0d19a01448aaa` and tree
`e1667892d0c05d9a8db8bac65e580fff54a876ea`, but the Codex worktree checkout
was detached rather than on branch `main`. Mandatory preflight item 1 therefore
failed as a whole. The session stopped before remote preflight, network access,
push, tag, ruleset, setting change, workflow dispatch, Docker or artifact work.
Its immutable record is
`evidence/PSCAN-06/CORRECTION-C2-BUILD-ONLY-LINUX-DOCKER-PROOF-GATE-FAILURE-001.md`,
SHA-256
`7EB1CB03E85A933C1BCFEB435FA85556170AE21F1112CFFA55041C2A733E6A5D`.

On 2026-09-06 the owner approved bounded main-branch execution Recovery R1.
The recovery changes only the launch context: a genuinely fresh task must use
Daybreak Blue `gpt-daybreak-blue-latest` at `xhigh` and run directly in the
saved project on exact branch `main` at the committed Recovery R1 authority
bundle. Detached worktrees, branch repair and execution from another root are
forbidden. Every original identity, preflight, zero-spend, one-day unsigned
artifact, no-signing, no-attestation, no-draft, no-publication and no-successor
condition remains mandatory. Because the failed session made no remote request
and no workflow dispatch, Recovery R1 permits at most one first dispatch after
all current preflight facts pass. Exact authority is recorded in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R1-AUTHORITY.md`.

This approval session imports evidence and records authority only. It performs
no remote preflight or mutation and does not execute the recovered gate.
PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains proposed and
unselected; PSCAN-08 remains inactive and ineligible.

Recovery R1 was subsequently attempted in genuinely fresh saved-project
session `01a07376-04a7-7aa1-8421-7d4e82a13465`. It reached the exact required
root, branch, authority commit and tree but failed closed at mandatory preflight
item 1 because `build/gitleaks/build.ps1` and
`build/gitleaks/collect-licenses.ps1` contained noncanonical mixed line endings.
It stopped before any remote request, workflow dispatch, Docker, artifact,
signing or publication action. The append-only failure is recorded in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R1-FAILURE-001.md`.

Under later owner direction, those two working-tree projections were restored
from their exact committed Git objects and now satisfy the canonical
LF-to-CRLF projection with zero bare LF or CR bytes. No committed source or
index identity changed. Exact local repair evidence is recorded in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-SOURCE-TRUST-REPAIR-001.md`.
This repair does not convert Recovery R1 into a pass and grants no retry or
remote authority. PSCAN-06 remains open and unaccepted overall.

Current post-repair local validation against exact evidence candidate
`a0ac587f97557b89beb3b61553fe621e80f26611` independently rechecked the complete
351-file saved-project tree, all 12 hostile source-trust states and the Windows
no-Docker image-admission/process-boundary matrix. All local checks passed and
the hostile boundary recorded zero Docker calls. Exact results are in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-LOCAL-POST-REPAIR-VALIDATION-001.md`.
This result does not supply actual-Linux, genuine-Docker, dependency,
reproducibility, workflow or artifact proof and grants no retry authority.

## Correction C2 main-branch execution Recovery R2

On 2026-09-08 the owner authorized one bounded Recovery R2 attempt and allowed
GitHub-hosted Docker when the required proof cannot be completed by the Docker
installation on the owner's PC. Local Docker cannot establish the required
GitHub workflow/ref, standard Linux-runner, run/job/step, one-day artifact and
remote read-back identities. Recovery R2 therefore authorizes only the exact
committed GitHub Actions build job for the remaining actual-Linux and genuine-
Docker proof; it does not authorize starting, stopping or reconfiguring local
Docker.

The approval session checked only `gh auth status --hostname github.com` as
requested. The GitHub CLI is installed and the intended account is active, but
the stored token is invalid. No repository query or remote mutation was
attempted. The owner must re-authenticate before the genuinely fresh Recovery
R2 execution task begins, and that task must repeat the complete current
preflight rather than reuse this failed check.

Exact bounded authority is recorded in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R2-AUTHORITY.md`.
This authority-recording session performs no login, push, tag, ruleset, setting,
workflow, Docker, dependency, build or artifact action. Recovery R2 permits at
most one USD 0 workflow dispatch after every preflight fact passes. Signing,
attestation, draft creation, publication, rerun and successor work remain
forbidden. PSCAN-06 remains open and unaccepted overall.

Recovery R2 was subsequently attempted in genuinely fresh saved-project
`main` session `01a07dbf-b19b-71d1-b40c-617f396b1c64` using Daybreak Blue at
`xhigh`. The exact local execution context and complete current authority chain
passed. The first current GitHub authentication command then returned exit code
`1` because the intended active account's stored token was invalid. Required
execution-context item 6 and mandatory preflight item 4 therefore failed. The
task stopped before any GitHub repository query, remote mutation, workflow
dispatch, Docker, build or artifact action.

Recovery R2 is terminally failed, consumed zero workflow dispatches and
authorizes no retry. Exact append-only evidence is recorded in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R2-FAILURE-001.md`.
Its SHA-256 is
`CACDEE7FA590FEA4714A921ECC000AF6F008F4E00515DEDE91C1BEE3E37BAC91`.
Valid authentication and a new exact owner decision are required before any
further proof-gate attempt. PSCAN-06 remains open and unaccepted overall;
PSCAN-07 remains proposed and unselected; PSCAN-08 remains inactive and
ineligible.

## Correction C2 main-branch execution Recovery R3

On 2026-09-09 the owner authorized one bounded Recovery R3 attempt from exact
Recovery R2 failure record SHA-256
`CACDEE7FA590FEA4714A921ECC000AF6F008F4E00515DEDE91C1BEE3E37BAC91` at
evidence commit `a4f5b85a3ac63ea31b4bb911be3445f2223b6db1`. The authority
preserved maximum spend USD 0, at most one workflow dispatch, one-day unsigned
workflow artifact only, and the prohibitions on local Docker, signing,
attestation, drafting, publication, rerun and successor work.

Recovery R3 ran in genuinely fresh saved-project `main` session
`01a08461-59eb-7721-b1f3-f11721e36bfd`. Exact local source trust, authority,
workflow bytes and host-level GitHub authentication passed. The current
authenticated credential could not read the required user billing, usage and
stop-usage state: every applicable endpoint returned HTTP `404` and reported
that the missing `user` scope was required. Public standard-runner minutes did
not prove artifact-storage allowance, overage and stop-usage state. Mandatory
preflight item 4 therefore failed.

The task stopped before remote mutation, workflow dispatch, Docker, build or
artifact action. Recovery R3 is terminally failed, consumed zero workflow
dispatches and authorizes no retry. Exact append-only evidence is recorded in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R3-FAILURE-001.md`.
Its SHA-256 is
`339E5D66B655B5753D239D5347F1BBAFF0E5858C4CCC289368196D014AE07C12`.
Any further attempt requires trusted current billing, artifact-storage
allowance and stop-usage evidence plus a new exact owner decision. PSCAN-06
remains open and unaccepted overall; PSCAN-07 remains proposed and unselected;
PSCAN-08 remains inactive and ineligible.

## Correction C2 main-branch execution Recovery R4

On 2026-09-11 the owner exactly authorized one bounded Recovery R4 attempt from
Recovery R3 failure record SHA-256
`339E5D66B655B5753D239D5347F1BBAFF0E5858C4CCC289368196D014AE07C12` at
evidence commit `6b7ae32f8768c1e4687e5134696688088f65386b`. Approval-time
read-only verification confirmed the authenticated account now has the `user`
scope, current Actions Linux and storage usage both have net amount `USD 0`,
and the account-level Actions budget is `USD 0` with stop usage enabled.

Recovery R4 preserves maximum spend `USD 0`, at most one workflow dispatch,
one unsigned uncompressed workflow artifact retained for one day, no local
Docker action, and the prohibitions on signing, attestation, drafting,
publication, rerun and successor work. Exact authority is recorded in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R4-AUTHORITY.md`.
This authority-recording session does not claim or execute Recovery R4. A
genuinely fresh saved-project `main` task must begin from the committed R4
authority bundle and repeat the complete fail-closed preflight before any
remote mutation or workflow action. PSCAN-06 remains open and unaccepted
overall; PSCAN-07 remains proposed and unselected; PSCAN-08 remains inactive
and ineligible.

Recovery R4 was then attempted in genuinely fresh saved-project `main` session
`01a08f7a-7d13-7ee1-be72-a997e4ff9cb6` using Daybreak Blue at `xhigh`. The
complete mandatory preflight passed, including every-byte local source trust,
workflow identity and permissions, current authenticated repository and
Actions controls, exact locked tags/rulesets, absent C2 signing gate/tag/run/
artifact, zero-net Actions Linux and storage usage, and the live account-level
`USD 0` Actions budget with stop usage enabled.

The ordered transaction fast-forwarded remote `main` exactly from
`3523e4409ebc53cc1931e3dcaf7d1eea74bb15ab` to Recovery R4 authority commit
`a82a3a64a04ee2d8b60757866c9864cd1b73b54b`. It then created the exact unsigned
local tag `release-tooling-v1.0.0-c2` at accepted candidate
`faef8435322c9096df09b56969662411f17356ea`. The first tag-push command failed
locally because its shell-expanded refspec was invalid. The no-retry condition
stopped the transaction without correcting or repeating the command.

Recovery R4 is terminally failed and consumed zero workflow dispatches and
zero artifacts. The remote C2 tag and ruleset remain absent; no Docker, build,
signing, attestation, draft, release or publication action occurred. The exact
local C2 tag remains preserved. Exact append-only evidence is recorded in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R4-FAILURE-001.md`.
Its SHA-256 is
`5F694FF128515A2B79A2BFA500F299FECAB8836763365E9D774DA5D847D68BEE`.
Any future attempt requires a new exact owner decision addressing remote
`main` at `a82a3a64`, the preserved local C2 tag, and the absent remote tag and
ruleset. PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains
proposed and unselected; PSCAN-08 remains inactive and ineligible.

## Correction C2 main-branch execution Recovery R5

On 2026-09-11 the owner exactly authorized one bounded Recovery R5 attempt from
Recovery R4 failure record SHA-256
`5F694FF128515A2B79A2BFA500F299FECAB8836763365E9D774DA5D847D68BEE` at
evidence commit `6fbeb633bf373d2f64d95d0ff93c328631a99a05`.

Recovery R5 preserves the exact split state: remote `main` remains at
`a82a3a64a04ee2d8b60757866c9864cd1b73b54b`; the exact local lightweight tag
`refs/tags/release-tooling-v1.0.0-c2` remains at accepted candidate
`faef8435322c9096df09b56969662411f17356ea`; and the remote C2 tag, C2
ruleset, C2 workflow run and C2 artifact remain absent. It authorizes no tag
recreation or remote `main` push.

After a complete genuinely fresh preflight, Recovery R5 permits at most one
non-force tag push with literal lowercase refspec
`refs/tags/release-tooling-v1.0.0-c2:refs/tags/release-tooling-v1.0.0-c2`, one
exact-name no-bypass update/deletion-denial tag ruleset, and at most one
workflow dispatch. Maximum spend remains `USD 0`; only one unsigned,
uncompressed artifact retained for one day is permitted. Local Docker,
signing, attestation, drafting, publication, rerun and successor work remain
forbidden.

Exact authority is recorded in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R5-AUTHORITY.md`.
This authority-recording session does not claim or execute Recovery R5. A
genuinely fresh saved-project `main` task using Daybreak Blue at `xhigh` must
begin from the committed R5 authority bundle and repeat the complete
fail-closed preflight before the single tag-push attempt or any other remote
action. PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains proposed
and unselected; PSCAN-08 remains inactive and ineligible.

Recovery R5 was then attempted in genuinely fresh saved-project `main` session
`01a08fae-2ee1-7553-a126-e635b70a34a6` using Daybreak Blue at `xhigh`. The
complete current preflight passed, including every-byte source trust, workflow
identity and permissions, authenticated repository and Actions controls,
locked predecessor tags/rulesets, the preserved split state, zero prior C2
runs/artifacts, zero-net Actions Linux and storage billing, and the live
account-level `USD 0` Actions budget with stop usage enabled.

The exact existing local C2 tag was pushed once using the sole literal
refspec, read back at accepted candidate
`faef8435322c9096df09b56969662411f17356ea`, and protected by active no-bypass
ruleset `22895383` denying update and deletion. Remote `main` remained at
`a82a3a64a04ee2d8b60757866c9864cd1b73b54b`. The exact workflow was then
dispatched once as run `34582887399`, attempt `1`, from the locked C2 tag.

The standard `ubuntu-24.04` `build` job passed its identity, checkout,
product-tag, exact-entrypoint and source-trust steps, then failed in
`Prove private admission and bounded native process cases without Docker`.
The clean isolated native matrix attempted a lowercase `$pid` assignment and
PowerShell failed because automatic `$PID` is read-only. All later Docker,
dependency, build, byte-comparison and artifact-transfer steps were skipped.
The signing job was skipped, and no artifact, release or deployment exists.

Recovery R5 is terminally failed and consumed its sole tag push and workflow
dispatch. Exact append-only evidence is recorded in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R5-FAILURE-001.md`.
Its SHA-256 is
`53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0`.
No rerun, repair, second push, second dispatch, signing, publication or
successor action is authorized. PSCAN-06 remains open and unaccepted overall;
PSCAN-07 remains proposed and unselected; PSCAN-08 remains inactive and
ineligible.

## Correction C2 iteration 006 approved boundary

On 2026-09-11 the owner exactly approved the local-only Correction C2
iteration 006 PowerShell PID-collision correction from immutable Recovery R5
failure record SHA-256
`53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0` at
evidence commit `da45bb19480663aa21e9a67754ad509834dcaed7`.

The bounded defect is the lowercase `$pid` assignment inside
`Get-LinuxSessionMembers` in `build/release/docker-execution.ps1`. PowerShell
variable names are case-insensitive, so that local name collides with read-only
automatic variable `$PID`. Iteration 006 may rename only that local parsed
process identifier and its direct identity/ledger uses, with a bounded
regression proving the collision is absent and the process identity semantics
remain intact. Every iteration-005 isolation and earlier containment, Docker,
release, schema and product boundary remains unchanged.

Exact scope, checks and exclusions are controlling in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHORITY.md`. Iteration 006 is
owner-approved and was subsequently claimed alone in genuinely fresh saved-
project session `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2` from exact authority
commit `de162e8c347c725a0f041d3c8b0f51df1211c12d`. Bounded local author
candidate `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`, tree
`de70410df5ec359b85d03262e46c8d2ed5f2d822`, renames only the colliding local
binding and its two dependent uses and adds the exact-source regression. The
complete available-host local no-Docker matrix and exact-commit source checks
pass. Fresh independent review then rejected local acceptance because the AST
regression detects a reserved `PID` assignment only when the assignment left
side is directly `VariableExpressionAst`. Valid `[int]$PiD = 1` instead has a
`ConvertExpressionAst` left side, reproduces the read-only automatic-variable
collision, and evades the candidate predicate. Exact evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REVIEW-REJECTED.md`. The
reviewer did not repair the candidate.

The same claimed implementation session subsequently refined only the rejected
source regression. Refined author candidate
`e6a767aab71db1d3f62063dded379b4701d2cb52`, tree
`288664447272f9543de1b69b8ca28c27c6e1e9ff`, preserves the production blob
byte-for-byte, unwraps typed assignment-target AST wrappers, and adds bound
untyped, typed and non-target adversarial predicate self-tests. Exact-commit
source trust, both parsers, exact two-path confinement, the strengthened PID
source/data-flow regression and the complete available-host no-Docker native
matrix pass. Evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REFINEMENT-001.md` and
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION-002.md`.
Fresh independent review then rejected local acceptance again. The refined
predicate correctly unwraps typed and attributed targets but not valid
parenthesized or multiple assignment targets. `($PiD) = 1`,
`$PiD, $other = 1, 2` and `($PiD, $other) = 1, 2` each reproduce the
read-only automatic-variable collision while the exact predicate reports zero
reserved assignments. Exact evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REVIEW-REJECTED-002.md`. The
reviewer did not repair the refined candidate.

Recovery R5 remains terminally failed and its sole tag push and workflow
dispatch remain consumed. This correction authorizes no remote read or
mutation, tag change, workflow dispatch or rerun, signing, publication or
successor action. Actual-Linux, genuine-Docker, dependency, build,
reproducibility and artifact-integrity proof remains open. PSCAN-06 remains
open and unaccepted overall; PSCAN-07 remains proposed and unselected;
PSCAN-08 remains inactive and ineligible.

The same claimed implementation session then refined only the regression a
second time while preserving both rejected candidates and rejection records.
Second refined author candidate
`763730403ae538850e9d806ce9db82c815517b2d`, tree
`05c0baba5e61edb808b918a2dd6c853f49b8955e`, leaves production blob
`2fb43bf89df352e9153ba5d7ad23fd59cad77749` unchanged. Its principled walker
follows only exact assignment-target typed/attributed, parenthesized
pipeline/command-expression and recursively nested multiple-target wrappers;
member/index and drive-qualified non-automatic targets remain controls.

The exact earlier untyped/typed cases, all three second-review counterexamples,
nested wrapper cases and RHS/string/comment/ledger/member/index controls are
bound to the production predicate. Exact-commit source trust, both parsers,
two-path confinement, committed-byte identities and the complete available-
host no-Docker native matrix pass. Evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REFINEMENT-002.md` and
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION-003.md`.
This is author evidence only; fresh independent review is still required and
no acceptance is claimed.

Fresh independent review then rejected local acceptance a third time. The
walker treats every true variable path whose final colon-delimited component
is `PID` as automatic variable `PID`. Legal braced variables such as
`${variable:env:PID}`, `${global:env:PID}` and `${local:foo:PID}` instead bind
distinct variable names, parse and execute without the automatic-variable
collision, but the exact committed predicate reports one reserved assignment
for each. The author-bound drive controls do not cover this valid multi-colon
family, so the built-in regression passes despite the false positive. Exact
evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REVIEW-REJECTED-003.md`. The
reviewer did not repair the candidate.

Recovery R5 remains terminally failed and its sole tag push and workflow
dispatch remain consumed. This refinement authorizes no remote read or
mutation, Docker, tag change, workflow dispatch or rerun, signing,
publication or successor action. Actual-Linux, genuine-Docker, dependency,
build, reproducibility and artifact-integrity proof remains open. PSCAN-06
remains open and unaccepted overall; PSCAN-07 remains proposed and unselected;
PSCAN-08 remains inactive and ineligible.

The same claimed implementation session then refined only the regression a
third time while preserving all three rejected candidates and review records.
Third refined author candidate
`cf1679f9bca24887340fa4060f37d1ebff21f305`, tree
`6cb0198c71881baff9c601c31ad6aa84217f4a49`, leaves production blob
`2fb43bf89df352e9153ba5d7ad23fd59cad77749` unchanged. The predicate now uses
public `VariablePath` scope flags and exact whole `UserPath` equality for only
the runtime-proven automatic PID forms. It no longer splits arbitrary legal
variable names on colons.

All earlier wrapper positives and controls remain bound to the exact production
predicate. New parser and isolated runtime cases prove unqualified,
`variable:`, global, script, local and private PID targets collide, while the
third-review multi-colon variables and env/function provider targets do not.
Exact-commit source trust, both parsers, two-path confinement, committed-byte
identities and the complete available-host no-Docker matrix pass. Evidence is
in `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REFINEMENT-003.md` and
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION-004.md`.
This is author evidence only; a fourth fresh independent review is required.

Fourth independent review then accepted bounded local iteration 006.
The reviewer extracted the exact predicate, verified the public
`VariablePath` contract against the installed parser, and exercised all prior
positives and controls plus independently mutated qualifier case, bracing,
escaping, extra-colon and separator placement, wrapper depth, tuple position,
provider paths and every installed-parser assignment-left family. Every
predicate result agreed with isolated runtime collision behavior; supported
non-PID lvalues were allowed and unsupported parser-invalid target shapes
failed closed. Production blob and exact identity/ledger data flow, both
parsers, source trust, path confinement, the complete available-host no-Docker
matrix and residual cleanup all pass. Exact evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-ACCEPTANCE.md`. This accepts
only bounded local iteration 006, not PSCAN-06 overall or any remote proof.

Recovery R5 remains terminally failed and consumed. No remote read or
mutation, Docker, tag change, workflow dispatch or rerun, signing,
publication or successor action is authorized. Actual-Linux, genuine-Docker,
dependency, build, reproducibility and artifact-integrity proof remains open.
PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains proposed and
unselected; PSCAN-08 remains inactive and ineligible.

## Correction C2 iteration 007 immutable R6 identity authority

The owner then directed autonomous, repository-consistent continuation and
authorized this authority session to select Correction C2 iteration 007 if an
immutable identity prerequisite proved necessary before Recovery R6. Exact
Git-object and source inspection proved that prerequisite.

The active protected tag `release-tooling-v1.0.0-c2` remains at pre-fix
candidate `faef8435322c9096df09b56969662411f17356ea`. The accepted
iteration-006 candidate `cf1679f9bca24887340fa4060f37d1ebff21f305` retains
the same committed C2 workflow and schema bytes. That workflow requires the
old tag, its SHA and its workflow ref simultaneously; schema 2.1, the builder
and verifier bind the same old identity. Dispatching the old tag reruns the
known defect, using the accepted candidate as input fails the SHA gate, using
a new ref fails the ref gate, and moving the protected tag would violate
immutable evidence.

Iteration 007 was owner-approved, unclaimed and local-only at its authority
commit. Fresh saved-project session
`01a09151-048a-7bf3-a3a6-2d0c26415f10` then claimed only iteration 007 from
the exact clean authority commit and implemented the bounded append-only
identity candidate. It
selects exact proposed identity `release-tooling-v1.0.0-c2-r6`, new workflow
path `.github/workflows/release-recovery-v1.0.0-c2-r6.yml` and additive
manifest schema 2.2. It may implement and locally verify only the append-only
identity/version roll-forward under the exact path and preservation rules in:

`evidence/PSCAN-06/CORRECTION-C2-ITERATION-007-AUTHORITY.md`

The implementation adds the new workflow and schema, rolls the active builder
and embedded verifier to 2.2, extends the generic verifier while preserving
2.0/2.1 behavior, adds cross-version/source-agreement tests and synchronizes
the living records. `docker-execution.ps1` changes only to materialize schema
2.2; the accepted iteration-006 PID predicate/data flow and Docker operation
table remain unchanged. This is author implementation, not acceptance.

No tag, Docker, dependency, remote, signing, publication or spending action is
authorized or performed. A separate fresh session must independently review
the exact committed candidate. Recovery R6 remains unauthorized until
iteration 007 is independently accepted and a later exact R6 authority
re-proves all live zero-spend remote gates. PSCAN-06 remains open and
unaccepted overall; PSCAN-07 remains proposed and unselected; PSCAN-08 remains
inactive and ineligible.

Fresh independent review session `01a0916e-95e6-7db0-9391-08232b2fc505`
inspected exact candidate `230e1761f7c45d9629cecf360e7498b33d64ab6f`,
tree `5d74628d76f95e150c0ea1de51b4892af2350a7d`, from direct authority parent
`b7a0fb211e4631f915e368174ae58022d1519182`. The new workflow and schema
match the selected append-only identity design, old workflow/schema bytes are
preserved, active identity sources agree, and the Docker file changes only by
the schema-2.2 materialization line. Exact source trust, path confinement,
PowerShell/JSON parsing and every available-host no-Docker regression passed.

Local acceptance was rejected because the committed integration test's final
workflow normalization globally replaces `2.2` with `2.1`. That also changes
the unchanged `actions/attest` comment `# v4.2.2` to `# v4.2.1` in only the
normalized new workflow, so the required preservation test deterministically
fails. Exact evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-007-REVIEW-REJECTED.md`. The
reviewer did not repair the candidate.

Recovery R5 remains terminally failed and consumed. Recovery R6 remains
unauthorized. No remote read or mutation, Docker, dependency acquisition, tag
change, workflow dispatch or rerun, signing, publication, spending or
successor action is authorized by this review. Actual-Linux, genuine-Docker,
dependency, build, reproducibility and artifact-integrity proof remains open.
PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains proposed and
unselected; PSCAN-08 remains inactive and ineligible.

## Correction C2 iteration 008 preservation-test authority

The owner selected iteration 008 as the minimal append-only correction after
the iteration-007 independent review. Candidate
`230e1761f7c45d9629cecf360e7498b33d64ab6f`, tree
`5d74628d76f95e150c0ea1de51b4892af2350a7d`, remains immutable and rejected
for one bounded reason: its workflow-preservation test globally replaces every
`2.2` substring, changing the unchanged `# v4.2.2` action comment and forcing
a false comparison failure. The workflow/schema identity design and all other
recorded no-Docker/static checks remain green evidence, not acceptance.

Iteration 008 is owner-approved, unclaimed and not implemented. It may change
only the workflow-normalization portion of
`tests/integration/supply-chain/release_test.go`, add new append-only
iteration-008 evidence and synchronize `README.md`, this task, the reading map
and tracker. The repair must normalize only the exact intended R6
workflow/schema-operation/reference tokens with proved occurrence counts,
preserve `# v4.2.2`, and keep every workflow, schema, builder, Docker,
verifier, product and acceptance-test byte unchanged.

Exact scope, byte identities, offline Go 1.27.1 validation route, acceptance
matrix and exclusions are controlling in:

`evidence/PSCAN-06/CORRECTION-C2-ITERATION-008-AUTHORITY.md`

This authority-recording session does not edit the test or execute Go. A
genuinely fresh session must claim only iteration 008 from the exact committed
authority bundle, perform the bounded repair, validate it offline and stop for
fresh independent review. Recovery R5 remains terminal and consumed; Recovery
R6 remains unauthorized. No Docker, dependency download, network, remote,
tag, workflow, signing, publication, spending or successor action is
authorized. PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains
proposed and unselected; PSCAN-08 remains inactive and ineligible.

Fresh saved-project session `01a091b5-687a-7392-b366-4b71f8bbbe8d` claimed
only iteration 008 from exact authority commit
`4f86651de67a05db8ae0076c0c155bd184082af3`. The exact source-trust,
authority, critical-byte, token-count and Go 1.27.1 toolchain gates passed. Its
bounded workflow-normalization repair was applied in the working tree, but the
required `gofmt -d` check returned native exit `1` with a 66-line diff outside
that permitted block. A separate clean materialization of the untouched
authority commit reproduced the same formatter result and four pre-existing
hunks.

The session stopped on that unrelated failure without broadening the repair.
The test attempt was reverted byte-for-byte to the authority commit; no
iteration-008 candidate, author-validation pass or acceptance exists. Exact
non-acceptance evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-008-IMPLEMENTATION.md` and
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-008-AUTHOR-VALIDATION.md`. A
separate fresh verification continuation reproduced the expected untouched
targeted failure, the full-package Windows synthetic Cosign failure, and a
passing targeted result under the bounded probe; the probe was restored and
the full package still failed solely at Cosign. A new exact owner decision is
required before another implementation attempt.
Recovery R5 remains terminal and consumed; Recovery R6 and all Docker,
dependency, remote, release and successor work remain unauthorized.

## Correction C2 iteration 009 single-file test authority

On 2026-09-12 the owner selected iteration 009 after iteration 008 stopped
with no candidate. The exact owner-selected boundary reconciles only the two
pre-existing within-file blockers with the already proved normalization
repair. Read-only diagnostic task
`01a091b7-90ec-7561-bd79-c8a0f91a522f` established that the synthetic
extensionless POSIX Cosign script cannot execute on Windows: executable
resolution returns `*exec.Error` wrapping `exec.ErrNotFound` before
`CreateProcess`. The production verifier behaves as designed.

Iteration 009 may modify only
`tests/integration/supply-chain/release_test.go`, add new iteration-009
evidence and synchronize the four named living-state documents. It must:

1. apply the exact iteration-008 five-token normalization with counts
   `1,1,8,3,1`, exact schema-assertion replacement, workflow byte equality,
   preserved `# v4.2.2` and absent `# v4.2.1`;
2. apply only the demonstrated 66-line/four-hunk Go 1.27.1 mechanical
   formatting projection outside the helper and normalization; and
3. replace only the POSIX synthetic fixture with a current-test-executable
   `TestMain` helper that admits the exact 18-argument, fixed-order,
   non-empty, temp-root-contained invocation under `COSIGN_YES=false`, while
   preserving every claim assertion and adding hostile near-miss coverage.

Exact scope, helper admission, offline/frozen-cache checks, Linux-amd64
cross-compile, acceptance matrix and exclusions are controlling in:

`evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-AUTHORITY.md`

This authority-recording session does not edit or execute the test. A
genuinely fresh session must claim only iteration 009, implement and validate
the bounded correction offline, commit its candidate and stop for genuinely
fresh independent review. Actual Linux runtime remains separately gated and
does not block this local test-only correction. Recovery R5 remains terminal
and consumed; Recovery R6 and all Docker, network, remote, tag, workflow,
signing, publication, spending and successor work remain unauthorized.
PSCAN-06 remains open and unaccepted overall.
