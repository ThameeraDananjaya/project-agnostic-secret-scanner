# Gitleaks Security Boundary

PSCAN-03's draft process boundary launches an exact absolute executable through
Go's argument-array API. It never invokes a shell. Executable and config
digests, runtime version, exact lowercase Git OIDs and local absolute paths are
validated before scanning. Engine output is bounded in private memory, supplied
only to a strict decoder, cleared, and never included in the returned result.

Success requires empty stderr, valid JSON, exit/result consistency and either
an empty finding array or fully redacted `Secret` fields. Timeout, capture
overflow, unexpected stderr, malformed output, nonstandard exit, a nonempty
clean report, an empty finding report or an unredacted candidate is explicit
non-pass.

These controls do not cure missing candidate bytes. The exact native Gitleaks
Git input omits binary blob contents from ordinary patch streams. The scanner
must therefore never report complete Git-history coverage for an affected
range. See `evidence/PSCAN-03/MATERIAL-GAP-001.md`.
