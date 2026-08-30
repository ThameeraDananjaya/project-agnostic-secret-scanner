# PSCAN-02 Acceptance Evidence

## Result

Accepted locally for the exact PSCAN-02 scanner-owned contract, request,
outcome, workspace, CLI-skeleton, and content-free serialization boundary at
`98768cce7e89f55bce2ed269c073985d60618dea`.

This acceptance does not claim an engine, authoritative scan coverage, release,
publication, consuming-project readiness, deployment, or final product
acceptance. It selects or activates no successor.

## Authority and history

- Task: `PSCAN-02` only
- Activation commit: `48a8bd0c30deff691a423fb558a7e45035a5fb24`
- Implementation commit: `b1b1bbcb7cabc3e82f9796a0ed8be0daa51a2021`
- First correction: `e82c78c3cbbdc041f1f8912b18332357f4551885`
- Accepted correction/HEAD:
  `98768cce7e89f55bce2ed269c073985d60618dea`
- Pre-change state: clean and detached at the exact activation commit
- Independent acceptance time: `2026-08-31T00:26:54.2558732+04:00`
- Controlling canonical PASS-SPEC-001: 57,769 bytes, 977 lines, SHA-256
  `78a1ad9a9abfde577a50ae9aa467b5b82f733162c8187ed58a2bc300952b692d`

`REVIEW-01.md` and `REVIEW-02.md` are immutable failed-closed review facts.
`CORRECTION-01.md` and `CORRECTION-02.md` record their bounded corrections. The
final independent review found no actionable acceptance findings.

## Accepted implementation files

1. `cmd/scanner-runner/main.go`
2. `cmd/scanner-runner/main_test.go`
3. `contracts/global-revocation/schema-1.0.json`
4. `contracts/release-manifest/schema-1.0.json`
5. `contracts/rule-pack/schema-1.0.json`
6. `contracts/scan-outcome/schema-1.0.json`
7. `contracts/scan-request/schema-1.0.json`
8. `docs/architecture/ARCHITECTURE.md`
9. `docs/architecture/SOURCE-LAYOUT.md`
10. `docs/decisions/DEC-001-SCHEMA-OWNERSHIP.md`
11. `docs/tasks/PSCAN-02.md`
12. `docs/tasks/TRACKER.md`
13. `docs/validation/VALIDATION-PLAN.md`
14. `evidence/PSCAN-02/CORRECTION-01.md`
15. `evidence/PSCAN-02/CORRECTION-02.md`
16. `evidence/PSCAN-02/IMPLEMENTATION.md`
17. `evidence/PSCAN-02/REVIEW-01.md`
18. `evidence/PSCAN-02/REVIEW-02.md`
19. `fixtures/adversarial/request/duplicate-keys.json`
20. `fixtures/adversarial/request/output-injection.json`
21. `fixtures/adversarial/workspace/cases.json`
22. `fixtures/clean/tracked-manifest.json`
23. `go.mod`
24. `internal/outcome/outcome.go`
25. `internal/outcome/serialize.go`
26. `internal/request/id.go`
27. `internal/request/input.go`
28. `internal/request/load.go`
29. `internal/request/model.go`
30. `internal/request/validate.go`
31. `internal/request/version.go`
32. `internal/workspace/workspace.go`
33. `tests/integration/contract/cli_linux_test.go`
34. `tests/integration/contract/cli_test.go`
35. `tests/integration/contract/cli_windows_test.go`
36. `tests/unit/outcome/outcome_test.go`
37. `tests/unit/request/path_linux_test.go`
38. `tests/unit/request/path_windows_test.go`
39. `tests/unit/request/request_test.go`
40. `tests/unit/workspace/platform_linux_test.go`
41. `tests/unit/workspace/platform_windows_test.go`
42. `tests/unit/workspace/workspace_test.go`

The closeout commit additionally changes only this acceptance record,
`docs/tasks/PSCAN-02.md`, and `docs/tasks/TRACKER.md`.

## Acceptance results

| Criterion | Expected | Actual |
|---|---|---|
| Path/action containment | PSCAN-02 allowlist only | Pass; 42 accepted files, zero outside paths |
| Request versions | corrupt/unknown major fail closed; compatible future minor bounded | Pass |
| Required request structure | missing/null/duplicate/unknown required semantics reject | Pass |
| Request identity/output | exactly one safe content-free object for every representable result | Pass |
| Git bindings | complete lowercase 40- or 64-hex OIDs; no abbreviations | Pass |
| Content bindings | exact lowercase 64-hex SHA-256 | Pass |
| Request time | canonical RFC 3339 UTC ending in `Z` | Pass; offset forms reject |
| Artifact-manifest digest | exact cross-language preimage | Pass; independently reproduced vector |
| State/reason/exit contract | 8 states, 24 reasons, exits 0/10/20/30/40 | Pass |
| Workspace isolation | deterministic attempt path, collision refusal, handle-bound root, owned cleanup | Pass |
| Injection/content boundary | no finding detail, candidate path, or control-text escape | Pass |
| Product independence | no consuming-project authority, identity, data, key, receipt, or statistic | Pass |
| Forbidden implementation scope | no engine/Git-range/extraction/policy/allowlist/workflow/remote code | Pass |
| Dependencies | Go 1.27.0 standard library only | Pass; no external module and no `go.sum` |
| Linux amd64 tests | full suite and vet | Pass at accepted correction |
| Windows targets | compile and platform-sensitive package proof | Compile pass at accepted correction; package suites passed at preceding correction |
| Git integrity | clean reviewed HEAD, valid ancestry, whitespace/connectivity checks | Pass |

## Exact limitations

- Host Application Control blocked fresh Windows test-binary and standalone
  runner execution at the accepted correction. This is not counted as a pass.
- The changed UTC and digest logic is platform-neutral and passed on Linux;
  current Windows targets compile. Earlier Windows in-process request,
  workspace, outcome, and CLI logic passed before that platform-neutral change.
- Windows root-replacement tests skip because the live open-directory handle
  prevents root rename on this host; Linux executes the replacement cases.
- The Linux race suite was not run because the available WSL distribution has
  no C compiler, and no tool was added to bypass that limitation.
- No engine is present by task design. A valid skeleton request therefore ends
  with `UNAVAILABLE_ENGINE` and exit 30.

These limitations do not establish final product or consuming-project
readiness and remain visible to successor planning.

## Tools, pins, and external effects

- Product toolchain pin: Go 1.27.0
- External Go dependencies: none
- Scanner or engine binary/source pin introduced: none
- Gitleaks or TruffleHog material downloaded/executed: none
- Local official Go archives and Graphify output: ignored development material,
  not committed
- Credentials, signing keys, receipts, remote resources, workflow changes,
  publication, spending, or consuming-project integration: none

## Reserved owner gates and successor state

Every remote, publication, signing, credential, spending, provider,
consuming-project, deployment, production, and go-live gate remains reserved.
PSCAN-03 through PSCAN-07 remain proposed and unselected. PSCAN-08 remains
inactive, unselected, and separately technically and legally gated.

## Closeout rule

The closeout commit containing this record must be followed by an exact clean
Git-status check. The final handoff reports that commit and status. No successor
is selected, activated, claimed, or begun by this evidence.
