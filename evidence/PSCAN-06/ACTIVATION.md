# PSCAN-06 Activation Evidence

- Task: `PSCAN-06`
- Event: exact owner activation in an activation-only session
- Owner command: `ACTIVATE PSCAN-06`
- Date: `2026-09-02`
- Branch: `main`
- Parent commit: `30d27856bdfd404fc190501be854b78bc5147f1f`
- Activation commit: this activation-bundle commit
- State after commit: activated, unclaimed, not implemented
- Pre-change Git status: clean
- Controlling authority: `PASS-OUTCOME-SPEC-001`
- PASS-OUTCOME-SPEC-001 SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Accepted PSCAN-05 closeout parent: exact activation parent above
- PSCAN-05 independent review SHA-256:
  `824B7F6F4B751E9DB077922298BE7578347EFEA99CBD63E395DEF76A5AB9E5EA`
- PSCAN-05 acceptance evidence SHA-256:
  `7097B52B5E378AB00B1A75B577E6A06994404E854AC0A15F766EB3C5119E2702`
- PSCAN-05 closeout evidence SHA-256:
  `115D05205E6B55BC9F9B9C8E7ED60DAD7085D64A42739517CFFAAEB9E7BBF2FB`
- Required implementation reading order:
  `docs/tasks/PSCAN-06-READING-MAP.md`
- Allowed implementation paths: exactly those listed in
  `docs/tasks/PSCAN-06.md`
- Forbidden implementation scope: exactly that listed in
  `docs/tasks/PSCAN-06.md`
- Official Go, Gitleaks, Cosign/Sigstore, GitHub Actions, Releases, immutable
  releases, artifact attestations, OIDC, settings, plan/cost, SBOM and licensing
  preflight: required fresh from primary sources in the implementation session
  before intake, build, execution, pinning or workflow change
- Read-only account, repository availability, visibility, permission,
  protection/ruleset, capability and zero-cost preflight: required fresh before
  proposing any exact remote action
- Separate exact owner approval after preflight: required before remote
  creation, push, settings change, workflow enablement or signing dry run
- Dependency, action, toolchain, scanner, rule or configuration pins
  introduced: none
- Product code, schemas, verifier behavior, build tooling, workflows, manifests,
  SBOMs, licences/notices, security/validation behavior or runbooks created or
  changed: none
- Scanner, dependency, toolchain, build, test, workflow or signing execution:
  none
- Gitleaks GitHub Action or TruffleHog material introduced or assessed: none
- Remote actions or resources: none; no Git remote is configured
- Credentials or signing keys created or handled: none
- Spending or paid capability: none
- Consuming-project source, policy, allowlist, receipt, revocation, identity or
  data introduced: none
- Successor selected or activated: none
- Successor policy: PSCAN-07 remains proposed and unselected; PSCAN-08 remains
  inactive, unselected, technically gated and AGPL owner-gated

This record authorizes only a fresh PSCAN-06 implementation session from the
exact activation commit. That session must complete the ordered reading map,
claim only PSCAN-06, perform the required live primary-source and read-only
account/capability preflights, modify only activated allowed paths,
independently validate reproducible Windows amd64 and Linux amd64 assets,
complete supply-chain evidence and offline verification, and stop after
PSCAN-06 closeout with every successor unselected. It must stop for separate
owner approval before any remote creation, push, settings change, workflow
enablement or signing dry run.

This activation does not authorize implementation in this session, dependency,
action or scanner intake, downloads, builds, tests, workflow execution, remote
creation or publication, GitHub settings, signing, credentials, spending,
consuming-project integration, final `v1.0.0` release publication, PSCAN-07
work, PSCAN-08 work, deployment, production, legal/compliance or go-live.
