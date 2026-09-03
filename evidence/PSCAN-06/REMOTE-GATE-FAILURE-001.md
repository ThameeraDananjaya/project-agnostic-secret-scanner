# PSCAN-06 Remote Gate Failure 001

## Outcome

The owner approved the exact `PSCAN-06 REMOTE AND SIGNING GATE` on 2026-09-03.
The authorized remote setup completed, but the single approved workflow run
failed closed during pinned dependency acquisition. PSCAN-06 remains unaccepted
and open. No rerun or corrective remote mutation is authorized by this record.

## Exact remote identity

- Repository: `ThameeraDananjaya/project-agnostic-secret-scanner`
- Repository ID: `1355442997`
- Repository node ID: `R_kgDOUMprNQ`
- Visibility: public
- Default branch: `main`
- Remote `main`: `c623813f68fa7cfadaf4029444d1b3027fbdd584`
- Locked tag `v1.0.0`: `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Workflow run ID: `33709197614`
- Workflow URL:
  `https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/actions/runs/33709197614`
- Run source ref: `refs/tags/v1.0.0`
- Run source commit: `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Run conclusion: `failure`

## Controls proved before the run

- Account-level Actions budget: USD 0; stop usage enabled.
- Actions: enabled only for the four exact full-SHA action identities in the
  workflow; all GitHub-owned and verified actions are not broadly allowed.
- Default workflow token: read-only; pull-request approval disabled.
- Artifact and log retention: one day.
- Repository variable:
  `PSCAN_RELEASE_GATE=PSCAN-06-SIGNING-APPROVED`.
- Immutable releases: enabled.
- `release-v1`: required reviewer `ThameeraDananjaya`; self-review permitted;
  administrator bypass disabled; only tag `v1.0.0` may deploy.
- Active `main` rules: deletion blocked, force pushes blocked, linear history
  required.
- Active `v1.0.0` tag rules: deletion and update blocked.

## Failure evidence

The immutable-invocation check and credential-free checkout passed. The pinned
Docker image and all three pinned source archives were acquired and their
digests passed. The acquisition container then repeatedly reported:

```text
go: writing go.mod cache: mkdir /gomodcache/cache: permission denied
```

`build/release/acquire.ps1` failed at line 106 with
`Pinned dependency acquisition failed`. The build, independent rebuild,
byte-comparison and artifact-transfer steps were skipped. The signing,
attestation and draft-release job was skipped.

## Negative-state read-back

- Workflow artifacts: zero.
- GitHub releases: zero.
- Draft releases: zero.
- Keyless signatures: zero.
- GitHub artifact attestations from this run: zero.
- No release was published.
- No paid usage was authorized or required.
- PSCAN-07 and PSCAN-08 remain unselected and inactive.

## Required next decision

A bounded correction must first reproduce and fix the Linux bind-mount
ownership behavior, add a Linux-host acquisition test that would have caught
this path, revalidate the complete candidate, and define a safe tag/version
strategy. The already-public, locked `v1.0.0` tag must not be moved, deleted or
silently replaced. Any correction, tag/version change, workflow rerun or remote
settings change requires a new exact owner decision.
