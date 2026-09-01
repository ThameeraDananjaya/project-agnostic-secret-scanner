# PSCAN-04 Live Primary-Source Preflight

## Scope and time

- Task: `PSCAN-04` only.
- Verified: `2026-09-02`, before dependency change, download, build, scanner
  execution or new pinning in this implementation session.
- Activation commit: `a21b030e1658f1f98ac4e4d001af12185d9ed311` on
  `main`; its parent is `14348b16d9ddb87f337391a4e47dd243f6e53bf2` and
  its tree is `f7846bf906eb28176ea8e69e0ee804fbf4f20e0b`.
- Start state: exact activation HEAD, clean worktree, strict Git object check
  passed, and the living tracker identifies PSCAN-04 as the sole activated
  task. PSCAN-05 and every other successor remain unselected and inactive.
- Controlling contract: `docs/spec/PASS-OUTCOME-SPEC-001.md`, SHA-256
  `8a034701867e366bf37adea6e4f47e3ffb7a53425463ceff4fb280bcb4296e74`.

## Archive classification and parsing boundary

The existing module boundary is retained without change. No new archive or OCI
runtime dependency is admitted by PSCAN-04.

- `github.com/mholt/archives` `v0.1.2`, exact release commit
  `39a94fddf9854a23dfad56d54ae5d5e70b90b885`, MIT licence. Existing module
  checksum: `h1:UBSe5NfYKHI1sy+S5dJsEsG9jsKKk8NJA4HCC+xTI4A=`; existing
  `go.mod` checksum: `h1:D7QzB6TIGZ8ALJU2bOl4zYeUVe0Ltm0x8ceT9iY8M+I=`.
- `github.com/h2non/filetype` `v1.1.3`, exact release commit
  `a977222ef40641c96b13eb033a9fb190089f6f46`, MIT licence. Existing module
  checksum: `h1:FKkx9QbD7HR/zjK1Ia5XiBsq9zdLi5Kf3zGyFTAFkGg=`; existing
  `go.mod` checksum: `h1:319b3zT68BvV+WRj7cwy856M2ehB3HqNOt6sy1HndBY=`.
- Go `archive/zip`, `archive/tar`, `compress/gzip`, `encoding/json` and
  `crypto/sha256` from the exact Go toolchain below are the selected parsers.
  This avoids admitting the broader encrypted and unsupported formats exposed
  by the generic archive library. Raw classification remains before any
  transformation and every unsupported class is non-pass.

Primary records:

- `https://github.com/mholt/archives/releases/tag/v0.1.2`
- `https://github.com/mholt/archives/commit/39a94fddf9854a23dfad56d54ae5d5e70b90b885`
- `https://raw.githubusercontent.com/mholt/archives/39a94fddf9854a23dfad56d54ae5d5e70b90b885/LICENSE`
- `https://github.com/h2non/filetype/releases/tag/v1.1.3`
- `https://github.com/h2non/filetype/commit/a977222ef40641c96b13eb033a9fb190089f6f46`
- `https://raw.githubusercontent.com/h2non/filetype/a977222ef40641c96b13eb033a9fb190089f6f46/LICENSE`
- `https://pkg.go.dev/archive/zip@go1.27.0`
- `https://pkg.go.dev/archive/tar@go1.27.0`
- `https://pkg.go.dev/compress/gzip@go1.27.0`

## OCI and Docker-save format boundary

- OCI Image Format `v1.1.1`, signed release/tag and exact signed commit
  `147f9c13cedb47a0c4d9a11a222961073d585877`, Apache-2.0 licence, is the
  selected OCI-layout authority. The implementation consumes no OCI project
  code. It validates `oci-layout`, `index.json`, manifest/config descriptors,
  `blobs/<algorithm>/<encoded>` paths and every referenced digest before
  admitting layer bytes.
- Docker-save input is admitted only as a strict tar container whose selected
  manifest, configuration and layers are present, content-bound and accounted
  for. Both legacy digest-JSON/tar-layer saves and Docker 29's content-addressed
  OCI-backed save form are bounded. The official `docker image save`
  documentation and live Docker 29.7.2 output are the format-source boundary;
  malformed or non-content-addressed extra material is non-pass.

Primary records:

- `https://github.com/opencontainers/image-spec/releases/tag/v1.1.1`
- `https://github.com/opencontainers/image-spec/commit/147f9c13cedb47a0c4d9a11a222961073d585877`
- `https://raw.githubusercontent.com/opencontainers/image-spec/147f9c13cedb47a0c4d9a11a222961073d585877/LICENSE`
- `https://raw.githubusercontent.com/opencontainers/image-spec/147f9c13cedb47a0c4d9a11a222961073d585877/image-layout.md`
- `https://raw.githubusercontent.com/opencontainers/image-spec/147f9c13cedb47a0c4d9a11a222961073d585877/manifest.md`
- `https://docs.docker.com/reference/cli/docker/image/save/`

## Gitleaks detector boundary

- Gitleaks `v8.30.1`, exact source commit
  `83d9cd684c87d95d656c1458ef04895a7f1cbd8e`, MIT licence.
- Exact source archive SHA-256:
  `6b2638a733b85619dc80bdf28e84e4fed7e526a761ab5c148fbf67695aea2115`.
- Exact source `go.mod` SHA-256:
  `607c140abf2a872e70423972d4dfc7fa658ebe10365d0ea995269ed292add7a3`.
