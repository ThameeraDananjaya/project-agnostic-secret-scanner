# PSCAN-06 Correction C2 iteration 004 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Iteration: `004`
- Title: Closed Docker execution and process-tree correction
- Event: exact owner approval of the bounded correction described below
- Owner command:
  `APPROVE PSCAN-06 C2 ITERATION-004 CLOSED DOCKER EXECUTION AND PROCESS-TREE CORRECTION`
- Date: `2026-09-05`
- Approval session: `01a06f21-4e4d-7511-a595-6443412348bd`
- Branch: `main`
- Required starting HEAD before this authority bundle:
  `b96e2696aee5eaa07b113a805e118a6741d52eeb`
- Required starting tree:
  `e52c3d37390d1838c529cbb0877e6a6483ce9fbe`
- Correction authority commit: this authority-bundle commit
- Pre-change tracked and untracked Git status: clean
- Preserved ignored workspace material: `graphify-out/**`; not read as product
  authority, modified, staged or included in this bundle
- State after commit: owner-approved; unclaimed; not implemented; not accepted
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

This is a bounded iteration inside the already-activated PSCAN-06 Correction
C2. It does not create, select or activate another task. A genuinely fresh
session must claim only Correction C2 iteration 004 from this committed
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
- Iteration-003 authority:
  `evidence/PSCAN-06/CORRECTION-C2-ITERATION-003-AUTHORITY.md`
- Iteration-003 authority SHA-256:
  `1CBC031133D17E2B05596E463A10AC2B934B56A329C5F4C1D81222AE67F7013D`
- Unaccepted iteration-003 implementation candidate:
  `dfbe897e9e47632ee5ca9437650bd62eafbaf341`
- Candidate tree:
  `b15586e8ec5b8b799425998920c60348f01e385b`
- Iteration-003 author-validation record:
  `evidence/PSCAN-06/CORRECTION-C2-ITERATION-003-AUTHOR-VALIDATION.md`
- Author-validation SHA-256:
  `544A3E89F8BB1006691E00D8FCF10EC210723E171F0D8B6F5EB77C064F3A5DED`
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

The original C2 and iteration-002/003 authorities, all candidates and evidence,
both locked tags, every existing manifest schema and all predecessor facts
remain immutable. Iteration-003 author validation is not acceptance and is
superseded only as a readiness claim by the missing execution-closure and
descendant-termination proof recorded here.

## Demonstrated gaps

### Docker execution remains open outside admission

Candidate `dfbe897` closes the bootstrap admission entrypoint, but later
release steps still invoke Docker through ordinary PowerShell command
resolution. `acquire.ps1`, `cache-canary.ps1`, `test-crlf-shell-payloads.ps1`
and `build.ps1` contain direct `& docker` calls; `build.ps1` contains two such
calls. Those executions do not necessarily use the fixed executable path,
publisher/file checks, recorded SHA-256, cleared environment, private Docker
configuration, bounded raw streams or terminal classification established by
`admit-image.ps1`.

`image-admission.ps1` also retains a second native Docker runner used by the
later scripts. It duplicates only part of the new boundary and resolves the
executable again. Successful admission therefore does not prove that every
subsequent container, acquisition and build operation is executed by the same
closed command boundary. An alias, function, PATH change, executable
replacement or unbounded later command must not be able to replace or bypass
the admitted Docker identity.

### Root exit is not complete process-tree proof

The iteration-003 runner calls `.Kill($true)` after a terminal condition and
then calls the root process object's bounded `WaitForExit`. That confirms the
root process exit only. Waiting for the two redirected read tasks to finish
proves pipe closure for those handles, but it does not enumerate and prove that
every child and grandchild observed during the command has terminated.

The no-Docker fixture checks recorded descendant identifiers after selected
cases, but the production runner retains no containment object or complete
membership ledger whose empty state is part of the command result. Descendant
creation can race with a tree snapshot or root exit. PID sampling after the
fact is not a durable proof of containment, ownership or complete termination.
Missing full-tree proof fails closed even when the root process and pipes close.

## Approved objective

