# PSCAN-06 Correction C2 iteration 005 implementation

Date: 2026-09-05  
Implementation session: `01a071ef-5eeb-7a91-8d82-1feb22a9c2f6`  
Authority commit: `13496a35eab71482f2e908bd0aac4b563941682d`

## Bounded change

`build/release/docker-execution.ps1` now treats an already loaded
`PscanNativeBoundary` as terminal ambient state. The check occurs immediately
after the one committed C# source literal and before `Add-Type`, every
`PscanNativeBoundary` method call, private Docker-boundary directory creation
and `Invoke-BoundDocker`. A clean process unconditionally compiles that literal
with `-PassThru`; both expected public types must be unique, come from the same
new assembly and be the exact types resolved by the current runtime. Compile or
identity failure has no fallback.

`build/release/test-docker-execution.ps1` no longer compiles and reuses the
production type in its orchestrator. It starts separate bounded, non-profile
PowerShell processes for:

- the exact-source clean Windows native matrix;
- a compatible fake capable of returning fabricated success;
- an older compatible stale fake;
- dot-source rejection; and
- the second-invocation rejection within the clean process.

Each hostile type records both boundary and would-be Docker dispatch calls. The
hostile processes must return nonzero with the exact production gate error,
zero calls, no emitted result, no fabricated field and no newly created
`pscan-docker-boundary-*` directory. The clean process begins with the type
absent, compiles the exact embedded source, reruns the complete iteration-004
Windows job/stream/UTF-8/timeout/descendant/replacement matrix, proves empty
membership, then proves a second production invocation rejects the now-loaded
exact type.

## Preserved surface

The implementation does not alter the finite nine-operation table, the pinned
image, executable binding, Windows suspended-start/job-object implementation,
Linux stopped PID-namespace implementation, private environment and
directories, `131072`-byte per-stream caps, `15000` ms command budget, `2000`
ms cleanup ceiling, schemas, tags, role separation, workflow, acquisition,
cache, CRLF, build, product, scanner or verifier behavior.

Reachable callers inspected before implementation were `acquire.ps1`,
`admit-image.ps1`, `build.ps1`, `cache-canary.ps1`,
`test-cache-boundary.ps1`, `test-crlf-shell-payloads.ps1` and
`test-image-admission.ps1`. No caller change was required.

## Pre-commit author checks

```text
PowerShell parse: PASS
JSON parse: PASS
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

These checks invoked no Docker, network, download, Go, Gitleaks, Cosign, build,
scanner, remote, signing, attestation, draft or publication action. Exact-commit
author validation and independent skeptical review remain separate. PSCAN-06
remains open and unaccepted; no successor was selected or activated.
