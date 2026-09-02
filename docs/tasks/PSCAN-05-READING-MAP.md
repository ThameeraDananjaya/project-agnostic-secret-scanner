# PSCAN-05 Reading Map

Read these authorities completely, in order, before claiming PSCAN-05 in a
fresh implementation session:

1. `AGENTS.md`
2. `CONSTITUTION.md`
3. `docs/spec/PASS-OUTCOME-SPEC-001.md`
4. `docs/spec/TRACEABILITY.md`
5. `docs/tasks/TRACKER.md`
6. `docs/tasks/PSCAN-05.md`
7. `docs/governance/TASK-LIFECYCLE.md`
8. `docs/governance/OWNER-GATES.md`
9. `docs/governance/EVIDENCE-AND-CLOSEOUT.md`
10. `docs/governance/RUNBOOK-REQUIREMENTS.md`
11. `docs/decisions/DEC-001-SCHEMA-OWNERSHIP.md`
12. `docs/decisions/DEC-002-OUTCOME-BASED-CONTRACT.md`
13. `docs/architecture/ARCHITECTURE.md`
14. `docs/architecture/SOURCE-LAYOUT.md`
15. `docs/architecture/GITLEAKS-ADAPTER.md`
16. `docs/architecture/ARTIFACT-NORMALIZATION.md`
17. `SECURITY.md`
18. `docs/security/THREAT-MODEL.md`
19. `docs/security/GITLEAKS-BOUNDARY.md`
20. `docs/security/ARTIFACT-BOUNDARY.md`
21. `docs/validation/VALIDATION-PLAN.md`
22. `docs/validation/GITLEAKS-COVERAGE-MAP.md`
23. `docs/validation/PSCAN-04-VALIDATION.md`
24. `docs/release/LICENSING.md`
25. `contracts/scan-request/schema-1.0.json`
26. `contracts/scan-outcome/schema-1.0.json`
27. `contracts/global-revocation/schema-1.0.json`
28. `contracts/release-manifest/schema-1.0.json`
29. `docs/tasks/PSCAN-02.md`
30. `evidence/PSCAN-02/ACTIVATION.md`
31. `evidence/PSCAN-02/IMPLEMENTATION.md`
32. `evidence/PSCAN-02/ACCEPTANCE.md`
33. `docs/tasks/PSCAN-04.md`
34. `evidence/PSCAN-04/ACTIVATION.md`
35. `evidence/PSCAN-04/PREFLIGHT.md`
36. `evidence/PSCAN-04/IMPLEMENTATION.md`
37. `evidence/PSCAN-04/VALIDATION.md`
38. `evidence/PSCAN-04/REVIEW.md`
39. `evidence/PSCAN-04/ACCEPTANCE.md`
40. `evidence/PSCAN-04/CLOSEOUT.md`
41. `evidence/PSCAN-05/ACTIVATION.md`

Then verify the exact activation commit, branch and clean Git status; the
activated allowed and forbidden paths; the absence of another selected,
activated or claimed task; and every reserved owner gate. Inspect every tracked
file under the allowed implementation paths before editing. Treat all
consuming-project policy, allowlist, receipt, revocation and schema instances as
untrusted caller inputs that must never become repository authority or retained
state.

Before any download, dependency change, build, scanner execution or pinning,
perform and record a fresh read-only primary-source preflight for every proposed
schema, canonicalization, signature/reference-verification or cryptographic
dependency and licence; the existing Gitleaks source, binary, configuration,
rules and detector behavior; and the Go toolchain/module support and licence
facts. Do not introduce key custody, receipt issuance or signing behavior.

Stop rather than claim or implement if an authority conflicts, DEC-001 cannot
remain consistent with PASS-OUTCOME-SPEC-001, the working tree contains
unexplained changes, the activation commit cannot be established as the
implementation base, a policy layer can weaken a higher layer, unknown or
ambiguous schema semantics can be admitted, credential-class material can be
allowlisted, receipt/revocation reference integrity cannot be verified without
custody, project isolation cannot be proved, or any missing, stale,
unsupported or untrustworthy condition would otherwise be treated as pass.