Establish one closed production Docker-execution boundary for every Docker
operation reachable from the Correction C2 workflow, from image admission
through cache proof, dependency acquisition and both network-disabled builds.
Every call must bind an allowed operation to the exact executable identity,
argument vector, environment, working directory, stream limits, time budget
and process containment. Release scripts must not invoke Docker through an
alias, function, command string, shell, PATH search or caller-provided runner.

Replace root-only cleanup with an operating-system-backed process-containment
boundary established before untrusted child code can create descendants. A
successful or rejected command must not return until bounded stream closure,
root exit and empty containment membership are all proved. If the platform
cannot establish, observe, terminate or prove the complete member set within
the fixed allowance, the command is terminal untrusted evidence and no later
action is permitted.

## Allowed implementation paths

Correction C2 iteration 004 may change only:

```text
.github/workflows/release-recovery-v1.0.0.yml
build/release/acquire.ps1
build/release/admit-image.ps1
build/release/build.ps1
build/release/cache-canary.ps1
build/release/docker-execution.ps1
build/release/image-admission.ps1
build/release/test-cache-boundary.ps1
build/release/test-crlf-shell-payloads.ps1
build/release/test-docker-execution.ps1
build/release/test-image-admission.ps1
docs/release/OFFLINE-VERIFICATION-RUNBOOK.md
docs/release/PSCAN-06-RELEASE-PLAN.md
docs/validation/PSCAN-06-VALIDATION.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/**
```

The two named Docker-execution files may be introduced only if a separate
closed production entrypoint and its no-Docker adversarial harness are needed.
This list supersedes the broader C2 and earlier-iteration path lists for
iteration 004 only. Every absent path is forbidden. A need for another path,
runtime dependency, binary or helper stops the iteration and returns to the
owner with exact evidence.

## Mandatory correction design

1. Inventory every Docker process start reachable from the recovery workflow.
   Route all of them through one committed closed production entrypoint. Remove
   direct `& docker`, `docker.exe`, shell-command, PATH/alias/function and
   duplicated native-runner execution from downstream release scripts.
2. The closed entrypoint must reject dot-sourcing and accept only a finite
   operation name with operation-specific typed inputs. It must construct the
   complete Docker argument list internally. It must not accept arbitrary
   arguments, an executable path, command runner, callback, script block,
   environment-selected implementation or mutable ambient override.
3. Bind the exact Docker executable used for every operation. Preserve the
   fixed Windows Docker Desktop and Linux `/usr/bin/docker` locations and the
   applicable publisher, ownership, mode, link and SHA-256 evidence. Close the
   validation-to-execution replacement window or reject when stable execution
   identity cannot be proved.
4. Use only a canonical private working directory and an empty, private Docker
   configuration directory created for the bounded run. Clear inherited
   environment state and add only the fixed minimum variables. No credential
   helper, plugin, registry login, alternate daemon/context or ambient Docker
   configuration may be consulted.
5. Preserve a finite operation table for engine inspection, exact-image
   inventory, the single approved pull, repository-digest inspection,
   container cache proof, container CRLF parsing, dependency acquisition and
   the two build/package invocations. Each operation must preserve its exact
   digest, mount, user, network, capability, filesystem and resource flags.
6. Establish operating-system-backed containment before the executable can run
   untrusted child code. On Windows, the proof must cover assignment and
   kill-on-close semantics for the complete job or an equally strong boundary.
   On Linux, it must cover a private process group/session or an equally strong
   boundary. A start-before-containment race is not accepted.
7. Maintain a bounded, identity-safe record of every containment member. On
   timeout, overflow, read failure, cancellation or other terminal condition,
   terminate the containment boundary and prove it empty. Root
   `WaitForExit`, pipe closure, PID reuse-prone sampling or best-effort tree
   kill alone is insufficient.
8. On normal exit, prove the root exit, both bounded streams closed, strict
   UTF-8 decoded where structured text is required, and the containment member
   set empty before returning evidence. An unexpected surviving helper or
   detached descendant is terminal even when Docker returns exit code zero.
9. Preserve the independent `131072`-byte stdout/stderr limits, `15000` ms
   complete monotonic command budget and fixed `2000` ms cleanup ceiling.
   Cleanup cannot extend or reclassify command success, and no terminal result
   may be retried or interpreted as image absence.
