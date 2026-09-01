# PSCAN-09 Activation Evidence

- Task: `PSCAN-09`
- Event: exact owner activation in an activation-only session
- Owner command: `ACTIVATE PSCAN-09`
- Date: `2026-09-01`
- Branch: `main`
- Parent commit: `7b866484da50319c91bd7b67c37637ceef888dab`
- Activation commit: this activation-bundle commit
- State after commit: activated, unclaimed, not implemented
- Pre-change Git status: clean
- Controlling authority during activation: PASS-SPEC-001 remains controlling
  until a fresh PSCAN-09 implementation is independently accepted and closed
- Canonical PASS-SPEC-001 SHA-256:
  `78A1AD9A9ABFDE577A50AE9AA467B5B82F733162C8187ED58A2BC300952B692D`
- External owner source digest: exact match
- PSCAN-03 prerequisite: rejected and closed fail-closed; not accepted; clean
  closeout parent `7b866484da50319c91bd7b67c37637ceef888dab`
- Owner direction bound by this activation: preserve PASS-SPEC-001 and its
  source record as immutable history; decommission controlling authority;
  establish an outcome-based successor permitting exact Git-object enumeration
  and deterministic coverage-proved projections; reconcile documentation and
  downstream planning; perform no scanner implementation
- Required implementation reading order:
  `docs/tasks/PSCAN-09-READING-MAP.md`
- Allowed implementation paths: exactly those listed in
  `docs/tasks/PSCAN-09.md`
- Forbidden implementation scope: exactly that listed in
  `docs/tasks/PSCAN-09.md`
- Dependency, action, toolchain, scanner, rule or configuration pins
  introduced: none
- Product code, schemas, adapters, Git preparation, projections, rules,
  fixtures, tests, build behavior or workflows created or changed: none
- Canonical contract, external source and historical evidence changed: none
- Gitleaks or TruffleHog material downloaded, introduced, built, assessed or
  executed: none
- Remote actions or resources: none
- Credentials or signing keys created or handled: none
- Spending or paid capability: none
- Consuming-project source, policy, identity or data introduced: none
- Successor selected or activated: none
- Successor policy: PSCAN-10, if proposed by the accepted transition, remains
  unselected; PSCAN-04 through PSCAN-07 remain unselected; PSCAN-08 remains
  inactive, unselected and separately technically and legally gated

This record authorizes only a fresh PSCAN-09 documentation/governance
implementation session from the exact activation commit. That session must
complete the ordered reading map, preserve immutable history, modify only
PSCAN-09 allowed paths, independently validate the complete authority
transition and stop after PSCAN-09 closeout with every successor unselected.

This activation does not itself decommission PASS-SPEC-001 and does not
authorize implementation in this session, scanner behavior, dependencies,
workflows, remote creation or publication, GitHub settings, signing,
credentials, spending, consuming-project integration, PSCAN-10 work, PSCAN-04
work or PSCAN-08 work.
