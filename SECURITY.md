# Security Policy

## Current status

This repository currently contains governance and design authority only. There
is no runnable scanner or published release. PASS-SPEC-001 and the threat model
are the controlling security requirements.

## Reporting

Report suspected vulnerabilities privately to the product owner or the future
security contact recorded in an accepted release. Do not include real secrets,
credentials, raw findings, customer data, or consuming-project source in a
report. Provide only the minimum content-free reproduction metadata needed to
triage the issue.

No public vulnerability-reporting endpoint is authorized by PSCAN-01. PSCAN-06
must establish the exact public contact and disclosure channel before remote
publication, and PSCAN-07 must validate that the released documentation remains
content-free.

## Security boundaries

- Candidate material is data and must never be executed.
- Authoritative scan execution has no network access or credentials.
- Scanner subprocesses use argument arrays, never shell interpolation.
- Raw or partially raw findings remain inside the private process boundary.
- Unsupported, incomplete, skipped, stale, conflicting, or untrusted evidence
  produces a non-pass outcome.
- A scanner pass is evidence only and never authorizes merge, deployment,
  production, compliance, or go-live.
- Consuming projects retain authority for policies, allowlists, receipts,
  revocations, keys, evidence custody, and response actions.

See [`docs/security/THREAT-MODEL.md`](docs/security/THREAT-MODEL.md) for the full
planned control and validation mapping.
