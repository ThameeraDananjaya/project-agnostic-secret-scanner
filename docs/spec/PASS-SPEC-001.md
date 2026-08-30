---
title: Project-Agnostic Secret Scanner — Full Design Requirements
document_id: PASS-SPEC-001
version: 1.0
status: owner-approved planning contract; implementation requires separately authorized work
date: 2026-08-30
owner: Project owner
primary_consumer: RestoHub
future_consumers: Any owner-approved software project
recommended_product_license: MIT
---

# Project-Agnostic Secret Scanner — Full Design Requirements

> **Canonical product contract.** This document defines what the new scanner product must build, protect, test, release, and prove. A downstream implementation may add stricter controls but must not weaken these requirements. Project-specific repositories remain authoritative for their own policy, allowlists, release receipts, revocations, credentials, deployment, and business decisions.

## 1. Purpose

The product will provide one independently versioned secret-scanning package that can be reused by RestoHub and future owner-controlled projects. It must scan local project material without sending source, findings, or possible credentials to an external scanning service. It must detect likely secrets early during implementation and perform a mandatory, deeper scan of every release candidate before staging or production promotion.

The scanner is a security evidence producer, not a deployment authority. A passing result never approves a pull request, deployment, production use, legal compliance, or go-live.

## 2. Product boundary

The product is a local command-line package and release supply chain. It is not a SaaS service, source-code host, credential validator, incident-response platform, per-finding database, deployment controller, or cross-project policy authority.

The initial external platform boundary is:

```text
Public Scanner Source Repository on GitHub
    -> trusted GitHub release workflow
    -> signed immutable scanner release
    -> Sigstore public signing proof
    -> consuming project downloads and verifies the exact release
    -> scanner executes locally with network disabled
```

The runtime boundary is:

```text
Consuming Project Adapter
    -> validated scan request
    -> Project-Agnostic Go Runner
    -> Gitleaks CLI primary engine
    -> TruffleHog OSS only for an approved coverage gap
    -> redaction firewall
    -> content-free outcome
    -> project-owned receipt issuer and verifier
```

GitHub and Sigstore are external platforms. Gitleaks, TruffleHog, Cosign, and the Go runner are software components, not online scanning services.

## 3. Governing decisions

The following choices are settled product direction:

- The scanner is a separate, independently versioned, project-agnostic product.
- GitHub native secret scanning is not the product and is not a dependency.
- Gitleaks CLI is the primary detection engine.
- TruffleHog OSS is a disabled-by-default fallback for proven coverage gaps, subject to legal approval.
- The orchestration runner is a thin Go command-line program licensed under MIT.
- The generic scanner repository is intended to be public and must contain no project data.
- Distribution uses exact-version GitHub Release assets; `latest` is prohibited.
- Production-consumable releases are immutable where GitHub makes that control available.
- Public scanner releases use Cosign keyless signing and Sigstore transparency evidence.
- Scanner execution is offline by default and cannot validate credentials with providers.
- GitHub-hosted Ubuntu execution is the initial authoritative CI model, with the scanner isolated in a network-disabled container.
- Each consuming project owns its policy, allowlists, receipt signing, receipt verification, revocations, evidence custody, deployment gates, and retention decisions.
- RestoHub receipts use a trust root distinct from the public scanner-release trust root.
- Findings are transient and fully redacted. There is no durable per-finding database.
- Missing, stale, incomplete, conflicting, unsupported, or untrustworthy evidence fails closed.
- Cross-project reporting is disabled by default and requires separate owner opt-in from every participating project.

## 4. Terminology

| Term | Definition |
|---|---|
| Scanner release | One exact signed product bundle containing the runner, primary engine, generic rules, schemas, licences, SBOM, tests, and integrity evidence. |
| Consuming project | RestoHub or another project invoking the scanner through its own adapter. |
| Project adapter | Project-owned integration that constructs scan requests, selects policy, invokes the runner, interprets outcomes, retries, and issues project receipts. |
| PR scan | Early scan of an implementation pull-request candidate. |
| Full-release scan | Mandatory scan of an exact frozen release candidate, including history, source, build context, and generated deployment artifacts. |
| Frozen candidate | A release commit and artifact set whose recorded bytes and digests cannot change without creating a new candidate. |
| Finding | A transient engine detection that may indicate a secret. |
| Outcome | Content-free result returned by the scanner: pass, fail, indeterminate, unavailable, or terminal. |
| Receipt | Project-issued, content-free, signed evidence that one exact full-release binding passed. |
| Revocation marker | Append-only evidence that a scanner release, receipt, key, exception, or other authority must no longer be trusted. |
| Binding | A recorded digest, version, revision, policy, artifact, or trust fact that must remain identical for evidence reuse. |

## 5. Roles and authority

| Role | Authority | Explicit limits |
|---|---|---|
| Product owner | Approves product scope, publication, legal posture, signing identity, spending, and production use. | Approval of design does not itself publish, spend, create credentials, or activate implementation. |
| Scanner maintainer | Implements and releases the generic scanner under approved tasks. | Cannot access project source, project keys, receipts, findings, or deployment authority. |
| Scanner release authority | Publishes an accepted scanner release through the trusted release workflow. | Cannot issue project release receipts. |
| Consuming-project owner | Owns project policy, allowlists, receipt trust, deployment gates, and exceptions. | Cannot weaken scanner safety invariants or treat a pass as deployment approval. |
| Project adapter maintainer | Integrates the scanner and maintains project-specific schemas and workflow code. | Cannot change the external scanner release or generic global revocations. |
| Receipt issuer | Issues a receipt only after independently verifying a passing outcome and every binding. | Cannot approve deployment or retroactively alter a receipt. |
| Receipt verifier | Accepts or rejects receipts and revocation chains. | Cannot issue or revoke evidence. |
| Incident/revocation authority | Emergency-revokes evidence after credible compromise or correctness failure. | Cannot restore revoked evidence; restoration requires new clean evidence. |
| Developer | Runs optional local scans and remediates candidates. | Cannot approve their own exception or bypass required scans. |
| Codex/automation | Implements bounded work, runs checks, reviews evidence, and reports status. | Cannot handle real credentials, approve reserved owner gates, or convert uncertainty into pass. |
| Auditor | Reads content-free evidence, chain proofs, versions, and retention records. | Cannot access raw findings or secret values through the audit interface. |

