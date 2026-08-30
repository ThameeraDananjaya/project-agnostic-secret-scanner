# Generated History Fixtures

History repositories are generated only in private test directories so no
`.git` authority, object database, hook or repository identity is committed.
The generators in `tests/unit/gitinput` and `tests/integration/gitleaks` create:

- a clean range base;
- in-range synthetic text and NUL-containing binary files;
- later in-range deletion of both files;
- branches, a merge and a rename;
- exact base, head, merge-base, ordered-history, parent-edge, blob and tracked
  tree bindings;
- deterministic raw and framed projection digests;
- a same-invocation Gitleaks coverage finding for every admitted projection.

The binary value is fabricated and is not a credential.
