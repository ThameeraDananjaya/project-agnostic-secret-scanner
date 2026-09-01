# PSCAN-10 Reading Map

Read these authorities completely, in order, before claiming PSCAN-10 in a
fresh implementation session:

1. `AGENTS.md`
2. `CONSTITUTION.md`
3. `docs/spec/PASS-OUTCOME-SPEC-001.md`
4. `docs/spec/TRACEABILITY.md`
5. `docs/tasks/TRACKER.md`
6. `docs/tasks/PSCAN-10.md`
7. `docs/governance/TASK-LIFECYCLE.md`
8. `docs/governance/OWNER-GATES.md`
9. `docs/governance/EVIDENCE-AND-CLOSEOUT.md`
10. `docs/governance/RUNBOOK-REQUIREMENTS.md`
11. `docs/decisions/DEC-002-OUTCOME-BASED-CONTRACT.md`
12. `docs/architecture/ARCHITECTURE.md`
13. `docs/architecture/GITLEAKS-ADAPTER.md`
14. `docs/architecture/SOURCE-LAYOUT.md`
15. `SECURITY.md`
16. `docs/security/THREAT-MODEL.md`
17. `docs/security/GITLEAKS-BOUNDARY.md`
18. `docs/validation/VALIDATION-PLAN.md`
19. `docs/validation/GITLEAKS-COVERAGE-MAP.md`
20. `docs/release/LICENSING.md`
21. `docs/tasks/PSCAN-03.md`
22. `evidence/PSCAN-03/MATERIAL-GAP-001.md`
23. `evidence/PSCAN-03/REVIEW-C1.md`
24. `evidence/PSCAN-03/CLOSEOUT-REJECTED.md`
25. `evidence/PSCAN-09/ACCEPTANCE.md`
26. `evidence/PSCAN-09/CLOSEOUT.md`
27. `evidence/PSCAN-10/ACTIVATION.md`

Then verify the exact activation commit, branch and clean Git status; the
activated allowed and forbidden paths; the absence of another selected,
activated or claimed task; and every reserved owner gate. Inspect every tracked
file under the allowed implementation paths before editing and treat rejected
PSCAN-03 material only as historical defect evidence.

Before any download, dependency change, rebuild, scanner execution or pinning,
perform and record a fresh read-only primary-source preflight for the exact
Gitleaks source, release artifacts, checksums, binary, configuration/rules,
detector span/stream/archive behavior and licence, plus the Go toolchain/module
support and licence facts. The separately licensed Gitleaks GitHub Action is
not an allowed substitute or dependency.

Stop rather than claim or implement if an authority or source conflicts, a
required primary source is unavailable, the working tree contains unexplained
changes, the activation commit cannot be established as the implementation
base, a finite detector span or complete streaming path cannot be proved, or
any missing/stale/unsupported condition would otherwise be treated as pass.
