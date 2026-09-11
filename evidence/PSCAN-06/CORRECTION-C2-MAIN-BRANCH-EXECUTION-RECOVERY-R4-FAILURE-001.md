# PSCAN-06 Correction C2 main-branch execution Recovery R4 failure 001

## Outcome

- Task: `PSCAN-06`
- Correction: `C2`
- Recovery: `R4`
- Gate: build-only actual-Linux and genuine-Docker proof
- Execution session: `01a08f7a-7d13-7ee1-be72-a997e4ff9cb6`
- Date: `2026-09-11`
- Terminal result: `FAIL_CLOSED_BEFORE_REMOTE_C2_TAG`
- Workflow dispatches consumed by Recovery R4: zero
- Artifact created or downloaded by Recovery R4: none
- Local Docker action: none
- Signing, attestation, draft, release or publication action: none
- Successor action: none

Recovery R4 passed the complete mandatory fresh preflight and began the single
authorized ordered transaction. Remote `main` fast-forwarded exactly from
`3523e4409ebc53cc1931e3dcaf7d1eea74bb15ab` to the committed Recovery R4
authority `a82a3a64a04ee2d8b60757866c9864cd1b73b54b`. The exact unsigned local
tag `release-tooling-v1.0.0-c2` was then created at accepted tooling candidate
`faef8435322c9096df09b56969662411f17356ea`.

The first tag-push command failed locally before creating the remote tag. Git
rejected the shell-expanded refspec with:

```text
fatal: invalid refspec 'refs/tags//tags/release-tooling-v1.0.0-c2'
```

The no-retry condition made this terminal. The command was not corrected or
repeated. No C2 tag ruleset was created, no setting was changed, and no
workflow was dispatched.

## Exact execution context

- Saved project root:
  `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`
- Branch: `main`
- Starting HEAD:
  `a82a3a64a04ee2d8b60757866c9864cd1b73b54b`
- Starting tree:
  `c7c5f8dedace8ca9ca4387bd915f79afef1dd665`
- Direct parent:
  `6b7ae32f8768c1e4687e5134696688088f65386b`
- Model: `gpt-daybreak-blue-latest`
- Reasoning effort: `xhigh`
- Fresh task source: agent-created task targeting the saved local project
  directly, not a worktree or detached checkout
- Preflight tracked, staged and non-ignored untracked state: clean
- Preserved ignored paths: 65 paths, all under `graphify-out/**`

The task/session metadata, saved-project identity, branch, commit, model and
reasoning effort were verified before the first repository action. The exact
root, tree, parent and clean state were then verified before any remote action.

## Local source-trust preflight

Complete index, path, mode, flag and every-byte comparison reported:

- tracked paths: 356
- raw-equal to committed blob: 338
- canonical whole-file LF-to-CRLF projections: 18
- unexpected tracked-byte mismatches: zero
- non-ignored untracked paths: zero
- unsupported Git source-trust configuration: none
- unexpected ignored paths: zero

Required identities and hashes passed:

- `PASS-OUTCOME-SPEC-001` SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- original C2 build-only gate authority SHA-256:
  `2A84E7109F151E55399870FB36D64AEDB2213E123976779CD1E9648AC1B9DAC2`
- iteration-005 authority SHA-256:
  `F866EDEDD850DE754C533FB3D9AC280D20C19CE72B9E45C0744FF9ACF64E49A9`
- iteration-005 local acceptance SHA-256:
  `1C9392E1899D1EB83C2DBDD673C0369042C5F6255D3D0BE9558C55281EDB72A9`
- Recovery R3 terminal failure SHA-256:
  `339E5D66B655B5753D239D5347F1BBAFF0E5858C4CCC289368196D014AE07C12`
- accepted candidate tree:
  `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`
- required evidence HEAD tree:
  `983028703e229abb082938d0e4417c506aac0d8c`
- required evidence HEAD direct parent:
  `faef8435322c9096df09b56969662411f17356ea`

## Workflow preflight

The workflow at accepted candidate
`faef8435322c9096df09b56969662411f17356ea` matched SHA-256
`C5F40F1B32E87C005FE33EE607AF7E3D19EE4F0173C21619E21158C31AA0BDB4`.
Its `build` job remained on standard `ubuntu-24.04`, had only
`contents: read`, and used only the recorded full-SHA action pins. Its sole C2
artifact remained named `pscan-v1.0.0-c2-unsigned-candidate` with
`retention-days: 1` and `compression-level: 0`.

