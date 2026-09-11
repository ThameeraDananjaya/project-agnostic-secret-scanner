package supplychain_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

func TestMain(m *testing.M) {
	os.Exit(testMainExitCode(os.Args[1:], os.Getenv, m.Run))
}

func testMainExitCode(arguments []string, getenv func(string) string, run func() int) int {
	argumentsPath, admitted := cosignHelperArgumentsPath(arguments, getenv)
	if !admitted {
		return run()
	}
	file, err := os.OpenFile(argumentsPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return 1
	}
	_, writeErr := file.WriteString(strings.Join(arguments, "\n") + "\n")
	closeErr := file.Close()
	if writeErr != nil || closeErr != nil {
		return 1
	}
	return 0
}

func cosignHelperArgumentsPath(arguments []string, getenv func(string) string) (string, bool) {
	if len(arguments) != 18 || getenv("COSIGN_YES") != "false" {
		return "", false
	}
	want := map[int]string{
		0: "verify-blob", 1: "--bundle", 3: "--trusted-root",
		5: "--certificate-identity", 7: "--certificate-oidc-issuer",
		9: "--certificate-github-workflow-repository", 11: "--certificate-github-workflow-ref",
		13: "--certificate-github-workflow-sha", 15: "--certificate-github-workflow-trigger",
	}
	for position, value := range want {
		if arguments[position] != value {
			return "", false
		}
	}
	for _, position := range []int{2, 4, 6, 8, 10, 12, 14, 16, 17} {
		if arguments[position] == "" {
			return "", false
		}
	}

	root := filepath.Clean(filepath.Dir(arguments[4]))
	if !plainDirectory(root) {
		return "", false
	}
	for _, position := range []int{2, 4, 17} {
		if !containedRegularFile(root, arguments[position]) {
			return "", false
		}
	}
	argumentsPath := filepath.Join(filepath.Dir(filepath.Clean(arguments[2])), "arguments.txt")
	if _, err := os.Lstat(argumentsPath); err == nil || !os.IsNotExist(err) {
		return "", false
	}
	return argumentsPath, true
}

func plainDirectory(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.IsDir() && info.Mode()&os.ModeSymlink == 0
}

func containedRegularFile(root, path string) bool {
	if !containedPath(root, path) {
		return false
	}
	info, err := os.Lstat(filepath.Clean(path))
	if err != nil || !plainRegularMode(info.Mode()) {
		return false
	}
	evaluatedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return false
	}
	evaluatedPath, err := filepath.EvalSymlinks(path)
	return err == nil && containedPath(evaluatedRoot, evaluatedPath)
}

func containedPath(root, path string) bool {
	absoluteRoot, err := filepath.Abs(filepath.Clean(root))
	if err != nil {
		return false
	}
	absolutePath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return false
	}
	relative, err := filepath.Rel(absoluteRoot, absolutePath)
	return err == nil && relative != "." && relative != ".." &&
		!filepath.IsAbs(relative) && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func plainRegularMode(mode os.FileMode) bool {
	return mode.IsRegular() && mode&os.ModeSymlink == 0
}

