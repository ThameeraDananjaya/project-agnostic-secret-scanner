# PSCAN-06 Correction C2 iteration 010 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Iteration: `010`
- Title: exact bundle-sibling helper output-path correction
- Event: owner selection of the minimal correction after independent finding
  I009-R01
- Decision date: `2026-09-12`
- Authority session: `01a092ac-14df-78f1-8604-0e2826b55c39`
- Branch: `main`
- Required starting HEAD before this authority bundle:
  `ddb67e1112ad0b263c98933cfad3b62263cf4ba4`
- Required starting tree:
  `0c2447bea9d4ac8a2559a687837cc7dad12ddf59`
- Required starting HEAD direct parent:
  `d89b033a6bc99c7e7886fffbbe677d047a479723`
- Correction authority commit: this resulting authority-bundle commit
- Pre-change tracked, staged and non-ignored untracked state: clean
- Preserved ignored workspace material: 66 entries, all under
  `graphify-out/**`; excluded from product authority and not refreshed,
  deleted, moved or staged
- State after commit: owner-approved; unclaimed; not implemented; not accepted
- Iteration 009 state: candidate independently rejected; no repair performed
- Recovery R5 state: terminally failed and consumed
- Recovery R6 state: not authorized and not executable until iteration 010 is
  independently accepted locally and a later exact R6 authority separately
  re-proves and approves every action-time gate
- Sole open task: `PSCAN-06`
- Successor selected or activated: none
- Spending and subscription ceiling: zero; no purchase or subscription action
  permitted

This is a bounded local iteration inside already activated PSCAN-06 Correction
C2. The owner's autonomous-local-completion direction selects the minimal
single-test-file correction below while preserving every stricter repository
gate. It does not approve online subscription purchases, nonzero spend or any
remote/runtime/release action.

This session records and commits authority only. It performs no implementation
and runs no Go command. A genuinely fresh implementation session whose
`CODEX_THREAD_ID` differs from this authority session must begin from the exact
resulting authority commit and claim only Correction C2 iteration 010.

## Controlling and immutable inputs

- `PASS-OUTCOME-SPEC-001` committed-byte SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Iteration-009 authority SHA-256:
  `CFDBE92EC51305FBEEBCDB47E6EA0BDA0CA03F95F6D5ABFAE518BF1DE378EB7E`
- Iteration-009 preflight SHA-256:
  `D3E9F272299EB9DD842C3048E4DC2FC1931A3F349A9F7D6D76C3705A95CC6EAA`
- Iteration-009 continuation authority SHA-256:
  `C0CCF3BD9A646872EBB98F460FB339ACAFEF013AE96F5E751A613F0D61C8D7D1`
- Iteration-009 continuation preflight SHA-256:
  `A499C832CF23E20F3E94BBC01B6CC5451B39A0FC34920F19B31811D40186DB2F`
- Iteration-009 continuation implementation SHA-256:
  `F0BD27002C461600A0B3108ABD283BD2231E4DD6931E512296CFCDC8F1DC6C2B`
- Iteration-009 continuation author-validation SHA-256:
  `FD53FAB4905C2D7CAB0D9B7E33318CEDDFA1EF6818A57938246224AB7C17921A`
- Iteration-009 independent rejection SHA-256:
  `26DC2E59E18B9F2125414A3E8BFDB5ED5244658118059025C4479D998D2CC896`
- Rejected candidate/tree:
  `d89b033a6bc99c7e7886fffbbe677d047a479723` /
  `a9ade332849f4345f2e7e8f032504dccc9e60e5a`
- Selected test-file Git blob / committed-byte SHA-256 / byte length:
  `106f2efacb44177fa08d62dd86afefa1f7a52d16` /
  `6388688F974D9C72C722E73E304CB984AA7269C17A9161CA26864FA2847C4FDB` /
  `25224`

