# Adversarial Git Fixture Matrix

PSCAN-03 tests generate hostile local Git inputs with candidate-controlled
metacharacters in paths, a fresh empty clone template, disabled inherited
configuration, hooks, attributes, pagers, external diffs, prompts and unsafe
protocols, reversed ancestry, and binary add/delete history.

No checked-in hook is executable and no fixture performs network access.
