# PSCAN-05 Live Primary-Source Preflight

## Scope and starting state

- Task: `PSCAN-05` only.
- Verified: `2026-09-02`, before dependency intake, download, build, scanner
  execution, test execution or new pinning in this implementation session.
- Activation commit: `50b418609c1f9927c0ecd5d740aba6d7bff11f58` on
  `main`; parent `1f0890878518de32a55ceb8d7b97430c4f3d2f2b`; tree
  `5b3e69591dc01466ae355aa7fae1631907e4857f`.
- Start state: exact activation HEAD, clean worktree, and PSCAN-05 is the sole
  activated task. Every successor remains unselected and inactive.
- Controlling contract SHA-256:
  `8a034701867e366bf37adea6e4f47e3ffb7a53425463ceff4fb280bcb4296e74`.
- Accepted PSCAN-04 closeout and acceptance SHA-256 values respectively:
  `62c45a85c02697e1c2878083277a56a3d810c3228712a0ea822aff2f81b3be2a`
  and
  `dbd2b2b6238faede70ba84b3fd5aa01f162fa944281352f51ee56f934accae77`.
- DEC-001 SHA-256:
  `df73ffb83071277765535b8f1d9db95c41f35b695b167b0c2958fa9f32a7f1b6`.

The complete ordered reading map and every tracked file in the PSCAN-05
implementation allowlist were inspected before editing.

## Schema and JSON boundary

The scanner remains on JSON Schema draft 2020-12 for its scanner-owned
contracts. The current JSON Schema core specification defines object members as
unordered and makes duplicate-name behavior undefined. It also treats unknown
schema keywords as annotations. PSCAN-05 therefore does not infer compatibility
from a generic JSON parser: project projection inputs reject duplicate member
names, enforce closed current-minor fields, reject unknown majors and unknown
required features, and admit a future minor only through an explicit bounded
compatibility window.

Primary record:

- `https://json-schema.org/draft/2020-12/json-schema-core`

## Canonicalization and signature boundary

RFC 8785 establishes that cryptographic hashing/signing requires an invariant
representation and defines JCS through I-JSON, ECMAScript primitive
serialization and deterministic property sorting. PSCAN-05 does not implement
or claim JCS compatibility. It avoids a new canonicalization dependency by
binding the SHA-256 of the exact admitted JSON bytes and verifying signatures
over a versioned, domain-separated, length-prefixed binary preimage containing
only declared schema, adapter, key and digest fields. Formatting changes thus
change the bound digest and cannot be normalized into the same authority.

RFC 8032 defines Ed25519. Go 1.27.0's standard-library `crypto/ed25519`
documents 32-byte public keys, 64-byte signatures and message verification; it
also documents that `Verify` panics for a wrong public-key length. PSCAN-05
checks key and signature sizes before verification. The scanner accepts only
caller-supplied public verification material, never a private key, never signs,
and never issues a receipt.

Primary records:

- `https://www.rfc-editor.org/rfc/rfc8785.html`
- `https://www.rfc-editor.org/rfc/rfc8032.html`
- `https://pkg.go.dev/crypto/ed25519@go1.27.0`

## Go toolchain and licence boundary

Official Go release history lists Go 1.27.0 as released on 2026-08-19 and Go
1.27.1 as released on 2026-09-01. Go 1.27.1 contains fixes to the compiler,
runtime, `encoding/json`, `os` and other packages. Go 1.27 remains a supported
major line. The repository's `go 1.27.0` language baseline and dependency graph
remain unchanged. The official download record publishes SHA-256
`a3911b5e0e1b1053f25ed0675f4c1c6aad1e2bfcf253df2b9be4caabd2edd95d`
for `go1.27.1.windows-amd64.zip` and
`63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445`
for `go1.27.1.linux-amd64.tar.gz`. The official BSD-style Go licence remains
compatible with this standard-library-only addition.

Primary records:

- `https://go.dev/doc/devel/release`
- `https://go.dev/dl/?mode=json`
- `https://go.dev/doc/go1.27`
- `https://go.dev/LICENSE`

## Existing detector and module boundary

Gitleaks remains exact version `8.30.1`, source commit
`83d9cd684c87d95d656c1458ef04895a7f1cbd8e`, under its MIT licence. Three
independently preserved Linux amd64 binaries each re-hashed to
`657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586`.
The exact source still contains the previously proved fragment and binary-skip
behavior. PSCAN-05 changes no engine, rule, configuration, adapter, artifact or
coverage behavior and does not need to execute Gitleaks.

The existing `github.com/mholt/archives` v0.1.2 and
`github.com/h2non/filetype` v1.1.3 module graph remains unchanged and is not
used by the new policy or verifier packages. No new external dependency,
licence obligation, action, scanner, rule or configuration pin is introduced.

Primary records:

- `https://github.com/gitleaks/gitleaks/tree/83d9cd684c87d95d656c1458ef04895a7f1cbd8e`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/LICENSE`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/go.mod`
- `https://raw.githubusercontent.com/gitleaks/gitleaks/83d9cd684c87d95d656c1458ef04895a7f1cbd8e/sources/file.go`

## Fail-closed decision

`PROCEED_TO_PSCAN_05_ONLY`.

The selected standard-library verifier, exact-byte digest binding, independent
schema-family compatibility model and project-owned trust boundary are
mutually consistent. Any missing digest, unsupported family/version,
ambiguous semantics, invalid signature/reference, expired or broadened
exception, credential-class exception, changed receipt binding, revocation,
rollback, divergence, conflict or cross-scope state is non-pass.

No download, dependency change, build, scanner/test execution, new pin, remote
action, credential or signing-key handling, receipt issuance, publication,
spending, TruffleHog work or successor action occurred before this record was
created.

## Post-preflight toolchain execution note

Both Go 1.27.1 archives were downloaded only after the proceed decision and
matched the official hashes above. The Windows `go` binary identified itself as
`go1.27.1 windows/amd64`, but Windows Application Control blocked the extracted
assembler for malicious-binary reputation before compilation. No bypass was
attempted. The verified Linux archive executed as `go1.27.1 linux/amd64` inside
the already-present local `docker-desktop` WSL environment and was used for
tests, vet and Windows-target compilation. This is an implementation-validation
fallback, not cross-platform acceptance.
