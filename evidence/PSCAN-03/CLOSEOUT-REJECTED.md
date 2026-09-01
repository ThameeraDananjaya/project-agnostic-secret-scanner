# PSCAN-03 Rejected Closeout

## Result

PSCAN-03 was closed fail-closed as `REJECTED` on 2026-09-01. It is not accepted.
The correction C1 candidate remains non-authoritative and does not establish a
usable scanner contract. No successor was selected, activated, claimed or
implemented.

## Bound history

- Activation commit: `7459ad43313002a70d3a082c0b220012bad72870`.
- Material-gap commit: `a8faad04d6b6ec255a5c6222ec27841e101049f9`.
- Rejected correction candidate: `6cd22a3615b21e5855f53af8c1d40cb51815ab27`.
- Independent rejection record commit:
  `a9035c5065f7e6738df9401cd5cd1c26b4c3462f`.
- Historical material-gap evidence SHA-256:
  `e1b4dd6140e63c7d5ba715f00b16de342e7dec8bf353d1d923fc0993b5e9a28c`.
- Historical material-gap Git blob:
  `15d91e1131ea1eaa544eceba10c19942f29632db`.

The material-gap record was not rewritten. The blocking facts, exact pins,
commands and adversarial reproduction remain in `REVIEW-C1.md` and
`MATERIAL-GAP-001.md`.

## Rejection basis

Independent review proved that the admitted-blob prefix marker could pass even
when a later Gitleaks fragment missed a synthetic finding. It also found
unbounded whole-match rules, incomplete archive classification, insufficiently
isolated binary proof, unbound resource limits and a conflict between the
projection design and PASS-SPEC-001 section 7.2's mandatory Gitleaks Git mode.
Those conditions make complete detector coverage unprovable under the current
contract and therefore fail closed.

## Owner direction and authority boundary

On 2026-09-01 the owner abandoned the contract-preserving option and directed
the project to decommission PASS-SPEC-001 and focus on the outcome-based option.
That decision establishes the next direction but does not authorize PSCAN-03 to
edit `AGENTS.md`, `CONSTITUTION.md` or `docs/spec/**`, which are outside its
allowed paths. PASS-SPEC-001 therefore remains controlling until a separate
activated governance task decommissions its authority while preserving it as
historical audit evidence.

Exactly one bounded successor candidate is proposed:

- PSCAN-09: preserve PASS-SPEC-001 and its source record as immutable history;
  decommission its controlling authority; establish an outcome-based successor
  contract permitting exact Git-object enumeration and deterministic,
  coverage-proved scan projections; reconcile governance, architecture,
  validation traceability and the downstream task sequence; perform no scanner
  implementation.

PSCAN-09 is proposed and unselected. PSCAN-04 remains unselected. PSCAN-08
remains inactive, unselected and unauthorized.

## Closeout validation

- Changed paths are restricted to PSCAN-03's existing allowlist:
  `docs/tasks/PSCAN-03.md`, `docs/tasks/TRACKER.md` and this evidence file.
- `git diff --check` passed before the closeout commit.
- Exact historical evidence hash and blob checks matched the values above.
- No scanner execution was needed for this documentation-only closeout. The
  rejected correction's pinned Docker validation does not override the
  independent rejection.
- Pins introduced by this closeout: none.

No TruffleHog assessment, download, integration or enablement occurred. No
credential, remote, publication, GitHub-setting, signing, spending, deployment,
archive/OCI normalization or consuming-project action occurred.