10. Preserve the iteration-002 three-state admission model, host-cache and
    host-only CRLF prerequisites, the one exact pull rule, post-admission
    `--pull=never --network none`, acquisition ordering and both byte-for-byte
    reproducibility builds.
11. Preserve manifest schema `2.1`, every earlier schema, product/tooling role
    separation, signature identity policy, locked tags and every C2 sandbox,
    credential, remote, publication and owner gate. This iteration changes no
    scanner, verifier, policy, consuming-project or release identity behavior.

## Required local checks and success criteria

1. Static inventory proves there is exactly one production Docker process
   creation boundary and zero ambient `docker` command invocations across every
   workflow-reachable release script.
2. Adversarial probes prove aliases, functions, PATH changes, environment
   variables, Docker configuration, command strings, arbitrary arguments,
   callbacks and dot-sourced state cannot replace the executable, operation or
   evidence path.
3. A replacement-race fixture proves the executed bytes retain the validated
   identity, or the operation rejects before execution. File-path validation
   followed by an independently mutable path lookup is insufficient.
4. The no-Docker state matrix still proves exact present/no-pull, conclusive
   absence/one exact pull and every untrusted/zero-pull outcome with no retry
   or later action.
5. Deterministic native fixtures cover parent exit with a live child, detached
   child and grandchild, child-held stdout/stderr, concurrent output flood,
   timeout during descendant creation, cleanup failure and containment-query
   failure. Every adverse case returns within the fixed total bound and proves
   an empty containment set or reports terminal uncertainty.
6. Immediate exit, nonzero exit, start failure, below/exact/above byte caps,
   simultaneous streams, split valid UTF-8, invalid/incomplete UTF-8 and pipe
   closure cases retain the iteration-003 behavior without truncation or
   unbounded waits.
7. The real closed process-boundary suite runs on Windows and actual Linux. A
   fake-result model, Docker Desktop Linux container or root-only PID check is
   not actual-Linux host process-tree proof.
8. Genuine Docker proof confirms every workflow Docker operation uses the
   closed entrypoint, exact executable identity and exact operation arguments.
   Any pull is limited to the already approved canonical digest and remains a
   separately gated later action; local author tests use no Docker or network.
9. Original C2 and prior-iteration source trust, host-cache, CRLF, admission,
   acquisition, reproducibility, schema/verifier, licensing, SBOM and offline
   verification checks pass without weakening any boundary.
10. Independent skeptical review confirms the complete Docker call inventory,
    stable executable binding and full descendant containment on both native
    platforms before any remote-gate proposal. Author validation cannot accept
    iteration 004 or PSCAN-06.

## Forbidden scope and reserved gates

- No implementation in this approval session.
- No modification of PASS-OUTCOME-SPEC-001, PASS-SPEC-001, their source record,
  TRACEABILITY, DEC-002, existing schemas, verifier/product behavior, scanner
  detection, rules, policy, allowlist, receipt, revocation or consuming-project
  behavior.
- No new runtime dependency, package, downloaded executable or trust root. A
  need for one returns to the owner.
- No movement, deletion, recreation or reuse of either locked tag; no C2
  tooling-tag creation and no replacement product version.
- No Docker command, image pull, image execution, dependency download,
  scanner/toolchain run, build or external network action in this approval
  session.
- No weakened digest pin, mutable/short image reference, registry login,
  retry-on-unknown behavior, privileged container, host networking, added
  capability, `chmod 777` or `--userns=host`.
- No credential, signing key, paid capability, TruffleHog work, push, setting,
  ruleset, workflow run, signing, attestation, draft, publication or other
  remote mutation.
- No deletion or movement of preserved ignored `graphify-out/**` material.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

After bounded local implementation and independent acceptance, a separate
exact owner decision is still required before a push, C2 tooling tag, remote
setting change or build-only workflow run. Signing, attestation, draft creation
and publication remain separately gated.

## Approval-session exclusions

This approval session introduces no workflow, Docker boundary, process runner,
build, test, schema, verifier or product implementation. It performs no Docker
command, image pull, dependency download, scanner/toolchain execution, build,
remote read/write, tag operation, workflow run, signing, attestation, draft or
release action. It creates no credential, key, remote resource or paid
commitment. Its only durable change is this bounded authority bundle and
synchronized living task state.
