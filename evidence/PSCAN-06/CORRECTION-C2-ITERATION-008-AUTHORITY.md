# PSCAN-06 Correction C2 iteration 008 authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Iteration: `008`
- Title: identity-specific workflow-preservation normalization repair
- Event: owner-directed selection of the minimal bounded correction after the
  iteration-007 independent rejection
- Date: `2026-09-11`
- Approval session: `01a0919e-5c8a-7f92-bece-1f18193a788f`
- Branch: `main`
- Required starting HEAD before this authority bundle:
  `b62394ef172cb92b6394ef313993a2a6a5613b5b`
- Required starting tree:
  `35fb2f9362325719a83f8644b9e3a035c0db92c4`
- Required starting HEAD direct parent:
  `230e1761f7c45d9629cecf360e7498b33d64ab6f`
- Correction authority commit: this authority-bundle commit
- Pre-change tracked, staged and non-ignored untracked Git status: clean
- Pre-change graphify-aware exact checkout verification: 380 tracked paths,
  362 raw-equal, 18 canonical whole-file LF-to-CRLF projections, exact index
  and flags, zero unsupported Git configuration, zero non-ignored untracked
  paths, and 66 preserved ignored paths all under `graphify-out/**`
- Preserved ignored workspace material: `graphify-out/**`; excluded from
  product authority and not refreshed, deleted, moved or staged
- State after commit: owner-approved; unclaimed; not implemented; not accepted
- Recovery R5 state: terminally failed and consumed
- Recovery R6 state: not authorized and not executable
- Sole open task: `PSCAN-06`
- Successor selected or activated: none

This is a bounded next iteration inside already activated PSCAN-06 Correction
C2. It selects no successor. This session records authority only. A genuinely
fresh session with an ID different from this approval session must claim only
Correction C2 iteration 008 from this committed bundle before implementation.

## Controlling and immutable inputs

- Controlling contract: `PASS-OUTCOME-SPEC-001`
- Contract committed-byte SHA-256:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`
- Iteration-007 authority committed-byte SHA-256:
  `A17520815841415AD06357682423CF7C635B32AAFAB7EA2B0303E2540FE709E8`
- Iteration-007 preflight committed-byte SHA-256:
  `31624B5BAC5DB0EA751EB739795BE302469BE543CD56F675D4FBCC2729733FAC`
- Iteration-007 implementation committed-byte SHA-256:
  `F251946635F474DD4AF39071EBBE14CED0F48FA52AD5FF5AF4D130822D59083B`
- Iteration-007 independent rejection committed-byte SHA-256:
  `E8B90DDE62848577E129B7E2857FE6D0BF56BBF691440A06636013BE458FAC21`
- Rejected iteration-007 candidate:
  `230e1761f7c45d9629cecf360e7498b33d64ab6f`
- Rejected iteration-007 candidate tree:
  `5d74628d76f95e150c0ea1de51b4892af2350a7d`
- Direct iteration-007 authority parent:
  `b7a0fb211e4631f915e368174ae58022d1519182`
- Recovery R5 terminal failure committed-byte SHA-256:
  `53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0`
- Accepted iteration-006 candidate:
  `cf1679f9bca24887340fa4060f37d1ebff21f305`
- Accepted iteration-006 candidate tree:
  `6cb0198c71881baff9c601c31ad6aa84217f4a49`

The iteration-007 candidate and its rejection remain immutable evidence. The
new R6 workflow, schema 2.2, builder, materializer and both verifier routes were
found coherent by independent review. All other recorded exact source-trust,
path-confinement, parsing and available-host no-Docker checks were green. This
authority does not reinterpret the rejected candidate as accepted.

## Exact bounded defect

The sole blocking defect is in
`TestIteration007RepositoryIdentityAgreementAndPreservation` within
`tests/integration/supply-chain/release_test.go`. Its workflow normalization
ends with:

```go
normalizedWorkflow = strings.ReplaceAll(normalizedWorkflow, "2.2", "2.1")
```

The new workflow contains two `2.2` substrings while the old workflow contains
one. Only the exact manifest-schema assertion is an intended version change.
The other occurrence is the unchanged action-pin comment `# v4.2.2`, present
identically in both workflows. The global replacement changes only the new
workflow's comment to `# v4.2.1`, so the preservation equality fails
deterministically even though the workflow derivation itself is correct.

## Approved objective

Repair only that preservation-test normalization. Normalize the intended R6
workflow identity, schema-operation and reference tokens to their C2 values
with exact token boundaries and exact occurrence assertions. Preserve the
unchanged action-version comment and every coherent iteration-007 product,
workflow, schema, builder, materializer, verifier and acceptance-test byte.

Author validation and a separate fresh independent review must prove the
previous false failure is removed without weakening the preservation test or
changing any implementation identity. Recovery R6 remains separately gated.

## Authority-session path boundary

This authority-recording session may change and commit only:

```text
README.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-008-AUTHORITY.md
```

It may not edit the test or any product, workflow, schema, build, verifier or
other evidence byte.

## Allowed implementation paths