## 6. Capability contract

### CAP-1 — Project-independent execution

- **Intent:** Any supported project can invoke the same scanner release through a stable local contract without adding project-specific logic to the scanner.
- **Success:** Two structurally different fixture projects use the same unmodified release while supplying independent policies, workspaces, and adapters.

### CAP-2 — Early PR scanning

- **Intent:** A project can scan an exact implementation candidate before merge.
- **Success:** The outcome is bound to the exact base, head, merge base, history range, tracked tree, scanner release, rules, policy, allowlist, and schema versions.

### CAP-3 — Full-release scanning

- **Intent:** A project can scan the complete frozen release candidate before staging or production.
- **Success:** The scan proves coverage of the required Git range, tracked source, build context, generated deployment files, container layers, and supported nested archives.

### CAP-4 — Offline operation

- **Intent:** Scanning can complete without provider calls, credential verification, update checks, telemetry, or source transmission.
- **Success:** The complete acceptance suite passes inside a network-disabled sandbox and a network canary proves outbound DNS and TCP are unavailable to the scanner process.

### CAP-5 — Primary engine orchestration

- **Intent:** The runner can invoke one pinned Gitleaks CLI release predictably and safely.
- **Success:** Exact source, binary, configuration, command, version, and output-format bindings are recorded, and any mismatch fails before source is scanned.

### CAP-6 — Controlled fallback coverage

- **Intent:** A project can invoke TruffleHog OSS only when an accepted coverage map proves Gitleaks cannot inspect a required input class.
- **Success:** Fallback execution is disabled by default, requires an explicit request binding, runs with provider verification and updates disabled, and cannot weaken a Gitleaks failure.

### CAP-7 — Artifact normalization

- **Intent:** Supported deployment artifacts can be presented as safe, deterministic local scan inputs.
- **Success:** ZIP, TAR, gzip/TGZ, OCI image layouts, Docker-save tar files, ordinary directories, and supported nested combinations are scanned or explicitly rejected as indeterminate; no content is silently skipped.

### CAP-8 — Redaction firewall

- **Intent:** No raw secret or detailed finding crosses the scanner output boundary.
- **Success:** Adversarial tests demonstrate that stdout, stderr, exceptions, reports, logs, job summaries, crash output, temporary files, and returned JSON contain no secret value or forbidden finding detail.

### CAP-9 — Project policy

- **Intent:** Each consuming project can define its own versioned scan scope, blocking rules, severity behavior, and permitted exception categories.
- **Success:** The same scanner release produces independently governed outcomes for two projects without sharing policy or state.

### CAP-10 — Governed allowlists

- **Intent:** A project owner can create narrow, expiring exceptions for proven non-secret material.
- **Success:** Every exception has an owner, rationale, exact scope, approval, creation time, expiry of at most 30 days, and invalidation rules; genuine credentials are always rejected from allowlisting.

### CAP-11 — Deterministic outcomes

- **Intent:** Identical validated inputs produce the same public outcome and bindings.
- **Success:** Repeated Windows and Linux fixture runs produce semantically identical content-free results, excluding timestamps, duration, and correlation identifiers.

### CAP-12 — Fail-closed behavior

- **Intent:** The scanner never reports pass when it cannot prove complete trusted coverage.
- **Success:** Unsupported schemas, unreadable files, skipped oversized inputs, engine errors, integrity failures, redaction failures, contradictory required engines, or missing artifacts return non-pass outcomes.

### CAP-13 — Content-free outcome contract

- **Intent:** CI systems and project adapters can consume a stable result without receiving findings.
- **Success:** The public JSON outcome and CLI exit code expose only approved status, reason code, binding digests, versions, attempt metadata, duration, and permitted action.

### CAP-14 — Supply-chain verification

- **Intent:** A consumer can verify exactly who built and published the scanner package and that no asset changed.
- **Success:** Offline verification proves the release manifest, every asset digest, exact release-workflow identity, signature bundle, licences, and SBOM before execution.

### CAP-15 — Versioned contracts

- **Intent:** Runner, engine, rules, request, outcome, release-manifest, policy, allowlist, receipt, and revocation formats can evolve independently.
- **Success:** Compatible readers accept supported minor versions, unknown major versions fail closed, and retirement obeys the approved overlap rules.

### CAP-16 — Receipt integration

- **Intent:** A consuming project can issue its own receipt after independently validating a full-release pass.
- **Success:** The scanner produces enough content-free binding evidence for the project issuer while possessing no project receipt-signing key and no authority to promote a release.

### CAP-17 — Revocation integration

- **Intent:** Scanner releases and project receipts can be invalidated without editing historical evidence.
- **Success:** Append-only global and project markers cause verifiers to reject affected evidence while preserving the original signed record.

### CAP-18 — Cross-project isolation

- **Intent:** No source, finding, policy, allowlist, cache, receipt, key, or identity leaks between consuming projects.
- **Success:** Parallel contamination tests prove independent workspaces, temporary stores, configuration, correlation IDs, and evidence.

### CAP-19 — Cross-platform delivery

- **Intent:** Supported Windows and Linux users can verify and run the same logical product contract.
- **Success:** GA releases include tested `windows/amd64` and `linux/amd64` assets; any additional architecture is labelled supported only after equivalent acceptance.

### CAP-20 — Safe release lifecycle

- **Intent:** Maintainers can update, roll back, retire, and revoke scanner releases without weakening consumers.
- **Success:** Consumers never auto-upgrade, roll back only to a still-supported unrevoked release, and reject retired or globally revoked dependencies according to policy.

## 7. Required product components

### 7.1 Go runner

The runner is the stable public interface. It must be small, deterministic, statically buildable where practical, and independent of any consuming project. It must invoke engine executables as subprocesses through argument arrays, never through interpolated shell commands. Candidate-controlled text must never become executable syntax.

The runner owns request validation, workspace setup, engine selection, input preparation, output capture, redaction, in-memory deduplication, outcome classification, resource accounting, and cleanup. It does not own retries, project receipt issuance, GitHub status publication, deployment, or credential response.

### 7.2 Gitleaks adapter

The primary adapter must:

