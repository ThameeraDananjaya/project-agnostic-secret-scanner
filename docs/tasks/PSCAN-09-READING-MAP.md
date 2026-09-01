# PSCAN-09 Reading Map

Read these authorities completely, in order, before claiming PSCAN-09 in a
fresh implementation session:

1. `AGENTS.md`
2. `CONSTITUTION.md`
3. `docs/spec/PASS-SPEC-001-SOURCE.md`
4. `docs/spec/PASS-SPEC-001.md`
5. `docs/spec/TRACEABILITY.md`
6. `docs/tasks/TRACKER.md`
7. `docs/tasks/PSCAN-09.md`
8. `docs/governance/TASK-LIFECYCLE.md`
9. `docs/governance/OWNER-GATES.md`
10. `docs/governance/EVIDENCE-AND-CLOSEOUT.md`
11. `docs/governance/RUNBOOK-REQUIREMENTS.md`
12. `docs/decisions/DEC-001-SCHEMA-OWNERSHIP.md`
13. `docs/architecture/ARCHITECTURE.md`
14. `docs/architecture/GITLEAKS-ADAPTER.md`
15. `docs/architecture/SOURCE-LAYOUT.md`
16. `SECURITY.md`
17. `docs/security/THREAT-MODEL.md`
18. `docs/security/GITLEAKS-BOUNDARY.md`
19. `docs/validation/VALIDATION-PLAN.md`
20. `docs/validation/GITLEAKS-COVERAGE-MAP.md`
21. `docs/release/LICENSING.md`
22. `docs/tasks/PSCAN-03.md`
23. `evidence/PSCAN-03/MATERIAL-GAP-001.md`
24. `evidence/PSCAN-03/REVIEW-C1.md`
25. `evidence/PSCAN-03/CLOSEOUT-REJECTED.md`
26. `evidence/PSCAN-09/ACTIVATION.md`
27. Owner-provided PASS-SPEC-001 source at the path recorded in
    `docs/spec/PASS-SPEC-001-SOURCE.md`

Then verify the exact activation commit, branch, clean Git status, canonical and
external contract hashes, historical-evidence hashes, PSCAN-09 allowed and
forbidden paths, the absence of another selected/active task and every reserved
owner gate.

Before editing, classify every potentially affected statement as current living
authority or immutable historical evidence. Preserve both PASS-SPEC-001 copies
and every prior task/evidence record. Stop rather than claim or implement if an
authority conflicts, a required source is unavailable, the working tree has
unexplained changes, the activation commit cannot be established as the base,
or the outcome replacement would require an owner decision broader than the
recorded mechanism-to-outcome direction.
