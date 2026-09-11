# PSCAN-06 Correction C2 iteration 010 preflight and claim

Date: `2026-09-12`
Implementation session: `01a092b2-d833-74f3-99c0-a2def0764153`
Authority session: `01a092ac-14df-78f1-8604-0e2826b55c39`

## Claim and authority

This genuinely fresh saved-project session claims only PSCAN-06 Correction C2
iteration 010. PSCAN-06 is the sole open task. Recovery R5 remains terminal
and consumed; Recovery R6 remains unauthorized. No Docker, dependency
acquisition, network, remote, tag, workflow, signing, publication, spending,
subscription or successor action is claimed or permitted.

The session read the repository agreement, README, controlling
`PASS-OUTCOME-SPEC-001`, live tracker, PSCAN-06 task and reading map, both
iteration-009 authorities, the iteration-009 preflight and continuation
preflight/implementation/author-validation records, the independent rejection,
the complete iteration-010 authority and the selected test file before editing.
The existing 108-node Graphify graph was queried read-only first for routing;
it predates iteration 010, so every authority and drift-prone fact was verified
against live files and Git objects.

The implementation session ID differs from the authority session ID. Branch,
HEAD, tree, parent and selected test blob matched exactly:

```text
branch=main
HEAD=ad77e68f8ae6710c24ba818bfb04c8ecfec5f4e4
tree=f74d2dbfbe5469216876a99de4f046752afb48a6
parent=ddb67e1112ad0b263c98933cfad3b62263cf4ba4
tests/integration/supply-chain/release_test.go blob=106f2efacb44177fa08d62dd86afefa1f7a52d16
iteration-010 authority SHA-256=9EAE28BC78743873C3C8286B36A0BD24A682CCA17BC5116155B7B497A0A423C2
```

The authority commit differs from its parent through exactly the five
authority-session paths:

```text
README.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/PSCAN-06.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-010-AUTHORITY.md
```

## Exact source trust

The index matched the complete HEAD tree by path, mode and object. Index flags
were exact, unsupported shortcut/partial-clone configuration was absent, the
nonignored worktree and index were clean, and all 66 ignored entries were
confined to preserved `graphify-out/**`:

```text
CURRENT_SOURCE_TRUST=PASS commit=ad77e68f8ae6710c24ba818bfb04c8ecfec5f4e4 tree=f74d2dbfbe5469216876a99de4f046752afb48a6 files=392 raw_equal=374 canonical_crlf=18 index=EXACT flags=EXACT nonignored_untracked=0 ignored_only_graphify=66
```

A fresh local-only clone used `--local --no-hardlinks --no-checkout`, disabled
checkout conversion with `core.autocrlf=false` and `core.eol=lf`, detached at
the exact authority commit and proved every tracked file raw-equal:

```text
MATERIALIZATION_SOURCE_TRUST=PASS commit=ad77e68f8ae6710c24ba818bfb04c8ecfec5f4e4 tree=f74d2dbfbe5469216876a99de4f046752afb48a6 files=392 raw_equal=392 canonical_crlf=0
materialization=C:\Users\ITDan\AppData\Local\Temp\pscan-06-i010-materialized-01a092b2-d833-74f3-99c0-a2def0764153
```

## Pinned offline tool and cache admission

The prescribed local tool identities passed before implementation:

```text
go version go1.27.1 windows/amd64
go.exe SHA-256 D3CCDB604EAFA6031133AEFE1A3DB24F0BB7362B857BC2125AC4E4C178B4B490
gofmt.exe SHA-256 AB730B8446C0C2369BE901F263FB7433EA0D2D053D3A3584CA543AAFF3B667CC
```

The literal frozen module cache has inherited `(OI)(CI)(RX)` access, no
reparse points, 34,447 files, 7,678 directories and 1,132,106,847 file bytes.
Its sorted relative-path/type/length/attributes/last-write metadata SHA-256
before use is
`7CC37B44BDEAAB4BEFB844FE707EE9A9A018CAB05650C1322A3D6A65173DEDFD`.
All validation will use fresh session-specific external `GOCACHE`, `GOTMPDIR`
and compile outputs with `GOPROXY=off`, `GOSUMDB=off`,
`GOTOOLCHAIN=local`, `GOENV=off`, `GOTELEMETRY=off` and `GOVCS=*:off`.

## Bounded implementation surface

The sole implementation-bearing path is
`tests/integration/supply-chain/release_test.go`. Only new iteration-010
evidence and the four authority-listed living documents may accompany it.
Every production path, other test file, predecessor evidence and the committed
iteration-010 authority are immutable. A need for another path stops
fail-closed.
