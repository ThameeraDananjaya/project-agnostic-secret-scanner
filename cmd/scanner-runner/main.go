package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/request"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/workspace"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, time.Now)) }

func run(args []string, stdout, stderr io.Writer, now func() time.Time) int {
	started := now().UTC()
	safeID, idErr := request.NewScanID()
	if idErr != nil {
		safeID = "00000000-0000-4000-8000-000000000000"
	}
	mode := "unknown"
	attempt := 1
	reason := outcome.ReasonFailInputIntegrity
	var bindings *outcome.Bindings
	var supersedes string

	flags := flag.NewFlagSet("scanner-runner", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	requestPath := flags.String("request", "", "")
	requestFD := flags.Int("request-fd", -1, "")
	workspaceRoot := flags.String("workspace-root", workspace.DefaultRoot(), "")
	flags.IntVar(&attempt, "attempt", 1, "")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 || attempt < 1 || attempt > 3 || ((*requestPath == "") == (*requestFD < 0)) {
		writeDiagnostic(stderr, "invalid invocation")
		return emit(stdout, safeID, mode, supersedes, reason, attempt, started, now().UTC(), bindings)
	}

	reader, closeReader, err := openRequest(*requestPath, *requestFD)
	if err != nil {
		writeDiagnostic(stderr, "required input unavailable")
		return emit(stdout, safeID, mode, supersedes, outcome.ReasonUnavailableRequiredInput, attempt, started, now().UTC(), bindings)
	}
	defer closeReader()
	value, err := request.Load(reader)
	if err != nil {
		reason = request.ReasonFor(err, outcome.ReasonFailInputIntegrity)
		writeDiagnostic(stderr, "request rejected")
		return emit(stdout, safeID, mode, supersedes, reason, attempt, started, now().UTC(), bindings)
	}
	safeID, mode, supersedes = value.ScanID, value.Mode, value.SupersedesScanID
	if err := request.ValidateAt(value, now().UTC()); err != nil {
		reason = request.ReasonFor(err, outcome.ReasonFailInputIntegrity)
		writeDiagnostic(stderr, "request rejected")
		return emit(stdout, safeID, mode, supersedes, reason, attempt, started, now().UTC(), bindings)
	}
	if err := request.ValidateBoundFiles(value); err != nil {
		reason = request.ReasonFor(err, outcome.ReasonFailInputIntegrity)
		writeDiagnostic(stderr, "request bindings rejected")
		return emit(stdout, safeID, mode, supersedes, reason, attempt, started, now().UTC(), bindings)
	}
	bindings = &outcome.Bindings{
		ScannerReleaseDigest:        value.ScannerReleaseDigest,
		EngineName:                  value.EngineBinding.Name,
		EngineVersion:               value.EngineBinding.Version,
		EngineBinaryDigest:          value.EngineBinding.BinaryDigest,
		AdapterVersion:              value.EngineBinding.AdapterVersion,
		RulePackDigest:              value.RulePackDigest,
		PolicyDigest:                value.PolicyDigest,
		AllowlistDigest:             value.AllowlistDigest,
		RequestSchemaVersion:        value.RequestSchemaVersion,
		SourceBaseCommit:            value.SourceBinding.BaseCommit,
		SourceHeadCommit:            value.SourceBinding.HeadCommit,
		SourceMergeBase:             value.SourceBinding.MergeBase,
		HistoryRangeDigest:          value.SourceBinding.HistoryRangeDigest,
		TrackedTreeDigest:           value.SourceBinding.TrackedTreeDigest,
		TrackedSourceManifestDigest: value.TrackedSourceManifest.Digest,
	}
	if value.BuildContextManifest != nil {
		bindings.BuildContextManifestDigest = value.BuildContextManifest.Digest
	}
	if value.ArtifactManifest != nil {
		bindings.ArtifactManifestDigest = value.ArtifactManifest.Digest
	}
	manager, err := workspace.NewManager(*workspaceRoot)
	if err != nil {
		writeDiagnostic(stderr, "workspace unavailable")
		return emit(stdout, safeID, mode, supersedes, outcome.ReasonUnavailableWorkspace, attempt, started, now().UTC(), bindings)
	}
	space, err := manager.Create(value.ScanID, attempt)
	if err != nil {
		writeDiagnostic(stderr, "workspace unavailable")
		return emit(stdout, safeID, mode, supersedes, outcome.ReasonUnavailableWorkspace, attempt, started, now().UTC(), bindings)
	}
	if err := space.Cleanup(); err != nil {
		writeDiagnostic(stderr, "workspace invariant failure")
		return emit(stdout, safeID, mode, supersedes, outcome.ReasonTerminalInternalInvariant, attempt, started, now().UTC(), bindings)
	}
	writeDiagnostic(stderr, "engine unavailable")
	return emit(stdout, safeID, mode, supersedes, outcome.ReasonUnavailableEngine, attempt, started, now().UTC(), bindings)
}

func openRequest(path string, fd int) (io.Reader, func(), error) {
	if path != "" {
		if err := request.ValidateLocalPath(path); err != nil {
			return nil, func() {}, err
		}
		file, err := os.Open(path)
		if err != nil {
			return nil, func() {}, err
		}
		return file, func() { _ = file.Close() }, nil
	}
	if fd < 3 {
		return nil, func() {}, fmt.Errorf("unsafe descriptor")
	}
	file := os.NewFile(uintptr(fd), "request")
	if file == nil {
		return nil, func() {}, fmt.Errorf("descriptor unavailable")
	}
	return file, func() { _ = file.Close() }, nil
}

func emit(w io.Writer, scanID, mode, supersedes string, reason outcome.ReasonCode, attempt int, start, end time.Time, bindings *outcome.Bindings) int {
	value, err := outcome.NewTerminal(scanID, mode, reason, attempt, start, end)
	if err != nil {
		value, _ = outcome.NewTerminal("00000000-0000-4000-8000-000000000000", "unknown", outcome.ReasonTerminalInternalInvariant, 1, start, end)
	}
	value.Bindings = bindings
	value.SupersedesScanID = supersedes
	if err := outcome.Serialize(w, value); err != nil {
		return 40
	}
	code, ok := outcome.ExitCode(value.State)
	if !ok {
		return 40
	}
	return code
}

func writeDiagnostic(w io.Writer, message string) {
	_, _ = io.WriteString(w, "scanner-runner: "+message+"\n")
}
