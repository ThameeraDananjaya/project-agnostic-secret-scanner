# PSCAN-06 Correction C2 iteration 006 author validation 002

Date: `2026-09-11`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`

## Exact refined candidate

- Commit: `e6a767aab71db1d3f62063dded379b4701d2cb52`
- Tree: `288664447272f9543de1b69b8ca28c27c6e1e9ff`
- Direct review-evidence parent:
  `86d953c75a1b3a68c7be4a6cc6fe35e4d797449b`
- Preserved rejected candidate:
  `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`
- Preserved rejected-candidate tree:
  `de70410df5ec359b85d03262e46c8d2ed5f2d822`
- Production blob: `2fb43bf89df352e9153ba5d7ad23fd59cad77749`
- Production blob SHA-256:
  `49211CA52336BE859CD0C976719B5FB1F209F9E1235F1DA4926C6C21B52905C2`
- Refined regression blob:
  `a3d85c5420199334e658efc9dc90152560b5f6bb`
- Refined regression blob SHA-256:
  `1259BA5BFB7C28EC6A05C0454D544AEDF176DC7994B3D8957563A9F09254AA6E`

The prior candidate and independent rejection remain immutable history. The
refined candidate is the direct child of the evidence-bearing rejection commit
and changes exactly two authority-permitted paths: the regression script and
the append-only refinement record. It does not change the production script.

## Review-finding correction

The same predicate used on the exact production AST now unwraps only
assignment-target `AttributedExpressionAst` wrappers, including typed
`ConvertExpressionAst`, before comparing the unqualified target variable name
to `PID` case-insensitively. It does not inspect right-hand-side references or
arbitrary descendants of non-variable targets.

Parser-only self-tests are inseparable from that predicate and reported:

```text
Docker execution iteration-006 PID source regression PASS untyped-mixed-case=REJECT typed-mixed-case=REJECT rhs-and-ledger-property=ALLOW data-flow=PASS
```

The cases prove `$pId = 1` and `[int]$PiD = 1` are detected, while right-hand-
side `$PID`, ledger data property `PID`, and a typed non-PID target do not
produce false positives. The existing exact source checks also continue to
prove the renamed `[int]` Linux process identifier feeds both the
`<process-id>:<start-time>` identity and ledger `PID` value without changing
the `StartTime` binding or ledger shape.

## Exact source, parse and path checks

The exact-commit source verifier reported:

```text
Exact source trust PASS commit=e6a767aab71db1d3f62063dded379b4701d2cb52 tree=288664447272f9543de1b69b8ca28c27c6e1e9ff files=365 raw_equal=347 canonical_crlf=18 unexpected_mismatch=0
index=EXACT index_flags=EXACT unsupported_config=0 nonignored_untracked=0 ignored_only_graphify=66
```

Both allowed implementation scripts parsed successfully. Exact range
`86d953c75a1b3a68c7be4a6cc6fe35e4d797449b..e6a767aab71db1d3f62063dded379b4701d2cb52`
passed `git diff --check` and exact two-path confinement. The production blob
was reverified byte-identical before both candidate validation runs.

## Local no-Docker native-boundary checks

The exact refined candidate reran the full available-host matrix:

```text
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

This covers the clean Windows job-object boundary, exact compiled source,
stream and UTF-8 limits, fixed timeout, child and grandchild cleanup, empty
containment membership, executable replacement rejection, second-invocation
rejection, and compatible/stale ambient-type rejection with zero fake calls,
zero Docker calls and zero trusted fabricated results.

Final inventory found zero surviving processes whose command line referenced
`native-fixture.ps1`, zero remaining `pscan-c2-iteration-005-*` fixture
directories and zero remaining `pscan-docker-boundary-*` directories. Tracked,
staged and non-ignored untracked state was clean; all 66 ignored paths remained
confined to preserved `graphify-out/**`.

Local tag identities remained unchanged:

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

This is author validation only. It does not overturn the prior rejection,
create local acceptance or claim independent review of the refined candidate.
A separate genuinely fresh skeptical task must inspect exact candidate
`e6a767aab71db1d3f62063dded379b4701d2cb52` and rerun proportionate local
no-Docker checks before any acceptance record.

Recovery R5 remains terminally failed and consumed. Actual-Linux,
genuine-Docker, dependency, build, reproducibility and artifact-integrity proof
remains open. PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains
proposed and unselected; PSCAN-08 remains inactive and ineligible. No successor
is selected, activated, claimed or worked.
