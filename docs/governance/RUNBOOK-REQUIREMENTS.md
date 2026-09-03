# Final Runbook Requirements

## Audience and format

PSCAN-06 and PSCAN-07 must deliver a plain-English Markdown runbook suitable for
maintainers, consuming-project adapter authors, reviewers and auditors. It must
include readable Mermaid diagrams for acquisition, scan execution, state/retry,
release, verification, revocation, rollback and incident escalation. Diagrams
supplement rather than replace exact text.

## Required input documentation

Document every CLI flag, environment variable, file descriptor, file path and
JSON request field, including:

- name, type, required/optional status and allowed values;
- mode applicability (`local`, `pr`, `release`);
- schema family and exact `major.minor` compatibility;
- freshness, size, path and normalization constraints;
- exact digest/version/trust binding semantics;
- default, minimum, maximum and who may change it;
- sensitivity and whether it may cross the scanner boundary;
- validation order and fail-closed reason for missing, duplicate, unsafe,
  conflicting or unsupported values;
- platform-specific representation where mechanics differ;
- synthetic valid and invalid examples.

The complete request payload documentation must cover `requestSchemaVersion`,
`scanId`, `mode`, scanner release, engine, rule, policy and allowlist bindings,
source binding, tracked-source manifest, release build-context and artifact
manifests, fallback requirement, limits, `offlineRequired`, `redactionMode`,
`requestedAt` and future compatible-field behavior.

It must also document the exact-object admission-ledger, raw-classifier,
preparation/projection, declared-profile and detector-inspection-proof bindings
introduced by PASS-OUTCOME-SPEC-001, including which conditions are non-pass.

## Required output documentation

Document stdout, stderr and process exit behavior. Stdout emits exactly one
`scan-outcome` JSON object. Document each allowed field, type, state
applicability, binding semantics and forbidden content. Cover outcome schema,
scan ID/mode, state, reason code, bindings/versions, adapter-supplied attempt,
start/end/duration, completed coverage classes, permitted action and optional
recovery reference.

Document exit mapping: 0 pass, 10 fail, 20 indeterminate, 30 unavailable and 40
terminal/internal invariant. Stderr is bounded and content-free. Explicitly
list forbidden findings, counts, locations, paths, detector/provider names,
authors, excerpts, hashes/fingerprints, allowlist details, annotations, SARIF,
verbose/progress/color and interactive output.

## Operational procedures

The runbook must include exact prerequisites, version/digest verification,
online acquisition versus offline execution separation, safe workspace/mount/
permission/network configuration, PR and release coverage, cancellation,
cleanup, retry eligibility, terminal escalation, content-free remediation,
release verification, update, compatible rollback, retirement, revocation,
evidence verification and project isolation.

Every procedure states inputs, outputs, authority, preconditions, postconditions,
failure modes, evidence produced, cleanup and rollback. Copyable examples use
synthetic values and exact pinned versions; they never use `latest`, pipe-to-
shell, real credentials or project-owned data.

## Authority warnings

The runbook must repeatedly distinguish product evidence from project
authority: pass does not approve merge, deployment, production, legal
compliance or go-live. It must state that projects own policies, allowlists,
receipts, revocations, keys, custody, deployment gates and retention.

## Verification gate

PSCAN-07 acceptance confirms every documented input/output matches the released
schemas and CLI, every command works on Windows amd64 and Linux amd64 as stated,
every Mermaid diagram renders, every link resolves, and no example exposes or
normalizes unsafe behavior.

## PSCAN-06 release verification supplement

`docs/release/OFFLINE-VERIFICATION-RUNBOOK.md` documents every release-verifier
parameter, stdout/exit behavior, release-manifest 1.1 field and asset payload;
Windows and Linux examples; fixed identity/tool digests; online quarantine and
offline custody; revocation refresh, rollback, retirement and operator
evidence; and Mermaid trust-flow and state diagrams. Remote creation, settings,
workflow, signing, draft and publication actions remain separately owner-gated.
