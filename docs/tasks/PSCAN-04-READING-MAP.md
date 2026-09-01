# PSCAN-04 Reading Map

Read these authorities completely, in order, before claiming PSCAN-04 in a
fresh implementation session:

1. `AGENTS.md`
2. `CONSTITUTION.md`
3. `docs/spec/PASS-OUTCOME-SPEC-001.md`
4. `docs/spec/TRACEABILITY.md`
5. `docs/tasks/TRACKER.md`
6. `docs/tasks/PSCAN-04.md`
7. `docs/governance/TASK-LIFECYCLE.md`
8. `docs/governance/OWNER-GATES.md`
9. `docs/governance/EVIDENCE-AND-CLOSEOUT.md`
10. `docs/governance/RUNBOOK-REQUIREMENTS.md`
11. `docs/decisions/DEC-001-SCHEMA-OWNERSHIP.md`
12. `docs/decisions/DEC-002-OUTCOME-BASED-CONTRACT.md`
13. `docs/architecture/ARCHITECTURE.md`
14. `docs/architecture/GITLEAKS-ADAPTER.md`
15. `docs/architecture/SOURCE-LAYOUT.md`
16. `SECURITY.md`
17. `docs/security/THREAT-MODEL.md`
18. `docs/security/GITLEAKS-BOUNDARY.md`
19. `docs/validation/VALIDATION-PLAN.md`
20. `docs/validation/GITLEAKS-COVERAGE-MAP.md`
21. `docs/release/LICENSING.md`
22. `contracts/scan-request/schema-1.0.json`
23. `contracts/scan-outcome/schema-1.0.json`
24. `docs/tasks/PSCAN-10.md`
25. `evidence/PSCAN-10/ACTIVATION.md`
26. `evidence/PSCAN-10/PREFLIGHT.md`
27. `evidence/PSCAN-10/IMPLEMENTATION.md`
28. `evidence/PSCAN-10/VALIDATION.md`
29. `evidence/PSCAN-10/REVIEW.md`
30. `evidence/PSCAN-10/ACCEPTANCE.md`
31. `evidence/PSCAN-10/CLOSEOUT.md`
32. `docs/tasks/PSCAN-03.md`
33. `evidence/PSCAN-03/MATERIAL-GAP-001.md`
34. `evidence/PSCAN-03/CORRECTION-C1.md`
35. `evidence/PSCAN-03/REVIEW-C1.md`
36. `evidence/PSCAN-03/CLOSEOUT-REJECTED.md`
37. `evidence/PSCAN-04/ACTIVATION.md`

Then verify the exact activation commit, branch and clean Git status; the
activated allowed and forbidden paths; the absence of another selected,
activated or claimed task; and every reserved owner gate. Inspect every tracked
file under the allowed implementation paths before editing. Treat PSCAN-03 and
correction C1 only as rejected historical defect evidence and do not reuse them
as implementation or acceptance authority.

Before any download, dependency change, build, scanner execution or pinning,
perform and record a fresh read-only primary-source preflight for the exact
archive/OCI parsing choices and licences; the existing Gitleaks source, binary,
configuration/rules and detector behavior; and the Go toolchain/module support
and licence facts. Candidate artifacts are data only and must never be executed,
built or used to install dependencies.

Stop rather than claim or implement if an authority or source conflicts, a
required primary source is unavailable, the working tree contains unexplained
changes, the activation commit cannot be established as the implementation
base, raw classification or byte-bound normalization cannot be proved, an
output/residual path cannot be placed behind the redaction and cleanup
boundaries, authoritative network/credential isolation cannot be proved, or any
missing, stale or unsupported condition would otherwise be treated as pass.
