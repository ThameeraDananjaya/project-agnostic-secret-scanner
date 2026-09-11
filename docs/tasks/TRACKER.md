# PSCAN Task Tracker

| Task | State | Current authority |
|---|---|---|
| PSCAN-01 | Completed and accepted locally | Activation `a85a64e`; accepted implementation `5019489` |
| PSCAN-02 | Completed and independently accepted locally | Activation `48a8bd0`; accepted implementation `98768cc` |
| PSCAN-03 | Rejected and closed fail-closed; not accepted | `evidence/PSCAN-03/CLOSEOUT-REJECTED.md` |
| PSCAN-04 | Completed and independently accepted locally | Exact activation `a21b030e1658f1f98ac4e4d001af12185d9ed311`; evidence-bearing closeout commit |
| PSCAN-05 | Completed and independently accepted locally | Exact activation `50b4186`; corrected evidence-bearing closeout commit |
| PSCAN-06 | Correction C2 iteration 006 candidate `cf1679f9bca24887340fa4060f37d1ebff21f305`, tree `6cb0198c71881baff9c601c31ad6aa84217f4a49`, remains independently accepted locally; Recovery R5 remains terminal; iteration 007 candidate `230e1761f7c45d9629cecf360e7498b33d64ab6f`, tree `5d74628d76f95e150c0ea1de51b4892af2350a7d`, remains independently rejected; iteration 008 stopped with no candidate; iteration 009 candidate `d89b033a6bc99c7e7886fffbbe677d047a479723`, tree `a9ade332849f4345f2e7e8f032504dccc9e60e5a`, remains independently rejected; iteration 010 is owner-approved, unclaimed, not implemented and not accepted as the minimal single-test-file correction for I009-R01; PSCAN-06 remains open and unaccepted overall | Exact iteration 010 authority is under `evidence/PSCAN-06/CORRECTION-C2-ITERATION-010-AUTHORITY.md`; no implementation in the authority session, R6, Docker, network, remote, workflow, tag, release or successor action is authorized |
| PSCAN-07 | Proposed; unselected; signing/publication gates reserved | PASS-OUTCOME-SPEC-001; follows PSCAN-06 |
| PSCAN-08 | Inactive; technically and legally gated | Material-gap evidence plus separate owner approval required |
| PSCAN-09 | Completed and independently accepted locally | Activation `f486989`; accepted closeout commit |
| PSCAN-10 | Completed and independently accepted locally | Activation `9053b37`; accepted correction `d10df99` |

PSCAN-02 completed its activated lifecycle within bounded paths and is accepted
locally. Its exact limitations remain recorded in `evidence/PSCAN-02`.
PSCAN-03 established a material Gitleaks binary-history gap; correction C1 then
failed independent fragment, archive-classification and proof-isolation review.
It is rejected, closed fail-closed and unaccepted. Its closeout selected or
activated no successor.

The owner abandoned the contract-preserving option on 2026-09-01 and activated
PSCAN-09. PSCAN-09 is now independently accepted and closed. It preserved
PASS-SPEC-001 and all PSCAN-01 through PSCAN-03 history, established
PASS-OUTCOME-SPEC-001 through DEC-002, and did not rehabilitate PSCAN-03.

