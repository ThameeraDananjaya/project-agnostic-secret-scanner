# PSCAN-03 Correction C1 Independent Review

## Result

`REJECT` on 2026-08-31. PSCAN-03 is not accepted or closed. Candidate commit
`6cd22a3615b21e5855f53af8c1d40cb51815ab27` remains non-authoritative and the
scanner CLI remains unchanged. PSCAN-04 is unselected. PSCAN-08 is inactive,
unselected and unauthorized.

## Bound candidate and preserved evidence

- Candidate parent: `a8faad04d6b6ec255a5c6222ec27841e101049f9`.
- Existing activation ancestor:
  `7459ad43313002a70d3a082c0b220012bad72870`.
- Historical gap evidence SHA-256:
  `e1b4dd6140e63c7d5ba715f00b16de342e7dec8bf353d1d923fc0993b5e9a28c`.
- Historical gap Git blob:
  `15d91e1131ea1eaa544eceba10c19942f29632db`.
- All candidate paths were inside the PSCAN-03 allowlist and the worktree was
  clean at review start.

## Exact pins

- Development/test image:
  `golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`.
- Gitleaks version: `8.30.1`; source commit:
  `83d9cd684c87d95d656c1458ef04895a7f1cbd8e`.
- Gitleaks Linux binary SHA-256:
  `c79361874b71d1b8a366773cc3cee1ade9159b1b0500e3456835fff065c7555a`.
- Gitleaks source archive SHA-256:
  `6b2638a733b85619dc80bdf28e84e4fed7e526a761ab5c148fbf67695aea2115`.
- Product/coverage config SHA-256:
  `2a9e75e17a09a1e3b06c6c232395cef85eb61448ef0298611162825f7da6c044`.
- Empty-ignore SHA-256:
  `3cdd737e34cb301ffc19861f903cf59441ff57453ddf70bd5b8c1bd63389aebe`.
- Candidate adapter/output binding: `1.1.0` / `json-v8.30.1`; finding exit
  code: `11`.

## Review commands

The independent reviewer used path- and commit-bounded Git inspection plus the
same pinned offline Docker controls recorded in `CORRECTION-C1.md`. Core state
checks were:

```powershell
git rev-parse HEAD
git rev-parse HEAD^
git merge-base --is-ancestor 7459ad43313002a70d3a082c0b220012bad72870 HEAD
git diff --name-only a8faad04d6b6ec255a5c6222ec27841e101049f9..6cd22a3615b21e5855f53af8c1d40cb51815ab27
git diff --check a8faad04d6b6ec255a5c6222ec27841e101049f9..6cd22a3615b21e5855f53af8c1d40cb51815ab27
git hash-object evidence/PSCAN-03/MATERIAL-GAP-001.md
Get-FileHash -Algorithm SHA256 evidence/PSCAN-03/MATERIAL-GAP-001.md
```

The material fragment-boundary reproduction ran with `--network none`, a
read-only root, all capabilities dropped, `no-new-privileges`, bounded PIDs,
memory and CPU, private tmpfs, the exact Gitleaks binary and the exact config.
Its synthetic file was 125,067 bytes: frame plus marker occupied 116 bytes, the
synthetic product canary began at byte 124,980, and Gitleaks split it at byte
125,000. The command was:

```powershell
docker run --rm --network none --read-only --cap-drop ALL --security-opt no-new-privileges --pids-limit 64 --memory 512m --cpus 1 --tmpfs /tmp:rw,nosuid,nodev,size=64m --tmpfs /work:rw,nosuid,nodev,size=4m -v "C:\Users\ITDan\AppData\Local\Temp\pscan-03-45abef31eed94a6aa82078c11c0eed9b\build-a\gitleaks-linux-amd64:/engine/gitleaks:ro" -v "${PWD}\rules\generic\gitleaks-v8.30.1.toml:/rules/config.toml:ro" -v "${PWD}\rules\generic\gitleaks-ignore-empty-v1.txt:/rules/ignore.txt:ro" golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452 sh -c 'mkdir /work/input; printf "PSCAN_GITLEAKS_PROJECTION_V1\nPSCAN_COVERAGE_MARKER_%064d\n" 0 > /work/input/split.dat; head -c 124864 /dev/zero | tr "\000" x >> /work/input/split.dat; printf "PSCAN_SYNTHETIC_SECRET_%064d" 0 >> /work/input/split.dat; /engine/gitleaks dir --config /rules/config.toml --gitleaks-ignore-path /rules/ignore.txt --ignore-gitleaks-allow --report-format json --report-path=- --redact=100 --no-banner --no-color --log-level error --exit-code 11 --max-archive-depth 0 --max-target-megabytes 1 --timeout 30 /work/input'
```

Observed private result: exactly one `pscan-projection-coverage` finding;
`pscan-synthetic-fixture-v1` was absent. The committed decoder therefore maps
the output to `PASS_NO_BLOCKING_FINDINGS`.

Pinned source establishes the boundary:

- `sources/file.go`: `defaultBufferSize = 100 * 1_000`.
- `sources/common.go`: `maxPeekSize = 25 * 1_000`.
- `sources/file.go` reads the 100,000-byte buffer and then calls
  `readUntilSafeBoundary`.
- Without a whitespace boundary, `sources/common.go` reads 25,000 more bytes
  and stops; the next fragment starts without overlap.

## Blocking findings

1. **Detector-fragment coverage is unproved.** One prefix marker proves that a
   file was opened, not that every later fragment was inspected. The exact
   reproduction above is a false pass.
2. **Finite-overlap remediation is not an exact proof.** The pinned TOML has
   unbounded whole-match rules, including generic API keys, private keys and
   curl credentials. A 25,000-byte overlap detects the 87-byte reproduction but
   another match can exceed the overlap and straddle both projections. Any
   finite overlap below the unbounded match span remains unprovable.
3. **Archive identity can be masked.** The frame hides ZIP/TAR/compression
   signatures while Gitleaks archive depth is zero. A source-complete raw
   classifier for every pinned archive family was not present.
4. **Binary detection was not independently proved.** The committed actual
   test combines deleted text and binary canaries. Aggregate fail can be caused
   by text alone.
5. **The rejected overlap experiment introduced unbound resource limits.** A
   hard fragment-count limit and 16 MiB report capture cannot cover the declared
   2 GiB PR and 10 GiB release profiles. That experiment was not committed.
6. **The candidate conflicts with the controlling engine-mode contract.**
   PASS-SPEC-001 section 7.2 requires Gitleaks Git mode for exact history/range
   scanning. Candidate C1 uses directory mode over a synthetic history
   projection. The correction authorization did not waive or replace this
   requirement.

## Fail-closed disposition

The smallest safe behavior with the current unbounded rule pack is
`INDETERMINATE_INCOMPLETE_COVERAGE` for any blob requiring more than one pinned
Gitleaks fragment. That behavior cannot satisfy a pass-capable 512 MiB file
path. A larger redesign needs a reviewed bounded-span rule contract, complete
raw archive classification, declared-profile resource accounting, independent
binary proof and new independent acceptance.

No design was established that simultaneously satisfies the mandatory native
Git-mode boundary and complete deleted-binary history coverage. Resolving that
contract conflict requires one consolidated owner decision; no successor or
fallback is inferred from this rejection.

No TruffleHog assessment, download, integration or enablement occurred. No
credential, remote, publication, GitHub-setting, signing, spend, deployment,
archive/OCI normalization or consuming-project action occurred.
