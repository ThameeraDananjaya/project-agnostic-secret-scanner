# Gitleaks Adapter Boundary

## Status

PSCAN-03 correction C1 candidate `6cd22a3` is independently rejected and is not
an authoritative pass-capable path. The historical `GITLEAKS-GAP-001` record
remains immutable and native Gitleaks patch mode remains prohibited from
producing pass.

## Corrected Git and tree boundary

```text
exact base/head/merge-base and ordered in-range commits
  -> raw parent-edge enumeration with renames disabled
  -> admit both regular-file pre-image and post-image blob OIDs
  -> enumerate every regular blob in the exact head tree
  -> reject unsafe paths, modes, types, conflicts and resource limits
  -> materialize byte-exact private blobs
  -> materialize path-preserving framed copies
       PSCAN_GITLEAKS_PROJECTION_V1
       deterministic per-blob coverage marker
       exact original blob bytes
  -> verify per-file and aggregate SHA-256 bindings
  -> one pinned Gitleaks invocation with product + coverage rules
  -> require exactly one coverage finding for every admitted projection
  -> classify product findings privately and return a content-free reason
```

The adapter binds engine `gitleaks`, version `8.30.1`, adapter `1.1.0`, output
`json-v8.30.1`, finding exit `11`, exact binary SHA-256
`c79361874b71d1b8a366773cc3cee1ade9159b1b0500e3456835fff065c7555a`,
product-rule SHA-256
`2a9e75e17a09a1e3b06c6c232395cef85eb61448ef0298611162825f7da6c044`,
and empty-ignore SHA-256
`3cdd737e34cb301ffc19861f903cf59441ff57453ddf70bd5b8c1bd63389aebe`.
It does not use the Gitleaks Action. Any binding mismatch stops before source is
scanned.

Git preparation accepts only absolute local paths and creates a non-local bare
clone with an empty template. The child environment disables inherited system
and global Git configuration, hooks, attributes, pagers, prompts, external
diffs and interactive helpers. Exact commit OIDs, ancestry, merge base, ordered
history, every parent edge and the tracked head tree are proven before a
projection is admitted. Candidate values are argument elements and data bytes,
never shell syntax. Native `ScanGitRange` is a fail-closed compatibility surface
that always returns `INDETERMINATE_INCOMPLETE_COVERAGE`.

Tracked-source manifests describe a sorted, unique and complete set of regular
files. Links, special files, traversal, extra files, missing files, size changes
and digest changes fail validation.

## Deferred boundary

The per-blob marker proves that Gitleaks opened a framed file, but it does not
prove that all internal fragments were inspected. Gitleaks v8.30.1 reads a
100,000-byte buffer and at most 25,000 additional bytes without overlap. A
synthetic product match split at byte 125,000 produced only the coverage
finding. Because the pinned product rules include unbounded whole-match spans,
no finite overlap is currently a complete proof. Raw archive classification is
also required before framing can be authoritative. Therefore this candidate
must not produce an accepted scanner path; see `evidence/PSCAN-03/REVIEW-C1.md`.
It also uses Gitleaks directory mode for projected history, conflicting with
PASS-SPEC-001 section 7.2's explicit Git-mode requirement.

The scanner CLI is intentionally unchanged because `cmd/scanner-runner/**` is
outside PSCAN-03's allowed paths. No public product path can claim the adapter
is integrated. Archive/OCI normalization, final PSCAN-04 redaction/cleanup
packaging, policy, allowlist, receipt, workflow and release remain outside this
task.