All iteration-009 authorities and candidate evidence remain immutable. The
complete pinned offline matrix is passing evidence, not acceptance. The exact
I009-R01 mismatch controls this iteration; no other candidate behavior is
reopened or reinterpreted.

## Blocking finding I009-R01

The current test-only helper derives its trust root from the directory of
argument 4 (`--trusted-root`) and correctly requires the trusted root, bundle
and manifest to be contained regular files without admitted symlink escape.
It then incorrectly derives `arguments.txt` from that trust-root directory.

The admitted domain includes a regular bundle in a nested directory beneath
the same trusted root. For that valid shape, the iteration-009 authority
requires output beside the lexical cleaned bundle path, but the candidate
returns and writes the distinct trusted-root-sibling output. The external-only
independent probe reproduced the mismatch. No production verifier defect,
workflow defect or artifact-normalizer defect exists.

## Selected objective and exact correction

Within `tests/integration/supply-chain/release_test.go` only:

1. preserve the existing admission root derived from
   `filepath.Clean(filepath.Dir(arguments[4]))`;
2. preserve `plainDirectory`, `containedRegularFile`, `containedPath`,
   `plainRegularMode` and their trusted-root, bundle and manifest admission
   semantics;
3. for an otherwise admitted invocation, derive the output exactly as
   `filepath.Join(filepath.Dir(filepath.Clean(arguments[2])), "arguments.txt")`;
4. perform the existing pre-existing-output `os.Lstat` rejection against that
   exact derived bundle-sibling path; and
5. preserve the `os.OpenFile` write with
   `os.O_WRONLY|os.O_CREATE|os.O_EXCL`, mode `0o600`, one argument per line,
   close/error handling and success exit `0`.

The output path is a lexical function of argument 2 after that argument has
passed the existing containment, regular-file and symlink admission. It must
not be derived from argument 4, an evaluated target, the manifest, current
directory, process environment or any new root. No broader helper admission,
production change or compatibility fallback is selected.

## Mandatory preservation

The correction and regression must preserve without weakening:

- exact length 18 and all exact fixed positions 0 through 17;
- non-empty values at positions 2, 4, 6, 8, 10, 12, 14, 16 and 17;
- exact `COSIGN_YES=false` and no custom environment sentinel;
- current-test-executable path and digest binding;
- the trusted-root-derived admission root and all contained regular-file,
  lexical escape, evaluated symlink/reparse and directory rejection;
- all six certificate identity, OIDC issuer, repository, ref, workflow SHA and
  trigger claim assertions;
- ordinary and near-miss fallback to exactly one `m.Run()` call with no
  recursive test execution and no output creation or modification;
- exclusive `arguments.txt` creation through `O_EXCL`;
- exact iteration-007 workflow normalization counts `1,1,8,3,1`, byte equality,
  `# v4.2.2` presence and `# v4.2.1` absence;
- the Go 1.27.1 formatter-clean candidate bytes outside this exact helper and
  regression correction; and
- every production, workflow, schema, builder, Docker materializer, verifier,
  artifact normalizer, other test and historical evidence byte.

## Required positive nested-bundle regression

The implementation must add a positive regression in the selected test file
using a test-owned temporary root whose paths are physically distinct:

```text
<root>/trusted-root.json
<root>/release-manifest.json
<root>/nested/release-manifest.sigstore.json
```

All three inputs must be regular files and the exact 18-position invocation
must use `COSIGN_YES=false`. The regression must independently prove:

1. `cosignHelperArgumentsPath` admits the invocation and returns only
   `filepath.Join(filepath.Dir(filepath.Clean(arguments[2])), "arguments.txt")`;
2. `testMainExitCode` returns success and writes the received arguments only
   to that exact bundle-sibling path;
3. the captured bytes retain all 18 ordered arguments and the six claim
   positions required by the authority;
4. `filepath.Join(filepath.Dir(filepath.Clean(arguments[4])), "arguments.txt")`
   is a distinct path and remains absent; and
