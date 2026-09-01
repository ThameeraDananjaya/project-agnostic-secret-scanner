# PASS-SPEC-001 to PASS-OUTCOME-SPEC-001 Transition Traceability

## Authority and interpretation

This matrix records the complete PSCAN-09 disposition of the immutable
PASS-SPEC-001 contract. PASS-SPEC-001 remains historical audit evidence and is
not edited. PASS-OUTCOME-SPEC-001 becomes controlling only upon independently
accepted and committed PSCAN-09 closeout.

`Retained` preserves meaning. `Strengthened` preserves meaning and adds proof.
`Reframed` replaces only a prescribed mechanism with equal-or-stronger outcome
proof. `Deferred` preserves a requirement for its named future task and makes no
implementation claim. `Owner-gated` preserves a separate human decision. No row
is waived or silently omitted.

## Capability transition

| ID | Predecessor requirement | Disposition | Successor proof and owner task |
|---|---|---|---|
| CAP-1 | One unmodified release serves two different projects with independent inputs/state | Retained | Project-neutral contracts plus two-project isolation; PSCAN-02/05/07 |
| CAP-2 | PR scan binds exact base/head/merge-base/range/tree and release/rule/project inputs | Reframed mechanism; requirement retained | Exact Git-object/range/tree admission ledger and byte-inspection proof; PSCAN-10, then PSCAN-07 |
| CAP-3 | Frozen release covers history, tree, build context, deployment artifacts, layers and archives | Strengthened | Exact object/artifact identity, raw class and every-byte inspection; PSCAN-10/04/07 |
| CAP-4 | Authoritative acceptance is offline with DNS/TCP unavailable | Retained | Credential-free network-disabled canaries; PSCAN-04/07 |
| CAP-5 | Pinned Gitleaks primary orchestration with exact bindings | Reframed mechanism only | Gitleaks remains primary; native Git mode is not privileged over proved exact-object input; PSCAN-10/06/07 |
| CAP-6 | Disabled fallback only for accepted Gitleaks gap and explicit request | Retained and owner-gated | Accepted material-gap map plus technical and AGPL approval; PSCAN-08 only if eligible |
| CAP-7 | Supported directory/archive/container inputs scan or explicitly return non-pass | Strengthened | Raw classification before transformation plus byte-bound normalization; PSCAN-04/07 |
| CAP-8 | No secret or forbidden finding detail crosses any output/residual surface | Retained | Multi-surface redaction and cleanup canaries; PSCAN-04/07 |
| CAP-9 | Project-owned policy remains isolated | Retained | Strict content-addressed projection and two-project proof; PSCAN-05/07 |
| CAP-10 | Narrow approved expiring exceptions; credentials cannot be allowlisted | Retained | Scope/expiry/mutation/credential rejection; PSCAN-05/07 |
| CAP-11 | Identical validated input yields equivalent deterministic outcomes | Retained | Repeat and Windows/Linux semantic parity including proof bindings; PSCAN-02/06/07 |
| CAP-12 | Unsupported, skipped, incomplete or untrusted conditions never pass | Strengthened | Includes unbound objects/bytes/classes/spans/profiles; PSCAN-02 through PSCAN-07 plus PSCAN-10 |
| CAP-13 | Stable content-free JSON/exit contract | Retained | Closed schema and forbidden-field tests; new proof digests only; PSCAN-02/04/07 |
| CAP-14 | Offline release identity/manifest/assets/licences/SBOM proof | Retained | Mutation and identity negative matrix; PSCAN-06/07 |
| CAP-15 | Independently versioned contracts and controlled retirement | Retained | Minor/major/retirement compatibility tests; PSCAN-02/05/06/07 |
| CAP-16 | Project can issue receipt without scanner receipt key or promotion authority | Retained | Synthetic external issuer and key-absence proof; PSCAN-05/07 |
| CAP-17 | Append-only revocation invalidates evidence without rewriting it | Retained | Authorized-source/effective-time matrix; PSCAN-05/07 |
| CAP-18 | No source/finding/policy/allowlist/cache/receipt/key/identity crosses projects | Retained | Parallel, sequential and reused-host contamination tests; PSCAN-02/05/07 |
| CAP-19 | Tested Windows amd64 and Linux amd64 delivery | Retained | Cross-platform build, verification and logical parity; PSCAN-06/07 |
| CAP-20 | No auto-update; safe rollback, retirement and revocation | Retained | Moving-reference denial and lifecycle tests; PSCAN-06/07 |

