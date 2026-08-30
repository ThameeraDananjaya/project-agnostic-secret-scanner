# Third-Party Notices

PSCAN-03 admits the exact Gitleaks v8.30.1 MIT licence and default rule config
for evaluation. No Gitleaks binary is committed or released. The exact source,
config, licence, toolchain and reproducible output digests are bound in
`build/gitleaks/manifest.json`.

Gitleaks is Copyright (c) 2019 Zachary Rice and is distributed under the MIT
licence. The exact upstream text is preserved at
`licenses/gitleaks/LICENSE-v8.30.1.txt`.

The 63 external Go modules compiled into the evaluated binary, with every
root-level licence and notice file and its digest, are preserved under
`licenses/gitleaks/modules/`. `manifest.json` in that directory binds the exact
module versions and evidence files. This inventory is intake evidence, not a
release SBOM or legal approval.

Before any third-party material is admitted, the responsible task must verify
the exact source revision, licence text, redistribution obligations, notices,
binary digest, and compatibility with the MIT-licensed runner. The resulting
notice and SBOM entries must bind to the exact accepted release asset.

Gitleaks acceptance is blocked by `GITLEAKS-GAP-001`. TruffleHog remains absent and may not
be assessed, downloaded, integrated, distributed, or enabled unless PSCAN-08
first becomes eligible through accepted material-gap evidence and separate
technical and AGPL owner approval.

See [`docs/release/LICENSING.md`](docs/release/LICENSING.md).
