# PASS-SPEC-001 Requirement Traceability

## Interpretation

This matrix preserves the controlling contract without claiming implementation.
“Boundary” names the target architecture and first authoritative source area;
“Test” names required future proof; “Task” names the bounded implementation or
acceptance owner. PSCAN-01 provides only this planning trace.

## Capability traceability

| ID | Preserved success requirement | Architecture and file boundary | Required proof | Task |
|---|---|---|---|---|
| CAP-1 | Two structurally different fixture projects use one unmodified release with independent policy, workspace and adapter | CLI/request/workspace/policy projection; `cmd/`, `internal/request`, `internal/workspace`, `internal/policy` | Two-project contract and isolation acceptance | PSCAN-02, PSCAN-05, PSCAN-07 |
| CAP-2 | PR outcome binds exact base, head, merge base, history range, tracked tree, release, rules, policy, allowlist and schemas | Request, Git input, outcome; `contracts/scan-request`, `internal/gitinput`, `contracts/scan-outcome` | Exact-range/tree mutation and binding tests | PSCAN-02, PSCAN-03, PSCAN-07 |
| CAP-3 | Full-release scan covers history, tracked source, build context, deployment files, container layers and nested archives | Request, Git input, artifact normalizer; `internal/gitinput`, `internal/artifact` | Frozen-candidate coverage and missing-input tests | PSCAN-02, PSCAN-03, PSCAN-04, PSCAN-07 |
| CAP-4 | Complete acceptance succeeds with DNS/TCP unavailable and no provider/update/telemetry/source transmission | Execution sandbox; `internal/workspace`, `build`, trusted workflow | Network-canary and credential-free acceptance | PSCAN-04, PSCAN-06, PSCAN-07 |
| CAP-5 | Exact Gitleaks source, binary, configuration, command, version and output binding; mismatch stops before scan | Engine orchestrator/Gitleaks adapter; `internal/engine/gitleaks`, release manifest | Pin/digest/format mutation tests | PSCAN-03, PSCAN-06, PSCAN-07 |
| CAP-6 | Fallback is disabled by default, requires accepted gap/request binding, is offline/unverified, and cannot weaken Gitleaks fail | Gated fallback boundary; `internal/engine/trufflehog` absent until eligible | Accepted coverage map; disabled/build/request gates; conflict test | PSCAN-07 eligibility proof; PSCAN-08 only if triggered |
| CAP-7 | Directories, ZIP, TAR, gzip/TGZ, OCI, Docker-save and supported nesting scan or explicitly return indeterminate | Artifact normalizer; `internal/artifact` | Format, traversal, collision, malformed and limit matrix | PSCAN-04, PSCAN-07 |
| CAP-8 | No secret or forbidden detail reaches output, logs, exceptions, crash data or temporary storage | Private channel/redaction/cleanup; `internal/redaction`, `internal/cleanup` | Canary sweep on every lifecycle/output path | PSCAN-04, PSCAN-07 |
| CAP-9 | Same release honors two independently governed project policies without shared state | Policy projection boundary; `internal/policy` | Two-project divergent-policy and contamination tests | PSCAN-05, PSCAN-07 |
| CAP-10 | Exceptions are owner-approved, exact, evidenced, at most 30 days, invalidating; genuine credentials never allowed | Project-owned allowlist projection; `internal/policy` | Expiry, scope, renewal, mutation and credential rejection | PSCAN-05, PSCAN-07 |
| CAP-11 | Identical inputs yield semantically identical outcomes across repeat, Windows and Linux runs | Request/outcome normalization; `internal/outcome` | Repeatability and cross-platform parity | PSCAN-02, PSCAN-06, PSCAN-07 |
| CAP-12 | Unsupported schema/input, unreadable/oversized/skip, engine/integrity/redaction/conflict/missing cases are non-pass | Validation, orchestration, normalizer, classifier | Fail-closed negative matrix | PSCAN-02 through PSCAN-07 |
| CAP-13 | JSON/exit expose only approved status, reason, bindings, versions, attempt, duration and action | CLI/outcome/redaction; `cmd`, `contracts/scan-outcome`, `internal/outcome` | Schema snapshots, exit mapping and forbidden-field tests | PSCAN-02, PSCAN-04, PSCAN-07 |
| CAP-14 | Offline verification proves manifest, assets, exact workflow identity, bundle, licences and SBOM | Release manifest/verifier/build/licences | Wrong identity, issuer, tag, bundle, digest, notice and SBOM tests | PSCAN-06, PSCAN-07 |
| CAP-15 | Independent contract versions accept compatible minor, reject unknown major and obey retirement overlap | All `contracts`, policy projection and compatibility registry | Compatibility, retirement and historical-reader tests | PSCAN-02, PSCAN-05, PSCAN-06, PSCAN-07 |
| CAP-16 | Outcome supplies enough bindings for project receipt while scanner has no receipt key or promotion authority | Outcome and reference verifier; scanner/project trust boundary | Mock external issuer with key-absence proof | PSCAN-02, PSCAN-05, PSCAN-07 |
| CAP-17 | Append-only global/project markers invalidate affected evidence without editing original records | Global revocation schema/reference chain verifier | Release/key/receipt/policy/exception revocation matrix | PSCAN-02, PSCAN-05, PSCAN-07 |
| CAP-18 | No source, finding, policy, allowlist, cache, receipt, key or identity leaks between projects | Workspace/policy/cleanup isolation | Parallel, sequential and reused-host contamination tests | PSCAN-02, PSCAN-05, PSCAN-07 |
| CAP-19 | GA assets support tested Windows amd64 and Linux amd64 with no unsupported architecture claim | Build/release compatibility matrix | Cross-platform build, verify and logical parity | PSCAN-06, PSCAN-07 |
| CAP-20 | No auto-update; rollback only to supported unrevoked release; retirement/revocation enforced | Release verifier and lifecycle docs | Update denial, rollback, retired and globally revoked tests | PSCAN-06, PSCAN-07 |

