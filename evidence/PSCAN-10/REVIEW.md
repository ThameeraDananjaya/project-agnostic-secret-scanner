# PSCAN-10 Independent Validation Review

## Verdict

`PASS` for PSCAN-10's exact-Git-object primary-coverage boundary after the
recorded correction pass. No blocking PSCAN-10 finding remains. This verdict
does not claim complete product or first-consumer readiness.

## Review boundary

- Activation base:
  `9053b37d799b19e3d98aeca0ae853296971ab09d`.
- Reviewed implementation commits:
  `47508b00b35d394412b5b2ffb2bf25b203bc9264` and
  `d10df991d72e3fcc40b378830b50d1f258e16654`.
- The implementation diff contains 30 paths, all within PSCAN-10's activated
  allowlist. Protected PASS specifications, their source record and PSCAN-01
  through PSCAN-03 and PSCAN-09 acceptance evidence have zero changed paths.
- PASS-OUTCOME-SPEC-001 remains 19,113 bytes, 311 lines and SHA-256
  `8a034701867e366bf37adea6e4f47e3ffb7a53425463ceff4fb280bcb4296e74`.
- The review was performed as a separate evidence pass against the committed
  diff, task contract, exact bindings and executed tests; implementation claims
  were not accepted as proof by themselves.

## Findings and corrections

The first review pass rejected acceptance for three proof defects:

1. Admission rows were emitted per reference instead of exactly once per
   unique Git object. The correction commit introduced a unique-object ledger
   with separately bound path/mode references and verifier checks for missing,
   duplicate or unbound rows.
2. Raw classification did not invoke the exact application-MIME matcher used
   by pinned Gitleaks. The correction invokes
   `github.com/h2non/filetype` v1.1.3 before preparation, alongside every
   registered `github.com/mholt/archives` v0.1.2 family, and fails closed for
   archive, compression, application/container and ambiguous classes.
3. Detector accounting and finite-witness rules required tightening. The
   correction applies one profile-wide timeout, bounds aggregate report bytes
   across both streams, and removes entropy gates whose shortened bounded
   witness could otherwise weaken a rewritten rule.

The complete bounded review and validation were repeated against correction
commit `d10df991d72e3fcc40b378830b50d1f258e16654`. No blocking finding remained.

## Independent checks

- `git diff --check` passed for the exact activation-to-correction range.
- The final rule pack has 224 parsed standard-library-regex rules, zero
  unbounded rules, conservative maximum UTF-8 span 4,020 bytes and SHA-256
  `cebfe007ae88a55540e42ccb8e042848d089fe2f59bcf604c5fd4904d6460792`.
- The exact source-built binaries were reproducible across two builds per
  platform: Windows amd64 SHA-256
  `b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178`;
  Linux amd64 SHA-256
  `657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586`.
- A network-disabled, read-only Linux container passed `go test -count=1
  ./...`, `go vet ./...`, all actual-engine class/history/boundary cases and
  the exact 512 MiB declared-maximum case after corrections.
- Nine Windows amd64 test packages compiled. The exact Windows detector
  executed version, clean and synthetic-finding canaries and agreed with Linux
  on exits `0` and `11` and rule identity.
- A fresh licence collection reproduced the complete 63-module manifest and
  64 licence-file tree. The manifest SHA-256 is
  `0b8c06685c7d16fd6d685ae0e5e481d85396d2e3c64671151f07634a7bb40b60`.
- Changed Markdown local links had zero broken targets. The final tracker keeps
  every successor unselected and PSCAN-08 inactive.

## Limitations and exclusions

- Windows host application control denied launch of newly generated Go test
  executables. This is recorded as a limitation, not a pass. Windows compile
  proof and direct execution of the exact detector binary are present; a fresh
  complete both-platform product-acceptance run remains PSCAN-07 work before a
  first consumer.
- Artifact/OCI normalization, final redaction, offline-container packaging,
  signing, publication, release and consumer readiness remain outside
  PSCAN-10. Unsupported raw classes remain explicit non-pass.
- No accepted evidence established a material required-class Gitleaks gap.
  PSCAN-08 therefore remains inactive and ineligible; no TruffleHog technical
  or AGPL decision was opened.

## Owner gates

No credential, signing, spend, remote, publication, settings, provider,
deployment, production, legal or go-live action occurred. All such gates remain
reserved. Review acceptance selects or activates no successor.
