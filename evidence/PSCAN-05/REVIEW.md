# PSCAN-05 Independent Review

## Review identity and starting state

- Candidate: `ab3626ac11e4a62c915e209f1b9f17f098950291`.
- Exact activation parent:
  `50b418609c1f9927c0ecd5d740aba6d7bff11f58`.
- Candidate tree: `7eb61ed9d71822a1562a409042ea7066f82cdc70`.
- Branch and initial state: `main`, exact candidate HEAD, clean worktree.
- Authority: `AGENTS.md`, `CONSTITUTION.md`, PASS-OUTCOME-SPEC-001,
  transition traceability, PSCAN-05 task/reading map, DEC-001 and the retained
  predecessor receipt, revocation, compatibility and isolation requirements.

The candidate's implementation, preflight and validation records were treated
as claims. Every one of its 32 changed paths was inspected. All candidate and
review corrections are inside the PSCAN-05 allowlist; protected contracts,
historical evidence, DEC-002, dependencies and successor files are unchanged.

## Corrected blocking findings

Independent adversarial review found and corrected six bounded defects:

1. Duplicate exception IDs could identify different rules/scopes, making
   receipt exception sets and revocation targets ambiguous.
2. The generic chain verifier admitted scanner-owned global records with
   project-only targets, non-digest targets, unregistered reasons or the wrong
   trust domain despite schema 1.1 rejecting them.
3. Receipt reuse did not bind the signed policy and allowlist projection
   identities, so a schema, adapter, trust-domain or key change could be missed
   when exact document bytes were unchanged.
4. Receipt reuse omitted predecessor request/outcome schema digests and the
   exact outcome timestamp.
5. A later receipt sequence could be accepted without proving the immediately
   preceding trusted receipt checkpoint, leaving a gap/divergence path.
6. The acceptance fixture locator embedded the Linux build path, so a valid
   Windows amd64 test executable could not run its adversarial corpus natively.

The corrections add unique exception identity, schema-specific global-record
validation, content-free signed-document reference digests, the missing
schema/time receipt bindings, a trusted previous-receipt anchor, and a safe
working-directory-first fixture lookup. They add no signing, private-key,
storage, promotion, network or project authority.

## Independent validation

The exact retained Go 1.27.1 archives were re-hashed before use:

- Windows amd64:
  `a3911b5e0e1b1053f25ed0675f4c1c6aad1e2bfcf253df2b9be4caabd2edd95d`.
- Linux amd64:
  `63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445`.

Authoritative Linux validation ran as root in a fresh `unshare --net` network
namespace with only a down loopback interface, no routes, DNS resolution exit
2, and public TCP reporting `Network unreachable`. The process environment was
rebuilt with `env -i`; `GOTOOLCHAIN=local`, `GOPROXY=off`, `GOSUMDB=off` and
explicit local caches prevented acquisition or telemetry.

- PSCAN-05 unit, integration, isolation and acceptance packages passed ten
  consecutive runs after the final correction.
- `go vet` passed for every touched package and test package.
- A complete `go test -count=1 ./...` offline regression run passed every
  PSCAN-05 package. Its only failures were the pre-existing WSL-on-Windows
  filesystem/Git boundary: Linux case semantics in one artifact test and
  Windows Git receiving Linux paths in `tests/unit/gitinput`. No failed package
  was changed by PSCAN-05.
- Windows amd64 test executables were cross-compiled in the same networkless
  namespace, then all five PSCAN-05 test binaries ran natively on Windows and
  passed ten consecutive runs. This satisfies native Windows execution without
  treating compile-only success as a pass.
- All five changed JSON documents parsed; `git diff --check` passed; all 32
  candidate paths were allowed; no protected or dependency path changed.
- Static review found no private-key/signing, network client, evidence store,
  project-instance retention, promotion or deployment API in the PSCAN-05
  boundary. Fixtures contain only synthetic identifiers and digests.

## Temporary material

Four exact verified PSCAN-05 temp directories were sent to the Windows Recycle
Bin after validation, so no live temp residue remains and recovery is possible:

- `pscan-05-go1.27.1`: 15,639 files, 320,138,439 bytes.
- `pscan-05-go1.27.1-linux`: 25,574 files, 1,107,743,017 bytes.
- `pscan-05-accept-go1.27.1`: 15,638 files, 241,207,079 bytes.
- `pscan-05-windows-tests`: 5 files, 29,343,232 bytes.

## Final review result

`NO BLOCKING FINDINGS`. PSCAN-05 is recommended for bounded local acceptance.
The result does not claim first-consumer readiness, a release, project
integration, deployment, production, complete PSCAN-07 acceptance or any
owner-reserved action. PSCAN-06 and every successor remain unselected.
