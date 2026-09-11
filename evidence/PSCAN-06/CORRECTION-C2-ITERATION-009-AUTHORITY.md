# PSCAN-06 Correction C2 iteration 009 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Iteration: `009`
- Title: cross-platform current-test-executable Cosign fixture and exact
  preservation normalization
- Event: owner selection of the bounded one-test-file correction after
  iteration 008 stopped with no candidate
- Decision date: `2026-09-12`
- Approval session: `01a091f1-3d3f-7a41-a541-f705cd29b687`
- Branch: `main`
- Required starting HEAD before this authority bundle:
  `c6ea1f08b28bb654023546f30156ff71bcfe4148`
- Required starting tree:
  `0a4f81c1f4a46f0f05d3f87f2a942ea3ce60f007`
- Required starting HEAD direct parent:
  `4f86651de67a05db8ae0076c0c155bd184082af3`
- Correction authority commit: this authority-bundle commit
- Pre-change tracked, staged and non-ignored untracked Git status: clean
- Pre-change graphify-aware exact checkout verification: 384 tracked paths,
  366 raw-equal, 18 canonical whole-file LF-to-CRLF projections, exact index
  and flags, zero unsupported Git configuration, zero non-ignored untracked
  paths, and 66 preserved ignored paths all under `graphify-out/**`
- Preserved ignored workspace material: `graphify-out/**`; excluded from
  product authority and not refreshed, deleted, moved or staged
- State after commit: owner-approved; unclaimed; not implemented; not accepted
- Iteration 008 state: claimed, stopped fail-closed, no candidate
- Recovery R5 state: terminally failed and consumed
- Recovery R6 state: not authorized and not executable
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

This is a bounded iteration inside already activated PSCAN-06 Correction C2.
It selects no successor. This session records and commits authority only. It
does not edit the test or execute an implementation. A genuinely fresh
session with an ID different from this approval session must claim only
Correction C2 iteration 009 from this committed bundle before implementation.

## Controlling and immutable inputs

- `PASS-OUTCOME-SPEC-001` committed-byte SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Iteration-007 authority SHA-256:
  `A17520815841415AD06357682423CF7C635B32AAFAB7EA2B0303E2540FE709E8`
- Iteration-007 implementation SHA-256:
  `F251946635F474DD4AF39071EBBE14CED0F48FA52AD5FF5AF4D130822D59083B`
- Iteration-007 independent rejection SHA-256:
  `E8B90DDE62848577E129B7E2857FE6D0BF56BBF691440A06636013BE458FAC21`
- Iteration-008 authority SHA-256:
  `1CA9CC9B8DDD6E76498149F266F0C61555E798ADC8963F18109EF981131AAE21`
- Iteration-008 preflight SHA-256:
  `A87B5B90862259FB752ED8AD397FDB5AD759FE47AA701D40C167B968C87AEA97`
- Iteration-008 implementation SHA-256:
  `F955119183BC53252B3537E33D028D38E39834A41B5FC1F8AF99C9C928637D62`
- Iteration-008 author-validation SHA-256:
  `8E2B8EEDE526A956D756831E0902B21BF47C864CEF32C8054A64A65CC362CE58`
- Rejected iteration-007 candidate/tree:
  `230e1761f7c45d9629cecf360e7498b33d64ab6f` /
  `5d74628d76f95e150c0ea1de51b4892af2350a7d`
- Recovery R5 terminal-failure SHA-256:
  `53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0`

Iteration 007 and its rejection remain immutable. Iteration 008 produced no
candidate: its exact normalization probe passed, but the untouched authority
commit independently reproduced both the pre-existing Go 1.27.1 formatter
diff and Windows synthetic-Cosign failure. This authority does not reinterpret
either prior iteration as accepted.

## Exact current test-file and diagnostic evidence

The only implementation-bearing path selected by this correction is
`tests/integration/supply-chain/release_test.go`. At the required starting HEAD
it has Git blob `60f8326cdde115a03e3219b04a3e1e40fbb43aa3`, committed-byte
SHA-256 `29E715D5534B389BCB65E046E84E5F2B2828F3D400BB664AFB13B55F12F377D0`,
17,768 bytes, 351 LF line endings, zero CRLF pairs and zero bare CR bytes.

Exact Go 1.27.1 `gofmt -d` produces a real pre-existing 66-line diff with
exactly four hunks. Iteration 009 admits only this demonstrated mechanical
projection across the complete file:

| Hunk | Permitted mechanical effect |
|---|---|
| `@@ -162,7 +162,9 @@` | expand the one-line `role swap` function literal |
| `@@ -193,10 +195,14 @@` | align the old-under-2.2 map and expand only its long path/identity literals |
| `@@ -209,10 +215,16 @@` | align the new-under-2.1 map and expand only its long ref/path/identity literals |
| `@@ -263,9 +275,9 @@` | whitespace-only column alignment in the three-entry source-token map |

