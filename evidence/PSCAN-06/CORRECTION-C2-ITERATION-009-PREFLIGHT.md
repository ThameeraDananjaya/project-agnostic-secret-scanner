# PSCAN-06 Correction C2 iteration 009 preflight and claim

Date: `2026-09-12`
Implementation session: `01a09228-1204-70b1-8b7e-1be433b2c496`
Authority approval session: `01a091f1-3d3f-7a41-a541-f705cd29b687`

## Claim and state

This genuinely fresh saved-project session claims only PSCAN-06 Correction C2
iteration 009. PSCAN-06 is the sole open task. Recovery R5 remains terminal
and consumed; Recovery R6 remains unauthorized. No Docker, network, remote,
tag, workflow, signing, publication, spending, subscription or successor
action is claimed or permitted.

The session read the repository agreement, controlling
`PASS-OUTCOME-SPEC-001`, live tracker, PSCAN-06 task and reading map, and the
complete iteration 007 through iteration 009 authority, implementation,
preflight, validation and rejection chain. The existing 108-node Graphify
graph was queried read-only for routing and was not refreshed or modified.

## Exact authority and source trust

- Branch: `main`.
- Starting authority commit:
  `a80350f93afad6e62fa82e2ca822c10207b1f89a`.
- Starting tree:
  `875fa97b57db58dbac343dd88f13861da79365d4`.
- Direct parent:
  `c6ea1f08b28bb654023546f30156ff71bcfe4148`.
- The implementation session differs from the approval session.
- Tracked, staged and non-ignored untracked state was clean.
- All 66 ignored paths were confined to preserved `graphify-out/**`.

The graphify-aware current-checkout verifier proved:

```text
CURRENT_SOURCE_TRUST=PASS commit=a80350f93afad6e62fa82e2ca822c10207b1f89a tree=875fa97b57db58dbac343dd88f13861da79365d4 files=385 raw_equal=367 canonical_crlf=18 index=EXACT flags=EXACT nonignored_untracked=0 ignored_only_graphify=66
```

A separate local-only, no-hardlink materialization with checkout conversion
disabled proved all committed inputs raw-equal:

```text
MATERIALIZATION_SOURCE_TRUST=PASS commit=a80350f93afad6e62fa82e2ca822c10207b1f89a tree=875fa97b57db58dbac343dd88f13861da79365d4 files=385 raw_equal=385 canonical_crlf=0
```

The eight controlling iteration-009 hashes and all nine critical workflow,
schema, builder, materializer, verifier and acceptance-test identities from
iteration 008 matched. The selected test file matched blob
`60f8326cdde115a03e3219b04a3e1e40fbb43aa3`, SHA-256
`29E715D5534B389BCB65E046E84E5F2B2828F3D400BB664AFB13B55F12F377D0`
and 17,768 committed bytes.

## Exact bounded defect evidence

The original R6 workflow contains the five selected tokens exactly
`1,1,8,3,1` times. It contains `# v4.2.2` exactly once and contains no
`# v4.2.1`.

The pinned Go 1.27.1 formatter reproduced the exact pre-existing projection:

```text
GOFMT_BASE_EXIT=1 lines=66 hunks=4
@@ -162,7 +162,9 @@
@@ -193,10 +195,14 @@
@@ -209,10 +215,16 @@
@@ -263,9 +275,9 @@
```

No formatter output outside those four authority-listed hunks is admitted.

## Offline toolchain and frozen cache

The exact toolchain identities passed:

```text
go version go1.27.1 windows/amd64
go.exe SHA-256 D3CCDB604EAFA6031133AEFE1A3DB24F0BB7362B857BC2125AC4E4C178B4B490
gofmt.exe SHA-256 AB730B8446C0C2369BE901F263FB7433EA0D2D053D3A3584CA543AAFF3B667CC
```

The literal frozen module cache retained inherited `(OI)(CI)(RX)` access and
the accepted inventory of 34,447 files, 7,678 directories and 1,132,106,847
file bytes. The accepted complete fingerprint is reused from iteration 008:
`F56283F5A82EF5FC37F00D01C0565357057D9FBEE1D95CAB125F6DF61B6C5441`.
`go mod verify` passed with native exit `0` and `all modules verified`.

Fresh external `GOCACHE` and `GOTMPDIR` paths are under
`C:/Users/ITDan/AppData/Local/Temp/pscan-06-i009-01a09228-1204-70b1-8b7e-1be433b2c496/`.
Every Go command will use the literal pinned executable and cache with
`GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local`, `GOENV=off`,
`GOTELEMETRY=off` and `GOVCS=*:off`.

## Bounded implementation surface

Implementation is confined to
`tests/integration/supply-chain/release_test.go`, new iteration-009 evidence
and the four authority-listed living-state documents. The authority file and
all pre-existing evidence are immutable. A need for any other implementation
path, formatter hunk, production verifier change or broader helper admission
stops fail-closed.
