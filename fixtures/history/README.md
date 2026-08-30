# Generated History Fixtures

History repositories are generated only in private test directories so no
`.git` authority, object database, hook or repository identity is committed.
The generator in `tests/unit/gitinput/git_test.go` creates:

- a clean range base;
- an in-range synthetic binary file with a NUL marker;
- a later in-range deletion of that file;
- exact base, head, merge-base, ordered-history and tracked-tree bindings.

The binary value is fabricated and is not a credential.
