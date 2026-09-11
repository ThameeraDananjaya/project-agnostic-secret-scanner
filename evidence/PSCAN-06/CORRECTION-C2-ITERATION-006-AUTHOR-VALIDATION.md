# PSCAN-06 Correction C2 iteration 006 author validation

Date: `2026-09-11`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`

## Exact candidate

- Commit: `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`
- Tree: `de70410df5ec359b85d03262e46c8d2ed5f2d822`
- Authority parent: `de162e8c347c725a0f041d3c8b0f51df1211c12d`
- Production blob: `2fb43bf89df352e9153ba5d7ad23fd59cad77749`
- Production blob SHA-256:
  `49211CA52336BE859CD0C976719B5FB1F209F9E1235F1DA4926C6C21B52905C2`
- Regression blob: `b0cc6c6d40259d934e2a4a7b13f8f3db2387b08c`
- Regression blob SHA-256:
  `3876319375E591FE42CFF5DBA4DC7AEBD242CDF270DC3E0872F2D50AC737BD7E`

The candidate is the direct child of the committed iteration-006 authority.
Its four changed paths are the two permitted scripts plus the preflight and
implementation evidence for this iteration.

## Exact source, format and path checks

The authority-compliant exact-commit source verifier reported:

```text
Exact source trust PASS commit=8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4 tree=de70410df5ec359b85d03262e46c8d2ed5f2d822 files=362 raw_equal=344 canonical_crlf=18 unexpected_mismatch=0
index=EXACT index_flags=EXACT unsupported_config=0 nonignored_untracked=0 ignored_only_graphify=66
```

Both implementation scripts parsed successfully. The committed range
`de162e8c347c725a0f041d3c8b0f51df1211c12d..8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`
passed `git diff --check` and exact four-path confinement.

The source-level production AST regression passed. It proves
`Get-LinuxSessionMembers` contains no assignment target whose unqualified
variable name equals `PID` under case-insensitive comparison. It also proves
the single `[int]` stat-derived `$linuxProcessIdentifier` binding feeds the
exact `<process-id>:<start-time>` identity and the ledger `PID` value, while
the ledger `StartTime` binding and two-field shape remain exact.

## Local no-Docker native-boundary checks

The exact committed scripts reran in new non-profile processes and reported:

```text
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

This covers the available-host clean native boundary, exact compiled source,
131072-byte stream bounds, UTF-8 behavior, 15000 ms timeout behavior, child and
grandchild cleanup, empty Windows job membership, executable replacement
rejection, second-invocation rejection, and compatible/stale ambient-type
rejection with zero fake or Docker calls and zero trusted fabricated results.

The final Windows process inventory found zero surviving processes launched
from `native-fixture.ps1`. No `pscan-c2-iteration-005-*` fixture directory
remained. Tracked, staged and non-ignored untracked state was clean; all 66
ignored paths remained confined to preserved `graphify-out/**`.

Local product and tooling tags remained unchanged:

- `v1.0.0` -> `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- `release-tooling-v1.0.0-c1` ->
  `3fb7592889820fa2739a4a53588e073689621809`
- `release-tooling-v1.0.0-c2` ->
  `faef8435322c9096df09b56969662411f17356ea`

## Limits and open proof

No Docker command, image or container action; network access; remote read or
mutation; dependency acquisition; build; workflow action; tag change;
credential action; signing; attestation; draft; release; publication; or
successor action occurred.

This is author validation only. It does not create local acceptance or claim
independent review. A separate genuinely fresh skeptical task must inspect the
exact candidate and rerun proportionate local no-Docker checks before any local
acceptance record.

Recovery R5 remains terminally failed and its sole tag push and workflow
dispatch remain consumed. Actual-Linux, genuine-Docker, dependency, build,
reproducibility and artifact-integrity proof remains open. PSCAN-06 remains
open and unaccepted overall; PSCAN-07 remains proposed and unselected;
PSCAN-08 remains inactive and ineligible. No successor is selected, activated,
claimed or worked.
