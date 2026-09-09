# PSCAN-06 Correction C2 main-branch execution Recovery R3 failure 001

## Outcome

The owner-approved Recovery R3 attempt failed closed during the mandatory
current remote preflight. The genuinely fresh saved-project `main` task passed
the exact local execution, source-trust, authority, workflow-byte and GitHub
authentication checks. Current GitHub account billing, usage and stop-usage
evidence could not be read with the authenticated host-level credential,
however. Every applicable user billing endpoint returned HTTP `404` and the
GitHub CLI stated that the operation requires the missing `user` scope.

Mandatory preflight item 4 therefore did not pass. The task stopped before any
remote mutation, workflow dispatch, Docker action, build or artifact action.
This is a terminal Recovery R3 failure, not a pass or a reusable preflight.
Recovery R3 authorizes no retry.

## Exact owner authority and execution identity

- Owner approval received: `2026-09-09`
- Recovery: `R3 from Recovery-R2 failure 001`
- Execution session: `01a08461-59eb-7721-b1f3-f11721e36bfd`
- Exact repository root:
  `C:\OFFICE-DATA\Projects\Ongoing\project-agnostic-secret-scanner`
- Branch: `main`
- Required and observed evidence commit:
  `a4f5b85a3ac63ea31b4bb911be3445f2223b6db1`
- Evidence-commit tree:
  `df865f84e7d207c728bc5bae821cbb1b9f3127b5`
- Evidence-commit parent / Recovery R2 authority commit:
  `07de8d6c5fc5f094f0b6c213c854176e6acf9ec0`
- Recovery R2 failure record SHA-256 supplied by the owner and reverified:
  `CACDEE7FA590FEA4714A921ECC000AF6F008F4E00515DEDE91C1BEE3E37BAC91`
- Maximum spend: `USD 0`
- Workflow dispatch limit: at most one; consumed by Recovery R3: zero
- Authorized artifact boundary: one unsigned, uncompressed workflow artifact
  retained for one day only
- Local Docker, signing, attestation, draft, publication, rerun and successor
  work: forbidden

The approval preserved the complete original proof-gate and Recovery R2
identity and safety boundary. It authorized no credential-scope change or
authentication refresh.

## Read-only local preflight evidence

- Exact root, `main`, HEAD, tree and parent matched the owner-approved evidence
  identity.
- The complete 354-path tree and index matched by path, mode and object.
- Working-tree bytes were exact for 336 tracked files and canonical whole-file
  LF-to-CRLF projections for 18 tracked files, with zero mismatches.
- Staged and non-ignored untracked state was clean. All 65 ignored paths were
  confined to preserved `graphify-out/**` material, which remained excluded
  from product authority and staging.
- Accepted tooling candidate
  `faef8435322c9096df09b56969662411f17356ea` resolved to tree
  `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`.
- Locked evidence HEAD `b6d341284cf63baa0502fcba319ea8cca3c7eb3a`
  resolved to tree `983028703e229abb082938d0e4417c506aac0d8c`
  with direct parent
  `faef8435322c9096df09b56969662411f17356ea`.
- The accepted workflow SHA-256 was exactly
  `C5F40F1B32E87C005FE33EE607AF7E3D19EE4F0173C21619E21158C31AA0BDB4`.
  Its build job retained `contents: read`, full-SHA action pins, standard
  `ubuntu-24.04`, one artifact, `retention-days: 1` and
  `compression-level: 0`.
- The signing job remained gated only by
  `PSCAN_RELEASE_C2_GATE=PSCAN-06-C2-SIGNING-APPROVED`.

## Read-only GitHub preflight evidence

- `gh auth status --hostname github.com` passed for active account
  `ThameeraDananjaya` using the host keyring. The reported OAuth scopes were
  `gist`, `read:org`, `repo` and `workflow`.
- Authenticated account ID was `50274860`.
- Repository identity was
  `ThameeraDananjaya/project-agnostic-secret-scanner`, numeric ID
  `1355442997`, public, unarchived, not a fork, with default branch `main`.
- Origin was the exact intended HTTPS URL.
- Actions was enabled with selected actions only, required full-SHA pinning,
  default workflow permission `read` and only the four exact committed action
  pins admitted.
- Repository variable `PSCAN_RELEASE_C2_GATE` was absent. The only repository
  variable was the distinct historical `PSCAN_RELEASE_GATE`.
- Locked tags `v1.0.0` and `release-tooling-v1.0.0-c1` resolved to their exact
  recorded commits. Their exact-name tag rulesets were active, denied update
  and deletion, and had no bypass actors.
- `release-tooling-v1.0.0-c2` was absent locally and remotely.
- Remote `main` was
  `3523e4409ebc53cc1931e3dcaf7d1eea74bb15ab`, a proven ancestor of the
  approved evidence commit, so a non-force fast-forward was possible.
- No C2 workflow run or artifact existed. The repository had no current
  Actions artifacts, release or deployment.

These passed read-only facts did not cure the billing evidence failure and
did not authorize mutation.

## Terminal billing and stop-usage evidence failure

The current authenticated CLI attempted the applicable user-level read-only
billing endpoints, including:

```text
GET /users/ThameeraDananjaya/settings/billing/usage
GET /users/ThameeraDananjaya/settings/billing/usage/summary
GET /users/ThameeraDananjaya/settings/billing/actions
```

Every request returned HTTP `404`, and `gh` reported:

```text
This API operation needs the "user" scope.
```

Current official GitHub documentation states that standard GitHub-hosted
runners, including `ubuntu-24.04`, are free and unlimited for public
repositories. It separately states that Actions artifact storage consumes a
plan allowance, that overage can be billed, and that a budget with stop-usage
enabled is the control against further metered usage. Therefore public
visibility and the standard runner were not sufficient to infer exact
zero-spend artifact storage. Current account usage, allowance, payment and
stop-usage state remained permission-denied and untrusted.

The governing authority says that any missing, permission-denied or otherwise
unknown billing or storage fact stops the gate before mutation. No credential
refresh or scope expansion was authorized or attempted.

## Actions not entered

- Remote `main` was not pushed or otherwise changed.
- No C2 tag or tag ruleset was created, pushed or changed.
- No Actions retention, repository, environment, branch, release or security
  setting was changed.
- No workflow was dispatched, rerun or cancelled. Recovery R3 consumed zero
  workflow dispatches.
- No local or hosted Docker action, dependency download, build, artifact upload
  or artifact download occurred.
- No signing, OIDC request, attestation, draft, publication, release, paid
  resource or spend occurred.
- No PSCAN-07, PSCAN-08 or successor work occurred.

## Status and required owner decision

Recovery R3 is terminally failed and authorizes no retry. PSCAN-06 remains open
and unaccepted overall. PSCAN-07 remains proposed and unselected. PSCAN-08
remains inactive and ineligible.

Any further attempt requires a new exact owner decision after current account
billing, artifact-storage allowance and stop-usage evidence can be read and
trusted. This record does not create that authority and does not authorize a
credential-scope change.