Correction C2 iteration 008 may change only:

```text
tests/integration/supply-chain/release_test.go
README.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-008-*.md (new files only)
```

The committed iteration-008 authority itself becomes immutable at this
bundle's commit and may not be edited by implementation or review. Every
pre-existing evidence file and every absent path is forbidden. A demonstrated
need for another implementation or test path stops fail-closed and returns
exact evidence to the owner.

## Mandatory repair design

1. Modify only the workflow-normalization portion of
   `TestIteration007RepositoryIdentityAgreementAndPreservation`, plus imports
   only if strictly required by that bounded test change. Do not change another
   test, fixture, product assertion or verification route.
2. Remove the global workflow replacement of every `2.2` substring. Do not
   replace it with another unbounded version, numeric, comment or substring
   transformation.
3. Normalize only this exact intended token set from the original new workflow
   text, proving the exact pre-replacement occurrence count for each:

| New-workflow token | Expected count | C2 normalization |
|---|---:|---|
| `gated-v1.0.0-c2-r6-recovery` | 1 | `gated-v1.0.0-c2-recovery` |
| `release-v1.0.0-c2-r6-recovery` | 1 | `release-v1.0.0-c2-recovery` |
| `release-tooling-v1.0.0-c2-r6` | 8 | `release-tooling-v1.0.0-c2` |
| `.github/workflows/release-recovery-v1.0.0-c2-r6.yml` | 3 | `.github/workflows/release-recovery-v1.0.0.yml` |
| `test "$(jq -r .manifestSchemaVersion dist/release-manifest.json)" = '2.2'` | 1 | `test "$(jq -r .manifestSchemaVersion dist/release-manifest.json)" = '2.1'` |

4. Prove the normalized new workflow equals the historical C2 workflow
   byte-for-byte and that the unchanged `# v4.2.2` comment remains present and
   unchanged. The test must fail if an expected token is absent, duplicated,
   over-broadly transformed or leaves any residual difference.
5. Preserve the existing schema 2.2 normalization and all historical
   preservation, cross-version, mutation and identity-agreement checks unless
   the exact test-only defect proves otherwise. This authority does not permit
   a looser comparison, ignored line, regular-expression wildcard, comment
   stripping or normalization of arbitrary version text.
6. Preserve these critical committed bytes exactly:

| Path | Git blob | Committed-byte SHA-256 |
|---|---|---|
| `.github/workflows/release-recovery-v1.0.0.yml` | `3aa42627628b8b5298d854d29b3880cb19dc36ff` | `C5F40F1B32E87C005FE33EE607AF7E3D19EE4F0173C21619E21158C31AA0BDB4` |
| `.github/workflows/release-recovery-v1.0.0-c2-r6.yml` | `ff5b55f32ced9f382b30e98746f9ed5e21f13a73` | `CFE3FB919CA26A7AE6CF2A8F56D22BA61B96242298E010940E7B863053C9FA1A` |
| `contracts/release-manifest/schema-2.1.json` | `f8b5232b8667636fec8f4f0c7f083e9a0297b0cb` | `CAA9CD26665CC3A3550AFFEA7490A0F1F277A69B10B1F1537787616E8AB973CE` |
| `contracts/release-manifest/schema-2.2.json` | `8c20c4943072f80d57d3d3b9a3b094ad6c524fca` | `347E424F23F48CF25A409328E3A2E6773F589104AC0129216977DE6D01ED872A` |
| `build/release/build.ps1` | `6f84f86cf6646caa1da0d2ff7567c9c7697549e1` | `717914E2D792D37C639F1D2AE5715A637DA0AE6E7CF469DF15DFED90901EFFD7` |
| `build/release/docker-execution.ps1` | `a1a3d616db37dcb22e401b0008f4c95a57b6b440` | `DB6D99F284316245A1D841547B4880BE6B72EE53F35CB88B0AC79C81249F0B9F` |
| `build/release/cmd/release-verifier/main.go` | `0ce02f38c0f57014609faaf292637d4dd0c3e3ff` | `999D3F85DE6E8BCBC053F726C7ED888CB6D41B2C6892125F301A021F5CE75CDF` |
| `internal/verify/release.go` | `fff07d8c30f1dcb676e0d18234b018f406c8431d` | `1020771E3F32E02BB4CF601FB2980A9656CB1DC5C4A5F4283434BCB23FBC7A7F` |
| `tests/acceptance/supply-chain/release_test.go` | `1379d07d40687ef0dca1b102c7fe873fec6cd268` | `B04BB8D575B59E29E7B9F34E94BCD418F45C54E1127ABE68A0AE22AA22B8C0A8` |

## Owner-supplied offline Go validation route

The owner supplied this exact local future validation route:

```text
Go executable:
C:\Users\ITDan\AppData\Local\Temp\pscan-06-local-toolchain-b7a0fb2-01a09168\toolchains\go1.27.1-windows-amd64\go\bin\go.exe

Frozen module cache:
C:\Users\ITDan\AppData\Local\Temp\pscan-06-local-toolchain-b7a0fb2-01a09168\gomodcache
```