5. no recursive `m.Run()` fallback occurs for the admitted invocation.

The same bounded regression surface must continue to prove non-entry and no
output creation or modification for wrong length, command, order/position,
empty values, absent/wrong `COSIGN_YES`, applicable bundle/trusted-root/
manifest escapes, symlink or reparse substitutions, directories and a
pre-existing exact bundle-sibling `arguments.txt`. Skips may cover only a
host capability that makes a specific symlink/reparse fixture impossible; a
skip cannot turn a capable hostile case into pass.

## Path boundaries

This authority-recording session may change, stage and commit only:

```text
README.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-010-AUTHORITY.md
```

Correction C2 iteration 010 implementation may change only:

```text
tests/integration/supply-chain/release_test.go
README.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-010-*.md (new files only)
```

The committed iteration-010 authority, every iteration-009 authority and
evidence record, `tests/unit/artifact/normalize_test.go`, every production
path, every other test file and every absent path are forbidden. A need for
another path stops fail-closed.

## Fresh implementation session and exact source trust

The implementation session must have a `CODEX_THREAD_ID` different from
`01a092ac-14df-78f1-8604-0e2826b55c39`. Before any edit or Go command it must:

1. start on branch `main` at the exact resulting authority commit and record
   its commit, tree and direct parent;
2. prove the authority commit's diff from
   `ddb67e1112ad0b263c98933cfad3b62263cf4ba4` contains exactly the five
   authority-session paths above;
3. prove an empty index and clean non-ignored tracked/untracked state, with all
   ignored work confined to preserved `graphify-out/**`;
4. re-read the controlling contract, live task/tracker/reading map, both
   iteration-009 authorities, candidate evidence, rejection and this authority;
5. prove exact current-checkout Git object/index/path/mode/raw-byte trust,
   allowing only the repository's already proved canonical whole-file
   LF-to-CRLF projections; and
6. create a separate fresh no-hardlink materialization of the exact authority
   commit with checkout conversion disabled and prove every tracked file
   raw-equal to its committed blob.

Any branch, identity, byte, path, mode, index, configuration, ignored-root,
materialization or session-freshness mismatch stops before implementation.

## Pinned offline validation route

Implementation and fresh independent review must each use only:

```text
Go: C:\Users\ITDan\AppData\Local\Temp\pscan-06-local-toolchain-b7a0fb2-01a09168\toolchains\go1.27.1-windows-amd64\go\bin\go.exe
Module cache: C:\Users\ITDan\AppData\Local\Temp\pscan-06-local-toolchain-b7a0fb2-01a09168\gomodcache
```

The Go SHA-256 must be
`D3CCDB604EAFA6031133AEFE1A3DB24F0BB7362B857BC2125AC4E4C178B4B490`;
`go version` must be exactly `go version go1.27.1 windows/amd64`; and adjacent
`gofmt.exe` SHA-256 must be
`AB730B8446C0C2369BE901F263FB7433EA0D2D053D3A3584CA543AAFF3B667CC`.

Admit the literal frozen module cache as read/execute-only input, prove its
access and accepted identity before use, run `go mod verify`, and prove the
complete file/directory/byte/metadata inventory unchanged after all commands.
Use fresh session-specific external `GOCACHE`, `GOTMPDIR` and compile outputs.
Set `GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local`, `GOENV=off`,
`GOTELEMETRY=off` and `GOVCS=*:off`. Any missing identity, denied required
read, successful cache write, inventory drift, attempted download need or
network requirement stops closed.

Assert every native exit code and require:

1. `gofmt` on the selected test file and empty final `gofmt -d` output;
2. targeted exact helper/claim capture pass;
3. targeted positive contained nested-bundle return/write/no-trusted-root-
   sibling regression pass;
4. targeted hostile escape/symlink/directory/pre-existing-output and all other
   near-miss non-entry regressions pass;
