# PSCAN-06 Correction C2 iteration 009 independent review rejection

Date: `2026-09-12`
Reviewer session: `01a0929e-4232-7ec0-a99a-894a8e90ae74`

## Verdict

`LOCAL_ACCEPTANCE_REJECTED_EXACT_HELPER_OUTPUT_PATH`

Candidate `d89b033a6bc99c7e7886fffbbe677d047a479723`, tree
`a9ade332849f4345f2e7e8f032504dccc9e60e5a`, is not independently accepted.
The complete required offline matrix passes, but one exact helper-admission
requirement does not. Passing tests do not override the authority mismatch.

This review does not repair the candidate. PSCAN-06 remains open and
unaccepted overall. Recovery R5 remains terminal and consumed; Recovery R6
remains unauthorized. No successor is selected or activated.

## Freshness and candidate identity

- Branch: `main`.
- Starting HEAD: `d89b033a6bc99c7e7886fffbbe677d047a479723`.
- Starting tree: `a9ade332849f4345f2e7e8f032504dccc9e60e5a`.
- Direct parent: `f4c77e239adbb83f5af90741373b22f0bda2c4da`.
- Reviewer session differs from the implementation, authority, approval and
  diagnostic sessions named in the Iteration 009 evidence.
- The starting nonignored tree and index were clean. Ignored paths were
  confined to preserved `graphify-out/**`.
- A fresh no-hardlink materialization with `core.autocrlf=false` and
  `core.eol=lf` reproduced all 390 tracked committed blobs byte-for-byte with
  zero mismatches and the same HEAD, tree and parent.

The candidate commit changes exactly these ten authority-listed paths:

```text
README.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/PSCAN-06.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-CONTINUATION-AUTHOR-VALIDATION.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-CONTINUATION-IMPLEMENTATION.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-CONTINUATION-PREFLIGHT.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-PREFLIGHT.md
tests/integration/supply-chain/release_test.go
tests/unit/artifact/normalize_test.go
```

No production path changed. `git diff --check` passed.

## Passing independent validation

The exact Go executable, version and formatter matched the authority:

```text
go.exe SHA-256 D3CCDB604EAFA6031133AEFE1A3DB24F0BB7362B857BC2125AC4E4C178B4B490
go version go1.27.1 windows/amd64
gofmt.exe SHA-256 AB730B8446C0C2369BE901F263FB7433EA0D2D053D3A3584CA543AAFF3B667CC
```

The frozen module cache retained `(OI)(CI)(RX)` access. Before and after the
authoritative review run it had 34,447 files, 7,678 directories,
1,132,106,847 file bytes and reviewer metadata SHA-256
`FC6A6D2F410B6223A161DE7E9D532E9B7D232FD3B77CD052AA48B4AE46780C4A`.
`go mod verify` returned exit `0` with `all modules verified`. The denied Go
stat-cache write after Linux compilation changed no cache entry or metadata.

All authoritative review commands used fresh reviewer-specific external
`GOCACHE` and `GOTMPDIR`, `GOPROXY=off`, `GOSUMDB=off`,
`GOTOOLCHAIN=local`, `GOENV=off`, `GOTELEMETRY=off` and `GOVCS=*:off`.
An initial wrapper was cut off by its local command timeout before producing
usable results; the complete targeted phase was restarted with a new fresh
cache/temp pair. The following authoritative rerun results all had native
exit `0`:

```text
go mod verify
gofmt -d tests/integration/supply-chain/release_test.go tests/unit/artifact/normalize_test.go
go test ./tests/integration/supply-chain -run '^(TestCosignCommandVerifierBindsEveryGitHubWorkflowClaim|TestCosignTestMainHelperRejectsNearMissesWithoutRecursiveExecution)$' -count=1 -v
go test ./tests/integration/supply-chain -run '^TestIteration007RepositoryIdentityAgreementAndPreservation$' -count=1 -v
go test ./tests/unit/artifact -run '^TestDirectoryLinksAndCaseCollisionsFailClosed/case_collision$' -count=1 -v
go test ./tests/unit/artifact -run '^TestDirectoryLinksAndCaseCollisionsFailClosed/directory_symlink$' -count=1 -v
go test -count=1 ./tests/unit/artifact
go test -count=1 ./tests/integration/supply-chain
go test -count=1 ./...
go vet ./...
```

The case-collision subtest skipped only after `os.SameFile` proved `A` and `a`
were the same physical object. The directory-symlink subtest passed its exact
`RejectUnsafe` assertion. No `runtime.GOOS` branch exists.

The workflow tokens occur exactly `1,1,8,3,1` times. Exact ordered
normalization equals the historical C2 workflow byte-for-byte; `# v4.2.2`
occurs once and `# v4.2.1` is absent. The parent file independently reproduced
the 66-line, four-hunk Go 1.27.1 formatting projection, and the candidate is
formatter-clean.

With `GOOS=linux`, `GOARCH=amd64` and `CGO_ENABLED=0`, `go list ./...`
returned 36 packages and `go test -c` compiled every package to fresh external
outputs. No result was executed.

## Blocking finding I009-R01

The original Iteration 009 authority requires the admitted helper to derive
`arguments.txt` beside `--bundle`. Candidate function
`cosignHelperArgumentsPath` instead sets `root` to the directory containing
argument 4 (`--trusted-root`) and returns `filepath.Join(root,
"arguments.txt")`.

This is not equivalent for the helper's admitted domain. The same function
allows the bundle to be any contained regular file below that root. A bundle
in a nested directory, with the trusted root and manifest as regular contained
files in the test temporary root, satisfies the exact 18-position argument
shape and `COSIGN_YES=false`. The candidate admits it but returns the trusted-
root sibling path rather than the bundle sibling path.

An external-only reviewer probe against the exact candidate used that fully
contained regular-file arrangement. The result was:

```text
arguments path = <temp-root>\arguments.txt
want beside --bundle at <temp-root>\nested\arguments.txt
PROBE_EXIT=1
```

The probe file existed only in a separate external diagnostic clone and is not
part of the candidate or this review commit. The candidate's twelve recorded
near-miss subtests pass, but none covers a contained nested bundle, so they do
not expose this broadened admitted case.

I009-A04 therefore fails: the current-test-executable helper does not satisfy
the exact output-path derivation required for every invocation it admits.
I009-A09 consequently cannot advance to acceptance.

## Boundaries preserved

This review performed no implementation edit, Docker action, dependency
download, network request, remote read or mutation, tag operation, workflow
dispatch, signing, attestation, draft, release, publication, spending,
subscription, Recovery R6 or successor action. No repair is authorized by
this review record.
