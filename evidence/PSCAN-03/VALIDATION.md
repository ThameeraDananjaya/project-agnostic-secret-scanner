# PSCAN-03 Validation Record

## Status

Stopped fail-closed on 2026-08-31. This is implementation and gap evidence,
not an acceptance record. No `ACCEPTANCE.md` exists for PSCAN-03.

## Authority and scope

- Activation commit: `7459ad43313002a70d3a082c0b220012bad72870`.
- Only paths allowed by `docs/tasks/PSCAN-03.md` changed.
- Canonical and external PASS-SPEC-001 remained byte-identical: 57,769 bytes,
  977 lines, SHA-256
  `78a1ad9a9abfde577a50ae9aa467b5b82f733162c8187ed58a2bc300952b692d`.
- No runner CLI, workflow, action, remote, credential, signing, publication,
  consuming-project material, TruffleHog material or successor activation was
  introduced.

## Exact intake and build checks

- Official release/tag/source/licence/security facts were revalidated before
  intake.
- Exact commit source archive, config, licence, `go.mod` and GoReleaser file
  hashes matched `build/gitleaks/manifest.json`.
- `build/gitleaks/build.ps1` reproduced both accepted evaluation hashes:
  Windows amd64
  `a9e923bdde0e353057f7b71c2b14f1e1b96016076f05fc81c93dd605f46525ea`
  and Linux amd64
  `c79361874b71d1b8a366773cc3cee1ade9159b1b0500e3456835fff065c7555a`.
- Two prior clean-room builds per platform were byte-identical.
- All 63 compiled external module identities were enumerated; all 64 preserved
  licence/notice files matched the 64 manifest digests. Gitleaks's own exact
  MIT licence is preserved separately.

## Passing implementation checks

The following command passed on Windows amd64 with Go 1.27.0:

```text
go test ./tests/unit/engine ./tests/unit/gitinput ./tests/integration/gitleaks
```

It covers private-output confinement, pre-execution digest refusal, bounded
capture, runtime-version refusal, strict JSON/exit/redaction classification,
timeout, argument metacharacters, safe bare cloning, exact ancestry/range/tree
binding, and complete regular-file manifests. No test was skipped.

`go vet` passed for the new implementation and PSCAN-03 test packages. The
PSCAN-03 acceptance test package compiled for both Windows amd64 and Linux
amd64.

## Non-passing complete suite

`go test ./...` is not a pass. The host's Application Control policy refused
freshly compiled Windows test/runner executables in
`tests/acceptance/gitleaks` and the pre-existing
`tests/integration/contract` process tests. This matches the accepted PSCAN-02
host limitation and is not treated as evidence of successful execution.

The exact Linux Gitleaks evaluation binary was executable in the local
`docker-desktop` WSL environment. Its SHA-256 matched the manifest, `version`
returned `8.30.1`, clean directory smoke returned exit 0, synthetic finding
smoke returned configured exit 11, reports were valid private JSON, stderr was
empty and the canary did not cross captured stdout/stderr. That environment has
no Git executable, so it could not repair or fully execute native Git-history
acceptance.

## Material gap reproduction

`TestSafeBareCloneAndExactRangeBinding` passed with Git
`2.55.0.windows.3`. Its generated binary add/delete fixture proves the exact
native `git log -p -U0` input form omits the fabricated candidate bytes. See
`MATERIAL-GAP-001.md`.

## Acceptance decision

Rejected/incomplete. Missing binary contents in an exact required Git-history
range are a material coverage gap. PSCAN-03 remains blocked and cannot support
`PASS_NO_BLOCKING_FINDINGS` for affected history. No successor is selected.
