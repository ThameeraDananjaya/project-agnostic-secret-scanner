# PSCAN-06 Correction C2 iteration 007 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Iteration: `007`
- Title: immutable Recovery R6 identity roll-forward
- Event: owner-directed selection of the minimal bounded prerequisite to any
  Recovery R6 dispatch
- Date: `2026-09-11`
- Approval session: `01a09141-1ba5-7931-bedc-bee2402fbf95`
- Branch: `main`
- Required starting HEAD before this authority bundle:
  `5ff14709520ba99b105ee515e5482455c8511efc`
- Required starting tree:
  `dfe8f3f9a5488dffa349fa49169f6acac15c469b`
- Required starting HEAD direct parent:
  `6ae41db6c52c7bb3efd627131dee79f5fb3bc851`
- Correction authority commit: this authority-bundle commit
- Pre-change tracked, staged and non-ignored untracked Git status: clean
- Pre-change graphify-aware exact checkout verification: 373 tracked paths,
  355 raw-equal, 18 canonical whole-file LF-to-CRLF projections, exact index
  and flags, zero unsupported Git configuration, zero non-ignored untracked
  paths, and 66 preserved ignored paths all under `graphify-out/**`
- Preserved ignored workspace material: `graphify-out/**`; excluded from
  product authority and not deleted, moved or staged
- State after commit: owner-approved; unclaimed; not implemented; not accepted
- Recovery R6 state: not authorized and not executable yet
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

The owner's operative direction is:

```text
Proceed autonomously with the repository-consistent minimal authority. If immutable identity analysis proves that a bounded Correction C2 iteration 007 must precede any Recovery R6 dispatch, record/commit that exact iteration-007 authority instead of forcing an invalid R6. You are authorized to make that governance decision within PSCAN-06, while preserving the activation-only/session boundary and zero-spend/no-subscription constraints.
```

The owner also granted the local implementation approval needed to continue
PSCAN-06 autonomously, while excluding online subscription purchases and every
nonzero-spend action. This bundle applies that direction narrowly: it records
local-only iteration 007 authority. It does not grant a remote Recovery R6
execution, create a tag, change a workflow, or spend money.

This is a bounded iteration inside the already activated PSCAN-06 Correction
C2. It does not create, select or activate PSCAN-07, PSCAN-08 or another task.
This session records authority only. A genuinely fresh session with an ID
different from the approval session must claim only Correction C2 iteration
007 from this committed bundle before implementation.

## Controlling and preserved authority

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract committed-byte SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Recovery R5 terminal failure:
  `evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R5-FAILURE-001.md`
- Recovery R5 failure committed-byte SHA-256:
  `53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0`
- Recovery R5 workflow run: `34582887399`, attempt `1`, conclusion `failure`
- Correction C2 iteration 006 local acceptance:
  `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-ACCEPTANCE.md`
- Iteration 006 acceptance committed-byte SHA-256:
  `9E045C7D22841EA89F98223848C8530396F910F729B2E25BAB20541C5F0CD6AC`
- Accepted iteration 006 candidate:
  `cf1679f9bca24887340fa4060f37d1ebff21f305`
- Accepted iteration 006 candidate tree:
  `6cb0198c71881baff9c601c31ad6aa84217f4a49`
- Iteration 006 evidence-bearing review parent:
  `6ae41db6c52c7bb3efd627131dee79f5fb3bc851`
- Iteration 006 acceptance commit:
  `5ff14709520ba99b105ee515e5482455c8511efc`
- Locked product tag: `v1.0.0` ->
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Locked Correction C1 tooling tag: `release-tooling-v1.0.0-c1` ->
  `3fb7592889820fa2739a4a53588e073689621809`
- Locked Correction C2 tooling tag: `release-tooling-v1.0.0-c2` ->
  `faef8435322c9096df09b56969662411f17356ea`
- Locked C2 tag ruleset: `22895383`, active, exact-name update/deletion
  denial, no bypass
- Preserved remote `main`:
  `a82a3a64a04ee2d8b60757866c9864cd1b73b54b`

Recovery R5, iteration 006, all three rejected iteration-006 candidates, every
earlier authority and evidence record, the product/C1/C2 tags, their rulesets,
existing schemas and remote facts remain immutable history. Iteration 007 does
not reinterpret Recovery R5 as a pass or restore its consumed tag push or
workflow dispatch.

## Live read-only state at selection

Approval-session read-only verification established:

- authenticated GitHub account `ThameeraDananjaya` with existing `workflow`
  and repository scopes;
- public repository `ThameeraDananjaya/project-agnostic-secret-scanner`,
  numeric repository ID `1355442997`, default branch `main`;
- remote `main` and all three locked tags exactly as recorded above;
- no local or remote `release-tooling-v1.0.0-c2-r6` tag;
- Recovery R5 run `34582887399` remains the only C2 run from the C2 tag and is
  terminally failed;
