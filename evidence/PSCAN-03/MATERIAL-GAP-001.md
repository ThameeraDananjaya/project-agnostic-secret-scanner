# PSCAN-03 Material Coverage Gap 001

## Decision

PSCAN-03 cannot be accepted. The exact pinned Gitleaks v8.30.1 Git-history
source path does not receive binary blob contents from an ordinary `git log -p`
range. This is a material required-history coverage gap. Every scan that
requires affected history must fail closed as
`INDETERMINATE_INCOMPLETE_COVERAGE`; it cannot produce `PASS`.

This evidence does not select, activate or authorize PSCAN-08. TruffleHog was
not assessed, downloaded, integrated or enabled.

## Exact engine binding

- Release: Gitleaks `v8.30.1`.
- Commit: `83d9cd684c87d95d656c1458ef04895a7f1cbd8e`.
- Source archive SHA-256:
  `6b2638a733b85619dc80bdf28e84e4fed7e526a761ab5c148fbf67695aea2115`.
- `sources/git.go` constructs `git -C <source> log -p -U0` plus the exact
  range supplied through `--log-opts`; it does not request or decode binary
  patch bodies.
- The upstream README independently describes Git scanning as `git log -p`.

## Reproduction

The generated test creates a clean base commit, adds an in-range file that
contains NUL bytes and a unique fabricated canary, commits its deletion, and
runs the same `git log -p -U0 <base>..<head>` input form used by the pinned
engine. The result contains Git's binary-difference marker but not the canary.
The candidate bytes therefore cannot reach Gitleaks's detector.

Test: `tests/unit/gitinput/git_test.go`,
`TestSafeBareCloneAndExactRangeBinding`.

Observed on 2026-08-31 with Git `2.55.0.windows.3`:

- exact three-commit fixture created;
- exact two-commit range and merge-base bound;
- ordered history and tracked-tree digests produced;
- binary-difference marker present;
- synthetic canary absent from the patch stream;
- test passed without publishing raw patch output.

## Materiality

PASS-SPEC-001 requires Git history, including deleted in-range content, to be
complete or explicitly non-pass. A candidate may place a credential-shaped
value in a binary-classified blob and delete it before the range head. Gitleaks
v8.30.1's native Git input will not inspect those bytes, so a clean engine
result cannot support the required history coverage claim.

The tracked-source directory mode does not cure deleted-history coverage:
deleted blobs are absent from the head tree. Archive/container normalization is
PSCAN-04 scope and was not used as an implied workaround.

## Other intake facts retained

- Exact Gitleaks config SHA-256:
  `e163e53b9e7e8a8511e77271e2b323ed057759542a6d988258afe3a1fa329caf`.
- Exact MIT licence SHA-256:
  `e3884b252b3bfc045e55be43a34d1e80da070bc6f804ac95bf4660e97d62ebc6`.
- Two clean Go 1.27.0 builds were byte-identical per platform:
  Windows amd64
  `a9e923bdde0e353057f7b71c2b14f1e1b96016076f05fc81c93dd605f46525ea`;
  Linux amd64
  `c79361874b71d1b8a366773cc3cee1ade9159b1b0500e3456835fff065c7555a`.
- The upstream tag resolves directly to an unsigned commit. The repository
  therefore records digest-pinned source provenance, not signed-tag provenance.
- Native Windows execution remains blocked by this host's Application Control
  policy. Linux directory-mode smoke evidence exists, but it cannot repair the
  Git-history coverage gap.

## Required owner decision

Select a new discussion task to decide whether PSCAN-03 should be redesigned to
materialize and scan every exact in-range Git blob, or whether the accepted gap
should make PSCAN-08 eligible for separate technical and AGPL owner review.
Neither direction is selected here.
