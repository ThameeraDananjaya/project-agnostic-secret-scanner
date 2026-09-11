# PSCAN-06 Correction C2 iteration 009 continuation authority

## Decision and state

- Task: `PSCAN-06`
- Correction: `C2`
- Iteration: `009`
- Event: owner approval of the narrow continuation required to unblock the
  inherited Windows unit-test fixture
- Decision date: `2026-09-12`
- Authority session: `01a0928b-b0af-77e2-aeca-3fe8278605ba`
- Diagnostic session: `01a09285-6aaa-7033-b65a-845a6521c67d`
- Original iteration 009 implementation session:
  `01a09228-1204-70b1-8b7e-1be433b2c496`
- Branch: `main`
- Required starting HEAD before this authority bundle:
  `a80350f93afad6e62fa82e2ca822c10207b1f89a`
- Required starting tree:
  `875fa97b57db58dbac343dd88f13861da79365d4`
- Continuation authority commit: this authority-bundle commit
- State after commit: owner-approved; unclaimed; not implemented; not
  accepted
- Sole open task: `PSCAN-06`
- Recovery R5 state: terminally failed and consumed
- Recovery R6 state: not authorized and not executable
- Successor selected or activated: none
- Spending and subscription ceiling: zero; no subscription action permitted

This is an append-only continuation authority inside already activated and
claimed PSCAN-06 Correction C2 iteration 009. It does not replace, weaken or
reinterpret the committed iteration 009 authority. It adds exactly one
implementation-bearing path for one inherited test-fixture defect. This
authority session records and commits governance only. It performs no
implementation, test retry or acceptance.

## Fail-closed starting state

Authority recording is permitted only when branch `main` is exactly at HEAD
`a80350f93afad6e62fa82e2ca822c10207b1f89a`, tree
`875fa97b57db58dbac343dd88f13861da79365d4`, the index is empty, and the only
non-ignored work is exactly:

| State | Path | Git blob | SHA-256 |
|---|---|---|---|
| modified, unstaged | `tests/integration/supply-chain/release_test.go` | `106f2efacb44177fa08d62dd86afefa1f7a52d16` | `6388688F974D9C72C722E73E304CB984AA7269C17A9161CA26864FA2847C4FDB` |
| untracked, unstaged | `evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-PREFLIGHT.md` | `739a299f3602764e2ee738ba3d42e60a296f0df6` | `D3E9F272299EB9DD842C3048E4DC2FC1931A3F349A9F7D6D76C3705A95CC6EAA` |

Both paths are implementation-session work. They must remain byte-for-byte
unchanged, unstaged and uncommitted throughout this authority session. They
are neither accepted nor incorporated into the authority commit. Every other
non-ignored tracked, staged or untracked path stops this authority session
fail-closed. Existing ignored `graphify-out/**` remains preserved and excluded
from product authority.

## Controlling authority

The following remain controlling and must be read with this continuation:

1. `AGENTS.md`;
2. `docs/spec/PASS-OUTCOME-SPEC-001.md`;
3. `docs/tasks/TRACKER.md`;
4. `docs/tasks/PSCAN-06.md`;
5. `docs/tasks/PSCAN-06-READING-MAP.md`;
6. `evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-AUTHORITY.md`; and
7. the carried-forward
   `evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-PREFLIGHT.md`.

The exact normalization, formatter, current-test-executable helper, offline
toolchain, frozen module-cache, validation, independent-review and exclusion
requirements in the original iteration 009 authority remain unchanged.

## Exact diagnostic evidence

Diagnostic session `01a09285-6aaa-7033-b65a-845a6521c67d` consumed exactly:

1. one clean-materialization targeted run; and
2. one live targeted run.

Both runs failed identically at
`tests/unit/artifact/normalize_test.go:507`. Clean and live
`tests/unit/artifact/normalize_test.go` bytes were identical. Clean and live
`internal/artifact/normalize.go` bytes were identical. Windows reported case
sensitivity disabled in both materialization roots and the temporary root.
Git blame shows the fixture has remained unchanged since PSCAN-04 commit
`1f0890878518de32a55ceb8d7b97430c4f3d2f2b`.

The inherited fixture creates directory names `A` and `a` and immediately
expects normalization to observe a case collision. On the diagnosed Windows
filesystems those names resolve to the same physical directory, so the
fixture cannot represent the adversarial condition it intends to test. The
diagnostic does not establish a production-normalizer defect and authorizes no
change to `internal/artifact/**`.

The two diagnostic runs are consumed evidence. They are not a complete matrix,
not acceptance, and not authority for another run in this authority session.

## Added implementation-bearing path and exact repair

This continuation adds exactly one implementation-bearing path:

```text
tests/unit/artifact/normalize_test.go
```

