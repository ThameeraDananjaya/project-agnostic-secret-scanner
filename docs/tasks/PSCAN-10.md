# PSCAN-10 Task Specification

## Title

Outcome-proved Gitleaks primary coverage for exact Git objects and deterministic
scan inputs.

## State

Activated by the owner's exact `ACTIVATE PSCAN-10` command on 2026-09-01.
Unclaimed and not implemented. This activation bundle authorizes only a fresh
implementation session from the exact activation commit; the session creating
the bundle must stop without claiming or implementing PSCAN-10.

## Objective

Implement the minimum primary-engine redesign required by
PASS-OUTCOME-SPEC-001: exact Git object/range/tree admission, raw input
classification before transformation, byte-conserving preparation, and proved
Gitleaks inspection for every admitted byte under bounded declared profiles.

## Preconditions

- PSCAN-09 accepted, committed and clean at activation parent
  `7c15c4574c06f34f98b79c07548baee16c0d176c`.
- Exact PSCAN-10 owner selection and activation recorded in
  `evidence/PSCAN-10/ACTIVATION.md`, followed by a fresh implementation session
  from that activation commit.
- Current exact Gitleaks source, binary, rules, detector span/stream/archive
  behavior, Go toolchain/module build and licence facts revalidated from primary
  sources before any intake or rebuild.
- A bounded design proves finite maximum rule spans or a complete streaming path
  before any pass-capable projection is implemented.

## Activation-only bundle

This session may change and commit only:

```text
docs/tasks/PSCAN-10.md
docs/tasks/PSCAN-10-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-10/ACTIVATION.md
```

No implementation, dependency intake, scanner/toolchain execution, test run or
successor work is authorized in the activation session.

## Allowed implementation paths

These paths are authoritative only for a fresh PSCAN-10 implementation session
from the exact activation commit:

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

- No modification of PASS-OUTCOME-SPEC-001, PASS-SPEC-001, their source record,
  or PSCAN-01 through PSCAN-03 and PSCAN-09 historical acceptance evidence.
- No Gitleaks GitHub Action and no TruffleHog assessment, legal work, download,
  material, integration, distribution or enablement.
- No artifact or OCI normalizer beyond raw classification needed to fail
  closed; no final redaction, offline-container, release or consumer-readiness
  claim.
- No project policy, allowlist, receipt, workflow, identity, finding,
  credential, customer data, repository identity, statistics or project data.
- No remote creation, push, publication, release, signing, GitHub-setting
  change, credential handling, spending, deployment, production or go-live.
- No reuse of rejected PSCAN-03 implementation or evidence as acceptance; it is
  historical defect evidence only.
- No selection, activation, claim or implementation of PSCAN-04, PSCAN-08 or
  any other successor.

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
