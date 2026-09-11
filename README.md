# Project-Agnostic Secret Scanner

This repository is the local product workspace for an independently versioned,
offline secret-scanning package. The product is governed by
PASS-OUTCOME-SPEC-001.

Current state: **PSCAN-06 Correction C2 iteration 006 is independently accepted
locally for its bounded PID-collision correction, but PSCAN-06 remains open and
unaccepted overall. Recovery R5 is terminal after its sole protected C2-tag
push and workflow dispatch failed before Docker or artifact creation. A direct
Recovery R6 is invalid because the immutable C2 tag, workflow and schema bind
the pre-fix tooling commit. Correction C2 iteration 007 candidate `230e176` is
independently rejected for its over-broad workflow normalization. Iteration 008
stopped with no candidate after reproducing a pre-existing 66-line Go 1.27.1
formatting projection and Windows-only synthetic Cosign fixture failure.
Correction C2 iteration 009 was claimed and its two exact bounded work paths
remain unstaged and uncommitted. A diagnostic continuation proved an inherited
Windows case-collision fixture is unrepresentable on the current case-
insensitive filesystems, while the test and production normalizer bytes match
the clean materialization. The owner-approved iteration 009 continuation adds
only `tests/unit/artifact/normalize_test.go` for a test-fixture-only repair.
The genuinely fresh continuation implemented that repair, preserved the
carried integration patch byte-for-byte, and passed the complete pinned
offline author matrix including Linux/amd64 compile-only coverage. This
evidence-bearing candidate is author-validated but not accepted; genuinely
fresh independent review remains required.
Recovery R6 remains unauthorized. PSCAN-07 is unselected and PSCAN-08 remains
inactive.**

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
verification. PSCAN-06 remains open while the append-only R6 identity,
actual-Linux/Docker build proof and separately owner-gated signing/release work
remain incomplete. No credential, long-lived signing key, signed release or
project integration is present.
