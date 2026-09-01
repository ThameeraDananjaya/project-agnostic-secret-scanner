# PSCAN-10 Live Primary-Source Preflight

## Scope and time

- Task: `PSCAN-10` only.
- Performed: `2026-09-01`, before dependency change, source or binary download,
  rebuild, scanner execution or new pinning in this implementation session.
- Activation base: `9053b37d799b19e3d98aeca0ae853296971ab09d` on
  `main`; preflight-start worktree status was clean.
- Sources: only the upstream Gitleaks repository/release records and official Go
  project records listed below. The separately licensed Gitleaks GitHub Action
  was not used, assessed or admitted.

## Gitleaks release, source and binary boundary

- The upstream release page still identifies `v8.30.1` as the latest release
  and binds it to commit
  `83d9cd684c87d95d656c1458ef04895a7f1cbd8e`.
- The release page publishes checksum-file SHA-256
  `061476c21adaf5441516f96f185c1a4706a83cd6329b9b38762271b3d4a52fae`
  and Linux x64 archive SHA-256
  `551f6fc83ea457d62a0d98237cbad105af8d557003051f41f3e7ca7b3f2470eb`.
- Upstream issue `#2086` records that the `v8.30.1` tag is not an ancestor of
  the default branch. Digest identity, not branch ancestry or signed-tag
  provenance, is therefore required.
- Upstream issue `#2164` records that a Windows x64 release archive previously
  did not match the published checksum file. A fresh download in this session
  produced Windows x64 archive SHA-256
  `d29144deff3a68aa93ced33dddf84b7fdc26070add4aa0f4513094c8332afc4e`,
  which matches the freshly downloaded checksum file. The historical mismatch
  still means release-asset identity cannot be inferred from the tag or asset
  name and is not admitted as PSCAN-10 build provenance.
- Upstream issue `#2170` records a `v8.30.1` no-detection report for a Homebrew
  arm64 build. This is not accepted as proof about the source-built Windows or
  Linux amd64 binaries, but it makes behavioral detection canaries mandatory
  and prevents version output or a digest alone from supporting pass.
- Disposition: no upstream release binary is authoritative for PSCAN-10.
  PSCAN-10 may use only locally rebuilt Windows amd64 and Linux amd64 binaries
  from the exact source commit, with exact source/build-input digests,
  reproducibility evidence and class-isolated behavioral detection proof.

Fresh intake evidence, collected only after this preflight record existed:

- checksum file SHA-256:
  `061476c21adaf5441516f96f185c1a4706a83cd6329b9b38762271b3d4a52fae`;
- source archive SHA-256:
  `6b2638a733b85619dc80bdf28e84e4fed7e526a761ab5c148fbf67695aea2115`;
- Linux x64 release archive SHA-256:
  `551f6fc83ea457d62a0d98237cbad105af8d557003051f41f3e7ca7b3f2470eb`;
- Windows x64 release archive SHA-256:
  `d29144deff3a68aa93ced33dddf84b7fdc26070add4aa0f4513094c8332afc4e`;
- exact source-tree `go.mod` SHA-256:
  `607c140abf2a872e70423972d4dfc7fa658ebe10365d0ea995269ed292add7a3`;
- exact source-tree `config/gitleaks.toml` SHA-256:
  `e163e53b9e7e8a8511e77271e2b323ed057759542a6d988258afe3a1fa329caf`;
- exact source-tree `LICENSE` SHA-256:
  `e3884b252b3bfc045e55be43a34d1e80da070bc6f804ac95bf4660e97d62ebc6`.

Primary records:

- `https://github.com/gitleaks/gitleaks/releases/tag/v8.30.1`
- `https://github.com/gitleaks/gitleaks/issues/2086`
- `https://github.com/gitleaks/gitleaks/issues/2164`
- `https://github.com/gitleaks/gitleaks/issues/2170`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/LICENSE`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/go.mod`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/.goreleaser.yml`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/config/gitleaks.toml`

## Detector span, stream, archive and skip boundary

The exact pinned source establishes:

