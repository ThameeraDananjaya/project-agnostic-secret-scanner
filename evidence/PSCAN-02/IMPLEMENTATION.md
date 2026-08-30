# PSCAN-02 Implementation Evidence

## Authority and claim

- Task: `PSCAN-02` only
- Activation commit: `48a8bd0c30deff691a423fb558a7e45035a5fb24`
- Activation parent: `67b519c23c9aa46e068cb7e1f8b66e245f675a5a`
- Activation subject: `chore(governance): activate PSCAN-02`
- Pre-change state: detached at the exact activation commit; zero Git status
  entries
- Claim: fresh PSCAN-02 implementation session after complete ordered reading
  map and official Go preflight
- Allowed and forbidden scope: exactly `docs/tasks/PSCAN-02.md`
- PASS-SPEC-001 external and canonical sources: byte-identical, 57,769 bytes,
  977 lines, SHA-256
  `78a1ad9a9abfde577a50ae9aa467b5b82f733162c8187ed58a2bc300952b692d`

## Official Go preflight

Verified current facts on 2026-08-30 from official Go primary sources only:

- The official [Go release history](https://go.dev/doc/devel/release) lists Go
  1.27.0 as released on 2026-08-19 and states that a major release is supported
  until two newer major releases exist.
- The official [download JSON](https://go.dev/dl/?mode=json) marks Go 1.27.0
  stable and publishes these exact archives and SHA-256 values:
  - `go1.27.0.windows-amd64.zip`:
    `f0c0a0d33ba94f4d2c5dbc887334ce678b21813504ddb3aafcb06e60a5a667c4`
  - `go1.27.0.linux-amd64.tar.gz`:
    `675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685`
- Both ignored local archives matched those values before extraction; their
  tools report `go version go1.27.0 windows/amd64` and
  `go version go1.27.0 linux/amd64`.
- The official [Go licence](https://go.dev/LICENSE) permits source and binary
  redistribution with its notice/conditions and disclaimer; official
  [source-install documentation](https://go.dev/doc/install/source) describes
  Go as distributed under a BSD-style licence.

Recommendation adopted for PSCAN-02: pin Go 1.27.0 in `go.mod`, use only the Go
standard library, commit no toolchain archive/binary, and defer reproducible
release-toolchain custody to PSCAN-06. This is a technical recommendation, not
publication or release approval.

## Implemented boundary

- Five scanner-owned, independently versioned JSON schema families.
- Eight states, 24 stable reason codes, permitted actions and exact terminal
  exits 0, 10, 20, 30 and 40.
- Bounded local-file/file-descriptor CLI with exactly one JSON stdout object and
  fixed content-free stderr diagnostics.
- Strict duplicate-key/payload/version/required-feature/freshness/mode/digest/
  limit/offline/redaction/fallback/path and bound-file validation.
- Fresh deterministic UUIDv4/attempt workspace with collision refusal,
  platform path checks, ownership marker, identity verification and bounded
  cleanup.
- Content-free typed serializer that has no representable finding-detail field.
- Valid pre-engine requests intentionally stop at `UNAVAILABLE_ENGINE` because
  no engine is allowed in PSCAN-02.

## Validation performed before implementation commit

| Check | Actual result |
|---|---|
| Linux amd64 `go test -count=1 ./...` | pass; CLI, integration, request, workspace and outcome suites |
| Linux amd64 `go vet ./...` | pass; no findings |
| Linux file-descriptor integration | pass |
| Windows request suite | pass, including corrupt/version/binding/path/digest/UUID cases |
| Windows workspace suite | pass, including UNC/root/link/collision/identity/parallel cases |
| Windows outcome suite | pass, including all state/reason/exit and serializer cases |
| Windows in-process CLI suite | pass, including valid, corrupt, unknown-major and injection cases |
| Windows standalone runner and all test binaries | compile pass with Go 1.27.0 |
| Windows standalone CLI process launch | host application control denied launch; not counted as pass |
| Linux standalone CLI process integration | pass, including exact exits and one-object output |
| Linux race suite | not run; WSL distribution has no C compiler and no tool was added |
| External Go modules | none; `go list -m all` reports only this product module |
| `go.sum` | absent because no external module is used |
| Product-code shell or engine process APIs | absent |
| Consuming-project identity/material in new product paths | absent |

The Windows process-launch limitation is isolated to this host policy. Windows
logic executed in Go test binaries and every Windows target compiled, but this
record does not claim a successful standalone Windows process launch.

## Scope exclusions and effects

- Gitleaks/TruffleHog material downloaded or executed: none
- Engine, Git-range, artifact, policy or allowlist implementation: none
- Project-owned schema authority or receipt/signing material: none
- Workflow, remote, credential, signing, publication or spending action: none
- Consuming-project source, identity, policy, data or statistics: none
- Toolchain archives and Graphify output: ignored local tooling only
- Successor selected or activated: none; PSCAN-03 remains proposed/unselected

## Independent review gate

The implementation commit is not accepted by this author record. A separate
review must inspect the bounded commit against PASS-SPEC-001 and PSCAN-02,
rerun proportionate checks, record limitations, and either accept or fail
closed before closeout.
