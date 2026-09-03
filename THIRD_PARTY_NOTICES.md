# Third-Party Notices

PSCAN-10 freshly revalidated the exact Gitleaks v8.30.1 MIT licence and bounded
rule config. PSCAN-06 locally assembles, but does not publish, the exact
source-built binaries. Source, config, licence, toolchain and reproducible
output digests are bound in
`build/gitleaks/manifest.json`.

Gitleaks is Copyright (c) 2019 Zachary Rice and is distributed under the MIT
licence. The exact upstream text is preserved at
`licenses/gitleaks/LICENSE-v8.30.1.txt`.

The 63 external Go modules compiled into the evaluated binary, with every
root-level licence and notice file and its digest, are preserved under
`licenses/gitleaks/modules/`. `manifest.json` in that directory binds the exact
module versions and evidence files. This inventory is intake evidence, not a
release SBOM or legal approval.

The PSCAN-10 packages directly use `github.com/mholt/archives` v0.1.2,
`github.com/h2non/filetype` v1.1.3 and their pinned transitive graph to classify
exact raw bytes before preparation. Those licence files are a subset of the
same 63-module inventory and were freshly reproduced without differences during
PSCAN-10.

The PSCAN-06 workflow source references exact commits of `actions/checkout`,
`actions/upload-artifact`, `actions/download-artifact` and `actions/attest`;
each upstream repository declares the MIT licence. It downloads Cosign v3.1.3
only by an exact asset URL plus SHA-256; Cosign declares Apache-2.0. Those tools
are build/signing infrastructure and are not linked into the runner or engine.
Their actions, execution, signing and distribution remain separately gated.

Before any third-party material is admitted, the responsible task must verify
the exact source revision, licence text, redistribution obligations, notices,
binary digest, and compatibility with the MIT-licensed runner. The resulting
notice and SBOM entries must bind to the exact accepted release asset.

PSCAN-10 found no material primary-engine gap for its admitted exact-Git
classes. TruffleHog remains absent and PSCAN-08 remains inactive and ineligible.

See [`docs/release/LICENSING.md`](docs/release/LICENSING.md).
