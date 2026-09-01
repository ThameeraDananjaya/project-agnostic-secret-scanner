# PSCAN-10 Task Specification

## Title

Outcome-proved Gitleaks primary coverage for exact Git objects and deterministic
scan inputs.

## State

Proposed. Unselected. Not activated. Not claimed. This planning label is not
authority to implement.

## Objective

Implement the minimum primary-engine redesign required by
PASS-OUTCOME-SPEC-001: exact Git object/range/tree admission, raw input
classification before transformation, byte-conserving preparation, and proved
Gitleaks inspection for every admitted byte under bounded declared profiles.

## Preconditions

- PSCAN-09 accepted, committed and clean.
- Exact PSCAN-10 owner selection and activation in an activation-only session,
  followed by a fresh implementation session from that activation commit.
- Current exact Gitleaks source, binary, rules, detector span/stream/archive
  behavior, Go toolchain/module build and licence facts revalidated from primary
  sources before any intake or rebuild.
- A bounded design proves finite maximum rule spans or a complete streaming path
  before any pass-capable projection is implemented.

## Provisional allowed paths

These paths become authoritative only in a future activation bundle:

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
docs/tasks/PSCAN-10.md
docs/tasks/TRACKER.md
evidence/PSCAN-10/**
```

## Forbidden scope

No Gitleaks GitHub Action; no TruffleHog assessment or material; no artifact/
OCI normalizer beyond raw classification needed to fail closed; no final
redaction/offline-container claim; no project policy, allowlist, receipt,
workflow, remote, signing, publication, credential, spending or project data.
No reuse of rejected PSCAN-03 evidence as acceptance.

## Required deliverables

- Exact base/head/merge-base/range/parent-edge/head-tree and Git-object identity
  proof with hostile Git behavior disabled.
- Private admission ledger binding every required object, byte length, digest,
  raw class, preparation, invocation and inspection proof.
- Raw archive/compression/container/ambiguous classification before framing or
  transformation, with unsupported classes explicit non-pass.
- Byte-conserving preparation with a proved mapping for every original byte.
- Exact detector buffer/fragment/stream/archive and maximum-rule-span evidence
  for the pinned binary and rules.
- Bounded named PR/release profiles with at/below/above resource tests.
- Separate text, binary, deleted-history, head-tree, merge, rename, boundary,
  archive-class, maximum-size and skip-path adversarial fixtures.
- Content-free evidence and an updated coverage map that identifies any
  remaining material required-class gap without selecting PSCAN-08.

## Acceptance

Pass is possible only when every required object and byte has exact admission
and actual detector-inspection proof. Prefix markers, file-open evidence,
aggregate failures, unproved fragments, masked raw classes, unbound limits and
unsupported inputs are non-pass. Exact binding mismatch stops before scanning;
candidate material never becomes executable syntax; Gitleaks remains primary;
no forbidden scope is introduced.

Stop after PSCAN-10 closeout with PSCAN-04 and every other successor unselected.
If accepted evidence establishes a material Gitleaks gap, stop for the separate
PSCAN-08 technical and AGPL owner decision; do not infer eligibility approval,
selection or activation.
