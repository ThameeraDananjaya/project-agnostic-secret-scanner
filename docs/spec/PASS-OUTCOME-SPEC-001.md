# Project-Agnostic Secret Scanner Outcome Contract

- Document ID: `PASS-OUTCOME-SPEC-001`
- Version: `1.0`
- Decision date: `2026-09-01`
- Status: controlling product contract through accepted PSCAN-09 closeout
- Historical predecessor: `PASS-SPEC-001` version 1.0

## 1. Authority transition

PASS-OUTCOME-SPEC-001 superseded PASS-SPEC-001 when PSCAN-09 was independently
accepted and closed. The reason for replacement is narrow:
PASS-SPEC-001 section 7.2 prescribed Gitleaks Git mode for exact history/range
scanning, while accepted historical gap evidence within rejected PSCAN-03
proved that the pinned Git-mode patch stream omits required deleted binary
bytes. An outcome contract must allow an
implementation to enumerate exact Git objects and use deterministic detector
projections when, and only when, complete detector inspection is proved.

PASS-SPEC-001, its source record, the owner-provided external source, and all
PSCAN-01 through PSCAN-03 records remain immutable historical evidence. The
PSCAN-03 implementation and correction C1 remain rejected and non-authoritative.
No sentence in this contract accepts, rehabilitates or incorporates that
candidate.

All predecessor capabilities, safety constraints, product boundaries,
non-goals, schema boundaries, owner gates and definition-of-done requirements
continue unless an explicit transition disposition says that a mechanism is
reframed. A mechanism may be replaced only by a proof obligation that is at
least as strong. Missing or ambiguous disposition fails closed.

## 2. Mission and product boundary

The product is one independently versioned, project-agnostic, local secret-
scanning package. It uses a thin MIT-licensed Go runner and Gitleaks CLI as its
primary detector. Authoritative execution is offline, credential-free and
candidate-as-data. Findings are transient and fully redacted. The public result
is a content-free outcome bound to exact inputs and proof evidence.

The product is not a hosted scanner, source host, credential validator,
incident-response platform, durable finding database, project policy authority,
receipt issuer, evidence custodian, merge controller, deployment controller,
compliance authority or go-live authority.

Consuming projects remain callers and own their policies, allowlists, receipt
trust, revocations, evidence custody, retention, deployment gates, credentials,
incident response and business decisions. The scanner repository contains no
consuming-project source, policy, allowlist, receipt, key, finding, credential,
customer data, repository identity, statistics or project-specific behavior.

### Standing owner decisions retained

- Product and intended future remote name:
  `project-agnostic-secret-scanner`.
- Default branch: `main`.
- Original runner and repository licence: MIT.
- Intended public repository owner: `ThameeraDananjaya`.
- Intended final release: signed public `v1.0.0`.
- Spending ceiling: zero.
- Intended signing design: Cosign keyless signing from the exact separately
  approved GitHub release workflow.
- Release maintainer and final human approver: `ThameeraDananjaya`.
- Final runbook: plain-English Markdown with clear Mermaid diagrams and complete
  input/output parameter and payload documentation.

These directions do not approve their action-time remote, settings, signing,
publication, credential or spending gates.

## 3. Non-negotiable invariants

1. The product remains independent of every consuming project.
2. GitHub native secret scanning and hosted scanning services are not product
   dependencies.
3. Gitleaks CLI remains the pinned primary detector.
4. TruffleHog remains absent and inactive unless accepted evidence proves a
   material required-class gap and the owner separately approves both the
   technical need and AGPL obligations.
5. Candidate material is data only and is never executed.
6. Scanner subprocesses use argument arrays, never interpolated shell commands.
7. Authoritative scanner execution has no network access or credentials.
8. Raw or partially raw findings never cross the private process boundary.
9. Unsupported, skipped, incomplete, stale, conflicting, unbound or untrusted
   coverage never produces pass.
10. Outcomes are content-free and cannot approve merge, deployment, production,
    legal compliance or go-live.
11. Projects own policy, allowlist, receipt, revocation, custody, deployment and
    retention authority.
12. Public scanner-release trust and project receipt trust remain separate.
13. Releases use exact versions and digests; `latest` and auto-update are
    prohibited.
14. All CAP-1 through CAP-20 and AT-01 through AT-31 requirements must pass
    before first-consumer readiness.
