# PSCAN-06 Task Specification

## Title

Reproducible cross-platform build, SBOM/licence bundle, keyless-signing release
workflow, immutable-release plan, and offline verifier guidance.

## State

Proposed. Unselected. Not activated. Not claimed. Remote, settings, spending,
signing and publication gates reserved.

## Objective

After current read-only preflight and exact action-time owner approval, establish
the trusted build/release boundary for exact Windows amd64 and Linux amd64 assets,
complete supply-chain evidence and offline consumer verification. PSCAN-06 does
not publish `v1.0.0`; final acceptance/publication belongs to PSCAN-07.

## Preconditions and owner gates

- PSCAN-05 accepted and clean; exact PSCAN-06 activation and fresh session.
- Revalidate official Go, Gitleaks, Cosign/Sigstore, GitHub Actions, immutable
  releases, artifact attestations, rules/settings, plan and cost facts.
- Read-only preflight must confirm exact account/repository/workflow identity,
  availability, proposed settings, permissions and zero-cost posture.
- Owner separately approves each remote creation, push/settings/workflow action
  and any signing dry run. Activation alone is not approval.

## Allowed paths

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

No secrets, long-lived signing keys, project receipt keys, consuming-project
integration, deployment gates, real candidate data, `latest`, mutable refs,
auto-update, pipe-to-shell, unpinned actions, TruffleHog, paid capability or
final release publication. No remote action without exact recorded approval.

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
