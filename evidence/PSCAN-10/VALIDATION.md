# PSCAN-10 Implementation Validation

## Exact bindings

- Activation base: `9053b37d799b19e3d98aeca0ae853296971ab09d`.
- Gitleaks source: `8.30.1`, commit
  `83d9cd684c87d95d656c1458ef04895a7f1cbd8e`.
- Build variant: Go standard-library regex, Go `1.27.0`, `CGO_ENABLED=0`.
- Windows amd64 binary SHA-256:
  `b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178`.
- Linux amd64 binary SHA-256:
  `657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586`.
- Rule pack SHA-256:
  `51e903dac96834d1dc87999e05971cf9a30c2a3ba3d8aaa57d34f494e8a942c9`.
- Empty ignore file SHA-256:
  `3cdd737e34cb301ffc19861f903cf59441ff57453ddf70bd5b8c1bd63389aebe`.
- Rule proof: 224 parsed rules, zero unbounded rules, conservative maximum
  UTF-8 byte span 4,020.
- Compiled-module licence manifest: 63 modules, LF-normalized SHA-256
  `0b8c06685c7d16fd6d685ae0e5e481d85396d2e3c64671151f07634a7bb40b60`;
  a fresh stdlib-variant collection reproduced the complete tree without
  differences.

## Executed validation

1. The repository build script rebuilt Windows/Linux amd64 from the exact
   source and Go inputs and reproduced both expected binary hashes.
2. Two independent source builds per platform were byte-identical.
3. Native Windows Gitleaks canaries returned clean exit `0`, synthetic-finding
   exit `11`, exact version `8.30.1`, and redacted candidate fields.
4. A network-disabled, read-only Linux amd64 container ran `go test -count=1
   ./...` with the exact source-built Gitleaks binary and pinned config. Every
   package passed.
5. The actual-engine suite separately passed text head, binary head,
   deleted-binary history, unchanged head tree, merge parent edge, rename,
   clean/out-of-range, skip-path, first byte, last byte, internal no-whitespace
   fragment boundary and exact 4,020-byte maximum-span cases.
6. The gated declared-maximum fixture passed an actual network-disabled
   Gitleaks run at exactly 512 MiB in 115.59 seconds. Profile unit tests reject
   every file/object/total/chunk/report/time/memory/process value above its
   named bound and accept the below/at values.
7. Nine Windows amd64 Go test packages compiled successfully. Their executable
   hashes were recorded in the session output. Host application control later
   denied launching newly generated test executables; this is a platform
   execution limitation, not converted to pass. The exact Windows Gitleaks
   binary itself executed both behavioral canaries.
8. Windows and network-disabled Linux binaries agreed on the clean/finding
   logical canaries (`0`/`11`) and the synthetic rule identity.

## Outcome-proof summary

- Range binding includes exact base, head, merge base, ordered commits, every
  parent edge, head-tree object OID and tree digest.
- Each Git blob is written from `cat-file`, size-checked, SHA-256 bound and
  independently re-hashed to its exact Git OID using hostile behavior-disabled
  Git commands.
- Each ledger row binds OID, mode, length, SHA-256, raw class, preparation
  version and every detector chunk.
- Raw classification calls every independently registered archive/compression
  family in the same `github.com/mholt/archives` v0.1.2 stack as the pinned
  detector. Archive, compression, container and ambiguous results are non-pass
  before preparation; ordinary text and binary are admitted.
- Detector payload windows are 90,000 bytes with 4,019-byte overlap. Complete
  prepared files stay below 100,000 bytes, so the pinned detector produces one
  fragment per file. Every chunk has a unique integrity-bound same-invocation
  witness, and continuity verification proves every original byte is mapped.
- Missing, duplicate, mutated, skipped, unbound, unsupported or over-limit
  input and missing/duplicate proof are explicit non-pass.

No credential, signature, spend, remote, publication, release, deployment,
TruffleHog action or successor action occurred. Candidate fixtures are
invalid-by-construction synthetic values only.
