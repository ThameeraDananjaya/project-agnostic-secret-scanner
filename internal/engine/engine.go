// Package engine provides a fail-closed process boundary for scanner engines.
// Candidate-controlled output is consumed only by a private decoder and is
// never included in the returned result.
package engine

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

const DefaultCaptureLimit = 16 << 20

type Command struct {
	Executable     string
	ExpectedDigest string
	Args           []string
	Dir            string
	Environment    []string
	Timeout        time.Duration
	CaptureLimit   int64
}

type PrivateOutput struct {
	Stdout []byte
	Stderr []byte
	Exit   int
}

type Decoder func(PrivateOutput) outcome.ReasonCode
type ContextDecoder func(context.Context, PrivateOutput) outcome.ReasonCode

type Result struct {
	Reason   outcome.ReasonCode
	ExitCode int
}

// RunPrivate verifies the executable immediately before launch, uses no shell,
// captures output in bounded private memory, and discards it before returning.
func RunPrivate(parent context.Context, command Command, decode Decoder) (result Result) {
	if decode == nil {
		return Result{Reason: outcome.ReasonTerminalInternalInvariant, ExitCode: -1}
	}
	return RunPrivateContext(parent, command, func(_ context.Context, output PrivateOutput) outcome.ReasonCode {
		return decode(output)
	})
}

func RunPrivateContext(parent context.Context, command Command, decode ContextDecoder) (result Result) {
	result = Result{Reason: outcome.ReasonTerminalInternalInvariant, ExitCode: -1}
	defer func() {
		if recover() != nil {
			result = Result{Reason: outcome.ReasonTerminalInternalInvariant, ExitCode: -1}
		}
	}()
	if err := VerifyRegularFileContext(parent, command.Executable, command.ExpectedDigest); err != nil {
		if errors.Is(parent.Err(), context.DeadlineExceeded) {
			return Result{Reason: outcome.ReasonIndeterminateTimeout, ExitCode: -1}
		}
		if errors.Is(parent.Err(), context.Canceled) {
			return Result{Reason: outcome.ReasonIndeterminateIncompleteCoverage, ExitCode: -1}
		}
		return Result{Reason: outcome.ReasonFailScannerIntegrity, ExitCode: -1}
	}
	if decode == nil || !filepath.IsAbs(command.Executable) || command.Timeout <= 0 {
		return Result{Reason: outcome.ReasonTerminalInternalInvariant, ExitCode: -1}
	}
	limit := command.CaptureLimit
	if limit <= 0 {
		limit = DefaultCaptureLimit
	}
	ctx, cancel := context.WithTimeout(parent, command.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command.Executable, command.Args...)
	cmd.Dir = command.Dir
	cmd.Env = append([]string(nil), command.Environment...)
	cmd.Stdin = nil
	var stdout, stderr boundedBuffer
	stdout.limit, stderr.limit = limit, limit
	defer stdout.Reset()
	defer stderr.Reset()
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	err := cmd.Run()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return Result{Reason: outcome.ReasonIndeterminateTimeout, ExitCode: -1}
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return Result{Reason: outcome.ReasonIndeterminateIncompleteCoverage, ExitCode: -1}
	}
	if stdout.overflow || stderr.overflow {
		return Result{Reason: outcome.ReasonIndeterminateResourceLimit, ExitCode: -1}
	}
	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			return Result{Reason: outcome.ReasonUnavailableRuntime, ExitCode: -1}
		}
		exitCode = exitErr.ExitCode()
	}
	reason := decode(ctx, PrivateOutput{
		Stdout: append([]byte(nil), stdout.Bytes()...),
		Stderr: append([]byte(nil), stderr.Bytes()...),
		Exit:   exitCode,
	})
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return Result{Reason: outcome.ReasonIndeterminateTimeout, ExitCode: -1}
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return Result{Reason: outcome.ReasonIndeterminateIncompleteCoverage, ExitCode: -1}
	}
	return Result{Reason: reason, ExitCode: exitCode}
}

