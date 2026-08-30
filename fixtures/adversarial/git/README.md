# Adversarial Git Fixture Matrix

PSCAN-03 tests generate hostile local Git inputs with candidate-controlled
metacharacters in paths, a fresh empty clone template, disabled inherited
configuration, hooks, attributes, pagers, external diffs, prompts and unsafe
protocols, reversed ancestry, binary add/delete history, merges, renames,
portable-path violations, case collisions, symlinks, gitlinks, per-file/count/
aggregate limits, missing files, unbound files and byte mutations.

No checked-in hook is executable and no fixture performs network access.
