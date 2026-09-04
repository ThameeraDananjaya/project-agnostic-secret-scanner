# Project-Agnostic Secret Scanner

This repository is the local product workspace for an independently versioned,
offline secret-scanning package. The product is governed by
PASS-OUTCOME-SPEC-001.

Current state: **PSCAN-06 Correction C2 is claimed alone and implemented in the
current bounded local candidate. Author validation is not independent
acceptance; actual fresh-runner proof and every remote/signing gate remain
open. PSCAN-07 is unselected and PSCAN-08 remains inactive.**

The public repository and locked product tag `v1.0.0` exist, but the authorized
release workflow failed closed before building and no signed release exists.
There is no consuming-project integration. A passing scanner result never
authorizes merge, deployment, production use, legal compliance, or go-live.

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

The product is a thin MIT-licensed Go runner with Gitleaks CLI as the
primary engine. Authoritative scanning will execute locally with no network or
credentials, treat candidate material as data, retain no findings, and return
only content-free outcomes. Projects remain the authority for their policies,
allowlists, receipts, keys, evidence custody, deployment gates, and retention.

PSCAN-09 established the controlling outcome contract; PSCAN-10 supplied the
accepted exact-object primary coverage; PSCAN-04 supplied accepted internal
artifact normalization; and PSCAN-05 supplied accepted policy/reference
verification. PSCAN-06 remains open until its correction is independently
accepted and separately owner-gated remote signing/release work succeeds. No
credential, long-lived signing key, signed release or project integration is
present.
