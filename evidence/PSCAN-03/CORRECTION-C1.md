# PSCAN-03 Correction Iteration C1

## Status

Implementation candidate complete on 2026-08-31; independent review and final
acceptance pending. Task `PSCAN-03` remains the only activated task.

## Authority and starting state

- Activation commit:
  `7459ad43313002a70d3a082c0b220012bad72870`.
- Required starting `main` HEAD:
  `a8faad04d6b6ec255a5c6222ec27841e101049f9`.
- Pre-change state: clean `main`.
- `7459ad4` was verified as an ancestor of `a8faad0`.
- Work is confined to the allowed paths in `docs/tasks/PSCAN-03.md`.
- Historical material-gap record SHA-256 remains
  `e1b4dd6140e63c7d5ba715f00b16de342e7dec8bf353d1d923fc0993b5e9a28c`;
  it was not rewritten.

## Correction decision

Native Gitleaks Git mode remains incapable of complete deleted-binary history
coverage and is retained only as a compatibility surface that always returns
`INDETERMINATE_INCOMPLETE_COVERAGE`.

The authoritative correction:

1. Revalidates exact base, head, merge base and ordered in-range commits.
2. Enumerates every parent edge with fixed raw Git metadata, renames disabled,
   external diffs disabled and text conversion disabled.
3. Admits both regular-file pre-image and post-image blob OIDs for each edge.
4. Independently enumerates every regular blob in the exact head tree.
5. Rejects unsafe paths, modes, object types, case conflicts, duplicates,
   individual-size limits, blob-count limits and aggregate-byte limits.
6. Writes a byte-exact private projection plus a deterministic path-preserving
   framed projection containing a unique coverage marker followed by every
   original blob byte.
7. Reopens every projected file and verifies per-file and aggregate SHA-256
   bindings before launch.
8. Runs product rules and the private coverage rule in the same pinned Gitleaks
   v8.30.1 invocation. Exactly one coverage finding per admitted file is
   mandatory. Missing, duplicate, unbound or malformed coverage is non-pass.

Gitleaks v8.30.1 remains the sole detector. No TruffleHog material was assessed,
downloaded, integrated or enabled.

## Exact pins

- Docker development/test image:
  `golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`.
- Container Go: `go1.27.0 linux/amd64`.
- Container Git: `2.39.5`.
- Gitleaks: `v8.30.1`, commit
  `83d9cd684c87d95d656c1458ef04895a7f1cbd8e`.
- Gitleaks Linux binary SHA-256:
  `c79361874b71d1b8a366773cc3cee1ade9159b1b0500e3456835fff065c7555a`.
- Upstream source archive SHA-256:
  `6b2638a733b85619dc80bdf28e84e4fed7e526a761ab5c148fbf67695aea2115`.
- Adapter: `1.1.0`; output binding: `json-v8.30.1`; finding exit: `11`.
- Product/coverage config SHA-256:
  `2a9e75e17a09a1e3b06c6c232395cef85eb61448ef0298611162825f7da6c044`.
- Explicit empty-ignore file SHA-256:
  `3cdd737e34cb301ffc19861f903cf59441ff57453ddf70bd5b8c1bd63389aebe`.

## Exact validation environment

The following PowerShell variables and Docker prefix were used:

```powershell
$root = (Get-Location).Path
$gitleaksHost = 'C:\Users\ITDan\AppData\Local\Temp\pscan-03-45abef31eed94a6aa82078c11c0eed9b\build-a\gitleaks-linux-amd64'
docker run --rm --network none --read-only --cap-drop ALL --security-opt no-new-privileges --pids-limit 512 --memory 3g --cpus 3 --tmpfs /tmp:rw,exec,nosuid,nodev,size=2g --tmpfs /root/.cache:rw,nosuid,nodev,size=2g -e GOTOOLCHAIN=local -e GOCACHE=/root/.cache/go-build -e PSCAN_GITLEAKS_BINARY=/engine/gitleaks -e PSCAN_GIT_BINARY=/usr/bin/git -e PSCAN_GITLEAKS_CONFIG=/src/rules/generic/gitleaks-v8.30.1.toml -e PSCAN_GITLEAKS_IGNORE=/src/rules/generic/gitleaks-ignore-empty-v1.txt -v "${root}:/src:rw" -v "${gitleaksHost}:/engine/gitleaks:ro" -w /src golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452 sh -c 'go test -count=1 ./... && go vet ./...'
```

Result: all repository tests and `go vet ./...` passed. No test was skipped in
the pinned actual-Gitleaks run. Scanner execution had Docker networking set to
`none`, a read-only container root, all Linux capabilities dropped,
`no-new-privileges`, bounded PIDs, memory and CPU, and tmpfs test/cache paths.

The race gate used the same Docker security and pin set:

```text
go test -race -count=1 ./tests/unit/gitinput ./tests/unit/engine ./tests/integration/gitleaks
```

Result: passed. Two pre-existing one-second helper-process timeouts were raised
to five seconds after race instrumentation made them expire; their security
expectations were unchanged.

The deterministic actual-Gitleaks acceptance command was run twice:

```text
go test -count=1 -run TestPinnedGitleaksProjectionCoverage -v ./tests/integration/gitleaks
```

Both runs produced the same bindings:

| Fixture | Plan digest | Exact-content digest | Coverage digest | Entries | Result |
|---|---|---|---|---:|---|
| Deleted text and NUL-binary history | `c1cae2d055692d7c3e46b4f494b183a72af9a64fa0a174607baf694f7add65a5` | `61667c379350bfe32e2c09a7916849a2127a83e210a1a80d9ccafb8d40698be2` | `9c15e5f828fdb3a5e3c8eb62ee3bb4f7997de5608136a7a3dfd4eb60df182e7e` | 6 | `FAIL_FINDING_DETECTED` |
| Clean exact range; finding only before base | `3333621df8e4cddfd0e4f5fca2325bbc4581ad996380a1a171434c1d2f115029` | `c63a9177254a8c38f581f7d881cb307fbe7f3d79c922228df3b196d21b839b93` | `03881ed266acc3a789fe8f4b13949a01d58359c97f00f7b7940d189c7d6810d4` | 3 | `PASS_NO_BLOCKING_FINDINGS` |
| Upstream path-allowlisted `.bin` | `5e668e3a08e7dd9ccb0e7c5868fa67daaff5c40d7c2035db4fdbd52ef1b283ea` | `b3222926c0224d654a2afe2857753cb80e4c46de80386099be133533bed924fb` | `62b2d3ef0713855c1d36908d1fef2ea2dff6075e67edea229343513098896d7f` | 3 | `INDETERMINATE_INCOMPLETE_COVERAGE` |

## Covered adversarial criteria

- deleted text and NUL-containing binary canaries;
- exact base/head/merge-base, ordered ranges and exact head trees;
- multi-parent merges and rename-as-delete/add behavior;
- out-of-range exclusion;
- hostile local and ambient Git configuration, attributes, external diff and
  hook canaries;
- unsafe portable paths, case collisions, symlinks and gitlinks;
- per-file, count and total-byte limits;
- missing, additional and mutated projections;
- exact engine/config/ignore/version/output bindings;
- private raw output, full engine redaction and content-free results;
- clean inputs and Gitleaks path-skip non-pass;
- repeated deterministic plan, exact-content and coverage digests.

## Exclusions and successor state

No archive/OCI normalization, final Docker packaging, public CLI integration,
credential handling, remote action, GitHub setting, signing, publication,
spending or consuming-project integration was performed. PSCAN-04 remains
unselected. PSCAN-08 remains inactive, unselected and unauthorized.
