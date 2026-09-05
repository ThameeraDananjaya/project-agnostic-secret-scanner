# PSCAN-06 Correction C2 iteration 005 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Iteration: `005`
- Title: Ambient-type isolation and tracker-reconciliation correction
- Event: exact owner approval of the bounded correction described below
- Owner command:
  `APPROVE PSCAN-06 C2 ITERATION-005 AMBIENT-TYPE-ISOLATION AND TRACKER-RECONCILIATION CORRECTION`
- Date: `2026-09-05`
- Approval session: `01a071da-6084-7d82-948c-00a2f2029cfd`
- Branch: `main`
- Required starting HEAD before this authority bundle:
  `fad4c3e12c655944cf8dfc3c46d622666e777fb7`
- Required starting tree:
  `e92c1e729b8aa2d6688dbadd6308c9340958ea9e`
- Correction authority commit: this authority-bundle commit
- Pre-change tracked and untracked Git status: clean
- Preserved ignored workspace material: `graphify-out/**`; not read as product
  authority, modified, staged or included in this bundle
- State after commit: owner-approved; unclaimed; not implemented; not accepted
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

This is a bounded iteration inside the already-activated PSCAN-06 Correction
C2. It does not create, select or activate another task. A genuinely fresh
session must claim only Correction C2 iteration 005 from this committed
authority bundle before implementation. That session must use an exact clean
authority checkout or another source-trust-compliant isolated worktree;
ordinary `git status` does not account for the preserved ignored material in
this checkout.

## Controlling and preserved authority

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Original Correction C2 authority:
  `evidence/PSCAN-06/CORRECTION-C2-AUTHORITY.md`
- Original authority SHA-256:
  `3CE8E20151993A8FC77890E51C5BF1557F68040CAFF2C3DB1A46237792A75D63`
- Iteration-002 authority SHA-256:
  `2FAD7986A1E4BF11D0F469CE2FA0334B6B0E7BE0A0C1372D50483F5846115173`
- Iteration-003 authority SHA-256:
  `1CBC031133D17E2B05596E463A10AC2B934B56A329C5F4C1D81222AE67F7013D`
- Iteration-004 authority:
  `evidence/PSCAN-06/CORRECTION-C2-ITERATION-004-AUTHORITY.md`
- Iteration-004 authority SHA-256:
  `6104EB836906D66C42F1ECDB0DD98BB2E9715F62D1A7E1516D7572B7076B9A75`
- Refined unaccepted iteration-004 candidate:
  `9583aa3d18310c2e9275c665f69eb7e5b4fb82a4`
- Refined candidate tree:
  `19233175319ff7220d38119fe296de4b632ce781`
- Refined author-validation record:
  `evidence/PSCAN-06/CORRECTION-C2-ITERATION-004-AUTHOR-VALIDATION-002.md`
- Refined author-validation SHA-256:
  `881FCA546E6E1D1AA4A94D8FAE213F5AC20AD68473A02C73DBEE674A05E21319`
- Review candidate and governance HEAD:
  `fad4c3e12c655944cf8dfc3c46d622666e777fb7`
- Review candidate tree:
  `e92c1e729b8aa2d6688dbadd6308c9340958ea9e`
- Locked product tag: `v1.0.0`
- Locked product-source commit:
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- Locked Correction C1 tooling tag: `release-tooling-v1.0.0-c1`
- Locked Correction C1 tooling commit:
  `3fb7592889820fa2739a4a53588e073689621809`
- Proposed but absent Correction C2 tooling tag:
  `release-tooling-v1.0.0-c2`
- Pinned build-image identity:
  `docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`

The original C2 and iteration-002/003/004 authorities, every candidate and
evidence record, both locked tags, every existing manifest schema and all
predecessor facts remain immutable. The iteration-004 implementation and
refinement are factual local author work but are independently rejected as
readiness evidence by the two blocking findings recorded below. Neither this
authority nor the tracker correction rewrites that history.

## Demonstrated gaps

### Ambient CLR type can replace the committed Windows boundary

At review candidate `fad4c3e`, `build/release/docker-execution.ps1` compiles its
committed C# boundary only when the current PowerShell process has no type named
`PscanNativeBoundary`. If a caller or previously executed script has already
loaded a compatible type with that name, the guard skips `Add-Type`. Production
then invokes the ambient type and trusts its returned exit, stream, terminal and
containment fields.

