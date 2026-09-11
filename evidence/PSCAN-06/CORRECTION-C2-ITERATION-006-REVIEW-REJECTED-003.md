# PSCAN-06 Correction C2 iteration 006 independent review rejection 003

Date: `2026-09-11`
Independent review session: `01a09003-9a4c-7862-967e-972112c9855d`
Approval session: `01a08fae-2ee1-7553-a126-e635b70a34a6`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`

## Outcome

- Result:
  `LOCAL_ACCEPTANCE_REJECTED_MULTI_COLON_VARIABLE_NAME_FALSE_POSITIVE`
- Second refined candidate commit:
  `763730403ae538850e9d806ce9db82c815517b2d`
- Second refined candidate tree:
  `05c0baba5e61edb808b918a2dd6c853f49b8955e`
- Direct second-rejection parent:
  `c4cb2bb15e8b5bd7be708c0500706a5e2ae54e6a`
- Evidence-bearing review parent HEAD:
  `b2d73fe4181629a3663e9a2e534fd701d483cdf4`
- Evidence-bearing review parent tree:
  `4eb409c0b281cc42371ad9be415e4646bb3aa95b`
- Production blob:
  `2fb43bf89df352e9153ba5d7ad23fd59cad77749`
- Production blob SHA-256:
  `49211CA52336BE859CD0C976719B5FB1F209F9E1235F1DA4926C6C21B52905C2`
- Second refined regression blob:
  `a57f93545df3902598e5f0119298d1bf77654935`
- Second refined regression blob SHA-256:
  `3AF57C4519353EA0DD13A844A664B9F3073D2FC88CBE95F159A728E52016C47D`

The second refinement closes the typed, parenthesized and multiple-target
bypasses reported by the first two independent reviews. Its exact predicate,
however, derives the alleged variable name by splitting every true variable
path on `:` and retaining only the last component. Legal braced variables may
contain further colons after a scope or `variable:` prefix. Such variables do
not bind automatic variable `PID`, but the predicate falsely rejects them when
their final component is `PID`. This mandatory non-target false positive
prevents local acceptance. The reviewer did not repair or amend the candidate.

## Admission and preserved history

The review ran in the direct saved checkout at
`C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`, not a
worktree. It verified branch `main`, the exact clean review-parent HEAD and
tree above, its direct candidate parent, the candidate tree and the
candidate's direct second-rejection parent. Git directory and common directory
both resolved to `.git`. The independent review session differs from the
approval and implementation sessions.

Tracked, staged and non-ignored untracked state was clean. All 66 ignored
paths were confined to preserved `graphify-out/**`; none was deleted, moved or
staged.

Both rejected predecessors and their evidence remain immutable:

- first rejected candidate:
  `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`;
- first rejection commit:
  `86d953c75a1b3a68c7be4a6cc6fe35e4d797449b`;
- first rejection-record SHA-256:
  `5986DC219707494A8A9115550F99ABD7F36BD25DE4B08E4D014A7B4DF795E095`;
- first refined and rejected candidate:
  `e6a767aab71db1d3f62063dded379b4701d2cb52`;
- second rejection commit:
  `c4cb2bb15e8b5bd7be708c0500706a5e2ae54e6a`; and
- second rejection-record SHA-256:
  `44FC5CC1F0F79EE5FB8B5F6D9EAA8423B7978E973EC64D9A6C89C8C985831C4D`.

The required second-refinement records independently matched:

- refinement-record SHA-256:
  `E01E9B31B6CA41149A4662A6B6CD86B4A5FB06A87C74E16750114006922F0C5C`;
- third author-validation SHA-256:
  `528DE194A2A90B87F1F46C1F97A19EDE517E43AAFB4D32DD610FAE2AB8CA9258`.

## Exact committed scope and data flow

The candidate changes exactly two authority-permitted paths:
`build/release/test-docker-execution.ps1` and the append-only second
refinement record. The review parent then changes only the three allowed
living task records and third author-validation record. Both ranges pass
`git diff --check`.

`build/release/docker-execution.ps1` remains byte-identical to all three
candidates. Its production blob and SHA-256 are recorded above. Independent
AST inspection confirms that exact `Get-LinuxSessionMembers` contains no
reserved PID assignment. It still binds one integer
`linuxProcessIdentifier` from the parsed stat PID, constructs exact identity
`"${linuxProcessIdentifier}:$startTime"`, and writes the exact two-field
ledger with `PID=$linuxProcessIdentifier` and `StartTime=$startTime`.

## Blocking adversarial finding

The exact committed target walker accepts a non-splatted
`VariableExpressionAst` when `VariablePath.IsVariable` is true, then records:

```powershell
($Node.VariablePath.UserPath -split ':')[-1]
```

`IsVariable` correctly excludes ordinary drive-qualified paths such as
`env:PID` and `function:PID`, but it does not prove that the complete variable
name is automatic `PID`. In braced syntax, the remaining text after a scope or
`variable:` prefix may itself contain colons.

The reviewer extracted the two exact predicate scriptblocks from the
committed regression AST and executed them against independently parsed
cases. Each case below is one valid `VariableExpressionAst`, has
`VariablePath.IsVariable=true`, and the exact predicate reports one reserved
assignment:

| Source | `VariablePath.UserPath` | Predicate count | Clean runtime |
|---|---|---:|---|
| `${variable:env:PID} = 1` | `variable:env:PID` | 1 | exit 0; no collision |
| `${global:env:PID} = 1` | `global:env:PID` | 1 | exit 0; no collision |
| `${local:foo:PID} = 1` | `local:foo:PID` | 1 | exit 0; no collision |
| `${script:foo:PID} = 1` | `script:foo:PID` | 1 | exit 0; no collision |
| `${private:foo:PID} = 1` | `private:foo:PID` | 1 | exit 0; no collision |

Every runtime probe ran in a separate `pwsh -NoProfile -NonInteractive`
process with terminating errors enabled and printed
`EXECUTED_WITHOUT_PID_COLLISION`. Direct `${variable:PID} = 1` and
`${global:PID} = 1` positive controls instead reproduced
`Cannot overwrite variable PID because it is read-only or constant.` and were
correctly classified.

The author-bound control covers `$env:PID` and `${function:PID}`, but none of
the valid multi-colon true-variable names above. The built-in regression
therefore passes despite this exact production-predicate false positive. The
mandatory review instruction requires drive-qualified variables that do not
bind automatic PID and every other non-target control to avoid false
positives. This defect blocks local acceptance even though the production
rename and data flow are correct.

## Passing independent checks

Both implementation scripts parsed with zero errors. Candidate two-path
confinement, review-parent four-path evidence confinement and both exact range
diff checks passed. The complete current-checkout every-byte verifier reported:

```text
Exact source trust PASS commit=b2d73fe4181629a3663e9a2e534fd701d483cdf4 tree=4eb409c0b281cc42371ad9be415e4646bb3aa95b files=369 raw_equal=351 canonical_crlf=18 unexpected_mismatch=0
index=EXACT index_flags=EXACT unsupported_config=0 nonignored_untracked=0 ignored_only_graphify=66
```

The complete available-host local no-Docker matrix passed:

```text
Docker execution iteration-006 PID source regression PASS untyped=REJECT typed=REJECT parenthesized=REJECT multi-target=REJECT nested-wrapper=REJECT member-index-drive-controls=ALLOW rhs-string-comment-ledger=ALLOW data-flow=PASS
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

The full matrix proves the unchanged clean Windows job-object boundary, exact
native source, stream limits, UTF-8 behavior, fixed timeout, descendant
cleanup, empty containment membership, executable replacement, second
invocation and hostile compatible/stale ambient-type rejection. It does not
cure the omitted multi-colon non-target controls. Final inventory found zero
surviving processes whose command line referenced `native-fixture.ps1`, zero
remaining `pscan-c2-iteration-005-*` directories and zero remaining
`pscan-docker-boundary-*` directories. The three local tag identities remained
unchanged.

## Preserved terminal boundary

No Docker command, network connection, remote read or mutation, dependency
acquisition, build, tag action, workflow dispatch or rerun, credential action,
signing, attestation, draft, release, deployment, publication or successor
work occurred.

Recovery R5 remains terminally failed; its sole tag push and workflow dispatch
remain consumed. Actual-Linux, genuine-Docker, dependency, build,
reproducibility and artifact-integrity proof remains open. Candidate
`763730403ae538850e9d806ce9db82c815517b2d` is not independently accepted.
PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains proposed and
unselected; PSCAN-08 remains inactive and ineligible. No successor is selected,
activated, claimed or worked.
