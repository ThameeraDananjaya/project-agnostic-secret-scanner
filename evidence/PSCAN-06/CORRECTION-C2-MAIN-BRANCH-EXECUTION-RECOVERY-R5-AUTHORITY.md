# PSCAN-06 Correction C2 main-branch execution Recovery R5 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Recovery: `R5`
- Gate: build-only actual-Linux and genuine-Docker proof
- Event: owner approval of one bounded split-tag execution recovery
- Owner approval received: `2026-09-11`
- Approval session: `01a0737f-d3fe-7c93-9d3e-db1f631c3df9`
- Required starting evidence commit before this bundle:
  `6fbeb633bf373d2f64d95d0ff93c328631a99a05`
- Required starting evidence tree:
  `1cac2f88f7373701448f311bd27bae04ed7299c6`
- Required starting parent:
  `a82a3a64a04ee2d8b60757866c9864cd1b73b54b`
- Recovery-authority commit: this authority-bundle commit
- Approval-session checkout: exact saved-project root on branch `main`
- Pre-change tracked, staged and non-ignored untracked Git status: clean
- Preserved ignored workspace material: `graphify-out/**`; excluded from
  product authority, staging and this commit
- State after commit: owner-approved; not claimed for execution; unexecuted
- Workflow dispatches consumed by Recoveries R1-R5: zero before R5 execution
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

The owner's exact standalone approval is:

```text
APPROVE PSCAN-06 C2 MAIN-BRANCH EXECUTION RECOVERY R5 FROM RECOVERY-R4 FAILURE 001 SHA256 5F694FF128515A2B79A2BFA500F299FECAB8836763365E9D774DA5D847D68BEE AT EVIDENCE COMMIT 6FBEB633BF373D2F64D95D0FF93C328631A99A05 PRESERVING REMOTE MAIN AT A82A3A64A04EE2D8B60757866C9864CD1B73B54B AND EXISTING LOCAL TAG REFS/TAGS/RELEASE-TOOLING-V1.0.0-C2 AT FAEF8435322C9096DF09B56969662411F17356EA WITH REMOTE TAG AND RULESET ABSENT ZERO-SPEND ONE TAG PUSH USING EXACT REFSPEC REFS/TAGS/RELEASE-TOOLING-V1.0.0-C2:REFS/TAGS/RELEASE-TOOLING-V1.0.0-C2 ONE WORKFLOW DISPATCH MAXIMUM ONE UNSIGNED UNCOMPRESSED ARTIFACT RETAINED ONE DAY NO-TAG-RECREATE NO-LOCAL-DOCKER NO-SIGNING NO-ATTESTATION NO-DRAFT NO-PUBLICATION NO-RERUN NO-SUCCESSOR
```

This is a bounded continuation inside the already activated PSCAN-06
Correction C2. It does not create, select or activate another task. This
session records Recovery R5 authority only. It does not claim or execute the
recovered gate.

## Preserved terminal Recovery R4 failure

Recovery R4 ran in genuinely fresh saved-project `main` session
`01a08f7a-7d13-7ee1-be72-a997e4ff9cb6`. Its complete mandatory preflight
passed. It fast-forwarded remote `main` exactly to its authority commit and
created the exact local lightweight C2 tooling tag at the accepted candidate.
Its first tag-push command failed locally because shell expansion produced an
invalid refspec. The no-retry gate stopped the transaction immediately.

The immutable failure record is:

`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R4-FAILURE-001.md`

Its SHA-256 was reverified before this bundle as:

`5F694FF128515A2B79A2BFA500F299FECAB8836763365E9D774DA5D847D68BEE`

Recovery R5 does not reinterpret Recovery R4 as a pass or reuse its preflight.
It changes only the retry state for the preserved split-tag condition and
permits one separately authorized attempt after every current preflight fact
passes.

## Exact preserved split state

Approval-time read-only verification established:

- local branch `main` is clean at evidence commit
  `6fbeb633bf373d2f64d95d0ff93c328631a99a05`;
- remote `main` is exactly
  `a82a3a64a04ee2d8b60757866c9864cd1b73b54b`;
- local `refs/tags/release-tooling-v1.0.0-c2` exists as a lightweight tag
  resolving directly to commit
  `faef8435322c9096df09b56969662411f17356ea` and tree
  `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`;
- remote `refs/tags/release-tooling-v1.0.0-c2` is absent;
- an exact-name C2 tag ruleset is absent;
- no C2 workflow run or C2 artifact exists;
- the authenticated repository remains public repository ID `1355442997` at
  `ThameeraDananjaya/project-agnostic-secret-scanner`; and
