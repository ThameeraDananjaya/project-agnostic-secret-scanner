# Evidence and Closeout Rules

## Required task evidence

Every implementation task records:

- task ID, activation commit and claim statement;
- activated specification digest or commit;
- allowed and forbidden paths/actions;
- pre-change branch/HEAD/status;
- exact files changed;
- exact toolchain, dependency, action, scanner and source pins introduced, or
  an explicit `none`;
- tests, adversarial checks and independent review performed;
- expected and actual outcomes, including limitations;
- evidence paths and relevant hashes;
- owner gates encountered or still reserved;
- post-commit HEAD/status;
- successor state as unselected.

Evidence is content-free. It never stores candidate source, findings, paths from
a consuming project, credentials, keys, receipt material or customer data.

For an authority transition, evidence also records predecessor canonical,
external-source and historical-record hashes; the explicit successor identity;
complete requirement dispositions; and agreement of every changed living
authority. Predecessor files remain immutable.

## Acceptance rules

Acceptance requires reproducible evidence for every task success criterion and
complete path/action containment. An implementation claim, passing subset,
workflow status, scanner result or clean Git status is insufficient alone.

Any missing, stale, unsupported, skipped, conflicting or untrusted condition
fails closed. Ordinary technical defects remain inside the same task until
corrected. A reserved owner gate stops work without inventing approval.

## Independent review

The reviewer reads the controlling contract and activated task, examines the
bounded diff rather than author summaries, reruns proportionate checks, and
adversarially tests security boundaries. The review records facts,
limitations, interpretations, recommendations and owner decisions separately.

## Closeout sequence

```text
implemented within activated bounds
    -> validation evidence complete
    -> independent read-only review
    -> accepted or failed closed
    -> path-bounded local commit
    -> clean worktree proof
    -> task marked completed
    -> stop with every successor unselected
```

No closeout activates, selects, claims or begins a successor. Remote push,
publication, release, signing, settings, spending, credentials and consuming-
project integration always require their own applicable authority.

## Recovery

Historical evidence is immutable. A correction creates a new bounded record
that references the earlier record; it does not rewrite an accepted or rejected
fact. Revocation is append-only and restoration requires new clean evidence.
Decommissioning a contract's current authority likewise preserves its bytes and
does not accept a rejected implementation produced under it.
