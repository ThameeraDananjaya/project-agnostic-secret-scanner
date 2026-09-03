# PSCAN-06 Correction C1 author validation 002

## Result

- Superseding candidate commit:
  `f24b832ebe5f6749aa0ab910e1b9be279065ebb1`
- Candidate tree: `74911002b78b0cf11192dd61eb8db321c1d04876`
- Authority commit:
  `3fb1b0a55dc4f48dd35464c63c768f497efbc89b`
- Result: `PARTIAL_AUTHOR_VALIDATION_PASS; ACTUAL_LINUX_HOST_PROOF_OPEN`
- Independent acceptance: not performed and not claimed
- Remote/signing/release state: not performed and not authorized

This record supersedes the candidate identity and deterministic digests in
`CORRECTION-C1-AUTHOR-VALIDATION.md`. The earlier record remains factual for
candidate `f4ad426`; `CORRECTION-C1-ITERATION-002.md` explains why its Ubuntu
wrong-owner evidence fixture was corrected before remote use.

## Fresh exact-candidate execution

From clean candidate `f24b832`, both of these pinned builds exited `0`:

```powershell
./build/release/build.ps1 -AcquisitionDirectory "$env:TEMP\pscan-06-c1-acquisition-20260904" -OutputDirectory "$env:TEMP\pscan-06-c1-build-c-20260904"
./build/release/build.ps1 -AcquisitionDirectory "$env:TEMP\pscan-06-c1-acquisition-20260904" -OutputDirectory "$env:TEMP\pscan-06-c1-build-d-20260904"
```

Each build used the completed schema-`2.0` acquisition ledger, networking
disabled, the module cache read-only, a read-only container root, all
capabilities dropped and no-new-privileges. Each independently passed:

- every Go unit, integration and acceptance package;
- `go vet -p=1 ./...`;
- Linux execution and Windows amd64 package/test compilation;
- product runner, Gitleaks engine and rule-pack accepted-digest checks; and
- schema-`2.0`, dual-identity, Cosign-claim and re-signed mutation tests.

`compare-builds.ps1` reported 30 byte-identical output files. The compiled
Linux verifier admitted the generated manifest, all 29 bound assets, exact
policy and bootstrap revocation evidence under a synthetic local signature
command. This verifies parser/policy/file binding only; no real project bundle,
OIDC request or signature was created.

## Exact identities

- Product source: tag `v1.0.0`, commit
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`, tree
  `217b711ddea51fd0ea7e808edd2e27fdecef8427`.
- Release tooling: proposed tag `release-tooling-v1.0.0-c1`, commit
  `f24b832ebe5f6749aa0ab910e1b9be279065ebb1`, tree
  `74911002b78b0cf11192dd61eb8db321c1d04876`.
- Workflow: `.github/workflows/release-recovery-v1.0.0.yml`, ref
  `refs/tags/release-tooling-v1.0.0-c1`, workflow SHA equal to the tooling
  commit, trigger `workflow_dispatch`.

## Exact deterministic digests

| Artifact | SHA-256 |
|---|---|
| `release-manifest.json` | `6fd7269107d2706f9b5e1cadba985a8372113f2443438ef8103f8d856d6a84e9` |
| Linux amd64 bundle | `091effdc3969f5c1705d30dacdd667a03cb8eba3c82eab069a3065d28f404cfa` |
| Windows amd64 bundle | `2eab88f8db5ff8e8d4d59c56b01bcef1b9944c4ff5665659c35741b51e332878` |
| Linux runner | `06043e9410e05a927356960f7b4ee64c9a073b7467e775452528a6813d97906e` |
| Windows runner | `1b74310e16e0df13cfd077c41338b93eabdc3d5959fb63dc256927ca51117543` |
| Linux Gitleaks | `657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586` |
| Windows Gitleaks | `b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178` |
| Rules | `cebfe007ae88a55540e42ccb8e042848d089fe2f59bcf604c5fd4904d6460792` |
| SPDX SBOM | `82a0efdc15f0dbc8da51dd427eb119063aadf4459cff0bb3a088c213ef0f2a7f` |
| `CHECKSUMS.sha256` | `ca50e16ba5fb48ec2b94d92738644d23cb4ae5b6124ecba0c35a880effad928b` |

Product runner, engine and rule digests remain identical to the accepted
product-source candidate.

## Honest remaining boundary

Windows Docker Desktop canary success and read-only rejection are executed
evidence. The actual-Linux positive, wrong-owner and read-only cases remain
`UNPROVEN_CURRENT_HOST`: only Docker Desktop's internal WSL distribution is
available and it explicitly forbids direct Docker CLI use. The corrected
Ubuntu workflow fixture is root-owned mode `0755`, so the mapped non-root
container cannot write while the host can prove the directory remains empty.
That workflow has not been run because remote execution is owner-gated.

PSCAN-06 remains open and unaccepted. Independent skeptical review, actual-
Linux execution and every push/tag/workflow/signing/attestation/release action
remain separate gates. The proposed tooling tag is absent; locked product tag
`v1.0.0` remains unchanged.
