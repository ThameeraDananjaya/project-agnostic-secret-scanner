# Gitleaks Primary Coverage Map

## Current status

No input class is currently accepted as pass-capable. PSCAN-03 and correction
C1 remain rejected historical evidence. This map states the proof required by
PASS-OUTCOME-SPEC-001 for proposed, unselected PSCAN-10; it does not claim
implementation or execute a scanner.

| Required class | Immutable historical fact | Required successor proof |
|---|---|---|
| Ordinary tracked head files | C1 enumerated head blobs but did not prove later detector fragments | Exact head-tree/object/byte ledger plus per-byte pinned-detector inspection |
| Text additions/deletions in exact range | Parent-edge OIDs were exact; a match split at byte 125,000 escaped | Exact pre/post admission plus first/last/internal/no-whitespace span-edge proof |
| Binary-classified head blobs | Bytes were framed; aggregate fail did not isolate binary detection | Raw binary classification before framing and binary-only clean/finding evidence |
| Binary-classified added/deleted history | Native patch omitted bytes; C1 projected objects without isolated proof | Exact object bytes plus deleted-binary-only inspection and detection proof |
| Merge and rename history | C1 deterministically enumerated parent edges/delete-add identities | Fresh exact range/edge/object proof under the accepted successor implementation |
| Out-of-range and absent-at-head object | C1 did not admit it | Fresh negative proof that it is neither admitted nor attributed while range stays bound |
| Detector path/type/fragment behavior | Prefix marker could pass while later fragment missed | Exact pinned skip/span/stream behavior and per-byte inspection; marker alone forbidden |
| File/count/total maximum | C1 preflighted some limits; overlap experiment introduced unbound limits | Named PR/release profiles with all resource bounds tested below/at/above |
| Archives/compression/container signatures | Framing could mask raw identity | Complete raw pre-transformation classification, ambiguity/polyglot handling and non-pass |
| Git metadata and hostile configuration | Fixed commands disabled many hostile features | Fresh proof for config/hooks/attributes/filters/pagers/diffs/textconv/prompts |
| Required detector input form | Native Git mode was incomplete; C1 directory mode conflicted with predecessor | Input form is mechanism-neutral only after complete outcome proof |
| Unsafe/special/unreadable/conflicting source | C1 rejected known unsafe classes | Fresh exhaustive admission rejection and class-isolated evidence |
| Rule maximum span | Pinned rules included unbounded whole matches | Finite proved span for every blocking rule or complete streaming proof; otherwise non-pass |
| Admission-to-inspection binding | Per-file prefix marker was not enough | One-to-one ledger rows binding every required original byte to actual inspection |

## Decision boundary

PSCAN-10 must fail closed unless every claimed class above is independently
proved under exact pins and declared profiles. Any remaining material required-
class gap is recorded as evidence and stops first-consumer readiness. It does
not select, activate or approve PSCAN-08. TruffleHog remains absent pending a
separate technical and AGPL owner decision.
