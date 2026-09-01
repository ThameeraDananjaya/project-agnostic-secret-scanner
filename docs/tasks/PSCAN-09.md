# PSCAN-09 Task Specification

## Title

Decommission PASS-SPEC-001 controlling authority and establish an outcome-based
successor contract.

## State

Completed and independently accepted locally on 2026-09-01. The implementation
began in a fresh session from exact activation commit
`f48698922c174411f05124b066772753a2b8dd1e`. No successor is selected,
activated, claimed or implemented.

## Objective

Preserve PASS-SPEC-001, its source record and its historical evidence as
immutable audit history; explicitly decommission its controlling authority;
establish one outcome-based successor contract that permits exact Git-object
enumeration and deterministic scan projections; and reconcile current
governance, architecture, security, validation, traceability and downstream
task planning without implementing scanner behavior.

## Authority transition boundary

- PASS-SPEC-001 remains the controlling product contract until PSCAN-09 is
  independently accepted and closed.
- `docs/spec/PASS-SPEC-001.md`, `docs/spec/PASS-SPEC-001-SOURCE.md` and the
  owner-provided external source remain byte-for-byte unchanged historical
  evidence with SHA-256
  `78A1AD9A9ABFDE577A50AE9AA467B5B82F733162C8187ED58A2BC300952B692D`.
- The successor contract must receive a distinct document identity and state
  exactly when and why it supersedes PASS-SPEC-001. It must not rewrite history
  or describe the rejected PSCAN-03 candidate as accepted.
- Every PASS-SPEC-001 capability, acceptance row, invariant, non-goal, schema
  boundary and owner gate must receive an explicit disposition in transition
  traceability. Mechanism-specific requirements may be replaced by outcome
  proofs. Any broader safety, product-boundary, legal or owner-gate reduction
  fails closed and returns to the owner rather than being inferred.
- The successor may permit exact Git-object enumeration and deterministic
  projections only when proof binds every admitted object and every byte to
  actual detector inspection. Prefix-only markers, unproved detector fragment
  spans, framing that masks raw input classification, unbound resource limits
  and aggregate tests that do not isolate the claimed class cannot prove pass.

## Preconditions

- PSCAN-03 is rejected, closed fail-closed and unaccepted at clean commit
  `7b866484da50319c91bd7b67c37637ceef888dab`.
- Exact owner activation of PSCAN-09 and a fresh implementation session from
  the exact activation commit.
- Canonical and external PASS-SPEC-001 bytes, line count and SHA-256 still match
  the source record before any authority transition.
- The PSCAN-03 rejection, material-gap evidence and correction review remain
  immutable and available.
- No other task is selected, activated or claimed.

## Allowed paths

```text
AGENTS.md
CONSTITUTION.md
README.md
SECURITY.md
docs/spec/PASS-OUTCOME-SPEC-001.md
docs/spec/TRACEABILITY.md
docs/decisions/DEC-002-OUTCOME-BASED-CONTRACT.md
docs/governance/TASK-LIFECYCLE.md
docs/governance/OWNER-GATES.md
docs/governance/EVIDENCE-AND-CLOSEOUT.md
docs/governance/RUNBOOK-REQUIREMENTS.md
docs/architecture/ARCHITECTURE.md
docs/architecture/GITLEAKS-ADAPTER.md
docs/architecture/SOURCE-LAYOUT.md
docs/security/THREAT-MODEL.md
docs/security/GITLEAKS-BOUNDARY.md
docs/validation/VALIDATION-PLAN.md
docs/validation/GITLEAKS-COVERAGE-MAP.md
docs/release/LICENSING.md
docs/tasks/PSCAN-04.md
docs/tasks/PSCAN-05.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-07.md
docs/tasks/PSCAN-08.md
docs/tasks/PSCAN-09.md
docs/tasks/PSCAN-10.md
docs/tasks/TRACKER.md
evidence/PSCAN-09/**
```

`docs/tasks/PSCAN-10.md` is reserved only for a proposed, unselected bounded
primary-coverage successor if the accepted transition analysis requires it.
Its inclusion does not select, activate or authorize PSCAN-10.

