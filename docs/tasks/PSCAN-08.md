# PSCAN-08 Task Specification

## Title

Optional TruffleHog fallback integration for an accepted material Gitleaks
coverage gap.

## State

Inactive. Unselected. Not eligible. Not activated. Not claimed.

## Eligibility and owner gates

PSCAN-08 cannot be selected or activated unless accepted PSCAN-10/07 evidence
proves a material required input class cannot be covered by Gitleaks. Work then
stops for one consolidated owner decision covering:

1. exact technical need and why normalizer/Gitleaks changes cannot close it;
2. current official TruffleHog source, release, licence and security evidence;
3. separate acceptance of AGPL obligations and distribution consequences;
4. exact build, request, release, rule and consumer boundary;
5. confirmation of zero spend and no provider verification/network access.

No legal assessment, download or exploratory integration occurs before material-
gap evidence exists. Eligibility is not approval; approval is not activation.

## Objective if separately authorized

Implement the minimum local-only fallback needed for the exact accepted gap,
behind both build-time and request-time disabled-by-default gates, without
weakening any Gitleaks result or project-independence boundary.

## Provisional allowed paths

These paths become authoritative only in a separately approved activation bundle:

```text
go.mod
go.sum
internal/engine/trufflehog/**
internal/engine/**
contracts/scan-request/**
contracts/scan-outcome/**
fixtures/adversarial/trufflehog/**
tests/unit/engine/trufflehog/**
tests/integration/trufflehog/**
tests/acceptance/trufflehog/**
build/trufflehog/**
licenses/trufflehog/**
THIRD_PARTY_NOTICES.md
docs/architecture/**
docs/security/**
docs/validation/**
docs/release/**
docs/tasks/PSCAN-08.md
docs/tasks/TRACKER.md
evidence/PSCAN-08/**
```

## Permanent controls

- Local Git, filesystem or local image-tar inputs only.
- Provider verification, caches, enumeration, telemetry, updates and network
  disabled.
- Same isolated network-disabled sandbox and private output capture.
- Exact approved gap/request/binary/adapter binding required.
- Any finding from either required engine fails; no engine can override another.
- Incomplete fallback coverage is indeterminate; unavailable required fallback
  is non-pass.
- Binary distributed only in a separately identifiable fallback bundle with
  complete licence evidence.

## Acceptance if activated

The exact gap is covered on both platforms, disabled gates are proven, offline
and redaction tests pass, conflict aggregation is fail-closed, licence and SBOM
evidence is complete, and no extra input class or provider behavior is admitted.

Until all prerequisites and gates exist, PSCAN-08 remains inactive.
