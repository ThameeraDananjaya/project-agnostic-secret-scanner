# PSCAN-06 Correction C2 iteration 004 containment refinement

Candidate `6f791646413bfde52a7f034f6219d92f6fb44c03` and its Windows author
validation remain factual for the exact committed bytes tested. A subsequent
post-integration design audit found that its Linux private-session ledger
sampled `/proc` every two milliseconds. Although the ledger used PID plus
start-time identity, a child could create a new session between samples. That
left a theoretical detached-descendant membership race and does not satisfy the
iteration-004 no-race requirement. Candidate `6f79164` is therefore not the
final iteration-004 readiness candidate.

The bounded refinement retains the private outer session but adds a new Linux
user/PID namespace using fixed root-owned system `unshare`. A fixed root-owned
shell becomes the namespace init and stops itself before Docker executes. The
parent records the init PID and kernel namespace identity, verifies the stopped
pre-execution state, then resumes it. The init executes the already held Docker
inode. A process cannot detach from its containing PID namespace; when the
namespace init is terminated or exits, the kernel terminates every remaining
member including nested namespaces. The boundary does not return until the
outer session is empty, the init `/proc` identity is absent, both streams close
and the supervisor exits.

The Linux detached fixture now remains alive until the command timeout so it
must be terminated by the namespace boundary. The actual-Linux fixture remains
unexecuted on the Windows author host and requires separately authorized Linux
execution before acceptance. Windows job-object code, the Docker operation
table, admission state model, arguments, payloads, schemas, tags, network
rules, stream limits and all remote/successor gates are unchanged.