The signing job remained conditional on exact repository variable
`PSCAN_RELEASE_C2_GATE=PSCAN-06-C2-SIGNING-APPROVED`. Current read-back found
that variable absent.

## Current remote and zero-spend preflight

Read-only GitHub verification established before mutation:

- active authenticated account: `ThameeraDananjaya`
- credential scopes: `gist`, `read:org`, `repo`, `user`, `workflow`
- repository: `ThameeraDananjaya/project-agnostic-secret-scanner`
- numeric repository ID: `1355442997`
- visibility: public
- default branch: `main`
- origin:
  `https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner.git`
- Actions: enabled with selected-actions policy and required SHA pinning
- workflow default permission: read
- artifact and log retention: one day
- selected actions: only the exact recorded `checkout`, `upload-artifact`,
  `download-artifact` and `attest` SHAs
- locked `v1.0.0` and `release-tooling-v1.0.0-c1` tags: exact, with active
  exact-name update/deletion-denial rulesets and no bypass actors
- local and remote `release-tooling-v1.0.0-c2`: absent before execution
- prior C2 workflow runs: zero
- prior C2 artifacts: zero
- repository releases and deployments: zero
- sole open PSCAN task: `PSCAN-06`

At `2026-09-11T08:16:49Z`, the current September 2026 Actions billing summary
reported:

- `actions_linux`: 447 minutes, gross `USD 2.682`, discount `USD 2.682`, net
  `USD 0`
- `actions_storage`: `0.066191896` GB-hours, gross `USD 0.0000221`, discount
  `USD 0.0000221`, net `USD 0`

The authenticated GitHub billing UI then visibly reported the account-level
Actions budget as `USD 0`, spent amount `USD 0`, and `Stop usage` as `Yes`.
Current official GitHub documentation listed `ubuntu-24.04` as a standard
GitHub-hosted runner for public repositories and stated that standard runners
are free and unlimited for public repositories.

## Ordered transaction and failure

### Step 1: exact remote main fast-forward

Passed. Remote `main` changed without force from:

`3523e4409ebc53cc1931e3dcaf7d1eea74bb15ab`

to:

`a82a3a64a04ee2d8b60757866c9864cd1b73b54b`

Immediate read-back matched the exact Recovery R4 authority commit.

### Step 2: exact C2 tooling tag and protection

The local unsigned lightweight tag was created and read back as:

- ref: `refs/tags/release-tooling-v1.0.0-c2`
- object type: `commit`
- commit: `faef8435322c9096df09b56969662411f17356ea`
- tree: `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`

The first push command then failed locally on the invalid refspec shown above.
No retry, correction or alternate push was attempted. Steps 2 remainder
through 7 remote execution were stopped.

## Fail-closed read-back

At `2026-09-11T08:21:46Z`, read-only verification established:

- remote `main`:
  `a82a3a64a04ee2d8b60757866c9864cd1b73b54b`
- remote `release-tooling-v1.0.0-c2`: absent
- local `release-tooling-v1.0.0-c2`: present at the exact accepted candidate
- C2 tag ruleset: absent
- C2 workflow run: absent
- C2 artifact: absent
- `PSCAN_RELEASE_C2_GATE`: absent
- artifact/log retention: one day
- tracked, staged and non-ignored untracked working-tree state before writing
  this evidence: clean

No workflow dispatch, Docker/image/container action, dependency acquisition,
build, artifact transfer, signing, OIDC, attestation, deployment, draft,
release or publication occurred. The only remote mutation was the exact
authorized fast-forward of `main`.

## Terminal authority state

Recovery R4 is terminally failed and consumed zero workflow dispatches. Its
no-rerun and no-retry terms authorize no correction, second tag push, tag
ruleset creation, workflow dispatch or promotion. The exact local C2 tag is
preserved as execution evidence and was not moved, deleted or recreated.

Any future attempt requires a new exact owner decision that addresses the
current split state: remote `main` at the R4 authority commit, exact local C2
tag present, remote C2 tag and ruleset absent, and zero C2 workflow runs or
artifacts. PSCAN-06 remains open and unaccepted overall. PSCAN-07 remains
proposed and unselected. PSCAN-08 remains inactive and ineligible. No successor
is selected, activated, claimed or worked.
