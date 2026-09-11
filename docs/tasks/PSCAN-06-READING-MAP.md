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
76. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-004-PREFLIGHT.md`
77. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-004-IMPLEMENTATION.md`
78. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-004-AUTHOR-VALIDATION.md`
79. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-004-REFINEMENT.md`
80. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-004-AUTHOR-VALIDATION-002.md`
81. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-005-AUTHORITY.md`
82. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-005-ACCEPTANCE.md`
83. `evidence/PSCAN-06/CORRECTION-C2-BUILD-ONLY-LINUX-DOCKER-PROOF-GATE-AUTHORITY.md`
84. `evidence/PSCAN-06/CORRECTION-C2-BUILD-ONLY-LINUX-DOCKER-PROOF-GATE-FAILURE-001.md`
85. `evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R1-AUTHORITY.md`
86. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHORITY.md`
87. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-PREFLIGHT.md`
88. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-IMPLEMENTATION.md`
89. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION.md`
90. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REVIEW-REJECTED.md`
91. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REFINEMENT-001.md`
92. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION-002.md`
93. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REVIEW-REJECTED-002.md`
94. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REFINEMENT-002.md`
95. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION-003.md`
96. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REVIEW-REJECTED-003.md`
97. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REFINEMENT-003.md`
98. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION-004.md`
99. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-ACCEPTANCE.md`

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

For Correction C2 iteration 005, begin from its exact committed authority
bundle and claim only that iteration in a genuinely fresh session from an exact
clean source-trust-compliant checkout. Treat refined iteration-004 candidate
`9583aa3d18310c2e9275c665f69eb7e5b4fb82a4`, governance HEAD `fad4c3e` and
their author evidence as immutable but independently rejected readiness
history. Prove that a preloaded or stale `PscanNativeBoundary` type is terminal
ambient state and cannot replace the exact committed native implementation,
execute, or fabricate trusted process evidence. Run clean and hostile cases in
separate non-profile processes, preserve every iteration-004 containment and
Docker-operation boundary, and reconcile the tracker to the factual claimed,
implemented/refined, independently rejected state. The iteration-005 approval
does not authorize implementation in the approval session, Docker execution,
network access, a push, C2 tooling tag, settings change, workflow run, signing,
attestation, draft release, publication or successor work.

The bounded iteration-005 local independent result is recorded in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-005-ACCEPTANCE.md`. It accepts only
the local ambient-type isolation and tracker-reconciliation objective at exact
candidate `faef8435322c9096df09b56969662411f17356ea`, tree
`3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`. PSCAN-06 remains open and
unaccepted overall pending actual Linux, genuine Docker/image/container,
dependency acquisition, complete builds, byte comparison and remote proof.
The local result authorizes none of those gates or any successor work.

The owner subsequently approved the exact Correction C2 build-only actual-
Linux and genuine-Docker proof gate for accepted tooling candidate
`faef8435322c9096df09b56969662411f17356ea` from evidence HEAD
`b6d341284cf63baa0502fcba319ea8cca3c7eb3a`. The bounded gate authority is
recorded in
`evidence/PSCAN-06/CORRECTION-C2-BUILD-ONLY-LINUX-DOCKER-PROOF-GATE-AUTHORITY.md`.
This approval session records authority only. A genuinely fresh execution
session must pass the complete current read-only preflight before the exact
fast-forward evidence push, immutable C2 tooling tag and single build-only
workflow dispatch. The C2 signing variable must be absent; signing,
attestation, draft creation, publication, any second run and successor work are
forbidden. The only permitted remote output is one uncompressed unsigned
workflow artifact retained for one day, at maximum cost USD 0.

The first proof-gate execution session failed closed before remote preflight
because its Codex worktree was detached rather than on branch `main`. Preserve
the exact failure record and its SHA-256 as immutable evidence. The owner then
approved Recovery R1 solely to permit a genuinely fresh Daybreak Blue `xhigh`
task to run directly in the saved project at the exact `main`-branch recovery
authority commit. Recovery R1 changes no proof, cost, signing, artifact,
publication or successor boundary. It permits at most one workflow dispatch,
and the failed detached session consumed zero dispatches. Exact scope is in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R1-AUTHORITY.md`.

Recovery R1 then failed closed at mandatory preflight item 1 because two
tracked working-tree files had noncanonical mixed line endings. Read the
append-only failure and subsequent local-only projection repair records in
order:

1. `evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R1-FAILURE-001.md`
2. `evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-SOURCE-TRUST-REPAIR-001.md`
3. `evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-LOCAL-POST-REPAIR-VALIDATION-001.md`

The repair changes no committed source and grants no retry, remote, Docker,
signing, publication or successor authority. The post-repair record proves the
current local source-trust and no-Docker regressions only; it is not remote,
actual-Linux, genuine-Docker or overall PSCAN-06 acceptance evidence.

The owner subsequently authorized one bounded Recovery R2 attempt. Read
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R2-AUTHORITY.md`
before any claim or action. It preserves the original proof-gate identities and
all zero-spend, one-dispatch, unsigned-artifact, no-signing, no-publication and
no-successor limits. The required proof cannot be replaced by local Docker
because GitHub workflow, runner and artifact identities are part of the
acceptance boundary; Recovery R2 therefore permits only the exact committed
GitHub Actions build job. The approval-session GitHub authentication check
failed because the stored token is invalid. Re-authentication must happen
before a genuinely fresh saved-project `main` execution task performs its own
complete current preflight. This approval session records authority only and
does not execute Recovery R2.

