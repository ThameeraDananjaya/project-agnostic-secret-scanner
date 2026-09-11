# PSCAN-06 Correction C2 iteration 007 implementation

Date: `2026-09-11`
Implementation session: `01a09151-048a-7bf3-a3a6-2d0c26415f10`
Starting authority commit: `b7a0fb211e4631f915e368174ae58022d1519182`
Starting authority tree: `bd2272041514675a8446bdde6733e359225c4e39`

## Result

`LOCAL_IMPLEMENTATION_COMPLETE_AUTHOR_VALIDATION_PENDING`.

The candidate implements only PSCAN-06 Correction C2 iteration 007's
append-only immutable R6 identity roll-forward. It creates no tag and performs
no Docker, dependency, network, remote, workflow, signing, publication,
spending or successor action.

## Bounded design

- New workflow `.github/workflows/release-recovery-v1.0.0-c2-r6.yml` derives
  from the immutable C2 workflow and changes only the selected R6
  name/concurrency, tag/ref, workflow path, certificate identity and schema
  2.2 assertion.
- New `contracts/release-manifest/schema-2.2.json` derives from immutable
  schema 2.1 and changes only version/title/identifier and the selected R6
  tag/ref/path/certificate identity.
- The active builder and embedded verifier now emit and require schema 2.2 and
  the R6 identity. The release set still includes schemas 2.0 and 2.1 and adds
  schema 2.2.
- The generic verifier preserves schema 2.0/C1 and schema 2.1/C2 behavior, adds
  schema 2.2/R6, and binds each version to its exact tag/ref/workflow and
  certificate identity.
- Integration and acceptance tests cover historical preservation, positive
  2.2 parsing/policy, old-under-new and new-under-old mixtures, single-field
  identity mutations, immutable predecessor hashes and source agreement across
  workflow, schema, builder and both verifiers.
- `docker-execution.ps1` adds only one exact copy of schema 2.2 to the release
  materialization. Its accepted `Get-LinuxSessionMembers` function, PID
  predicate/data flow, native boundary, Docker operation table and containment
  behavior are unchanged.

## Candidate paths

The implementation changes only authority-allowed paths:

```text
.github/workflows/release-recovery-v1.0.0-c2-r6.yml
contracts/release-manifest/schema-2.2.json
build/release/build.ps1
build/release/docker-execution.ps1
build/release/cmd/release-verifier/main.go
internal/verify/release.go
tests/integration/supply-chain/release_test.go
tests/acceptance/supply-chain/release_test.go
docs/decisions/DEC-004-C2-R6-IMMUTABLE-IDENTITY.md
docs/release/OFFLINE-VERIFICATION-RUNBOOK.md
docs/release/PSCAN-06-RELEASE-PLAN.md
docs/release/RELEASE-AUTHORITY.md
docs/validation/PSCAN-06-VALIDATION.md
README.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-007-PREFLIGHT.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-007-IMPLEMENTATION.md
```

`build/release/test-crlf-shell-payloads.ps1` required no edit: its existing
build-payload parser exercises the materialization payload after the bounded
schema-copy addition.

## Preserved authority

The old C2 workflow and schema 2.1 remain byte-identical to their authority
SHA-256 values. PASS contracts, historical evidence, all locked tags, product
source and scanner behavior remain unchanged. Ignored `graphify-out/**`
material remains preserved, excluded from product authority and unstaged.

## Preliminary local validation

The complete available-host no-Docker harness passed with explicit native
exit-code enforcement:

```text
Docker execution iteration-006 PID source regression PASS untyped=REJECT typed=REJECT parenthesized=REJECT multi-target=REJECT nested-wrapper=REJECT runtime-scopes=REJECT multi-colon-provider-controls=ALLOW member-index=ALLOW rhs-string-comment-ledger=ALLOW data-flow=PASS
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

PowerShell parsing passed for the changed builder and Docker materialization
route and for the unchanged CRLF payload harness. JSON parsing passed for
schemas 2.1 and 2.2. Exact normalization proved the new workflow and schema
are byte-identical to their predecessors after only the selected identity and
schema-version substitutions. The preserved predecessor and new-file SHA-256
values at this boundary were:

```text
old workflow  C5F40F1B32E87C005FE33EE607AF7E3D19EE4F0173C21619E21158C31AA0BDB4
old schema    CAA9CD26665CC3A3550AFFEA7490A0F1F277A69B10B1F1537787616E8AB973CE
new workflow  CFE3FB919CA26A7AE6CF2A8F56D22BA61B96242298E010940E7B863053C9FA1A
new schema    347E424F23F48CF25A409328E3A2E6773F589104AC0129216977DE6D01ED872A
```

Path confinement and `git diff --check` passed. All 19 changed or new paths
were authority-allowed, historical evidence had zero modifications, the index
had zero staged paths, and all 66 ignored paths remained confined to
`graphify-out/**`.

## Explicitly unproved checks

The retained, provenance-matching Go `1.27.0` tree has `go.exe` SHA-256
`7D828191BA32519A9C9361789AB647486236ED45C660889196C7770A8FF1985C`,
but it does not contain `bin/gofmt.exe`, Go formatter source or the standard
library source required by a fresh cache. An initial command attempted the
absent `gofmt.exe`; PowerShell reported command-not-found, but the surrounding
wrapper then printed a spurious `GOFMT=PASS_PROVEN_GO_1.27.0` line because it
had not asserted executable existence and native exit status. That output is
void and is not validation evidence. Subsequent attempts asserted existence
and every native exit code and stopped closed. No source file was formatted by
any of those attempts.

With toolchain selection fixed to local, proxy and checksum lookup disabled,
and fresh disposable caches, the bounded Go test attempt stopped at unavailable
module/standard-library material (`module lookup disabled by GOPROXY=off`). No
dependency was downloaded and no historical module or build cache was mutated.
Go formatting, parsing, compilation and tests are therefore unproved in this
session; the changed Go diff received only manual syntax/style inspection.

No YAML parser, `actionlint` executable or JSON Schema engine was already
available in the permitted local runtimes. YAML parsing and executable JSON
Schema validation are therefore unproved. The workflow's exact normalized
derivation from the preserved historical workflow, JSON parsing of schema 2.2
and its exact normalized derivation from schema 2.1 passed, but do not replace
those unavailable checks.

Author validation must run against the exact committed candidate from a clean
materialization. Independent acceptance remains reserved to a separate fresh
task. PSCAN-06 remains open and unaccepted overall, and Recovery R6 remains
unauthorized.
