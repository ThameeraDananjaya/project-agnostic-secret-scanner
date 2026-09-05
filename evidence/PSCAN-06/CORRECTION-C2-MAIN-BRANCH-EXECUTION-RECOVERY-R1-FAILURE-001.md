# PSCAN-06 Correction C2 main-branch execution Recovery R1 failure 001

## Outcome

The genuinely fresh saved-project execution session failed closed at mandatory
Recovery R1 preflight item 1 before any GitHub or repository-remote request,
push, tag, ruleset, setting change, workflow dispatch, Docker action, artifact
action or other remote mutation.

The checkout is at the exact required saved-project root, on branch `main`, at
the exact Recovery R1 authority commit and tree, with the required direct
parent. Tracked, staged and non-ignored untracked Git status is clean, and the
only ignored material is the expressly preserved `graphify-out/**` boundary.
However, exact every-byte inspection found two tracked files whose working-tree
bytes are neither raw-equal to their committed blobs nor a canonical LF-to-CRLF
projection. Mandatory preflight item 1 therefore fails as a whole. The session
did not repair, normalize, checkout, stage or otherwise modify either file.

PSCAN-06 remains open and unaccepted overall. This record authorizes no retry,
repair, rerun, signing, attestation, draft, publication or successor work.

## Exact session and authority identity

- Execution session: `01a07376-04a7-7aa1-8421-7d4e82a13465`
- Runtime model: `gpt-daybreak-blue-latest`
- Reasoning effort: `xhigh`
- Evidence time: `2026-09-05T21:33:34.5024805Z`
- Repository root:
  `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`
- Branch: `main`
- Recovery R1 authority commit:
  `ca996eed20294a1d36822dd453005a333218279a`
- Recovery R1 authority tree:
  `5b385603db85ddd0870f82250014bd8084534933`
- Recovery R1 direct parent:
  `268e8f339318a141d059115c9da0d19a01448aaa`
- Accepted tooling candidate:
  `faef8435322c9096df09b56969662411f17356ea`
- Accepted tooling tree:
  `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`
- Required evidence HEAD:
  `b6d341284cf63baa0502fcba319ea8cca3c7eb3a`
- Required evidence tree:
  `983028703e229abb082938d0e4417c506aac0d8c`
- Evidence HEAD direct parent: exact accepted tooling candidate
- Controlling contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Original proof-gate authority SHA-256:
  `2A84E7109F151E55399870FB36D64AEDB2213E123976779CD1E9648AC1B9DAC2`
- Preserved pre-dispatch failure SHA-256:
  `7EB1CB03E85A933C1BCFEB435FA85556170AE21F1112CFFA55041C2A733E6A5D`

The session verified its model and reasoning effort from its own local session
metadata before repository inspection. It then read `AGENTS.md`, the complete
current 85-file PSCAN-06 reading map, the controlling contract, traceability,
living task/tracker state, the original proof-gate authority, preserved failure
record, Recovery R1 authority and all authority linked by that map.

## Pre-evidence local state

- `git diff --quiet`: exit `0`
- `git diff --cached --quiet`: exit `0`
- Non-ignored untracked paths: `0`
- Ignored paths outside `graphify-out/**`: `0`
- Tracked paths enumerated from the exact authority tree: `349`
- Tracked paths raw-equal to their committed blob: `331`
- Tracked paths with canonical CRLF projection: `16`
- Tracked paths with noncanonical working-tree bytes: `2`
- Index paths with assume-unchanged, skip-worktree or unsupported flags: `0`
- Index paths with filesystem-monitor-valid or unsupported flags: `0`
- Unsupported Git shortcut, partial-clone or promisor configuration: `0`
- Git `core.autocrlf`: `true`

The exact candidate workflow blob equals the Recovery R1 workflow blob:
`3aa42627628b8b5298d854d29b3880cb19dc36ff`. Its worktree SHA-256 equals the
recorded value
`C5F40F1B32E87C005FE33EE607AF7E3D19EE4F0173C21619E21158C31AA0BDB4`.
Static inspection also confirmed that the build job uses only `contents: read`,
its two actions are full-SHA pinned, and its single artifact declares
`retention-days: 1` and `compression-level: 0`. These partial local facts do not
promote mandatory preflight item 1 or the complete preflight to pass.

## Terminal source-trust failure

The repository source-trust contract accepts a tracked working-tree file only
when its raw Git blob identity equals the committed tree entry, or, where
expressly allowed for a Windows checkout, every LF byte is projected to CRLF
without any mixed or other byte difference. The following two files fail both
admissible cases while Git's text clean filter maps them back to their expected
committed blobs:

### `build/gitleaks/build.ps1`

- Expected committed blob:
  `7d66c57c27767464f8a084211f37db50a85e4702`
- Raw working-tree blob:
  `be1241c8a2ba77f4140270750264d1dd54100893`
- Git-cleaned blob:
  `7d66c57c27767464f8a084211f37db50a85e4702`
- Working-tree bytes: `2933`
- CRLF sequences: `55`
- Bare LF bytes: `4`
- Bare CR bytes: `0`
- Git attribute: `text=auto`

### `build/gitleaks/collect-licenses.ps1`

- Expected committed blob:
  `6b5fd1e520c0d8be694461397e162b0447164ac8`
- Raw working-tree blob:
  `ef95129091f5c9c33900ca20423c42d1eaa4b645`
- Git-cleaned blob:
  `6b5fd1e520c0d8be694461397e162b0447164ac8`
- Working-tree bytes: `2479`
- CRLF sequences: `45`
- Bare LF bytes: `3`
- Bare CR bytes: `0`
- Git attribute: `text=auto`

Ordinary Git status is clean because Git's text filter normalizes both mixed-EOL
files before comparison. That clean status does not satisfy the stronger raw
source-trust requirement. No inference, repair or waiver is permitted.

## Actions not entered

Because mandatory preflight item 1 failed, the session did not enter GitHub or
repository-remote preflight items 2 through 6. It did not inspect or change
remote refs, variables, Actions policy, runners, billing, budgets, stop-usage
control, retention, tags, rulesets, workflow runs, artifacts, deployments or
releases. It did not push `main`, create or push the C2 tag, dispatch or rerun a
workflow, execute Docker, download an artifact, sign, attest, draft or publish.
Those remote facts remain unproved by this session, not inferred as absent.

An official OpenAI product-documentation lookup was attempted before repository
inspection solely to corroborate the locally recorded runtime-model identity,
as the owner required. It was not a GitHub, repository or proof-gate remote
preflight and did not consume the one-dispatch authority.

## Required next owner decision

The owner must decide whether to issue a new exact authorization for a fresh
saved-project execution after separately resolving the two noncanonical
working-tree projections. Nothing may repair or continue automatically from
this failure record. PSCAN-07 remains proposed and unselected. PSCAN-08 remains
inactive, unselected and ineligible. No successor was selected or activated.
