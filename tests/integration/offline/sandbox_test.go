package offline_test

import (
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestAuthoritativeSandboxHasNoNetworkCredentialsOrPrivilege(t *testing.T) {
	if os.Getenv("PSCAN_AUTHORITATIVE_DOCKER") != "1" {
		t.Skip("authoritative Docker-only test")
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}
	for _, iface := range interfaces {
		if iface.Name != "lo" {
			t.Fatalf("unexpected network interface: %s", iface.Name)
		}
	}
	if _, err := net.LookupHost("example.com"); err == nil {
		t.Fatal("DNS canary connected")
	}
	for _, address := range []string{"1.1.1.1:53", "169.254.169.254:80", "100.100.100.200:80"} {
		connection, err := net.DialTimeout("tcp", address, 250*time.Millisecond)
		if err == nil {
			_ = connection.Close()
			t.Fatalf("TCP/provider canary connected")
		}
	}
	for _, item := range os.Environ() {
		key := strings.ToUpper(strings.SplitN(item, "=", 2)[0])
		for _, forbidden := range []string{"TOKEN", "SECRET", "PASSWORD", "CREDENTIAL", "AWS_", "AZURE_", "GOOGLE_", "GITHUB_"} {
			if strings.Contains(key, forbidden) {
				t.Fatalf("inherited credential-like environment key: %s", key)
			}
		}
	}
	status, err := os.ReadFile("/proc/self/status")
	if err != nil || !strings.Contains(string(status), "CapEff:\t0000000000000000") || !strings.Contains(string(status), "NoNewPrivs:\t1") {
		t.Fatalf("privilege controls absent: %v", err)
	}
	if err := os.WriteFile("/pscan-root-write", []byte("x"), 0o600); err == nil {
		_ = os.Remove("/pscan-root-write")
		t.Fatal("root filesystem is writable")
	}
	assertCgroupBound(t, "/sys/fs/cgroup/pids.max", 256)
	memory, err := strconv.ParseInt(os.Getenv("PSCAN_SANDBOX_MEMORY_BYTES"), 10, 64)
	if err != nil || memory <= 0 {
		t.Fatal("declared sandbox memory is absent")
	}
	assertCgroupExact(t, "/sys/fs/cgroup/memory.max", memory)
	assertCgroupExact(t, "/sys/fs/cgroup/memory.swap.max", 0)
	if os.Getenv("PSCAN_SANDBOX_PROCESS_LIMIT") != "1" {
		t.Fatal("scanner subprocess limit is not declared as one")
	}
	if info, err := os.Stat("/src"); err != nil || !info.IsDir() {
		t.Fatal("read-only source mount absent")
	}
	probe := filepath.Join("/src", ".pscan-write-probe")
	if err := os.WriteFile(probe, []byte("x"), 0o600); err == nil {
		_ = os.Remove(probe)
		t.Fatal("source mount is writable")
	}
}

func assertCgroupExact(t *testing.T, path string, expected int64) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	value := strings.TrimSpace(string(raw))
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed != expected {
		t.Fatalf("%s=%s want %d", path, value, expected)
	}
}

func assertCgroupBound(t *testing.T, path string, maximum int64) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	value := strings.TrimSpace(string(raw))
	if value == "max" {
		t.Fatalf("%s is unbounded", path)
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed > maximum {
		t.Fatalf("%s=%s", path, value)
	}
}
