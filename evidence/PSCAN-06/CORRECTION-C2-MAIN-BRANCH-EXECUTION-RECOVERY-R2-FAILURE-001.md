# PSCAN-06 Correction C2 main-branch execution Recovery R2 failure 001

## Outcome

The owner-approved Recovery R2 attempt failed closed during the mandatory
current preflight. The fresh execution task reached the exact saved-project
`main` context and read the complete current PSCAN-06 authority chain, but
`gh auth status --hostname github.com` returned exit code `1`: the intended
account was active and its stored token was invalid.

Required execution-context item 6 and mandatory preflight item 4 therefore did
not pass. The task stopped before any GitHub repository API query, remote
mutation, workflow dispatch, Docker action, build or artifact action. This is a
terminal Recovery R2 failure, not a pass or a reusable preflight. Recovery R2
authorizes no retry.

## Exact execution and authority identity

- Evidence time (UTC): `2026-09-07T21:26:51.4957606Z`
- Execution session: `01a07dbf-b19b-71d1-b40c-617f396b1c64`
- Approval session: `01a0737f-d3fe-7c93-9d3e-db1f631c3df9`
- Model: `gpt-daybreak-blue-latest`
- Reasoning effort: `xhigh`
- Exact repository root:
  `C:\OFFICE-DATA\Projects\Ongoing\project-agnostic-secret-scanner`
- Branch: `main`
- Recovery R2 authority commit:
  `07de8d6c5fc5f094f0b6c213c854176e6acf9ec0`
- Recovery R2 authority tree:
  `b3dd2f2371bd7954d4b5a1104e1bed9323291a64`
- Recovery R2 authority parent:
  `31997817628cbcff8fa921ec24906f24cbd960e1`
- Locked evidence candidate:
  `faef8435322c9096df09b56969662411f17356ea`
- Locked evidence candidate tree:
  `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`
- Locked evidence head:
  `b6d341284cf63baa0502fcba319ea8cca3c7eb3a`
- Locked evidence-head tree:
  `983028703e229abb082938d0e4417c506aac0d8c`
- Locked evidence-head parent:
  `faef8435322c9096df09b56969662411f17356ea`
- Recovery R2 authority SHA-256:
  `CE16C537123F318485E3C1A7B91654B3593F4250EB0B358E3E5CBAE43FBEA762`
- Original proof-gate authority SHA-256:
  `2A84E7109F151E55399870FB36D64AEDB2213E123976779CD1E9648AC1B9DAC2`
- Controlling PASS contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`

The execution session differs from the approval session. The model and effort
were read from the execution task's own session metadata, not inferred from
prose.

## Read-only local preflight evidence

- `git rev-parse --show-toplevel` matched the exact required saved-project
  root.
- `git branch --show-current` returned `main`.
- `HEAD`, its tree and its single parent matched the Recovery R2 authority.
- The tracked working tree, index and untracked-file set were clean.
- The ignored `graphify-out/` directory was preserved and not regenerated or
  modified.
- The current PSCAN-06 reading map and all 88 paths enumerated or directly
  linked by it were read completely before the gate attempt, followed by the
  Recovery R2 authority.
- Origin was the intended HTTPS GitHub repository URL.
- The local Correction C2 tooling tag was absent.
- The two locked local release tags resolved to their recorded commit and tree
  identities.
- GitHub CLI resolved to `C:\WINDOWS\system32\gh.exe`.

These local facts do not prove any current GitHub repository, account, access,
ruleset, retention, workflow, billing, Actions or run fact.

## Terminal authentication failure

The first current GitHub authentication command was:

```text
gh auth status --hostname github.com
```

It returned exit code `1` and reported:

```text
github.com
  X Failed to log in to github.com account ThameeraDananjaya (default)
  - Active account: true
  - The token in default is invalid.
```

No credential value was read, printed, changed or handled. No login or logout
was attempted. Because current authentication failed, the task did not enter
the remaining remote preflight and did not infer any remote fact from local
state, historical evidence or the active account name.

## Actions not entered

- No GitHub repository API query or other authenticated remote read occurred.
- No fast-forward push, tag creation, ruleset change, retention change or other
  remote mutation occurred.
- No workflow was dispatched, rerun or cancelled. Recovery R2 consumed zero
  workflow dispatches; the earlier failed attempts also consumed zero.
- No local or hosted Docker action, dependency download, build, artifact upload
  or artifact download occurred.
- No signing, attestation, draft, publication, release, paid resource or spend
  occurred.
- No PSCAN-07, PSCAN-08 or successor work occurred.

## Status and required owner decision

Recovery R2 is terminally failed and authorizes no retry. PSCAN-06 remains open
and unaccepted overall. PSCAN-07 remains proposed and unselected. PSCAN-08
remains inactive and ineligible.

Before any further proof-gate attempt, the owner must restore valid GitHub
authentication outside this failed authority and issue a new exact bounded
owner decision. This record does not create that authority.
