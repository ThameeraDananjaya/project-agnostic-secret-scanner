# PSCAN-06 Correction C2 iteration 004 implementation

## Claim and boundary

- authority commit: `30d849c11a15a7ce35f182140d07dbe66350c135`
- authority tree: `6b13895a9d01fc2f36f55fd8184ec071a1193957`
- implementation session: `01a07160-ef63-7e80-a441-b2ce360cc23d`
- preflight commit: `13fdf6aad9c66ba0a21cc23fb7b52d50a775934f`
- claimed work: PSCAN-06 Correction C2 iteration 004 only
- implementation candidate: the commit containing this record
- independent acceptance: not performed and not claimed

The implementation was made in the isolated clean LF checkout recorded by the
preflight. The saved `main` checkout and its ignored `graphify-out/**` material
were not edited during implementation.

## Closed Docker execution

`build/release/docker-execution.ps1` is now the only production Docker
process-creation boundary. It rejects dot-sourcing and accepts only these
operations:

- engine inspection;
- exact-image inventory;
- the one approved exact-image pull;
- repository-digest inspection;
- container cache proof;
- container CRLF parsing;
- dependency acquisition;
- release build; and
- deterministic release packaging.

The complete Docker argument vector and all POSIX payloads are constructed in
that entrypoint. Downstream release scripts contain no ambient `& docker`,
`docker.exe`, PATH lookup, command string, executable parameter, arbitrary
argument parameter, callback, runner or duplicated native Docker process.
The cache, acquisition and build shell payload bytes are unchanged after LF
normalization from iteration 003.

Admission performs the only conditional pull and writes a bounded receipt only
after engine, inventory, optional pull, fresh digest inspection and the offline
container cache proof succeed. Acquisition and both builds require that exact
receipt, bind the same Docker SHA-256, freshly inspect the exact repository
digest and have no pull operation. The acquisition ledger repeats the bound
Docker digest and closed-boundary identity.

## Executable, environment and containment

Each operation opens the fixed Docker executable with sharing that denies
replacement, validates the held Windows Docker Inc signature or the held Linux
root ownership and non-mutable executable mode, hashes the held bytes, executes
that stable identity and re-hashes the same open handle afterward. Linux uses
the open inode through `/proc/<parent>/fd/<fd>` rather than reopening the
original mutable pathname.

Every invocation creates new empty permission-restricted working, temporary
and Docker-configuration directories. The inherited environment is cleared;
only the fixed platform system/temp/home values and fixed Docker daemon socket
are supplied. No PATH, context, credential helper, plugin, registry login or
ambient Docker configuration is admitted. The configuration directory must
remain empty and the private directory is boundedly removed.

On Windows the executable is created suspended, assigned to a new job with
kill-on-close semantics, recorded as a job member, and only then resumed. Job
membership is queried from the kernel and must be empty after root exit and
both streams close. On Linux the fixed root-owned `setsid` launcher establishes
a private session before the held Docker inode executes; session members are
recorded by PID plus `/proc` start time and known detached members remain in the
identity ledger until dead. Terminal cleanup kills the session and every
still-matching recorded member.

Both paths retain independent 131072-byte stdout and stderr limits, strict
UTF-8 after complete stream closure, a 15000 ms monotonic command budget and a
2000 ms terminal cleanup ceiling. Timeout, overflow, read, invalid UTF-8,
termination, pipe, membership or cleanup uncertainty is terminal and cannot be
retried or interpreted as image absence.

## Docker-free author checks before candidate commit

The following checks used only signed/system tooling and disposable fixtures;
they did not invoke Docker, a scanner, Go, a container, a network request or a
remote action:

- every release PowerShell file parsed;
- the embedded C# native boundary compiled;
- static inventory found zero ambient Docker invocations and no duplicated
  admission runner;
- the moved cache, acquisition and build POSIX payloads matched the authority
  bytes exactly after LF normalization;
- the complete present/absent/untrusted image-admission model passed; and
- Windows signed-PowerShell fixtures passed immediate/nonzero/start failure,
  below/exact/above stream caps, simultaneous streams, split/invalid/incomplete
  UTF-8, hang, child, grandchild, empty-job and held-file replacement cases.

The exact Windows fixture summary was:

```text
Docker execution iteration-004 PASS single-boundary=PASS job-assignment-before-resume=PASS kill-on-close=PASS streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT
Image admission iteration-004 PASS state-matrix=PASS parser-only=PASS closed-docker-boundary=PASS
```

Definitive exact-commit author validation is recorded separately after this
implementation commit exists. Actual Linux execution, genuine Docker operation
proof, dependency acquisition, complete builds and independent skeptical
review remain open. No Docker, network, tag, push, workflow, signing,
attestation, draft, publication, credential, paid capability, TruffleHog or
successor action occurred.