The owner exactly activated PSCAN-10 on 2026-09-01 as the bounded
primary-coverage successor. It was claimed in a fresh implementation session
from exact clean activation commit
`9053b37d799b19e3d98aeca0ae853296971ab09d`, independently accepted after the
correction commit `d10df991d72e3fcc40b378830b50d1f258e16654`, and closed locally. PSCAN-04
was exactly activated by the owner on 2026-09-01;
it was claimed alone on 2026-09-02 in a fresh implementation session from the
exact activation commit. PSCAN-04 is implemented, independently accepted and
closed locally through its evidence-bearing closeout commit. PSCAN-05 was
exactly activated by the owner on 2026-09-02 and claimed alone in a fresh
implementation session from exact activation commit
`50b418609c1f9927c0ecd5d740aba6d7bff11f58`. Candidate
`ab3626ac11e4a62c915e209f1b9f17f098950291` was independently reviewed,
corrected within PSCAN-05 paths, revalidated and accepted through its local
evidence-bearing closeout commit. PSCAN-06 was exactly activated by the owner
on 2026-09-02 and was claimed alone on 2026-09-03 in a fresh implementation
session from
exact clean activation commit
`058ffcd446c6431b2e1afeed769d02c7b1f307f8` after the complete ordered
reading map and current official/read-only preflight. Candidate
`a13c28fe7273bc8dc6545f97966a02889524eb4c` passed two clean network-disabled
builds, byte-for-byte reproduction, complete local tests and skeptical review.
The owner approved the exact remote/signing gate on 2026-09-03. The public
repository and approved controls were created and read back, `main` was pushed,
and locked tag `v1.0.0` was created at the validated candidate. Authorized run
`33709197614` failed closed during pinned dependency acquisition because its
Linux Docker bind-mounted Go module cache was not writable. No build artifact,
signature, attestation, draft release or publication was created. PSCAN-06 is
unaccepted and open. On 2026-09-03 the owner approved the bounded
contract-preserving Correction C1 recovery. Its authority bundle records a
Linux cache-ownership correction and a dual-identity release design that keeps
the public locked `v1.0.0` tag fixed at the accepted product-source candidate
while separately binding correction tooling and workflow identity. The
correction was claimed alone in a fresh session from exact clean authority
commit `3fb1b0a55dc4f48dd35464c63c768f497efbc89b`. Independent review rejected
candidate `f24b832ebe5f6749aa0ab910e1b9be279065ebb1` after reproducing CRLF
carriage-return corruption in its Docker POSIX shell payloads. Bounded
iteration 003 candidate `5ca77226ed3996a8267caf02da366d0beb915c8d` added runtime LF
normalization but was independently rejected when a complete CRLF-checkout
build consumed transformed working-tree integrity bytes and failed the pinned-
rule binding test. Bounded iteration 004 candidate
`a22579fd5af14473e1d49b5027f21ab591bbb589` removed checkout build inputs but
was independently rejected after assume-unchanged driver tampering bypassed
its status-only clean gate. Bounded iteration 005 candidate
`3fb7592889820fa2739a4a53588e073689621809` adds exact committed entrypoints and
complete hostile-state index/path/mode/raw-byte verification. Twelve isolated
adversarial cases reject before output or untrusted driver action, and two
independent complete builds from separate actual CRLF checkouts reproduce all
30 files byte-for-byte. The required actual-Linux host positive/wrong-owner
proof was attempted under the exact build-only gate, but run `33829598255`
failed closed before acquisition because the fresh Ubuntu runner lacked the
pinned Docker image required by the earlier CRLF regression. Both locked tags
remain unchanged; signing and all later steps were skipped. On 2026-09-04 the
owner approved bounded Correction C2 pinned-image bootstrap recovery. Its
authority preserves all historical evidence and existing schemas, permits only
an image-independent host canary and CRLF-normalization proof before the exact
canonical digest-pinned image bootstrap, and proposes a distinct immutable C2
tooling identity. Correction C2 was claimed alone in a fresh session from exact
clean authority commit `d4eca19e04862d660ca6ac0e9b64eec4fb06b61c` and implemented as a bounded
local candidate. Author validation is non-acceptance. Skeptical inspection then
found that generic Docker image-inspection failure was classified as image
absence and could reach the networked pull branch. The owner approved bounded
Correction C2 iteration 002 to make that admission state fail closed. The
iteration was claimed alone in fresh session
`01a06e1c-5c06-7763-958a-fdb623bae6aa` from exact clean authority commit
`281bea031bb6bbaf1be3059977074df2a88ecdf4` and implemented as a bounded local
author candidate. Its deterministic fake-engine validation is author evidence,
not acceptance. Inspection then found that mutable test callbacks remain in the
production admission file and that the nominal Docker output/time limits are
enforced only after unbounded stream capture and may enter an unbounded
post-timeout wait. On 2026-09-05 the owner approved bounded Correction C2
iteration 003 to establish a private production admission boundary and prove
live byte, wall-clock, pipe and process-tree bounds. The iteration was claimed
alone in fresh session `01a06e8c-6a86-74a2-819f-edcbc15fe6f8` from exact
authority commit `5762d1ef6a7708629bad5f5bb33eba21ca1cdcd4` in a separate
checkout with clean tracked, untracked and ignored state. Its bounded author
implementation closes the workflow entrypoint, removes all three mutable
production callbacks, and adds no-Docker native process fixtures for live
stream, UTF-8, timeout, pipe and descendant cleanup behavior. Exact-commit
Windows author validation passed source trust, 12 hostile source cases, the
host-only CRLF proof and the native process matrix. It is not independent
acceptance. Inspection then found that later acquisition, cache, CRLF and build
scripts still invoke ambient `docker`, bypassing the closed admission runner,
and that bounded root `WaitForExit` plus pipe closure does not prove the entire
descendant set terminated. On 2026-09-05 the owner approved bounded Correction
C2 iteration 004 to close every workflow-reachable Docker execution and require
operating-system-backed, empty-membership process containment. Fresh session
`01a07160-ef63-7e80-a441-b2ce360cc23d` claimed only that iteration from exact
authority commit `30d849c11a15a7ce35f182140d07dbe66350c135` and produced a
bounded local implementation candidate
`6f791646413bfde52a7f034f6219d92f6fb44c03`. Exact-commit Windows
no-Docker author fixtures and hostile-source proof passed, but a later audit
found a Linux between-sample detached-member race. Refined candidate
`9583aa3d18310c2e9275c665f69eb7e5b4fb82a4` replaces sampling-only
containment with a stopped, identity-recorded PID-namespace init before Docker
execution. It passed exact-commit Windows author validation. This is not
acceptance. Independent review then rejected its readiness because a preloaded
or stale CLR type named `PscanNativeBoundary` can replace the committed Windows
boundary and fabricate accepted process evidence. The review also identified
and this authority bundle corrects the former tracker-summary contradiction
that described iteration 004 as unclaimed and unimplemented. On 2026-09-05 the
owner approved bounded Correction C2 iteration 005 to isolate the exact
committed native type from ambient process state and add hostile
compatible/stale-type regressions. Fresh session
`01a071ef-5eeb-7a91-8d82-1feb22a9c2f6` claimed only iteration 005 from exact
authority commit `13496a35eab71482f2e908bd0aac4b563941682d` in a new local,
no-network, source-trust-compliant isolated clone. The bounded implementation
rejects any pre-existing native type before compilation or execution, verifies
the exact new compiled type identities, and runs the clean, compatible-hostile
and stale-hostile cases in separate non-profile processes. Local Windows author
fixtures and all iteration-004 local no-Docker regressions pass against exact
candidate `050f5862841a1fcf764ffc75d98a742968ac301e`. This is author validation,
not independent acceptance by itself. Separate reviewer task
`01a07210-d00a-7e21-8816-483c53027c21`, completed turn
`01a07210-d2af-7ca1-b40f-0a9ce6bde7fa`, found no blocking findings against
exact evidence-bearing candidate
`faef8435322c9096df09b56969662411f17356ea`, tree
`3cc6c6234d9cd318792c64ea9e6aa666f146ffb6`, and returned
`LOCAL_ACCEPTANCE_PASS_REMOTE_PROOF_OPEN`. Iteration 005 is independently
accepted locally for its bounded correction objective. PSCAN-06 remains open
and unaccepted overall pending actual Linux, genuine Docker/image/container,
dependency acquisition, complete builds, byte comparison and remote proof. The
review cleanup was blocked before execution by the host guard; no alternate
destructive method was used, and the three retained literal paths are recorded
in `evidence/PSCAN-06/CORRECTION-C2-ITERATION-005-ACCEPTANCE.md`. No correction
Docker, network or remote action was performed, and the local result does not
authorize the remote proof gate by itself. The owner then approved the exact
Correction C2 build-only actual-Linux and genuine-Docker proof gate for
accepted candidate `faef8435322c9096df09b56969662411f17356ea` from evidence
HEAD `b6d341284cf63baa0502fcba319ea8cca3c7eb3a`, with maximum spend USD 0 and
exactly one uncompressed unsigned workflow artifact retained for one day as
the only permitted remote output. The gate is owner-approved but unexecuted.
Its approval session performed no live remote preflight, push, tag, setting,
workflow, Linux/Docker, network or artifact action. A fresh execution session
must satisfy the complete fail-closed preflight before the exact fast-forward
evidence push, immutable C2 tooling tag and one build-only workflow dispatch.
The C2 signing variable must remain absent; signing, attestation, draft
creation, publication, rerun and successor work remain forbidden.
PSCAN-07 remains proposed and
unselected. PSCAN-08 remains inactive,
unselected and ineligible because
accepted evidence did not establish a material required-class Gitleaks gap;
separate technical and AGPL owner approval would also be required. No
successor to PSCAN-06 is selected or activated.

