# Gitleaks Primary Coverage Map

## Current status

PSCAN-03 and correction C1 remain rejected historical evidence. PSCAN-10
adapter `2.0.0` provides pass-capable primary coverage only for the exact Git
object classes below and only when every ledger/chunk/profile proof completes.
Artifact/OCI normalization and every later release boundary remain deferred.

| Required class | PSCAN-10 proof | Status |
|---|---|---|
| Ordinary tracked head files | Exact head-tree OID, raw object re-hash, unchanged-head finding fixture | Covered |
| Text additions/deletions | Exact parent-edge pre/post objects plus isolated text and deleted-history fixtures | Covered |
| Binary head/history | Pre-transform binary class plus isolated head and deleted-binary detector runs | Covered |
| Merge and rename | Every parent edge and delete/add identity; isolated detector fixtures | Covered |
| Out-of-range objects | Clean/out-of-range fixture proves absence from range/head ledger | Covered |
| Fragment boundaries | Single-fragment chunks, 4,019 overlap, first/last/internal/no-whitespace and exact 4,020-byte cases | Covered |
| File/count/total maximum | Named profiles; every bound below/at/above; actual 512 MiB detector run | Covered |
| Archive/compression/container | Exact pinned-family pre-transform classifier; supported families and ambiguity are non-pass | Fail-closed, normalization deferred |
| Hostile Git behavior | Fixed no-shell commands, sealed environment, no hooks/attributes/diffs/textconv/prompts | Covered |
| Rule maximum span | 224 parsed rules, zero unbounded, conservative maximum 4,020 bytes | Covered |
| Admission-to-inspection | One ledger row per object, complete overlap mapping and one same-run witness per single-fragment chunk | Covered |
| Artifacts and OCI layers | Raw classification only; no normalizer in PSCAN-10 | Deferred non-pass |

## Decision boundary

No material Gitleaks gap was established for PSCAN-10's admitted exact-Git
classes. Deferred artifact/OCI normalization is planned product scope, not
evidence of primary-engine inability. PSCAN-08 remains inactive and ineligible;
TruffleHog remains absent. No successor is selected or activated.