Recovery R2 then failed closed in fresh saved-project `main` session
`01a07dbf-b19b-71d1-b40c-617f396b1c64` at the first current GitHub
authentication command. Read the append-only terminal record before any new
PSCAN-06 claim or action:

`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R2-FAILURE-001.md`

The failure record SHA-256 is
`CACDEE7FA590FEA4714A921ECC000AF6F008F4E00515DEDE91C1BEE3E37BAC91`.

The failed attempt made no GitHub repository query or remote mutation,
consumed zero workflow dispatches and performed no Docker, build or artifact
action. Recovery R2 authorizes no retry. A further proof-gate attempt requires
restored authentication and a new exact owner decision. PSCAN-06 remains open
and unaccepted overall.

The owner then authorized one bounded Recovery R3 attempt from the exact
Recovery R2 failure record and evidence commit. Recovery R3 passed exact local
source trust, authority, workflow-byte and current host-level GitHub
authentication checks, but failed closed when all applicable user billing and
usage endpoints were permission-denied and reported that the authenticated
credential lacked the required `user` scope. Read the append-only terminal
record before any new PSCAN-06 claim or action:

`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R3-FAILURE-001.md`

The failure record SHA-256 is
`339E5D66B655B5753D239D5347F1BBAFF0E5858C4CCC289368196D014AE07C12`.

The failed attempt made read-only authenticated GitHub preflight queries but
no remote mutation, consumed zero workflow dispatches and performed no Docker,
build or artifact action. Recovery R3 authorizes no retry or credential-scope
change. A further proof-gate attempt requires trusted current billing,
artifact-storage allowance and stop-usage evidence plus a new exact owner
decision. PSCAN-06 remains open and unaccepted overall.

The owner then authorized one bounded Recovery R4 attempt after current
read-only billing and stop-usage evidence repaired the Recovery R3 blocker.
Read the exact committed authority before any action:

`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R4-AUTHORITY.md`

Recovery R4 passed the complete fresh preflight and fast-forwarded remote
`main` to the exact authority commit. It created the exact local C2 tooling tag
at the accepted candidate, then failed closed before creating the remote tag
because the first tag-push refspec was invalid. Read the append-only terminal
record before any new PSCAN-06 claim or action:

`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R4-FAILURE-001.md`

The failure record SHA-256 is
`5F694FF128515A2B79A2BFA500F299FECAB8836763365E9D774DA5D847D68BEE`.

Recovery R4 consumed zero workflow dispatches and produced no artifact. Remote
`main` is at the R4 authority commit; the exact local C2 tag is preserved; the
remote C2 tag and ruleset remain absent. Recovery R4 authorizes no correction,
retry or second push. A future attempt requires a new exact owner decision for
that split state. PSCAN-06 remains open and unaccepted overall.

The owner then authorized one bounded Recovery R5 attempt from the exact R4
failure and preserved split state. Read its exact committed authority before
interpreting the resulting remote tag or workflow run:

`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R5-AUTHORITY.md`

Recovery R5 passed the complete current preflight, pushed and locked the exact
C2 tooling tag, and dispatched the build-only workflow once. Run `34582887399`
failed in the clean isolated native matrix before genuine Docker, dependency
acquisition, builds or artifact transfer. Read the append-only terminal record
before any new PSCAN-06 claim or action:

`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R5-FAILURE-001.md`

The failure record SHA-256 is
`53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0`.

Recovery R5 consumed its sole tag push and sole workflow dispatch. The exact
remote C2 tag and active no-bypass update/deletion-denial ruleset remain in
place. No artifact, signing, attestation, draft, release, deployment or
publication occurred. Recovery R5 authorizes no rerun, repair, retry, second
push, second dispatch or successor work. A future attempt or correction
requires a new exact owner decision from this terminal state. PSCAN-06 remains
open and unaccepted overall.

The owner then exactly approved Correction C2 iteration 006 from the immutable
Recovery R5 failure record and evidence commit. Read the new bounded authority
before any claim, edit or local check:

`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHORITY.md`

Iteration 006 permits only the local PowerShell PID-collision correction in
the allowed production/test paths and synchronized PSCAN-06 evidence. A fresh
session must prove the lowercase local process-identifier binding no longer
collides with automatic `$PID`, while preserving the parsed identity, ledger
and every existing containment boundary. Local no-Docker checks and
independent review are required.

The iteration authorizes no implementation in the approval session, Docker,
network, remote read or mutation, tag change, workflow dispatch or rerun,
signing, publication or successor work. Recovery R5 remains terminal and
actual-Linux, genuine-Docker, dependency, build, reproducibility and artifact
proof remains open. PSCAN-06 remains open and unaccepted overall.

