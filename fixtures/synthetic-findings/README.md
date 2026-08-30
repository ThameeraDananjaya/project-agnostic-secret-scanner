# Synthetic Gitleaks Fixtures

Every value in this fixture family is fabricated for scanner validation and is
not a credential. `clean.txt` is the clean control. `finding.txt` uses the
invalid-by-construction `PSCAN_SYNTHETIC_SECRET_` acceptance marker recognized
only by the pinned PSCAN fixture rule.
Tests generate unique copies in private temporary directories; raw engine
reports never enter repository evidence or public process output.

The checked-in values must never be used for provider verification.