- the active GitHub credential is for `ThameeraDananjaya` with scopes `gist`,
  `read:org`, `repo`, `user` and `workflow`.

The capitalized owner-decision token denotes the already-established exact
lowercase Git ref identity. Git refs are case-sensitive. Execution must not
create an uppercase or differently cased ref. The sole executable refspec is
the following literal argument:

```text
refs/tags/release-tooling-v1.0.0-c2:refs/tags/release-tooling-v1.0.0-c2
```

Both sides passed local `git check-ref-format` before this bundle. This local
syntax check made no remote change and does not substitute for fresh preflight.

## Controlling and preserved authority

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Original build-only proof-gate authority:
  `evidence/PSCAN-06/CORRECTION-C2-BUILD-ONLY-LINUX-DOCKER-PROOF-GATE-AUTHORITY.md`
- Original gate-authority SHA-256:
  `2A84E7109F151E55399870FB36D64AEDB2213E123976779CD1E9648AC1B9DAC2`
- Recovery R4 authority:
  `evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R4-AUTHORITY.md`
- Accepted Correction C2 tooling candidate:
  `faef8435322c9096df09b56969662411f17356ea`
- Accepted tooling tree:
  `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`
- Required evidence HEAD:
  `b6d341284cf63baa0502fcba319ea8cca3c7eb3a`
- Required evidence tree:
  `983028703e229abb082938d0e4417c506aac0d8c`
- Candidate relation: direct parent of the required evidence HEAD
- Iteration-005 local independent verdict:
  `LOCAL_ACCEPTANCE_PASS_REMOTE_PROOF_OPEN`
- Recovery workflow: `.github/workflows/release-recovery-v1.0.0.yml`
- Workflow SHA-256:
  `C5F40F1B32E87C005FE33EE607AF7E3D19EE4F0173C21619E21158C31AA0BDB4`
- Locked product tag: `v1.0.0`
- Locked product-source commit:
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Locked Correction C1 tooling tag: `release-tooling-v1.0.0-c1`
- Locked Correction C1 tooling commit:
  `3fb7592889820fa2739a4a53588e073689621809`
- Exact approved C2 tooling tag name: `release-tooling-v1.0.0-c2`
- Exact pinned build image:
  `docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`
- Intended remote: `ThameeraDananjaya/project-agnostic-secret-scanner`
- Recorded remote repository ID: `1355442997`
- Maximum spend: `USD 0`
- Tag-push limit: at most one
- Workflow-dispatch limit: at most one
- Authorized artifact: at most one unsigned, uncompressed artifact retained
  for one day

Every predecessor authority, candidate, evidence record, locked tag, schema,
workflow byte, action/image/dependency pin, test boundary and product/tooling
role remains immutable. Recovery R5 weakens or waives no original condition.

## Required genuinely fresh execution context

The recovered gate must run in a genuinely fresh Codex task configured as
follows before its first repository or remote action:

1. Use the saved Project-Agnostic Secret Scanner project directly, not a Codex
   worktree, clone or detached checkout.
2. Use exact root
   `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`.
3. Require branch `main` at this committed Recovery R5 authority bundle and
   its exact tree, with direct parent
   `6fbeb633bf373d2f64d95d0ff93c328631a99a05`.
4. Use `gpt-daybreak-blue-latest` with reasoning effort `xhigh` and verify both
   from task/session metadata.
5. Require clean tracked, staged and non-ignored untracked state. Preserved
   ignored `graphify-out/**` is excluded from product authority and staging.
   No other unexplained ignored path is permitted.
6. Require an execution session ID different from approval session
   `01a0737f-d3fe-7c93-9d3e-db1f631c3df9`.
7. If any root, branch, commit, tree, parent, task, model, reasoning, source
   trust, GitHub identity or state fact is wrong or ambiguous, stop without
   repair, cleanup, authority copying or remote mutation.

This approval session is not the fresh execution session.

## Mandatory fresh preflight

Before mutation, the fresh task must fail closed unless current local and
remote read-only evidence proves all of the following:

1. The exact fresh execution context above and complete every-byte local
   source trust pass.
2. The Recovery R4 failure record and hash, preserved remote `main`, local C2
   tag object/commit/tree, absent remote C2 tag and absent C2 ruleset exactly
   match this authority.
3. Every predecessor authority, accepted candidate/evidence relation, locked
   tag and ruleset, workflow byte and committed build boundary remains exact.
4. The workflow build job remains on standard `ubuntu-24.04`, has only
   `contents: read`, uses only the recorded full-SHA action pins, and its sole
   C2 artifact still has exact name, `retention-days: 1` and
   `compression-level: 0`.