Fresh Daybreak Blue `xhigh` execution session
`01a072e4-50f8-7842-8c7d-28bb034a773f` then attempted only the approved
build-only proof gate from exact authority commit `268e8f3`, but failed closed
at mandatory preflight item 1 because its Codex worktree was detached rather
than on branch `main`. It stopped before remote preflight or any remote,
network, workflow, Docker or artifact action. The immutable failure record has
SHA-256 `7EB1CB03E85A933C1BCFEB435FA85556170AE21F1112CFFA55041C2A733E6A5D`.
The owner approved main-branch execution Recovery R1 on 2026-09-06. Recovery
R1 requires a genuinely fresh Daybreak Blue `xhigh` task running directly in
the saved project on exact branch `main` at the committed recovery authority.
It preserves every original preflight and permits at most one workflow
dispatch because the failed attempt consumed zero. The recovery is
owner-approved but unexecuted. Signing, attestation, draft creation,
publication, spend above USD 0 and successor work remain forbidden.

Recovery R1 was then attempted from the exact required saved-project `main`
checkout and failed closed at mandatory preflight item 1 because two tracked
PowerShell files had noncanonical mixed line endings. It stopped before any
remote request, workflow dispatch, Docker, artifact, signing or publication
action. The failure remains terminal and immutable. Later owner-directed local
maintenance restored only those two working-tree projections from exact
committed Git objects; they now satisfy the canonical LF-to-CRLF projection and
produce no tracked Git diff. This repair grants no retry or remote authority.
Current local post-repair validation passed complete tracked-byte inspection,
all 12 hostile source-trust cases and the Windows no-Docker
admission/process-boundary matrix with zero Docker calls. It does not supply
actual-Linux, genuine-Docker, dependency, reproducibility, workflow or artifact
proof and does not authorize another proof-gate execution.
PSCAN-06 remains open and unaccepted overall, with PSCAN-07 unselected and
PSCAN-08 inactive and ineligible.

On 2026-09-08 the owner authorized one bounded Recovery R2 attempt and allowed
GitHub-hosted Docker when the remaining proof cannot be completed by local
Docker. The remaining gate requires GitHub workflow/ref, actual-Linux runner,
run/job/step and one-day artifact identities, so local Docker cannot substitute
for it. Recovery R2 authorizes only the exact committed GitHub Actions build
job, at most one dispatch and maximum spend USD 0. It does not authorize local
Docker, signing, attestation, drafting, publication, rerun or successor work.