func TestCosignCommandVerifierBindsExecutableAndTrustedRootBytes(t *testing.T) {
	root := t.TempDir()
	cosignPath := filepath.Join(root, "cosign")
	trustedRootPath := filepath.Join(root, "trusted-root.json")
	cosignBytes := []byte("synthetic-cosign-executable")
	trustedRootBytes := []byte(`{"synthetic":"trusted-root"}`)
	if err := os.WriteFile(cosignPath, cosignBytes, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(trustedRootPath, trustedRootBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := verify.NewCosignCommandVerifier(cosignPath, verify.DigestBytes(cosignBytes), trustedRootPath, verify.DigestBytes(trustedRootBytes)); err != nil {
		t.Fatalf("exact Cosign and trusted-root bytes rejected: %v", err)
	}
	trustedRootBytes[0] ^= 1
	if err := os.WriteFile(trustedRootPath, trustedRootBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := verify.NewCosignCommandVerifier(cosignPath, verify.DigestBytes(cosignBytes), trustedRootPath, verify.DigestBytes([]byte(`{"synthetic":"trusted-root"}`))); !errors.Is(err, verify.ErrBindingMismatch) {
		t.Fatalf("mutated trusted-root bytes did not fail closed: %v", err)
	}
}

func TestCosignCommandVerifierBindsEveryGitHubWorkflowClaim(t *testing.T) {
	root := t.TempDir()
	argumentsPath := filepath.Join(root, "arguments.txt")
	trustedRootPath := filepath.Join(root, "trusted-root.json")
	manifestPath := filepath.Join(root, "release-manifest.json")
	bundlePath := filepath.Join(root, "release-manifest.sigstore.json")
	cosignPath, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	cosignBytes, err := os.ReadFile(cosignPath)
	if err != nil {
		t.Fatal(err)
	}
	trustedRoot := []byte(`{"synthetic":"trusted-root"}`)
	if err := os.WriteFile(trustedRootPath, trustedRoot, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(manifestPath, []byte("manifest"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bundlePath, []byte("bundle"), 0o600); err != nil {
		t.Fatal(err)
	}
	verifier, err := verify.NewCosignCommandVerifier(cosignPath, verify.DigestBytes(cosignBytes), trustedRootPath, verify.DigestBytes(trustedRoot))
	if err != nil {
		t.Fatal(err)
	}
	identity := releaseIdentityV2()
	if err := verifier.VerifyManifest(t.Context(), manifestPath, bundlePath, identity); err != nil {
		t.Fatalf("synthetic Cosign invocation rejected: %v", err)
	}
	raw, err := os.ReadFile(argumentsPath)
	if err != nil {
		t.Fatal(err)
	}
	arguments := string(raw)
	for _, required := range []string{
		"--certificate-identity\n" + identity.CertificateIdentity,
		"--certificate-oidc-issuer\n" + identity.OIDCIssuer,
		"--certificate-github-workflow-repository\n" + identity.Repository,
		"--certificate-github-workflow-ref\n" + identity.Ref,
		"--certificate-github-workflow-sha\n" + identity.WorkflowSHA,
		"--certificate-github-workflow-trigger\n" + identity.Trigger,
	} {
		if !strings.Contains(arguments, required) {
			t.Fatalf("required Cosign claim argument absent: %q in %q", required, arguments)
		}
	}
}

func TestCosignTestMainHelperWritesOnlyBesideContainedNestedBundle(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "nested")
	if err := os.Mkdir(nested, 0o700); err != nil {
		t.Fatal(err)
	}
	trustedRootPath := filepath.Join(root, "trusted-root.json")
	manifestPath := filepath.Join(root, "release-manifest.json")
	bundlePath := filepath.Join(nested, "release-manifest.sigstore.json")
	for path, contents := range map[string]string{
		trustedRootPath: "trusted-root",
		manifestPath:    "manifest",
		bundlePath:      "bundle",
	} {
		if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	identity := releaseIdentityV2()
	arguments := []string{
		"verify-blob",
		"--bundle", bundlePath,
		"--trusted-root", trustedRootPath,
		"--certificate-identity", identity.CertificateIdentity,
		"--certificate-oidc-issuer", identity.OIDCIssuer,
		"--certificate-github-workflow-repository", identity.Repository,
		"--certificate-github-workflow-ref", identity.Ref,
		"--certificate-github-workflow-sha", identity.WorkflowSHA,
		"--certificate-github-workflow-trigger", identity.Trigger,
		manifestPath,
	}
	getenv := func(name string) string {
		if name == "COSIGN_YES" {
			return "false"
		}
		return ""
	}
	bundleSibling := filepath.Join(filepath.Dir(filepath.Clean(arguments[2])), "arguments.txt")
	trustedRootSibling := filepath.Join(filepath.Dir(filepath.Clean(arguments[4])), "arguments.txt")
	if bundleSibling == trustedRootSibling {
		t.Fatal("nested bundle and trusted-root output paths are not distinct")
	}
	gotPath, admitted := cosignHelperArgumentsPath(arguments, getenv)
	if !admitted || gotPath != bundleSibling {
		t.Fatalf("contained nested bundle not admitted beside bundle: admitted=%v path=%q want=%q", admitted, gotPath, bundleSibling)
	}
	runs := 0
	if result := testMainExitCode(arguments, getenv, func() int {
		runs++
		return 73
	}); result != 0 || runs != 0 {
		t.Fatalf("admitted nested bundle did not use helper exclusively: result=%d runs=%d", result, runs)
	}
	raw, err := os.ReadFile(bundleSibling)
	if err != nil {
		t.Fatal(err)
	}
	if want := strings.Join(arguments, "\n") + "\n"; string(raw) != want {
		t.Fatalf("captured arguments differ from exact ordered invocation: got %q want %q", raw, want)
	}
	for position, want := range map[int]string{
		6: identity.CertificateIdentity, 8: identity.OIDCIssuer,
		10: identity.Repository, 12: identity.Ref,
		14: identity.WorkflowSHA, 16: identity.Trigger,
	} {
		if arguments[position] != want {
			t.Fatalf("claim position %d = %q, want %q", position, arguments[position], want)
		}
	}
	if _, err := os.Lstat(trustedRootSibling); !os.IsNotExist(err) {
		t.Fatalf("helper created distinct trusted-root-sibling output: %v", err)
	}
}

func TestCosignTestMainHelperRejectsNearMissesWithoutRecursiveExecution(t *testing.T) {
	identity := releaseIdentityV2()
	replaceWithSymlink := func(t *testing.T, arguments []string, position int) {
		t.Helper()
		target := filepath.Join(t.TempDir(), "regular-target")
		if err := os.WriteFile(target, []byte("target"), 0o600); err != nil {
			t.Fatal(err)
		}
		link := arguments[position]
		if err := os.Remove(link); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, link); err != nil {
			t.Skipf("file symlinks unavailable: %v", err)
		}
	}
	newFixture := func(t *testing.T) ([]string, string) {
		t.Helper()
		root := t.TempDir()
		trustedRootPath := filepath.Join(root, "trusted-root.json")
		manifestPath := filepath.Join(root, "release-manifest.json")
		bundlePath := filepath.Join(root, "release-manifest.sigstore.json")
		for path, contents := range map[string]string{
			trustedRootPath: "trusted-root",
			manifestPath:    "manifest",
			bundlePath:      "bundle",
		} {
			if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
				t.Fatal(err)
			}
		}
		return []string{
			"verify-blob",
			"--bundle", bundlePath,
			"--trusted-root", trustedRootPath,
			"--certificate-identity", identity.CertificateIdentity,
			"--certificate-oidc-issuer", identity.OIDCIssuer,
			"--certificate-github-workflow-repository", identity.Repository,
			"--certificate-github-workflow-ref", identity.Ref,
			"--certificate-github-workflow-sha", identity.WorkflowSHA,
			"--certificate-github-workflow-trigger", identity.Trigger,
			manifestPath,
		}, filepath.Join(root, "arguments.txt")
	}

	tests := map[string]func(*testing.T, []string){
		"wrong length": func(_ *testing.T, _ []string) {},
		"wrong command": func(_ *testing.T, arguments []string) {
			arguments[0] = "sign-blob"
		},
		"wrong flag position": func(_ *testing.T, arguments []string) {
			arguments[1], arguments[3] = arguments[3], arguments[1]
		},
		"empty value": func(_ *testing.T, arguments []string) {
			arguments[6] = ""
		},
		"bundle escape": func(t *testing.T, arguments []string) {
			outside := t.TempDir()
			arguments[2] = filepath.Join(outside, "release-manifest.sigstore.json")
			if err := os.WriteFile(arguments[2], []byte("bundle"), 0o600); err != nil {
				t.Fatal(err)
			}
		},
		"manifest escape": func(t *testing.T, arguments []string) {
			outside := t.TempDir()
			arguments[17] = filepath.Join(outside, "release-manifest.json")
			if err := os.WriteFile(arguments[17], []byte("manifest"), 0o600); err != nil {
				t.Fatal(err)
			}
		},
		"trusted-root escape": func(t *testing.T, arguments []string) {
			outside := t.TempDir()
			arguments[4] = filepath.Join(outside, "trusted-root.json")
			if err := os.WriteFile(arguments[4], []byte("trusted-root"), 0o600); err != nil {
				t.Fatal(err)
			}
		},
		"bundle symlink": func(t *testing.T, arguments []string) {
			replaceWithSymlink(t, arguments, 2)
		},
		"trusted-root symlink": func(t *testing.T, arguments []string) {
			replaceWithSymlink(t, arguments, 4)
		},
		"manifest symlink": func(t *testing.T, arguments []string) {
			replaceWithSymlink(t, arguments, 17)
		},
		"bundle directory": func(t *testing.T, arguments []string) {
			arguments[2] = filepath.Join(filepath.Dir(arguments[4]), "bundle-directory")
			if err := os.Mkdir(arguments[2], 0o700); err != nil {
				t.Fatal(err)
			}
		},
		"trusted-root directory": func(t *testing.T, arguments []string) {
			arguments[4] = filepath.Join(filepath.Dir(arguments[4]), "trusted-root-directory")
			if err := os.Mkdir(arguments[4], 0o700); err != nil {
				t.Fatal(err)
			}
		},
		"manifest directory": func(t *testing.T, arguments []string) {
			arguments[17] = filepath.Join(filepath.Dir(arguments[4]), "manifest-directory")
			if err := os.Mkdir(arguments[17], 0o700); err != nil {
				t.Fatal(err)
			}
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			arguments, argumentsPath := newFixture(t)
			if name == "wrong length" {
				arguments = arguments[:len(arguments)-1]
			} else {
				mutate(t, arguments)
			}
			runs := 0
			result := testMainExitCode(arguments, func(name string) string {
				if name == "COSIGN_YES" {
					return "false"
				}
				return ""
			}, func() int {
				runs++
				return 73
			})
			if result != 73 || runs != 1 {
				t.Fatalf("near miss entered helper: result=%d runs=%d", result, runs)
			}
			if _, err := os.Lstat(argumentsPath); !os.IsNotExist(err) {
				t.Fatalf("near miss created or modified arguments output: %v", err)
			}
		})
	}

	for name, value := range map[string]string{"absent COSIGN_YES": "", "wrong COSIGN_YES": "true"} {
		t.Run(name, func(t *testing.T) {
			arguments, argumentsPath := newFixture(t)
			runs := 0
			result := testMainExitCode(arguments, func(string) string { return value }, func() int {
				runs++
				return 73
			})
			if result != 73 || runs != 1 {
				t.Fatalf("environment near miss entered helper: result=%d runs=%d", result, runs)
			}
			if _, err := os.Lstat(argumentsPath); !os.IsNotExist(err) {
				t.Fatalf("environment near miss created or modified arguments output: %v", err)
			}
		})
	}

	t.Run("pre-existing exact output", func(t *testing.T) {
		arguments, argumentsPath := newFixture(t)
		original := []byte("preserve-existing-output")
		if err := os.WriteFile(argumentsPath, original, 0o600); err != nil {
			t.Fatal(err)
		}
		runs := 0
		result := testMainExitCode(arguments, func(name string) string {
			if name == "COSIGN_YES" {
				return "false"
			}
			return ""
		}, func() int {
			runs++
			return 73
		})
		if result != 73 || runs != 1 {
			t.Fatalf("pre-existing output entered helper: result=%d runs=%d", result, runs)
		}
		raw, err := os.ReadFile(argumentsPath)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(raw, original) {
			t.Fatalf("pre-existing output mutated: got %q want %q", raw, original)
		}
	})

	for name, mode := range map[string]os.FileMode{"symlink": os.ModeSymlink, "reparse-like irregular": os.ModeIrregular} {
		t.Run(name, func(t *testing.T) {
			if plainRegularMode(mode) {
				t.Fatalf("hostile file mode admitted: %v", mode)
			}
		})
	}
}

func TestReleaseManifestParserRejectsUnknownDuplicateAndUnsupportedMembers(t *testing.T) {
	if _, err := verify.ParseReleaseManifest([]byte(`{"schemaFamily":"scanner-release-manifest","schemaFamily":"scanner-release-manifest"}`)); !errors.Is(err, verify.ErrInvalidReference) {
		t.Fatalf("duplicate member was not rejected: %v", err)
	}
	manifest := minimumManifest()
	raw, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	raw = append(raw[:len(raw)-1], []byte(`,"unexpected":true}`)...)
	if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
		t.Fatalf("unknown member was not rejected: %v", err)
	}
	manifest.ManifestSchemaVersion = "3.0"
	raw, _ = json.Marshal(manifest)
	if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
		t.Fatalf("unknown schema major was not rejected: %v", err)
	}
}

