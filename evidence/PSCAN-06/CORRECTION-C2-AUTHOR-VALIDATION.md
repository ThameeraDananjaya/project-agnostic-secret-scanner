# PSCAN-06 Correction C2 author validation

Date: 2026-09-04

Implementation candidate: `e7faf0f81b3853e2090c75378bfbca568b52efad`

Candidate tree: `2088da833bf32ea9af227408dd69a5ac191213b4`

Result: `PARTIAL_AUTHOR_VALIDATION_PASS`; actual-Linux container, Go, complete
offline build and independent acceptance remain open.

## Passing local evidence

- Exact clean candidate checkout: PASS.
- C2 path allowlist: PASS for all 24 implementation-commit paths.
- `git diff --check`: PASS.
- PowerShell direct parse/binding checks for each changed release script: PASS.
- `schema-2.1.json` JSON parsing and exact C2 schema/tag/ref constants: PASS.
- Workflow order: PASS. Host-only CRLF and cache proof precede the sole image
  admission/pull boundary; offline container proofs precede dependency
  acquisition; both builds follow acquisition.
- Pull isolation: PASS. Exactly one `docker pull` statement exists in the C2
  release implementation and it receives only the constant canonical digest
  reference. `acquire.ps1` is invoked without image-pull authority.
- Forbidden-control scan: PASS; no privileged mode, host networking, mode
  `0777`, mutable image tag or unpinned workflow action was added.
- `test-image-admission.ps1`: PASS. Exact canonical evidence passes; short-name,
  mutable-tag, wrong-repository, wrong-digest, absent and ambiguous evidence
  cases reject without Docker.
- `test-crlf-shell-payloads.ps1 -Phase HostOnly`: PASS from a separate clean
  local clone at the full candidate commit. An actual CRLF checkout was proved,
  raw shell carriage returns rejected, normalized shell carriage returns were
  absent, and Docker was not invoked.

The first host-only harness attempt supplied the abbreviated commit and failed
closed at the 40-character identity gate. The second used an outer checkout
created with ambient Windows EOL conversion and failed closed on raw
`.editorconfig` bytes. The successful retry used the full commit and an outer
canonical-LF checkout; the harness itself then created and proved its required
actual CRLF fixture. Neither failed setup attempt changed the candidate or
used network access.

## Preservation proof

| Object | Post-candidate identity |
|---|---|
| locked product tag `v1.0.0` | commit `a13c28fe7273bc8dc6545f97966a02889524eb4c`; tree `217b711ddea51fd0ea7e808edd2e27fdecef8427` |
| locked C1 tooling tag | commit `3fb7592889820fa2739a4a53588e073689621809`; tree `7de9dc4c5725bf38cc80aa734e3c2bac0abd9762` |
| release-manifest schema 1.0 | SHA-256 `4d3f68236127ee6e2a8b908df84f28e57cd1d0682c363d6ad4cb5cf86ba1129e` |
| release-manifest schema 1.1 | SHA-256 `7b89d12749424f9094b849a9d653da580eba14dbf84cb150eb824b203acf3973` |
| release-manifest schema 2.0 | SHA-256 `aa6ae235938048cb7c5a026c5d476a6ac8d8e55f5d8761bd1e499c8d79980267` |

## Explicitly unproved

- The author host has no Go executable, so `go test`, `go vet`, Windows cross-
  compilation and verifier execution were not run locally.
- The Docker daemon is unavailable and the exact image was not locally
  present. Per authority, no pull was attempted. Therefore actual container
  cache/shell proofs, dependency acquisition, two offline builds and byte-for-
  byte reproduction remain unproved in this author session.
- No actual-Linux UID/GID positive/wrong-owner/read-only execution was
  performed on this Windows host.
- No independent skeptical review or acceptance was performed.
- No remote read/write, network acquisition, Docker pull, tag, push, setting,
  workflow run, signing, attestation, draft, publication, paid capability or
  TruffleHog action was performed.

This evidence cannot accept PSCAN-06. It supports only an independent review
of the committed local C2 candidate. Every fresh-runner and remote gate remains
fail-closed.
