# PSCAN-06 Correction C2 iteration 006 regression refinement 001

Date: `2026-09-11`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`
Starting review-evidence commit:
`86d953c75a1b3a68c7be4a6cc6fe35e4d797449b`

## Preserved rejection and production result

Independent review rejected prior candidate
`8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`, tree
`de70410df5ec359b85d03262e46c8d2ed5f2d822`, with result
`LOCAL_ACCEPTANCE_REJECTED_REGRESSION_AST_BYPASS`. The immutable review record
is `CORRECTION-C2-ITERATION-006-REVIEW-REJECTED.md`, SHA-256
`5986DC219707494A8A9115550F99ABD7F36BD25DE4B08E4D014A7B4DF795E095`.

That review independently accepted the bounded production rename as correct
but proved the regression predicate missed typed assignment targets whose
`AssignmentStatementAst.Left` is a `ConvertExpressionAst`. The rejected
candidate, its author evidence and the rejection remain unchanged history.

`build/release/docker-execution.ps1` is not modified by this refinement. Its
committed blob remains `2fb43bf89df352e9153ba5d7ad23fd59cad77749`, SHA-256
`49211CA52336BE859CD0C976719B5FB1F209F9E1235F1DA4926C6C21B52905C2`.

## Bounded regression correction

Only `build/release/test-docker-execution.ps1` changes in production/test
source. The assignment-target predicate now unwraps only
`AttributedExpressionAst` target wrappers, including
`ConvertExpressionAst`, until it reaches the assignable child. It classifies
the target as reserved only when that child is a `VariableExpressionAst` whose
unqualified name equals `PID` case-insensitively.

The same predicate is bound to three parser-only adversarial self-tests:

- untyped mixed-case `$pId = 1` must be rejected;
- typed mixed-case `[int]$PiD = 1` must be rejected; and
- right-hand-side `$PID`, ledger property name `PID`, and typed non-PID target
  must not be classified as assignments to automatic variable `PID`.

The predicate does not search arbitrary descendants of an assignment left
side, so an index expression or other non-variable target containing a `PID`
reference is not misclassified. The existing exact production source and
identity/ledger data-flow checks remain mandatory.

## Pre-commit checks

Both allowed implementation scripts parsed successfully. The unchanged
production blob was reverified before running the exact regression and full
available-host no-Docker matrix, which reported:

```text
Docker execution iteration-006 PID source regression PASS untyped-mixed-case=REJECT typed-mixed-case=REJECT rhs-and-ledger-property=ALLOW data-flow=PASS
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

`git diff --check` passed. Final inventory found zero surviving processes whose
command line referenced `native-fixture.ps1`, zero remaining
`pscan-c2-iteration-005-*` fixture directories and zero remaining
`pscan-docker-boundary-*` directories.

Exact starting source trust passed at review-evidence commit
`86d953c75a1b3a68c7be4a6cc6fe35e4d797449b`, tree
`5ad267189f20a58cfb83d11935bbe03df674eb80`, for 364 tracked files: 346 raw
equal, 18 authorized canonical CRLF projections and zero unexpected mismatch.
Exact-commit source trust and validation remain required after the refined
candidate is committed.

No Docker command, network connection, remote read or mutation, dependency
acquisition, build, tag action, workflow dispatch or rerun, credential action,
signing, attestation, draft, release, deployment, publication or successor
work occurred.

This is author refinement only, not independent acceptance. Recovery R5
remains terminally failed and all remaining proof boundaries stay open.
