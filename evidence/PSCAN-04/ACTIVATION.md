# PSCAN-04 Activation Evidence

- Task: `PSCAN-04`
- Event: exact owner activation in an activation-only session
- Owner command: `ACTIVATE PSCAN-04`
- Date: `2026-09-01`
- Branch: `main`
- Parent commit: `14348b16d9ddb87f337391a4e47dd243f6e53bf2`
- Activation commit: this activation-bundle commit
- State after commit: activated, unclaimed, not implemented
- Pre-change Git status: clean
- Controlling authority: `PASS-OUTCOME-SPEC-001`
- PASS-OUTCOME-SPEC-001 SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Accepted PSCAN-10 closeout parent: exact activation parent above
- PSCAN-10 closeout evidence SHA-256:
  `7E359C144431B4ADF2C233FD89755E03E97A73ADE2EDCDADC18F48D98E4C0F4B`
- PSCAN-10 acceptance evidence SHA-256:
  `196EE03E86FEABE57C13ED95CFB8DEB8601F4E140F1D10CA60B0BA03582AB817`
- Required implementation reading order:
  `docs/tasks/PSCAN-04-READING-MAP.md`
- Allowed implementation paths: exactly those listed in
  `docs/tasks/PSCAN-04.md`
- Forbidden implementation scope: exactly that listed in
  `docs/tasks/PSCAN-04.md`
- Archive/OCI parser and licence preflight: required fresh from primary sources
  in the implementation session before dependency intake, build, execution or
  pinning
- Existing Gitleaks, rules, preparation/proof boundary and Go toolchain/module
  preflight: required fresh from primary sources in the implementation session
  before scanner or test execution
- Dependency, action, toolchain, scanner, parser, rule or configuration pins
  introduced: none
- Product code, schemas, adapters, artifact handling, redaction, cleanup,
  sandbox behavior, fixtures, tests, build behavior or workflows created or
  changed: none
- Scanner, dependency, toolchain or test download, build or execution: none
- Gitleaks GitHub Action or TruffleHog material introduced or assessed: none
- Remote actions or resources: none
- Credentials or signing keys created or handled: none
- Spending or paid capability: none
- Consuming-project source, policy, identity, artifact or data introduced: none
- Successor selected or activated: none
- Successor policy: PSCAN-05 through PSCAN-07 remain proposed and unselected;
  PSCAN-08 remains inactive, unselected, technically gated and AGPL
  owner-gated

This record authorizes only a fresh PSCAN-04 implementation session from the
exact activation commit. That session must complete the ordered reading map,
claim only PSCAN-04, perform the required live primary-source preflights before
intake or execution, modify only activated allowed paths, independently validate
artifact normalization, redaction, cleanup and authoritative offline execution,
and stop after PSCAN-04 closeout with every successor unselected.

This activation does not authorize implementation in this session, dependency
or scanner intake, downloads, builds, tests, workflows, remote creation or
publication, GitHub settings, signing, credentials, spending, consuming-project
integration, PSCAN-05 work, PSCAN-08 work, deployment, production or go-live.