The authority session's requested read-only `gh auth status` check found the
intended active account but an invalid stored token. It made no repository API
query or remote mutation. Recovery R2 is owner-approved, not claimed and not
executed. A genuinely fresh saved-project `main` task may proceed only after
the owner re-authenticates and every current local and remote preflight fact
passes. Exact authority is in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R2-AUTHORITY.md`.
PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains proposed and
unselected; PSCAN-08 remains inactive and ineligible.

Recovery R2 was then attempted in genuinely fresh saved-project `main` session
`01a07dbf-b19b-71d1-b40c-617f396b1c64` using Daybreak Blue at `xhigh`. The
exact local execution context and authority chain passed, but the first current
GitHub authentication command returned exit code `1` because the active
intended account's stored token was invalid. Required execution-context item 6
and mandatory preflight item 4 failed. The task stopped before any GitHub
repository query, remote mutation, workflow dispatch, Docker, build or
artifact action. Recovery R2 is terminally failed, consumed zero workflow
dispatches and authorizes no retry. Exact evidence is in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R2-FAILURE-001.md`.
Its SHA-256 is
`CACDEE7FA590FEA4714A921ECC000AF6F008F4E00515DEDE91C1BEE3E37BAC91`.
A further attempt requires restored authentication and a new exact owner
decision. PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains
proposed and unselected; PSCAN-08 remains inactive and ineligible.

On 2026-09-09 the owner authorized one bounded Recovery R3 attempt from exact
Recovery R2 failure record SHA-256
`CACDEE7FA590FEA4714A921ECC000AF6F008F4E00515DEDE91C1BEE3E37BAC91` at
evidence commit `a4f5b85a3ac63ea31b4bb911be3445f2223b6db1`. The fresh
saved-project `main` task passed exact local source trust, authority,
workflow-byte and host-level GitHub authentication checks. Current user
billing, usage and stop-usage read-back was permission-denied because the
authenticated token lacked the required `user` scope. Mandatory preflight item
4 therefore failed. The task stopped before remote mutation, workflow
dispatch, Docker, build or artifact action. Recovery R3 is terminally failed,
consumed zero workflow dispatches and authorizes no retry. Exact evidence is in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R3-FAILURE-001.md`.
Its SHA-256 is
`339E5D66B655B5753D239D5347F1BBAFF0E5858C4CCC289368196D014AE07C12`.
A further attempt requires trusted current billing, artifact-storage allowance
and stop-usage evidence plus a new exact owner decision. PSCAN-06 remains open
and unaccepted overall; PSCAN-07 remains proposed and unselected; PSCAN-08
remains inactive and ineligible.

On 2026-09-11 the owner exactly authorized one bounded Recovery R4 attempt from
the immutable Recovery R3 failure record SHA-256
`339E5D66B655B5753D239D5347F1BBAFF0E5858C4CCC289368196D014AE07C12` at
evidence commit `6b7ae32f8768c1e4687e5134696688088f65386b`. Read-only
verification confirmed the active GitHub credential includes the required
`user` scope; current Actions Linux and storage usage each have net amount
`USD 0`; and the account-level Actions budget is `USD 0` with stop usage
enabled. Exact authority is recorded in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R4-AUTHORITY.md`.
Recovery R4 is owner-approved, not claimed and not executed. A genuinely fresh
saved-project `main` task must begin from the committed R4 authority bundle and
repeat every current preflight check. This session performs no remote mutation,
workflow dispatch, Docker, artifact, signing, attestation, draft, publication,
rerun or successor work. PSCAN-06 remains open and unaccepted overall;
PSCAN-07 remains proposed and unselected; PSCAN-08 remains inactive and
ineligible.

Recovery R4 was then attempted in genuinely fresh saved-project `main` session
`01a08f7a-7d13-7ee1-be72-a997e4ff9cb6` using Daybreak Blue at `xhigh`. The
complete local, workflow, GitHub identity, tag/ruleset, prior-run, current
billing, storage, budget and stop-usage preflight passed. Remote `main`
fast-forwarded exactly from `3523e4409ebc53cc1931e3dcaf7d1eea74bb15ab`
to R4 authority commit `a82a3a64a04ee2d8b60757866c9864cd1b73b54b`.
The exact unsigned local tag `release-tooling-v1.0.0-c2` was created at
accepted candidate `faef8435322c9096df09b56969662411f17356ea`, but the first
tag-push command failed locally because its shell-expanded refspec was invalid.
The no-retry gate stopped immediately. The command was not corrected or
repeated; no remote C2 tag or ruleset, workflow dispatch, Docker action, build,
artifact, signing, attestation, draft, release or publication occurred.
Recovery R4 is terminally failed and consumed zero dispatches. Exact evidence
is in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R4-FAILURE-001.md`,
SHA-256
`5F694FF128515A2B79A2BFA500F299FECAB8836763365E9D774DA5D847D68BEE`.
The exact local C2 tag remains preserved; remote `main` is at `a82a3a64`, while
the remote C2 tag and ruleset remain absent. Any future attempt requires a new
exact owner decision for that split state. PSCAN-06 remains open and unaccepted
overall; PSCAN-07 remains proposed and unselected; PSCAN-08 remains inactive
and ineligible.

On 2026-09-11 the owner exactly authorized one bounded Recovery R5 attempt from
the immutable Recovery R4 failure record SHA-256
`5F694FF128515A2B79A2BFA500F299FECAB8836763365E9D774DA5D847D68BEE` at
evidence commit `6fbeb633bf373d2f64d95d0ff93c328631a99a05`. The authority preserves
remote `main` at `a82a3a64a04ee2d8b60757866c9864cd1b73b54b`, the existing
local C2 tag at `faef8435322c9096df09b56969662411f17356ea`, and the absent
remote C2 tag, exact-name ruleset, C2 run and artifact.

Recovery R5 permits one genuinely fresh saved-project `main` task using
Daybreak Blue at `xhigh` to repeat the complete fail-closed preflight and, only
if it passes, push the preserved tag exactly once with literal lowercase
refspec
`refs/tags/release-tooling-v1.0.0-c2:refs/tags/release-tooling-v1.0.0-c2`.
It may then create or confirm the exact-name no-bypass tag ruleset and dispatch
the exact build-only workflow at most once. Maximum spend is `USD 0`; at most
one unsigned uncompressed artifact retained for one day is allowed. Tag
recreation, remote `main` mutation, local Docker, signing, attestation, draft,
publication, rerun and successor work are forbidden. Exact authority is in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R5-AUTHORITY.md`.
This authority-recording session performs no execution. PSCAN-06 remains open
and unaccepted overall; PSCAN-07 remains proposed and unselected; PSCAN-08
remains inactive and ineligible.

