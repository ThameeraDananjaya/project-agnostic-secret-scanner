# PSCAN-05 Task Specification

## Title

Policy precedence, allowlists, schema compatibility, receipt/revocation
reference verification, and project-isolation controls.

## State

Implemented as a bounded local candidate on 2026-09-02 from exact activation
commit `50b418609c1f9927c0ecd5d740aba6d7bff11f58`; awaiting independent
acceptance. No successor is selected or activated.

## Objective

Implement the strict projection boundary for project-owned policy and allowlist
inputs, enforce non-overridable precedence, provide custody-neutral receipt and
revocation reference verification, and prove that no state crosses projects.

## Preconditions

- PSCAN-04 accepted, closed and clean at activation parent
  `1f0890878518de32a55ceb8d7b97430c4f3d2f2b`.
- Exact PSCAN-05 owner selection and activation recorded in
  `evidence/PSCAN-05/ACTIVATION.md`, followed by a fresh implementation session
  from that activation commit.
- DEC-001 remains consistent with PASS-OUTCOME-SPEC-001 and its independent
  schema-family ownership boundary is preserved.

## Activation-only bundle

This session may change and commit only:

```text
docs/tasks/PSCAN-05.md
docs/tasks/PSCAN-05-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-05/ACTIVATION.md
```

No implementation, dependency intake, scanner/toolchain execution, test run or
successor work is authorized in the activation session.

## Allowed implementation paths

```text
go.mod
go.sum
internal/policy/**
internal/verify/**
internal/request/**
internal/outcome/**
internal/workspace/**
contracts/scan-request/**
contracts/scan-outcome/**
contracts/global-revocation/**
fixtures/adversarial/policy/**
fixtures/adversarial/allowlist/**
fixtures/adversarial/evidence/**
fixtures/adversarial/isolation/**
tests/unit/policy/**
tests/unit/verify/**
tests/integration/policy/**
tests/integration/isolation/**
tests/acceptance/policy/**
docs/architecture/**
docs/security/**
docs/validation/**
docs/decisions/DEC-001-SCHEMA-OWNERSHIP.md
docs/tasks/PSCAN-05.md
docs/tasks/TRACKER.md
evidence/PSCAN-05/**
```

## Forbidden scope

- No modification of PASS-OUTCOME-SPEC-001, PASS-SPEC-001, their source record,
  TRACEABILITY, DEC-002 or PSCAN-01 through PSCAN-04 and PSCAN-09/10 historical
  evidence.
- No consuming-project policy, allowlist, receipt or revocation instance, key,
  identity, source, repository, finding, credential, customer data, statistics
  or project-specific behavior.
- No project schema authority; no receipt issuance or signing; no receipt key
  custody; no evidence store; no promotion, deployment or retention authority.
- No TruffleHog assessment, legal work, download, material, integration,
  distribution or enablement.
- No GitHub workflow, remote creation, push, publication, release, signing,
  settings change, credential handling, spending, provider, deployment,
  production or go-live action.
- No selection, activation, claim or implementation of PSCAN-06, PSCAN-08 or
  any other successor.

## Deliverables

- Fixed precedence: scanner invariants, global revocations/protections, project-
  policy projection, then narrow allowlist exceptions.
- Digest/version/signature-aware projection with unknown-major and ambiguous-
  semantic rejection.
- Exception owner/approver/evidence/scope/create/expiry/invalidation rules,
  30-day maximum and mandatory credential-class rejection.
- Independent schema-family compatibility, overlap and retirement handling.
- Custody-neutral reference verifier for immutable receipts, append-only global
  and project revocations, bindings, 30-day staging/production reuse and chain
  rollback/divergence/conflict rejection.
- Parallel, sequential and reused-host isolation for workspaces, caches,
  correlation IDs, policy and evidence.

## Acceptance

No lower layer weakens a higher layer; invalid/expired/broadened exceptions fail;
real-looking credentials cannot be allowlisted; unknown majors and bad bindings
fail closed; any receipt binding change rejects reuse; revocation always wins;
rollback and divergent evidence block; two projects share no state. The scanner
never stores project instances, handles receipt keys or grants promotion.

Stop after PSCAN-05 closeout with PSCAN-06 unselected.
