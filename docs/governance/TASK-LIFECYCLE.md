# Task Lifecycle

```text
proposed
  -> owner selects exact task
  -> activation-only session creates and commits the activation bundle
  -> stop
  -> fresh session claims the activated task
  -> implement only allowed paths
  -> validate and independently review
  -> accept or fail closed
  -> close the task
  -> stop with every successor unselected
```

Only one task may be activated or implemented in a session. Activation never
means implementation, acceptance, publication, signing, deployment, or
successor selection.

An implementation session must begin from the exact activation commit and must
record:

- task ID and activation commit;
- allowed and forbidden paths;
- pre-change Git status;
- exact dependency/action/tool pins introduced;
- tests and adversarial checks run;
- exact evidence paths;
- post-change Git status;
- remaining owner gates;
- successor state as unselected.

If evidence is missing, inconsistent, or cannot prove the task's complete
boundary, the task is not accepted.

An authority-transition task additionally preserves every predecessor contract
and historical record byte-for-byte, assigns the successor a distinct identity,
maps every predecessor requirement explicitly and changes living authority only
through independently accepted closeout. A rejected implementation is never
rehabilitated by a governance transition.
