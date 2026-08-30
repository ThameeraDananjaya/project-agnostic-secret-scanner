# PSCAN-01 Task Specification

## Title

Repository governance, canonical specification, threat model, licence plan, and
release authority.

## State

Activated. Unclaimed. Not implemented. Implementation must begin in a fresh
session from this activation commit.

## Objective

Create the complete local governance and planning authority needed to control
PSCAN-02 through PSCAN-08 without implementing product code or exercising any
remote, spending, credential, signing, publication, or consuming-project gate.

## Allowed paths

```text
AGENTS.md
CONSTITUTION.md
README.md
SECURITY.md
LICENSE
THIRD_PARTY_NOTICES.md
.gitignore
.gitattributes
.editorconfig
.github/CODEOWNERS
docs/spec/**
docs/architecture/**
docs/governance/**
docs/security/**
docs/validation/**
docs/release/LICENSING.md
docs/release/RELEASE-AUTHORITY.md
docs/tasks/**
docs/decisions/**
evidence/PSCAN-01/**
```

## Forbidden paths

```text
cmd/**
internal/**
contracts/**
rules/**
fixtures/**
tests/**
build/**
go.mod
go.sum
.github/workflows/**
graphify-out/**
any scanner or third-party binary
```

## Forbidden actions

- Product-code or schema implementation.
- Dependency, toolchain, scanner, or action download.
- Scanner execution against any project data.
- Credential or signing-key creation or handling.
- Remote creation, push, PR, settings change, signing, or publication.
- Spending or paid capability enablement.
- TruffleHog assessment, download, integration, distribution, or enablement.
- Consuming-project data or behavior import.
- PSCAN-02 selection or activation.

## Required deliverables

- Byte-verified canonical `docs/spec/PASS-SPEC-001.md`.
- Complete CAP-1 through CAP-20 and acceptance-row traceability matrix.
- Repository constitution and working agreement.
- Full proposed source layout and component/trust-boundary architecture.
- Threat model and validation plan.
- Schema-ownership and policy-projection decision.
- Licence plan and release-authority plan.
- Task specifications for PSCAN-02 through PSCAN-08.
- Evidence and closeout rules.
- Plain-English runbook requirements, including Mermaid diagrams and complete
  input/output documentation, carried into later tasks.

## Checks

- Canonical contract hash and line-count verification.
- Requirement traceability completeness.
- Allowed/forbidden path validation.
- No consuming-project identifiers or data.
- No product code, workflows, binaries, dependencies, credentials, or secrets.
- Document link/reference validation.
- Licence-boundary review.
- Independent read-only review of requirement preservation and task bounds.
- Clean Git status after the accepted local commit.

## Success criteria

PSCAN-01 is complete only when the repository has enforceable governance, the
contract is canonical and completely traceable, the architecture/threat/validation
plans preserve every contract boundary, PSCAN-02 is bounded but unselected, all
owner gates remain explicit, exact acceptance evidence exists, and no forbidden
path or action occurred.
