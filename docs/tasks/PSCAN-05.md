# PSCAN-05 Task Specification

## Title

Policy precedence, allowlists, schema compatibility, receipt/revocation
reference verification, and project-isolation controls.

## State

Proposed. Unselected. Not activated. Not claimed.

## Objective

Implement the strict projection boundary for project-owned policy and allowlist
inputs, enforce non-overridable precedence, provide custody-neutral receipt and
revocation reference verification, and prove that no state crosses projects.

## Preconditions

- PSCAN-04 accepted and clean.
- Exact PSCAN-05 activation and fresh implementation session.
- DEC-001 remains consistent with PASS-OUTCOME-SPEC-001.

## Allowed paths

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

No project policy/allowlist/receipt instances, keys, identities, repositories or
statistics; no project schema authority; no receipt issuance/signing; no evidence
store; no deployment gate; no workflow, remote, release, spending or TruffleHog.

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
