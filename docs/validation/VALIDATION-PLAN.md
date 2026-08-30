# Validation Plan

## Validation principles

Validation is evidence-based, deterministic where possible, adversarial,
path-bounded and independent of implementation authorship. No single unit test,
scanner result, workflow status or signature proves the product complete.
Missing, stale, skipped, unsupported, conflicting or untrusted evidence is a
failure to accept.

All fixtures contain synthetic non-credentials. Tests must never scan real
project data or depend on online provider verification.

## Validation layers

| Layer | Purpose | Minimum evidence |
|---|---|---|
| Contract | Preserve every CAP, acceptance row, non-goal, boundary and gate | Canonical hash plus complete traceability audit |
| Static | Reject project identity, unsafe commands, forbidden outputs, unpinned references and authority drift | Path/content scans and schema validation |
| Unit | Prove validation, normalization, classification, redaction and cleanup invariants | Deterministic component tests |
| Integration | Prove trusted component boundaries and engine behavior | Generated fixtures with exact bindings |
| Adversarial | Challenge injection, traversal, bombs, leakage, conflicts, rollback and contamination | Negative-test corpus with non-pass assertions |
| Supply chain | Prove source-to-asset, manifest, licence, SBOM, identity and digest binding | Offline verifier evidence |
| Platform | Prove logical parity and supported execution | Windows amd64 and Linux amd64 results |
| Release | Prove immutable exact bundle and full acceptance | Signed release evidence and independent review |

## Required fixture families

Clean, synthetic findings, deleted history, range edges, large and binary files,
nested and malformed archives, archive traversal, symlinks/reparse points, case
collisions, OCI/Docker layers, output injection, network attempts, corrupted
assets/signatures, invalid policy and allowlists, retries, engine conflict,
cleanup lifecycle, receipt reuse/invalidation, revocation, evidence rollback,
parallel projects, hostile Git configuration, Windows paths and Linux paths.

Fixture generation must be reproducible, documented and credential-free.

## Mode and resource profiles

- PR: 15 minutes per attempt; tracked tree and exact range; no candidate build.
- Release: 60 minutes per attempt; frozen commit, history, build context,
  artifacts, images and supported nested archives.
- Retry: one initial attempt plus two fresh-workspace retries only for eligible
  transient indeterminate/unavailable states.
- Artifact defaults: depth 5; 100,000 entries; expanded bytes 2 GiB PR and
  10 GiB release; individual regular file 512 MiB; ratio 1,000:1.
- Performance profiles: small, medium and declared maximum on the documented
  standard authoritative runner class.

Crossing a limit is explicit indeterminate; files are never omitted to pass or
meet a timeout.

## Redaction oracle

Each adversarial run plants unique synthetic canaries in candidate values,
paths, metadata, encoded forms and engine crash output. Validation inspects
stdout, stderr, JSON, logs, summaries, exceptions, annotations, reports,
temporary storage and crash residue. The run is acceptable only when no value,
partial value, encoding, hash, reversible fingerprint, path, location,
detector/provider name, author, message, excerpt, count or allowlist detail
crosses the standard boundary.

## Network and permission proof

Authoritative execution runs after acquisition in a network-disabled sandbox
with no secrets or write-capable repository token, read-only inputs and scanner
root, dropped Linux capabilities, no-new-privileges, bounded resources and
memory-backed private output. DNS and TCP canaries plus a verification-capable
engine prove access is unavailable. Extra token rights or supplied secrets are
a workflow-policy failure.

## Determinism and parity

Repeated runs with identical validated inputs must produce semantically
identical public state, reason, bindings and coverage. Timestamps, duration and
random non-source-derived correlation IDs are excluded. Windows amd64 and Linux
amd64 must agree on the logical fixture outcome and forbidden-output contract.

## Supply-chain proof

Before execution, an offline-capable verifier checks exact release version,
manifest signature, every asset digest, repository/workflow/ref/issuer identity,
Cosign bundle, licences, SBOM, compatibility, retirement and global revocation.
Mutation, wrong identity, mutable reference, missing notice or revoked asset is
rejected before scanning.

## Independent acceptance

For each task, the reviewer must read the activated specification and governing
contract, inspect only the accepted diff and evidence, reproduce required
checks, confirm allowed paths and exclusions, and record limitations and owner
gates separately. The reviewer must not infer pass from author claims. PSCAN-07
requires complete fresh acceptance on both supported platforms before first
consumer readiness.

## PSCAN-01 documentation checks

PSCAN-01 acceptance performs read-only checks for:

1. activation commit and pre-change clean state;
2. canonical byte count, 977-line count and SHA-256 equality;
3. CAP-1 through CAP-20 present exactly once in traceability;
4. all 31 acceptance rows present exactly once in traceability;
5. every trace row maps architecture/file boundary, test and task;
6. only allowed paths changed and no forbidden path exists;
7. no product code, schema, workflow, binary, dependency or action introduced;
8. no secret-like or credential file introduced;
9. no consuming-project artifact or behavior introduced beyond the unchanged
   canonical owner contract;
10. document links and local references resolve;
11. MIT and third-party licence boundary is explicit;
12. every successor remains proposed/unselected and PSCAN-08 remains inactive;
13. exact acceptance evidence and remaining owner gates are recorded;
14. accepted local commit leaves the worktree clean.

The reproducible command evidence and results are recorded in
`evidence/PSCAN-01/ACCEPTANCE.md`; command output itself must remain
content-free and must not expose sensitive material.
