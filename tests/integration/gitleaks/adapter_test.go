package gitleaks_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine/gitleaks"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
)

func TestAdapterClassifiesPrivateOutputAndPreservesArguments(t *testing.T) {
	executable, executableDigest := buildFakeGitleaks(t)
	config := pinnedConfig(t)
	configDigest := digestFile(t, config)
	ignoreFile, privateHome := supportFiles(t)
	target := filepath.Join(t.TempDir(), "candidate;$(not-a-command) --help")
	if err := os.Mkdir(target, 0o700); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		mode   string
		reason outcome.ReasonCode
	}{
		{"clean", outcome.ReasonPassNoBlockingFindings},
		{"finding", outcome.ReasonFailFindingDetected},
		{"unredacted", outcome.ReasonIndeterminateRedactionUnproven},
		{"malformed", outcome.ReasonIndeterminateSchemaUnsupported},
		{"stderr", outcome.ReasonIndeterminateIncompleteCoverage},
	}
	for _, test := range tests {
		t.Run(test.mode, func(t *testing.T) {
			env := append(os.Environ(), "PSCAN_FAKE_MODE="+test.mode, "PSCAN_EXPECTED_TARGET="+target)
			adapter := gitleaks.Adapter{Binding: gitleaks.Binding{
				Executable: executable, ExecutableDigest: executableDigest,
				Config: config, ConfigDigest: configDigest,
				IgnoreFile: ignoreFile, IgnoreFileDigest: digestFile(t, ignoreFile), PrivateHome: privateHome, Environment: env,
			}}
			result := adapter.ScanDirectory(context.Background(), target, 2*time.Second, 8<<20)
			if result.Reason != test.reason {
				t.Fatalf("got %s, want %s", result.Reason, test.reason)
			}
		})
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(target), "not-a-command")); !os.IsNotExist(err) {
		t.Fatal("candidate target became executable syntax")
	}
}

func TestBindingAndRuntimeVersionMismatchFailBeforeScan(t *testing.T) {
	executable, executableDigest := buildFakeGitleaks(t)
	config := pinnedConfig(t)
	target := t.TempDir()
	ignoreFile, privateHome := supportFiles(t)
	adapter := gitleaks.Adapter{Binding: gitleaks.Binding{
		Executable: executable, ExecutableDigest: executableDigest,
		Config: config, ConfigDigest: digestFile(t, config),
		IgnoreFile: ignoreFile, IgnoreFileDigest: digestFile(t, ignoreFile), PrivateHome: privateHome,
		Environment: append(os.Environ(), "PSCAN_FAKE_MODE=wrong-version"),
	}}
	if result := adapter.ScanDirectory(context.Background(), target, time.Second, 1<<20); result.Reason != outcome.ReasonFailBindingMismatch {
		t.Fatalf("wrong runtime version did not fail binding: %#v", result)
	}
	adapter.Binding.ExecutableDigest = strings.Repeat("0", 64)
	if result := adapter.ScanDirectory(context.Background(), target, time.Second, 1<<20); result.Reason != outcome.ReasonFailScannerIntegrity {
		t.Fatalf("wrong executable digest did not fail before scan: %#v", result)
	}
}

func TestParentDeadlineIsExplicitTimeout(t *testing.T) {
	executable, executableDigest := buildFakeGitleaks(t)
	config := pinnedConfig(t)
	ignoreFile, privateHome := supportFiles(t)
	adapter := gitleaks.Adapter{Binding: gitleaks.Binding{
		Executable: executable, ExecutableDigest: executableDigest,
		Config: config, ConfigDigest: digestFile(t, config),
		IgnoreFile: ignoreFile, IgnoreFileDigest: digestFile(t, ignoreFile), PrivateHome: privateHome,
		Environment: append(os.Environ(), "PSCAN_FAKE_MODE=timeout"),
	}}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if result := adapter.ScanDirectory(ctx, t.TempDir(), 2*time.Second, 1<<20); result.Reason != outcome.ReasonIndeterminateTimeout {
		t.Fatalf("timeout was not explicit: %#v", result)
	}
}

func supportFiles(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	ignore := filepath.Join(dir, "empty.gitleaksignore")
	home := filepath.Join(dir, "home")
	if err := os.WriteFile(ignore, []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(home, 0o700); err != nil {
		t.Fatal(err)
	}
	return ignore, home
}

func pinnedConfig(t *testing.T) string {
	t.Helper()
	path, err := filepath.Abs(filepath.Join("..", "..", "..", "rules", "generic", "gitleaks-v8.30.1.toml"))
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func buildFakeGitleaks(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	source := filepath.Join(dir, "main.go")
	program := `package main
import("encoding/json";"fmt";"os";"strings";"time")
func main(){
 mode:=os.Getenv("PSCAN_FAKE_MODE")
 if len(os.Args)>1 && os.Args[1]=="version" { if mode=="wrong-version" {fmt.Print("8.29.0")} else {fmt.Print("8.30.1")}; return }
 if len(os.Args)==0 || os.Args[len(os.Args)-1]!=os.Getenv("PSCAN_EXPECTED_TARGET") && os.Getenv("PSCAN_EXPECTED_TARGET")!="" {fmt.Fprint(os.Stderr,"argument mismatch");os.Exit(2)}
	 if strings.HasPrefix(mode,"projection-") {
  paths:=strings.Split(os.Getenv("PSCAN_COVERAGE_PATHS"),"|")
  if mode=="projection-missing" && len(paths)>0 {paths=paths[:len(paths)-1]}
  findings:=[]map[string]string{}
  for _,p:=range paths {findings=append(findings,map[string]string{"Secret":"REDACTED","RuleID":"pscan-projection-coverage","File":p})}
  b,_:=json.Marshal(findings);fmt.Print(string(b));if len(findings)>0 {os.Exit(11)};return
 }
 switch mode {
 case "clean": fmt.Print("[]")
 case "projection-clean","projection-missing": fmt.Print("[]")
 case "finding": b,_:=json.Marshal([]map[string]string{{"Secret":"REDACTED","RuleID":"synthetic"}});fmt.Print(string(b));os.Exit(11)
 case "unredacted": b,_:=json.Marshal([]map[string]string{{"Secret":"PSCAN_SYNTHETIC_CANARY_NOT_A_CREDENTIAL"}});fmt.Print(string(b));os.Exit(11)
 case "malformed": fmt.Print("{not-json")
 case "stderr": fmt.Fprint(os.Stderr,"private synthetic diagnostic")
 case "timeout": time.Sleep(5*time.Second);fmt.Print("[]")
 default: fmt.Fprint(os.Stderr,"unknown mode");os.Exit(2)
 }}
`
	if err := os.WriteFile(source, []byte(program), 0o600); err != nil {
		t.Fatal(err)
	}
	name := "fake-gitleaks"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	executable := filepath.Join(dir, name)
	goExe := filepath.Join(runtime.GOROOT(), "bin", "go")
	if runtime.GOOS == "windows" {
		goExe += ".exe"
	}
	cmd := exec.Command(goExe, "build", "-trimpath", "-buildvcs=false", "-o", executable, source)
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("cannot build synthetic engine: %v (%s)", err, output)
	}
	return executable, digestFile(t, executable)
}

func digestFile(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
