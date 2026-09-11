# PSCAN-06 Correction C2 iteration 006 local independent acceptance

Date: `2026-09-11`
Independent review session: `01a09003-9a4c-7862-967e-972112c9855d`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`

## Result and boundary

`LOCAL_ACCEPTANCE_PASS_REMOTE_PROOF_OPEN`.

Correction C2 iteration 006 is independently accepted locally for its bounded
PowerShell automatic-PID collision correction and inseparable source
regression. The fourth skeptical review found no blocking defect after
independently exercising the exact committed predicate and installed parser
and runtime. This is not acceptance or completion of PSCAN-06 as a whole and
does not authorize Docker, network, remote, workflow, release or successor
action.

- Exact reviewed candidate:
  `cf1679f9bca24887340fa4060f37d1ebff21f305`.
- Exact reviewed tree:
  `6cb0198c71881baff9c601c31ad6aa84217f4a49`.
- Direct third-rejection parent:
  `cf6b61d73910d7525a2213877654dfb900fe2cb2`.
- Evidence-bearing review parent HEAD:
  `6ae41db6c52c7bb3efd627131dee79f5fb3bc851`.
- Evidence-bearing review parent tree:
  `26aac56fb182ce01e617c7c62487a51125a5673c`.
- Production blob:
  `2fb43bf89df352e9153ba5d7ad23fd59cad77749`.
- Production blob SHA-256:
  `49211CA52336BE859CD0C976719B5FB1F209F9E1235F1DA4926C6C21B52905C2`.
- Third refined regression blob:
  `39c08202e32d7fc45ac57d90df0a11ef65ab048a`.
- Third refined regression blob SHA-256:
  `E74D9EB349C4394D4AD8E8D6EB307FBD091D469791C84A690187FD1E869D8DAB`.
- Fourth author-validation SHA-256:
  `2FDD8FD057FCD9BC75CA0C84BCFFBBFCFFCBD78B7E3D517DFE77BFC2040BCF52`.
- Blocking findings: none.

## Admission and preserved history

The review ran in the direct saved checkout at
`C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`, not a
worktree. It verified branch `main`, exact clean review-parent HEAD and tree,
the direct candidate parent, exact candidate tree and direct third-rejection
parent. Git directory and common directory both resolved to `.git`. The
independent review session differs from the implementation session.

Tracked, staged and non-ignored untracked state was clean. All 66 ignored
paths were confined to preserved `graphify-out/**`; none was deleted, moved or
staged. The candidate changes exactly two authority-permitted paths:
`build/release/test-docker-execution.ps1` and the append-only third refinement
record. The review parent adds only the three living task records and fourth
author-validation record. Both exact ranges pass `git diff --check`.

All three rejected candidates and rejection records remain immutable:

- candidate `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`, rejected for a
  typed-target bypass; rejection SHA-256
  `5986DC219707494A8A9115550F99ABD7F36BD25DE4B08E4D014A7B4DF795E095`;
- candidate `e6a767aab71db1d3f62063dded379b4701d2cb52`, rejected for
  parenthesized and multiple-target bypasses; rejection SHA-256
  `44FC5CC1F0F79EE5FB8B5F6D9EAA8423B7978E973EC64D9A6C89C8C985831C4D`;
  and
- candidate `763730403ae538850e9d806ce9db82c815517b2d`, rejected for
  multi-colon variable-name false positives; rejection SHA-256
  `F21FD594E757747EE5D3F9865FBE2C33D1FECE7AF09B28AA50778D355A183373`.

The third refinement record independently matched SHA-256
`81B727EFCCC4BC270C2F40932205D7670412731D3075D1FBA56025DF22E9EFDB`.

## Exact predicate and runtime validation

The reviewer extracted the exact committed `isAutomaticPidVariablePath`,
`getAssignmentTargetVariableNames` and `isReservedPidAssignment` scriptblocks
from the regression AST rather than recreating the predicate. Reflection
against the installed `System.Management.Automation.VariablePath` type
independently confirmed that these getters are public:

- `UserPath`, `IsVariable`, `IsDriveQualified` and `DriveName`;
- `IsUnqualified`, `IsUnscopedVariable`, `IsGlobal`, `IsScript`, `IsLocal`
  and `IsPrivate`.

`UnqualifiedPath` exists but its getter is assembly-internal and is not used.
Independent parser probes confirmed that escaped braced spellings normalize
in `UserPath`, scope and provider flags classify the complete path, and legal
multi-colon names retain every colon.

Seventy-seven independently generated sources mutated scope/path case,
bracing, escaped characters, extra-colon placement, slash and backslash
placement, Unicode lookalikes, wrapper depth, compound operators, tuple
position and member/index context. Each source was classified by the exact
predicate and exercised in a separate `pwsh -NoProfile -NonInteractive`
process when runtime-comparable. Seventy-six cases agreed directly. One
over-parenthesized nested destructuring mutation was parser-accepted but the
runtime binder failed on its nested `IList` access before any lvalue write;
it was therefore treated as runtime-unsupported and correctly fail-closed,
not as an allowed non-PID control. The candidate's exact nested typed target
and simpler nested typed tuple variants reproduced the automatic-variable
collision. Every runtime-supported flagged target collided and every allowed
target did not.

Positive collision coverage included:

- unqualified and braced `$PID`, mixed case and escaped spellings;
- unbraced and braced `variable:`, `global:`, `script:`, `local:` and
  `private:` forms with mixed qualifier/path case;
- typed, multiply typed, attributed, `[ref]`, parenthesized, deeply nested,
  background and compound-assignment targets; and
- bare, parenthesized and recursively nested multiple-assignment targets with
  PID in different tuple positions.

Negative coverage included:

- `${variable:env:PID}`, `${global:env:PID}`, `${local:foo:PID}`,
  `${script:foo:PID}` and `${private:foo:PID}`;
- qualifier nesting, doubled colons, trailing components and slash/backslash
  path variants;
- `env:`, `function:`, unknown and workflow provider paths;
- PID prefix/suffix names, whitespace-bearing names and Unicode dotted/dotless
  I controls;
- member, dynamic-member, static-member, invoked-member and index lvalues,
  including typed and parenthesized forms and `$PID` used only as receiver or
  index; and
- right-hand-side references, hashtable/property keys, literal/expanded
  strings, comments and tuple controls.

Fourteen known supported non-PID lvalue probes covered every valid installed-
parser left-side family: `VariableExpressionAst`, typed/attributed expression,
exact parenthesized pipeline/command expression, recursively nested
`ArrayLiteralAst`, `MemberExpressionAst` and `IndexExpressionAst`. All returned
non-reserved without a walker exception. Twelve unsupported left-side probes
were parser-invalid and the walker failed closed for their array-expression,
subexpression, hashtable, script-block, string, constant, multi-command,
redirection, unary and ternary shapes. No known supported non-PID lvalue was
rejected.

## Production source and data flow

`build/release/docker-execution.ps1` remains byte-identical to every rejected
and refined candidate. Independent AST inspection found exactly one
`Get-LinuxSessionMembers` function, zero parser errors and no reserved PID
assignment. It reproduced exactly:

- one integer `linuxProcessIdentifier` binding from
  `[int]$stat.Substring(0, $stat.IndexOf(' '))`;
- one identity binding `"${linuxProcessIdentifier}:$startTime"`;
- one `$Ledger[$identity]` assignment with exactly two fields; and
- ledger values `PID=$linuxProcessIdentifier` and `StartTime=$startTime`.

The production rename therefore preserves the intended Linux process identity
and ledger flow.

## Source, parser and no-Docker validation

Both implementation scripts parsed with zero errors. The complete current-
checkout every-byte verifier reported:

```text
Exact source trust PASS commit=6ae41db6c52c7bb3efd627131dee79f5fb3bc851 tree=26aac56fb182ce01e617c7c62487a51125a5673c files=372 raw_equal=354 canonical_crlf=18 unexpected_mismatch=0
index=EXACT index_flags=EXACT unsupported_config=0 nonignored_untracked=0 ignored_only_graphify=66
```

The complete available-host local no-Docker matrix passed:

```text
Docker execution iteration-006 PID source regression PASS untyped=REJECT typed=REJECT parenthesized=REJECT multi-target=REJECT nested-wrapper=REJECT runtime-scopes=REJECT multi-colon-provider-controls=ALLOW member-index=ALLOW rhs-string-comment-ledger=ALLOW data-flow=PASS
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

This covers the unchanged clean Windows job-object boundary, exact native
source, 131072-byte stream limits, UTF-8 behavior, fixed 15000 ms timeout,
child and grandchild cleanup, empty containment membership, executable
replacement, second invocation and compatible/stale ambient-type rejection.
Final inventory found zero surviving processes whose command line referenced
`native-fixture.ps1`, zero remaining `pscan-c2-iteration-005-*` directories
and zero remaining `pscan-docker-boundary-*` directories. The three local tag
identities remained unchanged.

## Open proof and owner gates

No Docker command, image or container action; network connection; remote read
or mutation; dependency acquisition; build; tag action; workflow dispatch or
rerun; credential action; signing; attestation; draft; release; deployment;
publication; or successor action occurred.

Recovery R5 remains terminally failed and consumed. Actual-Linux,
genuine-Docker, dependency, build, reproducibility and artifact-integrity proof
remains open. PSCAN-06 remains open and unaccepted overall; only bounded local
Correction C2 iteration 006 is accepted. PSCAN-07 remains proposed and
unselected; PSCAN-08 remains inactive and ineligible. No successor is selected,
activated, claimed or worked.
