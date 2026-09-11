# PSCAN-06 Correction C2 iteration 006 implementation

Date: `2026-09-11`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`
Authority commit: `de162e8c347c725a0f041d3c8b0f51df1211c12d`

## Bounded change

`Get-LinuxSessionMembers` in `build/release/docker-execution.ps1` now binds the
parsed `/proc/<id>/stat` process identifier to `$linuxProcessIdentifier`
instead of lowercase `$pid`. Only the binding and its two dependent uses were
changed. The parsed type remains `[int]`, the start time remains `[uint64]`, the
identity remains `<process-id>:<start-time>`, and the ledger retains fields
`PID` and `StartTime` with the same values.

`build/release/test-docker-execution.ps1` now parses the exact production AST
and fails unless `Get-LinuxSessionMembers`:

- has no variable assignment target named `PID` under case-insensitive
  comparison;
- binds `$linuxProcessIdentifier` exactly once from the existing `[int]`
  `/proc` stat parse;
- uses that binding in the exact process identity; and
- uses that binding as the ledger `PID` value while preserving the two-field
  ledger shape and `StartTime` binding.

## Preserved surface

No `/proc` enumeration, malformed-record handling, session membership,
liveness, signaling, cleanup, time bound, stream bound, executable identity,
private environment, Windows job object, Linux stopped PID-namespace init,
Docker operation, image, schema, release identity, scanner, verifier, policy or
consuming-project behavior changed.

## Pre-commit author checks

Both allowed implementation scripts parsed successfully. The exact production
source regression and the complete existing isolated local no-Docker matrix
reported:

```text
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

`git diff --check` passed. A post-matrix Windows process inventory found zero
surviving processes launched from `native-fixture.ps1`, and no
`pscan-c2-iteration-005-*` fixture directory remained.

The clean claim boundary established exact source trust for all 360 tracked
files: 342 raw-equal files, 18 authorized canonical LF-to-CRLF projections and
zero unexpected mismatches. Exact-commit source trust, path confinement and a
fresh post-commit matrix remain required before author-validation evidence.

These checks invoked no Docker command, network, remote read or mutation,
dependency acquisition, build, workflow action, tag change, credential,
signing, attestation, draft, release, publication or successor action.

This is bounded author implementation only. It is not local acceptance or
independent review. Recovery R5 remains terminally failed; actual-Linux,
genuine-Docker, dependency, build, reproducibility and artifact-integrity proof
remains open. PSCAN-06 remains open and unaccepted overall.
