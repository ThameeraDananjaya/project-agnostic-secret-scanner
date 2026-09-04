# PSCAN-06 Correction C1 author validation 005

## Result and exact identities

- Superseding candidate commit:
  `3fb7592889820fa2739a4a53588e073689621809`
- Candidate tree: `7de9dc4c5725bf38cc80aa734e3c2bac0abd9762`
- Authority commit:
  `3fb1b0a55dc4f48dd35464c63c768f497efbc89b`
- Locked product commit:
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Locked product tree: `217b711ddea51fd0ea7e808edd2e27fdecef8427`
- Result: `PARTIAL_AUTHOR_VALIDATION_PASS; ACTUAL_LINUX_HOST_PROOF_OPEN`
- Independent acceptance: not performed and not claimed
- Remote/signing/release state: not performed and not authorized

Candidate `a22579fd5af14473e1d49b5027f21ab591bbb589` is rejected because
assume-unchanged driver tampering produced an empty porcelain status while the
modified working-tree build driver executed beyond its trust gate. The exact
reproduction and append-only repair history are in
`CORRECTION-C1-ITERATION-005.md`.

## Hostile source-state and bootstrap proof

Validation ran `build/release/test-source-trust.ps1` from an exact LF checkout
of the candidate. The positive checkout proved all 313 tracked files raw-equal
to `HEAD`. Twelve isolated adversarial clones each exited `1` before creating
build output and with no untrusted driver marker:

```text
assume-unchanged       REJECT output_files=0 untrusted_driver_action=ABSENT
skip-worktree          REJECT output_files=0 untrusted_driver_action=ABSENT
staged-change          REJECT output_files=0 untrusted_driver_action=ABSENT
unstaged-change        REJECT output_files=0 untrusted_driver_action=ABSENT
untracked-file         REJECT output_files=0 untrusted_driver_action=ABSENT
ignored-untracked-file REJECT output_files=0 untrusted_driver_action=ABSENT
fsmonitor-config       REJECT output_files=0 untrusted_driver_action=ABSENT
untracked-cache-config REJECT output_files=0 untrusted_driver_action=ABSENT
driver-tampering       REJECT output_files=0 untrusted_driver_action=ABSENT
path-mismatch          REJECT output_files=0 untrusted_driver_action=ABSENT
mode-mismatch          REJECT output_files=0 untrusted_driver_action=ABSENT
blob-mismatch          REJECT output_files=0 untrusted_driver_action=ABSENT
```

The exact positive summary was:

```text
Exact source trust PASS commit=3fb7592889820fa2739a4a53588e073689621809 tree=7de9dc4c5725bf38cc80aa734e3c2bac0abd9762 files=313 raw_equal=313 canonical_crlf=0
SOURCE-TRUST PASS adversarial_rejections=12 exact_positive=PASS zero_build_outputs=PASS untrusted_driver_action=ABSENT
```

The recovery workflow's independent Git-object bootstrap was also executed
locally with Git Bash. It materialized and raw-object-verified all 13 files
under `build/release` at the exact candidate:

```text
workflow exact-entry materialization PASS files=13 revision=3fb7592889820fa2739a4a53588e073689621809
```

## Actual forced-CRLF offline builds

Both final builds began from the exact LF source clone, created a separate
detached `core.autocrlf=true` clone, proved the pinned ignore integrity asset
was actually CRLF, and passed the complete source verifier. Each CRLF checkout
had 313 tracked paths: 291 were raw-equal to Git and 22 were exact canonical
CRLF projections. `BUILD-PROVENANCE.json` records
`workingTreeInputsUsed: false` and
`buildDriver: exact-git-object-materialization`.

Commands:

```powershell
& "$source/build/release/test-crlf-shell-payloads.ps1" -SourceRepository $source -SourceRevision 3fb7592889820fa2739a4a53588e073689621809 -WorkingDirectory "$env:TEMP/pscan-06-c1-005-crlf-a-3fb7592" -AcquisitionDirectory "$env:TEMP/pscan-06-c1-acquisition-iteration004-29e2b43" -BuildOutputDirectory "$env:TEMP/pscan-06-c1-005-build-a-3fb7592"
& "$source/build/release/test-crlf-shell-payloads.ps1" -SourceRepository $source -SourceRevision 3fb7592889820fa2739a4a53588e073689621809 -WorkingDirectory "$env:TEMP/pscan-06-c1-005-crlf-b-3fb7592" -AcquisitionDirectory "$env:TEMP/pscan-06-c1-acquisition-iteration004-29e2b43" -BuildOutputDirectory "$env:TEMP/pscan-06-c1-005-build-b-3fb7592"
& "$source/build/release/compare-builds.ps1" -FirstOutputDirectory "$env:TEMP/pscan-06-c1-005-build-a-3fb7592" -SecondOutputDirectory "$env:TEMP/pscan-06-c1-005-build-b-3fb7592"
```

Both wrappers reported:

```text
CRLF checkout regression PASS checkout-asset-CR=PROVED raw-shell-CR=REJECT normalized-shell-CR=ABSENT canary=PASS acquisition-syntax=PASS build-syntax=PASS read-only-cache=REJECT pinned-asset-integrity=PASS complete-build=PASS
```

