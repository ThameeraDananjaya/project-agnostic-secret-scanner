# PSCAN-06 Correction C2 iteration 009 continuation author validation

Date: `2026-09-12`
Continuation session: `01a09293-16a9-71f1-8e93-440952d83e2c`

## Offline route

Every Go command used the literal pinned executable
`C:\Users\ITDan\AppData\Local\Temp\pscan-06-local-toolchain-b7a0fb2-01a09168\toolchains\go1.27.1-windows-amd64\go\bin\go.exe`,
SHA-256 `D3CCDB604EAFA6031133AEFE1A3DB24F0BB7362B857BC2125AC4E4C178B4B490`,
with exact version `go version go1.27.1 windows/amd64`. Adjacent `gofmt.exe`
had SHA-256
`AB730B8446C0C2369BE901F263FB7433EA0D2D053D3A3584CA543AAFF3B667CC`.

The literal frozen `GOMODCACHE` retained `(OI)(CI)(RX)` access before and
after use. Fresh `GOCACHE`, `GOTMPDIR` and compile outputs were under:

```text
C:\Users\ITDan\AppData\Local\Temp\pscan-06-i009-cont-01a09293-16a9-71f1-8e93-440952d83e2c\
```

Every command used `GOPROXY=off`, `GOSUMDB=off`, `GOTOOLCHAIN=local`,
`GOENV=off`, `GOTELEMETRY=off` and `GOVCS=*:off`.

## Exact command results

All native exits were `0`:

```text
go mod verify
  all modules verified

gofmt -w tests/integration/supply-chain/release_test.go tests/unit/artifact/normalize_test.go
gofmt -d tests/integration/supply-chain/release_test.go tests/unit/artifact/normalize_test.go
  formatter-check output: 0 bytes

go test ./tests/integration/supply-chain -run '^(TestCosignCommandVerifierBindsEveryGitHubWorkflowClaim|TestCosignTestMainHelperRejectsNearMissesWithoutRecursiveExecution)$' -count=1 -v
  both top-level tests and all 12 near-miss subtests passed

go test ./tests/integration/supply-chain -run '^TestIteration007RepositoryIdentityAgreementAndPreservation$' -count=1 -v
  passed

go test ./tests/unit/artifact -run '^TestDirectoryLinksAndCaseCollisionsFailClosed/case_collision$' -count=1 -v
  parent passed; case-collision subtest skipped only after os.SameFile proof

go test ./tests/unit/artifact -run '^TestDirectoryLinksAndCaseCollisionsFailClosed/directory_symlink$' -count=1 -v
  parent and directory-symlink subtest passed

go test -count=1 ./tests/unit/artifact
go test -count=1 ./tests/integration/supply-chain
go test -count=1 ./...
go vet ./...
  all passed

GOOS=linux GOARCH=amd64 CGO_ENABLED=0
go list ./...
for each of 36 packages: go test -c -o <fresh-external-output> <package>
  all 36 packages cross-compiled; no result executed

git diff --check
  passed
```

The frozen module-cache inventory was unchanged before and after: 34,447
files, 7,678 directories, 1,132,106,847 file bytes and sorted
path/type/length/attributes/last-write metadata SHA-256
`B83C371750D2345E7A64055ECB00058EE9C47C77FAA47A48E46EA29C0F0CD59B`.
Go emitted one denied stat-cache write attempt after the compile-only sweep;
the `(RX)` ACL blocked it, every command still exited `0`, and the complete
inventory and metadata identity proved that no module-cache write occurred.
No dependency download or network access occurred.

This is author validation, not independent acceptance. A genuinely fresh
review must inspect the exact committed candidate and rerun the proportionate
complete offline matrix. PSCAN-06 remains open and unaccepted overall.
