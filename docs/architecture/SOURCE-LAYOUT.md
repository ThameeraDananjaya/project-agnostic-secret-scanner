# Proposed Source Layout

## Status

This is the complete target layout. Directories forbidden to PSCAN-01 are not
created by this task. Their future creation is authorized only by the exact
activated task named in the ownership table.

```text
/.github/workflows/             trusted CI and release workflows
/cmd/scanner-runner/            public CLI entry point
/internal/request/              request and schema validation
/internal/workspace/            isolated workspace lifecycle
/internal/gitinput/             safe Git-range and tree preparation
/internal/artifact/             archive and OCI normalization
/internal/engine/               engine interface and orchestration
/internal/engine/gitleaks/      primary adapter
/internal/engine/trufflehog/    disabled fallback adapter boundary
/internal/policy/               precedence and exception evaluation
/internal/redaction/            mandatory output firewall
/internal/outcome/              states, reason codes and exit mapping
/internal/cleanup/              transient-material destruction
/internal/verify/               release and evidence verification
/contracts/scan-request/        scanner-owned request schema
/contracts/scan-outcome/        scanner-owned outcome schema
/contracts/release-manifest/    scanner-owned release-manifest schema
/contracts/global-revocation/   scanner-owned global-revocation schema
/rules/generic/                 generic scanner rule pack only
/fixtures/clean/                known-safe synthetic input
/fixtures/synthetic-findings/   non-credential synthetic detections
/fixtures/history/              generated Git-history fixtures
/fixtures/archives/             generated safe/adversarial archives
/fixtures/containers/           generated OCI/Docker fixtures
/fixtures/adversarial/          injection, isolation and failure fixtures
/tests/unit/                    component-level tests
/tests/integration/             bounded component integration tests
/tests/acceptance/              complete product acceptance harness
/build/                         reproducible build and packaging definitions
/licenses/                      bundled release licence evidence
/docs/                          governance, architecture, security and runbooks
/evidence/                      task acceptance evidence; never findings
```

## Task ownership

| Boundary | First task allowed to create or materially implement it |
|---|---|
| Governance, architecture, threat, validation and task documents | PSCAN-01 |
| `cmd/`, request/outcome/workspace internals, initial scanner-owned schemas | PSCAN-02 |
| Gitleaks intake, engine orchestration, Git/file input and primary fixtures | PSCAN-03 |
| Artifact, redaction, cleanup and offline-sandbox implementation | PSCAN-04 |
| Policy projection, allowlist validation, compatibility and isolation | PSCAN-05 |
| Workflows, build, licences, SBOM, signing design and release verifier | PSCAN-06, subject to action-time owner gates |
| Full acceptance, parity, performance, lifecycle and v1 release evidence | PSCAN-07, subject to publication/signing gates |
| TruffleHog adapter/binary/fallback fixtures | PSCAN-08 only after material-gap evidence and separate technical/legal approval |

## Permanent exclusions

No path may contain consuming-project source, policy or allowlist instances,
receipt or revocation instances, signing or credential material, findings,
customer data, repository identities, project statistics, deployment behavior,
or cross-project state. Reference examples of project-owned formats must be
clearly non-authoritative and contain synthetic, non-identifying data only.
