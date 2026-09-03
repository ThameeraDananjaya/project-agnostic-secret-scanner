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
payloads from a Windows CRLF checkout. Correction C1 iteration 003 repaired
that defect as exact candidate
`5ca77226ed3996a8267caf02da366d0beb915c8d`. The forced-CRLF regression,
fresh acquisition, two network-disabled builds, 30-file byte comparison and
all callable checks pass. Actual-Linux host positive/wrong-owner proof and
independent skeptical acceptance remain open. Exact author evidence is in
`evidence/PSCAN-06/CORRECTION-C1-AUTHOR-VALIDATION-003.md`. PSCAN-06 remains
the sole open task and remains unaccepted.

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
