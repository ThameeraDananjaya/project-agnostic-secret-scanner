# PSCAN-06 Correction C2 main-branch execution Recovery R4 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Recovery: `R4`
- Gate: build-only actual-Linux and genuine-Docker proof
- Event: owner approval of one bounded post-billing-repair execution recovery
- Owner approval received: `2026-09-11`
- Approval session: `01a0737f-d3fe-7c93-9d3e-db1f631c3df9`
- Required starting evidence commit before this bundle:
  `6b7ae32f8768c1e4687e5134696688088f65386b`
- Required starting evidence tree:
  `6eab4570b5ca3f3114a2ca8f22024fff9b7d325d`
- Required starting parent:
  `a4f5b85a3ac63ea31b4bb911be3445f2223b6db1`
- Recovery-authority commit: this authority-bundle commit
- Approval-session checkout: exact saved-project root on branch `main`
- Pre-change tracked, staged and non-ignored untracked Git status: clean
- Preserved ignored workspace material: `graphify-out/**`; excluded from
  product authority, staging and this commit
- State after commit: owner-approved; not claimed for execution; unexecuted
- Workflow dispatches consumed by the original gate and Recoveries R1-R4:
  zero under the current C2 proof-gate authority
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

The owner's exact standalone approval is:

```text
APPROVE PSCAN-06 C2 MAIN-BRANCH EXECUTION RECOVERY R4 FROM RECOVERY-R3 FAILURE 001 SHA256 339E5D66B655B5753D239D5347F1BBAFF0E5858C4CCC289368196D014AE07C12 AT EVIDENCE COMMIT 6B7AE32F8768C1E4687E5134696688088F65386B USING A GENUINELY FRESH SAVED-PROJECT MAIN CHECKOUT ZERO-SPEND ONE WORKFLOW DISPATCH MAXIMUM ONE UNSIGNED UNCOMPRESSED ARTIFACT RETAINED ONE DAY NO-LOCAL-DOCKER NO-SIGNING NO-ATTESTATION NO-DRAFT NO-PUBLICATION NO-RERUN NO-SUCCESSOR
```

This is a bounded continuation inside the already activated PSCAN-06
Correction C2. It does not create, select or activate another PSCAN task. This
session records Recovery R4 authority and current read-only billing evidence
only. It does not claim or execute the recovered gate.

## Preserved terminal Recovery R3 failure

Recovery R3 ran in genuinely fresh saved-project `main` session
`01a08461-59eb-7721-b1f3-f11721e36bfd`. Exact local execution context, source
trust, authority, workflow bytes and host-level GitHub authentication passed.
The authenticated credential lacked the required `user` scope, so current
account billing, artifact-storage allowance and stop-usage evidence could not
be read. Recovery R3 stopped before remote mutation, workflow dispatch,
Docker, build or artifact action and consumed zero workflow dispatches.

The immutable failure record is:

`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R3-FAILURE-001.md`

Its SHA-256 was reverified before this bundle as:

`339E5D66B655B5753D239D5347F1BBAFF0E5858C4CCC289368196D014AE07C12`

Recovery R4 does not reinterpret Recovery R3 as a pass or reuse its preflight.
It permits one separately authorized attempt only after every current fresh
preflight fact passes.

## Current read-only billing repair evidence

Before this authority bundle was written, read-only verification established:

- `gh auth status --hostname github.com` passed for active account
  `ThameeraDananjaya` and reported scopes `gist`, `read:org`, `repo`, `user`
  and `workflow`.
- The current September 2026 Actions billing summary was readable. It reported
  447 Linux minutes, gross and discounted amount `USD 2.682`, and net amount
  `USD 0`.
- The same summary reported `0.066191896` Actions storage GB-hours, gross and
  discounted amount `USD 0.0000221`, and net amount `USD 0`.
- The authenticated GitHub billing UI displayed the account-level product
  budget for `Actions` as `USD 0`, current spent amount `USD 0`, and
  `Stop usage` as `Yes`.
- No billing, budget, payment, credential, repository, workflow or other
  setting was changed during verification.

