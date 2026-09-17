# Build-only upload storage admission

The R6 workflow retains its exact product source, schema 2.2, platform assets,
build controls and single unsigned one-day transfer. A new guard runs once
immediately before upload. Any rejection prevents the upload step.

## Limits and complete artifact

- Exactly the existing 32 top-level regular distribution files; no hidden extra,
  nested directory, symlink/reparse point, hardlink or empty file.
- Maximum total payload: 251,658,240 bytes (240 MiB).
- Upload ZIP envelope reserve: 16,777,216 bytes (16 MiB).
- Maximum admitted complete transfer: 268,435,456 bytes (256 MiB).
- Exact product/tag/tree and tooling commit/ref/workflow identity; every one of
  the 31 manifest assets' size and SHA-256 and the complete unique checksum set.
- Per-file identity, length and metadata checked before/after streaming hashes.

The pinned upload action uses Zlib level 0. For 32 short top-level filenames,
16 MiB conservatively covers ZIP local/central headers, even maximum legal
extra/comment fields, end/ZIP64 records and stored-DEFLATE framing. At the
240 MiB payload cap, even a 1% framing allowance plus 32 times 256 KiB and
1 MiB end metadata is under 16 MiB. The action's entry count, compression-level
0 and unchanged exact pin are prerequisites; this is not a general bound for
arbitrary uploader implementations. No compression saving is assumed.

## Admission record

`PSCAN_BUILD_STORAGE_ADMISSION` is a repository administrator-controlled JSON
variable, not a secret and not consuming-project data. A build executor may
issue it only after independently verifying account controls and a conservative
shared-storage balance. The guard cannot establish the truth of an externally
supplied storage measurement; an unsupported receipt must never be issued.

The exact required fields are defined by `RECORD_FIELDS` in
`build/release/storage_guard.py`. They bind schema, repository/owner, exact
tooling revision, evidence digest, issue/expiry epochs, included capacity,
occupied upper bound and other reserved bytes, zero-spend ceiling and explicit
Stop usage/headroom proof booleans. Numbers must be integers, not JSON booleans.
Unknown/duplicate fields, stale/future records, absent proof and conflicting
identity fail. Record lifetime and maximum age are one hour. Freshness is
checked again after hashing. Capacity cannot exceed a conservative 500,000,000
bytes, and proven headroom must accommodate the whole 256 MiB reservation.

Account observations, unrelated account object identifiers and credentials
remain outside the product source. Only an admitted account-level capacity
record may cross into the build. The current source does not create a record,
change account settings, grant credential scope, or contact a provider.

Budgets with Stop usage must remain USD 0. A rounded display or a partial package
inventory cannot justify setting `storage_headroom_verified=true`. Concurrent
consumption and the applicable billed storage model must be covered by the
external admission evidence. An absent record stops the transfer even if both
builds otherwise pass.

## Validation and preservation

Run `python -B build/release/test_storage_guard.py` for offline synthetic tests.
The Go workflow-preservation regression removes exactly one byte-equal guard
step, requires it immediately before upload, then compares the entire remaining
workflow with its historical identity-normalized predecessor. Old workflows,
schemas, product bytes and signing gate are preserved. This guard and its tests
do not establish Linux/Docker build success, native Windows execution or
PSCAN-06 acceptance.