15. No successor task activates automatically.
16. Every admitted Git object and every byte in each admitted object must be
    bound to evidence of actual detector inspection.
17. Raw input classification occurs before framing, transformation, extraction
    or detector projection and cannot be masked by those operations.
18. Detector stream, buffer, fragment, archive, decoding and maximum-match-span
    behavior is proved for the exact pinned detector and rule pack before pass.
19. Resource bounds are declared, versioned and proved against named PR and
    release profiles; exceeding any bound is non-pass and never skip-and-pass.
20. Acceptance evidence isolates every claimed input class and adversarial
    boundary; an aggregate failure cannot prove an individual class.

## 4. Outcome-proof model

### 4.1 Exact source identity

For PR mode, the request and outcome bind the exact base, head, merge base,
ordered in-range commits, every required parent edge, admitted pre-image and
post-image blob OID, and the complete regular-file head tree. For release mode,
the binding additionally covers the previous accepted release boundary, the
complete required history, frozen tracked tree, build-context manifest, every
deployment artifact and every supported nested member.

Enumeration uses trusted fixed Git commands with inherited configuration,
hooks, attributes, filters, pagers, external diffs, text conversion, prompts and
interactive helpers disabled. Object type, mode, size and object bytes are
verified against the exact OID. Links, gitlinks, special modes, unsafe paths,
unreadable objects, ambiguous identities, missing parents, unexpected objects
or inconsistent ranges are non-pass.

### 4.2 Admission ledger and byte conservation

Before detector invocation, the runner creates a private deterministic admission
ledger. Each row binds input class, source object or artifact identity, byte
length, SHA-256, raw classification, transformation version, projection digest,
detector invocation and inspection evidence. The ledger proves a bijection
between required inputs and admitted detector inputs. Duplicate, missing,
unbound or extra rows are non-pass.

Every transformation has a specified reversible or byte-conserving mapping.
For projections, proof must show where every original byte appears and which
detector stream/span inspected it. A prefix marker, file-open event, aggregate
count, exit code or clean report is never sufficient inspection evidence.

### 4.3 Raw classification before transformation

Raw bytes are classified before any prefix, frame, path-preserving wrapper,
decompression, decoding or normalization is applied. Classification covers all
exact input families recognized by the pinned detector and normalizer,
including archive/compression/container signatures and ambiguous/polyglot
cases. Unsupported, conflicting, encrypted, malformed or unclassified input is
explicitly non-pass. A transformation may not change the authority of the raw
classification or hide a supported/unsupported class.

### 4.4 Proved detector inspection

The exact detector source, binary, configuration, rule pack, ignore file,
command, version and output decoder are digest-bound. Acceptance must establish:

- exact buffer, peek, fragment, overlap, streaming and archive behavior;
- a finite maximum required match span for every admitted blocking rule, or a
  streaming proof that makes a finite span unnecessary;
- complete byte coverage at fragment boundaries, including the first and last
  possible byte and adversarial no-whitespace inputs;
- no size, type, path, archive, encoding or report limit can silently omit an
  admitted byte;
- one-to-one inspection evidence for each admission-ledger row; and
- a non-pass result for any detector behavior that cannot be established from
  the pinned implementation.

Unbounded whole-match rules cannot support a finite-overlap projection. They
must be replaced by bounded semantics with equivalent or stronger detection,
or by a proved streaming path, before that projection can pass.

### 4.5 Bounded profiles and class-isolated proof

The contract retains the predecessor required product profiles: PR timeout 15
minutes, release timeout 60 minutes, archive depth 5, 100,000 entries, 2 GiB PR
expanded bytes, 10 GiB release expanded bytes, 512 MiB regular file and 1,000:1
compression ratio. A consuming project may impose stricter policy limits. A
product-profile reduction requires an explicit owner-authorized successor
contract and cannot be inferred by an implementation; higher limits likewise
require a versioned reviewed profile and complete proof.

Tests name the exact profile and input class they prove. Text, binary, archive,
container, head-tree, deleted-history, merge, rename, boundary-span, maximum-size
and over-limit cases each require isolated clean/finding/non-pass evidence.
Combined evidence may supplement but never replace the isolated proof.

## 5. Capability contract and disposition