## Acceptance-row transition

| ID | Disposition | Required successor evidence | Task |
|---|---|---|---|
| AT-01 | Strengthened | Safe class-isolated fixture passes only with complete inspection proof | PSCAN-10/07 |
| AT-02 | Strengthened | Each claimed class independently fails on its synthetic finding without disclosure | PSCAN-10/04/07 |
| AT-03 | Reframed mechanism | Deleted in-range text and binary objects independently reach actual detector inspection and fail | PSCAN-10/07 |
| AT-04 | Strengthened | Exact range/object/tree identity; out-of-range object neither admitted nor attributed | PSCAN-10/07 |
| AT-05 | Strengthened | Every regular head-tree object/byte admitted and inspected or non-pass | PSCAN-10/07 |
| AT-06 | Retained | Complete manifest-bound build-context inspection or non-pass | PSCAN-04/07 |
| AT-07 | Retained | Complete deployment file/layer/member inspection or non-pass | PSCAN-04/07 |
| AT-08 | Strengthened | Each supported nested archive class separately detects its fixture | PSCAN-04/07 |
| AT-09 | Strengthened | Raw encrypted/malformed/unsafe/ambiguous class rejected before masking transformation | PSCAN-10/04/07 |
| AT-10 | Strengthened | Isolated at/below/above tests for every named resource bound; never skip/pass | PSCAN-10/04/07 |
| AT-11 | Retained | No value or forbidden metadata on standard/crash/residual surfaces | PSCAN-04/07 |
| AT-12 | Retained | Malicious names/messages still yield one valid content-free outcome | PSCAN-02/04/07 |
| AT-13 | Retained | DNS/TCP/provider canaries cannot connect | PSCAN-04/07 |
| AT-14 | Strengthened | One-byte mutation includes classifier/projection/ledger/proof formats | PSCAN-10/06/07 |
| AT-15 | Retained | Wrong identity/issuer/tag/ref/bundle rejects release | PSCAN-06/07 |
| AT-16 | Strengthened | Unknown-major/corrupt/duplicate/proof-incomplete request rejects or is non-pass | PSCAN-02/05/10/07 |
| AT-17 | Retained | Invalid policy signature/digest/projection/precedence fails closed | PSCAN-05/07 |
| AT-18 | Retained | Expired/broadened/credential exception rejects and fails | PSCAN-05/07 |
| AT-19 | Retained | At most two clean-workspace retries for eligible transient state, then terminal | PSCAN-02/07 |
| AT-20 | Retained | Deterministic finding does not retry | PSCAN-02/07 |
| AT-21 | Strengthened | Any finding fails; disagreement or incomplete inspection proof is indeterminate | PSCAN-07; PSCAN-08 only if eligible |
| AT-22 | Retained | Pass/fail/timeout/cancel/crash leave no transient material | PSCAN-04/07 |
| AT-23 | Strengthened | Receipt reuse requires all predecessor and new proof bindings identical within 30 days | PSCAN-05/07 |
| AT-24 | Strengthened | Any predecessor or new proof-binding mutation rejects reuse | PSCAN-05/07 |
| AT-25 | Retained | Scanner/key/receipt/policy/exception marker rejects affected evidence | PSCAN-05/07 |
| AT-26 | Retained | Truncation/divergence/duplicate/conflict blocks progression | PSCAN-05/07 |
| AT-27 | Retained | Parallel/sequential/reused-host runs expose no cross-project state | PSCAN-02/05/07 |
| AT-28 | Retained | Candidate hooks/attributes/filters/executables/workflows remain data | PSCAN-10/06/07 |
| AT-29 | Retained | Extra token rights, credentials or permissions fail workflow policy | PSCAN-06/07 |
| AT-30 | Strengthened | Windows/Linux outcomes and inspection claims are equivalent | PSCAN-06/07 |
| AT-31 | Strengthened | Named small/medium/maximum profiles complete; unsupported capacity is non-pass | PSCAN-10/04/07 |

## Invariant disposition

