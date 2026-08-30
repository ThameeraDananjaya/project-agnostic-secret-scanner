# PSCAN-02 Task Specification

## Title

External schemas, reason codes, CLI skeleton, deterministic workspace, and
content-free outcome contract.

## State

Activated by the owner on 2026-08-30. Claimed and implemented locally from the
exact activation commit. Independent acceptance is pending the bounded
implementation commit and review.

Implementation is authorized only in a fresh session from the exact PSCAN-02
activation commit. This activation session contains no implementation work.

## Objective

Implement the scanner-owned contract boundary and a standard-library-first Go
runner skeleton that validates requests, creates an isolated deterministic
workspace, emits exactly one content-free outcome and contains no engine,
artifact, policy or remote integration.

## Preconditions

- PSCAN-01 accepted and committed with a clean worktree.
- Exact standalone owner activation of PSCAN-02 in an activation-only session.
- Fresh implementation session from that activation commit.
- Current Go/toolchain primary-source support and licence preflight before pinning.

## Allowed paths

```text
go.mod
go.sum
cmd/scanner-runner/**
internal/request/**
internal/workspace/**
internal/outcome/**
contracts/scan-request/**
contracts/scan-outcome/**
contracts/release-manifest/**
contracts/global-revocation/**
contracts/rule-pack/**
fixtures/clean/**
fixtures/adversarial/request/**
fixtures/adversarial/workspace/**
tests/unit/request/**
tests/unit/workspace/**
tests/unit/outcome/**
tests/integration/contract/**
docs/architecture/**
docs/validation/**
docs/tasks/PSCAN-02.md
docs/tasks/TRACKER.md
docs/decisions/DEC-001-SCHEMA-OWNERSHIP.md
evidence/PSCAN-02/**
```

## Forbidden scope

No engine process or Gitleaks/TruffleHog material; no Git-range scanner; no
artifact extraction; no policy/allowlist evaluation; no project-owned schema
authority; no receipt signing; no workflow, remote, credential, signing,
publication, spend or consuming-project integration.

## Deliverables

- Independently versioned scanner-owned schemas with strict major/minor behavior.
- Complete state/reason/exit mapping from PASS-SPEC-001.
- Non-interactive CLI skeleton reading local request input and emitting exactly
  one JSON outcome.
- Pre-engine request rejection for unsafe paths, missing/duplicate/mismatched
  bindings, unsupported versions and corrupt payloads.
- Fresh workspace per project/attempt with non-source-derived correlation ID,
  bounded paths and lifecycle proof.
- Content-free serializer that cannot represent forbidden finding metadata.
- Unit/integration/adversarial contract evidence on Windows and Linux where
  platform mechanics differ.

## Acceptance

Unknown majors and corrupt requests reject/fail closed; compatible minor rules
are proven; every reason and exit code is exact; repeated logical inputs are
semantically deterministic; injection cannot escape JSON/stderr constraints;
workspaces do not collide; no engine or project authority is present; CAP-1,
CAP-2, CAP-11, CAP-12, CAP-13 and contract portions of CAP-15 through CAP-18 are
traceably advanced without claiming final product acceptance.

Stop after PSCAN-02 closeout with PSCAN-03 unselected.

## Implementation result

- Toolchain: Go 1.27.0; standard library only; no `go.sum`.
- Scanner-owned contracts: scan request, scan outcome, release manifest, global
  revocation and generic rule pack, each at independent schema version 1.0.
- Valid skeleton behavior: validate all pre-engine bindings, create and destroy
  the isolated workspace, then emit `UNAVAILABLE_ENGINE`/exit 30.
- Engine, Git input, artifact normalization, project policy/allowlist authority,
  workflow and remote boundaries remain absent.
- Exact implementation and acceptance evidence is under `evidence/PSCAN-02`.
- PSCAN-03 remains proposed and unselected.
