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

type Result struct {
	Reason   outcome.ReasonCode
	ExitCode int
}

// RunPrivate verifies the executable immediately before launch, uses no shell,
// captures output in bounded private memory, and discards it before returning.
func RunPrivate(parent context.Context, command Command, decode Decoder) Result {
	if err := VerifyRegularFile(command.Executable, command.ExpectedDigest); err != nil {
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
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	err := cmd.Run()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return Result{Reason: outcome.ReasonIndeterminateTimeout, ExitCode: -1}
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
	reason := decode(PrivateOutput{
		Stdout: append([]byte(nil), stdout.Bytes()...),
		Stderr: append([]byte(nil), stderr.Bytes()...),
		Exit:   exitCode,
	})
	stdout.Reset()
	stderr.Reset()
	return Result{Reason: reason, ExitCode: exitCode}
}

func VerifyRegularFile(path, expectedDigest string) error {
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
	if _, err := io.Copy(h, f); err != nil {
		return errors.New("bound file cannot be read")
	}
	actual := hex.EncodeToString(h.Sum(nil))
	if !constantStringEqual(actual, expectedDigest) {
		return fmt.Errorf("bound file digest mismatch")
	}
	return nil
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
