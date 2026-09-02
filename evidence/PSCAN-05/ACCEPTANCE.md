# PSCAN-05 Acceptance Evidence

## Result

`ACCEPT`. PSCAN-05 satisfies its activated policy projection, compatibility,
receipt/revocation reference and project-isolation scope after the six bounded
review corrections recorded in `REVIEW.md`. Acceptance is local and applies
only to PSCAN-05. No successor is selected or activated.

## Contract mapping

| Contract rows | Accepted evidence |
|---|---|
| CAP-9, AT-17 | Exact-byte digest plus public-key signature projection, independent family windows, strict current semantics and fixed scanner/global/project/exception precedence. |
| CAP-10, AT-18 | Unique owner-approved exception identities, exact scope/evidence/context, complete invalidation set, use-time expiry, project and 30-day maxima, and mandatory credential blocking. |
| CAP-15, AT-16 | Independent major/minor windows reject unknown, ambiguous or required future semantics and enforce signed 90-day/two-cycle retirement overlap with the 180-day ceiling. |
| CAP-16, AT-23, AT-24 | Custody-neutral receipt verification binds every predecessor and new proof fact, signed projection identities, schema digests, outcome time, exact 30-day deadline, previous receipt anchor and current head. |
| CAP-17, AT-25, AT-26 | Full genesis-to-expected-head global/project revocation chains, exact scanner-owned global schema authority, effective revocation precedence, and rollback/gap/divergence/duplicate/conflict rejection. |
| CAP-18, AT-27 | Per-call projections and chain results plus fresh parallel/sequential workspaces expose no shared policy/evidence state; the boundary defines no project cache, store, key or retained instance. |

## Accepted boundary

- Scanner-owned request 1.1 binds project projection and outcome-proof inputs;
  global-revocation 1.1 defines signed append-only scanner authority records.
- Project policy, allowlist, receipt and revocation formats remain project-owned
  caller inputs. The scanner stores no instance and owns no receipt schema.
- Ed25519 verification is public-key only. The API has no generation, signing,
  issuance, private-key, evidence-store, promotion or retention operation.
- Exact document, receipt, chain and exception-set projections are
  domain-separated and length-prefixed; changed binding facts fail closed.
- Parallel, sequential and native Windows/Linux validation passed for the
  complete PSCAN-05 package set.

## Evidence chain and limitations

- Activation: `ACTIVATION.md`.
- Author preflight/claims: `PREFLIGHT.md`, `IMPLEMENTATION.md`, `VALIDATION.md`.
- Independent corrections and reproduced validation: `REVIEW.md`.
- Closure and exclusions: `CLOSEOUT.md`.

The full repository WSL run retains two predecessor environment-specific
failures described in the review; every PSCAN-05 package passed, and native
Windows PSCAN-05 execution passed. This bounded acceptance does not establish
whole-product readiness or PSCAN-07 parity/performance/release acceptance.

No remote, push, publication, release, settings, signing, credential, spending,
provider, TruffleHog, consuming-project, deployment, legal, production or
go-live action occurred. Every such gate remains owner-controlled.
