# PSCAN-06 Correction C2 iteration 009 continuation implementation

Date: `2026-09-12`
Continuation session: `01a09293-16a9-71f1-8e93-440952d83e2c`
Starting authority commit: `f4c77e239adbb83f5af90741373b22f0bda2c4da`

## Result

`AUTHOR_CANDIDATE_COMPLETE_INDEPENDENT_REVIEW_PENDING`.

Only the newly authorized combined fixture in
`tests/unit/artifact/normalize_test.go` changed during this continuation. It
now has two subtests. The case-collision subtest creates `A` and `a`, calls
`os.Stat` on both, and skips only when `os.SameFile` proves both names resolve
to one filesystem object. On a capable filesystem it retains the exact
`artifact.RejectUnsafe` assertion and failure text. The directory-symlink
subtest retains its exact rejection assertion and only its existing
symlink-availability skip. No `runtime.GOOS` or OS-name branch exists.

The original implementation session's integration patch remained byte-for-
byte unchanged through continuation editing and Go 1.27.1 formatting:

| Path | Git blob | SHA-256 |
|---|---|---|
| `tests/integration/supply-chain/release_test.go` | `106f2efacb44177fa08d62dd86afefa1f7a52d16` | `6388688F974D9C72C722E73E304CB984AA7269C17A9161CA26864FA2847C4FDB` |
| `tests/unit/artifact/normalize_test.go` | `28794c7a30ca8e79148a3c32ca5277e3772ff0a9` | `C547659EBAE8166C7B6F3DEABB38F46A2364388DD19AA703F4F01388CD2806EA` |

The workflow source counts are exactly `1,1,8,3,1`; exact normalization equals
the historical C2 workflow byte-for-byte; the new workflow contains
`# v4.2.2` once and `# v4.2.1` zero times. The current-test-executable helper,
fixed 18-argument admission, `COSIGN_YES=false` gate, contained regular-file
checks, six claim assertions and hostile near-miss coverage remain intact.

No production, workflow, schema, builder, Docker materializer, verifier,
acceptance-test, authority or pre-existing evidence byte changed. Author
validation is recorded separately and is not acceptance.