Each build used the pinned image with `--network none`, a read-only root, all
capabilities dropped, no-new-privileges, and the completed module cache mounted
read-only. Both passed the focused pinned-rule test, all Go packages, serial
vet, Linux execution, Windows amd64 compilation, schema and mutation suites,
licence/SBOM checks, deterministic packaging and manifest generation. The
comparator reported:

```text
Independent network-disabled builds are byte-identical: 30 files
```

The reused completed acquisition cache was not modified by either offline
build. Its schema-2.0 ledger SHA-256 remains
`bae3aefd32afce6a6239d817a8b736f01b302d70d74182528d636c5b106577ff`
and records the pre-network semantic canary.

## Exact source materialization

| Source archive | Build A SHA-256 | Build B SHA-256 | Files |
|---|---|---|---:|
| Locked product | `f7d7607234a8f02533955875fd88de9671d58915e44bffa90ad98a3d23f07a79` | `f7d7607234a8f02533955875fd88de9671d58915e44bffa90ad98a3d23f07a79` | 288 |
| Correction tooling | `2afd41d6530ff672821b4c16d30d442237b570097eb34d0344f63035d36ace5d` | `2afd41d6530ff672821b4c16d30d442237b570097eb34d0344f63035d36ace5d` | 313 |

The tooling blob/path/mode manifests were also byte-identical. Their SHA-256
values are respectively
`c404276733ae5fcd9a55847039194c577b83f620e85838d57d5c0fc7292a628d`,
`cd6e1aaa9592d6a560a610d4dfc676facbcc205647634e818f8b8b99bb658c78`
and `de8f7f4b13a61870327ebce13bac8dc47f95d7efe8925bfead22545b32ecca9f`.
The product equivalents remain
`f9f5c508045ba4df88d7df5754ec8eccc6058ea9a9507249f828233d7929a765`,
`de76448cfa4bde1824160160a198dfe169a7e7302d9f94dd3d0740e31d93dea9`
and `ae4b9508c5783e9decc5f141f9842926b2e7b42ab7a756fb5a3f3a7c6f22dff7`.

## Generated verifier and deterministic digests

The generated Linux verifier admitted the exact manifest and all 29 assets in
the pinned network-disabled container with synthetic Cosign and trusted-root
files outside the release directory:

```text
verified release=v1.0.0 manifest_sha256=76d589bf50ce9ac5a0b9729862f51ba0da8200e160f2efa0b92c32184189be9f assets=29
```

Synthetic Cosign SHA-256 was
`306c6ca7407560340797866e077e053627ad409277d1b9da58106fce4cf717cb`;
synthetic trusted-root SHA-256 was
`0b29730cdf238cb9e4299d48acd1def08b8c21139880366f4b344715d831cfaf`.
This proves local command/parser/file binding only and creates no real
signature or credential.

| Artifact | SHA-256 |
|---|---|
| `release-manifest.json` | `76d589bf50ce9ac5a0b9729862f51ba0da8200e160f2efa0b92c32184189be9f` |
| Linux amd64 bundle | `08923f601ab796f5790b43e09abae62574e8a9360599e691050112f57079461e` |
| Windows amd64 bundle | `8559c50099dc43a80191c68a95347756eff9c54c148f6906c28bca36086f8065` |
| Linux runner | `06043e9410e05a927356960f7b4ee64c9a073b7467e775452528a6813d97906e` |
| Windows runner | `1b74310e16e0df13cfd077c41338b93eabdc3d5959fb63dc256927ca51117543` |
| Linux Gitleaks | `657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586` |
| Windows Gitleaks | `b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178` |
| Rules | `cebfe007ae88a55540e42ccb8e042848d089fe2f59bcf604c5fd4904d6460792` |
| Empty-ignore integrity asset | `3cdd737e34cb301ffc19861f903cf59441ff57453ddf70bd5b8c1bd63389aebe` |
| SPDX SBOM | `74f67421bc8d53dd143dee54cfa7bceb908c9951bc292c99852a5e6a09d23690` |
| `CHECKSUMS.sha256` | `b017a0ccab6f88842f63f845cb9517a6e11662a0cae90f9ca04eefaf58d969de` |

## Scope, preservation and remaining boundary

The authority-to-candidate diff contains 34 paths and zero paths outside the
approved Correction C1 boundary. From evidence HEAD `b4f226f` to the candidate,
no historical PSCAN-06 evidence blob changed. All 13 PowerShell files parsed,
all 18 tracked JSON files and 14 generated JSON files parsed, and all four
workflow actions remain pinned by complete 40-character commits.

The exact `PASS-OUTCOME-SPEC-001.md` blob remains
`dae97e7614e0ee4faf6f36e25a9b10d7b2d85885` with SHA-256
`8a034701867e366bf37adea6e4f47e3ffb7a53425463ceff4fb280bcb4296e74`
at both authority and candidate. The protected outcome/spec/source,
traceability and DEC-002 paths have no diff. The proposed local tooling tag is
absent, while `v1.0.0` remains locked to the product commit/tree above.

Actual general-purpose Linux UID/GID proof remains
`UNPROVEN_CURRENT_HOST`: WSL exposes only the internal `docker-desktop`
distribution. Native Windows executable execution remains
`UNPROVEN_SMART_APP_CONTROL`; Windows compilation passes. Independent
skeptical rereview remains open. No remote, push, tag, settings, workflow,
signing, attestation, draft/public release, credential, paid capability,
TruffleHog or successor action occurred. PSCAN-06 remains open and unaccepted.
