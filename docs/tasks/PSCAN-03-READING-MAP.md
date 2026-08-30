# PSCAN-03 Reading Map

Read these authorities completely, in order, before claiming PSCAN-03 in a
fresh implementation session:

1. `AGENTS.md`
2. `CONSTITUTION.md`
3. `docs/spec/PASS-SPEC-001-SOURCE.md`
4. `docs/spec/PASS-SPEC-001.md`
5. `docs/spec/TRACEABILITY.md`
6. `docs/tasks/TRACKER.md`
7. `docs/tasks/PSCAN-03.md`
8. `docs/governance/TASK-LIFECYCLE.md`
9. `docs/governance/OWNER-GATES.md`
10. `docs/governance/EVIDENCE-AND-CLOSEOUT.md`
11. `docs/governance/RUNBOOK-REQUIREMENTS.md`
12. `docs/decisions/DEC-001-SCHEMA-OWNERSHIP.md`
13. `docs/architecture/ARCHITECTURE.md`
14. `docs/architecture/SOURCE-LAYOUT.md`
15. `SECURITY.md`
16. `docs/security/THREAT-MODEL.md`
17. `docs/validation/VALIDATION-PLAN.md`
18. `docs/release/LICENSING.md`
19. `THIRD_PARTY_NOTICES.md`
20. `evidence/PSCAN-02/ACCEPTANCE.md`
21. `evidence/PSCAN-03/ACTIVATION.md`
22. Owner-provided PASS-SPEC-001 source at the path recorded in
    `docs/spec/PASS-SPEC-001-SOURCE.md`

Then verify the exact activation commit, branch, clean Git status, canonical
contract digest, PSCAN-03 allowed and forbidden paths, the accepted PSCAN-02
limitations, and every still-reserved owner gate. Before introducing any
Gitleaks material or pin, revalidate the current official release, source tag
and commit, licence, security posture, output format, reproducible-build
requirements, and the separation between the CLI/source and GitHub Action.

Stop rather than claim or implement if any authority conflicts, any required
source is unavailable, the working tree contains unexplained changes, the
activation commit cannot be established as the implementation base, or any
source, binary, rule, configuration, output-format, or licence binding cannot be
proved exactly.
