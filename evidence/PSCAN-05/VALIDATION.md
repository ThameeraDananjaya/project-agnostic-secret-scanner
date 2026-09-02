# PSCAN-05 Implementation Validation

## Toolchain

- Official Go `1.27.1` Windows amd64 archive SHA-256 matched
  `a3911b5e0e1b1053f25ed0675f4c1c6aad1e2bfcf253df2b9be4caabd2edd95d`.
  Windows Application Control blocked the extracted assembler before compile;
  no bypass was attempted.
- Official Go `1.27.1` Linux amd64 archive SHA-256 matched
  `63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445`
  and executed inside the existing local WSL environment.
- The repository module graph and `go.sum` were not changed.

## Passing checks

- PSCAN-05 unit, integration and acceptance packages passed ten consecutive
  runs. The predecessor request unit package passed a fresh regression run.
- `go vet` passed for `internal/policy`, `internal/verify`, `internal/request`,
  `internal/workspace` and every PSCAN-05 test package.
- Every repository package and test compiled for `windows/amd64` with
  `CGO_ENABLED=0`; the cross-target executables were not treated as executed
  tests.
- A repository-wide Linux run passed all packages except the pre-existing
  `tests/unit/gitinput` package. With Linux-native temporary storage, that
  package could not use the host Windows Git executable; with host temporary
  storage, its safe-clone boundary could not cross the WSL/Windows process
  boundary. The PSCAN-05 packages and all other packages passed.
- Both new JSON schemas parsed successfully. Diff whitespace validation passed.

The final Linux commands used the verified archive's `go` binary with
`GOTOOLCHAIN=local`, an isolated build/module cache and Linux-native temporary
storage:

```text
go test -count=10 ./tests/unit/policy ./tests/unit/verify ./tests/integration/policy ./tests/integration/isolation ./tests/acceptance/policy
go test ./tests/unit/request
go vet ./internal/policy ./internal/verify ./internal/request ./internal/workspace ./tests/unit/policy ./tests/unit/verify ./tests/integration/policy ./tests/integration/isolation ./tests/acceptance/policy
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go test -exec /bin/true ./...
```

The last command is compile-only: `/bin/true` prevents execution of the
cross-target Windows test binaries.

## Covered fail-closed cases

Exact-byte changes; malformed public keys/signatures; digest, adapter and trust
domain mismatch; duplicate JSON members; unknown majors and required features;
ambiguous minor/retirement declarations; invalid, expired, broadened or
misbound exceptions; credential classification; changed receipt bindings;
expired/future receipts; receipt-head rollback; chain truncation, sequence gap,
divergence, duplicate identity/target and revocation; parallel, sequential and
reused-host workspace/projection state.

## Limitations and state

This is author-side implementation validation, not independent acceptance.
Native Windows execution remains unavailable on this host because of the
recorded Application Control denial. No project instance, private key, signing,
receipt issuance, evidence store, promotion, remote action, scanner execution,
TruffleHog action or successor work occurred. PSCAN-05 remains implemented but
open pending independent acceptance; PSCAN-06 and PSCAN-07 remain unselected,
and PSCAN-08 remains inactive.