func VerifyRegularFile(path, expectedDigest string) error {
	return VerifyRegularFileContext(context.Background(), path, expectedDigest)
}

func VerifyRegularFileContext(ctx context.Context, path, expectedDigest string) error {
	if ctx == nil {
		return errors.New("invalid file binding")
	}
	if !filepath.IsAbs(path) || len(expectedDigest) != 64 || strings.ToLower(expectedDigest) != expectedDigest {
		return errors.New("invalid file binding")
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return errors.New("bound file is not a regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return errors.New("bound file cannot be opened")
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil || !os.SameFile(info, opened) {
		return errors.New("bound file identity changed")
	}
	h := sha256.New()
	if _, err := io.Copy(h, engineContextReader{ctx: ctx, reader: f}); err != nil {
		return errors.New("bound file cannot be read")
	}
	actual := hex.EncodeToString(h.Sum(nil))
	if !constantStringEqual(actual, expectedDigest) {
		return fmt.Errorf("bound file digest mismatch")
	}
	return nil
}

type engineContextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r engineContextReader) Read(value []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(value)
}

func SafeEnvironment(pathValue, privateHome string) []string {
	env := []string{
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_TERMINAL_PROMPT=0",
		"GIT_PAGER=cat",
		"GIT_ATTR_NOSYSTEM=1",
		"PAGER=cat",
		"PATH=" + pathValue,
		"HOME=" + privateHome,
		"XDG_CONFIG_HOME=" + filepath.Join(privateHome, ".config"),
		"GIT_CONFIG_GLOBAL=" + os.DevNull,
	}
	if runtime.GOOS == "windows" {
		env = append(env, "SYSTEMROOT="+os.Getenv("SYSTEMROOT"), "COMSPEC="+os.Getenv("COMSPEC"))
	}
	return env
}

// IsSealedEnvironment proves that authoritative execution receives only the
// scanner's explicit non-secret process variables. It never inherits the host
// environment wholesale.
func IsSealedEnvironment(environment []string) bool {
	allowed := map[string]bool{
		"GIT_CONFIG_NOSYSTEM": true, "GIT_TERMINAL_PROMPT": true,
		"GIT_PAGER": true, "GIT_ATTR_NOSYSTEM": true, "PAGER": true,
		"PATH": true, "HOME": true, "XDG_CONFIG_HOME": true,
		"GIT_CONFIG_GLOBAL": true, "SYSTEMROOT": true, "COMSPEC": true,
		"GOMEMLIMIT": true, "GOMAXPROCS": true,
	}
	seen := map[string]bool{}
	for _, item := range environment {
		key, value, ok := strings.Cut(item, "=")
		key = strings.ToUpper(key)
		if !ok || key == "" || !allowed[key] || seen[key] || strings.ContainsAny(value, "\r\n\x00") {
			return false
		}
		seen[key] = true
	}
	for _, required := range []string{"PATH", "HOME", "XDG_CONFIG_HOME", "GIT_CONFIG_GLOBAL", "GIT_CONFIG_NOSYSTEM", "GIT_TERMINAL_PROMPT"} {
		if !seen[required] {
			return false
		}
	}
	return true
}

func constantStringEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	var diff byte
	for i := range a {
		diff |= a[i] ^ b[i]
	}
	return diff == 0
}

type boundedBuffer struct {
	buf      bytes.Buffer
	limit    int64
	overflow bool
}

func (b *boundedBuffer) Bytes() []byte { return b.buf.Bytes() }
func (b *boundedBuffer) Len() int      { return b.buf.Len() }
func (b *boundedBuffer) Reset()        { b.buf.Reset(); b.overflow = false }

func (b *boundedBuffer) Write(p []byte) (int, error) {
	remaining := b.limit - int64(b.Len())
	if remaining <= 0 {
		b.overflow = true
		return len(p), nil
	}
	write := p
	if int64(len(write)) > remaining {
		write = write[:remaining]
		b.overflow = true
	}
	_, _ = b.buf.Write(write)
	return len(p), nil
}
