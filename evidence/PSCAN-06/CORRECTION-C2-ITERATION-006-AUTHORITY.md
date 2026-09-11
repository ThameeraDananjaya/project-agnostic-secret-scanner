# PSCAN-06 Correction C2 iteration 006 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Iteration: `006`
- Title: PowerShell PID-collision correction from Recovery R5 failure 001
- Event: exact owner approval of the bounded correction described below
- Owner command:
  `APPROVE PSCAN-06 C2 ITERATION-006 POWERSHELL-PID-COLLISION CORRECTION FROM RECOVERY-R5 FAILURE 001 SHA256 53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0 AT EVIDENCE COMMIT DA45BB19480663AA21E9A67754AD509834DCAED7 LOCAL-ONLY NO-DOCKER NO-NETWORK NO-REMOTE-MUTATION NO-TAG-CHANGE NO-WORKFLOW-DISPATCH NO-SIGNING NO-PUBLICATION NO-SUCCESSOR`
- Date: `2026-09-11`
- Approval session: `01a08fae-2ee1-7553-a126-e635b70a34a6`
- Branch: `main`
- Required starting HEAD before this authority bundle:
  `da45bb19480663aa21e9a67754ad509834dcaed7`
- Required starting tree:
  `2b3ea12ab339f255136dc9a319f5895c8c074e3c`
- Required starting HEAD direct parent:
  `fdd105492011e7acd9d05b63378642ec54a03289`
- Correction authority commit: this authority-bundle commit
- Pre-change tracked, staged and non-ignored untracked Git status: clean
- Preserved ignored workspace material: `graphify-out/**`; excluded from product
  authority and not staged or included in this bundle
- State after commit: owner-approved; unclaimed; not implemented; not accepted
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

This is a bounded iteration inside the already-activated PSCAN-06 Correction
C2. It does not create, select or activate another task. This approval session
records authority only and does not repair the defect. A genuinely fresh
session with an ID different from the approval session must claim only
Correction C2 iteration 006 from this committed authority bundle before any
implementation.

## Controlling and preserved authority

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Original Correction C2 authority:
  `evidence/PSCAN-06/CORRECTION-C2-AUTHORITY.md`
- Original authority SHA-256:
  `3CE8E20151993A8FC77890E51C5BF1557F68040CAFF2C3DB1A46237792A75D63`
- Iteration-005 authority:
  `evidence/PSCAN-06/CORRECTION-C2-ITERATION-005-AUTHORITY.md`
- Iteration-005 authority SHA-256:
  `F866EDEDD850DE754C533FB3D9AC280D20C19CE72B9E45C0744FF9ACF64E49A9`
- Iteration-005 local acceptance:
  `evidence/PSCAN-06/CORRECTION-C2-ITERATION-005-ACCEPTANCE.md`
- Iteration-005 local acceptance SHA-256:
  `1C9392E1899D1EB83C2DBDD673C0369042C5F6255D3D0BE9558C55281EDB72A9`
- Recovery R5 terminal failure:
  `evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R5-FAILURE-001.md`
- Recovery R5 failure SHA-256:
  `53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0`
- Recovery R5 evidence commit:
  `da45bb19480663aa21e9a67754ad509834dcaed7`
- Accepted local tooling candidate:
  `faef8435322c9096df09b56969662411f17356ea`
- Accepted local tooling candidate tree:
  `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`
- Locked product tag: `v1.0.0`
- Locked product-source commit:
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Locked Correction C2 tooling tag:
  `release-tooling-v1.0.0-c2`
- Recovery R5 workflow run: `34582887399`, attempt `1`, conclusion `failure`

Recovery R5 and all earlier authorities, candidates, evidence, tags, schemas,
rulesets and remote facts remain immutable history. This authority does not
reinterpret Recovery R5 as a pass, restore its consumed push or dispatch, or
authorize a second proof-gate attempt. The remote tag and ruleset facts above
are inherited only from the immutable R5 record; this local-only approval
session performs no network or remote read-back.

## Demonstrated defect

The Recovery R5 `ubuntu-24.04` build job passed its first six steps and failed
in the clean isolated no-Docker native matrix. In
`build/release/docker-execution.ps1`, function `Get-LinuxSessionMembers`
assigns a parsed `/proc/<id>/stat` process identifier to lowercase `$pid`.
PowerShell variable names are case-insensitive, so `$pid` resolves to automatic
variable `$PID`, which is read-only. The assignment therefore terminates with:

```text
Cannot overwrite variable PID because it is read-only or constant.
```

The failure occurred before genuine Docker, dependency acquisition, either
build, byte comparison or artifact transfer. It demonstrated only this
PowerShell identifier collision. It did not invalidate the independently
accepted iteration-005 ambient-type isolation or authorize broader process,
Docker, workflow or release changes.

## Approved objective

Remove the case-insensitive collision with PowerShell automatic variable
`$PID` by giving the local parsed Linux process identifier an unambiguous,
non-reserved name and updating only its directly dependent identity and ledger
uses. Preserve the exact `/proc` parsing, session membership filter, start-time
identity, ledger shape and downstream containment behavior.