- Pin one exact Gitleaks source revision and binary digest.
- Use Git mode for exact history/range scanning.
- Use directory/file mode for tracked trees, build contexts, normalized artifacts, and extracted container layers.
- Force complete redaction and suppress banners, color, verbose output, and interactive behavior.
- Capture all engine output inside private memory or memory-backed temporary storage.
- Treat target-size skips, timeouts, unsupported inputs, or incomplete coverage as non-pass.
- Never use the separately licensed Gitleaks GitHub Action.

### 7.3 TruffleHog fallback adapter

The fallback adapter must exist behind a build-time and request-time disabled gate. Its actual binary must not be distributed until the owner accepts the applicable AGPL obligations.

When enabled, it must:

- Use only local Git, filesystem, or local image-tar inputs.
- Disable credential verification, verification caches, remote enumeration, telemetry, updates, and provider access.
- Run inside the same network-disabled sandbox.
- Capture raw engine output before it reaches stdout or stderr.
- Treat any finding from either required engine as fail.
- Return indeterminate when required fallback coverage cannot complete.

### 7.4 Artifact normalizer

The normalizer must safely inspect files without executing them. It must reject path traversal, absolute archive paths, device files, unsafe links, duplicate-path ambiguity, case-collision ambiguity, encrypted archives, malformed headers, decompression bombs, and extraction outside the workspace.

Initial engineering limits are versioned policy defaults, not silent skips:

- Maximum nested archive depth: 5.
- Maximum archive entries per scan: 100,000.
- Maximum expanded bytes: 2 GiB for PR mode and 10 GiB for release mode.
- Maximum individual regular file: 512 MiB.
- Maximum declared compression ratio: 1,000:1.

Crossing a limit returns `INDETERMINATE_RESOURCE_LIMIT`. A consuming project may adopt stricter or higher reviewed limits through a new policy version; it cannot configure “skip and pass.”

### 7.5 Policy evaluator

Policy precedence is:

1. Non-overridable scanner safety invariants.
2. Global scanner-release revocations and mandatory base protections.
3. Consuming-project policy.
4. Narrow project allowlist exceptions.

Lower layers cannot weaken higher layers. Any credible private key, access token, password, connection credential, signing key, or live-looking authentication material is blocking by default. An unallowlisted blocking finding produces fail regardless of nominal severity.

### 7.6 Redaction firewall

The firewall is a mandatory process boundary between engine output and every external surface. It must redact 100% of candidate secret values and remove forbidden metadata before any output is written. It must operate even when an engine crashes or returns malformed output.

Standard CI and console status must not expose:

- Secret text, partial secret text, encodings, hashes, or reversible fingerprints.
- File or directory paths.
- Line or column numbers.
- Detector or provider names.
- Git author, email, commit message, or branch name.
- Finding excerpts or counts.
- Allowlist details.

An explicitly authorized local remediation session may reveal only the minimum redacted location needed to correct the project, never the candidate secret value. That session is project-owned, transient, inaccessible to ordinary CI, and outside receipt evidence.

### 7.7 Contract and schema package

The external product owns:

- `scan-request` schema.
- `scan-outcome` schema.
- `scanner-release-manifest` schema.
- `global-scanner-revocation` schema.
- Generic rule-pack schema.

Each consuming project owns:

- Project policy schema and instances.
- Project allowlist schema and instances.
- Project receipt schema and instances.
- Project receipt-revocation schema and instances.
- Project evidence-chain and custody rules.

The scanner repository may publish reference examples for project-owned schemas but cannot become their authority.

### 7.8 Fixture and acceptance harness

The product must include safe synthetic fixtures for clean source, fake secrets, deleted-history secrets, large files, binary files, nested archives, container layers, malformed archives, path traversal, symlinks, policy errors, expired exceptions, conflicting engines, output injection, network attempts, corrupted signatures, rollback, cross-project contamination, Windows paths, and Linux paths.

Fixtures must never contain real credentials.

## 8. Proposed source repository layout

```text
/.github/workflows/             trusted CI and release workflows
/cmd/scanner-runner/            public CLI entry point
/internal/request/              request/schema validation
/internal/workspace/            isolated workspace lifecycle
/internal/gitinput/             safe Git-range preparation
/internal/artifact/             archive and OCI normalization
/internal/engine/               engine interface and orchestration
/internal/engine/gitleaks/      primary adapter
/internal/engine/trufflehog/    disabled fallback adapter
/internal/policy/               precedence and exception evaluation
/internal/redaction/            mandatory output firewall
/internal/outcome/              states, reason codes, exit mapping
/internal/cleanup/              transient-material destruction
/contracts/scan-request/
/contracts/scan-outcome/
/contracts/release-manifest/
/contracts/global-revocation/
/rules/generic/
/fixtures/clean/
/fixtures/synthetic-findings/
/fixtures/history/
/fixtures/archives/
/fixtures/containers/
/fixtures/adversarial/
/tests/unit/
/tests/integration/
/tests/acceptance/
/build/                         reproducible build and packaging definitions
/licenses/                      runner and third-party licence evidence
/docs/                          integration, security, release, recovery, retirement
```

No directory may contain consuming-project source, policies, allowlists, receipts, keys, findings, repository names, customer data, or project statistics.

## 9. Scan request contract

The request must be schema validated before any engine starts. It contains:

| Field | Requirement |
|---|---|
| `requestSchemaVersion` | Exact independent schema version. |
| `scanId` | New random correlation identifier; not derived from source. |
| `mode` | `local`, `pr`, or `release`. |
| `scannerReleaseDigest` | Exact expected external release bundle digest. |
| `engineBinding` | Exact engine name, version, binary digest, and required adapter version. |
| `rulePackDigest` | Exact generic rule-pack digest. |
| `policyDigest` | Exact consuming-project policy digest. |
| `allowlistDigest` | Exact allowlist digest or explicit empty-set digest. |
| `sourceBinding` | Base/head/merge-base/tree and declared history-range digests appropriate to the mode. |
| `trackedSourceManifest` | Local transient path plus expected manifest digest. |
| `buildContextManifest` | Required in release mode; exact files and digest. |
| `artifactManifest` | Required in release mode; local paths, types, sizes, and expected digests. Paths are transient and never emitted. |
| `fallbackRequirement` | `disabled` or exact approved coverage reason and fallback binding. |
| `limits` | Time, depth, entry, expanded-size, file-size, memory, and CPU limits. |
| `offlineRequired` | Must be `true` for authoritative PR and release modes. |
| `redactionMode` | Must be `full`. |
| `requestedAt` | UTC timestamp used only for request freshness. |

