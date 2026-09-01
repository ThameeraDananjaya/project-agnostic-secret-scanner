# DEC-002: Outcome-Based Primary Coverage Contract

- Status: accepted with PSCAN-09 independent acceptance and closeout
- Date: 2026-09-01
- Authority: owner activation of PSCAN-09 at
  `f48698922c174411f05124b066772753a2b8dd1e`
- Successor contract: `PASS-OUTCOME-SPEC-001`

## Context

PASS-SPEC-001 section 7.2 required native Gitleaks Git mode for exact history
and range scanning. PSCAN-03 proved that the pinned native patch stream omitted
deleted binary object bytes. Correction C1 enumerated exact objects but was
independently rejected because a prefix marker did not prove later fragment
inspection, the rule pack included unbounded matches, framing masked raw archive
identity, resource limits were not bound, binary proof was not isolated, and
directory-mode projection conflicted with the controlling mechanism.

The owner abandoned the contract-preserving option and activated one bounded
governance transition. That decision changes only the mechanism authority needed
to permit a coverage-proved implementation. It does not waive safety, legal,
product-independence, redaction, offline, licensing, spending or owner gates.

## Decision

Adopt PASS-OUTCOME-SPEC-001 as the single controlling successor after PSCAN-09
is independently accepted and closed. Native Gitleaks Git mode is no longer
inherently authoritative. Exact Git-object enumeration and deterministic
projections may be authoritative only when exact range/object/tree identities,
every admitted byte, raw pre-transformation classification, pinned detector
span/stream behavior, bounded declared profiles and class-isolated adversarial
proof are all complete. Any unsupported or unproved condition is non-pass.

Gitleaks remains the primary detector and the runner remains a thin MIT-licensed
Go program. TruffleHog remains absent and inactive. PASS-SPEC-001 and all prior
records remain immutable history; PSCAN-03 and correction C1 remain rejected.

## Alternatives considered

### Retain native Git-mode authority

Rejected by the owner's recorded direction. It preserves wording but cannot
prove deleted binary history reaches the pinned detector and therefore cannot
meet complete history coverage.

### Accept correction C1 as sufficient

Rejected. Exact object enumeration alone does not prove detector inspection.
The independent false-pass reproduction and other blockers remain valid.

### Enable a fallback detector now

Rejected as premature and outside PSCAN-09. No accepted successor coverage map
has established the residual required-class gap. TruffleHog technical and AGPL
owner gates remain separate and no assessment is authorized.

### Replace Gitleaks or reduce supported classes/profiles

Not authorized. Either would broaden the owner decision or weaken a retained
product requirement. Such a decision must return to the owner.

## Weaknesses and consequences

- Exact-object preparation and inspection proofs add implementation and review
  complexity.
- A pinned detector or rule change can invalidate span/stream proof and requires
  fresh class-isolated acceptance.
- Unbounded whole-match semantics can make finite projections impossible; work
  must fail closed until bounded rules or a proved streaming design exists.
- Complete raw classification and declared maximum profiles may expose a
  material primary-engine gap. That is evidence for an owner decision, not
  permission to enable fallback or shrink claims silently.
- PSCAN-03 implementation cannot be reused as accepted authority. A new bounded
  task must implement and prove the successor design.

## Exact limits of the owner decision

This decision authorizes PSCAN-09 documentation/governance transition only. It
does not authorize scanner code, dependencies, rules, fixtures, tests,
workflows, binaries, downloads, tool execution, remote action, publication,
signing, settings, credential handling, spending, TruffleHog work or consuming-
project integration. It selects or activates no successor.
