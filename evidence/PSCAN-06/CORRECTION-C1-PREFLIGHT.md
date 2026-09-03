# PSCAN-06 Correction C1 implementation preflight

## Claim and authority

- Task: `PSCAN-06`
- Correction: `C1`
- Session date: `2026-09-04`
- Claimed authority commit:
  `3fb1b0a55dc4f48dd35464c63c768f497efbc89b`
- Authority tree: `544867396910969d20cfd2acd454ac0d69c1d7e4`
- Branch: `main`
- Initial status: clean
- Sole open task: `PSCAN-06`
- Successor selected or activated: none
- Decision: `PROCEED_TO_C1_ONLY`

This fresh session claims only the already-approved PSCAN-06 Correction C1.
It does not select or activate a task and does not claim independent
acceptance.

## Mandatory reading and bounded inspection

The implementing agent read all 62 entries in
`docs/tasks/PSCAN-06-READING-MAP.md` before implementation, including the
controlling contract, correction authority, historical remote failure and all
required code, schema, workflow, test and evidence files. It then inspected
all 35 files tracked under the correction's existing allowed implementation
paths. The narrower path list in
`evidence/PSCAN-06/CORRECTION-C1-AUTHORITY.md` is controlling for this work.

No predecessor contract, prior decision, historical evidence, existing
release-manifest schema, product detector/rule behavior or consuming-project
material is in scope.

## Immutable identities and repository read-back

- Repository: `ThameeraRA/project-agnostic-secret-scanner`
- Visibility: public
- Default branch: `main`
- Locked product tag: `v1.0.0`
- Local and remote product-tag commit:
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Locked product-source tree:
  `217b711ddea51fd0ea7e808edd2e27fdecef8427`
- Proposed correction-tooling tag: `release-tooling-v1.0.0-c1`
- Remote correction-tooling tag before implementation: absent
- Historical failed workflow run: `33709197614`
- Historical failure: Linux runner could not create the bind-mounted module
  cache's `cache` directory; no release/signing/publication result was created.

The product tag will not be moved, recreated or reinterpreted. Correction
tooling will remain a distinct identity and no tag is created locally or
remotely in this implementation session.

## Current primary-source verification

- GitHub documents `github.workflow_ref` as the workflow path plus ref and
  `github.workflow_sha` as the commit SHA of the workflow file:
  <https://docs.github.com/en/actions/reference/workflows-and-actions/contexts>
- Docker documents numeric `--user`, read-only roots, capability dropping and
  tmpfs ownership/mode controls:
  <https://docs.docker.com/reference/cli/docker/container/run>
  and <https://docs.docker.com/engine/storage/tmpfs/>.
- Sigstore's Cosign advisory identifies `v3.1.3` as the corrected release for
  the affected verification issue:
  <https://github.com/sigstore/cosign/security/advisories/GHSA-fx35-mq7g-6g98>.

The existing exact pins remain current and available:

- Go `1.27.1` Linux archive SHA-256:
  `63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445`
- Go `1.27.1` Windows archive SHA-256:
  `a3911b5e0e1b1053f25ed0675f4c1c6aad1e2bfcf253df2b9be4caabd2edd95d`
- Gitleaks `v8.30.1` source commit:
  `83d9cd684c87d95d656c1458ef04895a7f1cbd8e`
- Cosign `v3.1.3` release commit:
  `11926fa5bbbbde47e88fc006b625a17769b743b2`
- Cosign Linux amd64 asset SHA-256:
  `4629c757b7618056f8ddd7e2625ae9fdd94c0372a65049520bc7d9df9efc7f71`
- Cosign Windows amd64 asset SHA-256:
  `9fe59be0eca1271873ce019061335eb1ac419b7059202e797828467ddabe33be`

The existing immutable action pins resolved to their expected commits. Docker
Desktop `4.87.0` exposed a Linux/amd64 engine and already contained the exact
pinned Go build image. Native Windows `pwsh` and Docker were available; native
Go and a general-purpose Linux WSL distribution were not assumed.

## Planned correction boundary

The implementation will:

1. run the acquisition container with a read-only root, all capabilities
   dropped and no-new-privileges;
2. map the Linux host numeric UID/GID and use bounded cache/tmpfs modes;
3. run a write, same-filesystem atomic rename, exact read and delete canary
   before any dependency download and before any success ledger;
4. build twice with the completed cache read-only and networking disabled;
5. materialize the locked product source from its exact Git object while
   building correction tooling from its separate exact commit/tree;
6. add manifest schema `2.0` and fail-closed verification for both mandatory,
   non-interchangeable identities and the exact recovery workflow context; and
7. add the bounded recovery workflow, documentation and adversarial tests.

No dependency, archive, action or source pin changes are required. No remote
mutation, workflow execution, signing, attestation, release creation or
publication is authorized by this session.