Recovery R5 was then attempted in genuinely fresh saved-project `main` session
`01a08fae-2ee1-7553-a126-e635b70a34a6`. It passed the complete preflight,
pushed and locked the exact C2 tooling tag, and dispatched workflow run
`34582887399` exactly once. The clean isolated native matrix failed before
Docker, dependency acquisition, builds or artifact transfer because lowercase
local variable `$pid` collided case-insensitively with read-only automatic
variable `$PID`. Recovery R5 consumed its sole tag push and workflow dispatch
and is terminally failed. Exact evidence is in
`evidence/PSCAN-06/CORRECTION-C2-MAIN-BRANCH-EXECUTION-RECOVERY-R5-FAILURE-001.md`,
SHA-256
`53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0`.

On 2026-09-11 the owner exactly approved Correction C2 iteration 006 from that
immutable failure at evidence commit
`da45bb19480663aa21e9a67754ad509834dcaed7`. The iteration is limited to the
local PowerShell PID-collision repair, bounded no-Docker regression evidence
and synchronized PSCAN-06 records. Genuinely fresh saved-project session
`01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2` claimed only iteration 006 from exact
authority commit `de162e8c347c725a0f041d3c8b0f51df1211c12d` and produced
bounded author candidate `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`, tree
`de70410df5ec359b85d03262e46c8d2ed5f2d822`. The exact source regression and
complete available-host no-Docker native matrix pass. This is author evidence,
not independent acceptance. No Docker, network, remote read or mutation, tag
change, workflow dispatch or rerun, signing, publication or successor work was
performed or authorized. PSCAN-06 remains open and unaccepted overall;
PSCAN-07 remains proposed and unselected; PSCAN-08 remains inactive and
ineligible.

Fresh independent review session `01a09003-9a4c-7862-967e-972112c9855d`
then inspected exact candidate `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`
and reran the complete available-host no-Docker matrix. The production rename,
path confinement, source trust, native boundary, stream, timeout, descendant,
cleanup and hostile ambient-type checks passed. Local acceptance was rejected,
however, because the new regression recognizes a reserved `PID` assignment
only when the assignment left side is directly `VariableExpressionAst`.
PowerShell parses valid `[int]$PiD = 1` with a `ConvertExpressionAst` left
side; it reproduces the read-only automatic-variable collision while the
candidate predicate reports zero reserved assignments. Exact evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REVIEW-REJECTED.md`.

The reviewer did not repair the candidate. Recovery R5 remains terminally
failed and consumed; actual-Linux, genuine-Docker, dependency, build,
reproducibility and artifact-integrity proof remains open. PSCAN-06 remains
open and unaccepted overall; PSCAN-07 remains proposed and unselected;
PSCAN-08 remains inactive and ineligible. No successor is selected or
activated.

The same claimed implementation session subsequently refined only the rejected
regression under the existing iteration-006 authority. Refined author
candidate `e6a767aab71db1d3f62063dded379b4701d2cb52`, tree
`288664447272f9543de1b69b8ca28c27c6e1e9ff`, unwraps typed assignment-target
AST wrappers and binds the predicate to untyped, typed and non-target
adversarial self-tests. The production blob remains byte-identical. Exact-
commit source trust, path confinement, parsing, PID source/data-flow regression
and the complete available-host no-Docker matrix pass. Fresh independent review
rejected local acceptance because valid `($PiD) = 1`,
`$PiD, $other = 1, 2` and `($PiD, $other) = 1, 2` targets all reproduce the
read-only automatic-variable collision while the exact refined predicate
reports zero reserved assignments. Exact evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REVIEW-REJECTED-002.md`.

The reviewer did not repair the candidate. Both rejection records and their
candidates remain immutable. Recovery R5 remains terminally failed and
consumed; actual-Linux, genuine-Docker, dependency, build, reproducibility and
artifact-integrity proof remains open. PSCAN-06 remains open and unaccepted
overall; PSCAN-07 remains proposed and unselected; PSCAN-08 remains inactive
and ineligible. No successor is selected or activated.