Only the existing combined case-collision and directory-symlink fixture may be
changed in that file. The authorized change is exactly:

1. split the case-collision check and directory-symlink check into clear
   subtests;
2. after creating `A` and `a`, call `os.Stat` for both paths and compare the
   results with `os.SameFile`;
3. if `os.SameFile` reports the same physical object, skip only the
   unrepresentable case-collision subtest;
4. if the two paths are distinct physical objects, preserve the exact
   `artifact.RejectUnsafe` assertion and failure semantics;
5. do not use `runtime.GOOS` or any operating-system-name branch; and
6. preserve the directory-symlink test, its exact `RejectUnsafe` assertion and
   its existing symlink-availability skip.

No other unit-test change is selected. No production behavior, artifact
normalization logic, safety boundary or required case-collision assertion on a
capable filesystem may be weakened, skipped or reinterpreted.

## Expanded continuation path boundary

The fresh continuation session may change only:

```text
tests/integration/supply-chain/release_test.go
tests/unit/artifact/normalize_test.go
README.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-*.md (new files only)
```

This list adds only `tests/unit/artifact/normalize_test.go` to the original
iteration 009 implementation boundary. The committed original iteration 009
authority, this continuation authority, every other pre-existing evidence
file, every production path including `internal/artifact/**`, and every absent
path are forbidden. A need for any other implementation, test, product,
workflow, schema, builder, Docker, verifier or evidence path stops fail-closed.

## Fresh continuation and carry-forward rule

A genuinely fresh continuation session must have an ID different from the
authority session, diagnostic session and original implementation session. It
must begin from this resulting authority commit on branch `main` while carrying
forward exactly the two bound unstaged/uncommitted paths in the fail-closed
starting-state table. This is the sole narrow exception to a clean-working-tree
start.

Before any edit or command from the offline matrix, the continuation must prove:

1. exact branch, authority HEAD, tree and direct parent;
2. the authority commit contains exactly the five authority-session paths;
3. the index is empty;
4. the only non-ignored work is the two carried paths;
5. both carried paths still match their exact Git blob and SHA-256 identities;
6. all ignored work remains confined to preserved `graphify-out/**`; and
7. the current session ID differs from all three recorded sessions.

Any mismatch stops without editing or retrying. After that claim, the session
may edit only the expanded continuation path boundary. The two carried paths
cease to be preservation-only for the continuation and remain governed by the
original iteration 009 authority; this authority does not broaden their
semantics.

## Required validation and review

The continuation must rerun the complete offline matrix required by the
original iteration 009 authority, not merely the failed targeted command. It
must use the exact pinned Go 1.27.1 executable, frozen read-only/RX module
cache, fresh external `GOCACHE` and `GOTMPDIR`, and all recorded offline
environment locks. It must assert every native exit code and require:

1. empty final formatter output for both changed Go test files;
2. the targeted current-test-executable Cosign test pass;
3. the hostile/near-miss helper non-entry coverage pass;
4. the targeted iteration-007 preservation test pass;
5. the repaired case-collision subtest to skip only when `os.SameFile` proves
   the filesystem cannot represent distinct `A` and `a` objects;
6. the directory-symlink rejection test pass or retain only its existing
   availability skip;
7. the complete unit artifact package pass;
8. the complete `./tests/integration/supply-chain` package pass;
9. `go test ./...` pass;
10. `go vet ./...` pass; and
11. the applicable Linux-amd64 cross-compilation pass without execution.

Author validation is not acceptance. A genuinely fresh independent review
must inspect the exact committed continuation candidate, verify the
case-capability logic without an OS-name branch, re-prove path confinement and
rerun proportionate complete offline checks before any acceptance record.

## Authority-session commit boundary

This authority session may change, stage and commit only:

```text
README.md
docs/tasks/PSCAN-06.md
docs/tasks/PSCAN-06-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-06/CORRECTION-C2-ITERATION-009-CONTINUATION-AUTHORITY.md
```

It must create exactly one governance authority commit. The commit's direct
parent must be `a80350f93afad6e62fa82e2ca822c10207b1f89a`, and its diff must
contain exactly those five paths. The two bound dirty paths must remain exact,
unstaged and uncommitted after the commit.

## Exclusions and stop state

This authority session performs no Go command, formatting, test, vet,
cross-compilation, Docker action, dependency acquisition, network request,
remote read or mutation, tag operation, workflow dispatch, signing,
attestation, draft, release, publication, spending, subscription, Recovery R6
or successor action. It does not implement or modify either existing iteration
009 work path.

The future continuation receives no authority for Docker, network, remote,
tag, workflow, signing, publication, spending, subscription, Recovery R6 or
successor work. Recovery R5 remains terminal and consumed. PSCAN-06 remains
the sole open task and remains unaccepted overall. After the one authority
commit and exact preservation proof, this authority session stops.
