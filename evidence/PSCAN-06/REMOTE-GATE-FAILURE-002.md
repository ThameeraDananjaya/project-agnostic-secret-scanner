# PSCAN-06 Remote Gate Failure 002

## Outcome

The owner approved the exact
`PSCAN-06 C1 ITERATION-005 BUILD-ONLY LINUX PROOF GATE` for candidate
`3fb7592889820fa2739a4a53588e073689621809` on 2026-09-04. The approved
remote setup completed, but the single build-only workflow run failed closed
before dependency acquisition because a fresh GitHub-hosted Ubuntu runner did
not already contain the digest-pinned Docker image required by the
pre-acquisition CRLF regression. PSCAN-06 remains unaccepted and open. This
record authorizes no correction, rerun, signing, release or successor work.

## Exact remote identity and controls

- Repository: `ThameeraDananjaya/project-agnostic-secret-scanner`
- Visibility: public
- Default branch: `main`
- Remote `main`: `3523e4409ebc53cc1931e3dcaf7d1eea74bb15ab`
- Locked product tag `v1.0.0`:
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Locked correction-tooling tag `release-tooling-v1.0.0-c1`:
  `3fb7592889820fa2739a4a53588e073689621809`
- Correction-tooling tree:
  `7de9dc4c5725bf38cc80aa734e3c2bac0abd9762`
- Tooling-tag ruleset ID: `22241330`
- Tooling-tag rules: active deletion and update denial, no bypass actors,
  current user cannot bypass
- Workflow: `.github/workflows/release-recovery-v1.0.0.yml`
- Workflow ID: `349829683`
- Workflow run ID: `33829598255`
- Workflow run URL:
  `https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/actions/runs/33829598255`
- Trigger: `workflow_dispatch`
- Run source ref: `refs/tags/release-tooling-v1.0.0-c1`
- Run source commit: `3fb7592889820fa2739a4a53588e073689621809`
- Locked product input:
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Release-version input: `v1.0.0`
- Run start: `2026-09-04T02:27:40Z`
- Run completion: `2026-09-04T02:28:12Z`
- Run conclusion: `failure`

Before dispatch, read-back confirmed that repository variable
`PSCAN_RELEASE_C1_GATE` was absent. The only repository Actions variable was
the historical `PSCAN_RELEASE_GATE=PSCAN-06-SIGNING-APPROVED`. Therefore the
Correction C1 signing job was ineligible.

## Successful gates before failure

The build job passed these steps:

1. immutable dual-identity invocation validation;
2. credential-free checkout of the exact tooling revision;
3. locked product tag and tree verification;
4. exact committed release-entrypoint materialization and raw-object
   verification; and
5. all twelve hostile Git index/checkout source-state rejection cases.

## Failure evidence and diagnosis

The next step, `Prove CRLF checkout cannot corrupt Docker shell payloads`,
failed at `build/release/test-crlf-shell-payloads.ps1:83` with:

```text
Pinned image is required for the CRLF shell-payload regression
```

The script performs `docker image inspect` and rejects when the pinned image is
absent. The workflow intentionally runs this test before the network-enabled
`Acquire pinned public inputs after cache canary` step, while `acquire.ps1`
permits the image pull only in that later acquisition step. A fresh
`ubuntu-24.04` GitHub-hosted runner did not contain
`golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`.
The current workflow therefore has an image-bootstrap ordering contradiction
and cannot complete on the proved fresh-runner state.

This is not treated as a transient failure. The run stopped before the Linux
positive/wrong-owner/read-only cache proof, dependency acquisition, either
network-disabled build, byte comparison or artifact transfer.

## Negative-state read-back

- Build job: `failure`
- Signing/attestation/draft job: `skipped`
- Workflow artifacts: zero
- GitHub releases and draft releases: zero
- `release-v1` deployments from this run: zero
- No signature, attestation, draft or public release was created.
- No signing variable, credential, signing key, paid runner, paid capability,
  TruffleHog material or successor task was created or enabled.
- Both immutable tags and their no-bypass rules remain unchanged.

## Required next decision

Any recovery requires a new bounded PSCAN-06 correction that resolves the
pinned-image bootstrap without weakening the rule that every networked input
is explicit and digest-bound and without moving or deleting either locked tag.
The correction must define and independently validate a new immutable tooling
identity before proposing another remote run. No correction, new tooling tag,
settings change, workflow run, signing, attestation, draft release,
publication or PSCAN-07 work may begin without a new exact owner decision.