The same claimed implementation session then produced a second bounded
regression refinement after preserving both rejected candidates and review
records. Author candidate `763730403ae538850e9d806ce9db82c815517b2d`, tree
`05c0baba5e61edb808b918a2dd6c853f49b8955e`, keeps production blob
`2fb43bf89df352e9153ba5d7ad23fd59cad77749` byte-identical. The exact
assignment-left walker now covers typed/attributed, parenthesized and nested
multiple-target wrappers, including the three second-review counterexamples,
while member/index, drive-qualified, RHS, string/comment and ledger-property
cases remain negative controls. Exact source trust, parsing, two-path
confinement, data flow and the complete available-host no-Docker native matrix
pass. Evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REFINEMENT-002.md` and
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION-003.md`.

This is author evidence, not independent acceptance. Recovery R5 remains
terminally failed and consumed; actual-Linux, genuine-Docker, dependency,
build, reproducibility and artifact-integrity proof remains open. PSCAN-06
remains open and unaccepted overall; PSCAN-07 remains proposed and unselected;
PSCAN-08 remains inactive and ineligible. No successor is selected or
activated.

Fresh independent review session `01a09003-9a4c-7862-967e-972112c9855d`
then inspected exact second refined candidate
`763730403ae538850e9d806ce9db82c815517b2d`. Its production blob, exact data
flow, parsers, source trust, path confinement, complete available-host local
no-Docker matrix and residual cleanup checks pass. Local acceptance was
rejected, however, because the exact predicate reduces every true variable
path to the component after its last colon. Legal braced assignments such as
`${variable:env:PID} = 1`, `${global:env:PID} = 1` and
`${local:foo:PID} = 1` execute without colliding with automatic variable
`PID`, while the candidate predicate falsely reports one reserved assignment
for each. Exact evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REVIEW-REJECTED-003.md`.

The reviewer did not repair the candidate. All three rejection records and
their candidates remain immutable. Recovery R5 remains terminally failed and
consumed; actual-Linux, genuine-Docker, dependency, build, reproducibility and
artifact-integrity proof remains open. PSCAN-06 remains open and unaccepted
overall; PSCAN-07 remains proposed and unselected; PSCAN-08 remains inactive
and ineligible. No successor is selected or activated.

The same claimed implementation session then produced a third bounded
regression refinement. Author candidate
`cf1679f9bca24887340fa4060f37d1ebff21f305`, tree
`6cb0198c71881baff9c601c31ad6aa84217f4a49`, preserves production blob
`2fb43bf89df352e9153ba5d7ad23fd59cad77749`. Its exact predicate uses public
`VariablePath` scope flags and complete `UserPath` equality, retaining every
earlier assignment-target wrapper and non-target control while adding parser
and isolated runtime proof for actual automatic PID scopes versus the
multi-colon/provider controls from the third review. Exact source trust,
parsing, two-path confinement, data flow and the complete available-host
no-Docker native matrix pass. Evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-REFINEMENT-003.md` and
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-AUTHOR-VALIDATION-004.md`.

This is author evidence, not independent acceptance. Recovery R5 remains
terminally failed and consumed; actual-Linux, genuine-Docker, dependency,
build, reproducibility and artifact-integrity proof remains open. PSCAN-06
remains open and unaccepted overall; PSCAN-07 remains proposed and unselected;
PSCAN-08 remains inactive and ineligible. No successor is selected or
activated.

Fourth independent review session `01a09003-9a4c-7862-967e-972112c9855d`
then accepted exact third refined candidate
`cf1679f9bca24887340fa4060f37d1ebff21f305`, tree
`6cb0198c71881baff9c601c31ad6aa84217f4a49`, for bounded local iteration 006.
The reviewer extracted the exact committed predicate and independently
verified public `VariablePath` semantics, every earlier positive and control,
mutated case/scope/bracing/escaping/colon/separator forms, wrapper depth,
tuple positions, provider paths, member/index controls and all installed-
parser assignment-left families against isolated runtime behavior. No real
predicate/runtime divergence remained. Production source/data flow, both
parsers, exact source trust, two-path confinement, the complete available-host
no-Docker hostile/native matrix and residual cleanup pass. Exact evidence is
in `evidence/PSCAN-06/CORRECTION-C2-ITERATION-006-ACCEPTANCE.md`.

All three rejected candidates and review records remain immutable. This is
local acceptance of iteration 006 only. Recovery R5 remains terminally failed
and consumed; actual-Linux, genuine-Docker, dependency, build,
reproducibility and artifact-integrity proof remains open. PSCAN-06 remains
open and unaccepted overall; PSCAN-07 remains proposed and unselected;
PSCAN-08 remains inactive and ineligible. No successor is selected or
activated.

The owner then directed repository-consistent autonomous continuation and
authorized recording the minimal prerequisite instead of forcing an invalid
Recovery R6. Read-only Git-object analysis proved that the existing protected
`release-tooling-v1.0.0-c2` tag still points to pre-fix candidate
`faef8435322c9096df09b56969662411f17356ea`, while the committed workflow,
schema 2.1, builder and verifier require that exact tag/ref identity. The
accepted iteration-006 candidate cannot pass those SHA/ref gates, and moving
the protected tag would destroy immutable evidence.

Correction C2 iteration 007 is therefore owner-approved as an unclaimed,
local-only immutable identity roll-forward. It selects proposed tag
`release-tooling-v1.0.0-c2-r6`, new workflow path
`.github/workflows/release-recovery-v1.0.0-c2-r6.yml` and additive manifest
schema 2.2 for bounded implementation and independent local review. Exact
authority is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-007-AUTHORITY.md`. This authority
session creates none of those implementation artifacts and performs no remote,
Docker, dependency, tag, workflow, signing, publication or spending action.
Recovery R6 remains unauthorized until iteration 007 is implemented, validated
and independently accepted and a separate exact R6 authority records fresh
zero-spend remote facts. PSCAN-06 remains the sole open task; PSCAN-07 remains
proposed and unselected; PSCAN-08 remains inactive and ineligible.