- zero repository artifacts, releases and deployments;
- exact repository variable `PSCAN_RELEASE_C2_GATE` absent;
- Actions enabled with selected full-SHA actions, required SHA pinning,
  default read permission and one-day artifact/log retention; and
- current September 2026 GitHub billing API usage has net amount `USD 0`.

The account-budget UI could not be refreshed in this local-only selection
session. That is not used as remote authority. Any later R6 authority and
execution must freshly prove the account-level `USD 0` budget and stop-usage
control, current public-standard-runner terms, net billing, and every other
zero-spend prerequisite before mutation.

## Why direct Recovery R6 is invalid

The accepted iteration 006 candidate fixes the PowerShell PID collision but
does not change the accepted C2 release identity files. Committed Git-object
comparison proves both the old workflow and schema are byte-identical at
pre-fix candidate `faef8435322c9096df09b56969662411f17356ea` and accepted
iteration-006 candidate `cf1679f9bca24887340fa4060f37d1ebff21f305`:

- `.github/workflows/release-recovery-v1.0.0.yml` blob
  `3aa42627628b8b5298d854d29b3880cb19dc36ff`, committed-byte SHA-256
  `C5F40F1B32E87C005FE33EE607AF7E3D19EE4F0173C21619E21158C31AA0BDB4`;
- `contracts/release-manifest/schema-2.1.json` blob
  `f8b5232b8667636fec8f4f0c7f083e9a0297b0cb`, committed-byte SHA-256
  `CAA9CD26665CC3A3550AFFEA7490A0F1F277A69B10B1F1537787616E8AB973CE`.

The existing workflow requires all of these facts simultaneously:

```text
github.ref == refs/tags/release-tooling-v1.0.0-c2
github.sha == inputs.release_tooling_revision
github.workflow_ref == .../release-recovery-v1.0.0.yml@refs/tags/release-tooling-v1.0.0-c2
github.workflow_sha == inputs.release_tooling_revision
```

The protected C2 tag points to the pre-fix candidate. Dispatching it with the
accepted iteration-006 commit as input fails the SHA gate; dispatching it with
the old commit reruns the known defect. Moving or recreating the tag would
violate immutable evidence and active no-bypass protection. A different ref
fails the workflow-ref gate. Schema 2.1, the builder and the verifier also bind
the old tag, workflow path, ref and certificate identity. Therefore no valid
direct R6 dispatch exists from current committed state.

## Approved objective and selected future identities

Create a locally implemented, independently reviewed, append-only release
identity roll-forward that can later carry the accepted iteration-006 bytes
without changing any locked identity. The selected exact new identities are:

```text
tooling tag: release-tooling-v1.0.0-c2-r6
workflow path: .github/workflows/release-recovery-v1.0.0-c2-r6.yml
manifest schema: 2.2
workflow ref: refs/tags/release-tooling-v1.0.0-c2-r6
```

These identities are selected for local iteration-007 implementation only.
The tag, remote workflow and schema file do not yet exist. Selection is not
creation, tagging, remote authority, execution or acceptance.

The implementation must add schema 2.2 and the R6 workflow path rather than
editing or reinterpreting schema 2.1 or the existing C2 workflow. It may update
only the active builder, verifier and validation routes required to emit and
verify schema 2.2 with the selected exact tag/ref/path/certificate identity.
Historical schema 2.0/2.1 and C1/C2 verification behavior must remain accepted
under their original exact identities.

## Allowed implementation paths

Correction C2 iteration 007 may change only:

```text
.github/workflows/release-recovery-v1.0.0-c2-r6.yml
contracts/release-manifest/schema-2.2.json
build/release/build.ps1
build/release/docker-execution.ps1
build/release/test-crlf-shell-payloads.ps1
build/release/cmd/release-verifier/main.go
internal/verify/release.go
tests/integration/supply-chain/release_test.go
tests/acceptance/supply-chain/release_test.go
docs/decisions/DEC-004-C2-R6-IMMUTABLE-IDENTITY.md
docs/release/OFFLINE-VERIFICATION-RUNBOOK.md
docs/release/PSCAN-06-RELEASE-PLAN.md
docs/release/RELEASE-AUTHORITY.md
docs/validation/PSCAN-06-VALIDATION.md
README.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/**
```

Every absent path is forbidden. In particular, the implementation must not
change `.github/workflows/release-recovery-v1.0.0.yml`,
`contracts/release-manifest/schema-2.1.json`, acquisition inputs, Docker
operation tables, product source, scanner behavior or PASS contract files. A
demonstrated need for another path stops the iteration and returns exact
evidence to the owner.

## Mandatory implementation design

1. Add the selected R6 workflow as a new file derived from the committed C2
   workflow. Preserve job ordering, standard `ubuntu-24.04`, time limits,
   permissions, full-SHA actions, pinned inputs, offline build phases,
   artifact name/retention/compression and the disabled signing gate. Change
   only identity-bearing name, concurrency, tag/ref, workflow path and
   certificate-identity values required by the new immutable identity.
