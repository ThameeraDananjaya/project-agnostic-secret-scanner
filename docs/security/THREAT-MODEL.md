# Threat Model

## Scope and assets

This threat model covers the future local runner, engine adapters, input
preparation, artifact normalization, private finding channel, redaction,
content-free outcomes, release supply chain and reference verification. It does
not grant authority over consuming-project integration, receipts, deployment,
credential response or production.

Protected assets are candidate confidentiality, possible credential values,
finding metadata, scanner and engine integrity, rule and schema integrity,
outcome correctness, exact scan bindings, release identity, project isolation,
transient workspace destruction and evidence-chain trust.

## Trust assumptions

- Candidate content, names, Git metadata, archives and generated artifacts are
  hostile.
- The execution host may contain unrelated credentials or workspaces and must
  isolate them from the scanner.
- Engine output can contain raw secrets or malicious control text.
- Network, upstream releases, tags, caches and mutable remote state are
  untrusted until independently verified and pinned.
- A consuming project can supply malformed or incompatible policy data but
  cannot weaken product safety invariants.
- Maintainer, workflow and signing identity are separate authorities from
  project receipt identity.

## Threat and control register

| ID | Threat | Required preventive controls | Fail-closed result | Planned validation | Task |
|---|---|---|---|---|---|
| T-01 | Shell/workflow injection through paths, refs, messages or fields | Candidate-as-data; argument arrays; fixed executable paths; no candidate workflow execution | Reject or non-pass | Metacharacter and workflow-injection fixtures | PSCAN-02/10/07 |
| T-02 | Candidate code or dependency execution | No build/install/hooks; read-only inputs; trusted preparation only | Non-pass on required executable input | Executable fixture and hook canary | PSCAN-10/07 |
| T-03 | Malicious Git config, hooks, attributes, filters, pagers, diff or text conversion | Fresh safe Git config; disable inheritance and external helpers | Non-pass if safe preparation cannot be proven | Hostile Git configuration matrix | PSCAN-10/07 |
| T-04 | Path traversal, unsafe links, device files, reparse points or case collisions | Canonical path checks; root confinement; reject unsafe entry classes and ambiguity | `INDETERMINATE_UNSUPPORTED_INPUT` | Archive/filesystem adversarial fixtures on Windows and Linux | PSCAN-04/07 |
| T-05 | Archive bomb, malformed container or resource exhaustion | Depth, entry, file, expanded-size, compression, CPU, memory and process bounds; no silent skip | `INDETERMINATE_RESOURCE_LIMIT`, timeout or unsupported | Maximum, over-limit, malformed and bomb fixtures | PSCAN-04/07 |
| T-06 | Secret leakage in stdout, stderr, JSON, logs, errors, crashes or temporary files | Private capture; mandatory redaction firewall; bounded content-free diagnostics; memory-backed temporary output | `INDETERMINATE_REDACTION_UNPROVEN` or terminal | Canary sweep across every output and lifecycle path | PSCAN-04/07 |
| T-07 | Output/annotation injection | Strict outcome serialization; no engine text, paths, counts, names or annotations externally | Non-pass if serializer contract fails | Malicious names/messages/control-character fixtures | PSCAN-02/04/07 |
| T-08 | Engine/scanner/rule/schema substitution | Exact versions and digests; release manifest; signature/bundle; verify before scanning | Integrity/signature fail before scan | One-byte mutation and wrong-identity matrix | PSCAN-10/06/07 |
| T-09 | Mutable tag, rollback, resurrection, retirement bypass or revocation bypass | Exact versions; no `latest`; immutable release where available; signed lifecycle notices; global revocation wins | Reject asset/release | Mutable-reference, rollback, retired and revoked fixtures | PSCAN-06/07 |
| T-10 | Scanner exfiltration, verification, telemetry or update access | Acquisition/execution separation; network-disabled container; no credentials; DNS/TCP canary | Non-pass if isolation is unproven | DNS/TCP and verification-capable engine canaries | PSCAN-04/07 |
| T-11 | Ambient credential exposure from environment, Git, Docker or host | Empty allowlisted environment; no token; isolated mounts; no unrelated workspace; bounded process tree | Workflow-policy failure or terminal | Environment/mount/token permission inspection | PSCAN-04/06/07 |
| T-12 | Cross-project workspace, cache, policy, receipt or identity contamination | Fresh project/attempt workspace; no shared state; randomized non-source-derived scan ID; complete cleanup | Terminal/non-pass | Parallel, sequential and reused-host contamination tests | PSCAN-02/05/07 |
| T-13 | Incomplete coverage reported as pass | Exact admission ledger; every-byte inspection proof; explicit unsupported/limit reasons; fail-closed aggregation | Indeterminate/unavailable/fail | Missing, unreadable, skipped, conflicting, fragment-boundary and partial fixtures | PSCAN-10/04/07 |
| T-14 | Policy or allowlist weakens mandatory safety | Fixed precedence; version/digest validation; genuine credentials never allowlisted | `FAIL_POLICY_DENIED` or `FAIL_ALLOWLIST_INVALID` | Precedence, expiry, scope and credential exception tests | PSCAN-05/07 |
| T-15 | Schema confusion or semantic drift | Independent `major.minor`; unknown major rejection; immutable field meanings; signed retirement | Indeterminate/reject | Compatibility and corrupt-request matrix | PSCAN-02/05/07 |
| T-16 | Receipt authority crosses into scanner | Outcome-only boundary; no receipt keys; independent trust roots; reference verifier is custody-neutral | Reject mixed authority | Mock issuer proves key never enters scanner | PSCAN-05/07 |
| T-17 | Evidence rollback, truncation, duplicate sequence or divergent chain | Append-only markers; signed chain facts; monotonic/previous-hash verification; separate checkpoints | Reject evidence | Rollback/divergence/conflict/revocation fixtures | PSCAN-05/07 |
| T-18 | Retry overwrites or converts earlier outcome | New scan ID and clean workspace; earlier immutable; deterministic failures do not retry | Terminal after eligible exhaustion | Retry state-machine tests | PSCAN-02/07 |
| T-19 | Cleanup failure preserves findings or candidate material | Cleanup on success, fail, timeout, cancellation and crash; post-cleanup proof | Terminal/invariant failure | Lifecycle residue inspection | PSCAN-04/07 |
| T-20 | Public release or CI authority silently approves project action | Content-free evidence wording; CI-provider-neutral CLI; separate owner/project gates | No authorized action | Contract assertions and integration review | PSCAN-02/06/07 |
| T-21 | Framing or transformation masks a raw archive/binary/ambiguous class | Raw-byte classification before transformation; class authority bound in ledger | Unsupported/incomplete non-pass | Archive signatures, polyglots and framing-mask fixtures | PSCAN-10/04/07 |
| T-22 | Detector fragments skip an internal match while a prefix marker passes | Pinned span/stream proof; finite maximum rule spans or complete streaming; per-byte inspection evidence | Incomplete-coverage non-pass | First/last/internal/no-whitespace span-edge fixtures | PSCAN-10/07 |
| T-23 | Aggregate failure falsely proves a claimed class | Isolated clean/finding/non-pass fixtures per class and profile | Class remains unproved and non-pass | Text, binary, history, tree, archive and maximum-size isolation | PSCAN-10/04/07 |

