# Gitleaks Security Boundary

PSCAN-03 correction C1 launches an exact absolute executable through Go's
argument-array API. It never invokes a shell. Executable, product-config and
empty-ignore digests, runtime version, exact lowercase Git OIDs and local
absolute paths are validated before scanning. Engine output is bounded in
private memory, supplied only to a strict decoder, cleared, and never included
in the returned result.

Git patches are not trusted as complete input. The scanner enumerates every
exact in-range commit and each parent edge with raw, no-rename, no-external-diff
and no-textconv metadata. Both deleted pre-image and added/modified post-image
blob objects are admitted, and the exact head tree is enumerated independently.
Only regular blob modes and portable, non-conflicting paths are supported.
Unknown modes, links, gitlinks, unsafe names, missing objects, unreadable bytes,
count/byte limits, duplicate paths and case collisions fail closed.

Each blob has a byte-exact private copy and a deterministic framed scan copy
that preserves the original path suffix and every original byte. The pinned
product config contains a private coverage rule. Coverage and product rules run
in one Gitleaks v8.30.1 invocation, so any path/type/size skip removes the
required per-file marker and yields `INDETERMINATE_INCOMPLETE_COVERAGE`.

Success requires empty stderr, valid JSON, full engine redaction, exactly one
coverage finding per admitted projection, and no product finding. A product
finding returns fail. Timeout, capture overflow, unexpected stderr, malformed
output, nonstandard exit, missing or duplicate coverage, an unbound result path
or an unredacted candidate is explicit non-pass.

The historical native-patch defect remains valid evidence in
`evidence/PSCAN-03/MATERIAL-GAP-001.md`. Native Gitleaks Git input is not the
corrected authoritative path and cannot return pass. TruffleHog remains absent.
