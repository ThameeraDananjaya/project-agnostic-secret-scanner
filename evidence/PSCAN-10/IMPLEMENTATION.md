# PSCAN-10 Implementation Record

## Claim

PSCAN-10 alone was claimed on 2026-09-01 from exact clean activation commit
`9053b37d799b19e3d98aeca0ae853296971ab09d`, after the complete reading map and
live primary-source preflight were completed. No successor was selected.

## Pre-claim bounded-design proof

- Selected detector build: Gitleaks `8.30.1` exact source commit
  `83d9cd684c87d95d656c1458ef04895a7f1cbd8e`, standard-library regex variant,
  Go `1.27.0`, `CGO_ENABLED=0`, Windows/Linux amd64.
- Two fresh builds were byte-identical per platform: Windows amd64 SHA-256
  `b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178`;
  Linux amd64 SHA-256
  `657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586`.
- Exact version canary returned `8.30.1`; an invalid-by-construction synthetic
  finding returned the configured finding exit `11`; a clean input returned
  exit `0` and an empty JSON array.
- At claim, all 223 admitted regexes parsed under Go's selected standard-
  library regex engine with zero unbounded expressions and conservative maximum
  UTF-8 byte span 4,020. The final pack adds one invalid-by-construction exact-
  maximum-span fixture rule; all final 224 rules retain the same proved bound.
- The preparation design uses detector files smaller than the pinned 100,000
  byte base buffer and overlap of maximum span minus one. Every chunk contains
  one deterministic coverage witness and a byte-offset mapping. Therefore no
  product match of at most 4,020 bytes can cross every overlapping chunk, and
  each chunk is inspected as a single detector fragment.
- Raw classification is authoritative on the exact Git-object bytes before
  preparation. Recognized archive/compression/container signatures,
  conflicting signatures and malformed recognized inputs are explicit
  non-pass; admitted text and binary bytes are unchanged in the chunk payload.

This record is a claim and design gate, not acceptance. Final digest bindings,
tests, independent review, limitations and closeout follow in separate evidence
within this task.
