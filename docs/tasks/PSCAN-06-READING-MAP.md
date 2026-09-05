# PSCAN-06 Reading Map

Read these authorities completely, in order, before claiming PSCAN-06 in a
fresh implementation session:

1. `AGENTS.md`
2. `CONSTITUTION.md`
3. `docs/spec/PASS-OUTCOME-SPEC-001.md`
4. `docs/spec/TRACEABILITY.md`
5. `docs/tasks/TRACKER.md`
6. `docs/tasks/PSCAN-06.md`
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
24. `docs/validation/PSCAN-05-VALIDATION.md`
25. `docs/release/RELEASE-AUTHORITY.md`
26. `docs/release/LICENSING.md`
27. `LICENSE`
28. `THIRD_PARTY_NOTICES.md`
29. `go.mod`
30. `go.sum`
31. `.github/CODEOWNERS`
32. `build/gitleaks/manifest.json`
33. `licenses/gitleaks/LICENSE-v8.30.1.txt`
34. `licenses/gitleaks/modules/manifest.json`
35. `contracts/release-manifest/schema-1.0.json`
36. `contracts/global-revocation/schema-1.0.json`
37. `contracts/global-revocation/schema-1.1.json`
38. `docs/tasks/PSCAN-10.md`
39. `evidence/PSCAN-10/PREFLIGHT.md`
40. `evidence/PSCAN-10/IMPLEMENTATION.md`
41. `evidence/PSCAN-10/VALIDATION.md`
42. `evidence/PSCAN-10/REVIEW.md`
43. `evidence/PSCAN-10/ACCEPTANCE.md`
44. `evidence/PSCAN-10/CLOSEOUT.md`
45. `docs/tasks/PSCAN-04.md`
46. `evidence/PSCAN-04/PREFLIGHT.md`
47. `evidence/PSCAN-04/IMPLEMENTATION.md`
48. `evidence/PSCAN-04/VALIDATION.md`
49. `evidence/PSCAN-04/REVIEW.md`
50. `evidence/PSCAN-04/ACCEPTANCE.md`
51. `evidence/PSCAN-04/CLOSEOUT.md`
52. `docs/tasks/PSCAN-05.md`
53. `evidence/PSCAN-05/ACTIVATION.md`
54. `evidence/PSCAN-05/PREFLIGHT.md`
55. `evidence/PSCAN-05/IMPLEMENTATION.md`
56. `evidence/PSCAN-05/VALIDATION.md`
57. `evidence/PSCAN-05/REVIEW.md`
58. `evidence/PSCAN-05/ACCEPTANCE.md`
59. `evidence/PSCAN-05/CLOSEOUT.md`
60. `evidence/PSCAN-06/ACTIVATION.md`
61. `evidence/PSCAN-06/REMOTE-GATE-FAILURE-001.md`
62. `evidence/PSCAN-06/CORRECTION-C1-AUTHORITY.md`
63. `evidence/PSCAN-06/REMOTE-GATE-FAILURE-002.md`
64. `evidence/PSCAN-06/CORRECTION-C2-AUTHORITY.md`
65. `evidence/PSCAN-06/CORRECTION-C2-IMPLEMENTATION.md`
66. `evidence/PSCAN-06/CORRECTION-C2-AUTHOR-VALIDATION.md`
67. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-002-AUTHORITY.md`
68. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-002-PREFLIGHT.md`
69. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-002-IMPLEMENTATION.md`
70. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-002-AUTHOR-VALIDATION.md`
71. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-003-AUTHORITY.md`
72. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-003-PREFLIGHT.md`
73. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-003-IMPLEMENTATION.md`
74. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-003-AUTHOR-VALIDATION.md`
75. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-004-AUTHORITY.md`

Then verify the exact activation commit, branch and clean Git status; the
activated allowed and forbidden paths; the absence of another selected,
activated or claimed task; and every reserved owner gate. Inspect every tracked
file under the allowed implementation paths before editing. Treat workflows,
build tooling, source, binaries, rules, SBOMs, licences, notices, manifests,
signatures, bundles and attestations as untrusted until independently bound and
verified. Never introduce consuming-project authority or data.

Before any download, dependency or action change, build, scanner/tool execution,
pinning or signing experiment, perform and record a fresh read-only
primary-source preflight for the exact Go, Gitleaks, Cosign/Sigstore, GitHub
Actions, GitHub Releases, immutable releases, artifact attestations, OIDC,
rules/settings, plan/cost, SBOM and licensing facts. Bind every admitted tool,
action and asset to an exact immutable identity and independently verified
licence/notice evidence. Do not treat prior PSCAN-10 intake as current release
or distribution approval.

Before any remote creation, push, settings change, workflow enablement or
signing dry run, record the current read-only account, repository availability,
visibility, workflow identity, permissions, protection/ruleset, capability and
zero-cost facts; present the exact proposed action; and stop for separate owner
approval. Activation alone authorizes no remote or signing action.

Stop rather than claim or continue if an authority conflicts, the working tree
contains unexplained changes, the activation commit cannot be established as
the implementation base, a required primary source or exact pin is unavailable,
reproducibility or native Windows/Linux verification cannot be proved, any
licence/notice/SBOM/manifest/signature/identity/revocation binding is incomplete,
network or credential isolation cannot be proved, zero-spend cannot be
maintained, a remote gate is reached without exact approval, or any missing,
stale, conflicting, unsupported or untrustworthy condition would otherwise be
treated as pass.

For Correction C1, begin from the exact committed correction-authority bundle,
not from the original candidate or failed workflow run alone. Re-read the
locked product tag locally and remotely without mutation. Treat the product
source identity and correction-tooling/workflow identity as separate mandatory
roles, and verify that every implementation path is present in the narrower
Correction C1 path list in `docs/tasks/PSCAN-06.md`. The Correction C1 approval
does not authorize a push, tooling tag, settings change, workflow run, signing,
attestation, draft release or publication.

For Correction C2, begin from its exact committed authority bundle and treat
run `33829598255` as immutable failed evidence. Preserve both locked tags and
all existing manifest schemas. Admit only the exact canonical digest-pinned
build image through the separately bounded bootstrap. Before any pull, prove
host-cache semantics and CRLF rejection/normalization without Docker; after the
pull, prove the exact repository digest and run the container canary and shell
parser with `--pull=never --network none`. Preserve the later bounded dependency
acquisition and network-disabled builds. The C2 approval does not authorize a
push, C2 tooling tag, settings change, workflow run, signing, attestation,
draft release or publication.

For Correction C2 iteration 002, begin from its exact committed authority
bundle and claim only that iteration in a fresh session. Treat candidate
`e7faf0f81b3853e2090c75378bfbca568b52efad` and its partial author validation
as immutable but superseded readiness evidence. Prove a three-way pre-pull
state: exact present, conclusively absent, or untrusted/failed. Only conclusive
absence on an otherwise responsive engine may reach the single exact pull.
Every generic inspect failure, daemon or permission error, timeout, malformed
or ambiguous result must stop with zero pulls. The iteration-002 approval does
not authorize Docker execution, network access, implementation in the approval
session, a push, C2 tooling tag, settings change, workflow run, signing,
attestation, draft release, publication or successor work.

For Correction C2 iteration 003, begin from its exact committed authority
bundle and claim only that iteration in a fresh session from an exact clean
source-trust-compliant checkout. Treat candidate
`52f7ee22ab722d7590e2d2c8326e14a0e9670462` and its author validation as
immutable but superseded readiness evidence. Prove that the workflow-facing
production admission path cannot consult mutable callbacks or caller-provided
test doubles, and enforce stdout/stderr byte caps, wall-clock expiry, concurrent
pipe draining and complete process-tree cleanup while the native Docker process
is running. A fake-result case is not native process proof. The iteration-003
approval does not authorize implementation in the approval session, Docker
execution, network access, a push, C2 tooling tag, settings change, workflow
run, signing, attestation, draft release, publication or successor work.

For Correction C2 iteration 004, begin from its exact committed authority
bundle and claim only that iteration in a genuinely fresh session from an exact
clean source-trust-compliant checkout. Treat candidate
`dfbe897e9e47632ee5ca9437650bd62eafbaf341` and its author validation as
immutable but superseded readiness evidence. Prove that every Docker start
reachable from admission, cache proof, CRLF proof, acquisition and both builds
uses one exact closed operation boundary; ambient `docker`, duplicated runners
and arbitrary argument execution are forbidden. Prove operating-system-backed
containment established before child execution, bounded full-member
termination and empty membership on success and failure. Root exit, pipe
closure or post-hoc PID sampling alone is insufficient. The iteration-004
approval does not authorize implementation in the approval session, Docker
execution, network access, a push, C2 tooling tag, settings change, workflow
run, signing, attestation, draft release, publication or successor work.