These facts repair the permission-denied condition that stopped Recovery R3.
They are approval-time evidence only. The fresh Recovery R4 execution task
must read them again and fail closed if any current value is missing, stale,
contradictory or no longer proves maximum spend `USD 0`.

## Controlling and preserved authority

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Original build-only proof-gate authority:
  `evidence/PSCAN-06/CORRECTION-C2-BUILD-ONLY-LINUX-DOCKER-PROOF-GATE-AUTHORITY.md`
- Original gate-authority SHA-256:
  `2A84E7109F151E55399870FB36D64AEDB2213E123976779CD1E9648AC1B9DAC2`
- Recovery R2 authority:
  `evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R2-AUTHORITY.md`
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
- Locked product-source tree:
  `217b711ddea51fd0ea7e808edd2e27fdecef8427`
- Locked Correction C1 tooling tag: `release-tooling-v1.0.0-c1`
- Locked Correction C1 tooling commit:
  `3fb7592889820fa2739a4a53588e073689621809`
- Exact approved C2 tooling tag name: `release-tooling-v1.0.0-c2`
- Exact pinned build image:
  `docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`
- Intended remote: `ThameeraDananjaya/project-agnostic-secret-scanner`
- Recorded remote repository ID: `1355442997`
- Maximum spend: `USD 0`
- Workflow dispatch limit: at most one
- Authorized artifact: at most one unsigned, uncompressed artifact retained
  for one day

Every predecessor authority, candidate, evidence record, locked tag, schema,
workflow byte, action/image/dependency pin, test boundary and product/tooling
role remains immutable. Recovery R4 changes only the retry state after the
billing permission repair. It weakens or waives no original condition.

## Required genuinely fresh execution context

The recovered gate must run in a genuinely fresh Codex task configured as
follows before its first repository or remote action:

1. Use the saved Project-Agnostic Secret Scanner project directly, not a Codex
   worktree, clone or detached checkout.
2. The exact root must be
   `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`.
3. The checkout must be branch `main` at this committed Recovery R4 authority
   bundle and its exact tree, with direct parent
   `6b7ae32f8768c1e4687e5134696688088f65386b`.
4. Use `gpt-daybreak-blue-latest` with reasoning effort `xhigh` and verify both
   values from the task/session metadata.
5. Tracked, staged and non-ignored untracked state must be clean. Preserved
   ignored `graphify-out/**` is known and excluded from product authority,
   staging, remote evidence and execution decisions. No other unexplained
   ignored path is permitted.
6. If the task starts detached, off `main`, at another root or identity, with
   invalid GitHub authentication, or in a dirty or ambiguous state, stop
   without branch repair, cleanup, authority copying or remote mutation.

This approval session is not the fresh execution session. No execution task or
successor starts automatically from the authority commit.

## Mandatory fresh preflight

Before mutation, the fresh task must fail closed unless current local and
remote read-only evidence proves all of the following:

1. The execution context above is exact, including root, branch, Recovery R4
   authority commit/tree/parent, clean source trust, accepted candidate,
   evidence HEAD and required relation.
2. The workflow at the accepted candidate has the recorded exact bytes and
   SHA-256, all build-job actions remain full-SHA pinned, the build job has only
   `contents: read`, and the single artifact retains `retention-days: 1` and
   `compression-level: 0`.
3. Signing remains gated only by exact repository variable
   `PSCAN_RELEASE_C2_GATE=PSCAN-06-C2-SIGNING-APPROVED`; that variable is
   absent. If it exists with any value, stop without changing or dispatching.
4. The authenticated account, repository owner/name and numeric ID, public
   visibility, default branch, origin URL, Actions policy, standard runner
   availability, current account billing, Actions storage allowance, `USD 0`
   Actions budget and enabled stop-usage control are trusted and together prove
   the approved run can incur no charge.
5. Locked tags `v1.0.0` and `release-tooling-v1.0.0-c1` still resolve exactly
   to their recorded objects and retain active update/deletion denial with no
   bypass. The C2 tag is absent locally and remotely.
6. No workflow run, C2 tag or artifact consumed the original gate or
   Recoveries R1-R4; no other PSCAN task is active; PSCAN-07 remains proposed
   and unselected; PSCAN-08 remains inactive and ineligible.

