# PSCAN-06 Correction C1 iteration 005 contradiction and repair record

## Independent rejection carried forward

Candidate `a22579fd5af14473e1d49b5027f21ab591bbb589` is rejected as a
Correction C1 candidate. Independent review proved that its
`build/release/build.ps1` treated `git status --porcelain=v1` as sufficient
source trust. In an isolated clone, the tracked build driver was marked
assume-unchanged and modified. Git reported no status entries, and the
modified working-tree driver executed beyond the clean-source check.

Iteration 005 reproduced the false-clean path from the exact rejected commit:

```text
rejected_candidate=a22579fd5af14473e1d49b5027f21ab591bbb589
status_count=0
exit=1
untrusted_driver_marker=True
output_entries=4
```

The synthetic stop was deliberately placed after the historical status gate.
Its marker proves the untrusted driver reached that point; the four output
entries prove the rejected implementation also created build directories
before establishing trustworthy source state. Passing iteration 004 build,
CRLF and exact-archive evidence remains immutable history for `a22579f`, but
cannot override this false-clean path and cannot support acceptance.

## Bounded repair

Iteration 005 adds a hostile-config/index-resistant source verifier and makes
an exact committed launcher the only supported build entry. Before creating
the requested output directory or sourcing any working-tree script, the
launcher:

- clears Git directory, worktree, index, object-store and configuration
  environment overrides and disables replacement objects and lazy fetching;
- rejects filesystem-monitor, untracked-cache, ignore-stat, sparse-checkout,
  partial-clone and promisor configuration, while explicitly disabling Git
  filesystem-monitor, untracked-cache and ignore-stat shortcuts on all
  subsequent inspection commands;
- resolves one exact SHA-1 commit and tree and requires `HEAD` to equal it;
- compares every index stage-0 path, mode and object identity with the complete
  `HEAD` tree;
- rejects every assume-unchanged, skip-worktree, filesystem-monitor-valid or
  unsupported index flag;
- enumerates and rejects every untracked path, including ignored paths;
- rejects missing, non-file or reparse-point paths and unsupported Git modes;
- hashes every tracked working-tree file as raw Git blob bytes and compares it
  with `HEAD`; and
- permits a complete build from an actual Windows CRLF checkout only when each
  differing file is byte-for-byte the canonical LF-to-CRLF projection of its
  exact Git blob. Those projected checkout bytes are explicitly not trusted or
  consumed. The launcher rematerializes the entire tooling tree from raw Git
  blob objects and executes the build driver from that exact materialization.

The build repeats the same verifier before output creation and requires the
launcher's exact commit/tree binding. It constructs both deterministic source
archives from Git objects, preserves the iteration 004 path/mode/blob and
container extraction gates, and consumes no working-tree build input.

The recovery workflow first archives `build/release` from the exact workflow
SHA with conversion disabled, verifies every materialized entry with
`git hash-object --no-filters`, and executes all PowerShell entrypoints only
from that verified directory. The exact CRLF harness likewise runs its cache
canary and build launcher from committed materialized bytes; CRLF checkout
scripts are parsed only as adversarial data and are never sourced or executed.

Internal candidates `8c32eac` and `bdfd67c` were not promoted to final
evidence. The first established the full gate, and the second restored an
acquisition-path resolution accidentally removed during refactoring. Two
complete offline builds passed from `bdfd67c`, but review then strengthened
the untracked-file matrix to prove ignored residue also rejects. All complete
proofs were therefore repeated from the superseding exact candidate.

The commit containing this record is evidence only. Actual general-purpose
Linux UID/GID execution, native Windows execution, independent skeptical
rereview and all remote/signing/release gates remain open. No remote action,
tag, workflow run, signature, release, credential, paid capability,
TruffleHog or successor action is authorized or claimed.
