# Owner Gates

## Approved standing defaults

- Local Git repository bootstrap on branch `main`.
- Product name and future remote name:
  `project-agnostic-secret-scanner`.
- Intended GitHub account: `ThameeraDananjaya`.
- Final visibility: public.
- Product licence: MIT.
- Initial spending ceiling: zero.
- Intended final release: signed `v1.0.0`.
- Signing design: Cosign keyless signing from the exact approved GitHub release
  workflow.
- Release maintainer and final human approver: `ThameeraDananjaya`.
- Runbook: plain English Markdown with Mermaid illustrations and complete input
  and output contracts.

## Still reserved for action-time approval

The standing defaults do not authorize these actions in advance:

- Create the public GitHub repository or any other remote resource.
- Push, publish, sign, or create a release.
- Change GitHub Actions, branch, ruleset, environment, release, or protection
  settings.
- Incur any cost, paid plan, overage, storage charge, or provider commitment.
- Create or handle credentials, long-lived signing keys, or project receipt
  keys.
- Adopt the scanner in a consuming project or configure deployment/production
  gates.
- Enable cross-project reporting.
- Finalize the legal/compliance retention schedule.

Public repository creation remains scheduled for PSCAN-06 after a read-only
capability, settings, identity, and cost preflight and the owner's separate
approval.

## TruffleHog gate

PSCAN-08 remains inactive. If accepted tests prove a material Gitleaks coverage
gap, work must stop and return to the owner for a separate technical and AGPL
legal-approval decision before downloading, integrating, distributing, or
enabling TruffleHog. No legal assessment begins before that evidence exists.
