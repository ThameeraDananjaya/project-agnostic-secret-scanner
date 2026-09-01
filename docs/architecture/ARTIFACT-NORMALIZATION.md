# PSCAN-04 Artifact Normalization

## Authority and scope

PASS-OUTCOME-SPEC-001 controls this bounded implementation. PSCAN-04 adds no
dependency and no consuming-project behavior. Go 1.27.0 standard-library ZIP,
TAR, gzip and JSON parsers normalize directory, ZIP, TAR, gzip/TGZ, a bounded
OCI Image Layout 1.1.1 image subset and strict legacy or Docker 29
content-addressed save input. Candidate material is data only: the
normalizer never invokes a shell, builder, package manager or candidate file.

## Fail-closed flow

1. Validate absolute private roots and reject links, Windows reparse points,
   hard links, special files and unsafe identities. Copy admitted directory
   bytes into a private deterministic snapshot, then revalidate the source
   tree identity and digest before any pass-capable inspection.
2. Classify raw bytes before transformation. Conflicting, encrypted,
   unsupported or malformed classes stop as non-pass.
3. Enforce depth 5, 100,000 entries, 512 MiB individual files, 1,000:1
   compression, PR 2 GiB/release 10 GiB expanded bytes, PR 15-minute/release
   60-minute timeout, declared memory and one scanner process.
4. Preflight ZIP central directories and TAR headers against remaining entry,
   individual-file and expanded-byte budgets before parser allocation or
   member materialization; then validate checksums, terminal records, paths,
   duplicates, case collisions and member types before recursive admission.
5. Validate OCI layout/index/manifest/config/layer descriptors, ordered
   `rootfs.diff_ids` and exact SHA-256 blob identities. The admitted subset is
   image manifests, JSON image configuration and tar or gzip-tar layers;
   descriptor URLs, inline data, zstd, artifacts and unsupported optional
   semantics are non-pass.
6. Validate Docker-save manifest/config/layer semantics, configuration digest,
   ordered DiffIDs and content-addressed extras before layer recursion. Bind
   only blobs reachable from the selected image; reject conflicting repeated
   descriptors and enforce cumulative de-duplicated semantic entry and byte
   budgets across the selected manifest, configuration and layers.
7. Bind every admitted raw object and expanded member into a private ledger.
   Write numeric private projections in overlapping 90,000-byte payload chunks
   with 4,019-byte overlap, keeping every detector file below the pinned
   100,000-byte Gitleaks fragment boundary.
8. Verify raw and projection digests, continuity, every framed projection
   payload against its exact raw byte range, and absence of duplicate or extra
   detector inputs before the pinned detector runs.

Gitleaks archive recursion remains disabled. Plaintext members are normalized
by the scanner; raw container bytes and every expanded member have their own
coverage-bearing projection. The coverage rule and product rules run together,
so a skipped projection cannot become a pass.

## Lifecycle and external boundary

All detailed paths, digests, entries, detector output and panic values remain
private. The redaction package provides a typed, whitelisted content-free
diagnostic envelope and rejects extra fields, framing and candidate-controlled
payloads. It does not replace the authoritative outcome schema.
The lifecycle wrapper removes the complete workspace for pass, fail,
indeterminate, timeout, cancellation and recovered panic. A durable
scanner-owned marker plus an external sibling lease permits safe cleanup of
abandoned workspaces after process death while preserving unleased or
structurally unsafe material. Arbitrary hostile same-user filesystem mutation
is outside this process boundary because such a principal can modify scanner
code and evidence as well; candidate artifacts are never executed as that
principal.

Authoritative validation uses the digest-pinned Go image and exact source-built
Gitleaks binary with no network, a read-only root and source, read-only module
and detector mounts, all capabilities dropped, no-new-privileges, default
seccomp, explicit non-secret environment and bounded CPU, memory, PID and
temporary filesystems.

This is an internal PSCAN-04 component stage. The existing public runner is
outside this task's path allowlist and remains unavailable; PSCAN-04 therefore
does not claim public CLI wiring, consumer readiness, deployment readiness,
Windows runtime parity or general cross-platform acceptance.
