# PSCAN-10 Acceptance Evidence

## Result

`ACCEPT`. PSCAN-10 satisfies its activated exact-Git-object primary-coverage
boundary through the implementation and correction commits below. Acceptance
is local, evidence-based and bounded to PSCAN-10. No successor is selected or
activated.

## Bound state

- Task: PSCAN-10 only.
- Activation/base commit:
  `9053b37d799b19e3d98aeca0ae853296971ab09d`.
- Initial implementation commit:
  `47508b00b35d394412b5b2ffb2bf25b203bc9264`.
- Accepted correction commit:
  `d10df991d72e3fcc40b378830b50d1f258e16654`.
- Independent validation review: `evidence/PSCAN-10/REVIEW.md`.
- Detailed preflight, implementation and executed results:
  `PREFLIGHT.md`, `IMPLEMENTATION.md` and `VALIDATION.md` in this directory.

## Accepted capabilities

- Exact base, head, merge-base, ordered-commit, parent-edge and head-tree
  identity is bound before scanning with hostile Git behavior disabled.
- Every unique required Git blob has exactly one private admission-ledger row;
  all path/mode references, raw bytes, byte length, SHA-256, Git OID, raw class,
  preparation version and detector chunks are mutually verified.
- Raw archive, compression, application/container and ambiguous classification
  occurs on exact object bytes before transformation and is explicit non-pass.
- Admitted text and binary bytes are preserved in 90,000-byte payload windows
  with 4,019-byte overlap, exact byte mappings and per-chunk same-invocation
  inspection witnesses.
- The pinned 224-rule pack has a machine-proved conservative maximum span of
  4,020 UTF-8 bytes with zero unbounded rules.
- Named PR and release profiles bind per-object, object-count, aggregate-byte,
  detector-file, report, timeout, memory and process-count limits, with
  below/at/above evidence.
- Actual detector evidence covers text, binary, deleted history, head tree,
  merge, rename, clean and out-of-range inputs, skip paths, byte boundaries,
  maximum rule span and the exact 512 MiB declared maximum.

## Validation outcome

- Exact implementation diff: 30 paths; all activated; zero protected paths.
- `git diff --check`: pass.
- Network-disabled Linux `go test -count=1 ./...`: pass after correction.
- Network-disabled Linux `go vet ./...`: pass after correction.
- Exact 512 MiB actual-engine maximum: pass after correction in 128.38 seconds.
- Windows amd64 package compilation: pass for nine packages.
- Exact Windows detector canaries: version `8.30.1`, clean `0`, finding `11`.
- Windows/Linux logical detector parity: pass for clean/finding and rule
  identity.
- Exact reproducible source builds, rule digest, ignore digest and fresh licence
  inventory: pass.
- Changed Markdown local links: zero broken.

## Limitations

Windows host application control blocked execution of newly generated Go test
executables. That condition is not counted as a pass. Full fresh both-platform
product acceptance and first-consumer readiness remain outside PSCAN-10 and
must be established by separately selected later authority. Artifact/OCI
normalization and every final distribution, release and deployment concern
also remain non-pass or deferred.

## Interpretation

This acceptance proves the PSCAN-10 primary-engine boundary; it does not claim
that the complete scanner product is consumer-ready, distributable, signed,
published or deployed. No material required-class Gitleaks gap was established,
so PSCAN-08 eligibility was not satisfied.

## Owner decisions and gates

No remote, push, publication, release, signing, spending, credential, settings,
provider, TruffleHog, adoption, deployment, production, legal or go-live action
occurred. Every such decision remains owner-controlled. Acceptance creates no
successor authority.