func TestReleaseManifestV2RequiresDistinctCompleteIdentities(t *testing.T) {
	manifest := minimumManifestV2()
	raw, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := verify.ParseReleaseManifest(raw); err != nil {
		t.Fatalf("exact schema 2.0 manifest rejected: %v", err)
	}

	tests := map[string]func(*verify.ReleaseManifest){
		"old schema masquerade": func(m *verify.ReleaseManifest) { m.SourceRevision = m.ProductSource.Commit },
		"missing product role":  func(m *verify.ReleaseManifest) { m.ProductSource = nil },
		"missing tooling role":  func(m *verify.ReleaseManifest) { m.ReleaseTooling = nil },
		"swapped tags": func(m *verify.ReleaseManifest) {
			m.ProductSource.Tag, m.ReleaseTooling.Tag = m.ReleaseTooling.Tag, m.ProductSource.Tag
		},
		"workflow sha differs from tooling commit": func(m *verify.ReleaseManifest) {
			m.ReleaseTooling.WorkflowSHA = "1111111111111111111111111111111111111111"
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			candidate := minimumManifestV2()
			mutate(&candidate)
			raw, err := json.Marshal(candidate)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
				t.Fatalf("invalid identity arrangement was not rejected: %v", err)
			}
		})
	}
}

