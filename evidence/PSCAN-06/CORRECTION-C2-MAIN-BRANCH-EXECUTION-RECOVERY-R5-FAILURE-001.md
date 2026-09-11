# PSCAN-06 Correction C2 main-branch execution Recovery R5 failure 001

## Outcome

- Task: `PSCAN-06`
- Correction: `C2`
- Recovery: `R5`
- Gate: build-only actual-Linux and genuine-Docker proof
- Execution session: `01a08fae-2ee1-7553-a126-e635b70a34a6`
- Date: `2026-09-11`
- Terminal result: `FAIL_CLOSED_WORKFLOW_BUILD_NATIVE_MATRIX`
- Tag pushes consumed by Recovery R5: one
- Workflow dispatches consumed by Recovery R5: one
- Workflow run: `34582887399`, attempt `1`
- Workflow conclusion: `failure`
- Artifact created or downloaded by Recovery R5: none
- Local Docker action: none
- Signing, attestation, draft, release or publication action: none
- Successor action: none

Recovery R5 passed the complete mandatory fresh preflight, pushed the preserved
local C2 tooling tag using the sole authorized literal refspec, read back its
exact remote identity, and created and verified the exact no-bypass tag
ruleset. It then dispatched the exact build-only workflow once.

The exact `build` job began on standard `ubuntu-24.04` and passed its first six
steps, including source trust. It failed at step 7, `Prove private admission
and bounded native process cases without Docker`, before the workflow reached
host cache proof, image admission, genuine Docker operations, dependency
acquisition, either build, byte comparison or artifact transfer.

The no-rerun and no-repair boundary made the failure terminal. The workflow was
not rerun or repaired. No second tag push or dispatch occurred.

## Exact execution context

- Saved project root:
  `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`
- Branch: `main`
- Starting HEAD:
  `fdd105492011e7acd9d05b63378642ec54a03289`
- Starting tree:
  `62a8439477f59957f60a87bab00cfc33c6e68df8`
- Direct parent:
  `6fbeb633bf373d2f64d95d0ff93c328631a99a05`
- Model: `gpt-daybreak-blue-latest`
- Reasoning effort: `xhigh`
- Fresh task source: agent-created task targeting the saved local project
  directly, not a worktree, clone or detached checkout
- Approval session:
  `01a0737f-d3fe-7c93-9d3e-db1f631c3df9`
- Execution session differs from approval session: yes
- Preflight tracked, staged and non-ignored untracked state: clean
- Preserved ignored paths: 66 paths, all under `graphify-out/**`

Task metadata proved the required model and reasoning effort before the first
repository action. The exact root, branch, HEAD, tree, sole parent and clean
state were then verified before the first remote action.

## Local source-trust and authority preflight

Complete index, path, mode, flag and every-byte comparison reported:

- tracked paths: 358
- raw-equal to committed blob: 340
- canonical whole-file LF-to-CRLF projections: 18
- unexpected tracked-byte mismatches: zero
- non-ignored untracked paths: zero
- unsupported Git source-trust configuration: none
- unexpected ignored paths: zero

Required identities and hashes passed:

- Recovery R5 authority SHA-256:
  `FE5F7E76A45BCE3A568D71EEF9FDE7A5AFE9BDF7B79B2D602715C0ABABA2A2D2`
- Recovery R4 failure SHA-256:
  `5F694FF128515A2B79A2BFA500F299FECAB8836763365E9D774DA5D847D68BEE`
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
- accepted candidate:
  `faef8435322c9096df09b56969662411f17356ea`
- accepted candidate tree:
  `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`
- required evidence HEAD:
  `b6d341284cf63baa0502fcba319ea8cca3c7eb3a`
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
that variable absent. Repository Actions permissions remained enabled with the
selected-actions policy, required SHA pinning and default read permission.
Artifact and log retention remained one day.

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
- preserved remote `main`:
  `a82a3a64a04ee2d8b60757866c9864cd1b73b54b`
- local C2 tag type: lightweight, resolving directly to the accepted candidate
- remote C2 tag: absent
- C2 exact-name ruleset: absent
- prior C2 workflow runs: zero
- prior C2 artifacts: zero
- repository releases and deployments: zero
- exact repository variable `PSCAN_RELEASE_C2_GATE`: absent
- sole open task: `PSCAN-06`
- `PSCAN-07`: proposed and unselected
- `PSCAN-08`: inactive and ineligible

The locked `v1.0.0` and `release-tooling-v1.0.0-c1` tags resolved to their
recorded commits and retained active exact-name update/deletion-denial
rulesets with no bypass actors. Both halves of the sole R5 refspec passed
`git check-ref-format`.

Immediately before mutation, the September 2026 Actions billing summary
reported both `actions_linux` and `actions_storage` with net amount `USD 0`.
The authenticated GitHub billing UI visibly reported the account-level Actions
budget as `USD 0`, spent amount `USD 0`, and `Stop usage` as `Yes`. Current
official GitHub documentation listed `ubuntu-24.04` as a standard hosted runner
for public repositories and stated that standard runners are free and
unlimited for public repositories.

## Ordered transaction

### Step 1: exact existing-tag push

Passed. The already-existing local lightweight tag was not created, recreated,
moved or deleted. It was pushed once without force using the sole literal
refspec:

