package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"
)

func fixtureInfo() (*debug.BuildInfo, component) {
	c := components()[0]
	bi := &debug.BuildInfo{GoVersion: c.Go, Path: c.MainPath, Main: debug.Module{Path: c.Module, Version: "(devel)"}}
	for k, v := range map[string]string{"GOOS": c.OS, "GOARCH": "amd64", "GOAMD64": "v1", "CGO_ENABLED": "0", "-trimpath": "true", "-compiler": "gc", "-buildmode": "exe"} {
		bi.Settings = append(bi.Settings, debug.BuildSetting{Key: k, Value: v})
	}
	return bi, c
}

func TestDistinctCompiledVersionsRetainTheirOwnSumsAndLicences(t *testing.T) {
	path := "example.invalid/module"
	sum := "h1:Ra4+bf83h2ztPIQYNP99R6m+Y7KfnARDfID+a+vLl4s="
	first := &debug.Module{Path: path, Version: "v1.0.0", Sum: sum}
	second := &debug.Module{Path: path, Version: "v2.0.0+incompatible", Sum: sum}
	sums := map[string]string{path + "@" + first.Version: sum, path + "@" + second.Version: sum}
	licences := map[string][]licenceBinding{path + "@" + first.Version: {{"old/LICENSE", strings.Repeat("a", 64)}}, path + "@" + second.Version: {{"new/LICENSE", strings.Repeat("b", 64)}}}
	a, err := dependencyPackage(first, sums, licences)
	if err != nil {
		t.Fatal(err)
	}
	b, err := dependencyPackage(second, sums, licences)
	if err != nil {
		t.Fatal(err)
	}
	if a.ID == b.ID || a.Version != first.Version || b.Version != second.Version || a.Source == b.Source {
		t.Fatal("versions/licences collapsed")
	}
	delete(licences, path+"@"+second.Version)
	if _, err := dependencyPackage(second, sums, licences); err == nil || !strings.Contains(err.Error(), "no shipped licence") {
		t.Fatal("another version's licence substituted")
	}
	delete(sums, path+"@"+first.Version)
	if _, err := dependencyPackage(first, sums, licences); err == nil || !strings.Contains(err.Error(), "sum absent") {
		t.Fatal("missing source sum accepted")
	}
}
func TestMetadataRejectsIncompleteReplacedAndMixedComponents(t *testing.T) {
	mutations := map[string]func(*debug.BuildInfo){
		"toolchain":         func(b *debug.BuildInfo) { b.GoVersion = "go1.27.0" },
		"wrong root":        func(b *debug.BuildInfo) { b.Main.Path = engineModule },
		"wrong command":     func(b *debug.BuildInfo) { b.Path = engineModule },
		"root replacement":  func(b *debug.BuildInfo) { b.Main.Replace = &debug.Module{Path: "example.invalid/replaced"} },
		"missing setting":   func(b *debug.BuildInfo) { b.Settings = nil },
		"duplicate setting": func(b *debug.BuildInfo) { b.Settings = append(b.Settings, b.Settings[0]) },
		"vcs ambiguity": func(b *debug.BuildInfo) {
			b.Settings = append(b.Settings, debug.BuildSetting{Key: "vcs.modified", Value: "true"})
		},
		"missing sum": func(b *debug.BuildInfo) {
			b.Deps = []*debug.Module{{Path: "example.invalid/module", Version: "v1.0.0"}}
		},
		"replaced dep": func(b *debug.BuildInfo) {
			b.Deps = []*debug.Module{{Path: "example.invalid/module", Version: "v1.0.0", Sum: "h1:Ra4+bf83h2ztPIQYNP99R6m+Y7KfnARDfID+a+vLl4s=", Replace: &debug.Module{Path: "local"}}}
		},
		"duplicate dep": func(b *debug.BuildInfo) {
			m := &debug.Module{Path: "example.invalid/module", Version: "v1.0.0", Sum: "h1:Ra4+bf83h2ztPIQYNP99R6m+Y7KfnARDfID+a+vLl4s="}
			b.Deps = []*debug.Module{m, m}
		},
		"nil dep": func(b *debug.BuildInfo) { b.Deps = []*debug.Module{nil} },
	}
	good, c := fixtureInfo()
	if err := validateInfo(good, c); err != nil {
		t.Fatal(err)
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			b, c := fixtureInfo()
			mutate(b)
			if validateInfo(b, c) == nil {
				t.Fatal("untrusted metadata admitted")
			}
		})
	}
	for _, key := range []string{"GOOS", "GOARCH", "GOAMD64", "CGO_ENABLED", "-trimpath", "-compiler", "-buildmode"} {
		t.Run("wrong "+key, func(t *testing.T) {
			b, c := fixtureInfo()
			for i := range b.Settings {
				if b.Settings[i].Key == key {
					b.Settings[i].Value = "wrong"
				}
			}
			if validateInfo(b, c) == nil {
				t.Fatal("platform/build setting admitted")
			}
		})
	}
}
func TestReadRejectsMissingOversizeDirectoryAndLinks(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, []byte("input"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{root, filepath.Join(root, "missing")} {
		if _, err := readFile(p, 32); err == nil {
			t.Fatal("type/missing input admitted")
		}
	}
	if _, err := readFile(file, 2); err == nil {
		t.Fatal("oversize input admitted")
	}
	if err := os.Symlink(file, filepath.Join(root, "link")); err != nil {
		t.Log("symlink creation unavailable; directory/missing/size cases executed")
	} else if _, err := readFile(filepath.Join(root, "link"), 32); err == nil {
		t.Fatal("link admitted")
	}
}
func TestSourceSumRejectsDuplicateOrMalformed(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module "+scannerModule+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	good := "example.invalid/module v1.0.0 h1:Ra4+bf83h2ztPIQYNP99R6m+Y7KfnARDfID+a+vLl4s=\n"
	for _, body := range []string{good + good, "bad", "example.invalid/module v1.0.0 h1:invalid\n"} {
		if err := os.WriteFile(filepath.Join(root, "go.sum"), []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := readSource(root, "tooling", scannerModule, strings.Repeat("a", 40), strings.Repeat("b", 40)); err == nil {
			t.Fatal("bad source sum admitted")
		}
	}
}

// Uses real delivered bytes as data only. The optional config is supplied by
// the bounded local review harness; ordinary repository tests need no download.
func TestDeliveredSixBinaryClosureAndNegativeInputs(t *testing.T) {
	path := os.Getenv("PSCAN_SBOM_TEST_CONFIG")
	if path == "" {
		t.Skip("bounded delivered fixture not supplied")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var o options
	if json.Unmarshal(raw, &o) != nil {
		t.Fatal("config")
	}
	first, err := generate(o)
	if err != nil {
		t.Fatal(err)
	}
	second, err := generate(o)
	if err != nil || !bytes.Equal(first, second) {
		t.Fatal("nondeterministic generation", err)
	}
	if dest := os.Getenv("PSCAN_SBOM_PREVIEW"); dest != "" {
		if err := os.WriteFile(dest, first, 0600); err != nil {
			t.Fatal(err)
		}
	}
	var doc document
	if json.Unmarshal(first, &doc) != nil {
		t.Fatal("document")
	}
	packages := map[string]packageEntry{}
	modules := map[string]bool{}
	sources := map[string]bool{}
	binaryCount := 0
	for _, p := range doc.Packages {
		if _, ok := packages[p.ID]; ok {
			t.Fatal("duplicate component")
		}
		packages[p.ID] = p
		if strings.HasPrefix(p.ID, "SPDXRef-Module-") {
			modules[p.Name+"@"+p.Version] = true
		}
		if strings.HasPrefix(p.ID, "SPDXRef-Source-") {
			sources[p.Version] = true
		}
		if strings.HasPrefix(p.ID, "SPDXRef-Binary-") {
			binaryCount++
			if len(p.Checksums) != 1 {
				t.Fatal("binary checksum")
			}
		}
	}
	if binaryCount != 6 || len(modules) != 63 || len(sources) != 3 || !sources[productRevision] || !sources[engineRevision] || !sources[o.Revision] {
		t.Fatalf("incomplete roots/closure binaries=%d modules=%d roots=%d", binaryCount, len(modules), len(sources))
	}
	for _, key := range []string{"golang.org/x/crypto@v0.35.0", "golang.org/x/sys@v0.30.0", "golang.org/x/exp@v0.0.0-20250218142911-aa4b98e5adaa", "dario.cat/mergo@v1.0.1"} {
		if !modules[key] {
			t.Fatal("regressed missing dependency", key)
		}
	}
	for _, r := range doc.Relationships {
		if r.From != "SPDXRef-DOCUMENT" {
			if _, ok := packages[r.From]; !ok {
				t.Fatal("dangling source")
			}
		}
		if _, ok := packages[r.To]; !ok {
			t.Fatal("dangling target")
		}
	}
	// Independent edge oracle: each binary must name exactly its own compiled
	// dependencies plus its Go runtime. Windows/Linux differences remain visible.
	for _, c := range components() {
		b, e := os.ReadFile(filepath.Join(o.Binaries, c.Name))
		if e != nil {
			t.Fatal(e)
		}
		bi, e := binaryInfo(b, c)
		if e != nil {
			t.Fatal(e)
		}
		actual := map[string]bool{}
		for _, r := range doc.Relationships {
			if r.From == id("Binary", c.Name) && r.Type == "DEPENDS_ON" {
				actual[r.To] = true
			}
		}
		if len(actual) != len(bi.Deps)+1 {
			t.Fatal("per-binary dependency count", c.Name)
		}
		for _, d := range bi.Deps {
			if !actual[id("Module", d.Path+"@"+d.Version)] {
				t.Fatal("missing component-specific edge")
			}
		}
	}
	t.Run("old shipped SBOM fails compiled closure", func(t *testing.T) {
		raw, e := os.ReadFile(filepath.Join(o.Binaries, "sbom.spdx.json"))
		if e != nil {
			t.Fatal(e)
		}
		var old document
		if json.Unmarshal(raw, &old) != nil {
			t.Fatal("old document")
		}
		oldSet := map[string]bool{}
		for _, p := range old.Packages {
			oldSet[p.Name+"@"+p.Version] = true
		}
		missing := 0
		for key := range modules {
			if !oldSet[key] {
				missing++
			}
		}
		if missing != 40 {
			t.Fatalf("expected old defect 40, got %d", missing)
		}
	})
	for _, c := range components() {
		t.Run("missing "+c.Name, func(t *testing.T) {
			dir := t.TempDir()
			for _, other := range components() {
				if other.Name == c.Name {
					continue
				}
				b, _ := os.ReadFile(filepath.Join(o.Binaries, other.Name))
				if e := os.WriteFile(filepath.Join(dir, other.Name), b, 0600); e != nil {
					t.Fatal(e)
				}
			}
			bad := o
			bad.Binaries = dir
			if _, e := generate(bad); e == nil {
				t.Fatal("missing binary accepted")
			}
		})
	}
	t.Run("tampered fixed binary", func(t *testing.T) {
		c := components()[0]
		b, _ := os.ReadFile(filepath.Join(o.Binaries, c.Name))
		b[len(b)-1] ^= 1
		if _, e := binaryInfo(b, c); e == nil {
			t.Fatal("wrong binary digest accepted")
		}
	})
	t.Run("platform swap", func(t *testing.T) {
		c := components()[1]
		b, _ := os.ReadFile(filepath.Join(o.Binaries, "scanner-release-verifier-windows-amd64.exe"))
		if _, e := binaryInfo(b, c); e == nil {
			t.Fatal("wrong platform accepted")
		}
	})
	t.Run("malformed executable", func(t *testing.T) {
		if _, e := binaryInfo([]byte("not a Go executable"), components()[1]); e == nil {
			t.Fatal("malformed accepted")
		}
	})
	for _, name := range []string{"missing sum", "wrong sum", "licence text", "licence manifest"} {
		t.Run(name, func(t *testing.T) {
			bad := o
			if strings.Contains(name, "sum") {
				dir := t.TempDir()
				for _, f := range []string{"go.mod", "go.sum", "LICENSE"} {
					b, _ := os.ReadFile(filepath.Join(o.Engine, f))
					if f == "go.sum" {
						lines := strings.Split(string(b), "\n")
						for i, l := range lines {
							if strings.HasPrefix(l, "dario.cat/mergo v1.0.1 ") {
								if name == "missing sum" {
									lines[i] = ""
								} else {
									lines[i] = "dario.cat/mergo v1.0.1 h1:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
								}
							}
						}
						var nonempty []string
						for _, line := range lines {
							if line != "" {
								nonempty = append(nonempty, line)
							}
						}
						b = []byte(strings.Join(nonempty, "\n") + "\n")
					}
					if e := os.WriteFile(filepath.Join(dir, f), b, 0600); e != nil {
						t.Fatal(e)
					}
				}
				bad.Engine = dir
			} else {
				dir := t.TempDir()
				if e := os.CopyFS(dir, os.DirFS(o.Product)); e != nil {
					t.Fatal(e)
				}
				p := filepath.Join(dir, "licenses/gitleaks/modules/manifest.json")
				if name == "licence text" {
					p = filepath.Join(dir, "licenses/gitleaks/modules/dario.cat_mergo@v1.0.1/LICENSE")
				}
				if e := os.WriteFile(p, []byte("tampered"), 0600); e != nil {
					t.Fatal(e)
				}
				bad.Product = dir
			}
			_, e := generate(bad)
			if e == nil {
				t.Fatal("unbound source/licence accepted")
			}
			if strings.Contains(name, "sum") && !strings.Contains(e.Error(), "compiled module sum absent") {
				t.Fatalf("wrong rejection gate: %v", e)
			}
		})
	}
	for _, name := range []string{"revision", "tree", "created", "toolchain"} {
		t.Run("invalid "+name, func(t *testing.T) {
			bad := o
			switch name {
			case "revision":
				bad.Revision = productRevision
			case "tree":
				bad.Tree = strings.Repeat("0", 40)
			case "created":
				bad.Created = "yesterday"
			case "toolchain":
				bad.RunnerGo = o.EngineGo
			}
			if _, e := generate(bad); e == nil {
				t.Fatal("invalid identity accepted")
			}
		})
	}
}
