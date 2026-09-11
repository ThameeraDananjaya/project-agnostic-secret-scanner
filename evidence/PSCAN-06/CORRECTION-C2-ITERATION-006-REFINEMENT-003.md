# PSCAN-06 Correction C2 iteration 006 regression refinement 003

Date: `2026-09-11`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`
Starting review-evidence commit:
`cf6b61d73910d7525a2213877654dfb900fe2cb2`
Starting review-evidence tree:
`862019898de01ac81edb3814495fd754b58647ee`

## Preserved rejection history and production result

Independent review rejected second refined candidate
`763730403ae538850e9d806ce9db82c815517b2d`, tree
`05c0baba5e61edb808b918a2dd6c853f49b8955e`, with result
`LOCAL_ACCEPTANCE_REJECTED_MULTI_COLON_VARIABLE_NAME_FALSE_POSITIVE`.
The immutable third review record is
`CORRECTION-C2-ITERATION-006-REVIEW-REJECTED-003.md`, SHA-256
`F21FD594E757747EE5D3F9865FBE2C33D1FECE7AF09B28AA50778D355A183373`.

Candidates `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`,
`e6a767aab71db1d3f62063dded379b4701d2cb52` and
`763730403ae538850e9d806ce9db82c815517b2d`, their author evidence and all
three independent rejection records remain unchanged history. This refinement
does not modify `build/release/docker-execution.ps1`. Its committed blob
remains `2fb43bf89df352e9153ba5d7ad23fd59cad77749`, SHA-256
`49211CA52336BE859CD0C976719B5FB1F209F9E1235F1DA4926C6C21B52905C2`.

## VariablePath semantics and exact correction

Only `build/release/test-docker-execution.ps1` changes in test source. Parser
and reflection probes against the installed PowerShell runtime established:

- public `UserPath` preserves the complete path, including every colon;
- public `IsVariable` is true for unqualified and variable/scope-qualified
  variables, including legal multi-colon names after those qualifiers;
- public `IsDriveQualified` and `DriveName` identify provider/drive paths;
- public `IsUnqualified`, `IsUnscopedVariable`, `IsGlobal`, `IsScript`,
  `IsLocal` and `IsPrivate` classify the recognized variable scope;
- `UnqualifiedPath` exists but its getter is assembly-internal, not a supported
  public property available to the regression script.

For automatic PID, direct `$PID` has `IsUnqualified=true` and `UserPath=PID`;
`${variable:PID}` has `IsUnscopedVariable=true` and
`UserPath=variable:PID`; and the four explicit scope forms set their matching
public flag with exact whole paths `global:PID`, `script:PID`, `local:PID` or
`private:PID`. In contrast, `${global:env:PID}` is global but its exact
`UserPath` is `global:env:PID`, while `${local:foo:PID}` is local with exact
`UserPath=local:foo:PID`. Ordinary `env:` and `function:` targets are drive
qualified and are not variable paths.

The production-bound predicate now requires both the matching public scope
classification and case-insensitive equality of the entire expected
`UserPath`. It never splits on `:` or infers identity from a final path
segment. The existing assignment-left walker remains unchanged: typed,
attributed, parenthesized and recursively nested multiple targets are walked;
member/index targets terminate the walk; unsupported shapes fail closed.

## Executable predicate and runtime proof

Parser/predicate cases retain every earlier direct, typed, attributed,
parenthesized, nested, background and multiple-target positive. They add all
unqualified, `variable:`, global, script, local and private automatic PID
spellings, including braced forms. Negative cases retain member/index,
right-hand-side, property, string, comment and ledger controls and add:

- `${variable:env:PID} = 1`;
- `${global:env:PID} = 1`;
- `${local:foo:PID} = 1`;
- `${script:foo:PID} = 1`;
- `${private:foo:PID} = 1`;
- `$env:PID`, `${env:PID}`, `${function:PID}` and an unknown drive-qualified
  braced path.

The same regression executes 19 isolated `pwsh -NoProfile -NonInteractive`
runtime probes inside a function. Direct/braced unqualified, `variable:`,
global, script, local and private PID targets reproduce
`Cannot overwrite variable PID because it is read-only or constant.` Every
multi-colon and env/function provider control exits zero with
`EXECUTED_WITHOUT_PID_COLLISION`.

## Pre-commit checks

Both implementation scripts parsed with zero errors. The exact PID predicate,
runtime and data-flow regression plus the complete available-host no-Docker
matrix passed:

```text
Docker execution iteration-006 PID source regression PASS untyped=REJECT typed=REJECT parenthesized=REJECT multi-target=REJECT nested-wrapper=REJECT runtime-scopes=REJECT multi-colon-provider-controls=ALLOW member-index=ALLOW rhs-string-comment-ledger=ALLOW data-flow=PASS
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

`git diff --check` passed. Final inventory found zero surviving processes whose
command line referenced `native-fixture.ps1`, zero remaining
`pscan-c2-iteration-005-*` directories and zero remaining
`pscan-docker-boundary-*` directories. All three local tag identities remained
unchanged.

Exact starting source trust passed at review-evidence commit
`cf6b61d73910d7525a2213877654dfb900fe2cb2`, tree
`862019898de01ac81edb3814495fd754b58647ee`, for 370 tracked files: 352 raw
equal, 18 authorized canonical CRLF projections and zero unexpected mismatch.
Exact-commit source trust and validation remain required after the new
candidate is committed.

No Docker command, network connection, remote read or mutation, dependency
acquisition, build, tag action, workflow dispatch or rerun, credential action,
signing, attestation, draft, release, deployment, publication or successor
work occurred.

This is author refinement only, not independent acceptance. Recovery R5
remains terminally failed and all actual-Linux, genuine-Docker, dependency,
build, reproducibility, artifact-integrity, remote and successor gates remain
closed.