Unknown required fields, unsupported major versions, mismatched digests, missing manifests, duplicate artifacts, or unsafe local paths reject the request before scanning.

## 10. Scan input coverage

### 10.1 Local developer mode

Local mode is optional and advisory. It may scan the working tree or staged changes through a project adapter. It cannot issue a release receipt and cannot replace PR or release scanning.

### 10.2 PR mode

PR mode must bind and scan:

- Exact trusted base commit.
- Exact candidate head commit.
- Exact merge base.
- Required commit range and patches.
- Exact tracked source tree at the candidate head.
- Project policy and allowlist versions.
- Scanner, engine, rules, request, and outcome versions.

PR mode does not run candidate code, install candidate dependencies, execute candidate workflows, or build candidate artifacts.

### 10.3 Release mode

Release mode occurs after the release candidate is built and frozen by digest. It must bind and scan:

- Exact candidate release commit.
- Exact previous accepted release commit and complete intervening range; the first release scans all required reachable history.
- Exact tracked source tree.
- Exact clean build-context manifest.
- Every generated deployment artifact and digest.
- Exported OCI/Docker image tar files and every supported layer.
- Supported nested archives.
- Build provenance identifier.
- All scanner, rules, policy, allowlist, schema, and trust bindings.

The flow is:

```text
PR scan
    -> merge through normal project authority
    -> clean release build
    -> freeze commit and artifact digests
    -> mandatory full-release scan
    -> project validates pass and issues signed receipt
    -> staging verification
    -> production verification and separate owner approval
```

If any candidate byte, commit, artifact, policy, allowlist, rule pack, scanner release, or recorded binding changes, the old scan cannot cover the new candidate. The project must build a new candidate and run a new full-release scan.

## 11. Outcome state machine

```text
PENDING
  -> RUNNING
      -> PASS
      -> FAIL
      -> INDETERMINATE
      -> UNAVAILABLE

INDETERMINATE or UNAVAILABLE
  -> RETRYING
      -> RUNNING
      -> TERMINAL after retry exhaustion

PASS or FAIL or TERMINAL
  -> immutable run outcome

RECOVERY
  -> a new scanId that references the earlier scan
  -> never edits or replaces the earlier outcome
```

The project adapter, not the runner, owns retry scheduling. The approved default is one initial attempt plus two retries. PR timeout is 15 minutes per attempt. Full-release timeout is 60 minutes per attempt. Deterministic fail outcomes, invalid policies, bad signatures, and binding mismatches do not retry. Only plausibly transient indeterminate or unavailable outcomes retry in a clean workspace.

## 12. Reason-code contract

### Pass

- `PASS_NO_BLOCKING_FINDINGS`

### Fail

- `FAIL_FINDING_DETECTED`
- `FAIL_POLICY_DENIED`
- `FAIL_ALLOWLIST_INVALID`
- `FAIL_BINDING_MISMATCH`
- `FAIL_INPUT_INTEGRITY`
- `FAIL_SCANNER_INTEGRITY`
- `FAIL_SIGNATURE_INVALID`
- `FAIL_REVOKED_DEPENDENCY`

### Indeterminate

- `INDETERMINATE_UNSUPPORTED_INPUT`
- `INDETERMINATE_INCOMPLETE_COVERAGE`
- `INDETERMINATE_SCHEMA_UNSUPPORTED`
- `INDETERMINATE_CONFLICTING_RESULTS`
- `INDETERMINATE_REDACTION_UNPROVEN`
- `INDETERMINATE_RESOURCE_LIMIT`
- `INDETERMINATE_TIMEOUT`

### Unavailable

- `UNAVAILABLE_ENGINE`
- `UNAVAILABLE_RUNTIME`
- `UNAVAILABLE_REQUIRED_INPUT`
- `UNAVAILABLE_WORKSPACE`

### Terminal

- `TERMINAL_RETRY_EXHAUSTED`
- `TERMINAL_TIMEOUT_EXHAUSTED`
- `TERMINAL_UNRESOLVED_CONFLICT`
- `TERMINAL_INTERNAL_INVARIANT`

Reason codes are content-free and stable within a schema major version. A new reason code requires a compatible schema change; repurposing an existing code is forbidden.

## 13. CLI contract

The initial public interface is one non-interactive CLI. It reads the request from a file descriptor or local file and emits exactly one `scan-outcome` JSON object to stdout.

Exit codes are:

| Code | Meaning |
|---:|---|
| 0 | Pass |
| 10 | Fail |
| 20 | Indeterminate |
| 30 | Unavailable |
| 40 | Terminal/internal invariant failure |

Stderr is content-free and reserved for bounded operational diagnostics. Production mode has no verbose flag, interactive prompt, progress banner, colored output, GitHub annotation, SARIF upload, or finding report export.

The outcome contains:

- Outcome schema version.
- Scan ID and scan mode.
- State and terminal reason code.
- Exact content-free binding digests and versions.
- Attempt number supplied by the project adapter.
- Start/end/duration.
- Coverage classes completed.
- Permitted next action.
- Optional `supersedesScanId` for recovery.

It contains no findings, counts, locations, rule names, provider names, authors, excerpts, or allowlist details.

## 14. User-facing messages and permitted actions

| State | Required message meaning | Developer | Codex/automation | Owner |
|---|---|---|---|---|
| Pass | Scan passed for the recorded bindings; this does not authorize merge or deployment. | Continue normal review. | Verify other gates. | Retain approval authority. |
| Fail | Candidate did not satisfy policy; create a corrected candidate and new scan. | Remediate locally. | Assist through authorized redacted remediation. | Authorize credential response if needed. |
| Indeterminate | Complete trusted coverage could not be proven. | Do not bypass. | Diagnose contract/input cause. | Decide only genuine policy or scope gates. |
| Unavailable | Required scanner/runtime/input was unavailable. | Wait for bounded retry. | Restore availability without weakening controls. | Approve infrastructure or spending if required. |
| Retrying | A clean bounded retry is running. | Wait. | Preserve earlier outcome and new attempt evidence. | No action unless escalation occurs. |
| Terminal | Retries or safety checks are exhausted; progression is blocked. | Stop progression. | Produce a content-free diagnosis. | Decide a separately authorized recovery path. |

