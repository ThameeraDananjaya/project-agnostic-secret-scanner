# PSCAN-06 Local Validation Evidence

## Result

Local validation passed at `2026-09-03T02:19:47Z` for exact candidate commit
`a13c28fe7273bc8dc6545f97966a02889524eb4c` and tree
`217b711ddea51fd0ea7e808edd2e27fdecef8427`. This is local evidence only. It
does not establish remote repository identity/settings, a valid project
keyless signature, artifact attestations, a draft release or publication.

## Isolated build and test environment

- Docker Desktop 4.87.0, Engine 29.7.2.
- Build image:
  `golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`.
- Runner Go: official `go1.27.1` Linux amd64 archive, SHA-256
  `63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445`.
- Engine Go: official `go1.27.0` Linux amd64 archive, SHA-256
  `675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685`.
- Gitleaks source commit:
  `83d9cd684c87d95d656c1458ef04895a7f1cbd8e`, archive SHA-256
  `6b2638a733b85619dc80bdf28e84e4fed7e526a761ab5c148fbf67695aea2115`.
- Acquisition cache:
  `C:\Users\ITDan\AppData\Local\Temp\pscan-06-acquisition-c-20260903`.
- Both builds used Docker `--network none`, a read-only source mount, read-only
  module cache, read-only root filesystem, dropped capabilities,
  `no-new-privileges`, bounded processes/CPU/memory and bounded tmpfs.

Each build passed all packages with `go test -p=1 -count=1 ./...`, all-package
`go vet -p=1 ./...`, and compile-only Windows amd64
`go test -p=1 -exec /bin/true ./...`. Native Windows execution remains
`UNPROVEN_SMART_APP_CONTROL`; cross-compilation is not promoted as native proof.

## Reproducibility and asset identities

- Build A:
  `C:\Users\ITDan\AppData\Local\Temp\pscan-06-repro-a-a13c28f\dist`.
- Build B:
  `C:\Users\ITDan\AppData\Local\Temp\pscan-06-repro-b-a13c28f\dist`.
- All 29 files have identical names, sizes and SHA-256 values in A and B.
- Manifest assets: 28. Checksum rows: 27. SPDX packages: 95.
- Each Linux tar.gz and Windows ZIP has 84 unique, confined members. No
  absolute, drive-qualified, dot, dot-dot, duplicate or special tar member was
  admitted.

| Artifact | SHA-256 |
|---|---|
| `release-manifest.json` | `5781337e564be7ff243934f54a6186294d4b35d0fe0dbcd9ace8567d0b4cebcf` |
| Linux amd64 archive | `979efee7918fbae551d4fe976de0e43d64b663a8f3838680281932011195c8ff` |
| Windows amd64 archive | `d57281084205f8efca5492cbc23bd9f11e542781cd5d4715c756291d41f949fb` |
| Linux runner | `06043e9410e05a927356960f7b4ee64c9a073b7467e775452528a6813d97906e` |
| Windows runner | `1b74310e16e0df13cfd077c41338b93eabdc3d5959fb63dc256927ca51117543` |
| Linux release verifier | `3734ad243100d5ce46c18577ec29f8bc55b2b9a6dceba08a463d524947303df0` |
| Windows release verifier | `965ca7f75e45bfd2d894566280d5f0047279e19bc4272eb0267ec9aeac728a98` |
| Linux Gitleaks | `657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586` |
| Windows Gitleaks | `b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178` |
| SPDX SBOM | `73143956a9cf76be0913e8f28d8fe9cb0aa4ef51bdee9e7b9ebcc2b0dd6137b9` |
| `CHECKSUMS.sha256` | `42c22c32e27c9c4fb71b04265bc45286e3fdb6e4e235f183be5480fa135c94d3` |

Every manifest asset size/digest and every checksum row was independently
recomputed. All tracked JSON parsed. All PowerShell files passed AST parsing.
All four workflow actions are complete 40-character SHA pins.

## Adversarial and trust-root checks

The integration and acceptance supply-chain suites passed ten consecutive runs
under the same network-disabled/read-only-cache boundary. They reject unknown
or duplicate manifest members; one-byte runner, licence, SBOM, revocation and
bundle changes; wrong identity, issuer and tag; and mutated Cosign/trusted-root
bytes before scanner execution.

Pinned Cosign v3.1.3 SHA-256
`4629c757b7618056f8ddd7e2625ae9fdd94c0372a65049520bc7d9df9efc7f71`
verified its upstream bundle under Docker `--network none` with independently
acquired trusted-root SHA-256
`844a1c6de3986c9f02070266b25e0d1a2fa99ceccc89f6b9ad90aae47b62a16e`.
The definitive unsigned candidate was supplied to the scanner release verifier
with those exact external bytes and rejected with exit code 40 because the
required project bundle is absent.

## Scope and authority

The activation-to-candidate diff contains 23 files, all within PSCAN-06 allowed
paths. `PASS-OUTCOME-SPEC-001`, `PASS-SPEC-001`, `TRACEABILITY`, `DEC-002`,
PSCAN-05, PSCAN-09 and PSCAN-10 authority files are unchanged. The worktree was
clean and no Git remote existed at the candidate boundary.

## Honest failed-path record

Fail-closed iteration rejected an incomplete lazy Gitleaks module cache, a
scanner acquisition mode that would update read-only `go.sum`, concurrent
package timing distortion, execution-bit mutation on a read-only mount, the
obsolete Cosign `--offline` flag, a missing current attestation permission, an
out-of-scope test location and a PSCAN-07-valued signing gate. Each was
corrected before the definitive candidate; none is used as passing evidence.

## Conclusion

Every locally authorized PSCAN-06 deliverable and proof passes. Acceptance and
closeout fail closed because the separately approved remote/settings/workflow/
signing/attestation/draft proof is absent. See `evidence/PSCAN-06/GATE.md`.