No semantic edit is permitted in those hunks. Any other formatter hunk or
semantic change stops the iteration.

Independent diagnostic task `01a091b7-90ec-7561-bd79-c8a0f91a522f`
established that Windows `os/exec` returns `*exec.Error` wrapping
`exec.ErrNotFound` for the existing extensionless POSIX fixture before
`CreateProcess`. The production verifier is not defective. No product or
`internal/verify/**` correction is selected.

## Approved objective

Within the one selected test file only:

1. implement the exact iteration-008 five-token workflow-normalization repair;
2. apply only Go 1.27.1's demonstrated mechanical `gofmt` projection; and
3. replace the extensionless POSIX synthetic-Cosign fixture with a cross-
   platform `TestMain` helper backed by the current test executable.

Every product, workflow, schema, builder, Docker materializer, verifier,
acceptance-test and historical-evidence byte remains unchanged.

## Path boundaries

This authority-recording session may change and commit only:

```text
README.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-AUTHORITY.md
```

Correction C2 iteration 009 implementation may change only:

```text
tests/integration/supply-chain/release_test.go
README.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-*.md (new files only)
```

The committed iteration-009 authority becomes immutable. Every pre-existing
evidence file and every absent path is forbidden. A need for another code,
test, product or verifier path stops fail-closed.

## Mandatory workflow-normalization repair

Preserve the four identity-specific replacements and replace only the global
workflow `2.2` to `2.1` normalization with the exact full manifest-schema
assertion substitution. Assert these exact counts in the original new workflow:

| New-workflow token | Count | C2 normalization |
|---|---:|---|
| `gated-v1.0.0-c2-r6-recovery` | 1 | `gated-v1.0.0-c2-recovery` |
| `release-v1.0.0-c2-r6-recovery` | 1 | `release-v1.0.0-c2-recovery` |
| `release-tooling-v1.0.0-c2-r6` | 8 | `release-tooling-v1.0.0-c2` |
| `.github/workflows/release-recovery-v1.0.0-c2-r6.yml` | 3 | `.github/workflows/release-recovery-v1.0.0.yml` |
| `test "$(jq -r .manifestSchemaVersion dist/release-manifest.json)" = '2.2'` | 1 | `test "$(jq -r .manifestSchemaVersion dist/release-manifest.json)" = '2.1'` |

The normalized R6 workflow must equal the historical C2 workflow byte-for-
byte. `# v4.2.2` must remain present and `# v4.2.1` absent. Absence,
duplication, count drift, residual difference, comment stripping, wildcard
matching or equivalent broad normalization is forbidden.

## Mandatory current-test-executable helper

1. `TestCosignCommandVerifierBindsEveryGitHubWorkflowClaim` must use
   `os.Executable()`, read and hash those exact bytes, and give that exact path
   and digest to `verify.NewCosignCommandVerifier`. Preserve the exact trusted-
   root binding.
2. Add a test-only `TestMain` helper. It may intercept only an invocation whose
   `os.Args[1:]` length is 18 and whose exact positions are:

```text
0 verify-blob
1 --bundle                         2 nonempty bundle path
3 --trusted-root                   4 nonempty trusted-root path
5 --certificate-identity           6 nonempty certificate identity
7 --certificate-oidc-issuer        8 nonempty OIDC issuer
9 --certificate-github-workflow-repository  10 nonempty repository
11 --certificate-github-workflow-ref        12 nonempty workflow ref
13 --certificate-github-workflow-sha        14 nonempty workflow SHA
15 --certificate-github-workflow-trigger    16 nonempty workflow trigger
17 nonempty manifest path
```

3. Admission additionally requires `COSIGN_YES=false`; no custom environment
   sentinel is permitted. Cleaned bundle and manifest paths must be contained
   regular files inside the test-owned temporary root.
4. Only after exact admission may the helper derive `arguments.txt` beside
   `--bundle`, write the received arguments one per line, and exit `0`.
5. Every ordinary or near-miss invocation calls `m.Run()` and must not create
   or modify `arguments.txt`. Reject prefixes, reordered/missing/extra flags,
   empty values, wrong environment, path escapes and directory, symlink or
   reparse substitutions.
6. Preserve all six claim assertions: certificate identity, OIDC issuer,
   repository, ref, workflow SHA and trigger. Add hostile/near-miss non-entry
   coverage for wrong length, command, order/position, empty values, wrong or
   absent `COSIGN_YES`, and bundle/manifest containment, without recursive test
   execution or a custom environment sentinel.
7. No runtime compilation, shell script/invocation, `PATHEXT`, build-tag
   workaround, production verifier or `internal/verify/**` change is permitted.