| ID | Requirement | Disposition |
|---|---|---|
| INV-01 | Product independence | Retained unchanged |
| INV-02 | No GitHub-native/hosted scanner dependency | Retained unchanged |
| INV-03 | Gitleaks CLI is primary | Retained; only native Git-mode authority is reframed |
| INV-04 | TruffleHog absent until evidence plus separate approval | Retained unchanged |
| INV-05 | Candidate is data and never executed | Retained unchanged |
| INV-06 | Argument arrays; no shell interpolation | Retained unchanged |
| INV-07 | Authoritative execution has no network or credentials | Retained unchanged |
| INV-08 | Raw/partial findings remain private | Retained unchanged |
| INV-09 | Unsupported/incomplete/untrusted coverage is non-pass | Strengthened to objects, bytes, classes, spans and profiles |
| INV-10 | Outcome grants no merge/deploy/production/legal/go-live authority | Retained unchanged |
| INV-11 | Projects own policy/allowlist/receipt/revocation/custody/deployment/retention | Retained unchanged |
| INV-12 | Scanner release trust differs from project receipt trust | Retained unchanged |
| INV-13 | Exact versions/digests; no latest or auto-update | Retained unchanged |
| INV-14 | Complete capability and acceptance matrix before readiness | Retained unchanged |
| INV-15 | No automatic successor activation | Retained unchanged |
| INV-16 | Every admitted object and byte reaches actual detector inspection | New stricter proof |
| INV-17 | Raw classification precedes transformation | New stricter proof |
| INV-18 | Exact detector span/stream behavior is proved | New stricter proof |
| INV-19 | Resource profiles are named, bounded and tested | New stricter proof |
| INV-20 | Claimed classes have isolated adversarial evidence | New stricter proof |

## Component, schema and lifecycle boundaries

| ID | Predecessor boundary | Disposition |
|---|---|---|
| BND-01 | Thin deterministic MIT Go runner | Retained |
| BND-02 | Gitleaks adapter pinned, private, redacted and not the Action | Reframed only to permit outcome-proved exact-object input |
| BND-03 | TruffleHog compile/request disabled and undistributed until approval | Retained and owner-gated |
| BND-04 | Safe artifact normalizer with explicit limits | Retained; raw classification strengthened |
| BND-05 | Fixed safety/policy/allowlist precedence | Retained |
| BND-06 | Mandatory redaction firewall on all paths | Retained |
| BND-07 | Scanner owns five schema families; projects own policy/receipt families | Retained through DEC-001 |
| BND-08 | Synthetic fixture/acceptance harness, no credentials | Retained; class isolation strengthened |
| BND-09 | Validated request binds modes, assets, source, manifests and limits | Retained; adds ledger/classifier/preparation/profile/proof digests |
| BND-10 | PR and release exact-candidate coverage | Reframed from mandated native Git mode to exact-object outcome proof |
| BND-11 | Outcome state/retry/recovery history is immutable | Retained |
| BND-12 | Stable content-free reason/CLI/exit/message contract | Retained |
| BND-13 | Findings transient, redacted, non-durable and closed only by new scan | Retained |
| BND-14 | Project-owned policy/severity/allowlist and 30-day exception maximum | Retained |
| BND-15 | Independent major/minor compatibility and retirement overlap | Retained |
| BND-16 | Exact release intake, bundle, verification, rollback and retirement | Retained |
| BND-17 | Public release signing trust separated from project receipt trust | Retained |
| BND-18 | Receipt/revocation bindings and append-only evidence | Retained; new proof bindings join reuse/invalidation checks |
| BND-19 | Custody neutrality and cross-project isolation/reporting opt-in | Retained |
| BND-20 | Threat, CI, retention and performance controls | Retained; detector proof added to incomplete-coverage threat |
| BND-21 | Complete release documentation/runbook | Retained; documents new proof inputs and limits |

## Non-goal disposition

