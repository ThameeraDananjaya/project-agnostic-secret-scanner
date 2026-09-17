// The builder admits exact source archives before this generator reads the six
// compiled outputs. Input executables are data and are never executed here.
package main

import (
	"bytes"
	"crypto/sha256"
	"debug/buildinfo"
	"debug/elf"
	"debug/pe"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"sort"
	"strings"
	"time"
)

const (
	productRevision       = "a13c28fe7273bc8dc6545f97966a02889524eb4c"
	productTree           = "217b711ddea51fd0ea7e808edd2e27fdecef8427"
	engineRevision        = "83d9cd684c87d95d656c1458ef04895a7f1cbd8e"
	scannerModule         = "github.com/ThameeraDananjaya/project-agnostic-secret-scanner"
	engineModule          = "github.com/zricethezav/gitleaks/v8"
	licenceManifestDigest = "0b8c06685c7d16fd6d685ae0e5e481d85396d2e3c64671151f07634a7bb40b60"
)

var oidPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)
var slugPattern = regexp.MustCompile(`[^A-Za-z0-9._@+-]`)
var modulePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._~+!/-]*$`)

type options struct{ Binaries, Product, Tooling, Engine, RunnerGo, EngineGo, Revision, Tree, Created string }
type checksum struct {
	Algorithm string `json:"algorithm"`
	Value     string `json:"checksumValue"`
}
type reference struct {
	Category string `json:"referenceCategory"`
	Type     string `json:"referenceType"`
	Locator  string `json:"referenceLocator"`
}
type packageEntry struct {
	ID               string      `json:"SPDXID"`
	Name             string      `json:"name"`
	Version          string      `json:"versionInfo"`
	Download         string      `json:"downloadLocation"`
	FilesAnalyzed    bool        `json:"filesAnalyzed"`
	LicenseConcluded string      `json:"licenseConcluded"`
	LicenseDeclared  string      `json:"licenseDeclared"`
	Copyright        string      `json:"copyrightText"`
	Source           string      `json:"sourceInfo"`
	Checksums        []checksum  `json:"checksums,omitempty"`
	Refs             []reference `json:"externalRefs,omitempty"`
}
type relationship struct {
	From string `json:"spdxElementId"`
	Type string `json:"relationshipType"`
	To   string `json:"relatedSpdxElement"`
}
type document struct {
	Version       string         `json:"spdxVersion"`
	License       string         `json:"dataLicense"`
	ID            string         `json:"SPDXID"`
	Name          string         `json:"name"`
	Namespace     string         `json:"documentNamespace"`
	Creation      map[string]any `json:"creationInfo"`
	Packages      []packageEntry `json:"packages"`
	Relationships []relationship `json:"relationships"`
}
type licenceFile struct {
	File   string `json:"file"`
	SHA256 string `json:"sha256"`
}
type licenceModule struct {
	Module   string        `json:"module"`
	Version  string        `json:"version"`
	Licenses []licenceFile `json:"licenses"`
}
type licenceInventory struct {
	Schema  string          `json:"schemaVersion"`
	Scope   string          `json:"scope"`
	Count   int             `json:"moduleCount"`
	Modules []licenceModule `json:"modules"`
}
type licenceBinding struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
type sourceEvidence struct {
	Role, Module, Revision, Tree, GoModSHA256, GoSumSHA256 string
	Sums                                                   map[string]string `json:"-"`
}
type component struct{ Name, OS, Go, Role, Module, MainPath, ExpectedDigest string }

func components() []component {
	var out []component
	for _, platform := range []string{"linux", "windows"} {
		suffix := "-" + platform + "-amd64"
		if platform == "windows" {
			suffix += ".exe"
		}
		runnerHash, engineHash := "06043e9410e05a927356960f7b4ee64c9a073b7467e775452528a6813d97906e", "657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586"
		if platform == "windows" {
			runnerHash = "1b74310e16e0df13cfd077c41338b93eabdc3d5959fb63dc256927ca51117543"
			engineHash = "b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178"
		}
		out = append(out, component{"scanner-runner" + suffix, platform, "go1.27.1", "product", scannerModule, scannerModule + "/cmd/scanner-runner", runnerHash}, component{"scanner-release-verifier" + suffix, platform, "go1.27.1", "tooling", scannerModule, scannerModule + "/build/release/cmd/release-verifier", ""}, component{"gitleaks" + suffix, platform, "go1.27.0", "engine", engineModule, engineModule, engineHash})
	}
	return out
}
func main() {
	var o options
	flag.StringVar(&o.Binaries, "binaries", "", "directory containing the six fixed binary filenames")
	flag.StringVar(&o.Product, "product-source", "", "admitted exact product source directory")
	flag.StringVar(&o.Tooling, "tooling-source", "", "admitted exact tooling source directory")
	flag.StringVar(&o.Engine, "engine-source", "", "admitted digest-pinned Gitleaks source directory")
	flag.StringVar(&o.RunnerGo, "runner-go", "", "admitted Go 1.27.1 root")
	flag.StringVar(&o.EngineGo, "engine-go", "", "admitted Go 1.27.0 root")
	flag.StringVar(&o.Revision, "tooling-revision", "", "exact admitted tooling commit")
	flag.StringVar(&o.Tree, "tooling-tree", "", "exact admitted tooling tree")
	flag.StringVar(&o.Created, "created", "", "canonical UTC timestamp from tooling commit")
	output := flag.String("output", "", "new SPDX 2.3 JSON output path")
	flag.Parse()
	if flag.NArg() != 0 || *output == "" {
		fatal("complete deterministic arguments required")
	}
	raw, err := generate(o)
	if err != nil {
		fatal(err.Error())
	}
	f, err := os.OpenFile(*output, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		fatal(err.Error())
	}
	_, err = f.Write(raw)
	closeErr := f.Close()
	if err != nil || closeErr != nil {
		_ = os.Remove(*output)
		fatal("cannot write complete SBOM")
	}
}
func digest(raw []byte) string   { h := sha256.Sum256(raw); return hex.EncodeToString(h[:]) }
func info(v any) string          { raw, _ := json.Marshal(v); return string(raw) }
func id(kind, key string) string { return "SPDXRef-" + kind + "-" + digest([]byte(key)) }
func packageFor(kind, key, name, version string, provenance any) packageEntry {
	return packageEntry{ID: id(kind, key), Name: name, Version: version, Download: "NOASSERTION", LicenseConcluded: "NOASSERTION", LicenseDeclared: "NOASSERTION", Copyright: "NOASSERTION", Source: info(provenance)}
}

// Hashing and metadata parsing use one bounded byte snapshot. Symlinks on any
// path component reject. Exact-source admission/containment remains the builder's.
func readFile(path string, limit int64) ([]byte, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	for current := abs; ; current = filepath.Dir(current) {
		st, e := os.Lstat(current)
		if e != nil {
			return nil, e
		}
		if st.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("link rejected")
		}
		if parent := filepath.Dir(current); parent == current {
			break
		}
	}
	f, err := os.Open(abs)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil || !st.Mode().IsRegular() || st.Size() <= 0 || st.Size() > limit {
		return nil, fmt.Errorf("invalid input size/type")
	}
	raw, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil || int64(len(raw)) != st.Size() {
		return nil, fmt.Errorf("incomplete input")
	}
	return raw, nil
}
func readSource(root, role, module, revision, tree string) (sourceEvidence, error) {
	m, err := readFile(filepath.Join(root, "go.mod"), 1<<20)
	if err != nil {
		return sourceEvidence{}, err
	}
	normalized := bytes.ReplaceAll(m, []byte("\r\n"), []byte("\n"))
	if !bytes.HasPrefix(normalized, []byte("module "+module+"\n")) {
		return sourceEvidence{}, fmt.Errorf("source module mismatch")
	}
	if role == "engine" && digest(normalized) != "607c140abf2a872e70423972d4dfc7fa658ebe10365d0ea995269ed292add7a3" {
		return sourceEvidence{}, fmt.Errorf("engine source pin mismatch")
	}
	s, err := readFile(filepath.Join(root, "go.sum"), 2<<20)
	if err != nil {
		return sourceEvidence{}, err
	}
	sums := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(string(s)), "\n") {
		f := strings.Fields(line)
		if len(f) != 3 || !validSum(f[2]) {
			return sourceEvidence{}, fmt.Errorf("invalid source sum")
		}
		key := f[0] + "@" + f[1]
		if _, ok := sums[key]; ok {
			return sourceEvidence{}, fmt.Errorf("duplicate source sum")
		}
		sums[key] = f[2]
	}
	return sourceEvidence{role, module, revision, tree, digest(m), digest(s), sums}, nil
}
func validSum(s string) bool {
	if !strings.HasPrefix(s, "h1:") {
		return false
	}
	raw, err := base64.StdEncoding.DecodeString(s[3:])
	return err == nil && len(raw) == 32 && base64.StdEncoding.EncodeToString(raw) == s[3:]
}
func loadLicences(root string) (map[string][]licenceBinding, map[string][]licenceBinding, error) {
	raw, err := readFile(filepath.Join(root, "licenses/gitleaks/modules/manifest.json"), 1<<20)
	if err != nil {
		return nil, nil, err
	}
	if digest(raw) != licenceManifestDigest {
		return nil, nil, fmt.Errorf("licence manifest differs from fixed product")
	}
	var inv licenceInventory
	if json.Unmarshal(raw, &inv) != nil || inv.Schema != "1.0" || inv.Count != 63 || len(inv.Modules) != 63 {
		return nil, nil, fmt.Errorf("licence inventory")
	}
	byModule, byText := map[string][]licenceBinding{}, map[string][]licenceBinding{}
	for _, m := range inv.Modules {
		key := m.Module + "@" + m.Version
		if _, ok := byModule[key]; ok || len(m.Licenses) == 0 {
			return nil, nil, fmt.Errorf("duplicate/empty licence inventory")
		}
		slug := slugPattern.ReplaceAllString(key, "_")
		for _, l := range m.Licenses {
			if filepath.Base(l.File) != l.File || strings.ContainsAny(l.File, "/\\") || l.File == "." || l.File == ".." {
				return nil, nil, fmt.Errorf("unsafe licence filename")
			}
			path := "licenses/gitleaks/modules/" + slug + "/" + l.File
			text, err := readFile(filepath.Join(root, filepath.FromSlash(path)), 1<<20)
			if err != nil || digest(text) != l.SHA256 {
				return nil, nil, fmt.Errorf("licence text binding: %s", key)
			}
			binding := licenceBinding{path, l.SHA256}
			byModule[key] = append(byModule[key], binding)
			byText[l.SHA256] = append(byText[l.SHA256], binding)
		}
	}
	return byModule, byText, nil
}
func binaryInfo(raw []byte, c component) (*debug.BuildInfo, error) {
	if c.ExpectedDigest != "" && digest(raw) != c.ExpectedDigest {
		return nil, fmt.Errorf("fixed product/engine binary digest mismatch")
	}
	if c.OS == "linux" {
		f, e := elf.NewFile(bytes.NewReader(raw))
		if e != nil {
			return nil, e
		}
		defer f.Close()
		if f.Class != elf.ELFCLASS64 || f.Machine != elf.EM_X86_64 || f.Type != elf.ET_EXEC {
			return nil, fmt.Errorf("ELF platform mismatch")
		}
	} else {
		f, e := pe.NewFile(bytes.NewReader(raw))
		if e != nil {
			return nil, e
		}
		defer f.Close()
		if f.Machine != pe.IMAGE_FILE_MACHINE_AMD64 || f.Characteristics&pe.IMAGE_FILE_EXECUTABLE_IMAGE == 0 {
			return nil, fmt.Errorf("PE platform mismatch")
		}
	}
	bi, err := buildinfo.Read(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	if err = validateInfo(bi, c); err != nil {
		return nil, err
	}
	return bi, nil
}
func validateInfo(bi *debug.BuildInfo, c component) error {
	if bi.GoVersion != c.Go || bi.Path != c.MainPath || bi.Main.Path != c.Module || bi.Main.Version != "(devel)" || bi.Main.Replace != nil || bi.Main.Sum != "" {
		return fmt.Errorf("binary root identity mismatch")
	}
	settings := map[string]string{}
	for _, s := range bi.Settings {
		if _, ok := settings[s.Key]; ok {
			return fmt.Errorf("duplicate build setting")
		}
		settings[s.Key] = s.Value
	}
	for k, v := range map[string]string{"GOOS": c.OS, "GOARCH": "amd64", "GOAMD64": "v1", "CGO_ENABLED": "0", "-trimpath": "true", "-compiler": "gc", "-buildmode": "exe"} {
		if settings[k] != v {
			return fmt.Errorf("binary setting mismatch: %s", k)
		}
	}
	for k := range settings {
		if strings.HasPrefix(k, "vcs") {
			return fmt.Errorf("unexpected VCS build setting")
		}
	}
	seen := map[string]bool{}
	for _, d := range bi.Deps {
		if d == nil || !modulePattern.MatchString(d.Path) || strings.Contains(d.Path, "..") || !strings.HasPrefix(d.Version, "v") || strings.ContainsAny(d.Version, " \t\r\n\x00") || d.Replace != nil || !validSum(d.Sum) || seen[d.Path] {
			return fmt.Errorf("invalid, replaced or duplicate binary dependency")
		}
		seen[d.Path] = true
	}
	return nil
}
func generate(o options) ([]byte, error) {
	for _, p := range []string{o.Binaries, o.Product, o.Tooling, o.Engine, o.RunnerGo, o.EngineGo} {
		if p == "" {
			return nil, fmt.Errorf("missing input directory")
		}
	}
	if !oidPattern.MatchString(o.Revision) || !oidPattern.MatchString(o.Tree) || o.Revision == strings.Repeat("0", 40) || o.Tree == strings.Repeat("0", 40) || o.Revision == productRevision {
		return nil, fmt.Errorf("invalid distinct tooling identity")
	}
	when, err := time.Parse(time.RFC3339, o.Created)
	if err != nil || when.Format("2006-01-02T15:04:05Z") != o.Created {
		return nil, fmt.Errorf("canonical creation time required")
	}
	product, err := readSource(o.Product, "product", scannerModule, productRevision, productTree)
	if err != nil {
		return nil, err
	}
	tooling, err := readSource(o.Tooling, "tooling", scannerModule, o.Revision, o.Tree)
	if err != nil {
		return nil, err
	}
	engine, err := readSource(o.Engine, "engine", engineModule, engineRevision, "")
	if err != nil {
		return nil, err
	}
	sources := map[string]sourceEvidence{"product": product, "tooling": tooling, "engine": engine}
	licences, licenceTexts, err := loadLicences(o.Product)
	if err != nil {
		return nil, err
	}
	packages := map[string]packageEntry{}
	var relationships []relationship
	add := func(p packageEntry) error {
		if prior, ok := packages[p.ID]; ok && info(prior) != info(p) {
			return fmt.Errorf("conflicting component")
		}
		packages[p.ID] = p
		return nil
	}
	sourceIDs := map[string]string{}
	for _, s := range []sourceEvidence{product, tooling, engine} {
		p := packageFor("Source", s.Role+"@"+s.Revision, s.Module, s.Revision, s)
		var binding licenceBinding
		if s.Role == "engine" {
			binding = licenceBinding{"licenses/gitleaks/LICENSE-v8.30.1.txt", "e3884b252b3bfc045e55be43a34d1e80da070bc6f804ac95bf4660e97d62ebc6"}
		} else {
			binding = licenceBinding{"LICENSE.txt", "f7f424a6a20ec35897c7016fe3b12a5b8d92365f1d3bb4d9779e8a9e7661f62d"}
		}
		licencePath := filepath.Join(o.Product, filepath.FromSlash(binding.Path))
		if binding.Path == "LICENSE.txt" {
			licencePath = filepath.Join(o.Product, "LICENSE")
		}
		raw, e := readFile(licencePath, 1<<20)
		if e != nil || digest(raw) != binding.SHA256 {
			return nil, fmt.Errorf("root licence binding")
		}
		sourceRoot := map[string]string{"product": o.Product, "tooling": o.Tooling, "engine": o.Engine}[s.Role]
		sourceLicence, e := readFile(filepath.Join(sourceRoot, "LICENSE"), 1<<20)
		if e != nil || !bytes.Equal(sourceLicence, raw) {
			return nil, fmt.Errorf("root source/distributed licence mismatch: %s", s.Role)
		}
		p.Source = info(map[string]any{"source": s, "licence": binding, "admission": "exact source materialization enforced by build before generator"})
		p.LicenseDeclared = "MIT"
		if err := add(p); err != nil {
			return nil, err
		}
		sourceIDs[s.Role] = p.ID
	}
	runtimeIDs := map[string]string{}
	for _, r := range []struct{ Version, Root, ArchiveSHA256 string }{{"go1.27.1", o.RunnerGo, "63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445"}, {"go1.27.0", o.EngineGo, "675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685"}} {
		version, e := readFile(filepath.Join(r.Root, "VERSION"), 1024)
		if e != nil || strings.Split(strings.ReplaceAll(string(version), "\r\n", "\n"), "\n")[0] != r.Version {
			return nil, fmt.Errorf("toolchain version mismatch")
		}
		raw, e := readFile(filepath.Join(r.Root, "LICENSE"), 1<<20)
		if e != nil {
			return nil, e
		}
		h := digest(bytes.ReplaceAll(raw, []byte("\r\n"), []byte("\n")))
		bindings := licenceTexts[h]
		if len(bindings) == 0 {
			return nil, fmt.Errorf("Go runtime licence not shipped")
		}
		p := packageFor("Go", r.Version, "Go runtime and standard library", r.Version, map[string]any{"scope": "runtime and standard library linked into compiled components; compiler is build infrastructure", "expectedBuildToolchainLinuxArchiveSHA256": r.ArchiveSHA256, "licences": bindings})
		p.LicenseDeclared = "BSD-3-Clause"
		if err := add(p); err != nil {
			return nil, err
		}
		runtimeIDs[r.Version] = p.ID
	}
	for _, c := range components() {
		raw, e := readFile(filepath.Join(o.Binaries, c.Name), 64<<20)
		if e != nil {
			return nil, fmt.Errorf("%s: %w", c.Name, e)
		}
		bi, e := binaryInfo(raw, c)
		if e != nil {
			return nil, fmt.Errorf("%s: %w", c.Name, e)
		}
		componentVersion := "1.0.0"
		if c.Role == "engine" {
			componentVersion = "8.30.1"
		}
		p := packageFor("Binary", c.Name, c.Name, componentVersion, map[string]any{"scope": "compiled executable", "releaseVersion": "v1.0.0", "sourceRole": c.Role, "sourceRevision": sources[c.Role].Revision, "goVersion": bi.GoVersion, "os": c.OS, "arch": "amd64", "mainModule": bi.Main.Path, "mainPackage": bi.Path, "compiledExternalModules": len(bi.Deps)})
		p.Checksums = []checksum{{"SHA256", digest(raw)}}
		if err := add(p); err != nil {
			return nil, err
		}
		relationships = append(relationships, relationship{"SPDXRef-DOCUMENT", "DESCRIBES", p.ID}, relationship{p.ID, "GENERATED_FROM", sourceIDs[c.Role]}, relationship{p.ID, "DEPENDS_ON", runtimeIDs[c.Go]})
		for _, d := range bi.Deps {
			mod, err := dependencyPackage(d, sources[c.Role].Sums, licences)
			if err != nil {
				return nil, err
			}
			if err := add(mod); err != nil {
				return nil, err
			}
			relationships = append(relationships, relationship{p.ID, "DEPENDS_ON", mod.ID})
		}
	}
	doc := document{Version: "SPDX-2.3", License: "CC0-1.0", ID: "SPDXRef-DOCUMENT", Name: "project-agnostic-secret-scanner-v1.0.0", Namespace: "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/releases/v1.0.0/spdx/" + productRevision + "/" + o.Revision, Creation: map[string]any{"created": o.Created, "creators": []string{"Tool: scanner-release-sbom-2.0.0"}}}
	for _, p := range packages {
		doc.Packages = append(doc.Packages, p)
	}
	sort.Slice(doc.Packages, func(i, j int) bool { return doc.Packages[i].ID < doc.Packages[j].ID })
	sort.Slice(relationships, func(i, j int) bool { return info(relationships[i]) < info(relationships[j]) })
	doc.Relationships = relationships
	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(raw, '\n'), nil
}
func fatal(message string) { fmt.Fprintln(os.Stderr, "sbom:", message); os.Exit(1) }

func dependencyPackage(d *debug.Module, sums map[string]string, licences map[string][]licenceBinding) (packageEntry, error) {
	key := d.Path + "@" + d.Version
	if sums[key] != d.Sum {
		return packageEntry{}, fmt.Errorf("compiled module sum absent from exact source: %s", key)
	}
	binding := licences[key]
	if len(binding) == 0 {
		return packageEntry{}, fmt.Errorf("compiled dependency has no shipped licence: %s", key)
	}
	mod := packageFor("Module", key, d.Path, d.Version, map[string]any{"scope": "compiled external Go module", "goModuleSum": d.Sum, "licences": binding})
	// h1 hashes module-directory contents, not a distributable archive.
	parts := strings.Split(d.Path, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	mod.Refs = []reference{{"PACKAGE-MANAGER", "purl", "pkg:golang/" + strings.Join(parts, "/") + "@" + url.PathEscape(d.Version)}}
	return mod, nil
}