### CAP-6 readiness rule

PSCAN-07 must produce an accepted coverage map. If Gitleaks covers every required
input class, PSCAN-08 remains inactive and the fallback stays absent. If a
material gap exists, first-consumer readiness fails closed until the owner
separately approves both the technical need and AGPL obligations and PSCAN-08
implements and passes CAP-6. A planning label is not eligibility or approval.

## Acceptance-row traceability

| ID / contract area | Preserved expected result | Architecture and file boundary | Required test | Task |
|---|---|---|---|---|
| AT-01 Clean source | Known-safe fixtures pass | Fixture harness, orchestrator, classifier; `fixtures/clean`, `tests/acceptance` | Clean fixture matrix | PSCAN-03, PSCAN-07 |
| AT-02 Detection | Synthetic secrets fail without disclosure | Gitleaks/private channel/redaction | Synthetic detector and output-canary test | PSCAN-03, PSCAN-04, PSCAN-07 |
| AT-03 History | Added-then-deleted secret in range fails | Git input/Gitleaks Git mode; `fixtures/history` | Generated deleted-history repository | PSCAN-03, PSCAN-07 |
| AT-04 Range binding | Out-of-range finding is not attributed; exact range remains recorded | Request/Git input/outcome | Boundary commits and range-digest assertions | PSCAN-02, PSCAN-03, PSCAN-07 |
| AT-05 Tracked source | Manifest and tree are complete or non-pass | Request/Git input | Missing/unreadable/tree-mutation tests | PSCAN-02, PSCAN-03, PSCAN-07 |
| AT-06 Build context | Generated/untracked build inputs are complete or non-pass | Release request/artifact manifest | Manifest omission and digest mutation | PSCAN-02, PSCAN-04, PSCAN-07 |
| AT-07 Artifacts | Deployment files and image layers are complete or non-pass | Artifact normalizer | File, OCI and Docker-save coverage | PSCAN-04, PSCAN-07 |
| AT-08 Archives | Supported nested archive detects synthetic finding | Artifact normalizer; `fixtures/archives` | Supported-format/depth matrix | PSCAN-04, PSCAN-07 |
| AT-09 Unsupported archive | Encrypted, malformed or unsafe archive is indeterminate | Artifact normalizer | Encryption, header, traversal, link and collision fixtures | PSCAN-04, PSCAN-07 |
| AT-10 Limits | Oversized/bomb input is indeterminate, never skip/pass | Resource accounting/normalizer | Each declared limit at/below/above threshold | PSCAN-04, PSCAN-07 |
| AT-11 Redaction | Standard and crash outputs contain no secret or forbidden metadata | Private channel/redaction/cleanup | Multi-surface canary sweep | PSCAN-04, PSCAN-07 |
| AT-12 Output injection | Malicious names/messages yield valid content-free output | Serializer/redaction | Control, newline, markup and annotation injection | PSCAN-02, PSCAN-04, PSCAN-07 |
| AT-13 Offline | DNS/TCP unavailable; deterministic scan | Network-disabled sandbox | DNS/TCP and verification-capable engine canary | PSCAN-04, PSCAN-07 |
| AT-14 Integrity | Changed engine, runner or rule asset fails before scan | Release/engine binding verifier | One-byte mutation for every asset class | PSCAN-03, PSCAN-06, PSCAN-07 |
| AT-15 Signing | Wrong identity, issuer, tag or bundle rejects release | Offline release verifier | Signature-policy negative matrix | PSCAN-06, PSCAN-07 |
| AT-16 Schema | Unknown major or corrupt request is indeterminate/rejected | Request validator/contracts | Compatibility and corruption corpus | PSCAN-02, PSCAN-05, PSCAN-07 |
| AT-17 Policy | Invalid signature, digest or precedence fails closed | Policy projection/evaluator | Signature/digest/precedence mutations | PSCAN-05, PSCAN-07 |
| AT-18 Allowlist | Expired, broadened or real-credential exception is rejected and fails | Allowlist projection/evaluator | Expiry, scope and mandatory-class cases | PSCAN-05, PSCAN-07 |
| AT-19 Retry | Eligible transient condition gets two retries then terminal | Project-adapter contract/state machine | Three-attempt clean-workspace sequence | PSCAN-02, PSCAN-07 |
| AT-20 Deterministic fail | Confirmed synthetic finding does not retry | Outcome/retry contract boundary | Fail-state retry denial | PSCAN-02, PSCAN-07 |
| AT-21 Engine conflict | Either finding fails; otherwise incomplete/disagreement is indeterminate | Engine aggregation | Required-engine result cross-product | PSCAN-07; PSCAN-08 if fallback eligible |
| AT-22 Cleanup | Pass/fail/timeout/crash leaves no transient material | Workspace/cleanup | Lifecycle residue inspection | PSCAN-04, PSCAN-07 |
| AT-23 Receipt reuse | Identical staging/production bindings reuse within 30 days | Reference chain verifier | Exact-binding/time-bound mock receipt | PSCAN-05, PSCAN-07 |
| AT-24 Receipt invalidation | Any binding change rejects receipt | Reference chain verifier | One-field-at-a-time binding mutations | PSCAN-05, PSCAN-07 |
| AT-25 Revocation | Scanner, key, receipt, policy or exception marker rejects affected evidence | Revocation/reference verifier | Authorized-source and effective-time matrix | PSCAN-05, PSCAN-07 |
| AT-26 Evidence rollback | Truncated, divergent or conflicting chain blocks progression | Reference chain verifier | Sequence, hash, head and checkpoint attacks | PSCAN-05, PSCAN-07 |
| AT-27 Isolation | Parallel projects and reused host do not contaminate | Workspace/policy/cleanup | Concurrent/sequential two-project canaries | PSCAN-02, PSCAN-05, PSCAN-07 |
| AT-28 Untrusted PR | Execution/workflow injection remains data only | Git input/workspace/trusted workflow | Hook, executable and candidate-workflow canaries | PSCAN-03, PSCAN-06, PSCAN-07 |
| AT-29 Permissions | Extra token rights or supplied secrets fail workflow policy | Trusted workflow/sandbox | Permission and environment inspection | PSCAN-06, PSCAN-07 |
| AT-30 Windows/Linux | Logical fixture has equivalent outcome | Cross-platform build/outcome | Full parity corpus on amd64 | PSCAN-06, PSCAN-07 |
| AT-31 Performance | Approved profiles complete within mode timeout | Resource accounting/acceptance harness | Small, medium and declared-maximum profiles | PSCAN-04, PSCAN-07 |

## Other controlling boundaries

- Sections 1–5 are preserved by `CONSTITUTION.md`, `AGENTS.md`, architecture,
  lifecycle and owner-gate documents.
- Sections 7–26 are decomposed across the component architecture, source layout,
  threat model, schema decision, validation plan, licensing and release plans.
- Sections 28–29 are mandatory PSCAN-06 release-documentation and permanent
  non-goal checklists.
- Section 30 is represented by the bounded PSCAN-01 through PSCAN-08 task set.
- Section 31 remains the PSCAN-07 product definition of done and is not claimed
  by PSCAN-01 completion.
- Sections 32–33 remain explicit action-time gates and version-sensitive
  preflights; no planning document converts them into approval.
