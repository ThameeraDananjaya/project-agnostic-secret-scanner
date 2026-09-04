# PSCAN-06 Correction C2 iteration 002 author validation

Date: 2026-09-05

Authority commit: `281bea031bb6bbaf1be3059977074df2a88ecdf4`

Implementation candidate: `52f7ee22ab722d7590e2d2c8326e14a0e9670462`

Implementation tree: `475f186da2a2685740430afef905c633e6123ac7`

State: local author validation passed; not independently accepted

## Executed network-free checks

### Deterministic fake-engine boundary harness

Executed after the implementation candidate was committed:

```text
pwsh -NoProfile -File ./build/release/test-image-admission.ps1
Image admission iteration-002 PASS present=NO-PULL absent=ONE-PULL failures=ZERO-PULL identity-set=UNIQUE ordering=BOUND
```

The harness parsed all three iteration scripts and recorded every simulated
host prerequisite and Docker boundary call. It proved:

- a pre-existing exact canonical identity admits with zero pull calls;
- conclusive absence on a structured responsive engine permits exactly one
  pull of the exact full digest and then a fresh inspection;
- daemon, permission, timeout, invocation, protocol, stderr, oversized-output,
  malformed, ambiguous and deceptive-absence failures produce zero pull calls;
- null, scalar, empty, alias, mixed, duplicate, wrong-repository and
  wrong-digest identities reject;
- pull and post-pull failures do not retry; and
- every rejected state records zero later action, image execution, ledger,
  dependency acquisition or artifact action.

### Exact-candidate host-only CRLF proof

A clean local-only clone was created from the candidate with checkout bytes
pinned to LF, then the committed host-only CRLF test created its own deliberate
Windows CRLF fixture:

```text
pwsh -NoProfile -File ./build/release/test-crlf-shell-payloads.ps1 \
  -SourceRepository <clean-local-candidate-clone> \
  -SourceRevision 52f7ee22ab722d7590e2d2c8326e14a0e9670462 \
  -WorkingDirectory <isolated-system-temp-directory> -Phase HostOnly
CRLF host-only regression PASS checkout-asset-CR=PROVED raw-shell-CR=REJECT normalized-shell-CR=ABSENT docker=NOT_INVOKED
```

The first attempt against the main checkout stopped before testing because the
source-trust guard detected an old ignored Graphify metadata file. A first
disposable clone also stopped because inherited Windows checkout conversion
changed raw bytes. Neither stop invoked Docker. The unrelated ignored file was
preserved; the byte-mismatched clone was discarded; and the successful run used
a new clean LF-preserving local clone. Both successful-run temporary directories
were removed after the proof.

## Static and immutable-boundary reconciliation

- `git diff --check 281bea0..52f7ee2`: PASS.
- The authority-to-candidate path list contains only the 11 authorized paths.
- The recovery workflow and mandatory admission entry point contain no
  `AllowImagePull` switch.
- The implementation contains exactly one `pull` call site, inside
  `Invoke-ReleaseImageBootstrap`; it receives only the locked full digest.
- The historical acquisition switch remains in an unchanged file, but its
  request is mapped to a legacy parameter that throws before any Docker call.
- Locked tags remain:
  `v1.0.0` = `a13c28fe7273bc8dc6545f97966a02889524eb4c` and
  `release-tooling-v1.0.0-c1` =
  `3fb7592889820fa2739a4a53588e073689621809`.
- Immutable SHA-256 identities remain:

| Object | SHA-256 |
|---|---|
| PASS-OUTCOME-SPEC-001 | `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74` |
| original C2 authority | `3CE8E20151993A8FC77890E51C5BF1557F68040CAFF2C3DB1A46237792A75D63` |
| superseded C2 author validation | `FA4D3C9E7F767575916C885449399943B7F809070F5A1E58268FFA85693AF637` |
| release-manifest schema 1.0 | `4D3F68236127EE6E2A8B908DF84F28E57CD1D0682C363D6AD4CB5CF86BA1129E` |
| release-manifest schema 1.1 | `7B89D12749424F9094B849A9D653DA580EBA14DBF84CB150EB824B203ACF3973` |
| release-manifest schema 2.0 | `AA6AE235938048CB7C5A026C5D476A6AC8D8E55F5D8761BD1E499C8D79980267` |
| release-manifest schema 2.1 | `CAA9CD26665CC3A3550AFFEA7490A0F1F277A69B10B1F1537787616E8AB973CE` |

## Deliberate non-execution and remaining gates

No Docker command, image pull, image execution, network access, download,
dependency or scanner/toolchain run, build, real-data scan, remote read/write,
tag mutation, workflow run, signing, attestation, draft or publication action
was performed.

Therefore this record does not claim the original container-cache,
acquisition, reproducibility, schema/verifier, licensing, SBOM, signature,
actual Linux or end-to-end offline checks passed for this candidate. Those
checks, a fresh independent skeptical review and all reserved owner gates remain
open. This author validation does not accept iteration 002 or close PSCAN-06,
and it does not select or activate PSCAN-07, PSCAN-08 or any successor.
