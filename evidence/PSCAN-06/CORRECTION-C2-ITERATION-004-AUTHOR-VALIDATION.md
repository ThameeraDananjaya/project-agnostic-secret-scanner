# PSCAN-06 Correction C2 iteration 004 author validation

## Result and exact identity

- validation UTC: `2026-09-05T12:34:47Z`
- implementation candidate:
  `6f791646413bfde52a7f034f6219d92f6fb44c03`
- candidate tree: `6fa76a12b308b7b09400c3fa587cc284deb6873e`
- authority commit:
  `30d849c11a15a7ce35f182140d07dbe66350c135`
- result:
  `PARTIAL_AUTHOR_VALIDATION_PASS; ACTUAL_LINUX_AND_GENUINE_DOCKER_PROOF_OPEN`
- independent acceptance: not performed and not claimed
- Docker/network/remote/signing/release action: not performed

This validation applies only to the exact committed implementation candidate.
It does not accept Correction C2, PSCAN-06, a tooling tag, a remote workflow or
a release.

## Exact source and hostile-state proof

The committed source verifier reported:

```text
Exact source trust PASS commit=6f791646413bfde52a7f034f6219d92f6fb44c03 tree=6fa76a12b308b7b09400c3fa587cc284deb6873e files=338 raw_equal=338 canonical_crlf=0
SOURCE-TRUST PASS adversarial_rejections=12 exact_positive=PASS zero_build_outputs=PASS untrusted_driver_action=ABSENT
```

Assume-unchanged, skip-worktree, staged, unstaged, ordinary untracked, ignored
untracked, filesystem-monitor, untracked-cache, driver-tampering, path,
mode and blob cases each rejected with exit 1, zero output files and no
untrusted driver action.

An exact local forced-CRLF clone reported:

```text
CRLF host-only regression PASS checkout-asset-CR=PROVED raw-shell-CR=REJECT normalized-shell-CR=ABSENT docker=NOT_INVOKED
```

The moved cache, acquisition and build POSIX payloads were independently
extracted from the authority commit and candidate. Their normalized text,
length and SHA-256 values are identical:

| Payload | Length | SHA-256 |
|---|---:|---|
| cache canary | 445 | `fe1bd0ba2a19c726b1c77b5498d6a1aa73ebb298583fe0e299b45d04439c6312` |
| dependency acquisition | 1537 | `c16f3e630d2dde3d917e874f720e11e2d7b67d49c9a390572c5139fb6f3a2666` |
| release build | 6310 | `2b19f7fb2fb0084f266e47025206907fc8d8fdfc3a123d2ca3ab0755615f4632` |

## Closed-boundary and native Windows proof

Static inventory found no ambient `& docker` or `docker.exe` execution in a
workflow-reachable release script and no duplicate native admission runner.
All nine Docker operations are in the single typed operation table in
`docker-execution.ps1`. The pure image-admission helper contains parsing and
classification only.

The embedded native C# boundary compiled from the committed script. Signed
Windows PowerShell fixtures exercised the exact job-object runner with
immediate success, nonzero exit, start failure, stdout/stderr below and at the
131072-byte limit, overflow by one byte, simultaneous full streams, split
UTF-8, invalid and incomplete UTF-8, hang, child and grandchild cases. The root
was created suspended, assigned to a kill-on-close job before resume, and
kernel job membership was empty before every result returned. No fixture PID
remained live. Holding the validated executable handle rejected replacement.

Exact output:

```text
Docker execution iteration-004 PASS single-boundary=PASS job-assignment-before-resume=PASS kill-on-close=PASS streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT
Image admission iteration-004 PASS state-matrix=PASS parser-only=PASS closed-docker-boundary=PASS
```

The structurally separate no-Docker admission model passed exact-present with
zero pulls, conclusive absence with exactly one pull, and daemon, permission,
timeout, overflow, malformed, ambiguous, pull-failure, null, scalar, empty,
duplicate, alias, mixed, wrong-repository and wrong-digest rejections with no
later action.

## Scope and preservation

- 16 release PowerShell files parsed.
- 19 tracked JSON files parsed.
- The authority-to-candidate diff contains 18 paths and zero paths outside the
  iteration-004 allowlist.
- `git diff --check` passed.
- The controlling contract, original C2 and iteration-002/003 authority,
  iteration-003 author validation and manifest schemas 1.0, 1.1, 2.0 and 2.1
  retain all nine recorded SHA-256 values.
- `v1.0.0` remains commit `a13c28fe7273bc8dc6545f97966a02889524eb4c`
  and tree `217b711ddea51fd0ea7e808edd2e27fdecef8427`.
- `release-tooling-v1.0.0-c1` remains commit
  `3fb7592889820fa2739a4a53588e073689621809` and tree
  `7de9dc4c5725bf38cc80aa734e3c2bac0abd9762`.
- The proposed C2 tooling tag remains absent.

## Honest remaining boundary

This Windows author run did not execute the Linux private-session and detached-
member fixtures. It also did not invoke Docker, inspect an engine or image,
pull, run a container, acquire dependencies, execute either complete build,
compare distributables, run Go/scanner tests, or perform independent skeptical
review. Those checks are required before acceptance and may run only inside
their separately authorized local/remote boundaries. The Linux implementation
is present but remains `UNPROVEN_CURRENT_AUTHOR_HOST`.

No tag, push, setting, workflow, signing, attestation, draft, publication,
credential, paid capability, TruffleHog or successor action occurred.
PSCAN-06 remains open and unaccepted.
