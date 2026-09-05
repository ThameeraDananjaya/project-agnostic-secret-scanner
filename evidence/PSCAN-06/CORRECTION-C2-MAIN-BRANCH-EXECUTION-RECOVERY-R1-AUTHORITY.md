# PSCAN-06 Correction C2 main-branch execution Recovery R1 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Recovery: `R1`
- Gate: build-only actual-Linux and genuine-Docker proof
- Event: exact owner approval of the bounded pre-dispatch execution recovery
- Owner command:
  `APPROVE PSCAN-06 C2 MAIN-BRANCH EXECUTION RECOVERY R1 FROM PRE-DISPATCH FAILURE 001 SHA256 7EB1CB03E85A933C1BCFEB435FA85556170AE21F1112CFFA55041C2A733E6A5D AT AUTHORITY COMMIT 268e8f339318a141d059115c9da0d19a01448aaa USING A GENUINELY FRESH SAVED-PROJECT MAIN CHECKOUT ZERO-SPEND ONE WORKFLOW DISPATCH MAXIMUM NO-SIGNING NO-ATTESTATION NO-DRAFT NO-PUBLICATION`
- Date: `2026-09-06`
- Approval session: `01a053a9-298b-75e1-9d49-7d2261bbb951`
- Required starting authority commit before this recovery bundle:
  `268e8f339318a141d059115c9da0d19a01448aaa`
- Required starting authority tree:
  `e1667892d0c05d9a8db8bac65e580fff54a876ea`
- Recovery-authority commit: this authority-bundle commit
- Approval-session checkout: exact saved-project root on branch `main`
- Pre-change tracked and untracked Git status: clean
- Preserved ignored workspace material: `graphify-out/**`; excluded from
  product authority, staging and this commit
- State after commit: owner-approved; not claimed for execution; unexecuted
- Prior workflow dispatches under the build-only proof-gate authority: zero
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

This recovery is a bounded continuation inside the already activated PSCAN-06
Correction C2. It does not create, select or activate another PSCAN task. This
approval session records authority and preserved failure evidence only. It
performs no remote preflight or mutation and does not execute the recovered
gate.

## Preserved pre-dispatch failure

The first execution session, `01a072e4-50f8-7842-8c7d-28bb034a773f`, ran with
`gpt-daybreak-blue-latest` at `xhigh` from the exact gate-authority commit and
tree but in a detached Codex worktree. It failed closed at mandatory preflight
item 1 because the checkout itself was not on branch `main`.

The immutable imported record is
`evidence/PSCAN-06/CORRECTION-C2-BUILD-ONLY-LINUX-DOCKER-PROOF-GATE-FAILURE-001.md`
with SHA-256
`7EB1CB03E85A933C1BCFEB435FA85556170AE21F1112CFFA55041C2A733E6A5D`.
It proves that the session stopped before any GitHub or other remote preflight,
network request, push, tag, ruleset, setting change, workflow dispatch, Docker
action or artifact action. It created no remote resource and consumed no
workflow-dispatch authority. Recovery R1 does not erase, reinterpret or turn
that terminal result into a pass.

## Controlling and preserved authority

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Original build-only proof-gate authority:
  `evidence/PSCAN-06/CORRECTION-C2-BUILD-ONLY-LINUX-DOCKER-PROOF-GATE-AUTHORITY.md`
- Original gate-authority SHA-256:
  `2A84E7109F151E55399870FB36D64AEDB2213E123976779CD1E9648AC1B9DAC2`
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

Every predecessor authority, candidate, evidence record, locked tag, schema,
workflow byte, action/image/dependency pin, test boundary and product/tooling
role remains immutable. Recovery R1 changes only the admissible task-launch
context from a detached Codex worktree to the exact saved-project `main`
checkout. It does not weaken or waive any original preflight or execution
condition.

## Required fresh execution context

The recovered gate must run in a genuinely fresh Codex task configured as
follows before its first repository or remote action:

1. Use the saved Project-Agnostic Secret Scanner project directly, not a Codex
   worktree, clone or detached checkout.
2. The exact root must be
   `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`.
3. The checkout itself must be on branch `main` at this committed Recovery R1
   authority bundle and its exact tree, with direct parent
   `268e8f339318a141d059115c9da0d19a01448aaa`.
4. The task must use `gpt-daybreak-blue-latest` with reasoning effort `xhigh`
   and verify those values from its own task/session metadata.
5. Tracked, staged and untracked state must be clean. The preserved ignored
   `graphify-out/**` workspace material is known and excluded from product
   authority, staging, remote evidence and execution decisions. No other
   unexplained ignored path is permitted.
6. If the task starts detached, off `main`, at another root or identity, or in
   a dirty or ambiguous state, it must stop without switching branches,
   repairing, cleaning, copying authority to another checkout or making any
   remote request.

The approval session is not the fresh execution session. No successor starts
automatically from this authority commit.

## Mandatory fresh preflight

Before mutation, the fresh execution session must fail closed unless all of
the following are proved from current local and remote read-only evidence:

1. The required execution context above is exact, including root, branch,
   Recovery R1 authority commit/tree/parent, clean source trust, accepted
   candidate, evidence HEAD and their required relation.
2. The workflow at the accepted candidate has the recorded exact bytes and
   SHA-256, every build-job action remains full-SHA pinned, the build job has
   only `contents: read`, and the single artifact retains `retention-days: 1`
   and `compression-level: 0`.
3. The signing job remains gated only by exact repository variable
   `PSCAN_RELEASE_C2_GATE=PSCAN-06-C2-SIGNING-APPROVED`; that variable is
   absent. If it exists with any value, stop without changing or dispatching.
4. The authenticated account, repository owner/name and numeric ID, public
   visibility, default branch, origin URL, Actions policy, standard runner
   availability, account budget and stop-usage control are exactly trusted and
   permit the approved run for USD 0. Any billing, storage, overage, larger-
   runner or provider-commitment risk stops the gate.
5. Locked tags `v1.0.0` and `release-tooling-v1.0.0-c1` still resolve exactly
   to their recorded objects and retain active update/deletion denial with no
   bypass. The C2 tag is absent locally and remotely. Any conflicting, moved,
   deleted, mutable or ambiguous identity stops the gate.
6. No workflow run, C2 tag or artifact consumed the original proof-gate or
   this Recovery R1 one-dispatch authority; no other PSCAN task is active;
   PSCAN-07 remains proposed and unselected; and PSCAN-08 remains inactive and
   ineligible.

A stale, missing, unsupported, rate-limited, permission-denied, contradictory
or otherwise untrusted response is not absence or success. It stops the gate
before mutation. Passing preflight does not widen scope.

## Exactly authorized execution

Only after every preflight item passes, the fresh execution session may
perform this single ordered transaction:

1. Fast-forward remote `main`, without force, to the exact Recovery R1
   authority commit. Push no unrelated path, ref or commit.
2. Create exact tag `release-tooling-v1.0.0-c2` at accepted tooling candidate
   `faef8435322c9096df09b56969662411f17356ea`; push only that tag; create or
   confirm an active exact-name tag ruleset denying update and deletion with
   no bypass actors; and read back the tag object, commit, tree and rules.
3. Confirm repository Actions artifact/log retention is exactly one day. It
   may be reduced to one day only if trusted read-back proves a different
   value. No other repository, Actions, environment, branch, release or
   security setting may change.
4. Dispatch `.github/workflows/release-recovery-v1.0.0.yml` exactly once from
   `refs/tags/release-tooling-v1.0.0-c2` with only these inputs:

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
   unsigned uncompressed workflow artifact named
   `pscan-v1.0.0-c2-unsigned-candidate` retained for one day.
6. Observe that single run to a terminal state and read back its exact
   repository, workflow, ref, SHA, inputs, job, steps, runner, artifact,
   digest, retention, deployment and release facts. The exact unsigned
   artifact may be downloaded only for bounded independent verification and
   may not be copied to a release, public location or consuming project.
7. Record the complete pass or fail-closed outcome locally under PSCAN-06. Do
   not push post-run evidence, rerun, repair, retry, mutate or promote any
   result under this authority. Return to the owner for the next exact gate.

Every step is sequential and conditional. Any failure, timeout, cancellation,
unexpected artifact, missing proof, unexpected deployment, signing-job
eligibility or identity mismatch stops the transaction. Recovery R1 permits
at most one workflow dispatch in total and no second execution or dispatch,
even if the first appears transient or produces no artifact.

## Explicitly forbidden

- No product, workflow, build, test, schema, verifier or release-code change.
- No detached/worktree execution, branch repair, force push, history rewrite,
  tag movement/deletion/recreation, broad tag pattern, bypass actor or change
  to either locked existing tag.
- No creation, update or use of `PSCAN_RELEASE_C2_GATE`; no signing environment
  approval or deployment; the `sign-attest-and-draft` job must be skipped and
  ineligible.
- No signing, keyless OIDC request, signature, trusted-root action,
  attestation, draft or non-draft release, release asset, publication or
  immutable-release finalization.
- No credential or signing-key creation/handling, paid runner, paid storage,
  overage, subscription, provider commitment or spend above USD 0.
- No extra workflow run, rerun, retry, correction, tag, branch, PR, issue,
  repository, remote resource or setting change.
- No scanner adoption, consuming-project data, receipt, deployment,
  production, cross-project reporting, TruffleHog or real secret material.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

## Approval-session exclusions

This session imports the immutable failure record, creates this Recovery R1
authority record and synchronizes PSCAN-06 living governance only. It performs
no live remote preflight, remote read/write, push, tag or ruleset operation,
setting change, workflow dispatch, Docker/image/container action, network
acquisition, build, artifact creation/download, credential action, signing,
attestation, draft, release or publication. It creates no remote resource, key
or paid commitment. PSCAN-06 remains open and unaccepted overall; no successor
is selected or activated.