## Abuse cases

1. A candidate names a file using shell metacharacters and embeds a malicious
   workflow. The trusted runner passes the path as data, never evaluates the
   workflow, and returns content-free evidence.
2. An archive declares a small compressed size but expands past policy. The
   normalizer stops safely, destroys partial output and returns indeterminate.
3. An engine prints a synthetic secret to stderr and then crashes. Private
   capture and redaction prevent disclosure; if redaction cannot be proven the
   outcome is non-pass.
4. A second project runs on the same host after a failed scan. No workspace,
   cache, policy, identity or finding from the first run is observable.
5. A consumer provides an exact-looking but wrong release asset. Verification
   rejects it before candidate bytes are opened.
6. An expired exception attempts to suppress a mandatory credential class. The
   exception is invalid, the finding blocks, and no manual pass is available.

## Residual and owner-controlled risks

- Public repository, GitHub rules, Actions permissions, immutable-release and
  signing-identity configuration remain action-time PSCAN-06/07 owner gates.
- Receipt key custody, evidence store, deployment controls and incident response
  are consuming-project risks and remain outside product authority.
- TruffleHog legal and technical risk is not assessed by PSCAN-01; the component
  remains absent until the PSCAN-08 eligibility and owner gates are satisfied.
- The final retention schedule requires separate legal/compliance approval.

Any newly discovered material threat that cannot be controlled inside an
activated task fails that task closed and returns to the owner with bounded
alternatives.

## PSCAN-05 accepted controls

PSCAN-05 implements the independently accepted in-scope controls for T-12 and T-14
through T-17: transient per-call projections, fixed precedence, mandatory
credential-class blocking, exact exception context and expiry, independent
schema windows, public-key-only signature verification, exact receipt binding,
trusted receipt/global/project checkpoints, and append-only revocation-chain
verification. Its API contains no private-key, signing, receipt-issuance,
evidence-storage, promotion or project-identity authority. These controls do
not replace PSCAN-07 whole-product validation.