```text
refs/tags/release-tooling-v1.0.0-c2:refs/tags/release-tooling-v1.0.0-c2
```

Immediate remote read-back established:

- ref: `refs/tags/release-tooling-v1.0.0-c2`
- object type: `commit`
- commit: `faef8435322c9096df09b56969662411f17356ea`
- tree: `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`
- annotated tag object: absent
- remote `main`: unchanged at `a82a3a64a04ee2d8b60757866c9864cd1b73b54b`

### Step 2: exact C2 tag ruleset

Passed. Repository ruleset `22895383`, named
`Lock release-tooling-v1.0.0-c2 tag`, was created and read back as:

- target: `tag`
- enforcement: `active`
- sole include: `refs/tags/release-tooling-v1.0.0-c2`
- excludes: none
- rules: `deletion`, `update`
- bypass actors: none
- current user can bypass: `never`

### Step 3: one-day retention

Passed. Artifact and log retention read back as exactly one day. No setting
change was required or made.

### Step 4: single workflow dispatch

Passed. Exactly one dispatch was issued from tag
`release-tooling-v1.0.0-c2` with only:

```text
product_source_revision=a13c28fe7273bc8dc6545f97966a02889524eb4c
release_tooling_revision=faef8435322c9096df09b56969662411f17356ea
release_version=v1.0.0
```

The resulting run was:

- URL:
  `https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/actions/runs/34582887399`
- workflow ID: `349829683`
- repository ID: `1355442997`
- event: `workflow_dispatch`
- ref: `release-tooling-v1.0.0-c2`
- SHA: `faef8435322c9096df09b56969662411f17356ea`
- run attempt: `1`
- actor and triggering actor: `ThameeraDananjaya`

The workflow's immutable dual-identity step passed, and the `build` job-level
condition admitted the exact repository, tag, tooling SHA, workflow ref,
workflow SHA and three fixed inputs.

## Workflow failure

The sole executing job was `build`, job ID `103210318963`, on runner group
`GitHub Actions`, runner `GitHub Actions 1000000344`, with label
`ubuntu-24.04`. It started at `2026-09-11T09:11:01Z` and completed with
conclusion `failure` at `2026-09-11T09:11:33Z`.

These required steps passed before the failure:

1. `Set up job`
2. `Validate immutable dual-identity invocation`
3. `Check out exact tooling revision without persisted credentials`
4. `Verify locked product tag and tree`
5. `Materialize and verify exact committed release entrypoints`
6. `Reject hostile Git index and checkout states before build output`

Step 7, `Prove private admission and bounded native process cases without
Docker`, failed with exit code `1`. The failure originated at committed
`build/release/test-docker-execution.ps1` line 289 when its clean isolated
native matrix returned a PowerShell `WriteError`. The exact reported cause was:

```text
Cannot overwrite variable PID because it is read-only or constant.
```

The failing isolated path attempted a lowercase `$pid` assignment while
PowerShell variable names are case-insensitive and automatic `$PID` is
read-only. This record reports the observed committed failure only; Recovery R5
authorizes no repair.

All later host-cache, image-admission, Docker, dependency-acquisition, build,
comparison and artifact-transfer steps were skipped. The
`sign-attest-and-draft` job, ID `103210478044`, was skipped with zero steps and
no runner assigned.

## Final read-back

At `2026-09-11T09:13:30Z`, read-only verification established:

- run `34582887399`: completed, conclusion `failure`, attempt `1`
- C2 workflow-run count: one, exactly this run
- run artifacts: zero
- repository artifacts: zero
- releases: zero
- deployments: zero
- remote `main`: unchanged at
  `a82a3a64a04ee2d8b60757866c9864cd1b73b54b`
- remote C2 tag: exact accepted candidate
  `faef8435322c9096df09b56969662411f17356ea`
- C2 ruleset `22895383`: active, exact-name update/deletion denial, no bypass
- `PSCAN_RELEASE_C2_GATE`: absent
- artifact/log retention: one day
- local tracked, staged and non-ignored untracked state before writing this
  evidence: clean

The post-run Actions billing summary reported:

- `actions_linux`: 448 minutes, gross `USD 2.688`, discount `USD 2.688`, net
  `USD 0`
- `actions_storage`: `0.066584604` GB-hours, gross `USD 0.000022231`, discount
  `USD 0.000022231`, net `USD 0`

The refreshed account-level Actions budget UI still reported `USD 0` spent,
`USD 0` budget and `Stop usage` as `Yes`.

No artifact existed to download or independently verify. No local Docker,
remote-main push, signing, OIDC, attestation, deployment, draft, release,
publication, rerun, repair, second dispatch, second tag push or successor
action occurred.

## Terminal authority state

Recovery R5 is terminally failed. It consumed its sole tag push and sole
workflow dispatch and authorizes no rerun, repair, retry, second push, second
dispatch or promotion. The exact remote C2 tag and its active no-bypass
protection ruleset are preserved.

The actual-Linux execution context was reached, but the genuine-Docker,
dependency, build, reproducibility and artifact-integrity proof remains open.
PSCAN-06 remains open and unaccepted overall. Any future attempt or correction
requires a new exact owner decision from this terminal state. PSCAN-07 remains
proposed and unselected. PSCAN-08 remains inactive and ineligible. No successor
is selected, activated, claimed or worked.
