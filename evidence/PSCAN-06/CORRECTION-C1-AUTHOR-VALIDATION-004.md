# PSCAN-06 Correction C1 author validation 004

## Result

- Superseding candidate commit:
  `a22579fd5af14473e1d49b5027f21ab591bbb589`
- Candidate tree: `382a9fff0b116d7368b1a1b321a0f40d3ebe138f`
- Authority commit:
  `3fb1b0a55dc4f48dd35464c63c768f497efbc89b`
- Result: `PARTIAL_AUTHOR_VALIDATION_PASS; ACTUAL_LINUX_HOST_PROOF_OPEN`
- Independent acceptance: not performed and not claimed
- Remote/signing/release state: not performed and not authorized

Candidate `5ca77226ed3996a8267caf02da366d0beb915c8d` is rejected because an
independent complete build from a Windows CRLF checkout consumed transformed
working-tree integrity bytes and failed
`TestPinnedRuleAndCoverageIntegrityBindings`. The append-only contradiction
and repair history is in `CORRECTION-C1-ITERATION-004.md`. Earlier author
validation remains historical evidence for its exact candidate and cannot
override the reproduced failure. This validation applies only to candidate
`a22579f`.

## Exact Git-tree materialization

`build/release/build.ps1` no longer mounts the checkout as build input. It
proves the clean source repository, exact commit-to-tree identities and locked
`v1.0.0` product identity, then materializes deterministic commit archives
with checkout conversion disabled and tar umask `0022`. Independently generated
manifests read every expected file as raw bytes through `git cat-file blob`.

Before compilation, the pinned network-disabled container proved for each
archive:

- archive SHA-256 at the container boundary;
- successful extraction;
- exact sorted path set and file count;
- exact `100644`/`100755` Git modes as extracted `0644`/`0755`; and
- every extracted file SHA-256 equal to its raw Git blob.

The locked product tree contains 288 proved blobs. The correction-tooling tree
contains 309 proved blobs. Both independent builds produced identical source
archives:

| Exact source archive | Build A SHA-256 | Build B SHA-256 |
|---|---|---|
| Product `a13c28f` / `217b711` | `f7d7607234a8f02533955875fd88de9671d58915e44bffa90ad98a3d23f07a79` | `f7d7607234a8f02533955875fd88de9671d58915e44bffa90ad98a3d23f07a79` |
| Tooling `a22579f` / `382a9ff` | `293e753f7149e689875473127da0b899baff1591ed799ee54d772e724ae69937` | `293e753f7149e689875473127da0b899baff1591ed799ee54d772e724ae69937` |

The verifier, packaging and SBOM tools, tests, schema 2.0 and release
documentation came only from the verified tooling tree. Product runners,
rules, schemas, policies, licences and integrity assets came only from the
verified locked product tree. No integrity asset was normalized or rewritten.

A separate actual CRLF clone with a synthetic untracked file was rejected at
the clean-source gate with exit `1`; it produced zero output files.

## Actual CRLF acquisition and builds

A fresh acquisition was executed from a clean detached
`core.autocrlf=true` checkout created during iteration 004:

```powershell
./build/release/acquire.ps1 -CacheDirectory "$env:TEMP\pscan-06-c1-acquisition-iteration004-29e2b43" -AllowImagePull
```

It exited `0`. Its new schema-`2.0` ledger has SHA-256
`bae3aefd32afce6a6239d817a8b736f01b302d70d74182528d636c5b106577ff`,
records `write-atomic-rename-read-delete`, and records
`completedBeforeNetworkDependencyAcquisition: true`. It contains the three
digest-pinned public archives and completed module cache. Later candidate
refinements changed only offline build/materialization and evidence bytes; the
same completed cache was mounted read-only into both final-candidate builds.

Actual final-candidate commands:

```powershell
./build/release/test-crlf-shell-payloads.ps1 -WorkingDirectory "$env:TEMP\pscan-06-c1-iteration004-crlf-build-a-a22579f" -AcquisitionDirectory "$env:TEMP\pscan-06-c1-acquisition-iteration004-29e2b43" -BuildOutputDirectory "$env:TEMP\pscan-06-c1-build-iteration004-a-a22579f"
./build/release/test-crlf-shell-payloads.ps1 -WorkingDirectory "$env:TEMP\pscan-06-c1-iteration004-crlf-build-b-a22579f" -AcquisitionDirectory "$env:TEMP\pscan-06-c1-acquisition-iteration004-29e2b43" -BuildOutputDirectory "$env:TEMP\pscan-06-c1-build-iteration004-b-a22579f"
./build/release/compare-builds.ps1 -FirstOutputDirectory "$env:TEMP\pscan-06-c1-build-iteration004-a-a22579f" -SecondOutputDirectory "$env:TEMP\pscan-06-c1-build-iteration004-b-a22579f"
```

