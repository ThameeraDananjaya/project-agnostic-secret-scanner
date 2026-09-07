# PSCAN-06 Correction C2 main-branch execution Recovery R2 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Recovery: `R2`
- Gate: build-only actual-Linux and genuine-Docker proof
- Event: owner approval of one bounded post-repair execution recovery
- Owner direction, received 2026-09-08:
  1. allow GitHub-hosted Docker when the required task cannot be completed with
     the Docker installation on the owner's PC;
  2. authorize one `PSCAN-06 Recovery R2` attempt; and
  3. check GitHub access and provide the activation command if access is not
     valid.
- Approval session: `01a0737f-d3fe-7c93-9d3e-db1f631c3df9`
- Required starting authority commit before this bundle:
  `31997817628cbcff8fa921ec24906f24cbd960e1`
- Required starting authority tree:
  `da4f6f8cab8ad915af6fe85bde3498c3812bdb2b`
- Recovery-authority commit: this authority-bundle commit
- Approval-session checkout: exact saved-project root on branch `main`
- Pre-change tracked, staged and untracked Git status: clean
- Preserved ignored workspace material: `graphify-out/**`; excluded from
  product authority, staging and this commit
- State after commit: owner-approved; not claimed for execution; unexecuted
- Workflow dispatches under the original gate, Recovery R1 and Recovery R2:
  zero
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

This is a bounded continuation inside the already activated PSCAN-06
Correction C2. It does not create, select or activate another PSCAN task. This
session records the Recovery R2 authority and the requested read-only GitHub
authentication result only. It does not claim or execute the recovered gate.

## Preserved Recovery R1 failure and repair

Recovery R1 ran in genuinely fresh saved-project session
`01a07376-04a7-7aa1-8421-7d4e82a13465` from the correct root and branch but
failed closed at mandatory preflight item 1 because
`build/gitleaks/build.ps1` and `build/gitleaks/collect-licenses.ps1` had
noncanonical mixed working-tree line endings. It stopped before every remote,
workflow, Docker, artifact, signing and publication action.

The immutable Recovery R1 failure is:

`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R1-FAILURE-001.md`

SHA-256:

`1CC095F7FD16CF6CD8F004371D42F38BB24AFBEAF4F059B99992CFBB8B56BAF6`

The two projections were subsequently restored from their exact committed Git
objects without changing the source blobs or index. The immutable repair
record is:

`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-SOURCE-TRUST-REPAIR-001.md`

SHA-256:

`1044AC45D2D9ADE12DC6432AD69D9F3008F9D6D5397FEBDB14019C747BB910D1`

Post-repair validation then passed the complete tracked-byte audit, all 12
hostile source-trust cases and the Windows no-Docker
image-admission/process-boundary matrix with zero Docker calls. Its immutable
record is:

`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-LOCAL-POST-REPAIR-VALIDATION-001.md`

SHA-256:

`D953FDA39CECF4CD097C372568DCA7DAB944F725715720A3D44E20EC48DCCB92`

Recovery R2 does not reinterpret Recovery R1 as a pass. It permits one new
attempt only from this separately committed authority after every current
preflight fact passes.

## Why the proof uses GitHub-hosted Docker

Local Docker on the Windows PC cannot by itself prove the exact GitHub workflow
identity, standard `ubuntu-24.04` runner identity, workflow inputs, run/job/step
facts, one-day artifact retention or remote artifact read-back required by the
existing proof gate. Recovery R2 therefore uses the already committed GitHub
Actions build job for the actual-Linux and genuine-Docker proof. It does not
authorize starting, stopping or reconfiguring Docker on the owner's PC.

The GitHub-hosted path remains conditional on trusted current proof that the
standard runner, artifact transfer and all repository operations cost exactly
USD 0. Any billing, storage, overage, paid-runner or provider-commitment risk
stops the gate before mutation.

## Controlling and preserved authority

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Original build-only proof-gate authority:
  `evidence/PSCAN-06/CORRECTION-C2-BUILD-ONLY-LINUX-DOCKER-PROOF-GATE-AUTHORITY.md`
