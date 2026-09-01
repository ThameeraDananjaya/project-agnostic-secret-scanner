package artifact_test

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/artifact"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/engine/gitleaks"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/outcome"
	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/workspace"
)

const syntheticFinding = "PSCAN_SYNTHETIC_SECRET_83e08b4f9be9467ba5c614c51d7f2a9083e08b4f9be9467ba5c614c51d7f2a90"

func TestPinnedGitleaksArtifactProjection(t *testing.T) {
	binary := requiredPath(t, "PSCAN_GITLEAKS_BINARY")
	config := requiredPath(t, "PSCAN_GITLEAKS_CONFIG")
	ignore := requiredPath(t, "PSCAN_GITLEAKS_IGNORE")
	for _, test := range []struct {
		name  string
		build func(*testing.T, []byte) string
	}{
		{"directory", buildDirectory},
		{"zip", buildZIP},
		{"tar", buildTAR},
		{"gzip", buildGZIP},
		{"tgz", buildTGZ},
		{"oci-layout", buildOCI},
		{"docker-save", buildDockerSave},
	} {
		for _, content := range []struct {
			name   string
			value  string
			reason outcome.ReasonCode
		}{
			{"clean", "synthetic clean input", outcome.ReasonPassNoBlockingFindings},
			{"finding", syntheticFinding, outcome.ReasonFailFindingDetected},
		} {
			t.Run(test.name+"-"+content.name, func(t *testing.T) {
				input := test.build(t, []byte(content.value))
				home := filepath.Join(t.TempDir(), "home")
				if err := os.Mkdir(home, 0o700); err != nil {
					t.Fatal(err)
				}
				adapter := gitleaks.Adapter{Binding: gitleaks.Binding{
					Executable: binary, ExecutableDigest: digest(t, binary),
					Config: config, ConfigDigest: digest(t, config),
					IgnoreFile: ignore, IgnoreFileDigest: digest(t, ignore),
					PrivateHome: home, Environment: engine.SafeEnvironment(filepath.Dir(binary), home),
				}}
				workspaceRoot := filepath.Join(t.TempDir(), "workspaces")
				manager, err := workspace.NewManager(workspaceRoot)
				if err != nil {
					t.Fatal(err)
				}
				result := adapter.ScanArtifact(context.Background(), manager, "123e4567-e89b-42d3-a456-426614174000", 1, input, artifact.PRLimits())
				if err := manager.Close(); err != nil {
					t.Fatal(err)
				}
				if result.Reason != content.reason {
					t.Fatalf("got %s want %s", result.Reason, content.reason)
				}
				if bytes.Contains([]byte(result.Reason), []byte(content.value)) {
					t.Fatal("candidate value crossed private engine boundary")
				}
				entries, err := os.ReadDir(workspaceRoot)
				if err != nil || len(entries) != 0 {
					t.Fatalf("transient residue remains: %v %v", entries, err)
				}
			})
		}
	}
}

func buildDirectory(t *testing.T, value []byte) string {
	root := filepath.Join(t.TempDir(), "candidate")
	if err := os.Mkdir(root, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "payload.txt"), value, 0o600); err != nil {
		t.Fatal(err)
	}
	return root
}