No actor may manually convert a non-pass result into pass. Exceptions change versioned policy or allowlist authority and require a new scan.

## 15. Finding lifecycle

1. **Admission:** Engine output enters a private runner channel in memory or memory-backed temporary storage.
2. **Normalization:** The runner validates the engine format and maps it to an internal finding structure.
3. **Deduplication:** Findings are deduplicated only within the current run using an ephemeral per-run key. No reusable finding fingerprint is persisted.
4. **Redaction:** Secret values are removed before any information leaves the private channel.
5. **Policy evaluation:** Scanner safety, mandatory protections, project policy, and eligible exceptions are evaluated in precedence order.
6. **Remediation:** Standard output remains content-free. A separately authorized local remediation session may expose only minimum redacted location metadata.
7. **Credential response:** Suspected real credentials trigger a separate owner-authorized incident task. The scanner does not validate, revoke, rotate, test, or transmit credentials.
8. **Closure:** Only a new clean scan of the corrected exact candidate closes the failure. Manual closure is prohibited.
9. **Recurrence:** A later detection is a new event and receives a new scan ID. No prior dismissal suppresses it unless a still-valid versioned exception applies.
10. **Receipt impact:** No receipt can be issued while a fail, indeterminate, unavailable, or terminal result applies. Credible later evidence of compromise can revoke an earlier receipt.

There is no durable per-finding database. Temporary content-free remediation evidence is deleted seven days after closure and never later than 30 days after creation.

## 16. Policy, severity, and allowlists

Policies and allowlists are project-local, versioned, digest-bound inputs. They must not be stored in the public scanner repository.

Severity vocabulary is `critical`, `high`, `medium`, `low`, and `informational`, but blocking behavior is controlled by rule class and project policy. Credential-like material in mandatory classes is always blocking. Severity cannot downgrade a mandatory block.

An allowlist exception requires:

- Project-local exception ID.
- Named owner and approver.
- Exact rule and minimum possible scope.
- Evidence that the value is synthetic, public, invalid-by-construction, or otherwise not a credential.
- Creation and expiry timestamps.
- Maximum duration of 30 days.
- Review before renewal; renewal creates a new version.
- Automatic invalidation on scope, policy, rule, scanner, source, or evidence change.

Real credentials, signing keys, access tokens, passwords, private keys, production connection material, and unverified “probably inactive” values cannot be allowlisted.

## 17. Schema versioning and compatibility

Schema families evolve independently. Every family uses explicit `major.minor` identifiers.

- Compatible minor additions may be accepted only when unknown fields are safely ignorable and no meaning changes.
- Unknown major versions fail closed.
- Fields cannot change meaning within a major version.
- Deprecated versions remain readable for 90 days or two successful consuming-project release cycles, whichever is longer.
- The compatibility period has a 180-day hard ceiling unless the owner approves an extension.
- Retirement requires a signed notice, replacement version, migration guidance, and exact last-supported date.
- Historical public schemas and verification instructions remain available for the full evidence-retention period.

## 18. Scanner release lifecycle

### 18.1 Upstream intake

The scanner project must never tell consumers to download Gitleaks or TruffleHog directly. Maintainers perform a controlled intake:

1. Select an exact upstream tag and commit after current security, maintenance, and licence review.
2. Verify upstream source and available signatures/checksums.
3. Build in a clean trusted workflow using pinned toolchain and actions.
4. Run unit, fixture, parity, redaction, offline, supply-chain, and adversarial tests.
5. Produce runner and engine binaries.
6. Generate SBOM and complete licence evidence.
7. Generate the release manifest and asset checksums.
8. Sign through the exact approved keyless release identity.
9. Publish all assets together from a draft as one immutable release.

### 18.2 Release bundle

Every release contains:

- Windows amd64 runner and primary engine.
- Linux amd64 runner and primary engine.
- Generic rule pack.
- All external contract schemas.
- Release manifest covering every asset digest.
- SHA-256 checksum file.
- Sigstore/Cosign bundle.
- SBOM.
- MIT runner licence and all third-party notices.
- Compatibility matrix.
- Acceptance-test summary.
- Release notes, known limitations, and retirement status.
- Global revocation location and verification instructions.

TruffleHog is included only in a separately identifiable fallback bundle after legal approval. Its absence cannot be disguised; a required but unavailable fallback produces a non-pass outcome.

### 18.3 Distribution and verification

Consumers use an exact release version and asset digest. `latest`, mutable branch names, moving tags, direct installation scripts, pipe-to-shell installation, and unverified package-manager resolution are prohibited.

The consumer acquisition step may use network access. Before scan execution it must:

- Verify the immutable GitHub release when supported.
- Verify the Cosign bundle against the exact repository, workflow, tag/ref, and GitHub OIDC issuer.
- Verify the release manifest signature and every asset checksum.
- Verify the expected licence and SBOM are present.
- Check global revocation markers.
- Move the verified bundle into content-addressed local custody.

Scan execution begins only after acquisition completes and network access is removed.

### 18.4 Updates, rollback, and retirement

- No automatic updates.
- Upstream monitoring may propose a release but cannot publish one automatically.
- Every new release repeats independent acceptance.
- Rollback is allowed only to a supported, compatible, unrevoked signed release.
- Revocation always wins over pinning or rollback preference.
- A retired release remains verifiable for historical evidence but cannot issue new passing outcomes after its final-use date.

## 19. Signing and trust boundaries

### Public scanner releases

Use Cosign keyless signing from one exact protected GitHub release workflow. Verification policy must match the exact repository, workflow identity, release ref, and OIDC issuer. The Sigstore bundle is stored with the release for offline verification. Public transparency is acceptable because scanner releases contain no project data.

GitHub immutable-release attestations or artifact attestations may be additional evidence but are not the sole trust root. Private-repository artifact attestation features must not become a consuming-project requirement because availability depends on GitHub plan.

### Consuming-project receipts

The scanner release trust root cannot sign project receipts. Each project uses an independent trust root and custodian. RestoHub’s initial recommendation is an owner-managed Cosign key held offline; later KMS/HSM automation requires a separate provider, credential, and spending decision.

