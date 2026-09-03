# PSCAN-06 Reserved Owner Gate

## Current state

Candidate `a13c28fe7273bc8dc6545f97966a02889524eb4c` passed all locally authorized
implementation, reproducibility, test and review requirements. The intended
public repository does not exist, no remote is configured, and no settings,
workflow, tag, signature, attestation or draft-release action has occurred.
Those missing facts prevent PSCAN-06 acceptance and closeout.

## Consolidated exact owner decision block

```text
OWNER DECISION REQUIRED - PSCAN-06 REMOTE AND SIGNING GATE

Approve or reject this exact bounded, staged action:

1. Re-run read-only zero-cost and account/repository availability preflight.
2. Create only the public GitHub repository
   ThameeraDananjaya/project-agnostic-secret-scanner.
3. Configure and read back only these controls: default GITHUB_TOKEN read-only;
   required full-SHA actions only; release-v1 protected environment; exact
   main and v1.0.0 tag rules; immutable releases; one-day workflow artifact
   retention; zero-spend/budget guard; and repository variable
   PSCAN_RELEASE_GATE=PSCAN-06-SIGNING-APPROVED.
4. Add that repository as the only Git remote and push the exact local PSCAN-06
   history after re-reading and recording its commit and tree identities.
5. Create the exact v1.0.0 tag at the recorded PSCAN-06 candidate and invoke
   .github/workflows/release.yml only for that exact tag, commit and version.
6. Permit exactly one keyless GitHub-OIDC signing/attestation run and creation
   of one draft release containing all assets. Do not publish the release.
7. Read back repository ID/settings/rules/environment, workflow identity and
   claims, artifact digests, signature bundle, attestations and draft asset set;
   independently verify them and fail closed on any discrepancy.

Maximum approved cost: USD 0. Abort before any billable condition, provider
commitment, publication, credential creation or unsupported state.

This approval does not publish v1.0.0, does not select/activate/claim PSCAN-07
or PSCAN-08, and does not authorize TruffleHog or any other remote action.

Reply exactly one standalone line:
APPROVE PSCAN-06 REMOTE AND SIGNING GATE
or
REJECT PSCAN-06 REMOTE AND SIGNING GATE
```

An approval must be acted on only after a fresh action-time preflight. An
ambiguous or differently worded reply grants no authority.
