# Gitleaks Security Boundary

## Status

No accepted pass-capable Gitleaks path exists. PSCAN-03 native patch mode and
correction C1 remain immutable rejected history. PASS-OUTCOME-SPEC-001 controls
the future boundary and PSCAN-10 is proposed, unselected and unauthorized.

## Mandatory preparation controls

- Bind exact base, head, merge base, ordered range, every parent edge and the
  complete head tree.
- Admit only verified regular Git blob objects and required manifest/artifact
  inputs with exact identity, length and SHA-256.
- Disable inherited Git config, hooks, attributes, filters, pagers, prompts,
  external diffs, text conversion and interactive helpers.
- Classify raw bytes before any framing, extraction, decoding or projection.
- Reject unsafe paths, modes, types, links, ambiguity, missing/unreadable bytes,
  duplicates, conflicts and resource-limit breaches.
- Record a private admission-ledger row for every required input and byte-
  conserving transformation.

Candidate names and bytes are argument elements and data only, never executable
syntax. All inputs are read-only and transient.

## Mandatory detector controls

The exact Gitleaks source, binary, version, config, rule pack, ignore file,
command, output format and decoder are digest-bound. Before a pass, independent
evidence proves the exact buffer, fragment, overlap, streaming, archive,
decoding, path, type, size and report behavior of that pinned set. Every
blocking rule has a finite proved maximum span or a complete streaming proof.

Every original byte is bound to actual detector inspection. A prefix marker,
file-open event, aggregate count, exit code or clean report is insufficient.
Per-row evidence must be complete and unique. Any skip, timeout, overflow,
unexpected stderr, malformed output, missing/duplicate proof, unbound result or
redaction uncertainty is non-pass.

Engine output remains bounded in private memory or memory-backed temporary
storage, is decoded strictly, cleared and never emitted. Only the redaction
firewall may produce an external content-free outcome.

## Mandatory proof isolation

Named PR/release profiles bind object, byte, file, fragment, report, time,
memory and process limits. At/below/above cases and separate text, binary,
deleted-history, head-tree, merge, rename, span-boundary, raw-archive-class and
maximum-size fixtures are required. Aggregate failure does not prove a class.

## Preserved boundaries

Gitleaks remains primary; the Gitleaks GitHub Action is absent. Scanner
execution remains offline and credential-free. Native Git mode is not
inherently authoritative, but exact-object projection is not inherently
authoritative either: only complete outcome proof can support pass. TruffleHog
remains absent until accepted material-gap evidence and separate technical plus
AGPL owner approval.