Rotation records the old and new key identifiers and preserves old public keys. Compromise revokes affected evidence from the earliest credible compromise time. A replacement key cannot silently rehabilitate old evidence.

## 20. Receipt and revocation integration

The scanner emits a content-free outcome; the consuming project independently verifies it and issues the receipt.

A project receipt should bind at minimum:

- Receipt ID and schema version.
- Issuer and signing-key IDs.
- Scan ID, mode, and pass reason.
- Exact release commit and history range.
- Tracked-tree, build-context, and complete artifact-set digests.
- Build provenance identifier.
- Scanner release, runner, engine, rule-pack, policy, allowlist, and schema digests.
- Outcome digest and timestamp.
- Promotion deadline.
- Evidence-chain sequence and previous-record hash.
- Signature/bundle reference.

One full-release receipt may cover staging and production only when every recorded binding remains identical and production promotion occurs within 30 days. A rebuild, changed artifact, changed configuration contract, changed policy/allowlist, new scanner release, expired exception, revocation, or evidence conflict requires a new scan and receipt.

Authorized revocation sources are:

- Project owner: any project receipt.
- Project incident authority: emergency security or credential event.
- Receipt issuer: issuance or binding defect.
- Signing-key custodian: key compromise or loss of trust.
- Scanner release authority: global scanner-release defect or compromise.

Developers, ordinary automation, and verifiers may recommend rejection but cannot revoke unless explicitly delegated.

Non-technology revocation events include binding error, incomplete coverage discovered later, new credible finding in covered material, policy or exception invalidation, receipt issuance error, trust-authority error, evidence conflict, rollback/truncation, and incorrect reuse across changed staging/production bindings.

Receipts are immutable. Revocation is append-only and never deletes or edits the original receipt.

## 21. Evidence custody

The scanner product must remain custody-neutral but provide schemas and a reference chain verifier. Each project selects its own evidence store.

RestoHub’s initial model is a protected, hash-chained `security-evidence` branch with:

- Monotonic sequence numbers.
- Previous-record hashes.
- Signed receipt and revocation records.
- One restricted writer path.
- Independent staging and production last-seen checkpoints.
- Fail-closed handling of divergent heads, duplicate sequences, missing history, rollback, truncation, or conflicting records.

A stronger future immutable object store may replace active custody only through dual verification, a signed chain-head snapshot, migration markers in both stores, complete reconciliation, and continued historical readability of the old store.

## 22. Cross-project isolation and reporting

Every project has independent policy, allowlists, workspace, cache, temporary storage, receipt keys, receipts, revocations, custody, checkpoints, retention, and audit authority.

The only shareable items are signed scanner releases, generic rules, public schemas, synthetic fixtures, compatibility notices, and global scanner-release revocations.

Project source, findings, paths, repository identities, commits, policies, allowlists, receipt IDs, keys, customer data, and project-level statistics must never cross project boundaries.

Cross-project reporting is absent from the initial product. A future reporting extension requires:

- Explicit signed opt-in from each project owner.
- Immediate stop of future contribution after withdrawal.
- At least five participating projects for every reported aggregate cell.
- No project identity, path, detector, rule, receipt, policy, exception, or finding data.
- No authority to change any project’s policy.
- A provisional maximum 24-month aggregate retention, subject to legal review.

## 23. Security and threat model

The design must defend against:

- Malicious candidate source or filenames.
- Shell and workflow injection.
- Candidate code execution.
- Malicious Git configuration, hooks, attributes, filters, or external diff/text converters.
- Path traversal, symlink escape, hard links, reparse points, device files, and case collisions.
- Archive bombs and malformed containers.
- Resource exhaustion and intentional timeout.
- Engine output containing raw credentials.
- Output, log, exception, and annotation injection.
- Compromised or substituted engine/scanner assets.
- Mutable tags, rollback, repository resurrection, and evidence truncation.
- Cross-project cache or workspace contamination.
- A scanner binary attempting provider verification or exfiltration.
- Secrets available in environment variables, runner memory, Git credentials, Docker configuration, or unrelated workspaces.

Required controls include:

- Fresh workspace per attempt and project.
- Candidate treated only as data.
- No candidate dependency installation or execution.
- Subprocess argument arrays; no shell interpolation.
- Read-only input mounts and read-only scanner root filesystem.
- Network-disabled scanner container.
- Dropped Linux capabilities, no-new-privileges, bounded CPU/memory/process count, and memory-backed temporary output.
- No secrets or write-capable repository token in the scan job.
- Pinned third-party actions by full commit SHA.
- Trusted release and scan workflow definitions.
- Safe Git configuration with hooks, filters, pagers, external diff, text conversion, and inherited config disabled.
- Complete cleanup after each attempt.

## 24. GitHub and CI integration boundary

The scanner CLI is CI-provider neutral. GitHub-specific integration belongs to each consuming project’s adapter/workflow.

Recommended RestoHub future checks are:

- `RestoHub Secret Scan (PR)`
- `RestoHub Secret Scan (Release)`

The PR workflow must use trusted-base control logic, minimal read permissions, no secrets, and candidate-as-data handling. A workflow file alone is not a server-enforced merge barrier; the owner must separately configure required checks and branch protection.

The release workflow separates:

```text
online acquisition and integrity verification
    -> network-disabled scan execution
    -> content-free project verification
    -> separately authorized project receipt signing
    -> append-only evidence publication
```

Receipt signing credentials must never be present in the scanner execution job.

Email, issue creation, chat notifications, SARIF uploads, finding annotations, and dashboards remain disabled initially. Standard outputs are content-free GitHub Checks and content-free local console status.

## 25. Retention requirements

The scanner itself retains no findings. Consuming projects apply these approved defaults:

| Evidence | Retention |
|---|---|
| Temporary remediation evidence | Delete seven days after closure; absolute maximum 30 days. |
| PR content-free evidence | 180 days after PR closure. |
| Terminal/non-pass content-free evidence | 90 days after terminal resolution. |
| Project-local trend evidence | 24 months. |
| Cross-project aggregate, if ever approved | Provisional maximum 24 months. |
| Release receipts and revocation evidence | While release remains deployable/supported plus 24 months, with a three-year minimum. |

