# PSCAN-06 Validation Contract

The exact implementation candidate must pass from a clean commit using the
official Go 1.27.1 archive and full-digest Docker image:

1. revalidate every acquisition digest and populate a ledger only during the
   explicit networked phase;
2. build Gitleaks with Go 1.27.0 and compare both platform binaries with the
   accepted PSCAN-10 digests;
3. build the runner, release verifier and deterministic packaging/SBOM tools
   for Linux amd64 and Windows amd64 with Go 1.27.1;
4. pass all Go tests and all-package vet serially (`-p=1`) so the inherited
   100-millisecond deadline adversarial case is not distorted by unrelated
   concurrent package builds; compile every package/test for Windows amd64 with
   networking disabled and the module cache read-only;
5. build the complete release candidate twice from the same source commit and
   prove every file byte-identical;
6. parse every JSON contract and generated JSON asset, verify SPDX 2.3
   structure, validate archive membership/path safety, reconcile all licence
   files, and reject unexpected/unpinned workflow actions;
7. pass synthetic adversarial cases for manifest, asset, licence, SBOM,
   revocation, identity, issuer, tag and bundle mutation before scanner use;
8. confirm protected authority hashes and allowed-path scope remain unchanged;
9. perform a separate skeptical diff/evidence review.

Native Windows execution, GitHub repository/settings read-back, workflow
execution, signing, attestation, draft upload and publication are not replaced
by local success. Each remains `UNPROVEN` until its separately approved gate.

The final local execution of this contract is recorded in
`evidence/PSCAN-06/VALIDATION.md`. It binds candidate
`a13c28fe7273bc8dc6545f97966a02889524eb4c` and does not promote any reserved
remote or signing state to passed evidence.

## Correction C1 validation addendum

The correction tooling must additionally prove from its exact clean candidate:

1. the cache canary completes write, same-filesystem atomic rename, exact read
   and delete before any dependency download or success ledger;
2. an actual Linux numeric UID/GID positive case passes while wrong-owner and
   read-only bind mounts reject without acquired bytes or a ledger;
3. the acquisition container retains a read-only root, `--cap-drop ALL` and
   no-new-privileges and uses no privileged mode, added capability, mode `0777`
   or host user namespace;
4. schema `2.0` exactly binds the immutable product-source and correction-
   tooling roles, while every required-field omission, identity mutation, role
   swap, old-schema masquerade and unknown major rejects;
5. Cosign verification constrains repository, workflow ref, workflow SHA,
   trigger, certificate identity and GitHub OIDC issuer;
6. two independent network-disabled builds use the completed module cache
   read-only and produce byte-identical assets; and
7. a forced-CRLF checkout proves that the canary, acquisition and build Docker
   shell payloads are converted to LF immediately before invocation, raw
   carriage-return payloads reject, the exact normalized payloads parse in the
   pinned network-disabled image, and the read-only cache still fails closed;
   both complete reproducibility builds must execute from separate actual CRLF
   checkouts while all build inputs are first materialized from and verified
   against the exact product and correction-tooling Git trees, including exact
   path sets, modes and raw blob SHA-256 values; archive, tree, byte, cleanliness
   or extraction ambiguity must fail closed; no integrity asset normalization
   is permitted;
8. every release entrypoint used by the recovery workflow and CRLF harness is
   first materialized from the exact correction-tooling commit and verified
   against its Git blob identity; before output creation, the trusted launcher
   must reject all untracked files (including ignored files), every assume-
   unchanged or skip-worktree entry, filesystem-monitor/untracked-cache or
   partial-clone shortcuts, index path/mode/object differences, missing or
   unsupported paths/modes, and tracked raw-byte differences other than a
   byte-proved canonical LF-to-CRLF checkout projection that is never executed
   or consumed as a build input; exact Git-object materialization remains the
   sole build input; isolated adversarial cases must prove every rejection
   emits zero build files and executes no tampered driver action; and
9. the shipped runner, Gitleaks engines and rules retain the previously
   accepted product-source digests.

Local author validation remains non-acceptance. Independent skeptical review
and every remote/signing/publication owner gate remain open.
