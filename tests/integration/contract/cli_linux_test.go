//go:build linux

package contract_test

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

func TestLinuxRequestFileDescriptor(t *testing.T) {
	requestPath, workspaceRoot, _ := validRequestFile(t)
	file, err := os.Open(requestPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	command := exec.Command(runnerPath, "--request-fd", "3", "--workspace-root", workspaceRoot)
	command.ExtraFiles = []*os.File{file}
	var stdout bytes.Buffer
	command.Stdout = &stdout
	err = command.Run()
	var exit *exec.ExitError
	if !errors.As(err, &exit) || exit.ExitCode() != 30 || decodeOne(t, stdout.Bytes())["reasonCode"] != string(outcome.ReasonUnavailableEngine) {
		t.Fatalf("fd invocation failed: %v %s", err, stdout.Bytes())
	}
}
