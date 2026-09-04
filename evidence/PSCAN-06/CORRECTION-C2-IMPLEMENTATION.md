# PSCAN-06 Correction C2 implementation

Date: 2026-09-04

Authority base: `d4eca19e04862d660ca6ac0e9b64eec4fb06b61c`

State: bounded local author candidate; not independently accepted

## Implemented boundary

- Added a host cache canary that binds the invoking Linux numeric UID/GID and
  proves exact write, same-filesystem rename, read and cleanup without Docker.
- Split CRLF proof into a host-only phase and post-admission offline container
  phases. The host-only phase rejects raw carriage returns and proves the
  exact LF projection without inspecting or running an image.
- Added a dedicated image-admission boundary for only
  `docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`.
  It proves all exact acquisition cache paths first, permits a pull only when
  explicitly enabled, and requires post-admission Docker `RepoDigests` to bind
  the canonical repository and digest.
- Kept container cache and shell-parser canaries at `--pull=never --network
  none`; dependency acquisition remains isolated to the acquisition phase;
  both release builds remain `--network none` with a read-only completed cache.
- Added release-manifest schema `2.1` for proposed C2 tooling tag/ref/workflow/
  SHA/trigger identity. Verifier code preserves exact schema `2.0` C1 parsing
  and rejects cross-version masquerade and role swaps.
- Updated the release plan, validation contract, decision amendment, workflow
  gates and current task state without changing historical evidence.

## Deliberately unchanged

- locked product tag `v1.0.0` and locked C1 tooling tag;
- release-manifest schemas `1.0`, `1.1` and `2.0` byte-for-byte;
- accepted scanner, Gitleaks, rule, source-trust, sandbox, licensing, SBOM,
  revocation and exact-entrypoint controls;
- zero-spend, no-credential, no-TruffleHog, synthetic-only and fail-closed
  boundaries.

No remote or network action was performed. Local author validation and its
limitations are recorded separately and cannot accept PSCAN-06 or authorize a
C2 tag, pull, push, workflow run, signing, attestation, draft or publication.
