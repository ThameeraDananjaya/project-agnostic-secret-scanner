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
activated.

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
