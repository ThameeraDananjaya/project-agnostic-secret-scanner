// Package gitinput prepares and binds Git inputs without executing candidate
// text as command syntax.
package gitinput

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine"
)

const outputLimit = 32 << 20

var oidPattern = regexp.MustCompile(`^(?:[0-9a-f]{40}|[0-9a-f]{64})$`)

type Git struct {
	Executable string
	Digest     string
	Timeout    time.Duration
}

type RangeBinding struct {
	Base               string
	Head               string
	MergeBase          string
	FirstRelease       bool
	HistoryRangeDigest string
	TrackedTreeDigest  string
	HeadTreeOID        string
	CommitCount        int
	Commits            []string
}

// PrepareBare creates a new, non-local bare clone with an empty template and a
// closed configuration boundary. Both source and destination must be absolute
// local paths; URL-like sources are rejected.
func (g Git) PrepareBare(ctx context.Context, source, destination, privateRoot string) error {
	if err := g.verify(); err != nil {
		return err
	}
	if !safeAbsolute(source) || !safeAbsolute(destination) || !safeAbsolute(privateRoot) || source == destination {
		return errors.New("unsafe Git path")
	}
	if strings.Contains(source, "://") || strings.HasPrefix(source, "\\\\") {
		return errors.New("remote or UNC Git source is forbidden")
	}
	if _, err := os.Stat(destination); !errors.Is(err, os.ErrNotExist) {
		return errors.New("destination must not exist")
	}
	template := filepath.Join(privateRoot, "empty-template")
	home := filepath.Join(privateRoot, "home")
	if err := os.MkdirAll(template, 0o700); err != nil {
		return errors.New("cannot create private Git template")
	}
	if err := os.MkdirAll(home, 0o700); err != nil {
		return errors.New("cannot create private Git home")
	}
	args := []string{"-c", "protocol.file.allow=always", "clone", "--quiet", "--bare", "--no-local", "--no-hardlinks", "--no-tags", "--template=" + template, "--", source, destination}
	if _, err := g.run(ctx, "", home, args...); err != nil {
		return errors.New("safe bare clone failed")
	}
	configArgs := [][]string{
		{"-C", destination, "config", "--local", "core.hooksPath", os.DevNull},
		{"-C", destination, "config", "--local", "core.attributesFile", os.DevNull},
		{"-C", destination, "config", "--local", "core.pager", "cat"},
		{"-C", destination, "config", "--local", "pager.log", "false"},
		{"-C", destination, "config", "--local", "diff.external", ""},
		{"-C", destination, "config", "--local", "protocol.file.allow", "never"},
	}
	for _, command := range configArgs {
		if _, err := g.run(ctx, "", home, command...); err != nil {
			return errors.New("cannot seal bare repository configuration")
		}
	}
	return nil
}

// BindRange resolves exact commits, proves ancestry/merge-base, and hashes the
// exact ordered commit list and tracked tree representation.
func (g Git) BindRange(ctx context.Context, repository, base, head string, firstRelease bool, privateHome string) (RangeBinding, error) {
	if err := g.verify(); err != nil {
		return RangeBinding{}, err
	}
	if !safeAbsolute(repository) || !safeAbsolute(privateHome) || !oidPattern.MatchString(head) {
		return RangeBinding{}, errors.New("invalid Git binding")
	}
	commits, mergeBase, err := g.enumerateRange(ctx, repository, base, head, firstRelease, privateHome)
	if err != nil {
		return RangeBinding{}, err
	}
	binding := RangeBinding{Base: base, Head: head, MergeBase: mergeBase, FirstRelease: firstRelease}
	binding.CommitCount = len(commits)
	binding.Commits = append([]string(nil), commits...)
	binding.HistoryRangeDigest = digestDomain("pscan.git-history.v1", []byte(strings.Join(commits, "\x00")))
	treeRaw, err := g.run(ctx, "", privateHome, "-C", repository, "ls-tree", "-r", "-z", "--full-tree", head)
	if err != nil || len(treeRaw) == 0 {
		return RangeBinding{}, errors.New("cannot bind tracked tree")
	}
	binding.TrackedTreeDigest = digestDomain("pscan.git-tree.v1", treeRaw)
	treeOID, err := g.run(ctx, "", privateHome, "-C", repository, "rev-parse", "--verify", head+"^{tree}")
	if err != nil || !oidPattern.MatchString(strings.TrimSpace(string(treeOID))) {
		return RangeBinding{}, errors.New("cannot bind head tree object")
	}
	binding.HeadTreeOID = strings.TrimSpace(string(treeOID))
	return binding, nil
}

