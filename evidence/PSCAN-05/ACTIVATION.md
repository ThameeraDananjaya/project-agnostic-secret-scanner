# PSCAN-05 Activation Evidence

- Task: `PSCAN-05`
- Event: exact owner activation in an activation-only session
- Owner command: `ACTIVATE PSCAN-05`
- Date: `2026-09-02`
- Branch: `main`
- Parent commit: `1f0890878518de32a55ceb8d7b97430c4f3d2f2b`
- Activation commit: this activation-bundle commit
- State after commit: activated, unclaimed, not implemented
- Pre-change Git status: clean
- Controlling authority: `PASS-OUTCOME-SPEC-001`
- PASS-OUTCOME-SPEC-001 SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Accepted PSCAN-04 closeout parent: exact activation parent above
- PSCAN-04 closeout evidence SHA-256:
  `62C45A85C02697E1C2878083277A56A3D810C3228712A0EA822AFF2F81B3BE2A`
- PSCAN-04 acceptance evidence SHA-256:
  `DBD2B2B6238FAEDE70BA84B3FD5AA01F162FA944281352F51EE56F934ACCAE77`
- DEC-001 SHA-256:
  `DF73FFB83071277765535B8F1D9DB95C41F35B695B167B0C2958FA9F32A7F1B6`
- Required implementation reading order:
  `docs/tasks/PSCAN-05-READING-MAP.md`
- Allowed implementation paths: exactly those listed in
  `docs/tasks/PSCAN-05.md`
- Forbidden implementation scope: exactly that listed in
  `docs/tasks/PSCAN-05.md`
- Schema, canonicalization, signature/reference-verification, cryptographic and
  licence preflight: required fresh from primary sources in the implementation
  session before dependency intake, build, execution or pinning
- Existing Gitleaks, rules and Go toolchain/module preflight: required fresh
  from primary sources in the implementation session before scanner or test
  execution
- Dependency, action, toolchain, scanner, rule or configuration pins
  introduced: none
- Product code, schemas, policy behavior, allowlist behavior, verification,
  workspace behavior, fixtures, tests, build behavior or workflows created or
  changed: none
- Scanner, dependency, toolchain or test download, build or execution: none
- Receipt issuance, signing, key custody or evidence storage introduced: none
- Gitleaks GitHub Action or TruffleHog material introduced or assessed: none
- Remote actions or resources: none
- Credentials or signing keys created or handled: none
- Spending or paid capability: none
- Consuming-project source, policy, allowlist, receipt, revocation, identity or
  data introduced: none
- Successor selected or activated: none
- Successor policy: PSCAN-06 and PSCAN-07 remain proposed and unselected;
  PSCAN-08 remains inactive, unselected, technically gated and AGPL
  owner-gated

This record authorizes only a fresh PSCAN-05 implementation session from the
exact activation commit. That session must complete the ordered reading map,
claim only PSCAN-05, perform the required live primary-source preflights before
intake or execution, modify only activated allowed paths, independently validate
policy precedence, narrow allowlist exceptions, independent schema-family
compatibility, custody-neutral receipt/revocation reference verification and
project isolation, and stop after PSCAN-05 closeout with every successor
unselected.

This activation does not authorize implementation in this session, dependency
or scanner intake, downloads, builds, tests, workflows, remote creation or
publication, GitHub settings, signing, credentials, spending, consuming-project
integration, receipt issuance or key custody, PSCAN-06 work, PSCAN-08 work,
deployment, production or go-live.