- Original gate-authority SHA-256:
  `2A84E7109F151E55399870FB36D64AEDB2213E123976779CD1E9648AC1B9DAC2`
- Recovery R1 authority:
  `evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R1-AUTHORITY.md`
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
role remains immutable. Recovery R2 changes only the retry state after the
local projection repair. It weakens or waives no original condition.

## Requested GitHub access check

This approval session performed only the owner's requested read-only
authentication check:

- `gh` executable: `C:/WINDOWS/system32/gh.exe`
- configured origin:
  `https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner.git`
- command: `gh auth status --hostname github.com`
- result: nonzero; account `ThameeraDananjaya` is active but the stored token is
  invalid
- remote mutation, repository query and workflow dispatch: not attempted

The owner must re-authenticate outside this authority bundle. The fresh
execution session must repeat authentication validation and the complete
remote preflight; this result cannot be reused as a current pass.

## Required fresh execution context

The recovered gate must run in a genuinely fresh Codex task configured as
follows before its first repository or remote action:

1. Use the saved Project-Agnostic Secret Scanner project directly, not a Codex
   worktree, clone or detached checkout.
2. The exact root must be
   `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`.
3. The checkout must be branch `main` at this committed Recovery R2 authority
   bundle and its exact tree, with direct parent
   `31997817628cbcff8fa921ec24906f24cbd960e1`.
4. Use `gpt-daybreak-blue-latest` with reasoning effort `xhigh` and verify both
   values from the task/session metadata.
5. Tracked, staged and untracked state must be clean. Preserved ignored
   `graphify-out/**` is known and excluded from product authority, staging,
   remote evidence and execution decisions. No other unexplained ignored path
   is permitted.
6. If the task starts detached, off `main`, at another root or identity, with
   invalid GitHub authentication, or in a dirty or ambiguous state, stop
   without branch repair, cleanup, authority copying or remote mutation.

This approval session is not the fresh execution session. No execution task or
successor starts automatically from the authority commit.

## Mandatory fresh preflight

Before mutation, the fresh task must fail closed unless current local and
remote read-only evidence proves all of the following:

1. The execution context above is exact, including root, branch, Recovery R2
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
   availability, account budget and stop-usage control are trusted and permit
   the approved run for USD 0.
5. Locked tags `v1.0.0` and `release-tooling-v1.0.0-c1` still resolve exactly
   to their recorded objects and retain active update/deletion denial with no
   bypass. The C2 tag is absent locally and remotely.
6. No workflow run, C2 tag or artifact consumed the original gate, Recovery R1
   or Recovery R2 authority; no other PSCAN task is active; PSCAN-07 remains
   proposed and unselected; PSCAN-08 remains inactive and ineligible.

Missing, stale, permission-denied, rate-limited, contradictory, ambiguous or
otherwise untrusted evidence stops the gate before mutation. Passing preflight
does not widen scope.

## Exactly authorized execution

Only after every preflight item passes, the fresh task may perform this single
ordered transaction:

1. Fast-forward remote `main`, without force, to the exact Recovery R2
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
eligibility or identity mismatch stops the transaction. Recovery R2 permits at
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
  overage, subscription, provider commitment or spend above USD 0.
- No extra workflow run, rerun, retry, correction, tag, branch, PR, issue,
  repository, remote resource or setting change.
- No scanner adoption, consuming-project data, receipt, deployment,
  production, cross-project reporting, TruffleHog or real secret material.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

## Approval-session actions and exclusions

This session records the Recovery R2 authority and the requested read-only
`gh auth status` failure. It performs no login, token handling, repository API
query, push, tag or ruleset operation, setting change, workflow dispatch,
Docker/image/container action, network acquisition, build, artifact action,
signing, attestation, draft, release or publication. It creates no remote
resource, key or paid commitment.

PSCAN-06 remains open and unaccepted overall. PSCAN-07 remains proposed and
unselected. PSCAN-08 remains inactive and ineligible. No successor is selected
or activated.
