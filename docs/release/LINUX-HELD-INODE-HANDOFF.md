# Linux held-inode and namespace startup correction

Status: local candidate, not Linux runtime acceptance or a release build.

## Evidence and scope

Run35248154638 at f429621 successfully installed and selected the named profile,
retained ordinary credentials/context, and removed its owned policy. Its first
native case failed with exec EACCES for `/proc/2214/fd/169`. The journal was
unreadable. EACCES is observed; the precise denying kernel check was not traced.

The source-supported explanation is that `/proc/<parent>/fd` dereferencing is
ptrace-gated and the caller inside a new user namespace lacks ptrace capability
in the parent's namespace. Separately, namespace PID1 cannot use self-SIGSTOP
as the documented ancestor-origin forced stop. The old host `/proc` view also
does not correspond to PIDs returned inside the new PID namespace. These latter
issues are source findings, not additional observed run outcomes.

## Phase ordering

1. Preserve existing root-owned/non-writable Docker, setsid, unshare and shell
   admission. The only accepted target reference is the current caller's held
   `/proc/<current-pid>/fd/<number>`. Record its device/inode; never reopen an
   executable pathname as a fallback.
2. Start fixed setsid and a fixed shell with private redirected stdin and the
   original independently bounded stdout/stderr captures. The shell blocks on
   an initial token before opening the admitted descriptor or invoking unshare.
   Pin the supervisor using a pidfd, revalidate start time/parent/session and
   require it alive before releasing that first token.
3. The shell opens only the held reference as inherited descriptor3 and execs
   fixed unshare with same-user mapping, PID namespace, fork and kill-child.
   `--mount-proc --propagation private` adds a private mount namespace so /proc
   correctly describes the PID namespace. It does not remount the host's proc.
4. The init shell blocks on a second stdin token. Before releasing it, require
   the sole supervisor child, new user/PID/mount namespaces, stable parent,
   session and start time. Pin its identity, send SIGSTOP from the ancestor,
   observe state T, and compare inherited descriptor3's device/inode against
   the original held descriptor. All signals use pidfds, never numeric PID or
   process-group signaling after an identity check.
5. Only after complete admission, write the fixed second token while init is
   stopped, close owned stdin, and send SIGCONT through its pidfd. The init
   shell execs `/proc/self/fd/3` with the original argument-list boundaries.
   Malformed token or EOF before either token exits125 without target execution.
6. Require root exit, both completed streams, init death, absent init namespace
   path and empty outer session before returning trusted evidence. Original
   limits remain 131072 bytes independently, 15000ms complete process budget
   and 2000ms cleanup. No environment, operation or Docker daemon rule changes.

The public product entrypoint still accepts only its fixed Docker operations;
this introduces no general command runner or callback. The exclusive embedded
native source remains freshly compiled, with no ambient type fallback.

## Failure and cleanup matrix

| Failure | Required disposition |
| --- | --- |
| Unsupported host architecture/kernel pidfd API | Explicit failure, no PID-signaling fallback |
| Supervisor pin/identity failure | Close stdin; no initial token and no target execution |
| Open/map/mount/fork failure | Terminal denied startup; capture actual stderr and preserve missing namespace proof as uncertainty |
| Child disappearance, reuse, wrong parent/session/namespace | Refuse second token; close stdin and bounded cleanup |
| Stop, fd identity, token write or resume failure | Terminal; init remains the first pinned cleanup target |
| Timeout, overflow, stream failure, unexpected observation error | Close stdin; kill pinned init before supervisor, then pinned observed session members |
| Signal/poll/inspection failure or incomplete closure | Remain terminal; append cleanup uncertainty |
| Success or terminal return | Close owned stdin and dispose every acquired pidfd and process object |

pidfd acquisition compares process identity before/after opening and rejects
exited pins. ESRCH during signaling means that exact pinned process has exited;
it never redirects a signal to a reused PID. Other errno values fail. Poll is
nonblocking and accepts only termination/readable or reaped/hangup events.
The native bindings explicitly target the contract's Linux amd64 host matrix;
other architectures or missing syscall support fail instead of guessing an ABI.

## Fixture correction and preserved acceptance coverage

Release plan52-57 requires init death to kernel-terminate detached/nested
descendants and complete closure before trusted evidence. It does not require
a timeout when ordinary init exit already produces complete proven closure.
No additional authoritative terminal-on-normal-init-exit requirement was found
in the release/decision records searched for this correction.

Normal child/grandchild exit cases now require successful zero root exit and
complete closure. Added hold-child/hold-grandchild cases preserve explicit
timeout termination coverage for both topologies; the detached-session case,
hang, overflow, invalid UTF8 and cleanup-uncertainty rejection remain mandatory.
Atomic bounded markers prove that intended descendants actually started and
became ready inside the admitted namespace. Missing markers cannot vacuously
pass. Detached readiness must prove a new session. Namespace-local IDs are
never treated as host IDs. Surviving exact host ledger identities fail; death of
the pinned namespace init establishes kernel termination of all its namespace
descendants, including detached sessions. Policy removal alone proves none of
this. Windows fixture behavior is unchanged.

## Validation limits and sources

Local PowerShell parsing/native compilation, state/argument rejection tests and
the actual Windows CleanNative matrix passed during development. One temporary
UTF8 editing error was found by the parser and corrected before runtime tests;
the original euro-byte expectations and Windows branch remain intact.
There is no ordinary local Linux distribution; only stopped Docker Desktop WSL
was listed and left untouched. Actual Linux shell/unshare fd inheritance,
pidfd signaling, private proc mount, every fixture case and lifecycle closure
remain unproved until one separately reviewed invocation. Inert tests are not
substitutes for that proof. The named-profile driver separately records the
outer caller's mount namespace identity before/after the fixture.

Primary references (kernel source is explanatory, not an attestation of the
runner's exact Ubuntu kernel binary):

- https://man7.org/linux/man-pages/man5/proc_pid_fd.5.html
- https://raw.githubusercontent.com/torvalds/linux/v6.17/security/commoncap.c
- https://man7.org/linux/man-pages/man7/pid_namespaces.7.html
- https://raw.githubusercontent.com/torvalds/linux/v6.17/kernel/signal.c
- https://raw.githubusercontent.com/util-linux/util-linux/v2.39.3/sys-utils/unshare.c
- https://man7.org/linux/man-pages/man2/pidfd_open.2.html
- https://man7.org/linux/man-pages/man2/pidfd_send_signal.2.html
- https://raw.githubusercontent.com/torvalds/linux/v6.17/arch/x86/entry/syscalls/syscall_64.tbl