Both checkouts were clean, detached at `a22579f`, and reported `w/crlf` for
`rules/generic/gitleaks-ignore-empty-v1.txt`. Both wrappers reported:

```text
CRLF checkout regression PASS checkout-asset-CR=PROVED raw-shell-CR=REJECT normalized-shell-CR=ABSENT canary=PASS acquisition-syntax=PASS build-syntax=PASS read-only-cache=REJECT pinned-asset-integrity=PASS complete-build=PASS
```

Each build used `--network none`, the completed module cache read-only, a
read-only container root, all capabilities dropped and no-new-privileges. Each
independently passed these exact commands:

```text
go test -p=1 -count=1 -run '^TestPinnedRuleAndCoverageIntegrityBindings$' ./tests/acceptance/gitleaks
go test -p=1 -count=1 ./...
go vet -p=1 ./...
GOOS=windows GOARCH=amd64 go test -p=1 -exec /bin/true ./...
```

This includes all unit, integration and acceptance packages, Linux execution,
Windows amd64 compilation, strict schema `2.0`, dual-identity and Cosign claim
checks, and re-signed asset/identity/role/schema/revocation mutations. The
comparator reported `Independent network-disabled builds are byte-identical:
30 files`.

The generated manifest is schema `2.0`, binds 29 assets, and binds product
source `a13c28fe7273bc8dc6545f97966a02889524eb4c` / tree
`217b711ddea51fd0ea7e808edd2e27fdecef8427` separately from tooling
`a22579fd5af14473e1d49b5027f21ab591bbb589` / tree
`382a9fff0b116d7368b1a1b321a0f40d3ebe138f`.

The compiled Linux verifier admitted the exact generated manifest and all 29
assets under a synthetic local signature-command harness in the pinned,
network-disabled container:

```text
verified release=v1.0.0 manifest_sha256=45ec3df66552e7969237338b49b49cdc31b88bb74acbf35ad247fc40b3cb06b0 assets=29
```

Synthetic Cosign and trusted-root inputs were held outside the release
directory and digest-bound. This proves parser, policy, file and command
binding only; no real signature or project data was created.

## Exact deterministic digests

| Artifact | SHA-256 |
|---|---|
| `release-manifest.json` | `45ec3df66552e7969237338b49b49cdc31b88bb74acbf35ad247fc40b3cb06b0` |
| Linux amd64 bundle | `a6f56e4bd88eedffc42236e3946032db2c76b8b912876fe1f0f63659f51b2435` |
| Windows amd64 bundle | `b5f07d6bc9b7b2333f75bc7176459b39e4c40e76f9efa370928caac0cc87325d` |
| Linux runner | `06043e9410e05a927356960f7b4ee64c9a073b7467e775452528a6813d97906e` |
| Windows runner | `1b74310e16e0df13cfd077c41338b93eabdc3d5959fb63dc256927ca51117543` |
| Linux Gitleaks | `657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586` |
| Windows Gitleaks | `b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178` |
| Rules | `cebfe007ae88a55540e42ccb8e042848d089fe2f59bcf604c5fd4904d6460792` |
| Empty ignore integrity asset | `3cdd737e34cb301ffc19861f903cf59441ff57453ddf70bd5b8c1bd63389aebe` |
| SPDX SBOM | `a1abe7b047648318da0b87ca5e7b5803035ba6cd11919e528915efbae923029e` |
| `CHECKSUMS.sha256` | `21c3d4743e48a73fcdae64e83ca8aaa3311a5439b76f805eb9eec03b5cce4a08` |

Runner, engine, rules and ignore-file digests exactly match their accepted
product-source bindings.

## Scope, preservation and remaining boundary

The authority-to-candidate diff contains 30 paths, all inside the Correction
C1 allowlist. The exact committed bytes of `PASS-OUTCOME-SPEC-001.md` have the
same SHA-256 at authority and candidate:
`8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`.
From evidence HEAD `b5f0fb2` to candidate `a22579f`, the only PSCAN-06 evidence
path added or changed is the new append-only
`CORRECTION-C1-ITERATION-004.md`; every historical evidence blob remains
unchanged.

The proposed local tooling tag remains absent. Locked product tag `v1.0.0`
still resolves to commit `a13c28fe7273bc8dc6545f97966a02889524eb4c`
and tree `217b711ddea51fd0ea7e808edd2e27fdecef8427`.

Actual general-purpose Linux UID/GID proof remains
`UNPROVEN_CURRENT_HOST`: `wsl.exe --list --verbose` exposes only the internal
`docker-desktop` distribution. Native Windows executable execution remains
`UNPROVEN_SMART_APP_CONTROL`; Windows cross-compilation passes. Independent
skeptical rereview remains open. No remote, push, tag, settings, workflow,
signing, attestation, draft/public release, credential, paid capability,
TruffleHog or successor action occurred. PSCAN-06 remains open and unaccepted.