## Forbidden scope

- No change to `docs/spec/PASS-SPEC-001.md`,
  `docs/spec/PASS-SPEC-001-SOURCE.md`, the external owner source, PSCAN-01/02
  acceptance evidence, PSCAN-03 evidence, or PSCAN-01 through PSCAN-03 task
  specifications and reading maps.
- No product code, schemas, engine adapters, Git preparation, projections,
  rules, fixtures, tests, dependencies, binaries, build behavior or workflows.
- No scanner, dependency or tool download, build, execution or intake.
- No acceptance or rehabilitation of PSCAN-03 or correction C1.
- No TruffleHog assessment, legal work, download, integration or enablement.
- No consuming-project source, policy, allowlist, identity, receipt, key,
  finding, credential, customer data, statistics or project behavior.
- No remote creation, push, publication, release, signing, GitHub-setting
  change, credential handling, spending, deployment, production or go-live.
- No successor selection, activation, claim or implementation.

## Deliverables

- Immutable-history proof for PASS-SPEC-001, its source record, the external
  owner source and the rejected PSCAN-03 evidence chain.
- `PASS-OUTCOME-SPEC-001` as the single coherent successor contract, with an
  explicit authority transition, product boundary, safety invariants, outcome
  requirements, non-goals, acceptance definition and owner gates.
- DEC-002 recording the chosen outcome-based direction, rejected alternatives,
  weaknesses, transition rationale and exact limits of the owner decision.
- Transition traceability for every CAP-1 through CAP-20 item, AT-01 through
  AT-31 row and other governing requirement, including retained, reframed,
  deferred or explicitly owner-gated dispositions.
- A target primary-coverage architecture that admits exact Git-object
  enumeration and deterministic projections only with complete, byte-bound,
  class-isolated and adversarially proven detector coverage.
- Reconciled governance, security, validation and living documentation with no
  conflicting claim about which contract controls.
- A downstream planning sequence that removes accepted PSCAN-03 as a false
  prerequisite, keeps every successor proposed and unselected, and preserves
  PSCAN-08's separate material-gap, technical and AGPL owner gates.

## Acceptance

- The canonical, source-record and external PASS-SPEC-001 files remain exactly
  977 lines where applicable and retain the recorded SHA-256; historical task
  and evidence files are unchanged.
- Authority changes only through explicit accepted successor language; current
  living documents agree on one controller and historical documents remain
  clearly historical rather than silently rewritten.
- Transition traceability has no missing capability, acceptance row, invariant,
  non-goal, schema boundary or owner gate and does not hide a product reduction
  behind a mechanism change.
- The outcome contract resolves the PASS-SPEC-001 section 7.2 conflict without
  treating Gitleaks Git mode as inherently authoritative and without allowing
  incomplete detector coverage to pass.
- Architecture and validation require exact object/range/tree binding, raw
  input classification before transformation, proved detector span/stream
  behavior, bounded resources, class-isolated fixtures and non-pass for every
  unsupported or unproved condition.
- No scanner behavior, dependency, workflow, remote state, credential, signing,
  publication, spending or consuming-project material changes.
- Independent review proves complete path containment and records facts,
  limitations, interpretations, recommendations and owner gates separately.
- PSCAN-09 closes with every successor unselected and PSCAN-08 inactive unless
  separately made eligible and approved in a later owner-gated lifecycle.

Stop after PSCAN-09 closeout. Do not select or activate PSCAN-10, PSCAN-04 or
any other successor.

## Accepted closeout

PSCAN-09 established PASS-OUTCOME-SPEC-001 as the controlling successor through
DEC-002 and complete transition traceability. PASS-SPEC-001, its external source
and PSCAN-01 through PSCAN-03 records remain immutable history. PSCAN-03 and
correction C1 remain rejected. Exact implementation, independent review,
validation and closeout evidence is in `evidence/PSCAN-09`.

PSCAN-10 is proposed only. PSCAN-04 through PSCAN-07 remain proposed and
unselected. PSCAN-08 remains inactive, unselected, technically gated and AGPL
owner-gated. No successor activated automatically.
