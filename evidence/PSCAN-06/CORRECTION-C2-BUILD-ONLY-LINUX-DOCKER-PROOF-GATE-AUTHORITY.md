# PSCAN-06 Correction C2 build-only Linux and Docker proof gate authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Gate: build-only actual-Linux and genuine-Docker proof
- Event: exact owner approval of the bounded gate described below
- Owner command:
  `APPROVE PSCAN-06 C2 BUILD-ONLY LINUX AND DOCKER PROOF GATE FOR ACCEPTED TOOLING CANDIDATE faef8435322c9096df09b56969662411f17356ea FROM EVIDENCE HEAD b6d341284cf63baa0502fcba319ea8cca3c7eb3a ZERO-SPEND NO-SIGNING NO-ATTESTATION NO-DRAFT NO-PUBLICATION ONE-DAY UNSIGNED WORKFLOW ARTIFACT ONLY`
- Date: `2026-09-05`
- Approval session: `01a072d7-40e1-76a0-ac22-995d48e94d3e`
- Branch: `main`
- Required starting evidence HEAD before this authority bundle:
  `b6d341284cf63baa0502fcba319ea8cca3c7eb3a`
- Required starting evidence tree:
  `983028703e229abb082938d0e4417c506aac0d8c`
- Accepted Correction C2 tooling candidate:
  `faef8435322c9096df09b56969662411f17356ea`
- Accepted tooling tree:
  `3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`
- Candidate relation: direct parent of the required evidence HEAD
- Gate-authority commit: this authority-bundle commit
- Pre-change tracked and untracked Git status: clean
- Preserved ignored workspace material: `graphify-out/**`; excluded from
  product authority, staging and this commit. Graph routing deterministically
  regenerated only its ignored vocabulary cache from the unchanged graph.
- Gate state after commit: owner-approved; not claimed for execution; no live
  remote preflight performed; no push, tag, setting or workflow action taken
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

This approval records a bounded action-time authority bundle only. It does not
execute the gate in the approval session. A genuinely fresh execution session
must begin from the exact committed authority bundle, claim only this gate,
read the complete current PSCAN-06 authority and evidence chain, and pass the
fresh fail-closed preflight below before any remote mutation or workflow run.

## Controlling and preserved authority

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Original Correction C2 authority:
  `evidence/PSCAN-06/CORRECTION-C2-AUTHORITY.md`
- Original C2 authority SHA-256:
  `3CE8E20151993A8FC77890E51C5BF1557F68040CAFF2C3DB1A46237792A75D63`
- Iteration-005 authority:
  `evidence/PSCAN-06/CORRECTION-C2-ITERATION-005-AUTHORITY.md`
- Iteration-005 authority SHA-256:
  `F866EDEDD850DE754C533FB3D9AC280D20C19CE72B9E45C0744FF9ACF64E49A9`
- Iteration-005 local independent acceptance:
  `evidence/PSCAN-06/CORRECTION-C2-ITERATION-005-ACCEPTANCE.md`
- Local-acceptance SHA-256:
  `1C9392E1899D1EB83C2DBDD673C0369042C5F6255D3D0BE9558C55281EDB72A9`
- Local independent verdict: `LOCAL_ACCEPTANCE_PASS_REMOTE_PROOF_OPEN`
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

The candidate, evidence HEAD, every predecessor commit and evidence record,
both locked existing tags, manifest schemas, product/tooling role separation,
image and action pins, workflow, test and release boundaries remain immutable.
This gate supplies missing actual-Linux, genuine-Docker, acquisition, complete-
build, byte-comparison and remote evidence only. A successful run is not by
itself independent acceptance or PSCAN-06 completion.

## Mandatory fresh preflight

Before mutation, the fresh execution session must fail closed unless all of
the following are proved from current local and remote read-only evidence:

1. The exact repository root, `main` branch, gate-authority commit and tree are
   present in an exact clean source-trust-compliant checkout; the accepted
   candidate and evidence HEAD resolve to the identities above with the stated
   parent relation.
2. The workflow at the accepted candidate has exactly the recorded bytes and
   SHA-256, all build-job actions remain full-SHA pinned, the build job retains
   `contents: read`, and its single artifact uses `retention-days: 1` and
   `compression-level: 0`.
3. The signing job remains gated only by exact repository variable
   `PSCAN_RELEASE_C2_GATE=PSCAN-06-C2-SIGNING-APPROVED`; that variable is
   absent. If it exists with any value, stop without changing or dispatching.
