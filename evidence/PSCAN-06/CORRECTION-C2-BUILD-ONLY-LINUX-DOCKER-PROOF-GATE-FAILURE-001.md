# PSCAN-06 Correction C2 build-only Linux and Docker proof gate failure 001

## Outcome

The fresh execution session failed closed at mandatory local preflight item 1
before any remote preflight, network access, push, tag, ruleset, setting change,
workflow dispatch, Docker action, artifact action or other remote mutation.

The authority requires the exact clean gate-authority checkout on branch
`main`. This execution checkout is detached: Git reports
`branch.head (detached)` and `HEAD (no branch)`. The session did not switch,
repair or reinterpret the checkout. The remaining preflight and authorized
transaction were not entered.

PSCAN-06 remains open and unaccepted overall. This record authorizes no retry,
repair, rerun, signing, attestation, draft, publication or successor work.

## Exact session and authority identity

- Execution session: `01a072e4-50f8-7842-8c7d-28bb034a773f`
- Runtime model: `gpt-daybreak-blue-latest`
- Reasoning effort: `xhigh`
- Gate approval session: `01a072d7-40e1-76a0-ac22-995d48e94d3e`
- Evidence time: `2026-09-05T18:50:41.026Z`
- Session scope: only the owner-approved PSCAN-06 Correction C2 build-only
  actual-Linux and genuine-Docker proof gate
- Gate-authority commit:
  `268e8f339318a141d059115c9da0d19a01448aaa`
- Gate-authority tree:
  `e1667892d0c05d9a8db8bac65e580fff54a876ea`

The session read the required current entry authorities completely:
`AGENTS.md`, `docs/spec/PASS-OUTCOME-SPEC-001.md`,
`docs/tasks/PSCAN-06.md`, `docs/tasks/PSCAN-06-READING-MAP.md`,
`docs/tasks/TRACKER.md`, and
`evidence/PSCAN-06/CORRECTION-C2-BUILD-ONLY-LINUX-DOCKER-PROOF-GATE-AUTHORITY.md`.

## Read-only local evidence

- Repository root:
  `C:/Users/ITDan/.codex/worktrees/b31c/project-agnostic-secret-scanner`
- Required branch: `main`
- Actual current branch: none; detached HEAD
- Git porcelain branch facts:
  `branch.oid 268e8f339318a141d059115c9da0d19a01448aaa` and
  `branch.head (detached)`
- Current HEAD:
  `268e8f339318a141d059115c9da0d19a01448aaa`
- Current tree:
  `e1667892d0c05d9a8db8bac65e580fff54a876ea`
- Current HEAD direct parent:
  `b6d341284cf63baa0502fcba319ea8cca3c7eb3a`
- Required evidence HEAD tree:
  `983028703e229abb082938d0e4417c506aac0d8c`
- Required evidence HEAD direct parent:
  `faef8435322c9096df09b56969662411f17356ea`
- Accepted candidate tree:
  `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`
- Pre-evidence tracked diff: none (`git diff --quiet` exit `0`)
- Pre-evidence staged diff: none (`git diff --cached --quiet` exit `0`)
- Pre-evidence tracked, untracked and ignored status entries: none
- Local `refs/heads/main` resolves to the required authority commit, but Git's
  worktree inventory proves that branch is checked out in the separate primary
  worktree `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`.
  That does not make this fresh execution checkout a `main` checkout.

## Terminal preflight failure

Mandatory preflight item 1 requires the exact repository root, `main` branch,
gate-authority commit and tree in an exact clean source-trust-compliant
checkout. The commit, tree, evidence relation and candidate identity resolve,
but the required branch condition is false. The item therefore fails as a
whole and cannot be promoted from partial evidence to pass.

Because the first mandatory item failed, this session performed no GitHub or
other remote preflight and made no remote request. It did not inspect or change
remote refs, tags, rulesets, variables, Actions settings, retention, billing,
runners, prior runs or artifacts. It did not dispatch the workflow or execute
Docker. Those facts remain unproved by this session, not inferred as absent.

## Required next owner decision

The owner must decide whether to issue a new exact authorization for a genuinely
fresh execution session whose current checkout itself is on exact branch
`main` at the gate-authority commit and tree. Nothing may retry, repair or
continue automatically from this failure record. PSCAN-07 remains proposed and
unselected; PSCAN-08 remains inactive and ineligible.
