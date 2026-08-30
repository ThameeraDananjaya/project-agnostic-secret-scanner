# Gitleaks Adapter Boundary

## Status

PSCAN-03 draft implementation; not accepted. Material gap
`GITLEAKS-GAP-001` blocks completion.

## Realized draft boundary

```text
exact regular executable/config + SHA-256
  -> private runtime version probe (must be 8.30.1)
  -> fixed argument array (no shell)
  -> forced JSON, full engine redaction, no banner/color, error-only logs
  -> bounded private stdout/stderr
  -> strict exit/schema/redaction consistency classification
  -> content-free reason code only
```

The adapter binds engine `gitleaks`, version `8.30.1`, adapter `1.0.0`, output
`json-v8.30.1`, finding exit `11` and the exact default config digest. It does
not use the Gitleaks Action. The executable and config are re-opened as regular
files and hashed before launch. A wrong digest or runtime version fails before
candidate scanning.

Git preparation accepts only absolute local paths and creates a non-local bare
clone with an empty template. The child environment disables inherited system
and global Git configuration, hooks, attributes, pagers, prompts, external
diffs and interactive helpers. Exact commit OIDs, ancestry, merge base, ordered
history digest and tracked-tree digest are proven before an engine range is
formed. Candidate values are argument elements, never shell syntax.

Tracked-source manifests describe a sorted, unique and complete set of regular
files. Links, special files, traversal, extra files, missing files, size changes
and digest changes fail validation.

## Non-realized boundary

The existing scanner CLI is intentionally unchanged because
`cmd/scanner-runner/**` is outside PSCAN-03's allowed paths. No public product
path can claim the draft adapter is integrated. Archive/OCI normalization,
final redaction firewall, offline execution, policy, allowlist, receipt,
workflow and release remain outside this task.