The release-receipt schedule is provisional pending applicable legal/compliance review. A legal hold or owner-approved longer requirement may extend retention; automatic shortening below the minimum is prohibited.

## 26. Performance and reliability requirements

- PR scan timeout: 15 minutes per attempt.
- Full-release scan timeout: 60 minutes per attempt.
- Attempts: one initial attempt plus two retries for eligible transient outcomes.
- The runner must stop promptly on cancellation and still perform content-free cleanup.
- The runner must not silently omit files to meet a timeout.
- The test suite must include small, medium, and declared maximum fixture profiles.
- Performance acceptance is measured on the documented standard GitHub-hosted Ubuntu runner class used by the release.
- Projects exceeding declared limits receive indeterminate and must explicitly revise policy or execution capacity.
- Windows and Linux outcomes must be semantically identical for the same logical fixture.

## 27. Acceptance and test matrix

| Area | Required test | Expected result |
|---|---|---|
| Clean source | Known-safe fixtures | Pass |
| Detection | Realistic synthetic secret fixtures | Fail without disclosure |
| History | Secret added then deleted within range | Fail |
| Range binding | Finding outside declared range | Not attributed; exact range remains recorded |
| Tracked source | Manifest and tree coverage | Complete or non-pass |
| Build context | Generated/untracked build inputs | Complete or non-pass |
| Artifacts | Deployment files and image layers | Complete or non-pass |
| Archives | Supported nested archive | Detect synthetic finding |
| Unsupported archive | Encrypted/malformed/unsafe archive | Indeterminate |
| Limits | Oversized/decompression-bomb fixture | Indeterminate, never skip/pass |
| Redaction | All standard and crash outputs | No secret or forbidden metadata |
| Output injection | Malicious filenames/messages | Content-free valid output |
| Offline | DNS/TCP canary and verification-capable engine | No network; deterministic scan |
| Integrity | Changed engine/runner/rule asset | Fail before scan |
| Signing | Wrong identity, issuer, tag, or bundle | Reject release |
| Schema | Unknown major/request corruption | Indeterminate/reject |
| Policy | Invalid signature/digest/precedence | Fail closed |
| Allowlist | Expired, broadened, or real credential | Reject exception and fail |
| Retry | Eligible transient failure | Two retries then terminal |
| Deterministic fail | Confirmed synthetic finding | No retry |
| Engine conflict | Required engines disagree or one incomplete | Fail if either finds; otherwise indeterminate |
| Cleanup | Successful, failed, timeout, and crash runs | No transient material remains |
| Receipt reuse | All recorded staging/production bindings identical | Reuse allowed within 30 days |
| Receipt invalidation | Any binding changes | Reject receipt |
| Revocation | Scanner, key, receipt, policy, or exception marker | Reject affected evidence |
| Evidence rollback | Truncated/divergent/conflicting chain | Block progression |
| Isolation | Parallel projects and reused runner host | No contamination |
| Untrusted PR | Candidate attempts execution/workflow injection | Candidate remains data only |
| Permissions | Extra token rights or supplied secrets | Workflow-policy failure |
| Windows/Linux | Identical logical fixture | Equivalent outcome |
| Performance | Approved fixture profiles | Complete within mode timeout |

## 28. Required release documentation

Every accepted scanner release must document:

- Supported operating systems and architectures.
- Exact runner, engine, rule-pack, and schema versions.
- Supported Git, file, archive, and image input classes.
- Known unsupported inputs and fail-closed behavior.
- Complete CLI and JSON contract.
- Offline sandbox requirements.
- Supply-chain verification procedure.
- Consumer integration sequence.
- Policy/allowlist separation.
- Redaction and remediation boundary.
- Update, rollback, retirement, and revocation procedure.
- Licence and SBOM evidence.
- Security contact and content-free vulnerability-reporting process.

## 29. Non-goals

- GitHub native secret scanning or GitHub Secret Protection replacement features.
- Hosted scanning, cloud source upload, provider verification, or credential testing.
- Real credential storage, rotation, revocation, validation, or incident response.
- Project deployment, merge approval, production approval, compliance certification, or go-live authority.
- Shared cross-project policy, allowlists, receipt keys, receipts, findings, or caches.
- Durable per-finding database.
- Initial GUI, web server, dashboard, email, issue, chat, or SARIF reporting.
- General vulnerability, dependency, malware, SAST, DAST, or infrastructure scanning.
- Automatic scanner updates or direct upstream consumption.
- Automatic exception approval.
- Production KMS/HSM selection.

## 30. Implementation sequence

Each task must be independently bounded, verified, accepted, and closed before a successor is selected. Proposed IDs are planning labels until the new project records them:

1. `PSCAN-01` — Repository governance, canonical specification, threat model, licence plan, and release authority.
2. `PSCAN-02` — External schemas, reason codes, CLI skeleton, deterministic workspace, and content-free outcome contract.
3. `PSCAN-03` — Controlled Gitleaks source intake, build, adapter, Git/file scanning, and fixture coverage.
4. `PSCAN-04` — Artifact normalizer, OCI/archive coverage, redaction firewall, cleanup, and network-disabled execution.
5. `PSCAN-05` — Policy precedence, allowlists, schema compatibility, and project-isolation controls.
6. `PSCAN-06` — Reproducible cross-platform build, SBOM/licence bundle, Cosign keyless signing, immutable release, and verifier guidance.
7. `PSCAN-07` — Full adversarial acceptance, Windows/Linux parity, performance, revocation, retirement, and first signed v1 release.
8. `PSCAN-08` — Optional TruffleHog fallback integration only after legal approval and a proven coverage requirement.

RestoHub integration starts only after `PSCAN-07` produces an independently accepted signed release. RestoHub’s broad T11.1A work should then be split into its own project adapter, PR gate, release/receipt gate, and end-to-end acceptance tasks.

## 31. Definition of done

The product is ready for its first consuming project only when:

- All CAP-1 through CAP-20 success criteria pass.
- The full acceptance matrix passes on Windows amd64 and Linux amd64.
- The release contains no project-specific data or real credentials.
- Source, runner, engine, rules, schemas, SBOM, licences, checksums, and signatures are mutually bound.
- The exact release is immutable and verifiable offline.
- The scanner completes authoritative tests with network disabled.
- No raw or partially raw finding reaches any external output.
- No unsupported input or skipped content can produce pass.
- Global revocation and retirement procedures are tested.
- A sample project adapter independently verifies a full-release pass and issues a mock project-owned receipt without giving its key to the scanner.
- Independent review confirms the scanner has no merge, deployment, business-policy, or credential-response authority.

