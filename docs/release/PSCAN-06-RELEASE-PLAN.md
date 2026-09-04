# PSCAN-06 Release Control Plan

Correction C2 uses a separate manual-dispatch recovery workflow, exact product
source and exact correction-tooling identities. Its build job has read-only
contents permission and runs only from proposed tooling tag
`release-tooling-v1.0.0-c2`. Its signing/draft job is skipped unless repository
variable `PSCAN_RELEASE_C2_GATE` exactly equals
`PSCAN-06-C2-SIGNING-APPROVED`, is protected by environment `release-v1`, and
the product source remains locked tag `v1.0.0`, commit
`a13c28fe7273bc8dc6545f97966a02889524eb4c`. Creating the tooling tag or gate
and running the workflow remain separately owner-gated.

Build acquisition accepts only digest-pinned Go/Gitleaks archives and the exact
canonical Docker image repository plus digest. Before any image pull, it maps
the Linux host numeric UID/GID, proves host cache write, atomic rename, read and
delete, proves wrong-owner/read-only rejection, and proves CRLF normalization
without Docker. Image admission permits one explicit pull only for
`docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`,
then re-inspects the canonical repository digest. Container cache and shell
parser proofs follow with `--pull=never --network none`.
Compilation, tests, vet, Windows cross-compilation and two byte-for-byte builds
then run with networking disabled and the module cache read-only. Workflow-
transfer evidence is uncompressed and retained one day; it is not a release.

The gated job verifies Cosign `v3.1.3` by SHA-256, requests one GitHub OIDC
identity, signs schema-`2.1` manifest bytes, and verifies exact repository,
workflow ref, workflow SHA, trigger, certificate identity and issuer against an
explicit authenticated trusted root. It then
creates SBOM-bound GitHub attestations and creates one draft containing every
asset. It has no publish command. Publication needs a separate PSCAN-07 owner
decision after independent acquisition and verification.

Cosign v3.1.3 removed the legacy `--offline` flag. The workflow acquires an
authenticated Sigstore trusted-root document through pinned Cosign/TUF and
passes it explicitly to `verify-blob`. Consumer verification acquires the root
independently online, records its digest, then runs the same bundle verification
inside a network-disabled boundary.

```mermaid
flowchart LR
    A[Accepted local commit] --> B{Remote setup approval}
    B -->|no| H[Hold locally]
    B -->|yes| C[Public repo plus read-back controls]
    C --> D{Workflow and signing approval}
    D -->|no| H
    D -->|yes| E[Two offline reproducible builds]
    E --> F[Environment approval, keyless sign, attest]
    F --> G[Draft with all assets]
    G --> I{Independent verify and publication approval}
    I -->|no| H
    I -->|yes PSCAN-07 only| J[Publish immutable v1.0.0]
```

Full action pins:

- `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1`
- `actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a`
- `actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c`
- `actions/attest@1e69f48acb82d1966a394da916b4c1698aa569d6`

The attestation job alone receives `id-token: write`, `attestations: write` and
`artifact-metadata: write`; its draft upload also needs `contents: write`.
These permissions are absent from the build job.

No Gitleaks action, mutable action tag, signing key, stored cloud credential,
larger runner, paid feature, auto-update, TruffleHog material, automatic
publication or consuming-project data is present.