| ID | Non-goal | Disposition |
|---|---|---|
| NG-01 | GitHub Secret Protection replacement | Retained |
| NG-02 | Hosted scan/source upload/provider verification | Retained |
| NG-03 | Credential storage, validation, rotation, revocation or incident response | Retained |
| NG-04 | Merge, deployment, production, compliance or go-live authority | Retained |
| NG-05 | Shared project policy, allowlist, receipt, key, finding or cache | Retained |
| NG-06 | Durable per-finding database | Retained |
| NG-07 | Initial GUI, server, dashboard, notifications or SARIF | Retained |
| NG-08 | General vulnerability/dependency/malware/SAST/DAST/infrastructure scan | Retained |
| NG-09 | Automatic updates or direct unverified upstream consumption | Retained |
| NG-10 | Automatic exception approval | Retained |
| NG-11 | Production KMS/HSM selection | Retained and owner-gated |
| NG-12 | Exact-object preparation as a new authority or source-export feature | New clarification |

## Definition-of-done disposition

| ID | Requirement | Disposition |
|---|---|---|
| DOD-01 | All twenty capability successes | Retained |
| DOD-02 | Full 31-row acceptance on Windows/Linux amd64 | Retained |
| DOD-03 | No project-specific data or real credentials in release | Retained |
| DOD-04 | Source/runner/engine/rules/schemas/SBOM/licences/checksums/signatures bound | Retained; classifier/projection/ledger/proof bindings added |
| DOD-05 | Exact immutable release verifiable offline | Retained |
| DOD-06 | Authoritative tests network-disabled | Retained |
| DOD-07 | No raw or partial finding externally | Retained |
| DOD-08 | No unsupported or skipped input can pass | Strengthened with object/byte/span/class/profile proof |
| DOD-09 | Revocation and retirement tested | Retained |
| DOD-10 | Synthetic external adapter issues mock receipt without scanner key | Retained |
| DOD-11 | Independent review confirms no project/deployment/credential authority | Retained |
| DOD-12 | Two unrelated fixtures consume the same release and block on unproved coverage | Retained and strengthened by inspection proof |

## Owner-gate disposition

| ID | Reserved decision/action | Disposition |
|---|---|---|
| OG-01 | Create remote repository and final account/name | Retained; no action in PSCAN-09 |
| OG-02 | Push, public visibility, publication or release | Retained; no action |
| OG-03 | GitHub Actions/ruleset/branch/environment/release/protection settings | Retained; no action |
| OG-04 | Plan, cost, overage, storage or provider commitment | Retained zero-spend gate |
| OG-05 | Release maintainer and exact keyless workflow identity | Retained |
| OG-06 | TruffleHog technical need and AGPL obligations | Retained separately; no assessment |
| OG-07 | Signing credential or project receipt key | Retained; no handling |
| OG-08 | KMS/HSM or other paid custody provider | Retained |
| OG-09 | Cross-project reporting | Retained |
| OG-10 | Consuming-project adoption | Retained |
| OG-11 | Required checks, merge/deployment/staging/production/go-live | Retained |
| OG-12 | Final legal/compliance retention schedule | Retained |

## Historical and planning disposition

| ID | Requirement | Disposition |
|---|---|---|
| HIST-01 | PASS-SPEC-001 canonical and external bytes remain 977 lines and SHA-256 `78A1AD...B692D` | Immutable historical evidence; not edited |
| HIST-02 | PASS-SPEC-001 source record remains immutable | Preserved byte-for-byte |
| HIST-03 | PSCAN-01 and PSCAN-02 accepted records remain immutable | Preserved byte-for-byte; their implemented boundaries remain historical facts |
| HIST-04 | PSCAN-03 gap, correction and rejection remain immutable | Preserved; PSCAN-03 remains rejected and non-authoritative |
| HIST-05 | No accepted PSCAN-03 prerequisite is inferred | Removed from living successor planning; replaced by proposed PSCAN-10 |
| HIST-06 | Open version-sensitive facts require live primary-source preflight before pinning | Retained for the responsible future task |
| HIST-07 | Required release documentation remains mandatory | Retained for PSCAN-06/07 and extended to outcome-proof inputs |
| HIST-08 | Proposed/selected/activated/claimed/accepted/completed remain distinct | Retained; every successor remains unselected |

## Downstream planning

PSCAN-10 is proposed and unselected as the bounded primary-coverage successor.
PSCAN-04 through PSCAN-07 remain proposed and unselected, with living
prerequisites updated to require accepted PSCAN-10 primary coverage rather than
accepted PSCAN-03. PSCAN-08 remains inactive, unselected and separately gated.
No transition row selects, activates, claims or implements any successor.
