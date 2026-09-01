# PSCAN-03 Task Specification

## Title

Controlled Gitleaks source intake, build, adapter, Git/file scanning, and
fixture coverage.

## State

Rejected and closed fail-closed on 2026-09-01. Correction C1 candidate
`6cd22a3` was independently rejected on 2026-08-31. Exact blob enumeration is
implemented in that non-authoritative candidate, but the same-run per-blob
marker does not prove complete Gitleaks fragment coverage. PSCAN-03 is not
accepted and its scanner changes remain non-authoritative.

Implementation was authorized only in a fresh session from the exact PSCAN-03
activation commit. The activation session contained no implementation work.

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

## Historical stop and correction C1

`evidence/PSCAN-03/MATERIAL-GAP-001.md` remains the immutable proof that native
Gitleaks `git log -p -U0` input omits deleted binary blob bytes. Correction C1
does not rewrite that evidence or use the incomplete input form. It enumerates
exact in-range commits, every parent-edge pre/post blob and the head tree,
materializes deterministic path-preserving binary-safe projections, and uses a
same-invocation Gitleaks coverage rule to make every skip non-pass.

Independent review proved that Gitleaks v8.30.1 can split a product-rule match
at its 125,000-byte no-whitespace boundary while reporting the single prefix
coverage marker. The pinned rule pack also contains unbounded whole-match
expressions, so no finite overlap establishes a complete maximum-span proof.
Recognized archives can likewise be masked by framing unless raw archive
classification is complete. These are acceptance blockers recorded in
`evidence/PSCAN-03/REVIEW-C1.md`; the correction candidate is non-authoritative.
The candidate also replaces PASS-SPEC-001 section 7.2's required Gitleaks Git
mode for exact history/range scanning with directory mode over projections.
The correction authorization did not waive that controlling requirement.

Native patch mode remains fail-closed. TruffleHog is absent. No successor is
selected or activated. On 2026-09-01 the owner abandoned the contract-preserving
option and directed future work toward an outcome-based replacement contract.
That direction does not change PASS-SPEC-001 inside PSCAN-03 because its
governance paths are outside this task's allowlist. PSCAN-09 is proposed only
for that bounded governance replacement. PSCAN-04 remains unselected and
PSCAN-08 remains inactive and owner-gated. The rejected closeout is recorded in
`evidence/PSCAN-03/CLOSEOUT-REJECTED.md`.
