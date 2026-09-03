# PSCAN-06 Validation Contract

The exact implementation candidate must pass from a clean commit using the
official Go 1.27.1 archive and full-digest Docker image:

1. revalidate every acquisition digest and populate a ledger only during the
   explicit networked phase;
2. build Gitleaks with Go 1.27.0 and compare both platform binaries with the
   accepted PSCAN-10 digests;
3. build the runner, release verifier and deterministic packaging/SBOM tools
   for Linux amd64 and Windows amd64 with Go 1.27.1;
4. pass all Go tests, all-package vet and compile every package/test for Windows
   amd64 with networking disabled and the module cache read-only;
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
