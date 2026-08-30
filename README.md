# Project-Agnostic Secret Scanner

This repository is the local product workspace for an independently versioned,
offline secret-scanning package. The product is governed by PASS-SPEC-001.

Current state: **PSCAN-01 implementation prepared locally; acceptance closeout
is pending. Every successor remains unselected and PSCAN-08 remains inactive**.

There is no runnable scanner, release, remote repository, or consuming-project
integration yet. A passing scanner result will never authorize merge,
deployment, production use, legal compliance, or go-live.

## Current authority

- Canonical contract: [`docs/spec/PASS-SPEC-001.md`](docs/spec/PASS-SPEC-001.md)
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

No product code, schema, workflow, scanner binary, dependency, credential,
signing key, public release, remote repository, or project integration exists
in this repository yet.