5. Exact repository variable `PSCAN_RELEASE_C2_GATE` remains absent.
6. Current authenticated account, repository identity, Actions policy,
   standard-runner availability, billing, artifact-storage allowance,
   account-level `USD 0` Actions budget and enabled stop-usage control are
   trusted and together prove maximum spend `USD 0`.
7. No workflow run or artifact consumed any earlier C2 proof-gate attempt; no
   release or deployment exists; PSCAN-06 remains the sole open task;
   PSCAN-07 remains proposed and unselected; and PSCAN-08 remains inactive and
   ineligible.
8. Both halves of the sole literal refspec pass `git check-ref-format`, its
   source resolves to the preserved local candidate, and its destination is
   absent remotely.

Missing, stale, permission-denied, rate-limited, contradictory, ambiguous or
otherwise untrusted evidence stops the gate before mutation. Passing preflight
does not widen scope.

## Exactly authorized execution

Only after every preflight item passes, the fresh task may perform this single
ordered transaction:

1. Push the already-existing local tag exactly once, without force, using one
   literal refspec argument and no wildcard or variable expansion:

   ```text
   git push -- origin 'refs/tags/release-tooling-v1.0.0-c2:refs/tags/release-tooling-v1.0.0-c2'
   ```

   Do not create, recreate, move, delete, retag or otherwise mutate the local
   tag. Push no branch or other ref.
2. Read back the remote tag object, commit and tree. Create or confirm one
   active exact-name tag ruleset denying update and deletion with no bypass,
   then read back its exact scope and enforcement.
3. Confirm Actions artifact/log retention is exactly one day. It may be
   reduced to one day only if trusted read-back proves a different value.
   Change no other setting.
4. Dispatch `.github/workflows/release-recovery-v1.0.0.yml` exactly once from
   `refs/tags/release-tooling-v1.0.0-c2` with only:

   ```text
   product_source_revision=a13c28fe7273bc8dc6545f97966a02889524eb4c
   release_tooling_revision=faef8435322c9096df09b56969662411f17356ea
   release_version=v1.0.0
   ```

5. Permit only the workflow `build` job on standard `ubuntu-24.04`. It may run
   its exact committed actual-Linux checks, genuine Docker engine/image/
   container operations, one bounded canonical-digest pull only after
   conclusive absence, pinned public dependency acquisition, two
   network-disabled builds, every-byte comparison and transfer of exactly one
   unsigned uncompressed artifact named
   `pscan-v1.0.0-c2-unsigned-candidate`, retained for one day.
6. Observe the single run to terminal state and read back its exact repository,
   workflow, ref, SHA, inputs, job, steps, runner, artifact, digest, retention,
   deployment and release facts. Download only that unsigned artifact for
   bounded independent verification; do not copy it to a release, public
   location or consuming project.
7. Record the complete pass or fail-closed outcome locally under PSCAN-06. Do
   not push post-run evidence, rerun, repair, retry, mutate or promote a result
   under this authority. Return to the owner for the next exact gate.

Every step is sequential and conditional. Any failure, timeout, cancellation,
unexpected artifact, missing proof, unexpected deployment, signing-job
eligibility or identity mismatch stops the transaction. Recovery R5 permits
at most one tag push and one workflow dispatch, with no second attempt.

## Explicitly forbidden

- No product, workflow, build, test, schema, verifier or release-code change.
- No remote `main` push or mutation; it remains at the preserved R4 authority.
- No local Docker start, stop, pull, image, container or other Docker action on
  the owner's PC.
- No tag recreation, movement, deletion, force push, alternate refspec,
  wildcard, broad tag pattern, case variant or bypass actor.
- No creation, update or use of `PSCAN_RELEASE_C2_GATE`; no signing environment
  approval or deployment action.
- No signing, OIDC, key, certificate, attestation, SBOM promotion, draft,
  release, publication or consumer-project action.
- No paid resource, spend above `USD 0`, workflow rerun, second dispatch,
  second tag push, repair attempt or successor work.
- No post-run evidence push or public artifact outside the single private-to-
  workflow unsigned artifact transfer already bounded above.

## Closure boundary

Recovery R5 can prove only the remaining build-only actual-Linux,
genuine-Docker, reproducibility and artifact-integrity gate. A passing run does
not authorize signing, attestation, drafting, publication or final PSCAN-06
acceptance. A failing run is terminal under this authority. In either case the
fresh task records one local evidence outcome and returns to the owner. No
successor is selected or activated.
