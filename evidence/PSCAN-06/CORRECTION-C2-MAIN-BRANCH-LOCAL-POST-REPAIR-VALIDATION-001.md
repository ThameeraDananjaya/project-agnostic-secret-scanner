# PSCAN-06 Correction C2 main-branch local post-repair validation 001

## Result and boundary

Local, zero-cost, no-network and no-Docker validation passed against exact
evidence candidate `a0ac587f97557b89beb3b61553fe621e80f26611`, tree
`c62698df4d9b10bb549171a34ad4369ab1a7f707`.

This is current local regression evidence after the working-tree projection
repair. It is not independent PSCAN-06 acceptance, does not rehabilitate the
terminal Recovery R1 failure, and grants no retry, remote, Docker, signing,
attestation, draft, publication or successor authority.

## Exact starting state

- Repository: `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`
- Branch: `main`
- HEAD: `a0ac587f97557b89beb3b61553fe621e80f26611`
- Tree: `c62698df4d9b10bb549171a34ad4369ab1a7f707`
- Tracked, staged and non-ignored untracked status: clean
- Ignored material outside the preserved `graphify-out/**` boundary: none
- Accepted Correction C2 tooling candidate:
  `faef8435322c9096df09b56969662411f17356ea`
- Accepted tooling tree:
  `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`
- Required prior evidence HEAD:
  `b6d341284cf63baa0502fcba319ea8cca3c7eb3a`
- The required evidence HEAD is an ancestor of this evidence candidate.

## Complete tracked-byte inspection

The current saved-project checkout was inspected without trusting ordinary Git
status. The exact `HEAD` tree and index contained 351 regular blobs:

- Raw-equal working-tree files: `333`
- Canonical LF-to-CRLF projections: `18`
- Noncanonical or mismatched files: `0`
- Non-ignored untracked paths: `0`
- Ignored paths outside `graphify-out/**`: `0`
- Index/path/mode mismatch: none
- Unsupported index flags: none

The exact LF-clone positive case then reported:

```text
Exact source trust PASS commit=a0ac587f97557b89beb3b61553fe621e80f26611 tree=c62698df4d9b10bb549171a34ad4369ab1a7f707 files=351 raw_equal=351 canonical_crlf=0
```

## Hostile source-trust matrix

`build/release/test-source-trust.ps1` ran from a new local no-hardlink,
no-network LF clone and rejected all 12 named hostile states before build
output or untrusted driver action:

1. assume-unchanged
2. skip-worktree
3. staged change
4. unstaged change
5. untracked file
6. ignored untracked file
7. filesystem-monitor configuration
8. untracked-cache configuration
9. build-driver tampering
10. path mismatch
11. mode mismatch
12. index/blob mismatch

Terminal result:

```text
SOURCE-TRUST PASS adversarial_rejections=12 exact_positive=PASS zero_build_outputs=PASS untrusted_driver_action=ABSENT
```

The unique temporary clone and fixture directory was resolved beneath the
operating-system temporary directory, checked for path containment and removed.

## No-Docker admission and process-boundary matrix

`build/release/test-image-admission.ps1` ran locally. It uses synthetic engine
and image evidence and invokes the no-Docker Windows native-boundary fixtures;
it does not start or stop Docker.

```text
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
Image admission iteration-004 PASS state-matrix=PASS parser-only=PASS closed-docker-boundary=PASS
```

## Static release identity checks

- Recovery workflow SHA-256:
  `C5F40F1B32E87C005FE33EE607AF7E3D19EE4F0173C21619E21158C31AA0BDB4`
- Workflow Git blob at the accepted tooling candidate and this evidence
  candidate: `3aa42627628b8b5298d854d29b3880cb19dc36ff`
- Build job permissions: `contents: read`
- Build runner: `ubuntu-24.04`
- Checkout and upload-artifact actions: full-SHA pinned
- Unsigned artifact retention: one day
- Unsigned artifact compression level: zero
- Signing job remains gated by exact variable/value
  `PSCAN_RELEASE_C2_GATE=PSCAN-06-C2-SIGNING-APPROVED`
- Local locked product tag `v1.0.0` resolves to
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`.
- Local locked C1 tooling tag `release-tooling-v1.0.0-c1` resolves to
  `3fb7592889820fa2739a4a53588e073689621809`.
- Local C2 tooling tag remains absent.

These are local static facts only. No remote fact was queried or inferred.

## Remaining proof and gates

Actual Linux execution, genuine Docker engine/image/container evidence,
dependency acquisition, two complete builds, every-byte comparison, workflow
identity and artifact read-back remain unperformed after Recovery R1. No new
proof-gate execution is authorized by this validation. PSCAN-06 remains open
and unaccepted overall; PSCAN-07 remains proposed and unselected; PSCAN-08
remains inactive and ineligible.