Add or refine a bounded local no-Docker regression that proves the production
function no longer assigns to `PID` in any case spelling and that the parsed
identifier still feeds the process identity and ledger `PID` field. Rerun the
applicable local no-Docker PowerShell parse, isolation and native-boundary
checks without contacting a network or invoking Docker.

## Allowed implementation paths

Correction C2 iteration 006 may change only:

```text
build/release/docker-execution.ps1
build/release/test-docker-execution.ps1
docs/release/PSCAN-06-RELEASE-PLAN.md
docs/validation/PSCAN-06-VALIDATION.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/**
```

This list supersedes broader C2 and earlier-iteration path lists for iteration
006 only. Every absent path is forbidden. A demonstrated need to change a
workflow, another build or release script, source contract, schema, dependency
or helper stops this iteration and returns to the owner with exact evidence.

## Mandatory correction design

1. In `Get-LinuxSessionMembers`, replace the local lowercase `$pid` binding
   with one explicit non-reserved process-identifier name. Update only the
   exact string identity and ledger-property value that consume that binding.
2. Preserve the parsed value type as `[int]`, the start time as `[uint64]`, the
   identity format `<process-id>:<start-time>`, and the ledger object fields
   `PID` and `StartTime`. Property name `PID` is data and must not be renamed.
3. Do not change `/proc` enumeration, malformed-record handling, session
   membership selection, liveness checks, process signaling, cleanup or time
   bounds to mask the collision.
4. Add a regression over the exact production source proving that
   `Get-LinuxSessionMembers` has no assignment target named `PID` under
   case-insensitive comparison and still binds its parsed process identifier
   into both the identity key and ledger `PID` value.
5. Preserve iteration-005 ambient-type rejection and every iteration-004
   executable-identity, private-environment, Windows job-object, Linux stopped
   PID-namespace-init, stream-limit, time-limit and process-tree property.
6. Preserve the finite nine-operation Docker table, pinned image identity,
   exact release identities, schema `2.1`, earlier schemas and all scanner,
   verifier, policy and consuming-project behavior.
7. Keep the correction local-only. Any result requiring Linux, Docker, a
   workflow, a remote read, dependency acquisition or a network connection is
   explicitly not proved by this iteration.

## Required local checks and success criteria

1. PowerShell parsing succeeds for both allowed implementation scripts.
2. Static inspection and the bounded regression prove that the production
   function assigns no variable named `PID` in any case spelling.
3. The regression proves the renamed parsed identifier remains the source of
   both the `<process-id>:<start-time>` identity and ledger `PID` value.
4. The existing local no-Docker hostile compatible-type, stale-type, clean
   native-boundary, stream, timeout, descendant and cleanup matrix passes on
   the available host without surviving marked processes.
5. Exact source-trust, diff/path confinement and `git diff --check` pass. No
   absent path changes and no ignored `graphify-out/**` material is staged.
6. Independent skeptical review examines the exact bounded diff and reruns
   proportionate local no-Docker checks before any local acceptance record.
7. Recovery R5 remains terminally failed. Actual-Linux, genuine-Docker,
   dependency, build, reproducibility and artifact-integrity proof remains
   open, and no new proof-gate execution is inferred from local success.
8. PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains proposed and
   unselected; PSCAN-08 remains inactive and ineligible.

## Forbidden scope and reserved gates

- No implementation in this approval session.
- No modification of PASS-OUTCOME-SPEC-001, PASS-SPEC-001, their source record,
  TRACEABILITY, DEC-002, existing schemas, scanner/verifier behavior, policy,
  allowlist, receipt, revocation or consuming-project behavior.
- No change to a workflow, image-admission path, Docker operation table,
  acquisition, cache, CRLF, build, comparison, signing or publication logic.
- No Docker command, image pull, image execution, dependency download,
  scanner/toolchain run, build or external network action.
- No remote read or mutation, push, fetch, tag creation, tag movement, tag
  deletion, ruleset or setting action, workflow dispatch, rerun or retry.
- No credential, signing key, paid capability, TruffleHog work, signing,
  attestation, draft, release, deployment or publication action.
- No deletion, movement or staging of preserved ignored `graphify-out/**`
  material.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

Any future actual-Linux, genuine-Docker or workflow proof attempt requires a
separate exact owner decision after this correction is implemented, validated
and independently accepted locally. This authority supplies no such attempt.

## Fresh implementation-session requirements

The implementation session must:

1. have a session ID different from
   `01a08fae-2ee1-7553-a126-e635b70a34a6`;
2. begin on exact branch `main` at this committed authority bundle with a clean
   tracked, staged and non-ignored untracked state;
3. establish exact every-byte source trust or an already-authorized canonical
   whole-file projection before editing;
4. read the controlling contract, PSCAN-06 task, reading map, tracker, this
   authority and the immutable R5 failure completely;
5. claim only PSCAN-06 Correction C2 iteration 006 and record its bounded
   preflight before implementation; and
6. stop fail-closed on any identity, path, scope or authority mismatch.

## Approval-session exclusions

This approval session introduces no production code, test, workflow, Docker
boundary, build, schema, verifier or product implementation. It performs no
Docker command, network request, remote read/write, tag operation, workflow
run, dependency acquisition, build, signing, attestation, draft, release,
publication or successor action. Its only durable change is this bounded
authority bundle and synchronized living PSCAN-06 task state.