func TestReleaseManifestV21BindsCorrectionC2AndPreservesC1Parsing(t *testing.T) {
	c1 := minimumManifestV2()
	raw, err := json.Marshal(c1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := verify.ParseReleaseManifest(raw); err != nil {
		t.Fatalf("historical schema 2.0 C1 manifest rejected: %v", err)
	}

	c2 := minimumManifestV21()
	raw, err = json.Marshal(c2)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := verify.ParseReleaseManifest(raw); err != nil {
		t.Fatalf("schema 2.1 C2 manifest rejected: %v", err)
	}

	for name, mutate := range map[string]func(*verify.ReleaseManifest){
		"C1 tag under 2.1": func(m *verify.ReleaseManifest) { m.ReleaseTooling.Tag = "release-tooling-v1.0.0-c1" },
		"C1 ref under 2.1": func(m *verify.ReleaseManifest) { m.ReleaseTooling.WorkflowRef = "refs/tags/release-tooling-v1.0.0-c1" },
		"role swap": func(m *verify.ReleaseManifest) {
			m.ProductSource.Tag, m.ReleaseTooling.Tag = m.ReleaseTooling.Tag, m.ProductSource.Tag
		},
	} {
		t.Run(name, func(t *testing.T) {
			candidate := minimumManifestV21()
			mutate(&candidate)
			raw, _ := json.Marshal(candidate)
			if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
				t.Fatalf("invalid C2 identity was not rejected: %v", err)
			}
		})
	}
}

