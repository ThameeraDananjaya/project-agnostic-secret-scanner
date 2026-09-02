# PSCAN-05 Implementation Record

## Sole-task claim

PSCAN-05 alone is claimed on `2026-09-02` from exact clean activation commit
`50b418609c1f9927c0ecd5d740aba6d7bff11f58`, after the complete ordered
reading and live primary-source preflight in `PREFLIGHT.md`.

PASS-OUTCOME-SPEC-001 remains the controlling contract. PSCAN-06, PSCAN-07,
PSCAN-08 and every other successor remain unselected, inactive, unclaimed and
unimplemented.

## Activated boundary

Implementation is restricted to the allowed paths in `docs/tasks/PSCAN-05.md`.
The forbidden specifications, traceability, historical task/evidence records,
engine/rule/artifact/redaction/cleanup code, public runner, workflows, release,
remote, provider, credential, signing, receipt-issuance and consuming-project
boundaries remain untouched.

## Initial design decision

- Project policy and allowlist JSON remain transient caller inputs. The scanner
  verifies exact-byte digests, declared independent schema versions, detached
  signatures and a minimal versioned internal projection; it does not publish
  or retain a project schema or instance.
- Signatures cover explicit domain-separated binary preimages, not an
  implementation-defined JSON serialization. Only public verification keys are
  accepted. No signing API or private-key type enters product code.
- Precedence is fixed in code: scanner mandatory protections, global
  revocation/protection facts, project policy projection, then eligible narrow
  exceptions. A lower layer can only add denial.
- Receipt and revocation verification operates over caller-owned in-memory
  references and checkpoints. It returns a decision and retains no evidence,
  receipt, project identity or trust key.
- All caches and correlation state are instance-local or attempt-local; no
  package-global project state is introduced.

## Implemented candidate

- Scanner-owned request schema 1.1 carries exact policy/allowlist signed
  bindings plus admission-ledger, raw-classifier, preparation,
  resource-profile and inspection-proof bindings while retaining 1.0 input
  compatibility.
- The policy package rejects duplicate keys, unknown required semantics,
  unsupported family windows, bad signatures/digests/adapters, ambiguous
  versions and malformed projections. Exceptions bind distinct owner and
  approver, rationale, exact rule/class/scope/evidence, scanner/rule/policy/
  source context, creation/expiry and the complete invalidation set.
- Evaluation rechecks expiry and project-specific duration at use time.
  Scanner credential classification and known credential classes are mandatory
  blocks; revocations precede project denial and a valid exact exception is the
  only final-layer suppression.
- The verification package accepts public keys only, validates exact signed
  document messages, independent family windows and signed retirement notices,
  and rejects ambiguous minor/retirement declarations.
- Receipt reuse binds the full predecessor and new proof identity, a maximum
  30-day deadline and an exact trusted receipt head. Full global/project
  revocation chains start at genesis and must reach the declared trusted heads;
  rollback, truncation, sequence gaps, previous-hash divergence, duplicate
  records/targets, conflicts and effective revocations reject.
- Signed global revocation schema 1.1 preserves the scanner-owned family.
  Project receipt, allowlist, policy, evidence-chain and revocation schemas and
  instances remain outside this repository.
- Synthetic adversarial fixtures and tests cover policy, allowlist, evidence
  and parallel/sequential/reused-host isolation without any consuming-project
  identity, secret, credential, receipt, policy instance or private key.

Validation results and limitations are recorded in `VALIDATION.md`. This
implementation claim is not independent acceptance and selects no successor.