The world-change success signal is: RestoHub and a second unrelated fixture project can independently consume the same exact signed scanner release, keep all project data isolated, detect synthetic secrets in PR and complete release inputs, emit only content-free outcomes, and block promotion whenever trusted coverage cannot be proven.

## 32. Owner gates that remain operationally separate

The design direction is settled, but these actions still require explicit owner authority at the time they occur:

- Create the new repository and select its final name/account.
- Make the repository public or publish any release.
- Change GitHub Actions, branch, release, environment, or protection settings.
- Approve current GitHub plan availability or paid usage.
- Establish release maintainers and the exact keyless workflow identity.
- Accept TruffleHog AGPL obligations and fallback distribution.
- Create or handle any signing credential or consuming-project receipt key.
- Select KMS/HSM or another paid custody provider.
- Enable cross-project reporting.
- Adopt the scanner in a consuming project.
- Configure required checks, merge barriers, deployment gates, staging, production, or go-live.
- Approve the final legal/compliance retention schedule.

## 33. Open implementation facts to verify at project start

- Current supported Gitleaks release, source revision, licence, security posture, and reproducible build requirements.
- Current TruffleHog release and legal approval status before enabling fallback.
- Current Go stable/toolchain support policy.
- Current Cosign/Sigstore bundle format and offline-verification procedure.
- Current GitHub immutable-release, artifact-attestation, Actions, and plan behavior for the selected account.
- Exact public repository name, owner account, protected release workflow identity, and release approvers.

These are version-sensitive preflight facts, not permission to change the design silently. A material incompatibility must return to the owner with alternatives and a recommendation.

## 34. Reference sources

Product implementation must revalidate current primary sources before pinning a version:

- Gitleaks repository and CLI documentation: <https://github.com/gitleaks/gitleaks>
- Gitleaks releases: <https://github.com/gitleaks/gitleaks/releases>
- TruffleHog repository, offline flag, licence, and signed release evidence: <https://github.com/trufflesecurity/trufflehog>
- Sigstore Cosign blob signing: <https://docs.sigstore.dev/cosign/signing/signing_with_blobs/>
- Sigstore Cosign verification: <https://docs.sigstore.dev/cosign/verifying/verify/>
- GitHub immutable releases: <https://docs.github.com/en/code-security/concepts/supply-chain-security/immutable-releases>
- GitHub artifact attestations: <https://docs.github.com/en/actions/concepts/security/artifact-attestations>
- GitHub Actions secure-use guidance: <https://docs.github.com/en/actions/reference/security/secure-use>

RestoHub planning provenance is fully absorbed into this document from its accepted T11.1A tracker specification, acceptance evidence, Direction Guard workflow, and owner-approved scanner discussions. Those RestoHub artifacts remain RestoHub authority; they are not dependencies of the new scanner repository.

## 35. Exact prompt for the new project

Copy this prompt into the new scanner-project session and attach this Markdown document:

```text
You are working in the new project for the Project-Agnostic Secret Scanner.

Read the attached “PROJECT-AGNOSTIC-SECRET-SCANNER-DESIGN-REQUIREMENTS.md” completely and treat it as the controlling product contract. Preserve every capability, constraint, non-goal, security boundary, schema boundary, acceptance test, owner gate, and definition-of-done requirement. You may add stricter controls, but you must not weaken, reinterpret, omit, or silently replace any requirement.

The product must remain independent of RestoHub and every other consuming project. RestoHub is the first consumer, not product authority. Do not copy RestoHub source, policies, allowlists, receipts, keys, findings, credentials, customer data, repository identity, or project-specific behavior into this scanner project.

Use Gitleaks CLI as the primary engine. Use a thin MIT-licensed Go runner. Keep TruffleHog OSS disabled by default and do not bundle or enable it until the owner separately approves its legal obligations and a proven coverage gap. Do not use GitHub native secret scanning or any hosted scanning service. Scanner execution must be network-disabled, findings must be transient and fully redacted, and every incomplete or untrustworthy condition must fail closed.

Start by inspecting the new repository’s live authority, Git state, available GitHub/account facts, and existing work. Preserve all existing user changes. Revalidate current official primary sources for Gitleaks, TruffleHog, Go, Cosign/Sigstore, GitHub Releases, immutable releases, artifact attestations, Actions security, licensing, and any plan/cost dependency. Clearly distinguish verified current facts, recommendations, assumptions, and owner-reserved decisions.

Then produce a requirement-preservation matrix mapping every CAP-1 through CAP-20 requirement and every acceptance-matrix row to the proposed architecture, file boundary, test, and implementation task. Identify contradictions or missing prerequisites and resolve ordinary technical choices yourself. Return one consolidated set of genuine owner gates rather than a sequence of small questions.

Adopt the bounded implementation sequence PSCAN-01 through PSCAN-08 from the specification. Work on exactly one activated task per session. No successor activates automatically. The first session must prepare the project-grounded architecture, repository layout, task specifications, threat model, validation plan, and proposed PSCAN-01 activation boundary; it must not implement product code, create credentials, publish the repository, change GitHub settings, incur spending, enable TruffleHog, or perform remote actions unless the owner explicitly authorizes those actions.

For every later implementation task: keep changes inside its approved paths, use pinned dependencies and actions, add the required unit/integration/adversarial tests, independently verify the result, record exact evidence, preserve unrelated work, and stop after that task’s acceptance without activating its successor.

At the end of this first session, provide:
1. The requirement-preservation matrix.
2. The complete proposed architecture and source layout.
3. The exact PSCAN-01 scope, allowed paths, forbidden paths, checks, exclusions, and success criteria.
4. All current account, legal, publication, credential, spending, signing, and GitHub-setting gates in one owner decision block.
5. A readiness verdict.
6. A copy-ready activation prompt only if the new project’s governance makes PSCAN-01 eligible; label it “not executed.”

Do not create or publish a repository, create or handle signing keys, enable paid services, change GitHub settings, download or execute a scanner against real project data, handle real credentials, perform remote actions, or activate a task during this initial design session.
```