- `sources/file.go` uses a 100,000-byte base buffer.
- `sources/common.go` reads at most 25,000 additional bytes to seek a safe
  boundary and introduces no overlap with the next fragment.
- `sources/file.go` identifies archives before ordinary-file fragmentation,
  silently returns when archive depth is exceeded, and silently skips detected
  application/binary MIME types.
- `detect/detect.go` can skip complete paths through global allowlists and can
  skip fragments through `MaxTargetMegaBytes`; a clean report alone therefore
  cannot prove inspection.
- `sources/git.go` uses `git log -p -U0`, skips deletions and skips non-archive
  binary patch files; native Git mode cannot satisfy the exact-object outcome
  contract.
- the upstream release build enables the `gore2regex` tag; the source also
  provides a standard-library regex build. PSCAN-10 must bind the chosen build
  variant explicitly and prove its behavior.

Primary records:

- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/sources/file.go`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/sources/common.go`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/sources/git.go`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/detect/detect.go`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/regexp/wasilibs_regex.go`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/regexp/stdlib_regex.go`

## Configuration, rules and licence boundary

- The exact upstream default configuration contains 222 product rules and
  unbounded regular expressions. It also contains global path allowlists that
  skip binary and other required input classes. It cannot support a finite-
  overlap or every-byte pass claim unchanged.
- Any admitted PSCAN-10 rule pack must retain or strengthen the blocking
  semantics, remove pass-capable path skips, provide a machine-checked finite
  maximum match span for every blocking rule and bind the resulting digest.
- The exact Gitleaks source licence is MIT, copyright 2019 Zachary Rice. The
  separately licensed Gitleaks GitHub Action remains outside scope.

## Go toolchain and module boundary

- Official Go release history records Go `1.27.0`, released 2026-08-19, as the
  current major release on the preflight date.
- Official policy supports a major release until two newer major releases
  exist. Go 1.27 is therefore supported on the preflight date.
- Go 1.27 retains Go 1 compatibility and supports pure-Go static binaries with
  `CGO_ENABLED=0`.
- Gitleaks `v8.30.1` declares `go 1.24.11`; Go 1.27.0 satisfies that module
  requirement. Its exact module graph and all compiled-module licence files
  must be regenerated and digest-bound for the selected build variant.
- The Go toolchain is distributed under the official BSD-style licence. No Go
  module is added to the scanner runner; `go.mod` remains standard-library-only
  unless a separately evidenced need is established.

Fresh official download metadata binds the Windows amd64 Go 1.27.0 archive to
SHA-256 `f0c0a0d33ba94f4d2c5dbc887334ce678b21813504ddb3aafcb06e60a5a667c4`.
That archive was re-hashed, freshly extracted for this task and produced exact
version output `go version go1.27.0 windows/amd64`; the extracted `go.exe` has
SHA-256 `7d828191ba32519a9c9361789ab647486236ed45c660889196c7770a8ff1985c`.

Primary records:

- `https://go.dev/doc/devel/release`
- `https://go.dev/doc/go1.27`
- `https://go.dev/doc/toolchain`
- `https://go.dev/LICENSE`

## Fail-closed preflight result

`PROCEED_TO_BOUNDED_DESIGN`, subject to all of the following before claim:

1. the freshly rechecked source/archive and Go archive digests above remain
   bound through build and closeout evidence;
2. the selected source-built binaries remain bound to the byte-identical
   cross-build and detection-canary results in `IMPLEMENTATION.md`;
3. the admitted 224-rule pack remains machine-proved with zero unbounded rules
   and exact conservative maximum span 4,020 bytes;
4. raw archive/compression/container/ambiguous classes are classified before
   preparation and unsupported classes are explicit non-pass; and
5. no release-asset checksum, Homebrew behavior, marker, clean exit or
   aggregate finding is treated as proof for a supported input class.

No download, source intake, binary intake, build, scanner execution, dependency
change, remote action, credential handling, signing, spending, publication,
TruffleHog work or successor action occurred while producing this record.