| ID | Disposition under PASS-OUTCOME-SPEC-001 |
|---|---|
| CAP-1 | Retained: two structurally different projects use one unmodified release with isolated project-owned inputs and state. |
| CAP-2 | Reframed mechanism only: exact PR bindings and complete required bytes remain; Git mode is not privileged over proved exact-object inspection. |
| CAP-3 | Retained and strengthened by admission-ledger proof for history, tree, build context, artifacts, layers and archives. |
| CAP-4 | Retained: complete authoritative acceptance runs without DNS, TCP, providers, updates, telemetry, credentials or source transmission. |
| CAP-5 | Reframed mechanism only: pinned Gitleaks remains primary, while an adapter may use an outcome-proved input form instead of native Git mode. |
| CAP-6 | Retained and owner-gated: fallback remains absent until accepted material-gap evidence and separate technical plus AGPL approval. |
| CAP-7 | Retained and strengthened by raw pre-transformation classification and byte-bound normalization. |
| CAP-8 | Retained: no raw value or forbidden finding detail crosses any external or residual surface. |
| CAP-9 | Retained: policy remains project-owned and isolated behind a strict projection boundary. |
| CAP-10 | Retained: narrow approved exceptions expire within 30 days and genuine credentials remain non-allowlistable. |
| CAP-11 | Retained: identical validated inputs and proof bindings yield semantically identical outcomes on Windows and Linux. |
| CAP-12 | Retained and strengthened: any unproved object, byte, class, detector span or resource profile is non-pass. |
| CAP-13 | Retained: public JSON and exit expose only approved content-free fields and exact proof/binding digests. |
| CAP-14 | Retained: exact offline release, manifest, asset, identity, bundle, licence and SBOM verification. |
| CAP-15 | Retained: independent contract versions, unknown-major failure and controlled retirement overlap. |
| CAP-16 | Retained: projects can issue receipts from content-free bindings without any receipt key or promotion authority entering the scanner. |
| CAP-17 | Retained: append-only global/project revocations invalidate evidence without rewriting history. |
| CAP-18 | Retained: no project source, finding, policy, allowlist, cache, receipt, key or identity crosses projects. |
| CAP-19 | Retained: tested Windows amd64 and Linux amd64 assets with no unsupported architecture claim. |
| CAP-20 | Retained: no automatic update, constrained rollback, and enforced retirement/revocation. |

## 6. Public request, outcome and schema boundaries

The scanner continues to own `scan-request`, `scan-outcome`,
`scanner-release-manifest`, `global-scanner-revocation` and generic rule-pack
schemas. Projects continue to own policy, allowlist, receipt,
receipt-revocation, evidence-chain, custody and retention schemas and instances.
DEC-001 remains valid.

The request retains all predecessor bindings and additionally binds the
admission-ledger schema/digest, raw-classifier version/digest, projection or
stream preparation version/digest, declared resource-profile identifier and
inspection-proof format. Unknown majors, ambiguous semantics, missing required
proof fields, duplicate keys, mismatched digests or unsafe paths fail before
scanning.

The outcome remains exactly one content-free JSON object with the predecessor
states, reason codes and exit mapping. It may expose only digests/versions for
the new proof bindings, never object names, paths, finding counts, detector
names, excerpts or inspection detail.

## 7. Acceptance matrix disposition

