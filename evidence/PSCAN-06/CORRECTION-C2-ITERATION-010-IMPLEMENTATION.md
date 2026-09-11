# PSCAN-06 Correction C2 iteration 010 implementation

Date: `2026-09-12`
Implementation session: `01a092b2-d833-74f3-99c0-a2def0764153`
Starting authority commit: `ad77e68f8ae6710c24ba818bfb04c8ecfec5f4e4`

## Result

`AUTHOR_CANDIDATE_COMPLETE_INDEPENDENT_REVIEW_PENDING`.

The sole implementation-bearing path is
`tests/integration/supply-chain/release_test.go`. For an otherwise exactly
admitted invocation, `cosignHelperArgumentsPath` now derives output exactly as:

```go
filepath.Join(filepath.Dir(filepath.Clean(arguments[2])), "arguments.txt")
```

The trusted admission root remains derived from cleaned argument 4. All
existing root, containment, regular-file, evaluated-symlink, exact-position,
environment, claim, exclusive-create and non-recursion behavior remains in
place.

A new positive regression uses regular files at
`<root>/trusted-root.json`, `<root>/release-manifest.json` and
`<root>/nested/release-manifest.sigstore.json`. It proves direct helper
admission returns only `<root>/nested/arguments.txt`, `testMainExitCode`
returns success without `m.Run`, the exact ordered 18 arguments and all six
claim positions are captured there, and the distinct trusted-root-sibling
output remains absent.

The hostile regression now exercises bundle, trusted-root and manifest file
symlinks; trusted-root escape; all three input-directory substitutions; and a
pre-existing exact output whose bytes must remain unchanged. Existing wrong
length, command, flag position, empty value, bundle/manifest escape,
`COSIGN_YES`, irregular-mode and non-recursion cases remain.

After pinned formatting, the selected test has working blob
`1c8ba4957bd82f1092f054a47a4a7207ea7de655`, SHA-256
`D2A65243613A9C531D90AEA662106404A520F148B17435A454511C7DFA6BBED4`
and 30,311 bytes. No production, workflow, schema, builder, Docker
materializer, verifier, unit-test, acceptance-test, authority or predecessor
evidence file changed.

Author validation is recorded separately. This result is not acceptance.
