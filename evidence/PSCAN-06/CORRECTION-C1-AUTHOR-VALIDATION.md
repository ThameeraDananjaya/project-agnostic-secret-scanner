# PSCAN-06 Correction C1 author validation

## Result

- Candidate commit:
  `f4ad42664eabedba61226bccdea2bc64e1156b7e`
- Candidate tree: `bcd4336c80a547fccef9e8c3cd9c1f1a61d26642`
- Base authority commit:
  `3fb1b0a55dc4f48dd35464c63c768f497efbc89b`
- Result: `PARTIAL_AUTHOR_VALIDATION_PASS; ACTUAL_LINUX_HOST_PROOF_OPEN`
- Independent acceptance: not performed and not claimed
- Remote/signing/release state: not performed and not authorized

The implementation candidate is complete and all checks callable on the author
host passed. The author host is Windows with Docker Desktop and has no supported
general-purpose Linux WSL distribution. The mandatory actual-Linux positive and
wrong-owner host-mapping executions therefore remain open. The pinned recovery
workflow contains those Ubuntu-host tests, but running it is a separately
owner-gated remote mutation. This record does not convert workflow design into
executed evidence and does not establish correction acceptance or remote-gate
eligibility.

## Acquisition and cache boundary

Fresh command:

```powershell
./build/release/acquire.ps1 -CacheDirectory "$env:TEMP\pscan-06-c1-acquisition-20260904"
```

Actual result: exit `0`. The canary completed before all archive/module
downloads. Only after full acquisition did `acquisition-ledger.json` appear with:

- schema `2.0`;
- semantics `write-atomic-rename-read-delete`;
- `completedBeforeNetworkDependencyAcquisition: true`; and
- host mode `windows-docker-desktop-default-user`.

The same canary passed against a fresh writable Windows Docker Desktop bind
mount. A separately mounted read-only cache rejected at the first write and
left no canary or ledger. Static inspection confirmed the acquisition and
canary containers retain `--read-only`, `--cap-drop ALL` and
`no-new-privileges`, while the Linux branch maps numeric `id -u`/`id -g`, checks
host ownership and uses mode `0700`. No privileged mode, capability addition,
mode `0777` or host user namespace is present.

Actual Linux host result: `UNPROVEN_CURRENT_HOST`. The positive, wrong-owner and
read-only Linux harness is `build/release/test-cache-boundary.ps1` and is called
by the recovery workflow on `ubuntu-24.04` before acquisition.

## Clean reproducible builds

The following commands each ran from the exact clean candidate with the fresh
completed cache:

```powershell
./build/release/build.ps1 -AcquisitionDirectory "$env:TEMP\pscan-06-c1-acquisition-20260904" -OutputDirectory "$env:TEMP\pscan-06-c1-build-a-20260904"
./build/release/build.ps1 -AcquisitionDirectory "$env:TEMP\pscan-06-c1-acquisition-20260904" -OutputDirectory "$env:TEMP\pscan-06-c1-build-b-20260904"
./build/release/compare-builds.ps1 -FirstOutputDirectory "$env:TEMP\pscan-06-c1-build-a-20260904" -SecondOutputDirectory "$env:TEMP\pscan-06-c1-build-b-20260904"
```

Actual results: both builds exit `0`; the comparison reports 30 byte-identical
files. Both builds used `--network none`, a read-only module-cache bind, a
read-only container root, all capabilities dropped and no-new-privileges.

The generated manifest binds:

- product tag/commit/tree: `v1.0.0`,
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`,
  `217b711ddea51fd0ea7e808edd2e27fdecef8427`;
- tooling tag/commit/tree: `release-tooling-v1.0.0-c1`,
  `f4ad42664eabedba61226bccdea2bc64e1156b7e`,
  `bcd4336c80a547fccef9e8c3cd9c1f1a61d26642`; and
- recovery workflow/ref/SHA/trigger:
  `.github/workflows/release-recovery-v1.0.0.yml`,
  `refs/tags/release-tooling-v1.0.0-c1`, candidate commit above,
  `workflow_dispatch`.

The compiled Linux verifier admitted the exact generated manifest, all 29
bound assets and bootstrap revocation evidence using a synthetic local
signature-command harness. This proves parser/policy/file binding, not a real
signature. The exact unsigned candidate has no project signature bundle.

## Tests and adversarial results

Inside each pinned, network-disabled build container:

- all Go unit, integration and acceptance packages passed;
- `go vet -p=1 ./...` passed;
- Windows amd64 package/test compilation passed;
- Linux execution passed;
- schema `2.0` exact parsing passed; and
- re-signed mutations of both product/tooling tags, commits, trees, workflow,
  ref, SHA and trigger rejected, as did missing identities, role swapping,
  old-schema masquerade and unknown major.

The Cosign command test proved the verifier supplies exact certificate filters
for repository, workflow ref, workflow SHA, trigger, certificate identity and
OIDC issuer. No scanner ran during offline release verification.

## Exact deterministic digests

| Artifact | SHA-256 |
|---|---|
| `release-manifest.json` | `226d72d96c573f6f96a591b0998303eee46556a29ea85bb3f10a603517b8a846` |
| Linux amd64 bundle | `489a34b2b93b73ef51a1ca24d8bbe73d529a66188ce54f6b9719d8eecaf2ab25` |
| Windows amd64 bundle | `29a47208e228ee270747fb84eefb7f083ca92e88f20bd988d380cc1d9933578d` |
| Linux runner | `06043e9410e05a927356960f7b4ee64c9a073b7467e775452528a6813d97906e` |
| Windows runner | `1b74310e16e0df13cfd077c41338b93eabdc3d5959fb63dc256927ca51117543` |
| Linux Gitleaks | `657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586` |
| Windows Gitleaks | `b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178` |
| Rules | `cebfe007ae88a55540e42ccb8e042848d089fe2f59bcf604c5fd4904d6460792` |
| SPDX SBOM | `ac67d3d811abff2dec060601b5d4cfa538358f970586de49a8e487a1eac4c8c0` |
| `CHECKSUMS.sha256` | `9976d684851720ec3b0461ac4404d79cd22a701a7aaf9425f191a8824f0bd461` |

Runner, engine and rule digests exactly match the previously accepted product-
source values.

## Scope and preservation

The authority-to-candidate diff contains 22 paths, all inside the Correction C1
allowlist. Controlling PASS-OUTCOME-SPEC-001 SHA-256 remains
`8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`.
Existing release-manifest schemas `1.0` and `1.1`, predecessor decisions and
historical evidence are unchanged. No remote action, tag creation, workflow
run, signing, attestation, draft/public release, credential, paid capability,
consumer data, TruffleHog work or successor work occurred.

## Remaining gates

1. Execute the actual Linux positive, wrong-owner and read-only-cache cases on
   an approved Linux host and preserve command/output evidence.
2. Perform independent skeptical review without relying on this summary.
3. Only after complete local acceptance may the owner consider an exact remote
   correction gate. The proposed tooling tag does not yet exist.