func TestReleaseManifestV22BindsR6AndRejectsCrossVersionMixtures(t *testing.T) {
	for name, manifest := range map[string]verify.ReleaseManifest{
		"historical 2.0 C1": minimumManifestV2(),
		"historical 2.1 C2": minimumManifestV21(),
		"current 2.2 R6":    minimumManifestV22(),
	} {
		t.Run(name, func(t *testing.T) {
			raw, err := json.Marshal(manifest)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := verify.ParseReleaseManifest(raw); err != nil {
				t.Fatalf("valid versioned manifest rejected: %v", err)
			}
		})
	}

	for name, mutate := range map[string]func(*verify.ReleaseManifest){
		"old tag under 2.2": func(m *verify.ReleaseManifest) { m.ReleaseTooling.Tag = "release-tooling-v1.0.0-c2" },
		"old ref under 2.2": func(m *verify.ReleaseManifest) { m.ReleaseTooling.WorkflowRef = "refs/tags/release-tooling-v1.0.0-c2" },
		"old path under 2.2": func(m *verify.ReleaseManifest) {
			m.ReleaseTooling.Workflow = ".github/workflows/release-recovery-v1.0.0.yml"
		},
		"old identity under 2.2": func(m *verify.ReleaseManifest) {
			m.ReleaseIdentity.CertificateIdentity = "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-recovery-v1.0.0.yml@refs/tags/release-tooling-v1.0.0-c2"
		},
	} {
		t.Run(name, func(t *testing.T) {
			candidate := minimumManifestV22()
			mutate(&candidate)
			raw, _ := json.Marshal(candidate)
			if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
				t.Fatalf("2.2 cross-version identity was not rejected: %v", err)
			}
		})
	}

	for name, mutate := range map[string]func(*verify.ReleaseManifest){
		"new tag under 2.1": func(m *verify.ReleaseManifest) { m.ReleaseTooling.Tag = "release-tooling-v1.0.0-c2-r6" },
		"new ref under 2.1": func(m *verify.ReleaseManifest) {
			m.ReleaseTooling.WorkflowRef = "refs/tags/release-tooling-v1.0.0-c2-r6"
		},
		"new path under 2.1": func(m *verify.ReleaseManifest) {
			m.ReleaseTooling.Workflow = ".github/workflows/release-recovery-v1.0.0-c2-r6.yml"
		},
		"new identity under 2.1": func(m *verify.ReleaseManifest) {
			m.ReleaseIdentity.CertificateIdentity = "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-recovery-v1.0.0-c2-r6.yml@refs/tags/release-tooling-v1.0.0-c2-r6"
		},
	} {
		t.Run(name, func(t *testing.T) {
			candidate := minimumManifestV21()
			mutate(&candidate)
			raw, _ := json.Marshal(candidate)
			if _, err := verify.ParseReleaseManifest(raw); !errors.Is(err, verify.ErrInvalidReference) {
				t.Fatalf("2.1 cross-version identity was not rejected: %v", err)
			}
		})
	}
}