4. The authenticated account, repository owner/name and numeric ID, public
   visibility, default branch, origin URL, Actions policy, standard runner
   availability, account budget and stop-usage control are exactly trusted and
   permit the approved run for USD 0. Any billing, storage, overage, larger-
   runner or provider-commitment risk stops the gate.
5. Locked tags `v1.0.0` and `release-tooling-v1.0.0-c1` still resolve exactly
   as above and retain active update/deletion denial with no bypass. The C2 tag
   is absent locally and remotely. Any conflicting, moved, deleted, mutable or
   ambiguous identity stops the gate.
6. No previous C2 run or artifact already consumed this one-run authority, no
   other PSCAN task is active, PSCAN-07 remains proposed and unselected, and
   PSCAN-08 remains inactive and ineligible.

Preflight evidence is not permission to widen scope. A stale, missing,
unsupported, rate-limited, permission-denied or otherwise unknown response is
not absence or success and stops before mutation.

## Exactly authorized execution

After every preflight item passes, the fresh execution session may perform only
this ordered transaction:

1. Fast-forward remote `main`, without force, to the exact gate-authority
   commit descended directly from required evidence HEAD `b6d3412`. Stage or
   push no unrelated path, ref or commit.
2. Create exact tag `release-tooling-v1.0.0-c2` at accepted tooling candidate
   `faef8435322c9096df09b56969662411f17356ea`, push only that tag, create or
   confirm an active exact-name tag ruleset that denies update and deletion
   with no bypass actors, and read back the tag object, commit, tree and rules.
3. Confirm repository Actions artifact/log retention is exactly one day. It
   may be reduced to one day only if the current trusted read-back proves a
   different value; no other repository, Actions, environment, branch, release
   or security setting may change.
4. Dispatch `.github/workflows/release-recovery-v1.0.0.yml` exactly once from
   `refs/tags/release-tooling-v1.0.0-c2` with only these inputs:

   ```text
   product_source_revision=a13c28fe7273bc8dc6545f97966a02889524eb4c
   release_tooling_revision=faef8435322c9096df09b56969662411f17356ea
   release_version=v1.0.0
   ```

5. Permit only the workflow `build` job to run on the standard
   `ubuntu-24.04` runner. This authorizes that job's exact committed actual-
   Linux checks, genuine Docker engine/image/container operations, one bounded
   canonical digest pull if and only if image admission proves conclusive
   absence, pinned public dependency acquisition, two network-disabled builds,
   every-byte comparison and transfer of exactly one unsigned uncompressed
   workflow artifact named `pscan-v1.0.0-c2-unsigned-candidate` for one day.
6. Observe the single run to a terminal state and read back its exact
   repository, workflow, ref, SHA, input, job, step, runner, artifact, digest,
   retention, deployment and release facts. Downloading that exact unsigned
   artifact for bounded independent verification is read-only and permitted;
   it may not be copied to a release, public location or consuming project.
7. Record the complete pass or fail-closed outcome locally under PSCAN-06. Do
   not push post-run evidence, rerun, repair, retry, mutate or promote any
   result under this authority. Return to the owner for the next exact gate.

Every step is sequential and conditional. A failure, timeout, cancellation,
unexpected artifact, missing proof, unexpected deployment, signing-job
eligibility or identity mismatch stops the transaction. This authority permits
no second dispatch, even if the first run appears transient or produces no
artifact.

## Explicitly forbidden

- No product, workflow, build, test, schema, verifier or release-code change.
- No force push, history rewrite, tag movement/deletion/recreation, broad tag
  pattern, bypass actor or change to either locked existing tag.
- No creation, update or use of `PSCAN_RELEASE_C2_GATE`; no signing environment
  approval or deployment; the `sign-attest-and-draft` job must be skipped and
  ineligible.
- No signing, keyless OIDC request, signature, Sigstore trusted-root action,
  attestation, draft or non-draft release, release asset, publication or
  immutable-release finalization.
- No credential or signing key creation/handling, paid runner, paid storage,
  overage, subscription, provider commitment or spend above USD 0.
- No extra workflow run, rerun, retry, correction, tag, branch, PR, issue,
  repository, remote resource or setting change.
- No scanner adoption, consuming-project data, receipt, deployment,
  production, cross-project reporting, TruffleHog or real secret material.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

## Approval-session exclusions

This approval session creates only this bounded authority bundle and
synchronizes living PSCAN-06 state. It performs no live remote preflight,
remote read/write, push, tag or ruleset operation, setting change, workflow
dispatch, Docker/image/container action, network acquisition, build, artifact
creation/download, credential action, signing, attestation, draft, release or
publication. It creates no remote resource, key or paid commitment. PSCAN-06
remains open and unaccepted overall; no successor is selected or activated.