func (g Git) enumerateRange(ctx context.Context, repository, base, head string, firstRelease bool, privateHome string) ([]string, string, error) {
	resolvedHead, err := g.resolveCommit(ctx, repository, head, privateHome)
	if err != nil || resolvedHead != head {
		return nil, "", errors.New("head commit mismatch")
	}
	var revListArgs []string
	mergeBase := ""
	if firstRelease {
		if base != "" {
			return nil, "", errors.New("first release must not bind a base")
		}
		revListArgs = []string{"-C", repository, "rev-list", "--reverse", "--topo-order", head}
	} else {
		if !oidPattern.MatchString(base) || base == head {
			return nil, "", errors.New("invalid base commit")
		}
		resolvedBase, resolveErr := g.resolveCommit(ctx, repository, base, privateHome)
		if resolveErr != nil || resolvedBase != base {
			return nil, "", errors.New("base commit mismatch")
		}
		mergeBaseRaw, runErr := g.run(ctx, "", privateHome, "-C", repository, "merge-base", "--", base, head)
		if runErr != nil {
			return nil, "", errors.New("cannot establish merge base")
		}
		mergeBase = strings.TrimSpace(string(mergeBaseRaw))
		if mergeBase != base {
			return nil, "", errors.New("base is not an ancestor of head")
		}
		revListArgs = []string{"-C", repository, "rev-list", "--reverse", "--topo-order", base + ".." + head}
	}
	commitsRaw, err := g.run(ctx, "", privateHome, revListArgs...)
	if err != nil {
		return nil, "", errors.New("cannot enumerate exact history range")
	}
	commits, err := parseCommitList(commitsRaw)
	if err != nil || len(commits) == 0 {
		return nil, "", errors.New("empty or invalid history range")
	}
	return commits, mergeBase, nil
}

func (g Git) resolveCommit(ctx context.Context, repository, oid, privateHome string) (string, error) {
	out, err := g.run(ctx, "", privateHome, "-C", repository, "rev-parse", "--verify", oid+"^{commit}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (g Git) verify() error {
	if g.Timeout <= 0 {
		return errors.New("invalid Git timeout")
	}
	return engine.VerifyRegularFile(g.Executable, g.Digest)
}

func (g Git) run(parent context.Context, directory, privateHome string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(parent, g.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, g.Executable, args...)
	cmd.Dir = directory
	cmd.Env = gitEnvironment(filepath.Dir(g.Executable), privateHome)
	cmd.Stdin = nil
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return nil, errors.New("Git command failed")
	}
	if ctx.Err() != nil || stdout.Len() > outputLimit || stderr.Len() != 0 {
		return nil, errors.New("Git command produced incomplete or unexpected output")
	}
	return append([]byte(nil), stdout.Bytes()...), nil
}

func gitEnvironment(pathValue, privateHome string) []string {
	env := engine.SafeEnvironment(pathValue, privateHome)
	env = append(env,
		"GIT_EXTERNAL_DIFF=",
		"GIT_DIFF_OPTS=--no-ext-diff",
		"GIT_SSH_COMMAND=",
		"GIT_ASKPASS=",
		"SSH_ASKPASS=",
		"GIT_CONFIG_COUNT=3",
		"GIT_CONFIG_KEY_0=core.hooksPath", "GIT_CONFIG_VALUE_0="+os.DevNull,
		"GIT_CONFIG_KEY_1=core.attributesFile", "GIT_CONFIG_VALUE_1="+os.DevNull,
		"GIT_CONFIG_KEY_2=diff.external", "GIT_CONFIG_VALUE_2=",
	)
	return env
}

func parseCommitList(raw []byte) ([]string, error) {
	lines := strings.Fields(string(raw))
	seen := make(map[string]bool, len(lines))
	for _, line := range lines {
		if !oidPattern.MatchString(line) || seen[line] {
			return nil, errors.New("invalid commit list")
		}
		seen[line] = true
	}
	return lines, nil
}

func digestDomain(domain string, payload []byte) string {
	h := sha256.New()
	h.Write([]byte(domain))
	h.Write([]byte{0})
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

func safeAbsolute(path string) bool {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path || strings.ContainsRune(path, '\x00') {
		return false
	}
	if runtime.GOOS == "windows" && strings.HasPrefix(path, "\\\\") {
		return false
	}
	return true
}