func buildZIP(t *testing.T, value []byte) string {
	path := filepath.Join(t.TempDir(), "candidate.zip")
	t.Helper()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	entry, err := w.Create("payload.txt")
	if err == nil {
		_, err = entry.Write(value)
	}
	if closeErr := w.Close(); err == nil {
		err = closeErr
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func buildTAR(t *testing.T, value []byte) string {
	path := filepath.Join(t.TempDir(), "candidate.tar")
	writeTARFile(t, path, map[string][]byte{"payload.txt": value})
	return path
}

func buildGZIP(t *testing.T, value []byte) string {
	path := filepath.Join(t.TempDir(), "candidate.gz")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	w := gzip.NewWriter(f)
	_, err = w.Write(value)
	if closeErr := w.Close(); err == nil {
		err = closeErr
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func buildTGZ(t *testing.T, value []byte) string {
	tarPath := buildTAR(t, value)
	raw, err := os.ReadFile(tarPath)
	if err != nil {
		t.Fatal(err)
	}
	return buildGZIP(t, raw)
}

func buildOCI(t *testing.T, value []byte) string {
	root := filepath.Join(t.TempDir(), "oci")
	layerPath := filepath.Join(t.TempDir(), "layer.tar")
	writeTARFile(t, layerPath, map[string][]byte{"payload.txt": value})
	layer, _ := os.ReadFile(layerPath)
	layerSum := sha256.Sum256(layer)
	config, _ := json.Marshal(map[string]any{"architecture": "amd64", "os": "linux", "rootfs": map[string]any{"type": "layers", "diff_ids": []string{"sha256:" + hex.EncodeToString(layerSum[:])}}})
	configDescriptor := descriptorFor(config, "application/vnd.oci.image.config.v1+json")
	layerDescriptor := descriptorFor(layer, "application/vnd.oci.image.layer.v1.tar")
	manifest, _ := json.Marshal(map[string]any{"schemaVersion": 2, "config": configDescriptor, "layers": []any{layerDescriptor}})
	manifestDescriptor := descriptorFor(manifest, "application/vnd.oci.image.manifest.v1+json")
	index, _ := json.Marshal(map[string]any{"schemaVersion": 2, "manifests": []any{manifestDescriptor}})
	writePath(t, filepath.Join(root, "oci-layout"), []byte(`{"imageLayoutVersion":"1.0.0"}`))
	writePath(t, filepath.Join(root, "index.json"), index)
	for _, raw := range [][]byte{config, layer, manifest} {
		sum := sha256.Sum256(raw)
		writePath(t, filepath.Join(root, "blobs", "sha256", hex.EncodeToString(sum[:])), raw)
	}
	return root
}

func buildDockerSave(t *testing.T, value []byte) string {
	layerPath := filepath.Join(t.TempDir(), "layer.tar")
	writeTARFile(t, layerPath, map[string][]byte{"payload.txt": value})
	layerTar, _ := os.ReadFile(layerPath)
	layer := gzipRaw(t, layerTar)
	diffSum := sha256.Sum256(layerTar)
	config, _ := json.Marshal(map[string]any{"architecture": "amd64", "os": "linux", "rootfs": map[string]any{"type": "layers", "diff_ids": []string{"sha256:" + hex.EncodeToString(diffSum[:])}}})
	configSum := sha256.Sum256(config)
	layerSum := sha256.Sum256(layer)
	configName := "blobs/sha256/" + hex.EncodeToString(configSum[:])
	layerName := "blobs/sha256/" + hex.EncodeToString(layerSum[:])
	imageManifest, _ := json.Marshal(map[string]any{"schemaVersion": 2,
		"config": descriptorFor(config, "application/vnd.oci.image.config.v1+json"),
		"layers": []any{descriptorFor(layer, "application/vnd.oci.image.layer.v1.tar+gzip")}})
	imageManifestSum := sha256.Sum256(imageManifest)
	imageManifestName := "blobs/sha256/" + hex.EncodeToString(imageManifestSum[:])
	index, _ := json.Marshal(map[string]any{"schemaVersion": 2,
		"manifests": []any{descriptorFor(imageManifest, "application/vnd.oci.image.manifest.v1+json")}})
	manifest, _ := json.Marshal([]map[string]any{{"Config": configName, "RepoTags": []string{"synthetic:latest"}, "Layers": []string{layerName}}})
	path := filepath.Join(t.TempDir(), "docker-save.tar")
	writeTARFile(t, path, map[string][]byte{"manifest.json": manifest,
		"oci-layout": []byte(`{"imageLayoutVersion":"1.0.0"}`), "index.json": index,
		configName: config, layerName: layer, imageManifestName: imageManifest})
	return path
}

func gzipRaw(t *testing.T, raw []byte) []byte {
	t.Helper()
	var buffer bytes.Buffer
	w := gzip.NewWriter(&buffer)
	if _, err := w.Write(raw); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}

func writeTARFile(t *testing.T, path string, values map[string][]byte) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	w := tar.NewWriter(f)
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		value := values[name]
		header := &tar.Header{Name: name, Mode: 0o600, Size: int64(len(value)), Typeflag: tar.TypeReg}
		if err = w.WriteHeader(header); err == nil {
			_, err = w.Write(value)
		}
		if err != nil {
			break
		}
	}
	if closeErr := w.Close(); err == nil {
		err = closeErr
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		t.Fatal(err)
	}
}

func descriptorFor(raw []byte, mediaType string) map[string]any {
	sum := sha256.Sum256(raw)
	return map[string]any{"mediaType": mediaType, "digest": "sha256:" + hex.EncodeToString(sum[:]), "size": len(raw)}
}

func writePath(t *testing.T, path string, raw []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		t.Fatal(err)
	}
}

func requiredPath(t *testing.T, name string) string {
	t.Helper()
	value := os.Getenv(name)
	if value == "" {
		t.Skip(name + " is required")
	}
	if !filepath.IsAbs(value) {
		t.Fatal(name + " must be absolute")
	}
	return value
}

func digest(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
