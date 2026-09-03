# PSCAN-06 Correction C1 author validation 003

## Result

- Superseding candidate commit:
  `5ca77226ed3996a8267caf02da366d0beb915c8d`
- Candidate tree: `e65ac8ca9b1c5dc9a8a70e26c122aee3e84d32d2`
- Authority commit:
  `3fb1b0a55dc4f48dd35464c63c768f497efbc89b`
- Result: `PARTIAL_AUTHOR_VALIDATION_PASS; ACTUAL_LINUX_HOST_PROOF_OPEN`
- Independent acceptance: not performed and not claimed
- Remote/signing/release state: not performed and not authorized

Candidate `f24b832ebe5f6749aa0ab910e1b9be279065ebb1` is rejected because
independent review reproduced CRLF corruption in its Docker POSIX shell
payloads. `CORRECTION-C1-ITERATION-003.md` preserves that rejection and the
historical preflight repository-name contradiction without rewriting either
earlier record. This validation applies only to candidate `5ca7722`.

## CRLF checkout regression

Actual command from clean candidate `5ca7722`:

```powershell
pwsh -NoProfile -File build/release/test-crlf-shell-payloads.ps1
```

Actual result: exit `0` with:

```text
CRLF shell-payload regression PASS raw-CR=REJECT normalized-CR=ABSENT canary=PASS acquisition-syntax=PASS build-syntax=PASS read-only-cache=REJECT
```

The test creates CRLF copies of `cache-canary.ps1`, `acquire.ps1` and
`build.ps1`, proves each extracted multiline payload contains carriage
returns, proves an unnormalized payload fails the LF assertion, converts it,
asserts no carriage return remains, and passes the exact normalized payload to
the full-digest pinned image with networking disabled for POSIX shell parsing.
It then executes the actual cache-canary function from the CRLF fixture. The
writable canary passes; a read-only bind rejects at its first write and leaves
no canary residue or acquisition ledger. Every Docker shell invocation is
covered, including the one-line deterministic packaging payload.

Runtime normalization converts CRLF and lone carriage returns to LF and then
rejects any remaining carriage return or NUL immediately before Docker is
invoked. The accepted read-only root, dropped capabilities,
no-new-privileges, resource bounds, network boundaries and embedded shell
commands are unchanged. The recovery workflow now runs this regression before
Linux cache ownership tests or acquisition.

## Fresh acquisition

Actual command:

```powershell
./build/release/acquire.ps1 -CacheDirectory "$env:TEMP\pscan-06-c1-acquisition-iteration003-5ca7722" -AllowImagePull
```

Actual result: exit `0`. All three public input digests passed and the
acquisition payload completed inside image
`golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`.
The newly created schema-`2.0` ledger has SHA-256
`db77fffe482389e06f5afbe1bc8fa3e72a54e047546655da8fc94fb48568bb99`,
records `write-atomic-rename-read-delete`, and records
`completedBeforeNetworkDependencyAcquisition: true`. It contains exactly the
three pinned archives and a newly populated module cache. The author-host mode
is `windows-docker-desktop-default-user`; it is not Linux UID/GID evidence.

## Independent network-disabled builds

Actual commands from the clean exact candidate:

```powershell
./build/release/build.ps1 -AcquisitionDirectory "$env:TEMP\pscan-06-c1-acquisition-iteration003-5ca7722" -OutputDirectory "$env:TEMP\pscan-06-c1-build-iteration003-a-5ca7722"
./build/release/build.ps1 -AcquisitionDirectory "$env:TEMP\pscan-06-c1-acquisition-iteration003-5ca7722" -OutputDirectory "$env:TEMP\pscan-06-c1-build-iteration003-b-5ca7722"
./build/release/compare-builds.ps1 -FirstOutputDirectory "$env:TEMP\pscan-06-c1-build-iteration003-a-5ca7722" -SecondOutputDirectory "$env:TEMP\pscan-06-c1-build-iteration003-b-5ca7722"
```

Both builds exited `0`. Each used `--network none`, the completed module cache
read-only, a read-only container root, all capabilities dropped and
no-new-privileges. Each independently passed:

- `go test -p=1 -count=1 ./...` across every package;
- `go vet -p=1 ./...`;
- Linux amd64 build and execution;
- Windows amd64 package and test compilation;
- strict release-manifest schema `2.0` and 29-asset reconciliation;
- dual product/tooling identity, Cosign claim and re-signed mutation tests;
- product runner, Gitleaks engine and rule-pack accepted-digest checks; and
- licence, SBOM, compatibility and offline verifier checks.

The comparator reported `Independent network-disabled builds are
byte-identical: 30 files`. Both manifests bind product source
`a13c28fe7273bc8dc6545f97966a02889524eb4c` / tree
`217b711ddea51fd0ea7e808edd2e27fdecef8427` and correction tooling
`5ca77226ed3996a8267caf02da366d0beb915c8d` / tree
`e65ac8ca9b1c5dc9a8a70e26c122aee3e84d32d2` in distinct roles.

## Exact deterministic digests

| Artifact | SHA-256 |
|---|---|
| `release-manifest.json` | `421e1ef4541a648b96631fa47aa95a603b91c9ae2b96ee1ef6bfdc8c7d18260e` |
| Linux amd64 bundle | `825f6c90902736de24189211e6d4ff00c66e32357508bb0b343165c7140d8889` |
| Windows amd64 bundle | `cf467b4cc5a225bba1240e2e261e5dc9acd5adc299593829f122ed51397d5d21` |
| Linux runner | `06043e9410e05a927356960f7b4ee64c9a073b7467e775452528a6813d97906e` |
| Windows runner | `1b74310e16e0df13cfd077c41338b93eabdc3d5959fb63dc256927ca51117543` |
| Linux Gitleaks | `657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586` |
| Windows Gitleaks | `b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178` |
| Rules | `cebfe007ae88a55540e42ccb8e042848d089fe2f59bcf604c5fd4904d6460792` |
| SPDX SBOM | `9cbab9f5889967b8fea7c763dfd2d4af0c44210fa628e650753ef6858f94225b` |
| `CHECKSUMS.sha256` | `a0098f984cf0c6a71fff3a591556c7f7c86853065c41e96d02fe1d7e92c2f300` |

The runner, engine and rule digests exactly match the accepted product-source
values.

## Scope, preservation and remaining boundary

The authority-to-candidate diff contains 28 paths, all inside the Correction
C1 allowlist. The exact committed bytes of `PASS-OUTCOME-SPEC-001.md` have the
same SHA-256 at authority and candidate:
`8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`.
Candidate and prior evidence commit have identical blobs for historical
`CORRECTION-C1-PREFLIGHT.md` and
`CORRECTION-C1-AUTHOR-VALIDATION-002.md`; neither was rewritten.

The verified repository identity is
`ThameeraDananjaya/project-agnostic-secret-scanner`. The proposed local tooling
tag remains absent and locked product tag `v1.0.0` still resolves to product
commit `a13c28fe7273bc8dc6545f97966a02889524eb4c` and tree
`217b711ddea51fd0ea7e808edd2e27fdecef8427`.

Actual general-purpose Linux UID/GID proof remains
`UNPROVEN_CURRENT_HOST`: `wsl.exe --list --verbose` exposes only the internal
`docker-desktop` distribution. Native Windows executable execution remains
`UNPROVEN_SMART_APP_CONTROL`; Windows cross-compilation passes. Independent
skeptical rereview remains open. No remote, push, tag, settings, workflow,
signing, attestation, draft/public release, credential, paid capability,
TruffleHog or successor action occurred. PSCAN-06 remains open and unaccepted.