Fresh independent review session `01a0916e-95e6-7db0-9391-08232b2fc505`
then inspected exact iteration-007 candidate
`230e1761f7c45d9629cecf360e7498b33d64ab6f`, tree
`5d74628d76f95e150c0ea1de51b4892af2350a7d`. Exact authority, path
confinement, committed workflow/schema identities, builder/verifier agreement,
schema 2.0/2.1 preservation, Docker materialization confinement, the accepted
PID/data-flow bytes, source trust and every available-host no-Docker matrix
passed.

Local acceptance was rejected because the committed workflow-preservation Go
test globally replaces every `2.2` substring with `2.1`. It therefore changes
the unchanged action comment `# v4.2.2` to `# v4.2.1` only in the normalized
new workflow and deterministically fails its equality assertion. Exact evidence
is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-007-REVIEW-REJECTED.md`.

The reviewer did not repair the candidate. Recovery R5 remains terminally
failed and consumed; Recovery R6 remains unauthorized. Actual-Linux,
genuine-Docker, dependency, build, reproducibility and artifact-integrity
proof remains open. PSCAN-06 remains open and unaccepted overall; PSCAN-07
remains proposed and unselected; PSCAN-08 remains inactive and ineligible. No
successor is selected or activated.

On 2026-09-11 the owner selected Correction C2 iteration 008 as the minimal
append-only correction to the sole iteration-007 review finding. Exact
authority is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-008-AUTHORITY.md`.

Iteration 008 may edit only the workflow-normalization portion of
`tests/integration/supply-chain/release_test.go`, add new iteration-008
evidence and synchronize the four named living-state documents. It must
replace the global workflow `2.2` normalization with exact, occurrence-bound
normalization of the single manifest-schema assertion while preserving the
unchanged `# v4.2.2` action comment and every coherent iteration-007 workflow,
schema, builder, Docker, verifier, product and acceptance-test byte.

The owner supplied an exact local Go 1.27.1 executable and frozen module cache
for future offline implementation/review checks. A fresh session must re-prove
that toolchain identity, use fresh external build/temp caches and the exact
offline environment locks recorded in the authority. This authority session
does not execute Go or edit the test. Recovery R5 remains terminal and
consumed; Recovery R6 remains unauthorized. PSCAN-06 remains the sole open
task and unaccepted overall; no successor is selected or activated.

Fresh saved-project session `01a091b5-687a-7392-b366-4b71f8bbbe8d` claimed
only iteration 008 from exact authority commit
`4f86651de67a05db8ae0076c0c155bd184082af3`, tree
`4ef9e344c1849c42509a61077e6a509fae8ef355`. Exact source trust, authority
hashes, critical bytes, token counts and the supplied Go 1.27.1 identity
passed. The required formatter check then returned exit `1` with a 66-line
diff outside the authorized workflow-normalization block. A separate clean
materialization of the untouched authority commit reproduced the same result.

The bounded test attempt was reverted exactly and no candidate survives. The
claiming session stopped before Go tests or vet rather than broaden the
correction. A separate fresh verification continuation then reproduced the
targeted authority failure, the full-package Windows synthetic Cosign failure,
and a passing targeted result under the bounded probe; the probe was restored
and the full package still failed solely at Cosign. Exact evidence is in the
iteration-008 preflight, implementation and author-validation records. Further
implementation requires a new exact owner decision reconciling both baseline
failures with the within-file scope. Recovery R5
remains terminal and consumed; Recovery R6 and all Docker, dependency, remote,
release and successor work remain unauthorized.

On 2026-09-12 the owner selected Correction C2 iteration 009 as the exact
bounded decision for those two baseline failures. It permits only
`tests/integration/supply-chain/release_test.go`, new iteration-009 evidence
and the four named living-state documents. The future candidate must combine
the exact iteration-008 five-token normalization, the demonstrated Go 1.27.1
66-line/four-hunk mechanical formatting projection, and a cross-platform
current-test-executable `TestMain` helper replacing the POSIX-only synthetic
Cosign fixture. Read-only diagnostic task
`01a091b7-90ec-7561-bd79-c8a0f91a522f` proved Windows fails executable
resolution with `exec.ErrNotFound` before `CreateProcess`; no production
verifier defect exists.