| ID | Required successor result |
|---|---|
| AT-01 | Known-safe class-isolated fixtures pass only with complete inspection proof. |
| AT-02 | Synthetic findings fail without disclosure for every separately claimed detector/input class. |
| AT-03 | Added-then-deleted text and binary objects in the exact range independently fail after exact-object admission. |
| AT-04 | Out-of-range objects are not admitted or attributed; exact range/object/tree identities remain bound. |
| AT-05 | Every regular head-tree object and byte is in the admission ledger and inspected, or the result is non-pass. |
| AT-06 | Every generated/untracked build-context byte is manifest-bound and inspected, or non-pass. |
| AT-07 | Every deployment file, container layer and supported member is byte-bound and inspected, or non-pass. |
| AT-08 | Each supported nested archive class independently detects its synthetic finding. |
| AT-09 | Encrypted, malformed, unsafe, ambiguous or unsupported raw archives are indeterminate before masking transformation. |
| AT-10 | At/below/above each declared resource boundary is isolated; over-limit is indeterminate, never skip/pass. |
| AT-11 | Standard, error, timeout and crash surfaces contain no secret or forbidden metadata. |
| AT-12 | Malicious names and messages still yield one valid content-free outcome. |
| AT-13 | DNS/TCP/provider canaries prove deterministic network-disabled execution. |
| AT-14 | One-byte changes to runner, engine, rules, classifier, projection, ledger or proof formats fail before pass. |
| AT-15 | Wrong release identity, issuer, tag, ref or bundle is rejected. |
| AT-16 | Unknown-major, corrupt, duplicate-key or proof-incomplete requests are rejected/non-pass. |
| AT-17 | Invalid project-policy signature, digest, projection or precedence fails closed. |
| AT-18 | Expired, broadened or credential-bearing exceptions are rejected and fail. |
| AT-19 | Eligible transient states receive at most two clean-workspace retries, then terminal. |
| AT-20 | A deterministic synthetic finding is never retried. |
| AT-21 | Any required-engine finding fails; otherwise disagreement or incomplete proof is indeterminate. |
| AT-22 | Pass, fail, timeout, cancellation and crash leave no transient candidate/finding/proof material. |
| AT-23 | Receipt reuse requires every predecessor and new proof binding identical and within 30 days. |
| AT-24 | Any predecessor or new proof binding change rejects receipt reuse. |
| AT-25 | Scanner, key, receipt, policy or exception revocation rejects affected evidence. |
| AT-26 | Truncated, divergent, duplicate or conflicting evidence chains block progression. |
| AT-27 | Parallel, sequential and reused-host projects expose no cross-project state. |
| AT-28 | Hooks, attributes, filters, executables and candidate workflows remain data only. |
| AT-29 | Extra token rights, credentials or permissions are workflow-policy failures. |
| AT-30 | Windows amd64 and Linux amd64 give equivalent outcomes and inspection claims. |
| AT-31 | Small, medium and declared-maximum profiles complete within bound; unsupported capacity is non-pass. |

## 8. Non-goals

The predecessor non-goals are retained: no GitHub Secret Protection substitute,
hosted scanning, provider verification, credential storage/rotation/revocation,
incident response, merge/deployment/production/compliance/go-live authority,
shared project policy/state, finding database, initial GUI/server/dashboard/
notifications/SARIF, general vulnerability scanning, automatic updates,
automatic exceptions or production KMS/HSM selection.

Exact-object enumeration and deterministic projections are preparation
mechanisms, not new product capabilities. They do not authorize arbitrary Git
history rewriting, source export, durable object storage, detector replacement,
weaker rules, skipped classes or execution of candidate content.

## 9. Definition of done

First-consumer readiness requires every CAP-1 through CAP-20 and AT-01 through
AT-31 result on Windows amd64 and Linux amd64; exact mutual release bindings;
offline verification and execution; complete admission-ledger and detector-
inspection proof for every claimed input byte; no leaked finding material; no
unsupported input or skipped content producing pass; tested revocation and
retirement; a synthetic external project receipt flow without scanner receipt
keys; and independent confirmation that the scanner has no project authority.

The world-change signal remains two structurally unrelated fixture projects
using the same exact signed release, detecting synthetic findings in PR and
complete frozen-release inputs, emitting only content-free outcomes and blocking
progress whenever coverage or trust cannot be proved.

## 10. Owner gates

The following remain separately owner-controlled at action time: remote/repo
creation; push/publication/release; GitHub workflow, ruleset, environment,
branch, protection or release settings; any spending; release maintainer and
keyless workflow identity; any credential or signing/receipt key; TruffleHog
technical need and AGPL obligations; KMS/HSM or other provider; cross-project
reporting; consuming-project adoption; required checks, deployment, staging,
production or go-live; and the final legal/compliance retention schedule.

The mechanism-to-outcome decision authorizes none of those actions.

## 11. Bounded implementation sequence

PSCAN-10 is the proposed, unselected task for the primary-coverage redesign and
proof harness required by this contract. PSCAN-04 through PSCAN-07 remain
proposed and unselected; their prerequisites are reconciled to require accepted
PSCAN-10 primary coverage rather than accepted PSCAN-03. PSCAN-08 remains
inactive and separately technically and legally gated. No task is selected or
activated by this contract.