func TestIteration007RepositoryIdentityAgreementAndPreservation(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", "..", ".."))
	read := func(path string) string {
		t.Helper()
		raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatal(err)
		}
		return string(raw)
	}

	oldSchema := read("contracts/release-manifest/schema-2.1.json")
	if verify.DigestBytes([]byte(oldSchema)) != "caa9cd26665cc3a3550affea7490a0f1f277a69b10b1f1537787616e8ab973ce" {
		t.Fatal("historical schema 2.1 bytes changed")
	}
	newSchema := read("contracts/release-manifest/schema-2.2.json")
	normalizedSchema := strings.ReplaceAll(newSchema, "release-tooling-v1.0.0-c2-r6", "release-tooling-v1.0.0-c2")
	normalizedSchema = strings.ReplaceAll(normalizedSchema, ".github/workflows/release-recovery-v1.0.0-c2-r6.yml", ".github/workflows/release-recovery-v1.0.0.yml")
	normalizedSchema = strings.ReplaceAll(normalizedSchema, "2.2", "2.1")
	if !reflect.DeepEqual([]byte(normalizedSchema), []byte(oldSchema)) {
		t.Fatal("schema 2.2 differs from schema 2.1 outside the selected identity/version fields")
	}

	oldWorkflow := read(".github/workflows/release-recovery-v1.0.0.yml")
	if verify.DigestBytes([]byte(oldWorkflow)) != "c5f40f1b32e87c005fe33ee607af7e3d19ee4f0173c21619e21158c31aa0bdb4" {
		t.Fatal("historical C2 workflow bytes changed")
	}
	newWorkflow := read(".github/workflows/release-recovery-v1.0.0-c2-r6.yml")
	normalizedWorkflow := newWorkflow
	for _, replacement := range []struct {
		from  string
		to    string
		count int
	}{
		{"gated-v1.0.0-c2-r6-recovery", "gated-v1.0.0-c2-recovery", 1},
		{"release-v1.0.0-c2-r6-recovery", "release-v1.0.0-c2-recovery", 1},
		{"release-tooling-v1.0.0-c2-r6", "release-tooling-v1.0.0-c2", 8},
		{".github/workflows/release-recovery-v1.0.0-c2-r6.yml", ".github/workflows/release-recovery-v1.0.0.yml", 3},
		{`test "$(jq -r .manifestSchemaVersion dist/release-manifest.json)" = '2.2'`, `test "$(jq -r .manifestSchemaVersion dist/release-manifest.json)" = '2.1'`, 1},
	} {
		if count := strings.Count(newWorkflow, replacement.from); count != replacement.count {
			t.Fatalf("R6 workflow token %q count = %d, want %d", replacement.from, count, replacement.count)
		}
		normalizedWorkflow = strings.ReplaceAll(normalizedWorkflow, replacement.from, replacement.to)
	}
	actionPin := "uses: actions/attest@1e69f48acb82d1966a394da916b4c1698aa569d6 # v4.2.2"
	for name, workflow := range map[string]string{"old": oldWorkflow, "new": newWorkflow, "normalized": normalizedWorkflow} {
		if count := strings.Count(workflow, actionPin); count != 1 {
			t.Fatalf("%s workflow action pin count = %d, want 1", name, count)
		}
		if strings.Contains(workflow, "# v4.2.1") {
			t.Fatalf("%s workflow contains forbidden action-pin normalization", name)
		}
	}
	if normalizedWorkflow != oldWorkflow {
		t.Fatal("R6 workflow differs from the C2 workflow outside the selected identity/version fields")
	}

	for path, required := range map[string][]string{
		"build/release/build.ps1":                    {"release-tooling-v1.0.0-c2-r6", ".github/workflows/release-recovery-v1.0.0-c2-r6.yml", "manifestSchemaVersion='2.2'", "schema-release-manifest-2.1.json", "schema-release-manifest-2.2.json"},
		"build/release/cmd/release-verifier/main.go": {"release-tooling-v1.0.0-c2-r6", ".github/workflows/release-recovery-v1.0.0-c2-r6.yml", `ManifestSchemaVersion: "2.2"`},
		"internal/verify/release.go":                 {"release-tooling-v1.0.0-c2", "release-tooling-v1.0.0-c2-r6", ".github/workflows/release-recovery-v1.0.0.yml", ".github/workflows/release-recovery-v1.0.0-c2-r6.yml"},
	} {
		source := read(path)
		for _, token := range required {
			if !strings.Contains(source, token) {
				t.Fatalf("%s does not contain required identity token %q", path, token)
			}
		}
	}
}