This authority session confirmed only that both literal paths exist and that
the current `go.exe` SHA-256 is
`D3CCDB604EAFA6031133AEFE1A3DB24F0BB7362B857BC2125AC4E4C178B4B490`.
It did not execute Go, enumerate or mutate the module cache, or treat the
toolchain as implementation evidence.

The fresh implementation and independent review sessions must each re-prove
the literal executable path, SHA-256 and exact `go version` result `go version
go1.27.1 windows/amd64`; use the literal module-cache path only as frozen read-only
input; take exact before/after evidence that it did not change; use fresh
external `GOCACHE` and `GOTMPDIR` directories outside the repository; and set:

```text
GOPROXY=off
GOSUMDB=off
GOTOOLCHAIN=local
GOENV=off
GOTELEMETRY=off
GOVCS=*:off
```

No fallback `go` from `PATH`, automatic toolchain acquisition, module download,
network lookup or mutation of the frozen module cache is allowed. Any missing,
changed, not demonstrably frozen, incomplete or provenance-ambiguous input stops the
check rather than weakening it.

## Acceptance matrix

| ID | Required result |
|---|---|
| I008-A01 | A genuinely fresh session claims only iteration 008 from the exact clean committed authority bundle; its session ID differs from `01a0919e-5c8a-7f92-bece-1f18193a788f`. |
| I008-A02 | Candidate diff changes only the exact test normalization, new append-only iteration-008 evidence and the four named living-state documents. Every historical evidence file and ignored `graphify-out/**` path is preserved and unstaged. |
| I008-A03 | The five exact intended workflow substitutions have the required source counts; no global `2.2` replacement or equivalent broad transformation remains in workflow normalization. |
| I008-A04 | The normalized R6 workflow equals the historical C2 workflow byte-for-byte; `# v4.2.2` remains unchanged; the former sole residual `# v4.2.1` is absent. |
| I008-A05 | All critical workflow, schema, builder, Docker materializer, verifier and acceptance-test identities in this authority remain byte-identical. The rejected iteration-007 candidate and review record remain immutable. |
| I008-A06 | With the exact offline Go route and explicit native exit-code checks, formatting is clean; the exact preservation test, the complete integration package, `go test ./...` and `go vet ./...` pass without network or module-cache mutation. Any unavailable input or unrelated failure stops fail-closed and is not repaired under this authority. |
| I008-A07 | Exact graphify-aware checkout verification, a separate exact clean materialization, path confinement, `git diff --check`, parsing and the proportionate available-host no-Docker/static checks pass with no surviving marked process or fixture. |
| I008-A08 | Author validation is recorded as non-acceptance. A separate genuinely fresh independent review inspects the exact committed candidate and reruns the targeted preservation proof plus proportionate offline checks before any local acceptance record. |
| I008-A09 | Recovery R5 remains terminal and consumed; Recovery R6 remains unauthorized; no Docker, dependency download, network, remote, tag, workflow, signing, publication, spend or successor action occurs. PSCAN-06 remains open and unaccepted overall. |

## Forbidden scope and reserved gates

- No implementation or test edit in this authority session.
- No edit to any workflow, schema, builder, Docker/materialization route,
  verifier, acceptance test, product source, scanner, PASS contract,
  historical decision or pre-existing evidence file.
- No weakening, deletion, skipping or broad normalization of the preservation,
  cross-version, mutation or identity-agreement tests.
- No Docker command, dependency download, network action, remote read or
  mutation, push, fetch, tag action, ruleset/setting action, workflow dispatch,
  rerun or retry.
- No credential, signing key, paid capability, subscription, nonzero spend,
  TruffleHog work, signing, attestation, draft, release, deployment or
  publication action.
- No deletion, movement, refresh or staging of preserved ignored
  `graphify-out/**`.
- No PSCAN-07, PSCAN-08 or successor selection, activation, claim or work.

## Fresh implementation-session requirements

The implementation session must:

1. have a session ID different from
   `01a0919e-5c8a-7f92-bece-1f18193a788f`;
2. begin on exact branch `main` at this committed authority bundle with clean
   tracked, staged and non-ignored untracked state;
3. preserve all ignored `graphify-out/**` and establish both graphify-aware
   current-checkout trust and a separate exact clean materialization;
4. read the controlling contract, PSCAN-06 task, reading map, tracker,
   iteration-007 authority/implementation/rejection and this authority
   completely;
5. claim only PSCAN-06 Correction C2 iteration 008 and record its bounded
   preflight before editing; and
6. stop fail-closed on any identity, path, token-count, preservation,
   toolchain, offline-test or authority mismatch.

## Approval-session exclusions

This approval session introduces no product, implementation, test, workflow,
schema, verifier or build-logic change. It performs no Go test, Docker command,
dependency acquisition, network or remote action, tag operation, workflow run,
signing, attestation, draft, release, deployment, publication, purchase,
subscription or successor action. Its only durable changes are this bounded
authority bundle and synchronized living PSCAN-06 state.
