# PSCAN-06 Correction C2 iteration 006 regression refinement 002

Date: `2026-09-11`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`
Starting review-evidence commit:
`c4cb2bb15e8b5bd7be708c0500706a5e2ae54e6a`
Starting review-evidence tree:
`3d9a63a4ed4918452206a9c77df626f4db59bc3e`

## Preserved rejection history and production result

Independent review rejected refined candidate
`e6a767aab71db1d3f62063dded379b4701d2cb52`, tree
`288664447272f9543de1b69b8ca28c27c6e1e9ff`, with result
`LOCAL_ACCEPTANCE_REJECTED_PARENTHESES_AND_MULTI_TARGET_AST_BYPASS`.
The immutable second review record is
`CORRECTION-C2-ITERATION-006-REVIEW-REJECTED-002.md`, SHA-256
`44FC5CC1F0F79EE5FB8B5F6D9EAA8423B7978E973EC64D9A6C89C8C985831C4D`.

The prior candidate `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`, first
refined candidate, their author evidence and both independent rejection
records remain unchanged history. This refinement does not modify
`build/release/docker-execution.ps1`. Its committed blob remains
`2fb43bf89df352e9153ba5d7ad23fd59cad77749`, SHA-256
`49211CA52336BE859CD0C976719B5FB1F209F9E1235F1DA4926C6C21B52905C2`.

## Exact assignment-target correction

Only `build/release/test-docker-execution.ps1` changes in test source. Its PID
predicate now walks only the exact `AssignmentStatementAst.Left` target:

- `AttributedExpressionAst` unwraps through `Child`;
- `ParenExpressionAst` unwraps through its single-element `PipelineAst` and
  exact `CommandExpressionAst.Expression`;
- `ArrayLiteralAst` recursively walks its target elements, including nested
  multiple-assignment targets;
- `VariableExpressionAst` contributes a target name only when it is not
  splatted and `VariablePath.IsVariable` is true; and
- `MemberExpressionAst` and `IndexExpressionAst` are terminal lvalues whose
  receiver, member or index descendants are not themselves assigned.

Every other target shape fails closed. A multi-element pipeline or redirected
command expression also fails closed instead of being searched. The predicate
rejects when any collected true variable-lvalue name equals `PID`
case-insensitively. It does not use `FindAll` over assignment-left descendants.
Drive-qualified `env:` and `function:` targets have
`VariablePath.IsVariable=false` and are not mistaken for automatic variable
`PID`.

The same production predicate is bound to executable parser self-tests for
the earlier untyped and typed mixed-case positives, the three exact review
counterexamples `($PiD) = 1`, `$PiD, $other = 1, 2` and
`($PiD, $other) = 1, 2`, nested typed/parenthesized/scoped targets, nested
multiple-target elements and a valid background parenthesized target.
Negative probes cover direct, typed and parenthesized member/index targets;
drive-qualified non-automatic targets; right-hand-side `$PID`; string and
comment text; ledger property `PID`; and typed non-PID assignment.

## Independent local AST enumeration

Direct parser inspection established the valid left-side family used by the
walker:

1. true variable lvalues are `VariableExpressionAst`, including braced,
   scope-qualified and compound assignments;
2. type constraints are `ConvertExpressionAst`, an
   `AttributedExpressionAst` wrapper;
3. parentheses are `ParenExpressionAst -> PipelineAst ->
   CommandExpressionAst`, and can nest around a variable, member, index or
   multiple target;
4. multiple assignment is `ArrayLiteralAst`, whose elements can recursively
   contain the preceding wrappers or another array literal; and
5. property and index lvalues terminate as `MemberExpressionAst` and
   `IndexExpressionAst`.

Separate non-profile processes confirmed 20 representative direct, braced,
scope-qualified, typed, attributed, parenthesized, nested, compound,
background and multiple-target PID forms all fail with
`Cannot overwrite variable PID because it is read-only or constant.` Direct
and parenthesized member/index targets and `env:`/`function:` drive targets
executed without that collision.

`ArrayExpressionAst`, `SubExpressionAst`, `UsingExpressionAst`, splatted and
unary expressions, hashtables, script blocks, strings, constants, multi-command
pipelines and redirected command expressions all produced parser errors as
assignment left sides. They therefore cannot bind automatic PID in a parsed
production script and are not traversed by the walker.

## Pre-commit checks

Both implementation scripts parsed with zero errors. The exact PID source and
data-flow regression plus the complete available-host no-Docker matrix passed:

```text
Docker execution iteration-006 PID source regression PASS untyped=REJECT typed=REJECT parenthesized=REJECT multi-target=REJECT nested-wrapper=REJECT member-index-drive-controls=ALLOW rhs-string-comment-ledger=ALLOW data-flow=PASS
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

`git diff --check` passed. Final inventory found zero surviving processes whose
command line referenced `native-fixture.ps1`, zero remaining
`pscan-c2-iteration-005-*` fixture directories and zero remaining
`pscan-docker-boundary-*` directories. The three local tag identities remained
unchanged.

Exact starting source trust passed at review-evidence commit
`c4cb2bb15e8b5bd7be708c0500706a5e2ae54e6a`, tree
`3d9a63a4ed4918452206a9c77df626f4db59bc3e`, for 367 tracked files: 349 raw
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
