# PSCAN-06 Correction C2 iteration 008 preflight and claim

Date: `2026-09-11`
Implementation session: `01a091b5-687a-7392-b366-4b71f8bbbe8d`
Authority approval session: `01a0919e-5c8a-7f92-bece-1f18193a788f`

## Claim and project state

This genuinely fresh saved-project session claims only PSCAN-06 Correction C2
iteration 008. PSCAN-06 is the sole open task. Recovery R5 remains terminally
failed and its sole protected C2-tag push and workflow dispatch remain
consumed. Recovery R6 is not authorized. PSCAN-07 remains proposed and
unselected; PSCAN-08 remains inactive and ineligible. No successor is
selected, activated, claimed or worked.

Iteration 008 repairs only the deterministic false failure in
`TestIteration007RepositoryIdentityAgreementAndPreservation`. It does not
reinterpret the rejected iteration-007 candidate as accepted and does not
change any product, workflow, schema, builder, materializer, verifier or
acceptance-test byte.

## Starting identity and source trust

- Saved-project root:
  `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`.
- Direct saved checkout: Git directory and common directory both resolved to
  `.git`.
- Branch: `main`.
- Exact starting authority commit:
  `4f86651de67a05db8ae0076c0c155bd184082af3`.
- Exact starting tree:
  `4ef9e344c1849c42509a61077e6a509fae8ef355`.
- Exact direct parent:
  `b62394ef172cb92b6394ef313993a2a6a5613b5b`.
- This implementation session differs from the authority approval session.
- Tracked and staged diffs were empty; non-ignored untracked status contained
  zero paths.
- All 66 ignored paths were confined to preserved `graphify-out/**`; none was
  refreshed, deleted, moved or staged.

The complete graphify-aware every-byte check proved:

```text
Exact graphify-aware source trust PASS commit=4f86651de67a05db8ae0076c0c155bd184082af3 tree=4ef9e344c1849c42509a61077e6a509fae8ef355 files=381 raw_equal=363 canonical_crlf=18 unexpected_mismatch=0
index=EXACT index_flags=EXACT unsupported_config=0 nonignored_untracked=0 ignored_only_graphify=66
```

The check used the repository source-trust tree, index, index-flag, Git-object
and canonical whole-file LF-to-CRLF algorithms. It verified SHA-1 object
format, exact complete tree/index path, mode and object identity, ordinary and
fsmonitor flags, absence of shortcut/sparse/promisor configuration, every
tracked working-tree byte and the explicit ignored-Graphify exception only.

## Readings, routing and immutable inputs

Before this claim the session read completely: `AGENTS.md`,
`PASS-OUTCOME-SPEC-001`, `TRACKER.md`, `PSCAN-06.md`, the PSCAN-06 reading map,
the iteration-007 authority, preflight, implementation and independent
rejection, and the iteration-008 authority. The project goal and tracker were
checked at the start.

The existing 108-node Graphify graph was queried read-only and routed the work
to the repository agreement, PSCAN tracker/task, release authority and supply-
chain validation surfaces. The graph predates iterations 007 and 008, so it
was used only as a navigation aid and live committed Git evidence controls.
No Graphify refresh, reflection, result save, deletion, movement or staging
occurred.

The committed-byte SHA-256 identities matched:

- `PASS-OUTCOME-SPEC-001`:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`;
- iteration-007 authority:
  `A17520815841415AD06357682423CF7C635B32AAFAB7EA2B0303E2540FE709E8`;
- iteration-007 preflight:
  `31624B5BAC5DB0EA751EB739795BE302469BE543CD56F675D4FBCC2729733FAC`;
- iteration-007 implementation:
  `F251946635F474DD4AF39071EBBE14CED0F48FA52AD5FF5AF4D130822D59083B`;
- iteration-007 independent rejection:
  `E8B90DDE62848577E129B7E2857FE6D0BF56BBF691440A06636013BE458FAC21`;
  and
- Recovery R5 terminal failure:
  `53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0`.

All nine critical workflow, schema, builder, Docker materializer, verifier and
acceptance-test Git blobs and SHA-256 identities listed in the iteration-008
authority matched exactly.

## Exact defect boundary

The original R6 workflow contains the five authority-listed tokens exactly
`1`, `1`, `8`, `3` and `1` times, respectively. It contains exactly one
`# v4.2.2` action-pin comment and zero `# v4.2.1` comments. The current test
retains the four identity-specific replacements followed by one forbidden
global `strings.ReplaceAll(normalizedWorkflow, "2.2", "2.1")` call. No
normalization or product byte was changed during preflight.

Implementation is confined to:

```text
tests/integration/supply-chain/release_test.go
README.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-008-*.md (new files only)
```

The iteration-008 authority and every pre-existing evidence file are
immutable. A need for any other path stops fail-closed.

## Offline Go route and frozen module-cache baseline

The exact owner-supplied executable exists at:

```text
C:\Users\ITDan\AppData\Local\Temp\pscan-06-local-toolchain-b7a0fb2-01a09168\toolchains\go1.27.1-windows-amd64\go\bin\go.exe
```

Its SHA-256 is
`D3CCDB604EAFA6031133AEFE1A3DB24F0BB7362B857BC2125AC4E4C178B4B490`, and
the exact native result was `go version go1.27.1 windows/amd64`, exit `0`.
The adjacent `gofmt.exe` exists with SHA-256
`AB730B8446C0C2369BE901F263FB7433EA0D2D053D3A3584CA543AAFF3B667CC`.

The frozen read-only module-cache input is:

```text
C:\Users\ITDan\AppData\Local\Temp\pscan-06-local-toolchain-b7a0fb2-01a09168\gomodcache
```

Its exact pre-use inventory contains 34,447 files, 7,678 directories and
1,132,106,847 file bytes. The SHA-256 fingerprint over sorted path, type,
length, attributes, last-write time and every file SHA-256 is
`F56283F5A82EF5FC37F00D01C0565357057D9FBEE1D95CAB125F6DF61B6C5441`.

Fresh session-specific external paths were created only at:

```text
GOCACHE=C:\Users\ITDan\AppData\Local\Temp\pscan-06-i008-01a091b5-687a-7392-b366-4b71f8bbbe8d\gocache
GOTMPDIR=C:\Users\ITDan\AppData\Local\Temp\pscan-06-i008-01a091b5-687a-7392-b366-4b71f8bbbe8d\gotmp
```

Every Go command must use the literal executable and frozen module-cache path
with `GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local`, `GOENV=off`,
`GOTELEMETRY=off`, `GOVCS=*:off` and explicit native exit-code enforcement.
No fallback executable, automatic toolchain acquisition, dependency download,
network lookup or module-cache mutation is permitted.

## Preserved boundary

No Docker command, dependency download, network action, remote read or
mutation, tag or workflow action, signing, publication, purchase,
subscription, spending or successor work occurred. Any identity, path,
token-count, critical-byte, toolchain, offline-test or source-trust mismatch
will stop this iteration rather than broaden it.
