# Gitleaks Adapter Boundary

## Current authority

PASS-OUTCOME-SPEC-001 controls this target. No accepted pass-capable Gitleaks
adapter exists. PSCAN-03 native patch input and correction C1 remain rejected,
non-authoritative historical candidates. Their exact evidence remains immutable
in `evidence/PSCAN-03`.

Gitleaks remains the primary detector and the separately licensed Gitleaks
GitHub Action remains prohibited. PSCAN-10 implements adapter `2.0.0` from its
exact activation boundary; no successor is selected.

## Historical defect boundary

Pinned Gitleaks v8.30.1 native Git mode used a textual `git log -p` stream that
omitted deleted binary blob bytes. Correction C1 enumerated exact pre/post/head
objects and framed them for directory mode, but independent review proved a
false pass at the detector's 125,000-byte fragment boundary. A prefix marker
proved only file opening, not later fragment inspection. The rule pack contained
unbounded whole-match rules, framing masked raw archive identity, resource
limits were unbound and aggregate evidence did not isolate binary detection.
Neither path is authoritative.

## Target exact-object and inspection boundary

```text
exact base, head, merge base, ordered range and every parent edge
  -> enumerate required pre-image/post-image blob OIDs and complete head tree
  -> verify exact object type, mode, size, OID and bytes
  -> classify raw bytes before framing, extraction or transformation
  -> create private admission-ledger row for every object and required artifact
  -> apply only byte-conserving, versioned preparation with an exact mapping
  -> invoke one pinned Gitleaks binary/config/rule/ignore/output contract
  -> bind every original byte to actual detector stream/span inspection
  -> require complete per-row proof and class-isolated adversarial acceptance
  -> privately classify product findings and emit one content-free outcome
```

The implemented v2 preparation writes each exact raw blob to a private ledger,
classifies it with every registered `github.com/mholt/archives` v0.1.2 family,
and admits only ordinary text or binary. It then writes unchanged payload bytes
into 90,000-byte overlapping chunks. Adjacent chunks overlap by 4,019 bytes;
each complete detector file is below the pinned 100,000-byte Gitleaks base
buffer and contains a unique same-invocation inspection witness. Original path
suffixes are retained for path-bound rules.

Trusted fixed Git commands disable inherited configuration, hooks, attributes,
filters, pagers, prompts, external diffs, text conversion and interactive
helpers. Missing parents, unexpected types or modes, links, gitlinks, unsafe or
ambiguous paths, duplicate identities, unreadable objects, mismatched bytes,
unclassified input, unsupported raw class or resource-limit breach is non-pass.

## Detector proof requirements

The exact source, binary, version, config, rule pack, ignore file, command,
output format and decoder are digest-bound. A pass-capable design proves the
pinned detector's buffer, peek, fragment, overlap, streaming, archive, decoding,
path/type/size and report behavior. Every blocking rule has a finite proved
maximum required span, or the detector path has a complete streaming proof.
Unbounded matches cannot be justified by finite overlap. The admitted rule pack
contains 224 rules, parses under the selected Go standard-library regex engine,
has no unbounded expression and has conservative maximum UTF-8 byte span 4,020.

Inspection evidence is one-to-one with admission-ledger rows and proves the
first, last and every internal byte under the declared profile. A marker,
file-open event, output count, aggregate failure or clean exit is insufficient.
Timeout, capture overflow, unexpected output, malformed decoding, missing or
duplicate proof, detector skip, unbound result or redaction uncertainty is
explicit non-pass.

## Raw classification and bounded profiles

Classification occurs on exact raw bytes and covers every archive,
compression, container, text, binary, ambiguous and polyglot family recognized
by the pinned stack. Framing cannot hide or change that authority. Unsupported,
encrypted, malformed or conflicting classes fail before a projection can pass.

Declared PR and release profiles bind 512 MiB per blob, 100,000 objects,
2/10 GiB total bytes, 125,000 detector files, 256 MiB private report capture,
15/60 minute time, 2/4 GiB Go memory limit and one detector process. At/below/above boundary tests and separate text,
binary, deleted-history, head-tree, merge, rename, archive-class and span-edge
fixtures are mandatory. Combined failures do not prove an individual class.

## Deferred boundaries

PSCAN-10 may implement only its future activated boundary. Artifact/OCI
normalization, final redaction/cleanup/offline packaging, project policy,
allowlists, receipts, workflows and release remain later tasks. TruffleHog
remains absent and inactive pending accepted gap evidence and separate technical
plus AGPL owner approval.