Iteration 006 was subsequently claimed alone in genuinely fresh saved-project
session `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2` from exact authority commit
`de162e8c347c725a0f041d3c8b0f51df1211c12d`. Bounded author candidate
`8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`, tree
`de70410df5ec359b85d03262e46c8d2ed5f2d822`, passed exact-commit source trust,
the case-insensitive PID-assignment regression and the complete available-host
local no-Docker native matrix. Read these records before review:

1. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-PREFLIGHT.md`
2. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-IMPLEMENTATION.md`
3. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION.md`

This author result is not local acceptance. Fresh independent review rejected
local acceptance because the AST regression misses a valid typed
case-insensitive `PID` assignment target whose left side is
`ConvertExpressionAst`. Read the bounded finding before any new iteration-006
claim, repair or readiness statement:

`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REVIEW-REJECTED.md`

The production rename and available-host no-Docker matrix passed, but the
mandatory regression proof did not. No repair is authorized by the review.
Recovery R5 remains terminal, actual-Linux and every genuine-Docker,
dependency, build, reproducibility, artifact and remote gate remains open, and
no successor is selected or activated.

The same claimed implementation session subsequently refined only the
regression under the existing iteration-006 authority. Read the append-only
refinement and second author-validation records:

1. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REFINEMENT-001.md`
2. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION-002.md`

Refined author candidate `e6a767aab71db1d3f62063dded379b4701d2cb52`,
tree `288664447272f9543de1b69b8ca28c27c6e1e9ff`, preserves the accepted
production blob and adds typed-target-aware, adversarially self-tested AST
coverage. Exact author checks pass. Fresh independent review nevertheless
rejected local acceptance because valid parenthesized and multiple `PID`
assignment targets reproduce the automatic-variable collision while evading
the exact refined predicate. Read the second bounded rejection before any
further iteration-006 claim, repair or readiness statement:

`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REVIEW-REJECTED-002.md`

The earlier candidate and both rejection records remain immutable. No repair
is authorized by the review. Recovery R5 remains terminal, all actual-Linux,
Docker, dependency, build, reproducibility, artifact and remote proof remains
open, and every workflow, tag, signing, publication and successor gate remains
closed.

The same claimed implementation session subsequently performed a second
bounded regression refinement under the existing iteration-006 authority.
Read its append-only refinement and third author-validation records:

1. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REFINEMENT-002.md`
2. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION-003.md`

Second refined author candidate
`763730403ae538850e9d806ce9db82c815517b2d`, tree
`05c0baba5e61edb808b918a2dd6c853f49b8955e`, preserves production blob
`2fb43bf89df352e9153ba5d7ad23fd59cad77749`. Its exact assignment-left walker
covers typed/attributed, parenthesized pipeline/command-expression and nested
multiple-target wrappers without searching member/index descendants. The
three second-review counterexamples and broader positive/control probes are
bound to the exact production predicate. Exact author checks and the complete
available-host no-Docker matrix pass.

Fresh independent review rejected this candidate for local acceptance because
the predicate splits every true variable path on `:` and compares only its
last component. Legal braced names such as `${variable:env:PID}`,
`${global:env:PID}` and `${local:foo:PID}` execute without binding automatic
variable `PID`, but the exact predicate falsely reports each as reserved. Read
the third bounded rejection before any further iteration-006 claim, repair or
readiness statement:

`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REVIEW-REJECTED-003.md`

All three candidates and rejection records remain immutable. No repair is
authorized by the review. Recovery R5 remains terminal; all actual-Linux,
Docker, dependency, build, reproducibility, artifact, remote and successor
gates remain closed.

The same claimed implementation session subsequently performed a third
bounded regression refinement under the existing iteration-006 authority.
Read its append-only refinement and fourth author-validation records:

1. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REFINEMENT-003.md`
2. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION-004.md`

Third refined author candidate
`cf1679f9bca24887340fa4060f37d1ebff21f305`, tree
`6cb0198c71881baff9c601c31ad6aa84217f4a49`, preserves production blob
`2fb43bf89df352e9153ba5d7ad23fd59cad77749`. Public `VariablePath` scope
flags plus exact complete `UserPath` equality now distinguish runtime-proven
automatic PID forms from legal multi-colon variables and provider paths. All
earlier typed, parenthesized, nested, multiple-target and non-target cases
remain covered, and the complete exact-candidate no-Docker matrix passes.

Fourth independent review accepted this exact candidate for bounded
local iteration 006. The exact committed predicate agreed with isolated
runtime behavior across every prior case and independent case/scope, bracing,
escape, colon/separator, wrapper, tuple, provider, member/index and supported-
lvalue mutations. Production data flow, parsers, source trust, confinement,
the full available-host no-Docker matrix and cleanup pass. Read the exact
bounded acceptance before interpreting status:

`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-ACCEPTANCE.md`

Every earlier candidate and rejection remains immutable. This is local
acceptance of iteration 006 only. Recovery R5 remains terminal; all actual-
Linux, Docker, dependency, build, reproducibility, artifact, remote and
successor gates remain closed, and PSCAN-06 remains open and unaccepted
overall.
