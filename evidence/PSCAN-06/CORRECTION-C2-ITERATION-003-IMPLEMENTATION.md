# PSCAN-06 Correction C2 iteration 003 implementation

## Status

- Task: `PSCAN-06`
- Correction: `C2`, iteration `003`
- Authority commit:
  `5762d1ef6a7708629bad5f5bb33eba21ca1cdcd4`
- Session: `01a06e8c-6a86-74a2-819f-edcbc15fe6f8`
- Candidate: the implementation commit containing this record
- State: bounded author implementation; not accepted

## Implemented boundary

- `admit-image.ps1` is a closed workflow entrypoint and rejects dot-sourcing.
  Its parameter surface contains only the four committed host inputs and no
  callback, script block, command runner, executable selector or pull switch.
- The three iteration-002 mutable callback seams are absent from production.
  Host cache and host-only CRLF prerequisites execute from committed sibling
  scripts before Docker identity or commands are considered.
- The entrypoint internally constructs exactly four Docker operation argument
  lists. Only conclusive structured absence after a separately proved responsive
  engine reaches one exact canonical digest pull. Dot-sourceable
  `image-admission.ps1` contains no pull operation.
- Docker resolves only from the fixed Docker Desktop Windows path or
  `/usr/bin/docker`. Windows requires a valid Docker publisher signature;
  Linux rejects a symbolic link and group/world writable mode. The executable
  SHA-256 is recorded in the successful entrypoint result.
- The inherited child environment is cleared. Each command uses
  `ProcessStartInfo.ArgumentList`, concurrently drains raw stdout and stderr,
  enforces independent 131072-byte caps while reading, and decodes with strict
  deterministic UTF-8 only after exit and complete pipe closure.
- A monotonic 15000 ms budget covers start, reads, exit and closure. Timeout,
  overflow, read or exit uncertainty triggers complete-tree termination. A
  fixed 2000 ms cleanup grace bounds termination and pipe closure; no failure
  becomes success, absence or retry eligibility.
- The test-only model preserves the complete iteration-002 state matrix without
  entering the production admission decision. Native no-Docker fixtures use the
  already-present signed PowerShell host and temporary files only.

## Changed paths

- `.github/workflows/release-recovery-v1.0.0.yml`
- `build/release/admit-image.ps1`
- `build/release/image-admission.ps1`
- `build/release/test-image-admission.ps1`
- `docs/release/OFFLINE-VERIFICATION-RUNBOOK.md`
- `docs/release/PSCAN-06-RELEASE-PLAN.md`
- `docs/validation/PSCAN-06-VALIDATION.md`
- `docs/tasks/PSCAN-06.md`
- `docs/tasks/TRACKER.md`
- `evidence/PSCAN-06/CORRECTION-C2-ITERATION-003-PREFLIGHT.md`
- `evidence/PSCAN-06/CORRECTION-C2-ITERATION-003-IMPLEMENTATION.md`

All are within the exact 11 allowed path groups. No dependency, module,
executable, download, schema, tag, key, credential or trust root was added.

## Preserved boundaries

The exact canonical image, schema `2.1`, every earlier schema, locked product
and C1 tooling tags, dual identities, keyless signature policy, workflow order,
pre-admission host proofs, post-admission `--pull=never --network none`, later
acquisition and reproducibility order, sandbox controls and every remote,
publication, spending, credential and TruffleHog owner gate remain unchanged.

Author checks are recorded separately after execution from the exact committed
candidate. Windows author proof is not actual-Linux proof or independent
acceptance. PSCAN-06 remains open and no successor is selected or activated.