func minimumManifest() verify.ReleaseManifest {
	digest := verify.DigestBytes([]byte("x"))
	assets := []verify.ReleaseAsset{}
	for _, path := range []string{"runner-linux", "runner-windows.exe", "engine-linux", "engine-windows.exe", "rules.toml", "schema.json", "tests.json", "limitations.md", "revocations.json", "checkpoint.json"} {
		assets = append(assets, verify.ReleaseAsset{Path: path, Kind: "documentation", OS: "none", Arch: "none", Size: 1, SHA256: digest})
	}
	return verify.ReleaseManifest{
		SchemaFamily: "scanner-release-manifest", ManifestSchemaVersion: "1.1", ReleaseVersion: "v1.0.0",
		SourceRevision: "0123456789abcdef0123456789abcdef01234567", SourceTree: "89abcdef0123456789abcdef0123456789abcdef",
		RunnerVersion: "1.0.0", GoToolchainVersion: "go1.27.1",
		RunnerBindings: []verify.PlatformDigest{{OS: "linux", Arch: "amd64", Path: "runner-linux", SHA256: digest}, {OS: "windows", Arch: "amd64", Path: "runner-windows.exe", SHA256: digest}},
		EngineBindings: []verify.EngineBinding{
			{Name: "gitleaks", Version: "8.30.1", SourceRevision: "83d9cd684c87d95d656c1458ef04895a7f1cbd8e", OS: "linux", Arch: "amd64", Path: "engine-linux", SHA256: digest},
			{Name: "gitleaks", Version: "8.30.1", SourceRevision: "83d9cd684c87d95d656c1458ef04895a7f1cbd8e", OS: "windows", Arch: "amd64", Path: "engine-windows.exe", SHA256: digest},
		},
		RulePack:        verify.NamedDigest{Path: "rules.toml", SHA256: digest},
		SchemaBindings:  []verify.SchemaBinding{{Family: "scan-request", Version: "1.1", Path: "schema.json", SHA256: digest}},
		ReleaseIdentity: releaseIdentity(),
		Compatibility:   verify.ReleaseCompatibility{MinimumRunnerVersion: "1.0.0", SupportedOperatingSystems: []string{"linux", "windows"}, SupportedArchitectures: []string{"amd64"}, TestSummaryAsset: "tests.json", LimitationsAsset: "limitations.md"},
		Revocation:      verify.ReleaseRevocation{DiscoveryLocation: "https://api.github.com/repos/ThameeraDananjaya/project-agnostic-secret-scanner/releases?per_page=100", SnapshotAsset: "revocations.json", CheckpointAsset: "checkpoint.json", SchemaVersion: "1.1", MaximumSnapshotAgeHours: 24},
		Assets:          assets, CreatedAt: "2026-09-03T00:00:00Z",
	}
}

