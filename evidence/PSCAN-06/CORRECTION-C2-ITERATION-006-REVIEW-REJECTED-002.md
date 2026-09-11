# PSCAN-06 Correction C2 iteration 006 independent review rejection 002

Date: `2026-09-11`
Independent review session: `01a09003-9a4c-7862-967e-972112c9855d`
Approval session: `01a08fae-2ee1-7553-a126-e635b70a34a6`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`

## Outcome

- Result:
  `LOCAL_ACCEPTANCE_REJECTED_PARENTHESES_AND_MULTI_TARGET_AST_BYPASS`
- Refined candidate commit:
  `e6a767aab71db1d3f62063dded379b4701d2cb52`
- Refined candidate tree:
  `288664447272f9543de1b69b8ca28c27c6e1e9ff`
- Direct rejection-evidence parent:
  `86d953c75a1b3a68c7be4a6cc6fe35e4d797449b`
- Evidence-bearing review parent HEAD:
  `4b98d9bae4ae07ad2416e590986482e75fd66569`
- Evidence-bearing review parent tree:
  `ca3b9eb4b546049d61bcb8e5a57208e71c5c9216`
- Production blob:
  `2fb43bf89df352e9153ba5d7ad23fd59cad77749`
- Production blob SHA-256:
  `49211CA52336BE859CD0C976719B5FB1F209F9E1235F1DA4926C6C21B52905C2`
- Refined regression blob:
  `a3d85c5420199334e658efc9dc90152560b5f6bb`
- Refined regression blob SHA-256:
  `1259BA5BFB7C28EC6A05C0454D544AEDF176DC7994B3D8957563A9F09254AA6E`

The refinement closes the previously reported typed-target bypass, but its
exact predicate still misses valid parenthesized and multiple assignment
targets that write to case-insensitive automatic variable `PID`. This
mandatory adversarial coverage gap prevents local acceptance. The reviewer
did not repair or amend the candidate.

## Admission and preserved history

The review ran in the direct saved checkout at
`C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`, not a
worktree. It verified branch `main`, exact clean HEAD and tree above, its direct
candidate parent, candidate tree and direct rejection parent. Git directory
and common directory both resolved to `.git`. The review task remained
different from the authority approval and implementation tasks.

Tracked, staged and non-ignored untracked state was clean. All 66 ignored
paths were confined to preserved `graphify-out/**`; none was deleted, moved or
staged.

The prior candidate and rejection remain immutable:

- rejected candidate:
  `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`;
- rejected-candidate tree:
  `de70410df5ec359b85d03262e46c8d2ed5f2d822`;
- first rejection commit:
  `86d953c75a1b3a68c7be4a6cc6fe35e4d797449b`;
- first rejection record SHA-256:
  `5986DC219707494A8A9115550F99ABD7F36BD25DE4B08E4D014A7B4DF795E095`.

Required new records independently matched:

- refinement record SHA-256:
  `BB8D156D43224EF516CAD9ED84A977908E4FFE693A55B4A25416093E27B665A3`;
- second author-validation SHA-256:
  `8764AD600E0E1085BC9EE5FD2434215CF39E0D80624B29E57017D93CB6CD1E05`.

## Exact committed scope and data flow

The refined candidate changes exactly two authority-permitted paths:
`build/release/test-docker-execution.ps1` and the append-only refinement
record. The evidence-bearing review parent then changes only the three allowed
living task records and second author-validation record.

`build/release/docker-execution.ps1` has the same blob and SHA-256 as both the
prior rejected candidate and the refined candidate. The production rename is
therefore unchanged and remains bounded to the local parsed process identifier
and its two dependent uses.

Independent AST inspection confirmed that the regression parses the exact
`docker-execution.ps1` path, selects exact function
`Get-LinuxSessionMembers`, and reports no current reserved assignment. The
single parsed identifier binding remains:

```text
[int]$stat.Substring(0, $stat.IndexOf(' '))
```

The identity remains
`"${linuxProcessIdentifier}:$startTime"`; the exact ledger assignment retains
only `PID` and `StartTime`; and their values remain
`$linuxProcessIdentifier` and `$startTime`. The data-flow requirement passes.

## Blocking adversarial finding

The refined predicate unwraps only `AttributedExpressionAst` assignment-left
wrappers. After unwrapping, it classifies only a direct
`VariableExpressionAst`. Valid PowerShell assignment targets are not limited
to those two shapes.

The reviewer extracted the exact two predicate scriptblocks from the committed
refined regression AST and executed that exact predicate against independently
parsed cases. These valid colliding forms all executed in separate non-profile
processes and produced
`Cannot overwrite variable PID because it is read-only or constant.`:

| Case | Assignment-left AST | Exact predicate count |
|---|---|---:|
| `($PiD) = 1` | `ParenExpressionAst` | 0 |
| `$PiD, $other = 1, 2` | `ArrayLiteralAst` | 0 |
| `($PiD, $other) = 1, 2` | `ParenExpressionAst` | 0 |

The first case is a direct parenthesized spelling of the same automatic
variable assignment. The latter two show the same omission for a multiple
assignment target, both bare and parenthesized. Each parses without error,
causes the runtime collision, and evades the exact candidate predicate.

Positive controls succeeded: direct upper- and mixed-case, script-scoped,
typed `ConvertExpressionAst`, nested attributed/typed, and compound assignment
targets were classified and reproduced the collision. Negative controls for a
right-hand-side `$PID`, ledger property name `PID`, string/comment text,
member target, index target, parenthesized non-PID target and typed non-PID
target produced no false positives.

The mandatory review instruction explicitly requires parenthesized valid
assignment-target forms that could collide with automatic `$PID` to reject.
The result above is therefore a blocking coverage failure even though the
current production source and its identity/ledger flow are correct.

## Passing independent checks

Both implementation scripts parsed successfully. The complete current
checkout every-byte verifier reported:

```text
Exact source trust PASS commit=4b98d9bae4ae07ad2416e590986482e75fd66569 tree=ca3b9eb4b546049d61bcb8e5a57208e71c5c9216 files=366 raw_equal=348 canonical_crlf=18 unexpected_mismatch=0
index=EXACT index_flags=EXACT unsupported_config=0 nonignored_untracked=0 ignored_only_graphify=66
```

Candidate two-path confinement and `git diff --check` passed. Independent
committed-byte hashing reproduced the production and refined-regression
SHA-256 values above.

The complete available-host local no-Docker matrix passed:

```text
Docker execution iteration-006 PID source regression PASS untyped-mixed-case=REJECT typed-mixed-case=REJECT rhs-and-ledger-property=ALLOW data-flow=PASS
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

These passing checks do not cure the parenthesized and multiple-target bypass.

## Preserved terminal boundary

No Docker command, network connection, remote read or mutation, dependency
acquisition, build, tag action, workflow dispatch or rerun, credential action,
signing, attestation, draft, release, deployment, publication or successor
work occurred.

Recovery R5 remains terminally failed; its sole tag push and workflow dispatch
remain consumed. Actual-Linux, genuine-Docker, dependency, build,
reproducibility and artifact-integrity proof remains open. Refined candidate
`e6a767aab71db1d3f62063dded379b4701d2cb52` is not independently accepted.
PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains proposed and
unselected; PSCAN-08 remains inactive and ineligible. No successor is selected
or activated.
