# PSCAN-06 Implementation Evidence

## Claim

- Task: `PSCAN-06` only.
- Session type: fresh implementation session, not activation.
- Claimed UTC: `2026-09-03T00:15:00Z`.
- Claim baseline: exact clean activation commit
  `058ffcd446c6431b2e1afeed769d02c7b1f307f8` on `main`.
- Controlling authority: `PASS-OUTCOME-SPEC-001`.
- Reading-map completion: all 60 listed authorities read completely before
  claim.
- Allowed paths: exactly the PSCAN-06 allowed implementation paths.
- Remote, settings, workflow, signing, credential, spending and publication
  gates: reserved; no action performed.
- TruffleHog: not downloaded, assessed, integrated, distributed or enabled.
- Successors: none selected, activated, claimed or implemented.

This record claims PSCAN-06 and no other task. Implementation evidence,
commands, deterministic artifact identities, tests and limitations are appended
only after execution. Missing evidence fails closed.

## Local implementation candidate

- Added release-manifest schema 1.1 with exact Windows/Linux runner and engine
  bindings, rule/schema bindings, fixed repository owner/workflow/tag/OIDC
  identity, compatibility, revocation and a complete asset inventory.
- Added a strict duplicate/unknown-member rejecting release verifier. It
  verifies exact manifest/bundle bytes, pinned Cosign executable bytes, fixed
  identity policy, path confinement, regular-file identity before/after reads,
  every asset length/digest, no unexpected files, and structural initial
  revocation evidence.
- Added deterministic tar.gz/ZIP and SPDX 2.3 generators plus two-phase
  acquisition/build scripts. Network is available only for exact digest-checked
  acquisition; builds use `--network none`, a read-only module cache, exact Go
  archives and the full-digest container image.
- Added a manual exact-tag workflow with four full-SHA action pins. Build jobs
  have `contents: read`; signing/draft permissions exist only in the separately
  gated job. The workflow has no publication command.
- Added full release verification and scanner input/output runbooks with Mermaid
  acquisition, release, execution, state, revocation, rollback and retirement
  diagrams.
- Added integration and synthetic adversarial tests for duplicate/unknown
  schemas and one-byte runner, licence, SBOM, revocation and bundle changes plus
  wrong identity, issuer and tag.

Pre-candidate validation passed on Go 1.27.1 in the pinned network-disabled
container: `go test -count=1 ./...`, `go vet ./...`, and compile-only
`GOOS=windows GOARCH=amd64 go test -exec /bin/true ./...`. The full clean-tree
two-build reproduction and exact artifact identities are recorded after this
candidate is committed.

No remote, settings, workflow, signing, credential, spend, publication,
TruffleHog or successor action occurred.