A hostile compatible type can therefore fabricate a successful result without
running the exact committed native boundary or the validated Docker executable.
A stale type left in a long-lived process after an earlier script version has
the same substitution effect. The current adversarial harness repeats the same
conditional type-loading pattern, so it can silently exercise the ambient type
instead of proving that the committed implementation ran. This violates the
iteration-004 prohibition on mutable ambient override and is release-blocking
even though the committed Windows fixture suite is otherwise green.

### Living tracker contradicts the implementation history

The PSCAN-06 summary row in `docs/tasks/TRACKER.md` says iteration 004 is
unclaimed and unimplemented. The detailed tracker history, PSCAN-06 task record
and committed evidence instead prove that iteration 004 was claimed in fresh
session `01a07160-ef63-7e80-a441-b2ce360cc23d`, implemented, refined and
author-validated. Independent review then rejected its readiness because of the
ambient-type substitution path. Conflicting lifecycle state is untrustworthy
and fails closed.

The correct state before iteration-005 implementation is: iteration 004 was
claimed and locally implemented/refined, its author checks are non-acceptance,
independent review rejected it, and actual Linux, genuine Docker, complete
build, remote and release proof remain open. PSCAN-06 is open and unaccepted.

## Approved objective

Make the exact committed native boundary the only admissible implementation in
each production process. If the expected native type already exists before the
committed source is compiled, fail closed before any method on that type, any
Docker executable or any later release action can run. Never reuse, inspect as
sufficient, or dispatch through a preloaded or stale type.

Move the native-boundary tests to process-isolated execution so a clean case
compiles and exercises the exact source under test, while a hostile case first
loads a compatible fake type and proves that production rejects before fake or
Docker execution and before fabricated evidence can be accepted. Preserve the
iteration-004 Windows job-object boundary, Linux stopped PID-namespace init,
typed Docker operation table, executable identity, private environment,
stream/time limits and all earlier C2 constraints.

Reconcile the PSCAN-06 tracker summary and detailed history to one accurate
lifecycle state. The documentation correction must preserve every historical
commit and evidence record and must not promote author validation to
independent acceptance.

## Allowed implementation paths

Correction C2 iteration 005 may change only:

```text
build/release/docker-execution.ps1
build/release/test-docker-execution.ps1
docs/release/PSCAN-06-RELEASE-PLAN.md
docs/validation/PSCAN-06-VALIDATION.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/**
```

This list supersedes the broader C2 and earlier-iteration path lists for
iteration 005 only. Every absent path is forbidden. A need to change another
workflow, release script, source file, contract, schema, dependency or helper
stops the iteration and returns to the owner with exact evidence.

## Mandatory correction design

1. Before compiling the committed native-boundary source, test whether the
   exact expected type name already exists in the current process. Any existing
   definition is terminal ambient state. Reject it without invoking, reflecting
   over for compatibility, replacing or reusing it.
2. Compile the single native-boundary source embedded in the exact committed
   production script for a clean process. Compilation failure or an unexpected
   post-compilation type identity is terminal. No fallback implementation or
   previously loaded definition is permitted.
3. Do not accept a caller-provided type name, C# source, assembly, callback,
   script block, runner, executable path or environment-selected implementation.
   Preserve the existing rejection of dot-sourcing and arbitrary arguments.
4. Run clean and hostile native-boundary cases in separate, newly started,
   non-profile PowerShell processes. The harness must not preload and then reuse
   the production type in its own long-lived process.
5. Add a hostile compatible `PscanNativeBoundary` fixture that can return
   fabricated success and records if any fake method is called. Invoke the
   production entrypoint after preload and prove terminal rejection, zero fake
   calls, zero Docker calls and no accepted result.
6. Add a stale-definition case representing an older compatible boundary and
   prove the same rejection. A type that merely has the expected name must
   never become trusted because its methods or result shape look compatible.
7. Preserve the finite nine-operation Docker table, pinned image, executable
   binding, Windows suspended-start/job-object containment, Linux stopped
   PID-namespace-init containment, private directories/environment,
   `131072`-byte per-stream caps, `15000` ms complete command budget and `2000`
   ms cleanup ceiling. Any exact code-context edit must be necessary only for
   the type-isolation gate and must not weaken another property.
