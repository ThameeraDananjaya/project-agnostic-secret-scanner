# PSCAN-06 Correction C2 iteration 005 author validation

Date: 2026-09-05
Implementation session: `01a071ef-5eeb-7a91-8d82-1feb22a9c2f6`

## Exact candidate

- Commit: `050f5862841a1fcf764ffc75d98a742968ac301e`
- Tree: `53e38f8b2169a28a398162265b6e5a3dfad91cfe`
- Authority parent:
  `13496a35eab71482f2e908bd0aac4b563941682d`
- Isolated checkout:
  `C:\Users\ITDan\AppData\Local\Temp\pscan-06-c2-i005-preclaim-e7177ab9e4a0425c82b69a2ebf6ae713`
- Candidate status before and after validation: zero tracked, untracked or
  ignored entries.

The isolated checkout was created locally with `--no-hardlinks`, without
network access, and materialized with LF bytes. The saved `main` checkout and
its preserved ignored `graphify-out/**` material were not used for
implementation.

## Exact source and format checks

```text
Exact source trust PASS commit=050f5862841a1fcf764ffc75d98a742968ac301e tree=53e38f8b2169a28a398162265b6e5a3dfad91cfe files=344 raw_equal=344 canonical_crlf=0
elapsed_ms=2990

PowerShell parse and contract JSON parse PASS
powershell_files=19 elapsed_ms=171

Committed diff/path check PASS
changed_paths=8 outside_allowed=0 elapsed_ms=77
```

The committed range changes only the two boundary scripts, the authorized
release/task/validation documents and iteration-005 evidence. `git diff
--check` passed.

## Hostile source-trust and CRLF checks

All twelve cases rejected with nonzero status, zero build output and no
untrusted driver action:

```text
assume-unchanged, skip-worktree, staged-change, unstaged-change,
untracked-file, ignored-untracked-file, fsmonitor-config,
untracked-cache-config, driver-tampering, path-mismatch, mode-mismatch,
blob-mismatch
SOURCE-TRUST PASS adversarial_rejections=12 exact_positive=PASS zero_build_outputs=PASS untrusted_driver_action=ABSENT
elapsed_ms=43610
```

The real CRLF fixture reported:

```text
CRLF host-only regression PASS checkout-asset-CR=PROVED raw-shell-CR=REJECT normalized-shell-CR=ABSENT docker=NOT_INVOKED
elapsed_ms=7763
```

## Native boundary and admission checks

The direct isolated boundary suite reported:

```text
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
elapsed_ms=77119
```

The image-admission state matrix independently reran the complete boundary
suite and reported:

```text
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
Image admission iteration-004 PASS state-matrix=PASS parser-only=PASS closed-docker-boundary=PASS
elapsed_ms=77093
```

Static inventory found the finite nine-operation table, zero ambient Docker
calls in workflow-reachable scripts, zero duplicate native runners and the one
closed boundary referenced by all eight production/test callers. Product tag
`v1.0.0` remained at `a13c28fe7273bc8dc6545f97966a02889524eb4c`;
C1 tooling tag `release-tooling-v1.0.0-c1` remained at
`3fb7592889820fa2739a4a53588e073689621809`; proposed C2 tooling tag
`release-tooling-v1.0.0-c2` remained absent. A final Windows process inventory
found zero surviving processes marked with `native-fixture.ps1`.

## Limits and open gates

This validation used no Docker command, image, container, network, download,
Go, Gitleaks, Cosign, build, scanner, remote read/write, tag mutation, workflow,
credential, signing, attestation, draft or publication action.

This is author validation only. Independent skeptical review has not accepted
the candidate. Actual Linux, genuine Docker/image/container execution,
dependency acquisition, complete builds, byte comparison, remote proof, C2
tooling-tag creation, signing, attestation, draft and publication remain open
and separately owner-gated. PSCAN-06 remains open; PSCAN-07 remains proposed
and unselected; PSCAN-08 remains inactive and ineligible. No successor was
selected or activated.
