# PSCAN-06 Correction C2 iteration 005 preflight and claim

Date: 2026-09-05
Implementation session: `01a071ef-5eeb-7a91-8d82-1feb22a9c2f6`
Delegating session: `01a053a9-298b-75e1-9d49-7d2261bbb951`
Authority approval session: `01a071da-6084-7d82-948c-00a2f2029cfd`

## Claim

This fresh implementation session claims only PSCAN-06 Correction C2 iteration
005. It does not claim, select, activate or implement PSCAN-07, PSCAN-08 or any
other successor. The controlling authority is
`CORRECTION-C2-ITERATION-005-AUTHORITY.md`; its SHA-256 at the claim boundary was
`f866ededd850de754c533fb3d9ac280d20c19ce72b9e45c0744ff9acf64e49a9`.

## Starting identity and source trust

- Isolated implementation checkout:
  `C:\Users\ITDan\AppData\Local\Temp\pscan-06-c2-i005-preclaim-e7177ab9e4a0425c82b69a2ebf6ae713`
- The checkout is a new local `--no-hardlinks`, no-network clone of the saved
  repository, configured for LF materialization.
- Exact starting commit:
  `13496a35eab71482f2e908bd0aac4b563941682d`.
- Exact starting tree:
  `3b0e185e31293aed26b27165274e09303bd164cc`.
- `git status --porcelain=v1 --untracked-files=all` returned zero entries before
  the claim.
- A new non-profile PowerShell probe reported `PscanNativeBoundary` absent.
- The implementation session differs from both the delegating session and the
  authority approval session.

The exact every-byte trust check reported:

```text
Exact source trust PASS commit=13496a35eab71482f2e908bd0aac4b563941682d tree=3b0e185e31293aed26b27165274e09303bd164cc files=342 raw_equal=342 canonical_crlf=0
```

The saved checkout remains at the exact authority commit with no tracked or
index diff. Its pre-existing ignored `graphify-out/**` material is preserved and
is not inspected, deleted, moved or copied. Integration into saved `main` is
permitted only after final revalidation and only by a verified ancestry-safe
local fast-forward with no drift or overlap.

## Locked identities and gates

- Product tag `v1.0.0` remains at
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`.
- C1 tooling tag `release-tooling-v1.0.0-c1` remains at
  `3fb7592889820fa2739a4a53588e073689621809`.
- Proposed C2 tooling tag `release-tooling-v1.0.0-c2` is absent.
- No Docker, network, download, Go, Gitleaks, Cosign, build, scanner, remote,
  signing, attestation, draft, publication or successor action was performed.

The complete ordered reading map and every tracked file in the iteration-005
allowed path surface were read before this claim. Reachable production callers
of `docker-execution.ps1` were enumerated before implementation. Missing or
out-of-scope changes remain terminal and return to the owner.
