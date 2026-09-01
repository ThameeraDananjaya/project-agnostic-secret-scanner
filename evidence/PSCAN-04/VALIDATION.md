# PSCAN-04 Executed Validation

## Bound identities

- Activation commit: `a21b030e1658f1f98ac4e4d001af12185d9ed311`.
- Controlling contract SHA-256:
  `8a034701867e366bf37adea6e4f47e3ffb7a53425463ceff4fb280bcb4296e74`.
- Go image: `golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`;
  `go version go1.27.0 linux/amd64`.
- Gitleaks `8.30.1` source commit:
  `83d9cd684c87d95d656c1458ef04895a7f1cbd8e`; mounted Linux amd64 binary
  SHA-256 `657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586`.
- Rule SHA-256:
  `cebfe007ae88a55540e42ccb8e042848d089fe2f59bcf604c5fd4904d6460792`;
  empty-ignore SHA-256:
  `3cdd737e34cb301ffc19861f903cf59441ff57453ddf70bd5b8c1bd63389aebe`.

The dated primary sources, links, licences, versions, commits, digests and
compatibility decisions are recorded in `PREFLIGHT.md`.

## Authoritative full validation

Both final commands used the exact read-only module cache and exact verified
Gitleaks binary:

```powershell
& .\build\sandbox\validate.ps1 -ModuleCache 'C:\Users\ITDan\AppData\Local\Temp\pscan-04-modcache-a21b030e' -GitleaksBinary 'C:\Users\ITDan\AppData\Local\Temp\pscan-10-live-preflight-20260901\build-a\gitleaks-linux-amd64' -Profile pr
& .\build\sandbox\validate.ps1 -ModuleCache 'C:\Users\ITDan\AppData\Local\Temp\pscan-04-modcache-a21b030e' -GitleaksBinary 'C:\Users\ITDan\AppData\Local\Temp\pscan-10-live-preflight-20260901\build-a\gitleaks-linux-amd64' -Profile release
```

Both passed `go test -count=1 ./...` and silent `go vet ./...` after the final
`/tmp` `noexec` correction. The PR profile asserted an exact 4 GiB memory
cgroup and zero swap; the release profile asserted an exact 8 GiB memory
cgroup and zero swap. Both asserted loopback-only networking, failed DNS,
public-TCP and provider-metadata canaries, read-only root/source/cache/detector,
no credential-like environment, zero capabilities, no-new-privileges, a
256-PID harness limit and the declared one-scanner-process boundary.

The first implementation run had moved Go build temporaries to executable
`/work` because legacy tests cannot execute from `noexec` storage. A final
consistency audit found `/tmp` lacked the separately documented `noexec` mount
option; this was corrected and both profiles were repeated successfully.

## Race and cross-compile checks

The following network-disabled pinned-container race command passed:

```text
go test -race -p=1 -count=1 ./tests/unit/artifact ./tests/unit/cleanup ./tests/unit/redaction ./tests/integration/artifact
```

Observed package results were artifact unit `1.193s`, cleanup unit `1.034s`,
redaction unit `2.070s` and artifact integration `11.781s`.

The final Windows amd64 compile-only invocation used `CGO_ENABLED=0`,
`GOOS=windows`, `GOARCH=amd64`, no network, a read-only root/source/module
cache, dropped capabilities, no-new-privileges, 4 GiB memory with zero swap,
and isolated tmpfs build/output caches. It compiled but did not execute:

| Package | Generated test-binary SHA-256 |
|---|---|
| `tests/unit/artifact` | `47c96bad751b52880f6088d021277ff5ae2905e21492a8abb524d807a2a0edad` |
| `tests/unit/cleanup` | `ec57238e8629e4a36301e47ce0d8316d17ff797e9241f89040e8e93808d99b99` |
| `tests/unit/workspace` | `003784656978d10ce470c8a51cbe460013053b7930df2030e04e65e40e8e3aed` |

The first compile-only invocation failed before compilation because the
read-only root had no Go build-cache mount. The corrected invocation added an
isolated `/gobuild` tmpfs and passed. Windows runtime tests remain unexecuted
because Windows Smart App Control blocks locally generated executables; this is
not counted as Windows runtime or cross-platform acceptance.

## Coverage and adversarial outcomes

- Actual pinned-Gitleaks clean and synthetic-finding cases passed separately
  for directory, ZIP, TAR, gzip, TGZ, OCI layout and current Docker 29 save.
- Exact ledger/projection tests verify every raw and expanded supported payload
  range, chunk continuity, overlap, per-payload digest and absence of missing,
  duplicate or extra detector input.
- Unsafe path/link/reparse/device/hard-link, duplicate/case collision,
  encrypted, malformed, unsupported, ambiguous and below/at/above resource
  cases fail closed. Aggregate OCI/Docker semantic budgets and descriptor
  conflicts are included.
- Raw, partial, Base64, URL-Base64, hexadecimal and SHA-256 redaction canaries
  were absent from pass, finding, error, timeout, cancellation, panic and
  writer-error surfaces. Duplicate/unknown keys and framing injection were
  rejected.
- Cleanup tests passed for pass, fail, indeterminate, timeout, cancellation,
  recovered panic and scanner-owned abandoned-workspace recovery; forged,
  unleased or structurally unsafe material is preserved rather than deleted.

All test inputs are synthetic. No candidate artifact was built, installed or
executed. No real project data or credential was handled.

## Repository checks and limitations

`git diff --check` passed. `go.mod`, `go.sum` and `THIRD_PARTY_NOTICES.md` are
unchanged. All changed paths are within the PSCAN-04 allowlist. The final
staged-path and clean-worktree checks are performed around the local closeout
commit because the commit cannot contain its own identity.

A non-mutating `gofmt -d` check over all 32 changed Go files passed inside the
same pinned, network-disabled Go image. The host command could not start because
`gofmt` is not exposed on the Windows host `PATH`; this is not counted as a
failed product or runtime test.

Acceptance is limited to the bounded internal PSCAN-04 component and the
documented strict OCI/Docker subset. It does not prove a public runner,
consumer, deployment, production, Windows-runtime or complete cross-platform
boundary. Hostile same-user filesystem mutation remains outside the documented
process boundary.
