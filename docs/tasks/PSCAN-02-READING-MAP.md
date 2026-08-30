# PSCAN-02 Reading Map

Read these authorities completely, in order, before claiming PSCAN-02 in a
fresh implementation session:

1. `AGENTS.md`
2. `CONSTITUTION.md`
3. `docs/spec/PASS-SPEC-001-SOURCE.md`
4. `docs/spec/PASS-SPEC-001.md`
5. `docs/spec/TRACEABILITY.md`
6. `docs/tasks/TRACKER.md`
7. `docs/tasks/PSCAN-02.md`
8. `docs/governance/TASK-LIFECYCLE.md`
9. `docs/governance/OWNER-GATES.md`
10. `docs/governance/EVIDENCE-AND-CLOSEOUT.md`
11. `docs/governance/RUNBOOK-REQUIREMENTS.md`
12. `docs/decisions/DEC-001-SCHEMA-OWNERSHIP.md`
13. `docs/architecture/ARCHITECTURE.md`
14. `docs/architecture/SOURCE-LAYOUT.md`
15. `docs/security/THREAT-MODEL.md`
16. `docs/validation/VALIDATION-PLAN.md`
17. `evidence/PSCAN-01/ACCEPTANCE.md`
18. `evidence/PSCAN-02/ACTIVATION.md`
19. Owner-provided PASS-SPEC-001 source at the path recorded in
    `docs/spec/PASS-SPEC-001-SOURCE.md`

Then verify the exact activation commit, branch, clean Git status, canonical
contract digest, PSCAN-02 allowed and forbidden paths, current official Go
toolchain support and licence facts before pinning, and every still-reserved
owner gate.

Stop rather than claim or implement if any authority conflicts, any required
source is unavailable, the working tree contains unexplained changes, or the
activation commit cannot be established as the implementation base.
