# PSCAN-03 Task Specification

## Title

Controlled Gitleaks source intake, build, adapter, Git/file scanning, and
fixture coverage.

## State

Activated by the owner on 2026-08-31. Unclaimed. Not implemented.

Implementation is authorized only in a fresh session from the exact PSCAN-03
activation commit. This activation session contains no implementation work.

## Objective

Admit one currently supported, licence-reviewed, exact Gitleaks CLI source and
binary binding; implement safe argument-array orchestration for exact Git ranges,
tracked trees and ordinary directories/files; and prove primary-engine coverage
with synthetic fixtures.

## Preconditions

- PSCAN-02 accepted and clean.
- Exact PSCAN-03 activation and fresh implementation session.
- Current official Gitleaks release, source revision, licence, security posture,
  output format and reproducible-build requirements revalidated.

## Allowed paths

```text
go.mod
go.sum
internal/gitinput/**
internal/engine/**
internal/engine/gitleaks/**
rules/generic/**
fixtures/synthetic-findings/**
fixtures/history/**
fixtures/adversarial/git/**
tests/unit/gitinput/**
tests/unit/engine/**
tests/integration/gitleaks/**
tests/acceptance/gitleaks/**
build/gitleaks/**
licenses/gitleaks/**
THIRD_PARTY_NOTICES.md
docs/architecture/**
docs/security/**
docs/validation/**
docs/release/LICENSING.md
docs/tasks/PSCAN-03.md
docs/tasks/TRACKER.md
evidence/PSCAN-03/**
```

## Forbidden scope

No Gitleaks GitHub Action; no TruffleHog assessment or material; no archive/OCI
normalizer; no final redaction or offline-container claim; no project policy,
allowlist, receipt, workflow, remote, signing, publication or project data.

## Deliverables

- Exact upstream source/tag/commit, licence, source checksum, built binary digest,
  adapter version, rule/config digest and output-format binding.
- Safe Git preparation with inherited config, hooks, filters, pagers, external
  diff and text conversion disabled.
- Git mode for exact ranges/history and file/directory mode for manifest inputs.
- Private capture with forced engine redaction, no banners/color/verbose/interactive
  behavior and explicit non-pass for skip/timeout/format/integrity errors.
- Synthetic clean, finding, deleted-history, binary/large-file, range and hostile
  Git fixtures; no real credentials.
- Coverage map identifying any required input class not inspectable by Gitleaks.

## Acceptance

Exact binding mismatch fails before scanning; candidate text never becomes
executable syntax; deleted in-range synthetic secrets fail; out-of-range facts
remain accurately bound; tracked source is complete or non-pass; all raw engine
output remains private; licence evidence is complete. A possible coverage gap is
recorded only as evidence and does not select, activate or authorize PSCAN-08.

Stop after PSCAN-03 closeout with PSCAN-04 unselected.
