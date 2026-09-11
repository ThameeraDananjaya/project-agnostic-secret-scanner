# PSCAN-06 Correction C2 iteration 010 author validation

Date: `2026-09-12`
Implementation session: `01a092b2-d833-74f3-99c0-a2def0764153`

## Offline route and immutable inputs

Every Go command used the literal pinned executable:

```text
C:\Users\ITDan\AppData\Local\Temp\pscan-06-local-toolchain-b7a0fb2-01a09168\toolchains\go1.27.1-windows-amd64\go\bin\go.exe
go version go1.27.1 windows/amd64
go.exe SHA-256 D3CCDB604EAFA6031133AEFE1A3DB24F0BB7362B857BC2125AC4E4C178B4B490
gofmt.exe SHA-256 AB730B8446C0C2369BE901F263FB7433EA0D2D053D3A3584CA543AAFF3B667CC
```

The literal frozen module cache retained inherited `(OI)(CI)(RX)` access.
Fresh `GOCACHE`, `GOTMPDIR` and compile outputs were under:

```text
C:\Users\ITDan\AppData\Local\Temp\pscan-06-i010-validation-01a092b2-d833-74f3-99c0-a2def0764153\
```

Every command set `GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local`,
`GOENV=off`, `GOTELEMETRY=off`, `GOVCS=*:off` and `GOWORK=off`. No dependency
download or network action occurred.

Before and after validation, the cache inventory was exactly 34,447 files,
7,678 directories, 1,132,106,847 file bytes, zero reparse points and sorted
relative-path/type/length/attributes/last-write metadata SHA-256
`7CC37B44BDEAAB4BEFB844FE707EE9A9A018CAB05650C1322A3D6A65173DEDFD`.
The Linux compile sweep emitted one denied stat-cache write attempt; the RX
ACL rejected it, every compile command still exited zero, and the complete
inventory proved no cache mutation.

## Exact command results

All asserted native exits were `0`:

```text
go mod verify
  all modules verified

gofmt -w tests/integration/supply-chain/release_test.go
gofmt -d tests/integration/supply-chain/release_test.go
  formatter-check output: 0 bytes

go test ./tests/integration/supply-chain -run '^(TestCosignCommandVerifierBindsEveryGitHubWorkflowClaim|TestCosignTestMainHelperWritesOnlyBesideContainedNestedBundle|TestCosignTestMainHelperRejectsNearMissesWithoutRecursiveExecution)$' -count=1 -v
  all three top-level tests passed
  nested-bundle return/write/ordered-capture/no-fallback/no-trusted-root-sibling passed
  all 18 hostile/environment/pre-existing/mode subtests passed with no skip

go test ./tests/integration/supply-chain -run '^TestIteration007RepositoryIdentityAgreementAndPreservation$' -count=1 -v
  passed

go test ./tests/unit/artifact -run '^TestDirectoryLinksAndCaseCollisionsFailClosed/(case_collision|directory_symlink)$' -count=1 -v
  case_collision skipped only after os.SameFile capability proof
  directory_symlink passed its exact rejection assertion

go test -count=1 ./tests/unit/artifact
go test -count=1 ./tests/integration/supply-chain
go test -count=1 ./...
go vet ./...
  all passed
```

Independent byte checks proved the workflow normalization counts exactly
`1,1,8,3,1`; the normalized workflow equals the historical C2 workflow
byte-for-byte; `# v4.2.2` occurs once and `# v4.2.1` is absent.
`git diff --check` passed.

With `GOOS=linux`, `GOARCH=amd64` and `CGO_ENABLED=0`, `go list ./...`
enumerated exactly 36 packages. `go test -c -o <fresh-external-output>`
returned zero for every enumerated package; 22 packages produced test
binaries and no compiled result was executed.

## Status

This is author validation only. The resulting one-commit candidate is not
accepted until a genuinely fresh independent reviewer materializes the exact
commit, inspects the helper domain and reruns the complete pinned offline
matrix without repair. Recovery R5 remains terminal and consumed; Recovery R6
and every Docker, network, remote, tag, workflow, signing, publication,
spending, subscription and successor action remain unauthorized.