2. Add schema 2.2 as a new file derived from schema 2.1. Change only schema
   version/title/identifier and the exact release-tooling tag, workflow path,
   ref and certificate identity. Preserve every field, type, required member,
   additional-property denial and product-source binding.
3. Update the active builder and embedded release verifier to emit and require
   schema 2.2 and the selected R6 identity. Do not change product tag/commit,
   toolchain/dependency pins, asset contents, redaction or public outcome
   behavior.
4. Extend the generic verifier to recognize schema 2.2 with the same or
   stricter exact dual-identity checks. Preserve schema 2.0 and 2.1 behavior.
   Reject cross-version mixtures, including old tag/ref/path under 2.2 and new
   tag/ref/path under 2.1.
5. Change `docker-execution.ps1` only as needed to materialize the new schema.
   The independently accepted `Get-LinuxSessionMembers` function, its
   iteration-006 exact predicate/data flow, process boundary, Docker operation
   table and every containment invariant must remain unchanged.
6. Add path-bounded tests proving exact identity agreement across workflow,
   schema, builder, verifier and trust policy, plus negative cross-version and
   single-field mutation cases. Preserve all historical schema tests.
7. Record an append-only decision/evidence chain. Do not edit historical
   decision, authority, failure, rejection or acceptance records.

## Required local checks and success criteria

1. Exact graphify-aware current-checkout source verification passes with exact
   index, flags and tracked bytes, zero non-ignored untracked paths and every
   ignored path confined to preserved `graphify-out/**`.
2. A separate exact clean temporary clone or materialization of the committed
   candidate passes the product source-trust verifier without ignored files.
   The main checkout's preserved graphify material must not be deleted or
   moved to manufacture this result.
3. JSON schema, YAML, PowerShell and Go parsing succeeds for every changed
   implementation path.
4. Offline, no-download tests prove schema 2.0/2.1 preservation, exact schema
   2.2 identity, builder/verifier agreement and cross-version rejection. If a
   required dependency is not already available locally, stop without
   downloading it and record the unproved check.
5. The complete available-host local no-Docker iteration-006 PID regression
   and iteration-005 native/hostile matrix pass with no surviving marked
   process or temporary fixture.
6. Exact diff/path confinement and `git diff --check` pass. No existing tag,
   old workflow, schema 2.1, historical evidence or ignored graphify path is
   changed or staged.
7. Independent skeptical review in a separate genuinely fresh session examines
   the exact candidate and reruns proportionate local no-Docker/offline checks
   before any iteration-007 acceptance record.
8. Recovery R5 remains terminal. Recovery R6 remains unauthorized until local
   iteration 007 is independently accepted and a separate exact R6 authority
   records the accepted candidate/commit/tree and fresh zero-spend remote
   state.
9. PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains proposed and
   unselected; PSCAN-08 remains inactive and ineligible.

## Forbidden scope and reserved gates

- No implementation in this authority session.
- No edit to the existing C2 workflow, schema 2.1, PASS-OUTCOME-SPEC-001,
  PASS-SPEC-001, TRACEABILITY, historical decision/evidence files, product
  source, scanner, policy, allowlist, receipt or consuming-project behavior.
- No Docker command, image pull, image execution, dependency download,
  external scanner/toolchain acquisition, build requiring network, or network
  action in iteration 007.
- No remote mutation, push, fetch, tag creation/movement/deletion, ruleset or
  setting action, workflow dispatch, rerun or retry.
- No credential, signing key, paid capability, subscription, nonzero spend,
  TruffleHog work, signing, attestation, draft, release, deployment or
  publication action.
- No deletion, movement or staging of preserved ignored `graphify-out/**`.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

## Fresh implementation-session requirements

The implementation session must:

1. have a session ID different from
   `01a09141-1ba5-7931-bedc-bee2402fbf95`;
2. begin on exact branch `main` at this committed authority bundle with a clean
   tracked, staged and non-ignored untracked state;
3. preserve all ignored `graphify-out/**` and use the two-part source-trust
   method required above;
4. read the controlling contract, PSCAN-06 task, reading map, tracker, this
   authority, R5 terminal failure and iteration-006 acceptance completely;
5. claim only PSCAN-06 Correction C2 iteration 007 and record its bounded
   preflight before implementation; and
6. stop fail-closed on any identity, path, scope, source-trust, offline-test or
   authority mismatch.

## Approval-session exclusions

This approval session introduces no product, implementation, test, workflow,
schema, verifier or build-logic change. It performs no Docker command,
dependency acquisition, remote mutation, tag operation, workflow run, signing,
attestation, draft, release, deployment, publication, purchase, subscription
or successor action. Its only durable changes are this bounded authority
bundle and synchronized living PSCAN-06 state.