5. targeted iteration-007 workflow normalization and `# v4.2.2` preservation
   pass;
6. targeted artifact case-collision and directory-symlink preservation tests
   pass or retain only their already authorized host-capability skips;
7. `go test -count=1 ./tests/unit/artifact` pass;
8. `go test -count=1 ./tests/integration/supply-chain` pass;
9. `go test -count=1 ./...` pass;
10. `go vet ./...` pass; and
11. `GOOS=linux GOARCH=amd64 CGO_ENABLED=0` compile-only coverage for every
    package to fresh external outputs, with no compiled result executed.

Actual Linux runtime remains reserved for a separately authorized Recovery R6
and is not substituted by cross-compilation. Recovery R6 cannot be considered
or authorized from this iteration until the bounded candidate is independently
accepted locally; local acceptance itself still grants no R6 action authority.

## Candidate and independent-review gate

The implementation session must record preflight before editing, make only the
selected correction, run the complete offline matrix, update only new
iteration-010 evidence and living documents, and create exactly one candidate
commit. That commit must have the authority commit as direct parent and contain
only the iteration-010 implementation boundary. Author validation is not
acceptance.

A genuinely fresh independent review session with a different
`CODEX_THREAD_ID` must materialize and inspect that exact candidate commit,
prove its parent/tree/path confinement and immutable-source identities, review
the exact helper domain and bundle-sibling derivation, independently exercise
the positive nested-bundle and hostile/pre-existing cases, and rerun the
proportionate complete pinned offline matrix before recording acceptance or
rejection. The reviewer may not repair the candidate.

## Acceptance matrix

| ID | Required result |
|---|---|
| I010-A01 | A genuinely fresh session claims only iteration 010 from the exact clean committed authority; its ID differs from the authority session. |
| I010-A02 | Exact current-checkout source trust, separate raw-equal materialization, ignored-root confinement and authority-commit five-path confinement pass. |
| I010-A03 | The selected test file is the sole implementation-bearing path; every production, other-test and pre-existing evidence byte is preserved. |
| I010-A04 | Every admitted invocation derives output exactly beside lexical cleaned argument 2 while retaining the trusted-root-derived admission root and all existing containment/regular-file/symlink checks. |
| I010-A05 | The contained nested-bundle regression proves exact helper return, exclusive write, ordered capture, no fallback and absence of the distinct trusted-root-sibling output. |
| I010-A06 | Escapes, symlinks/reparse substitutions, directories, pre-existing exact output and all existing near misses reject without output creation/modification or recursion. |
| I010-A07 | Exact 18 positions, `COSIGN_YES=false`, no custom sentinel, all six claims, `O_EXCL`, workflow counts, byte equality and `# v4.2.2` preservation pass. |
| I010-A08 | Pinned Go 1.27.1, frozen-cache invariance, `go mod verify`, formatter, targeted, both complete packages, repository test/vet and all-package Linux-amd64 compile-only gates pass offline. |
| I010-A09 | Exactly one bounded candidate commit is author-validated, then a genuinely fresh independent reviewer inspects and tests it without repair. |
| I010-A10 | Recovery R5 remains terminal; Recovery R6 remains unauthorized; no Docker, network, remote, tag, workflow, signing, publication, spend, subscription or successor action occurs. |

## Authority-session commit and stop boundary

This authority session must create exactly one governance commit whose direct
parent is `ddb67e1112ad0b263c98933cfad3b62263cf4ba4` and whose diff contains
exactly the five authority-session paths. It must prove the resulting commit,
tree, parent, exact path list, empty index and clean non-ignored final state,
with ignored work still confined to `graphify-out/**`, then stop.

This session performs no Go command, formatting, test, vet, compilation,
Docker action, dependency acquisition, network request, remote read or
mutation, tag operation, workflow dispatch, signing, attestation, draft,
release, publication, spending, subscription, Recovery R6 or successor action.
It selects or activates no successor and makes no acceptance claim.
