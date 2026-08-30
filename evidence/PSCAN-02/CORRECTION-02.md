# PSCAN-02 Correction 02

## Authority and result

This correction addresses only the blocker and ambiguity recorded by the
independent `REVIEW-02.md`. Task and path authority remain PSCAN-02.

- Schema and runtime now require canonical RFC 3339 UTC text ending in `Z`.
- Positive, negative, and non-canonical `+00:00` offsets fail input integrity.
- The artifact-manifest digest now uses an explicit domain-separated,
  length-prefixed binary preimage with big-endian integers, documented in the
  architecture and referenced by the schema.
- No successor or forbidden implementation scope is introduced.

## Validation

| Check | Result |
|---|---|
| Linux amd64 `go test -count=1 ./...` | Pass |
| Linux amd64 `go vet ./...` | Pass |
| Canonical UTC positive case and `+04:00`, `-05:00`, `+00:00` negative cases | Pass on Linux |
| Artifact digest framing and determinism regressions | Pass on Linux |
| Windows request, contract and runner targets | Compile pass |
| Fresh Windows test-binary execution | Host Application Control blocked the new binary hashes; not counted as pass |

Fresh independent acceptance is required before closeout. The Windows
standalone-execution limitation remains explicit and no bypass was attempted.
