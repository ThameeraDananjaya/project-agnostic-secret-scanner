# PSCAN-07 Task Specification

## Title

Full adversarial acceptance, Windows/Linux parity, performance, revocation,
retirement, and first signed v1 release.

## State

Proposed. Unselected. Not activated. Not claimed. Signing and publication gates
reserved.

## Objective

Run the complete PASS-OUTCOME-SPEC-001 acceptance matrix on the frozen release candidate,
independently verify first-consumer readiness, establish lifecycle evidence and,
only after exact owner approvals, sign and publish immutable `v1.0.0`.

## Preconditions and owner gates

- PSCAN-06 accepted and clean; exact PSCAN-07 activation and fresh session.
- Current release/source/workflow identities and account/settings/cost facts
  revalidated; exact final candidate frozen.
- Owner separately approves the exact signing workflow invocation and exact
  publication. Activation or test success is not publication approval.

## Allowed paths

```text
fixtures/**
tests/acceptance/**
tests/integration/**
docs/release/**
docs/security/**
docs/validation/**
docs/governance/RUNBOOK-REQUIREMENTS.md
docs/tasks/PSCAN-07.md
docs/tasks/TRACKER.md
SECURITY.md
THIRD_PARTY_NOTICES.md
evidence/PSCAN-07/**
```

Separately approved remote release assets are workflow outputs, not repository
paths, and do not broaden the allowed local file set.

Product defects discovered during acceptance fail PSCAN-07 closed and return to
the appropriate new bounded correction task; they do not authorize unbounded
source changes inside PSCAN-07.

## Forbidden scope

No consuming-project integration or real data; no receipt keys or deployment
authority; no weakening, skipped coverage or manual pass; no unapproved remote,
settings, signing, publication or spend; no TruffleHog unless PSCAN-08 has first
become eligible, separately approved, implemented and accepted.

## Deliverables

- Fresh proof for all CAP-1 through CAP-20 and AT-01 through AT-31 on the exact
  candidate, including every-byte detector inspection, raw classification,
  Windows/Linux parity and declared performance profiles.
- Coverage map and explicit CAP-6 decision: no material gap and fallback absent,
  or fail-closed referral to PSCAN-08 before readiness.
- Mock external project adapter/issuer proof using synthetic data and an
  ephemeral test-only trust root that never enters scanner execution or release.
- Complete runbook, release docs, compatibility, limitations, update, rollback,
  retirement, global revocation and content-free vulnerability process.
- Independent release-manifest, asset, identity, signature bundle, licence, SBOM,
  immutability and offline-verification evidence.
- If and only if separately approved: exact signed public `v1.0.0` and publication
  evidence. Otherwise stop before signing/publication without claiming readiness.

## Acceptance

Every capability and matrix row passes without exception; exact frozen assets
are mutually bound, immutable where supported and verifiable offline; no project
data or forbidden output exists; both platforms are equivalent; unsupported or
untrusted conditions cannot pass; global lifecycle controls work; independent
review confirms the scanner has no merge/deployment/business/credential authority.

Stop after PSCAN-07 closeout. No integration task or successor is selected.
