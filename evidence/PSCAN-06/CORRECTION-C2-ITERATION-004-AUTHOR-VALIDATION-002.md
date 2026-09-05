# PSCAN-06 Correction C2 iteration 004 author validation 002

## Superseding author-readiness result

- refined candidate: `9583aa3d18310c2e9275c665f69eb7e5b4fb82a4`
- refined tree: `19233175319ff7220d38119fe296de4b632ce781`
- authority: `30d849c11a15a7ce35f182140d07dbe66350c135`
- prior candidate: `6f791646413bfde52a7f034f6219d92f6fb44c03`
- result:
  `PARTIAL_AUTHOR_VALIDATION_PASS; ACTUAL_LINUX_AND_GENUINE_DOCKER_PROOF_OPEN`
- independent acceptance: not performed and not claimed

The prior author-validation record remains exact historical evidence for
candidate `6f79164`. The containment-refinement record explains why its Linux
sampling design is not the final readiness candidate. This record applies only
to refined candidate `9583aa3`.

## Exact committed checks

The source verifier reported 340 tracked files, all 340 raw-equal, zero CRLF
projections and exact tree `19233175319ff7220d38119fe296de4b632ce781`.
All 12 hostile source cases rejected with zero output and no untrusted driver
action. The forced-CRLF host-only proof passed with actual CR bytes, raw-shell
rejection, LF normalization and `docker=NOT_INVOKED`.

The complete no-Docker admission state matrix and Windows job-object fixtures
passed unchanged:

```text
Docker execution iteration-004 PASS single-boundary=PASS job-assignment-before-resume=PASS kill-on-close=PASS streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT
Image admission iteration-004 PASS state-matrix=PASS parser-only=PASS closed-docker-boundary=PASS
```

Sixteen release PowerShell scripts parsed. Static inventory across the
authority-to-refined-candidate 20-path diff found zero paths outside the
iteration allowlist, zero ambient Docker invocations and zero duplicate native
runner definitions. `git diff --check` passed. The Docker operation table,
payload bytes, Windows containment implementation, schemas, locked tags and
all owner gates are unchanged from candidate `6f79164`.

## Closed Linux design and remaining proof

The refined Linux boundary uses fixed root-owned system `setsid`, `unshare`
and the resolved system shell. The outer session starts the user/PID namespace.
The namespace init stops before Docker, so the parent can bind its PID,
start-time ledger and kernel namespace identity before sending `SIGCONT`. The
init then executes the already held Docker inode. PID-namespace init death is a
kernel boundary for all members, including a process that creates a new session
or nested PID namespace. Trusted return requires supervisor/root exit, both
streams closed, empty outer session and disappearance of the recorded namespace
init identity.

This implementation removes the sampling escape window, but the actual Linux
native fixtures were not runnable on the Windows author host and remain
`UNPROVEN_CURRENT_AUTHOR_HOST`. Genuine Docker engine/image/container
operations, acquisition, the two complete builds, byte comparison and
independent skeptical review also remain open. No Docker, network, tag, push,
workflow, signing, attestation, draft, publication, credential, paid feature,
TruffleHog or successor action occurred. PSCAN-06 remains open and unaccepted.