- Exact upstream configuration SHA-256:
  `e163e53b9e7e8a8511e77271e2b323ed057759542a6d988258afe3a1fa329caf`.
- Exact source licence SHA-256:
  `e3884b252b3bfc045e55be43a34d1e80da070bc6f804ac95bf4660e97d62ebc6`.
- Admitted 224-rule product configuration SHA-256:
  `cebfe007ae88a55540e42ccb8e042848d089fe2f59bcf604c5fd4904d6460792`;
  conservative maximum match span 4,020 bytes. Empty ignore-file SHA-256:
  `3cdd737e34cb301ffc19861f903cf59441ff57453ddf70bd5b8c1bd63389aebe`.
- Selected exact source-built outputs: Windows amd64 SHA-256
  `b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178`;
  Linux amd64 SHA-256
  `657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586`.
- Exact source inspection reconfirmed the 100,000-byte base fragment,
  additional 25,000-byte boundary read, no fragment overlap, archive depth and
  size skips, MIME skips and allowlist skips. A clean detector result is never
  sufficient; PSCAN-04 must supply deterministic projection coverage evidence.

Primary records:

- `https://github.com/gitleaks/gitleaks/releases/tag/v8.30.1`
- `https://github.com/gitleaks/gitleaks/tree/83d9cd684c87d95d656c1458ef04895a7f1cbd8e`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/LICENSE`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/sources/file.go`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/sources/common.go`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/detect/detect.go`

## Go toolchain boundary

- Go `1.27.0`, released `2026-08-19`, is current and supported on the
  verification date. The repository remains bound to `go 1.27.0`.
- Source archive SHA-256:
  `7002401cd36acaadf3e05c43a94e6aac63f50ac1250ddad6c94b40e797364b0e`.
- Linux amd64 archive SHA-256:
  `675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685`.
- Windows amd64 archive SHA-256:
  `f0c0a0d33ba94f4d2c5dbc887334ce678b21813504ddb3aafcb06e60a5a667c4`.
- Licence: official Go BSD-style licence. Go 1 compatibility applies; the
  selected source-built Gitleaks module declares Go 1.24.11 and is compatible
  with this toolchain.

Primary records:

- `https://go.dev/dl/`
- `https://go.dev/doc/devel/release`
- `https://go.dev/doc/go1.27`
- `https://go.dev/LICENSE`

## Authoritative Docker execution boundary

- Pinned test image: `golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`.
- Required controls: `--network none`; `--read-only`; source and detector bind
  mounts `ro`; writable bounded `tmpfs`; `--cap-drop ALL`;
  `--security-opt no-new-privileges`; default seccomp retained;
  `--pids-limit`; `--memory`; and `--cpus`.
- Only explicitly declared non-secret variables may enter the container. No
  credential directory, Docker socket, host environment, secret store or
  provider configuration may be mounted or inherited.
- The official none-network driver documentation states that only loopback is
  available. The run reference and storage documentation establish the
  read-only root and bind-mount controls. Resource and seccomp documentation
  establish the process, memory, CPU and syscall boundaries. These controls
  still require live inspect and DNS/TCP/provider canaries before acceptance.

Primary records:

- `https://docs.docker.com/engine/network/drivers/none/`
- `https://docs.docker.com/reference/cli/docker/container/run/`
- `https://docs.docker.com/engine/storage/bind-mounts/`
- `https://docs.docker.com/engine/containers/resource_constraints/`
- `https://docs.docker.com/engine/security/seccomp/`

## Fail-closed result

`PROCEED_TO_PSCAN_04_ONLY`.

The selected parsers, identities, licences, compatibility boundary and Docker
controls are complete and mutually consistent. PSCAN-04 may now be claimed and
implemented without changing dependencies. Acceptance remains prohibited until
the pinned Docker image identity, engine controls, detector identity, every-byte
coverage, redaction, isolation and cleanup behavior are proved live. Any
identity mismatch, missing reference, unexpected input class, uncontrolled
output, or incomplete lifecycle proof is non-pass.

No download, dependency change, build, scanner execution, new pin, remote
action, credential handling, signing, spending, publication, TruffleHog work or
successor action occurred before this record was created.

## Post-preflight live runtime verification

After the proceed decision, read-only checks recorded Docker client/server
`29.7.2`, Docker Desktop `4.87.0`, Engine commit `6a43e3d`, containerd `v2.2.5`
commit `e53c7c1516c3b2bff98eb76f1f4117477e6f4e66`, runc `1.3.6` commit
`491b69ba`, and Linux amd64 engine/runtime identity. The selected image inspect
returned only the exact pinned digest and `amd64/linux`.

An official-image format probe used only the already-pinned Go image. Docker
29.7.2 emitted a 311,128,576-byte save with seven selected layers, a seven-entry
`rootfs.diff_ids` configuration, `amd64/linux`, `rootfs.type=layers`, and the
content-addressed configuration path
`blobs/sha256/037ded92fa0bffa00c184cd9f5c02e2f60f59a47c2fdc14a632a79eabde14cd8`.
The probe archive was removed immediately after structural inspection; it was
not scanned, admitted as a synthetic fixture or retained.

The task-specific Go module cache was populated in a separate acquisition
container after this preflight, using the existing `go.sum`, Go checksum
verification, `proxy.golang.org` and `sum.golang.org`. Every authoritative build
then used `GOPROXY=off`, `GOSUMDB=off`, `--network none` and that cache mounted
read-only.