Missing, stale, permission-denied, rate-limited, contradictory, ambiguous or
otherwise untrusted evidence stops the gate before mutation. Passing preflight
does not widen scope.

## Exactly authorized execution

Only after every preflight item passes, the fresh task may perform this single
ordered transaction:

1. Fast-forward remote `main`, without force, to the exact Recovery R4
   authority commit. Push no unrelated path, ref or commit.
2. Create exact tag `release-tooling-v1.0.0-c2` at accepted tooling candidate
   `faef8435322c9096df09b56969662411f17356ea`; push only that tag; create or
   confirm an active exact-name tag ruleset denying update and deletion with no
   bypass; and read back the tag object, commit, tree and rules.
3. Confirm Actions artifact/log retention is exactly one day. It may be reduced
   to one day only if trusted read-back proves a different value. Change no
   other repository, Actions, environment, branch, release or security setting.
4. Dispatch `.github/workflows/release-recovery-v1.0.0.yml` exactly once from
   `refs/tags/release-tooling-v1.0.0-c2` with only:

   ```text
   product_source_revision=a13c28fe7273bc8dc6545f97966a02889524eb4c
   release_tooling_revision=faef8435322c9096df09b56969662411f17356ea
   release_version=v1.0.0
   ```

5. Permit only the workflow `build` job on the standard `ubuntu-24.04` runner.
   This authorizes its exact committed actual-Linux checks, genuine Docker
   engine/image/container operations, one bounded canonical-digest pull only
   after conclusive absence, pinned public dependency acquisition, two
   network-disabled builds, every-byte comparison and transfer of exactly one
   unsigned uncompressed artifact named
   `pscan-v1.0.0-c2-unsigned-candidate`, retained for one day.
6. Observe the single run to a terminal state and read back its exact
   repository, workflow, ref, SHA, inputs, job, steps, runner, artifact,
   digest, retention, deployment and release facts. Download the exact unsigned
   artifact only for bounded independent verification; do not copy it to a
   release, public location or consuming project.
7. Record the complete pass or fail-closed outcome locally under PSCAN-06. Do
   not push post-run evidence, rerun, repair, retry, mutate or promote any
   result under this authority. Return to the owner for the next exact gate.

Every step is sequential and conditional. Any failure, timeout, cancellation,
unexpected artifact, missing proof, unexpected deployment, signing-job
eligibility or identity mismatch stops the transaction. Recovery R4 permits at
most one workflow dispatch and no second execution or dispatch.

## Explicitly forbidden

- No product, workflow, build, test, schema, verifier or release-code change.
- No local Docker start, stop, pull, image or container action on the owner's
  PC under this authority.
- No detached/worktree execution, branch repair, force push, history rewrite,
  tag movement/deletion/recreation, broad tag pattern or bypass actor.
- No creation, update or use of `PSCAN_RELEASE_C2_GATE`; no signing environment
  approval or deployment. The `sign-attest-and-draft` job must be skipped and
  ineligible.
- No signing, OIDC request, signature, trusted-root action, attestation, draft,
  release asset, publication or immutable-release finalization.
- No credential or signing-key creation/handling, paid runner, paid storage,
  overage, subscription, provider commitment or spend above `USD 0`.
- No extra workflow run, rerun, retry, correction, tag, branch, PR, issue,
  repository, remote resource or setting change.
- No scanner adoption, consuming-project data, receipt, deployment,
  production, cross-project reporting, TruffleHog or real secret material.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

## Approval-session actions and exclusions

This session records the Recovery R4 authority and performs only read-only
local identity, GitHub authentication, billing API and billing-UI checks. It
performs no push, tag or ruleset operation, setting change, workflow dispatch,
Docker/image/container action, network acquisition, build, artifact action,
signing, attestation, draft, release or publication. It creates no remote
resource, key or paid commitment.

PSCAN-06 remains open and unaccepted overall. PSCAN-07 remains proposed and
unselected. PSCAN-08 remains inactive and ineligible. No successor is selected
or activated.
