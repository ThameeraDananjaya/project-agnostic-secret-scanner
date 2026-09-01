# Project-Agnostic Secret Scanner

This repository is the local product workspace for an independently versioned,
offline secret-scanning package. The product is governed by
PASS-OUTCOME-SPEC-001.

Current state: **PSCAN-09 completed and independently accepted locally. PSCAN-10
and every other successor remain proposed/unselected; PSCAN-08 remains
inactive**.

There is no runnable scanner, release, remote repository, or consuming-project
integration yet. A passing scanner result will never authorize merge,
deployment, production use, legal compliance, or go-live.

## Current authority

- Controlling contract: [`docs/spec/PASS-OUTCOME-SPEC-001.md`](docs/spec/PASS-OUTCOME-SPEC-001.md)
- Historical predecessor: [`docs/spec/PASS-SPEC-001.md`](docs/spec/PASS-SPEC-001.md)
- Requirement traceability: [`docs/spec/TRACEABILITY.md`](docs/spec/TRACEABILITY.md)
- Architecture: [`docs/architecture/ARCHITECTURE.md`](docs/architecture/ARCHITECTURE.md)
- Threat model: [`docs/security/THREAT-MODEL.md`](docs/security/THREAT-MODEL.md)
- Validation plan: [`docs/validation/VALIDATION-PLAN.md`](docs/validation/VALIDATION-PLAN.md)
- Task tracker: [`docs/tasks/TRACKER.md`](docs/tasks/TRACKER.md)
- Owner gates: [`docs/governance/OWNER-GATES.md`](docs/governance/OWNER-GATES.md)

## Product boundary

The future product is a thin MIT-licensed Go runner with Gitleaks CLI as the
primary engine. Authoritative scanning will execute locally with no network or
credentials, treat candidate material as data, retain no findings, and return
only content-free outcomes. Projects remain the authority for their policies,
allowlists, receipts, keys, evidence custody, deployment gates, and retention.

The PSCAN-02 contract/workspace skeleton and rejected non-authoritative PSCAN-03
candidate exist in repository history; there is no accepted pass-capable
scanner, workflow, release, credential, signing key, remote repository or
project integration. Exact-object primary coverage is proposed only as
unselected PSCAN-10.
