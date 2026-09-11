# PSCAN-06 Correction C2 iteration 006 independent review rejection

Date: `2026-09-11`
Independent review session: `01a09003-9a4c-7862-967e-972112c9855d`
Approval session: `01a08fae-2ee1-7553-a126-e635b70a34a6`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`

## Outcome

- Result: `LOCAL_ACCEPTANCE_REJECTED_REGRESSION_AST_BYPASS`
- Candidate commit:
  `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`
- Candidate tree:
  `de70410df5ec359b85d03262e46c8d2ed5f2d822`
- Direct authority parent:
  `de162e8c347c725a0f041d3c8b0f51df1211c12d`
- Evidence-bearing review parent HEAD:
  `df67e36f5c46a061588877dff72ec246ff803bde`
- Evidence-bearing review parent tree:
  `bff8e857cdff2dce8e3a982f8818d78eead4ff4f`
- Production blob:
  `2fb43bf89df352e9153ba5d7ad23fd59cad77749`
- Production blob SHA-256:
  `49211CA52336BE859CD0C976719B5FB1F209F9E1235F1DA4926C6C21B52905C2`
- Regression blob:
  `b0cc6c6d40259d934e2a4a7b13f8f3db2387b08c`
- Regression blob SHA-256:
  `3876319375E591FE42CFF5DBA4DC7AEBD242CDF270DC3E0872F2D50AC737BD7E`

The production rename is bounded and correct in the exact candidate bytes,
but the required regression does not reject every case-insensitive assignment
to PowerShell automatic variable `PID`. This mandatory coverage gap prevents
local acceptance. The reviewer did not repair or amend the candidate.

## Admission and authority

The review began in the direct saved checkout at
`C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`, not a
worktree. It verified branch `main`, exact clean HEAD and tree above, direct
parent candidate identity, candidate tree and authority parent, and a task ID
different from both prior sessions. Git directory and common directory both
resolved to `.git`.

Tracked, staged and non-ignored untracked state was clean. All 66 ignored
paths were confined to preserved `graphify-out/**`; none was deleted, moved or
staged.

Required hashes independently matched:

- iteration-006 authority SHA-256:
  `231A1E8C964CAB8FD73DA1667EFAAEE62AC86C509B0A425A29DC03C6756B3B36`;
- immutable Recovery R5 failure SHA-256:
  `53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0`;
- author validation SHA-256:
  `FCF9DB55593A809D0DE79D2E7EF268CBE5C931DEF57222907F1CBA28E7ACC784`;
- controlling contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`.

The candidate changes exactly four authority-permitted paths: the production
and regression scripts plus iteration-006 preflight and implementation
evidence. The later author-evidence commit changes only the three permitted
living task records and the author-validation record. The production and
regression blobs at review parent HEAD are byte-identical to the candidate.

## Production inspection

The exact production diff renames only lowercase `$pid` to
`$linuxProcessIdentifier` at its binding and two dependent uses inside
`Get-LinuxSessionMembers`. The parsed type remains `[int]`; start time remains
`[uint64]`; identity remains `<process-id>:<start-time>`; and ledger fields
remain `PID` and `StartTime` with the same source values.

No `/proc` enumeration, record parsing, membership filter, liveness check,
signal, cleanup, stream bound, time bound, executable identity, private
environment, Windows job-object, Linux PID-namespace or Docker-operation
production code changed. The exact diff therefore preserves the previously
accepted local production properties.

## Blocking regression finding

The regression enumerates `AssignmentStatementAst` nodes but treats one as a
reserved `PID` assignment only when `AssignmentStatementAst.Left` is directly
a `VariableExpressionAst`. This is not exhaustive for valid PowerShell
assignment syntax.

The independent in-memory adversarial probe used the case variant:

```powershell
[int]$PiD = 1
```

The probe established all of the following:

- PowerShell parses the assignment without a parse error;
- executing it in a new non-profile process exits nonzero with
  `Cannot overwrite variable PID because it is read-only or constant.`;
- its assignment left side is `ConvertExpressionAst`;
- that left side contains variable expression `PiD`; and
- the candidate's exact reserved-assignment predicate returns count `0` and
  would not reject it.

The mandatory design requires the regression to prove that
`Get-LinuxSessionMembers` has no assignment target named `PID` under
case-insensitive comparison. A valid colliding typed assignment that the guard
does not see defeats that proof. The current production source contains no
such assignment and the direct rename fixes the observed R5 line, but passing
the current exact text is insufficient when the required regression itself
has this syntactic bypass.

## Passing independent checks

Both allowed implementation scripts parsed successfully. The complete current
checkout every-byte verifier reported:

```text
Exact source trust PASS commit=df67e36f5c46a061588877dff72ec246ff803bde tree=bff8e857cdff2dce8e3a982f8818d78eead4ff4f files=363 raw_equal=345 canonical_crlf=18 unexpected_mismatch=0
index=EXACT index_flags=EXACT nonignored_untracked=0 ignored_only_graphify=66
```

Candidate range path confinement and `git diff --check` passed. Independent
blob hashing reproduced both author-recorded script SHA-256 values.

The complete available-host local no-Docker matrix passed:

```text
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

This covered the clean Windows job-object boundary, exact native source,
131072-byte stream limits, UTF-8 cases, 15000 ms timeout behavior, child and
grandchild cleanup, executable replacement, second invocation and compatible
and stale ambient types. Final inventory found zero surviving processes whose
command line referenced `native-fixture.ps1`, zero remaining
`pscan-c2-iteration-005-*` fixture directories and zero remaining
`pscan-docker-boundary-*` directories.

These passing checks do not cure the mandatory regression coverage gap.

## Preserved terminal boundary

No Docker command, network connection, remote read or mutation, dependency
acquisition, build, tag action, workflow dispatch or rerun, credential action,
signing, attestation, draft, release, deployment, publication or successor
work occurred.

Recovery R5 remains terminally failed; its sole tag push and workflow dispatch
remain consumed. Actual-Linux, genuine-Docker, dependency, build,
reproducibility and artifact-integrity proof remains open. Candidate
`8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4` is not independently accepted.
PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains proposed and
unselected; PSCAN-08 remains inactive and ineligible. No successor is selected
or activated.