## Required offline validation

Implementation and fresh independent review must each use only:

```text
Go: C:\Users\ITDan\AppData\Local\Temp\pscan-06-local-toolchain-b7a0fb2-01a09168\toolchains\go1.27.1-windows-amd64\go\bin\go.exe
Module cache: C:\Users\ITDan\AppData\Local\Temp\pscan-06-local-toolchain-b7a0fb2-01a09168\gomodcache
```

The Go SHA-256 must be
`D3CCDB604EAFA6031133AEFE1A3DB24F0BB7362B857BC2125AC4E4C178B4B490`;
`go version` must be exactly `go version go1.27.1 windows/amd64`; and adjacent
`gofmt.exe` SHA-256 must be
`AB730B8446C0C2369BE901F263FB7433EA0D2D053D3A3584CA543AAFF3B667CC`.
Use fresh external `GOCACHE`/`GOTMPDIR` and set `GOPROXY=off`, `GOSUMDB=off`,
`GOTOOLCHAIN=local`, `GOENV=off`, `GOTELEMETRY=off`, `GOVCS=*:off`.

Admit the frozen module cache as bounded read/execute-only input. Prove its
literal path and read-only/RX access before and after use, run `go mod verify`,
and prove no write occurred. Reuse iteration 008's complete admission
fingerprints instead of repeatedly hashing all 1.13 GB. Any missing/changed
bounded identity, failed read-only control, failed verification, unexpected
write or download need stops closed.

Assert every native exit code and require:

1. empty final `gofmt -d` output for the selected test file;
2. targeted current-test-executable Cosign pass;
3. hostile/near-miss helper non-entry pass;
4. targeted iteration-007 preservation pass;
5. complete `./tests/integration/supply-chain` package pass;
6. `go test ./...` pass;
7. `go vet ./...` pass; and
8. applicable Linux-amd64 cross-compilation pass without executing the result.

Actual Linux runtime remains reserved for separately gated Recovery R6. Its
absence does not block this test-only local correction and cross-compilation
does not substitute for it.

## Acceptance matrix

| ID | Required result |
|---|---|
| I009-A01 | A genuinely fresh session claims only iteration 009 from this exact clean committed authority; its ID differs from the approval session. |
| I009-A02 | The test file is the only implementation-bearing path; only new iteration-009 evidence and four living documents accompany it; historical evidence and `graphify-out/**` are preserved. |
| I009-A03 | Counts `1,1,8,3,1`, exact schema substitution, workflow byte equality, `# v4.2.2` presence and `# v4.2.1` absence pass. |
| I009-A04 | Current test executable path/digest binding and exact 18-argument/environment/path helper admission pass; all six claim assertions remain. |
| I009-A05 | Ordinary and hostile near misses do not enter the helper or write its output. |
| I009-A06 | Only the exact four-hunk mechanical formatter projection occurs outside helper/normalization; no other semantic edit occurs. |
| I009-A07 | Exact offline toolchain, frozen-cache read-only/RX, `go mod verify`, format, targeted, package, repository test/vet and Linux-amd64 cross-compile gates pass without download or cache mutation. |
| I009-A08 | Exact checkout trust, clean materialization, path confinement, parsing and `git diff --check` pass. |
| I009-A09 | Author validation is non-acceptance; a genuinely fresh independent review inspects the exact committed candidate and reruns proportionate checks. |
| I009-A10 | R5 remains terminal; R6 remains unauthorized; no Docker, dependency download, network, remote, tag, workflow, signing, publication, spend or successor action occurs. |

## Forbidden scope and fresh-session requirements

No test edit or implementation occurs in this authority session. No product,
workflow, schema, builder, Docker route, production verifier, acceptance test,
PASS contract, historical decision or pre-existing evidence may change. No
semantic edit outside helper/normalization, formatter output beyond the four
hunks, runtime compilation, shell/PATHEXT/build-tag workaround, custom helper
environment sentinel, Docker, download, network, remote, tag, workflow,
signing, publication, spend, R6 or successor action is authorized.

The fresh implementation session must differ from
`01a091f1-3d3f-7a41-a541-f705cd29b687`, begin on clean `main` at this committed
bundle, preserve `graphify-out/**`, prove current-checkout trust and a separate
clean materialization, read all named controlling/prior records and the live
test, claim only iteration 009, record preflight before editing, and stop on
any identity, scope, formatter, helper, cache, test or authority mismatch.

## Approval-session exclusions

This session introduces no product, implementation, test, helper, workflow,
schema, verifier or build-logic change. Its read-only checks executed no Go
test, compilation, Docker, dependency, network or remote action. Its only
durable changes are this authority and the four synchronized living documents.
