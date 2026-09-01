# PSCAN-04 Implementation Record

## Sole-task claim

PSCAN-04 alone is claimed on `2026-09-02` from exact clean activation commit
`a21b030e1658f1f98ac4e4d001af12185d9ed311`, after the complete ordered reading
and live primary-source preflight in `PREFLIGHT.md`.

PASS-OUTCOME-SPEC-001 remains the controlling contract. PSCAN-05 and every
other successor remain unselected, inactive, unclaimed and unimplemented.

## Bounded design decision

PSCAN-04 will use the existing Go 1.27.0 module graph and strict standard-library
ZIP, TAR, gzip and JSON parsing. OCI Image Format v1.1.1 is a data-format
authority, not a runtime dependency. The broad existing archive dependencies
remain classification compatibility inputs only; they do not admit encrypted
or unsupported formats.

Implementation, validation commands and exact evidence are appended only under
the PSCAN-04 allowed paths. Protected specifications, traceability, decisions,
history, predecessor evidence and successor records remain untouched.

## Implemented boundary

- `internal/artifact`: strict directory, ZIP, TAR, gzip/TGZ, OCI-layout and
  Docker-save validation, deterministic normalization, exact private ledger and
  overlapping detector projections.
- `internal/redaction`: non-authoritative typed diagnostic whitelist, schema
  validator, panic conversion and raw/partial/encoded/hash canary oracle; the
  authoritative scanner outcome schema remains unchanged.
- `internal/cleanup` and `internal/workspace`: lifecycle cleanup plus
  scanner-marker-bound restart recovery after process death.
- `internal/engine` and `internal/engine/gitleaks`: bounded decoder-panic
  recovery, sealed environment validation and pass-capable verified artifact
  projection scanning with the pinned detector. `ScanArtifact` runs
  normalization and scanning under one cleanup-owned attempt deadline.
- `build/sandbox`: digest-pinned, network-disabled, read-only, capability-
  dropped, no-new-privileges and resource-bounded authoritative test runner.
- Synthetic unit, integration and acceptance tests under only PSCAN-04 allowed
  paths. No candidate file is executed and no real project data is used.

No module, licence or third-party notice changed because no dependency was
added, removed or upgraded.

## Live runtime identity and corrections

Docker client/server `29.7.2` was verified against Docker Desktop `4.87.0` on
Linux amd64. The local image resolved exactly to
`sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`
and reported `go version go1.27.0 linux/amd64`.

A mutable temporary path named `built/gitleaks-linux-amd64` did not match the
accepted digest and was rejected before execution. Independent preserved
`build-a`, `build-b` and `script-build` copies each matched exact SHA-256
`657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586`.
Only `build-a` was mounted read-only for PSCAN-04 validation.

The first full container run proved the new tests but exposed a harness-only
failure: legacy contract helpers attempted execution from the intentionally
`noexec` `/tmp`. `TMPDIR` was moved to the bounded executable `/work` tmpfs;
`/tmp` remains `noexec`. The corrected full run passed.

Docker 29.7.2's current save output for the already-pinned Go image was
inspected as official format evidence. It used content-addressed
`blobs/sha256/*` configuration and gzip layers plus OCI root metadata, with
seven selected layers and seven ordered DiffIDs. The strict Docker-save subset
was corrected to support that current form as well as legacy digest-JSON/tar
layers. The 311,128,576-byte probe was removed after inspection and was not
used as candidate data or retained.

The implementation remains a bounded internal component stage. PSCAN-04 does
not modify the out-of-allowlist public runner and makes no public CLI, consumer,
deployment, production, Windows-runtime or complete cross-platform readiness
claim.
