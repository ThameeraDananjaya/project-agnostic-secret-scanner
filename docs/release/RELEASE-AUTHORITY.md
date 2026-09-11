# Release Authority Plan

## Standing direction

- Intended public repository owner: `ThameeraDananjaya`.
- Intended product/repository name: `project-agnostic-secret-scanner`.
- Intended first production-consumable release: signed GitHub `v1.0.0`.
- Product licence: MIT.
- Spending ceiling: zero.
- Release maintainer and final human approver: `ThameeraDananjaya`.
- Signing design: Cosign keyless signing from one exact approved protected
  GitHub release workflow.

These are design defaults, not action authorization.

Correction C1 preserves `v1.0.0` at product-source commit
`a13c28fe7273bc8dc6545f97966a02889524eb4c`. Any recovery must use separately
accepted immutable tooling tag `release-tooling-v1.0.0-c1` and bind its exact
commit, tree, workflow path/ref/SHA and trigger independently from the product
tag/commit/tree. Neither identity may stand in for the other. Creating that tag,
enabling its gate or running the recovery workflow requires a later exact owner
decision.

## Current local authority status

The Correction C1 paragraph above records its original planning boundary. The
locked product tag remains `v1.0.0` at
`a13c28fe7273bc8dc6545f97966a02889524eb4c`; the locked C1 tooling tag remains
at `3fb7592889820fa2739a4a53588e073689621809`; and the protected C2 tooling tag
remains at pre-fix candidate
`faef8435322c9096df09b56969662411f17356ea`.

Recovery R5 consumed its only tag push and workflow dispatch and terminally
failed before Docker or artifact creation. Iteration 006 is independently
accepted locally for its bounded PID correction. Iteration 007 is claimed and
locally implemented only to add schema 2.2 plus proposed tag/ref/workflow
identity `release-tooling-v1.0.0-c2-r6`. It creates no tag and authorizes no
workflow run. Recovery R6 requires independent local acceptance followed by a
separate exact authority and fresh zero-spend remote preflight. Signing,
attestation, draft creation and publication remain separately owner-gated.

## Separation of authorities

| Authority | May do | May not do |
|---|---|---|
| Product owner | Approve publication, legal posture, spend, identity and final release | Convert an unaccepted build into accepted evidence |
| Scanner maintainer | Prepare bounded source and evidence under activated tasks | Publish, sign or change remote controls without action-time approval |
| Release workflow | Build/test/package/sign the exact authorized revision after gates pass | Receive project receipt keys or candidate source from consumers |
| Release approver | Approve exact accepted release invocation | Issue project receipts or approve project deployment |
| Consumer verifier | Acquire and verify exact release into local custody | Trust `latest`, mutable refs or incomplete verification |
| Project receipt issuer | Sign exact project-owned full-release evidence | Use scanner release identity or grant deployment authority |

## Planned release flow

```mermaid
flowchart TD
    Accepted[Accepted local product revision] --> Preflight[Read-only identity, settings, plan, cost, licence and capability preflight]
    Preflight --> OwnerGate{Exact owner approval?}
    OwnerGate -->|no| Stop[Stop; no remote change]
    OwnerGate -->|yes| Remote[Create/configure approved public repository]
    Remote --> Build[Trusted pinned workflow builds exact source]
    Build --> Accept[Cross-platform adversarial acceptance]
    Accept --> Draft[Draft exact v1.0.0 assets together]
    Draft --> Sign[Cosign keyless exact workflow identity]
    Sign --> Verify[Independent manifest, digest, licence, SBOM and bundle verification]
    Verify --> PublishGate{Exact publication approval?}
    PublishGate -->|no| Hold[Hold without publication]
    PublishGate -->|yes| Publish[Publish immutable release where supported]
```

## PSCAN-06 action-time gates

Before any remote creation or settings change, perform a current read-only
preflight of account identity, repository availability, visibility, plan/cost,
Actions, rulesets/protection, immutable releases, artifact attestations, OIDC,
workflow permissions and exact public documentation. Present the exact proposed
changes and any cost. The owner must separately approve the identified action.

PSCAN-01 authorizes none of the following: remote creation, push, PR, settings
change, workflow enablement, signing, credential handling, public artifact,
release or spending.

## Release workflow requirements

- Trusted definition; all third-party actions pinned by full commit SHA.
- Minimal permissions and protected exact workflow/ref identity.
- Clean reproducible build from exact accepted source and pinned toolchain.
- Windows amd64 and Linux amd64 runner plus primary engine assets.
- Generic rules, scanner-owned schemas, release manifest, checksums, bundle,
  SBOM, licences/notices, compatibility, acceptance summary and lifecycle docs.
- Draft-first all-assets-together publication; exact version only; no `latest`.
- Immutable release protection where current verified GitHub capability permits.
- Independent verification of repository, workflow, ref/tag, OIDC issuer,
  manifest signature and every asset before publication.
- Global revocation location, rollback and retirement procedure.

## Separate project boundary

Online acquisition and verification end before network-disabled execution.
Project receipt signing occurs only after a project independently validates the
content-free full-release pass. Receipt credentials are never present in a
scanner execution job. Scanner pass, GitHub check and signed product release do
not approve merge, deployment, production, compliance or go-live.
