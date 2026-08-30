# Gitleaks v8.30.1 Coverage Map

| Required class | Corrected primary-engine observation | PSCAN-03 disposition |
|---|---|---|
| Ordinary tracked head files | Exact `ls-tree -r -z --full-tree` blob enumeration and framed Gitleaks directory projection | Every regular blob is bound or the plan is non-pass |
| Text additions/deletions in exact Git range | Raw parent-edge OIDs admit both pre-image and post-image blobs | Deleted text canary detected; exact base/head/merge-base/range remain bound |
| Binary-classified blobs present at head | Binary bytes are preserved after a deterministic text frame | NUL-containing blob reaches the same Gitleaks invocation and is detected |
| Binary-classified blob added/deleted only in history | Exact blob-object projection bypasses incomplete textual patches | Deleted NUL-containing history canary detected without fallback |
| Merge and rename history | Every in-range commit and every merge-parent edge is enumerated; renames are represented as delete/add | Deterministic parent/path-preserving projections proved |
| Finding outside declared range and absent at head | Neither an in-range edge nor the head tree admits the object | Clean exact-range test passes without misattribution |
| Gitleaks path/type skip | Same-run coverage rule lacks the required per-file marker | `INDETERMINATE_INCOMPLETE_COVERAGE`; never pass |
| File/count/total above declared maximum | Rejected during object-size/count/aggregate preflight | Resource-limit non-pass before Gitleaks launch |
| Nested archives and container layers | Adapter forces archive depth zero | PSCAN-04 only; unavailable in PSCAN-03 |
| Git metadata | Only trusted fixed Git commands enumerate identities and blobs | Candidate metadata is not executed or emitted |
| Symlink, gitlink, special, unsafe, unreadable, conflicting or unbound source | Admission and integrity gates reject before or during projection verification | Explicit non-pass |

`MATERIAL-GAP-001.md` remains the immutable reproduction for ordinary native
patch input. Correction C1 does not use that incomplete input form and does not
select or authorize a fallback engine. PSCAN-08 remains inactive.