func minimumManifestV2() verify.ReleaseManifest {
	manifest := minimumManifest()
	manifest.ManifestSchemaVersion = "2.0"
	manifest.SourceRevision = ""
	manifest.SourceTree = ""
	manifest.ProductSource = &verify.ProductSourceIdentity{Tag: "v1.0.0", Commit: "a13c28fe7273bc8dc6545f97966a02889524eb4c", Tree: "217b711ddea51fd0ea7e808edd2e27fdecef8427"}
	manifest.ReleaseTooling = &verify.ReleaseToolingIdentity{
		Tag: "release-tooling-v1.0.0-c1", Commit: "3fb1b0a55dc4f48dd35464c63c768f497efbc89b", Tree: "544867396910969d20cfd2acd454ac0d69c1d7e4",
		Workflow: ".github/workflows/release-recovery-v1.0.0.yml", WorkflowRef: "refs/tags/release-tooling-v1.0.0-c1", WorkflowSHA: "3fb1b0a55dc4f48dd35464c63c768f497efbc89b", Trigger: "workflow_dispatch",
	}
	manifest.ReleaseIdentity = releaseIdentityV2()
	return manifest
}

func minimumManifestV21() verify.ReleaseManifest {
	manifest := minimumManifestV2()
	manifest.ManifestSchemaVersion = "2.1"
	manifest.ReleaseTooling.Tag = "release-tooling-v1.0.0-c2"
	manifest.ReleaseTooling.WorkflowRef = "refs/tags/release-tooling-v1.0.0-c2"
	manifest.ReleaseIdentity.Ref = manifest.ReleaseTooling.WorkflowRef
	manifest.ReleaseIdentity.CertificateIdentity = "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-recovery-v1.0.0.yml@refs/tags/release-tooling-v1.0.0-c2"
	return manifest
}

func minimumManifestV22() verify.ReleaseManifest {
	manifest := minimumManifestV21()
	manifest.ManifestSchemaVersion = "2.2"
	manifest.ReleaseTooling.Tag = "release-tooling-v1.0.0-c2-r6"
	manifest.ReleaseTooling.Workflow = ".github/workflows/release-recovery-v1.0.0-c2-r6.yml"
	manifest.ReleaseTooling.WorkflowRef = "refs/tags/release-tooling-v1.0.0-c2-r6"
	manifest.ReleaseIdentity.Workflow = manifest.ReleaseTooling.Workflow
	manifest.ReleaseIdentity.Ref = manifest.ReleaseTooling.WorkflowRef
	manifest.ReleaseIdentity.CertificateIdentity = "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-recovery-v1.0.0-c2-r6.yml@refs/tags/release-tooling-v1.0.0-c2-r6"
	return manifest
}

func releaseIdentity() verify.ReleaseIdentity {
	return verify.ReleaseIdentity{Repository: "ThameeraDananjaya/project-agnostic-secret-scanner", RepositoryOwnerID: 50274860, Workflow: ".github/workflows/release.yml", Ref: "refs/tags/v1.0.0", OIDCIssuer: "https://token.actions.githubusercontent.com", CertificateIdentity: "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release.yml@refs/tags/v1.0.0"}
}

func releaseIdentityV2() verify.ReleaseIdentity {
	return verify.ReleaseIdentity{
		Repository: "ThameeraDananjaya/project-agnostic-secret-scanner", RepositoryOwnerID: 50274860,
		Workflow: ".github/workflows/release-recovery-v1.0.0.yml", Ref: "refs/tags/release-tooling-v1.0.0-c1",
		WorkflowSHA: "3fb1b0a55dc4f48dd35464c63c768f497efbc89b", Trigger: "workflow_dispatch",
		OIDCIssuer:          "https://token.actions.githubusercontent.com",
		CertificateIdentity: "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-recovery-v1.0.0.yml@refs/tags/release-tooling-v1.0.0-c1",
	}
}
