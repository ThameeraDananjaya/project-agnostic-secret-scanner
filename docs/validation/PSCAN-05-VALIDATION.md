# PSCAN-05 Validation Boundary

## Acceptance mapping

| PSCAN-05 requirement | Implemented validation |
|---|---|
| Fixed non-overridable precedence | Unit matrix for binding failure, revocation-first denial, mandatory credential blocking, project denial and exact exception |
| Signed project projection | Exact-byte digest test, RFC 8032 public vector, trust-domain/adapter/family/version rejection, duplicate/unknown member rejection |
| Narrow expiring exceptions | Distinct owner/approver, unique exception identity, exact scope/context/evidence, complete invalidation set, 30-day hard maximum, project maximum and use-time expiry |
| Independent compatibility | Unknown-major, duplicate-minor, conflicting-retirement, two-cycle, 90-day and current-major replacement checks |
| Receipt/revocation reference verification | Every-binding comparison including signed projection/schema/time facts, 30-day deadline, trusted previous-receipt anchor and current head, full genesis-to-trusted-head revocation chains, scanner-owned global-schema authority checks, changed binding, revocation, rollback and divergence |
| Project isolation | Parallel and sequential workspace identities plus concurrent transient policy projections |
| No custody or authority expansion | Static API/content audit: public keys only; no signing, private key, issuance, store, promotion, project schema or project instance |

## Required execution

Implementation validation uses the exact Go 1.27.1 Linux amd64 archive whose
official SHA-256 is recorded in `evidence/PSCAN-05/PREFLIGHT.md`. The bounded
PSCAN-05 unit, integration and acceptance packages must pass repeatedly;
`go vet` must pass for touched packages; the previous request tests must remain
green; all packages and tests must compile for Windows amd64; JSON contracts
must parse; and static path/content/protected-hash checks must pass.

Native Windows test execution and fresh independent acceptance remain required
before PSCAN-05 may close. When host application control prevents execution of
the unsigned Go assembler, network-disabled Linux cross-compilation may produce
the Windows amd64 test executables, but those executables must then run natively
and pass on Windows; compile-only success is not a substitute.
