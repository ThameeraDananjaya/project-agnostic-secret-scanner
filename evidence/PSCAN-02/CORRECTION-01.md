# PSCAN-02 Correction 01

## Authority

- Corrects only the five findings in `REVIEW-01.md`.
- Baseline implementation commit:
  `b1b1bbcb7cabc3e82f9796a0ed8be0daa51a2021`.
- Task and path authority remain exactly `PSCAN-02`.
- No engine, Git-range, artifact extraction, policy, allowlist, workflow,
  remote, credential, signing, publication, spending, or consuming-project
  scope was added.

## Corrections

- Request identity becomes outcome identity only after structural validation;
  staged serialization falls back to one fixed, valid, content-free terminal
  outcome before the single stdout write.
- Commit and source-revision fields now accept only complete lowercase 40-hex
  SHA-1 or 64-hex SHA-256 Git object IDs. Content digests remain exactly
  64-hex SHA-256.
- Raw JSON presence validation now enforces all schema-required request fields,
  nested fields, release `firstRelease`, and artifact entry structure before
  typed decoding.
- Local request paths are opened only after regular-file, non-link, and payload
  bounds checks, then verified by open-handle identity. Descriptor input accepts
  bounded regular files; pipe input is accepted only when a read deadline can
  be established, and directories/devices/other unsafe types fail closed.
- Workspace operations use Go's directory-handle-bound `os.Root`, retain the
  original root identity, revalidate the path and open handle before creation
  and cleanup, and perform all child mutation through the bound root.

## Validation

| Check | Result |
|---|---|
| Linux amd64 `go test -count=1 ./...` | Pass |
| Linux amd64 `go vet ./...` | Pass |
| Linux standalone CLI integration | Pass, including safe invalid identity output and descriptor rejection |
| Linux FIFO request-path regression | Pass without opening or blocking on the FIFO |
| Windows in-process CLI suite | Pass, including invalid scan ID, mode and recovery ID |
| Windows request suite | Pass, including required presence, 40/64 Git IDs, link and descriptor rejection |
| Windows outcome suite | Pass, including 40/64 Git IDs and abbreviation rejection |
| Windows workspace suite | Pass; root-rename cases skip because the live directory handle prevents rename |
| Windows standalone targets | Compile pass; host Application Control still blocks newly linked standalone process launch |
| External Go modules | None |
| Product-code engine or shell process APIs | None |

Fresh independent acceptance remains required. PSCAN-03 remains proposed and
unselected.