8. Reconcile the tracker summary with its detailed record: iteration 004 was
   claimed and locally implemented/refined, its author evidence is not
   acceptance, independent review rejected its readiness, and all actual-Linux,
   genuine-Docker, complete-build, remote and release gates remain open.
9. Preserve schema `2.1`, every earlier schema, product/tooling role separation,
   signature identity policy, locked tags and every C2 sandbox, credential,
   network, remote, publication and owner gate. This iteration changes no
   scanner, verifier, policy, consuming-project or release identity behavior.

## Required local checks and success criteria

1. Static inspection proves the production entrypoint fails closed on a
   pre-existing `PscanNativeBoundary` before `Add-Type`, before any static method
   dispatch and before Docker execution.
2. The hostile compatible-type process returns nonzero within the fixed local
   test bound, records zero fake-boundary calls and zero Docker calls, emits no
   trusted success object and cannot have its fabricated fields accepted.
3. The stale-type process proves the same terminal behavior. A second invocation
   in a process where the type is already loaded must also reject rather than
   silently reuse process state.
4. A clean isolated process compiles the exact embedded source and reruns the
   full Windows native job-object, stream, UTF-8, timeout, pipe, descendant,
   replacement-race and cleanup matrix with zero surviving marked processes.
5. The test harness contains no conditional ambient-type reuse path and proves
   that its clean and hostile subprocesses start with the intended type state.
6. All iteration-004 local no-Docker regressions pass: PowerShell and JSON parse,
   exact source trust, all twelve hostile-source cases, actual-CRLF host-only
   proof with `docker=NOT_INVOKED`, Docker-boundary fixtures, image-admission
   state matrix, static single-entrypoint inventory and diff/path checks.
7. The tracker contains no contradictory current PSCAN-06 state and names the
   new iteration-005 authority as current. PSCAN-07 remains proposed and
   unselected; PSCAN-08 remains inactive and ineligible.
8. Independent skeptical review reproduces the hostile preload attempt and
   confirms the exact production source—not ambient process state—backs every
   accepted local boundary result before any Linux/Docker proof proposal.

Actual Linux host execution, genuine Docker/image/container operations,
dependency acquisition, complete builds, byte comparison and remote proof are
not local success criteria for this bounded implementation and remain open.

## Forbidden scope and reserved gates

- No implementation in this approval session.
- No modification of PASS-OUTCOME-SPEC-001, PASS-SPEC-001, their source record,
  TRACEABILITY, DEC-002, existing schemas, verifier/product behavior, scanner
  detection, rules, policy, allowlist, receipt, revocation or consuming-project
  behavior.
- No changes to the recovery workflow, acquisition, cache, CRLF, image-admission
  or build scripts. A demonstrated need for any such change returns to the
  owner.
- No new runtime dependency, package, downloaded executable or trust root.
- No movement, deletion, recreation or reuse of either locked tag; no C2
  tooling-tag creation and no replacement product version.
- No Docker command, image pull, image execution, dependency download,
  scanner/toolchain run, build or external network action in this approval
  session or during bounded local implementation/author validation.
- No weakened digest pin, mutable/short image reference, registry login,
  retry-on-unknown behavior, privileged container, host networking, added
  capability, `chmod 777` or `--userns=host`.
- No credential, signing key, paid capability, TruffleHog work, push, setting,
  ruleset, workflow run, signing, attestation, draft, publication or other
  remote mutation.
- No deletion or movement of preserved ignored `graphify-out/**` material.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

After bounded local implementation and independent acceptance, a separate
exact owner decision is still required before any actual-Linux or Docker proof,
push, C2 tooling tag, remote setting change or build-only workflow run. Signing,
attestation, draft creation and publication remain separately gated.

## Approval-session exclusions

This approval session introduces no production runner, harness, workflow,
Docker boundary, build, test, schema, verifier or product implementation. It
performs no Docker command, image pull, dependency download, scanner/toolchain
execution, build, remote read/write, tag operation, workflow run, signing,
attestation, draft or release action. It creates no credential, key, remote
resource or paid commitment. Its only durable change is this bounded authority
bundle and synchronized living task state.
