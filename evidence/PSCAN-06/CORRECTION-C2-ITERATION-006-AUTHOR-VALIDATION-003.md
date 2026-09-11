# PSCAN-06 Correction C2 iteration 006 author validation 003

Date: `2026-09-11`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`

## Exact second refined candidate

- Commit: `763730403ae538850e9d806ce9db82c815517b2d`
- Tree: `05c0baba5e61edb808b918a2dd6c853f49b8955e`
- Direct second-review-evidence parent:
  `c4cb2bb15e8b5bd7be708c0500706a5e2ae54e6a`
- Preserved first rejected candidate:
  `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`
- Preserved first refined/rejected candidate:
  `e6a767aab71db1d3f62063dded379b4701d2cb52`
- Second review-record SHA-256:
  `44FC5CC1F0F79EE5FB8B5F6D9EAA8423B7978E973EC64D9A6C89C8C985831C4D`
- Production blob: `2fb43bf89df352e9153ba5d7ad23fd59cad77749`
- Production blob SHA-256:
  `49211CA52336BE859CD0C976719B5FB1F209F9E1235F1DA4926C6C21B52905C2`
- Second refined regression blob:
  `a57f93545df3902598e5f0119298d1bf77654935`
- Second refined regression blob SHA-256:
  `3AF57C4519353EA0DD13A844A664B9F3073D2FC88CBE95F159A728E52016C47D`
- Second refinement-record blob:
  `22bfc2a15532126a314c0246211e83c0c6bb3f36`
- Second refinement-record SHA-256:
  `E01E9B31B6CA41149A4662A6B6CD86B4A5FB06A87C74E16750114006922F0C5C`

The candidate is the direct child of the evidence-bearing second rejection
commit and changes exactly two authority-permitted paths: the regression test
and the append-only second refinement record. The production script is
byte-identical to both rejected candidates. Both prior candidates, their
author evidence and both independent rejection records remain immutable.

## Assignment-target walker and adversarial proof

The predicate used against the exact production AST walks only
`AssignmentStatementAst.Left`. It unwraps typed/attributed targets, exact
parenthesized pipeline/command-expression wrappers and recursively nested
`ArrayLiteralAst` target elements. It collects only non-splatted true variable
lvalues, using `VariablePath.IsVariable`, and compares their unqualified names
to `PID` case-insensitively. Member and index lvalues terminate the walk;
their receiver, member, index and other descendants are never searched.
Unexpected shapes, multi-command pipelines and redirected command expressions
fail closed.

The inseparable executable parser/predicate self-tests reported:

```text
Docker execution iteration-006 PID source regression PASS untyped=REJECT typed=REJECT parenthesized=REJECT multi-target=REJECT nested-wrapper=REJECT member-index-drive-controls=ALLOW rhs-string-comment-ledger=ALLOW data-flow=PASS
```

This directly covers earlier untyped/typed mixed-case forms and the three
second-review counterexamples: `($PiD) = 1`,
`$PiD, $other = 1, 2` and `($PiD, $other) = 1, 2`. Further positives cover
nested typed, parenthesized, scoped, background and multiple-target wrappers.
Controls cover direct, parenthesized and typed member/index targets;
drive-qualified `env:`/`function:` targets; right-hand-side PID references;
strings; comments; ledger property `PID`; and typed non-PID assignment.

Independent AST enumeration found all valid tested left-side routes to be true
variable lvalues, typed/attributed wrappers, exact parenthesized pipeline and
command-expression wrappers, recursively nested multiple-target array
literals, or terminal member/index lvalues. Twenty representative collision
forms reproduced the read-only automatic-variable error in separate
non-profile processes. Non-lvalue expression families and multi-command or
redirected parenthesized forms produced parser errors; drive-qualified,
member and index controls did not collide. Exact enumeration details are in
`CORRECTION-C2-ITERATION-006-REFINEMENT-002.md`.

The production data-flow proof continues to identify exactly one renamed
integer binding from the parsed stat PID, exact identity
`"${linuxProcessIdentifier}:$startTime"`, exact two-field ledger shape, and
ledger bindings `PID=$linuxProcessIdentifier` and `StartTime=$startTime`.

## Exact source, parse and confinement checks

The exact-candidate source verifier reported:

```text
Exact source trust PASS commit=763730403ae538850e9d806ce9db82c815517b2d tree=05c0baba5e61edb808b918a2dd6c853f49b8955e files=368 raw_equal=350 canonical_crlf=18 unexpected_mismatch=0
index=EXACT index_flags=EXACT unsupported_config=0 nonignored_untracked=0 ignored_only_graphify=66
```

Both implementation scripts parsed with zero errors. Exact range
`c4cb2bb15e8b5bd7be708c0500706a5e2ae54e6a..763730403ae538850e9d806ce9db82c815517b2d`
passed `git diff --check` and exact two-path confinement. Committed-byte hashing
reproduced the production, regression and refinement-record identities above.

## Local no-Docker native-boundary checks

The complete exact-candidate available-host matrix passed:

```text
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

This covers the clean Windows job-object boundary, exact native source,
131072-byte stream limits, UTF-8 behavior, fixed 15000 ms timeout, child and
grandchild cleanup, empty containment membership, executable replacement,
second invocation and compatible/stale ambient types. Final inventory found
zero surviving processes whose command line referenced `native-fixture.ps1`,
zero remaining `pscan-c2-iteration-005-*` directories and zero remaining
`pscan-docker-boundary-*` directories.

Tracked, staged and non-ignored untracked state was clean. All 66 ignored
paths remained confined to preserved `graphify-out/**`. Local tags remained:

- `v1.0.0` -> `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- `release-tooling-v1.0.0-c1` ->
  `3fb7592889820fa2739a4a53588e073689621809`
- `release-tooling-v1.0.0-c2` ->
  `faef8435322c9096df09b56969662411f17356ea`

## Limits and open proof

No Docker command, image or container action; network connection; remote read
or mutation; dependency acquisition; build; tag action; workflow dispatch or
rerun; credential action; signing; attestation; draft; release; deployment;
publication; or successor action occurred.

This is author validation only. It does not create acceptance or independently
approve the candidate. A separate genuinely fresh skeptical task must inspect
exact candidate `763730403ae538850e9d806ce9db82c815517b2d` and rerun
proportionate local no-Docker checks before any acceptance record.

Recovery R5 remains terminally failed and consumed. Actual-Linux,
genuine-Docker, dependency, build, reproducibility and artifact-integrity proof
remains open. PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains
proposed and unselected; PSCAN-08 remains inactive and ineligible. No successor
is selected, activated, claimed or worked.