Exact authority is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-AUTHORITY.md`. This session
records authority only. Iteration 009 is owner-approved, unclaimed, not
implemented and not accepted. Recovery R5 remains terminal and consumed;
Recovery R6 and all Docker, network, remote, release, spending and successor
work remain unauthorized.

Iteration 009 was then claimed in genuinely fresh implementation session
`01a09228-1204-70b1-8b7e-1be433b2c496`. Its modified
`tests/integration/supply-chain/release_test.go` and untracked
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-PREFLIGHT.md` remain exactly
bound, unstaged and uncommitted. Diagnostic session
`01a09285-6aaa-7033-b65a-845a6521c67d` consumed exactly one clean-
materialization targeted run and one live targeted run; both failed
identically at `tests/unit/artifact/normalize_test.go:507`. Clean and live
copies of that test and `internal/artifact/normalize.go` were byte-identical,
and Windows reported case sensitivity disabled in both roots and the temp
root. Git blame preserves the fixture unchanged from PSCAN-04 commit
`1f0890878518de32a55ceb8d7b97430c4f3d2f2b`.

On 2026-09-12 the owner approved the narrow continuation recorded in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-CONTINUATION-AUTHORITY.md`.
It adds exactly one implementation-bearing path,
`tests/unit/artifact/normalize_test.go`, for a test-fixture-only split of the
case-collision and directory-symlink checks. The case-collision subtest must
use `os.Stat` and `os.SameFile` after creating `A` and `a`, skip only when the
two names are the same physical object, and otherwise preserve the exact
`RejectUnsafe` assertion. `runtime.GOOS` is forbidden; the symlink assertion
and existing availability skip remain. A genuinely fresh continuation session
must start from the committed continuation authority while carrying exactly
the two bound dirty paths, rerun the complete offline iteration 009 matrix and
stop for independent review. This authority session performs no retry, Go,
Docker, network, remote, tag, workflow, signing, publication, spending,
subscription, R6 or successor action.

Fresh continuation session `01a09293-16a9-71f1-8e93-440952d83e2c` then
claimed only iteration 009 from exact continuation-authority commit
`f4c77e239adbb83f5af90741373b22f0bda2c4da`, while carrying the two exact
bound work paths. It split only the inherited combined unit fixture, using
`os.Stat` and `os.SameFile` to skip only an unrepresentable case-collision
subtest while retaining the capable-filesystem and directory-symlink
`RejectUnsafe` assertions. The carried integration patch remained byte-for-
byte unchanged. Exact tool/cache admission, formatting, targeted tests, both
complete packages, repository-wide tests, vet and compile-only Linux/amd64
coverage for all 36 packages passed. This evidence-bearing candidate is author
validated, not accepted; genuinely fresh independent review remains pending.
Recovery R5 remains terminal and consumed, Recovery R6 remains unauthorized,
and no successor is selected or activated.

Fresh independent review session `01a0929e-4232-7ec0-a99a-894a8e90ae74`
then inspected exact candidate `d89b033a6bc99c7e7886fffbbe677d047a479723`,
tree `a9ade332849f4345f2e7e8f032504dccc9e60e5a`. Exact checkout,
materialization, ten-path confinement, formatting, five-token workflow
normalization, `# v4.2.2` preservation, fixture capability checks, complete
packages, repository tests, vet and all 36 Linux/amd64 compile-only packages
passed with the pinned offline toolchain and unchanged frozen cache.

Local acceptance was rejected because the helper admits a contained nested
`--bundle` but derives `arguments.txt` beside `--trusted-root` rather than
beside that bundle as the authority requires. An external-only exact-candidate
probe reproduced the divergent path. Exact evidence is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-REVIEW-REJECTED.md`. The
reviewer did not repair the candidate. Recovery R5 remains terminal and
consumed; Recovery R6 remains unauthorized. PSCAN-06 remains open and
unaccepted overall, and no successor is selected or activated.

On 2026-09-12 the owner selected Correction C2 iteration 010 as the minimal
bounded correction for I009-R01. Exact authority is in
`evidence/PSCAN-06/CORRECTION-C2-ITERATION-010-AUTHORITY.md`. It permits only
`tests/integration/supply-chain/release_test.go`, new iteration-010 evidence
and the four living-state documents. For an otherwise exactly admitted helper
invocation, `arguments.txt` must be derived beside the lexical cleaned
`--bundle` path at argument 2, never beside `--trusted-root`. A positive
contained nested-bundle regression must prove the helper returns and writes
only that exact bundle-sibling path and creates no trusted-root-sibling output;
hostile escapes, symlinks, directories and pre-existing output must continue
to reject.

This authority session records governance only. A genuinely fresh
implementation session must start from the resulting authority commit, use the
pinned offline Go 1.27.1 and frozen module cache with fresh external caches,
produce exactly one bounded candidate commit, and stop for genuinely fresh
independent review. Recovery R5 remains terminal and consumed; Recovery R6 and
all Docker, network, remote, tag, workflow, signing, publication, spending,
subscription and successor actions remain unauthorized. PSCAN-06 remains the
sole open task and remains unaccepted overall. Recovery R6 remains unavailable
until iteration 010 is independently accepted locally and a separate exact R6
authority is later recorded.
