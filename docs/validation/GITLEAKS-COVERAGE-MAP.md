# Gitleaks v8.30.1 Coverage Map

| Required class | Correction C1 observation | PSCAN-03 disposition |
|---|---|---|
| Ordinary tracked head files | Exact `ls-tree -r -z --full-tree` blob enumeration and framed Gitleaks directory projection | Enumeration is implemented; detector-fragment completeness is unproved, so candidate is not accepted |
| Text additions/deletions in exact Git range | Raw parent-edge OIDs admit both pre-image and post-image blobs | Exact range binding passed, but a match split at byte 125,000 escaped product detection |
| Binary-classified blobs present at head | Binary bytes are preserved after a deterministic text frame | Aggregate text-plus-binary fail did not isolate binary detection |
| Binary-classified blob added/deleted only in history | Exact blob-object projection bypasses incomplete textual patches | Exact bytes are projected, but independent binary-only Gitleaks proof is absent from the committed candidate |
| Merge and rename history | Every in-range commit and every merge-parent edge is enumerated; renames are represented as delete/add | Deterministic parent/path-preserving projections proved |
| Finding outside declared range and absent at head | Neither an in-range edge nor the head tree admits the object | Clean exact-range test passes without misattribution |
| Gitleaks path/type/fragment skip | Same-run coverage rule requires one prefix marker per file | Path skip is non-pass; later internal-fragment omission can still falsely pass |
| File/count/total above declared maximum | Rejected during object-size/count/aggregate preflight | Resource-limit non-pass before Gitleaks launch |
| Nested archives and container layers | Framing masks raw signatures while adapter forces archive depth zero | Unsupported in PSCAN-03 and must fail before framing; source-complete classification is absent |
| Git metadata | Only trusted fixed Git commands enumerate identities and blobs | Candidate metadata is not executed or emitted |
| Required engine mode | History blobs are passed to Gitleaks directory mode | Conflicts with PASS-SPEC-001 section 7.2 Git-mode requirement |
| Symlink, gitlink, special, unsafe, unreadable, conflicting or unbound source | Admission and integrity gates reject before or during projection verification | Explicit non-pass |

`MATERIAL-GAP-001.md` remains the immutable reproduction for ordinary native
patch input. Correction C1 does not use that input form, but its replacement is
also rejected. No fallback engine is selected or authorized. PSCAN-08 remains
inactive.
