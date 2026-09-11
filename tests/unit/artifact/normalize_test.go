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
	"strings"
	"testing"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/artifact"
)

func limits() artifact.Limits { return artifact.TestLimits(5, 1000, 32<<20, 8<<20, 1000, time.Minute) }

func writeFile(t *testing.T, path string, value []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, value, 0o600); err != nil {
		t.Fatal(err)
	}
}

func zipBytes(t *testing.T, values map[string][]byte) []byte {
	t.Helper()
	var buffer bytes.Buffer
	w := zip.NewWriter(&buffer)
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	for _, name := range names {
		entry, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := entry.Write(values[name]); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}

func tarBytes(t *testing.T, entries []tar.Header, values [][]byte) []byte {
	t.Helper()
	var buffer bytes.Buffer
	w := tar.NewWriter(&buffer)
	for index := range entries {
		header := entries[index]
		if err := w.WriteHeader(&header); err != nil {
			t.Fatal(err)
		}
		if index < len(values) && values[index] != nil {
			if _, err := w.Write(values[index]); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}

func gzipBytes(t *testing.T, raw []byte) []byte {
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

func normalize(t *testing.T, input string, selected artifact.Limits) artifact.Result {
	t.Helper()
	output := filepath.Join(t.TempDir(), "normalized")
	result, err := artifact.Normalize(context.Background(), input, output, selected)
	if err != nil {
		t.Fatal(err)
	}
	if err := artifact.Verify(result); err != nil {
		t.Fatal(err)
	}
	return result
}

func TestSupportedDirectoryZIPTARTGZAndByteCoverage(t *testing.T) {
	root := filepath.Join(t.TempDir(), "input")
	writeFile(t, filepath.Join(root, "plain.txt"), []byte("synthetic clean text"))
	innerTar := tarBytes(t, []tar.Header{{Name: "nested.txt", Mode: 0o600, Size: 18, Typeflag: tar.TypeReg}}, [][]byte{[]byte("synthetic finding!")})
	writeFile(t, filepath.Join(root, "nested.zip"), zipBytes(t, map[string][]byte{"inner.tar": innerTar}))
	writeFile(t, filepath.Join(root, "bundle.tgz"), gzipBytes(t, innerTar))
	result := normalize(t, root, limits())
	if result.Class != artifact.ClassDirectory || len(result.Entries) < 7 || result.ExpandedBytes <= 0 {
		t.Fatalf("incomplete result: class=%s entries=%d bytes=%d", result.Class, len(result.Entries), result.ExpandedBytes)
	}
	for _, entry := range result.Entries {
		if len(entry.Chunks) == 0 || entry.Chunks[0].Start != 0 || entry.Chunks[len(entry.Chunks)-1].End != entry.Size {
			t.Fatal("entry lacks complete projection")
		}
	}
	probe := filepath.Join(result.ProbeRoot, filepath.FromSlash(result.ExpectedProbeFiles[0]))
	if err := os.WriteFile(probe, []byte("tampered"), 0o600); err != nil {
		t.Fatal(err)
	}
	if artifact.Verify(result) == nil {
		t.Fatal("tampered projection verified")
	}
}

func TestDockerSaveNormalization(t *testing.T) {
	layer := tarBytes(t, []tar.Header{{Name: "app/config.txt", Mode: 0o600, Size: 16, Typeflag: tar.TypeReg}}, [][]byte{[]byte("synthetic layer!")})
	layerSum := sha256.Sum256(layer)
	config, _ := json.Marshal(map[string]any{"architecture": "amd64", "os": "linux", "rootfs": map[string]any{"type": "layers", "diff_ids": []string{"sha256:" + hex.EncodeToString(layerSum[:])}}})
	configSum := sha256.Sum256(config)
	configName := hex.EncodeToString(configSum[:]) + ".json"
	manifest, _ := json.Marshal([]map[string]any{{"Config": configName, "RepoTags": []string{"synthetic:latest"}, "Layers": []string{"layer/layer.tar"}}})
	archive := tarBytes(t, []tar.Header{
		{Name: "manifest.json", Mode: 0o600, Size: int64(len(manifest)), Typeflag: tar.TypeReg},
		{Name: configName, Mode: 0o600, Size: int64(len(config)), Typeflag: tar.TypeReg},
		{Name: "layer/layer.tar", Mode: 0o600, Size: int64(len(layer)), Typeflag: tar.TypeReg},
	}, [][]byte{manifest, config, layer})
	path := filepath.Join(t.TempDir(), "docker-save.tar")
	writeFile(t, path, archive)
	result := normalize(t, path, limits())
	if result.Class != artifact.ClassDockerSave || len(result.Entries) < 5 {
		t.Fatalf("docker save not completely expanded: %s %d", result.Class, len(result.Entries))
	}
}

func TestDocker29ContentAddressedSaveNormalization(t *testing.T) {
	layerTar := tarBytes(t, []tar.Header{{Name: "app/value.txt", Mode: 0o600, Size: 9, Typeflag: tar.TypeReg}}, [][]byte{[]byte("synthetic")})
	layer := gzipBytes(t, layerTar)
	diffSum := sha256.Sum256(layerTar)
	config, _ := json.Marshal(map[string]any{"architecture": "amd64", "os": "linux", "rootfs": map[string]any{"type": "layers", "diff_ids": []string{"sha256:" + hex.EncodeToString(diffSum[:])}}})
	configSum := sha256.Sum256(config)
	layerSum := sha256.Sum256(layer)
	configName := "blobs/sha256/" + hex.EncodeToString(configSum[:])
	layerName := "blobs/sha256/" + hex.EncodeToString(layerSum[:])
	imageManifest, _ := json.Marshal(map[string]any{"schemaVersion": 2,
		"config": descriptor("application/vnd.oci.image.config.v1+json", config),
		"layers": []any{descriptor("application/vnd.oci.image.layer.v1.tar+gzip", layer)}})
	imageManifestSum := sha256.Sum256(imageManifest)
	imageManifestName := "blobs/sha256/" + hex.EncodeToString(imageManifestSum[:])
	index, _ := json.Marshal(map[string]any{"schemaVersion": 2, "manifests": []any{descriptor("application/vnd.oci.image.manifest.v1+json", imageManifest)}})
	manifest, _ := json.Marshal([]map[string]any{{"Config": configName, "RepoTags": nil, "Layers": []string{layerName}}})
	values := map[string][]byte{
		"manifest.json": manifest, "oci-layout": []byte(`{"imageLayoutVersion":"1.0.0"}`), "index.json": index,
		configName: config, layerName: layer, imageManifestName: imageManifest,
	}
	names := []string{"manifest.json", "oci-layout", "index.json", configName, layerName, imageManifestName}
	headers, bodies := make([]tar.Header, 0, len(names)), make([][]byte, 0, len(names))
	for _, name := range names {
		headers = append(headers, tar.Header{Name: name, Mode: 0o600, Size: int64(len(values[name])), Typeflag: tar.TypeReg})
		bodies = append(bodies, values[name])
	}
	path := filepath.Join(t.TempDir(), "docker-29-save.tar")
	writeFile(t, path, tarBytes(t, headers, bodies))
	result := normalize(t, path, limits())
	if result.Class != artifact.ClassDockerSave {
		t.Fatalf("modern Docker save class=%s", result.Class)
	}
	invalidImageManifest := []byte(`{"schemaVersion":2}`)
	invalidSum := sha256.Sum256(invalidImageManifest)
	invalidName := "blobs/sha256/" + hex.EncodeToString(invalidSum[:])
	invalidIndex, _ := json.Marshal(map[string]any{"schemaVersion": 2, "manifests": []any{descriptor("application/vnd.oci.image.manifest.v1+json", invalidImageManifest)}})
	invalidValues := map[string][]byte{
		"manifest.json": manifest, "oci-layout": []byte(`{"imageLayoutVersion":"1.0.0"}`), "index.json": invalidIndex,
		configName: config, layerName: layer, invalidName: invalidImageManifest,
	}
	headers, bodies = headers[:0], bodies[:0]
	for _, name := range []string{"manifest.json", "oci-layout", "index.json", configName, layerName, invalidName} {
		headers = append(headers, tar.Header{Name: name, Mode: 0o600, Size: int64(len(invalidValues[name])), Typeflag: tar.TypeReg})
		bodies = append(bodies, invalidValues[name])
	}
	invalidPath := filepath.Join(t.TempDir(), "invalid-docker-29-save.tar")
	writeFile(t, invalidPath, tarBytes(t, headers, bodies))
	_, err := artifact.Normalize(context.Background(), invalidPath, filepath.Join(t.TempDir(), "out"), limits())
	if !artifact.IsCode(err, artifact.RejectMalformed) {
		t.Fatalf("malformed Docker 29 image manifest got %v", err)
	}
	conflictIndex, _ := json.Marshal(map[string]any{"schemaVersion": 2, "manifests": []any{
		descriptor("application/vnd.oci.image.manifest.v1+json", imageManifest),
		descriptor("application/vnd.oci.image.index.v1+json", imageManifest),
	}})
	conflictValues := map[string][]byte{
		"manifest.json": manifest, "oci-layout": []byte(`{"imageLayoutVersion":"1.0.0"}`), "index.json": conflictIndex,
		configName: config, layerName: layer, imageManifestName: imageManifest,
	}
	headers, bodies = headers[:0], bodies[:0]
	for _, name := range []string{"manifest.json", "oci-layout", "index.json", configName, layerName, imageManifestName} {
		headers = append(headers, tar.Header{Name: name, Mode: 0o600, Size: int64(len(conflictValues[name])), Typeflag: tar.TypeReg})
		bodies = append(bodies, conflictValues[name])
	}
	conflictPath := filepath.Join(t.TempDir(), "conflicting-docker-29-save.tar")
	writeFile(t, conflictPath, tarBytes(t, headers, bodies))
	_, err = artifact.Normalize(context.Background(), conflictPath, filepath.Join(t.TempDir(), "out"), limits())
	if !artifact.IsCode(err, artifact.RejectMalformed) {
		t.Fatalf("conflicting Docker descriptor got %v", err)
	}
}

func TestOCILayoutNormalizationAndDigestBinding(t *testing.T) {
	root := filepath.Join(t.TempDir(), "oci")
	layer := tarBytes(t, []tar.Header{{Name: "app.txt", Mode: 0o600, Size: 15, Typeflag: tar.TypeReg}}, [][]byte{[]byte("synthetic layer")})
	layerSum := sha256.Sum256(layer)
	config, _ := json.Marshal(map[string]any{"architecture": "amd64", "os": "linux", "rootfs": map[string]any{"type": "layers", "diff_ids": []string{"sha256:" + hex.EncodeToString(layerSum[:])}}})
	configDescriptor := descriptor("application/vnd.oci.image.config.v1+json", config)
	layerDescriptor := descriptor("application/vnd.oci.image.layer.v1.tar", layer)
	manifest, _ := json.Marshal(map[string]any{"schemaVersion": 2, "config": configDescriptor, "layers": []any{layerDescriptor}})
	manifestDescriptor := descriptor("application/vnd.oci.image.manifest.v1+json", manifest)
	index, _ := json.Marshal(map[string]any{"schemaVersion": 2, "manifests": []any{manifestDescriptor}})
	writeFile(t, filepath.Join(root, "oci-layout"), []byte(`{"imageLayoutVersion":"1.0.0"}`))
	writeFile(t, filepath.Join(root, "index.json"), index)
	for _, raw := range [][]byte{config, layer, manifest} {
		sum := sha256.Sum256(raw)
		writeFile(t, filepath.Join(root, "blobs", "sha256", hex.EncodeToString(sum[:])), raw)
	}
	result := normalize(t, root, limits())
	if result.Class != artifact.ClassOCILayout || len(result.Entries) < 5 {
		t.Fatalf("OCI layout incomplete: %s %d", result.Class, len(result.Entries))
	}
	writeFile(t, filepath.Join(root, "extra"), []byte("unbound"))
	_, err := artifact.Normalize(context.Background(), root, filepath.Join(t.TempDir(), "bad"), limits())
	if !artifact.IsCode(err, artifact.RejectMalformed) {
		t.Fatal("unbound OCI material was admitted")
	}
}

func TestOCIAndDockerSemanticExpansionBudgetsFailBeforePass(t *testing.T) {
	content := bytes.Repeat([]byte("x"), 4096)
	layerTar := tarBytes(t, []tar.Header{{Name: "large", Mode: 0o600, Size: int64(len(content)), Typeflag: tar.TypeReg}}, [][]byte{content})
	layerGZIP := gzipBytes(t, layerTar)
	diffSum := sha256.Sum256(layerTar)
	config, _ := json.Marshal(map[string]any{"architecture": "amd64", "os": "linux", "rootfs": map[string]any{"type": "layers", "diff_ids": []string{"sha256:" + hex.EncodeToString(diffSum[:])}}})

	t.Run("oci", func(t *testing.T) {
		root := filepath.Join(t.TempDir(), "oci")
		configDescriptor := descriptor("application/vnd.oci.image.config.v1+json", config)
		layerDescriptor := descriptor("application/vnd.oci.image.layer.v1.tar+gzip", layerGZIP)
		manifest, _ := json.Marshal(map[string]any{"schemaVersion": 2, "config": configDescriptor, "layers": []any{layerDescriptor}})
		manifestDescriptor := descriptor("application/vnd.oci.image.manifest.v1+json", manifest)
		index, _ := json.Marshal(map[string]any{"schemaVersion": 2, "manifests": []any{manifestDescriptor}})
		layout := []byte(`{"imageLayoutVersion":"1.0.0"}`)
		writeFile(t, filepath.Join(root, "oci-layout"), layout)
		writeFile(t, filepath.Join(root, "index.json"), index)
		for _, raw := range [][]byte{config, layerGZIP, manifest} {
			sum := sha256.Sum256(raw)
			writeFile(t, filepath.Join(root, "blobs", "sha256", hex.EncodeToString(sum[:])), raw)
		}
		rawBytes := int64(len(layout) + len(index) + len(config) + len(layerGZIP) + len(manifest))
		_, err := artifact.Normalize(context.Background(), root, filepath.Join(t.TempDir(), "out"), artifact.TestLimits(5, 100, rawBytes+512, 1<<20, 1000, time.Minute))
		if !artifact.IsCode(err, artifact.RejectResource) {
			t.Fatalf("OCI aggregate expansion got %v", err)
		}
	})

	t.Run("docker-save", func(t *testing.T) {
		configSum := sha256.Sum256(config)
		configName := hex.EncodeToString(configSum[:]) + ".json"
		manifest, _ := json.Marshal([]map[string]any{{"Config": configName, "RepoTags": []string{"synthetic:latest"}, "Layers": []string{"layer/layer.gz"}}})
		outer := tarBytes(t, []tar.Header{
			{Name: "manifest.json", Mode: 0o600, Size: int64(len(manifest)), Typeflag: tar.TypeReg},
			{Name: configName, Mode: 0o600, Size: int64(len(config)), Typeflag: tar.TypeReg},
			{Name: "layer/layer.gz", Mode: 0o600, Size: int64(len(layerGZIP)), Typeflag: tar.TypeReg},
		}, [][]byte{manifest, config, layerGZIP})
		path := filepath.Join(t.TempDir(), "docker-save.tar")
		writeFile(t, path, outer)
		rawBytes := int64(len(outer) + len(manifest) + len(config) + len(layerGZIP))
		_, err := artifact.Normalize(context.Background(), path, filepath.Join(t.TempDir(), "out"), artifact.TestLimits(5, 100, rawBytes+512, 1<<20, 1000, time.Minute))
		if !artifact.IsCode(err, artifact.RejectResource) {
			t.Fatalf("Docker aggregate expansion got %v", err)
		}
	})
}

func descriptor(media string, raw []byte) map[string]any {
	sum := sha256.Sum256(raw)
	return map[string]any{"mediaType": media, "digest": "sha256:" + hex.EncodeToString(sum[:]), "size": len(raw)}
}

func TestUnsafeArchiveFormsFailClosed(t *testing.T) {
	tests := []struct {
		name string
		raw  []byte
		code artifact.RejectCode
	}{
		{"zip-traversal", zipBytes(t, map[string][]byte{"../escape": []byte("x")}), artifact.RejectUnsafe},
		{"zip-absolute", zipBytes(t, map[string][]byte{"/escape": []byte("x")}), artifact.RejectUnsafe},
		{"zip-case-collision", zipBytes(t, map[string][]byte{"A.txt": []byte("a"), "a.txt": []byte("b")}), artifact.RejectUnsafe},
		{"tar-symlink", tarBytes(t, []tar.Header{{Name: "link", Linkname: "target", Typeflag: tar.TypeSymlink}}, nil), artifact.RejectUnsafe},
		{"tar-hardlink", tarBytes(t, []tar.Header{{Name: "link", Linkname: "target", Typeflag: tar.TypeLink}}, nil), artifact.RejectUnsafe},
		{"tar-device", tarBytes(t, []tar.Header{{Name: "dev", Typeflag: tar.TypeChar}}, nil), artifact.RejectUnsafe},
		{"tar-absolute", tarBytes(t, []tar.Header{{Name: "/absolute", Mode: 0o600, Size: 1, Typeflag: tar.TypeReg}}, [][]byte{[]byte("x")}), artifact.RejectUnsafe},
		{"tar-traversal", tarBytes(t, []tar.Header{{Name: "../escape", Mode: 0o600, Size: 1, Typeflag: tar.TypeReg}}, [][]byte{[]byte("x")}), artifact.RejectUnsafe},
		{"unsupported-7z", []byte{'7', 'z', 0xbc, 0xaf, 0x27, 0x1c, 0}, artifact.RejectUnsupported},
		{"unsupported-rar", []byte{'R', 'a', 'r', '!', 0x1a, 0x07, 0x00}, artifact.RejectUnsupported},
		{"unsupported-bzip2", []byte{'B', 'Z', 'h', '9', 0}, artifact.RejectUnsupported},
		{"unsupported-xz", []byte{0xfd, '7', 'z', 'X', 'Z', 0}, artifact.RejectUnsupported},
		{"unsupported-zstd", []byte{0x28, 0xb5, 0x2f, 0xfd, 0}, artifact.RejectUnsupported},
		{"unsupported-elf", []byte{0x7f, 'E', 'L', 'F', 2, 1, 1}, artifact.RejectUnsupported},
		{"unsupported-pdf", []byte{'%', 'P', 'D', 'F', '-', '1', '.', '7'}, artifact.RejectUnsupported},
		{"malformed-zip", []byte{'P', 'K', 3, 4, 0}, artifact.RejectMalformed},
		{"malformed-gzip", []byte{0x1f, 0x8b, 0x08, 0}, artifact.RejectMalformed},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "candidate")
			writeFile(t, path, test.raw)
			output := filepath.Join(t.TempDir(), "normalized")
			_, err := artifact.Normalize(context.Background(), path, output, limits())
			if !artifact.IsCode(err, test.code) {
				t.Fatalf("got %v", err)
			}
			if _, statErr := os.Stat(output); !os.IsNotExist(statErr) {
				t.Fatal("failed normalization left residue")
			}
		})
	}
}

func TestEncryptedZIPAndConcatenatedGZIPFailClosed(t *testing.T) {
	raw := zipBytes(t, map[string][]byte{"value": []byte("synthetic")})
	for offset := 0; offset+8 < len(raw); offset++ {
		if bytes.Equal(raw[offset:offset+4], []byte{'P', 'K', 3, 4}) {
			raw[offset+6] |= 1
		} else if bytes.Equal(raw[offset:offset+4], []byte{'P', 'K', 1, 2}) {
			raw[offset+8] |= 1
		}
	}
	path := filepath.Join(t.TempDir(), "encrypted.zip")
	writeFile(t, path, raw)
	_, err := artifact.Normalize(context.Background(), path, filepath.Join(t.TempDir(), "out"), limits())
	if !artifact.IsCode(err, artifact.RejectUnsafe) {
		t.Fatalf("encrypted ZIP got %v", err)
	}
	concatenated := append(gzipBytes(t, []byte("one")), gzipBytes(t, []byte("two"))...)
	path = filepath.Join(t.TempDir(), "multi.gz")
	writeFile(t, path, concatenated)
	_, err = artifact.Normalize(context.Background(), path, filepath.Join(t.TempDir(), "out"), limits())
	if !artifact.IsCode(err, artifact.RejectMalformed) {
		t.Fatalf("concatenated gzip got %v", err)
	}
}

func TestResourceBoundariesAndCancellationNeverPass(t *testing.T) {
	root := filepath.Join(t.TempDir(), "input")
	writeFile(t, filepath.Join(root, "one"), []byte("1234"))
	writeFile(t, filepath.Join(root, "two"), []byte("5678"))
	atLimit := artifact.TestLimits(5, 2, 8, 4, 1000, time.Minute)
	if _, err := artifact.Normalize(context.Background(), root, filepath.Join(t.TempDir(), "at"), atLimit); err != nil {
		t.Fatalf("exact boundary rejected: %v", err)
	}
	for name, selected := range map[string]artifact.Limits{
		"entries":  artifact.TestLimits(5, 1, 8, 4, 1000, time.Minute),
		"expanded": artifact.TestLimits(5, 2, 7, 4, 1000, time.Minute),
		"file":     artifact.TestLimits(5, 2, 8, 3, 1000, time.Minute),
	} {
		t.Run(name, func(t *testing.T) {
			_, err := artifact.Normalize(context.Background(), root, filepath.Join(t.TempDir(), "over"), selected)
			if !artifact.IsCode(err, artifact.RejectResource) {
				t.Fatalf("over-limit got %v", err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := artifact.Normalize(ctx, filepath.Join(root, "one"), filepath.Join(t.TempDir(), "cancel"), limits())
	if !artifact.IsCode(err, artifact.RejectCancelled) {
		for _, code := range []artifact.RejectCode{artifact.RejectUnsupported, artifact.RejectUnsafe, artifact.RejectMalformed, artifact.RejectResource, artifact.RejectTimeout, artifact.RejectCancelled, artifact.RejectInvariant} {
			if artifact.IsCode(err, code) {
				t.Fatalf("cancel got code %s", code)
			}
		}
		t.Fatalf("cancel got untyped %v", err)
	}
}

func TestArchiveEntryAndExpansionLimitsRejectBeforeMaterialization(t *testing.T) {
	archives := map[string][]byte{
		"zip": zipBytes(t, map[string][]byte{"one": []byte("1"), "two": []byte("2")}),
		"tar": tarBytes(t, []tar.Header{
			{Name: "one", Mode: 0o600, Size: 1, Typeflag: tar.TypeReg},
			{Name: "two", Mode: 0o600, Size: 1, Typeflag: tar.TypeReg},
		}, [][]byte{[]byte("1"), []byte("2")}),
	}
	for name, raw := range archives {
		t.Run(name+"-entries", func(t *testing.T) {
			path := filepath.Join(t.TempDir(), name)
			writeFile(t, path, raw)
			output := filepath.Join(t.TempDir(), "normalized")
			_, err := artifact.Normalize(context.Background(), path, output, artifact.TestLimits(5, 2, 1<<20, 1<<20, 1000, time.Minute))
			if !artifact.IsCode(err, artifact.RejectResource) {
				t.Fatalf("entry overflow got %v", err)
			}
			if _, err := os.Stat(output); !os.IsNotExist(err) {
				t.Fatal("preflight rejection left material")
			}
		})
		t.Run(name+"-expanded", func(t *testing.T) {
			path := filepath.Join(t.TempDir(), name)
			writeFile(t, path, raw)
			_, err := artifact.Normalize(context.Background(), path, filepath.Join(t.TempDir(), "normalized"), artifact.TestLimits(5, 100, int64(len(raw))+1, 1<<20, 1000, time.Minute))
			if !artifact.IsCode(err, artifact.RejectResource) {
				t.Fatalf("expanded overflow got %v", err)
			}
		})
	}
}

func TestCompressionRatioAndNestingDepth(t *testing.T) {
	raw := zipBytes(t, map[string][]byte{"bomb": bytes.Repeat([]byte("0"), 64<<10)})
	path := filepath.Join(t.TempDir(), "bomb.zip")
	writeFile(t, path, raw)
	selected := artifact.TestLimits(5, 100, 1<<20, 1<<20, 2, time.Minute)
	_, err := artifact.Normalize(context.Background(), path, filepath.Join(t.TempDir(), "out"), selected)
	if !artifact.IsCode(err, artifact.RejectResource) {
		t.Fatalf("ratio bomb got %v", err)
	}
	nested := []byte("leaf")
	for range 6 {
		nested = zipBytes(t, map[string][]byte{"nested.zip": nested})
	}
	path = filepath.Join(t.TempDir(), "deep.zip")
	writeFile(t, path, nested)
	_, err = artifact.Normalize(context.Background(), path, filepath.Join(t.TempDir(), "out"), limits())
	if !artifact.IsCode(err, artifact.RejectResource) {
		t.Fatalf("depth overflow got %v", err)
	}
}

func TestPublicErrorsContainNoCandidateMaterial(t *testing.T) {
	canary := "PSCAN_REDACTION_CANARY_9f0c"
	path := filepath.Join(t.TempDir(), canary+".zip")
	writeFile(t, path, zipBytes(t, map[string][]byte{"../" + canary: []byte(canary)}))
	_, err := artifact.Normalize(context.Background(), path, filepath.Join(t.TempDir(), "out"), limits())
	if err == nil || strings.Contains(err.Error(), canary) || err.Error() != "artifact rejected" {
		t.Fatalf("unsafe public error: %q", err)
	}
}

func TestDirectoryHardLinkRejectedWhenSupported(t *testing.T) {
	root := filepath.Join(t.TempDir(), "input")
	original := filepath.Join(root, "original")
	writeFile(t, original, []byte("synthetic"))
	if err := os.Link(original, filepath.Join(root, "alias")); err != nil {
		t.Skipf("hard links unavailable: %v", err)
	}
	_, err := artifact.Normalize(context.Background(), root, filepath.Join(t.TempDir(), "out"), limits())
	if !artifact.IsCode(err, artifact.RejectUnsafe) {
		t.Fatalf("hard-linked directory input got %v", err)
	}
}

func TestTrailingArchivePayloadsAreMalformed(t *testing.T) {
	zipRaw := append(zipBytes(t, map[string][]byte{"value": []byte("synthetic")}), []byte("trailing")...)
	tarRaw := append(tarBytes(t, []tar.Header{{Name: "value", Mode: 0o600, Size: 1, Typeflag: tar.TypeReg}}, [][]byte{[]byte("x")}), []byte("trailing")...)
	for name, raw := range map[string][]byte{"zip": zipRaw, "tar": tarRaw} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), name)
			writeFile(t, path, raw)
			_, err := artifact.Normalize(context.Background(), path, filepath.Join(t.TempDir(), "out"), limits())
			if !artifact.IsCode(err, artifact.RejectMalformed) {
				t.Fatalf("trailing payload got %v", err)
			}
		})
	}
}

func TestCandidateExecutableIsDataOnly(t *testing.T) {
	root := filepath.Join(t.TempDir(), "input")
	marker := filepath.Join(t.TempDir(), "must-not-exist")
	script := []byte("#!/bin/sh\nprintf executed > " + marker + "\n")
	path := filepath.Join(root, "candidate.sh")
	writeFile(t, path, script)
	if err := os.Chmod(path, 0o700); err != nil {
		t.Fatal(err)
	}
	_ = normalize(t, root, limits())
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("candidate artifact executed")
	}
}

func TestDirectoryLinksAndCaseCollisionsFailClosed(t *testing.T) {
	t.Run("case collision", func(t *testing.T) {
		root := filepath.Join(t.TempDir(), "input")
		upper := filepath.Join(root, "A")
		lower := filepath.Join(root, "a")
		writeFile(t, filepath.Join(upper, "one"), []byte("one"))
		writeFile(t, filepath.Join(lower, "two"), []byte("two"))
		upperInfo, err := os.Stat(upper)
		if err != nil {
			t.Fatal(err)
		}
		lowerInfo, err := os.Stat(lower)
		if err != nil {
			t.Fatal(err)
		}
		if os.SameFile(upperInfo, lowerInfo) {
			t.Skip("filesystem cannot represent distinct case-colliding directories")
		}
		_, err = artifact.Normalize(context.Background(), root, filepath.Join(t.TempDir(), "out"), limits())
		if !artifact.IsCode(err, artifact.RejectUnsafe) {
			t.Fatalf("case-colliding directories got %v", err)
		}
	})

	t.Run("directory symlink", func(t *testing.T) {
		target := filepath.Join(t.TempDir(), "target")
		writeFile(t, target, []byte("target"))
		root := filepath.Join(t.TempDir(), "links")
		if err := os.Mkdir(root, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, filepath.Join(root, "link")); err != nil {
			t.Skipf("symlinks unavailable: %v", err)
		}
		_, err := artifact.Normalize(context.Background(), root, filepath.Join(t.TempDir(), "out"), limits())
		if !artifact.IsCode(err, artifact.RejectUnsafe) {
			t.Fatalf("directory symlink got %v", err)
		}
	})
}

func TestProductProfilesAreExactAndCannotBeExpanded(t *testing.T) {
	pr, release := artifact.PRLimits(), artifact.ReleaseLimits()
	if !pr.Valid() || !release.Valid() || pr.MaxDepth != 5 || pr.MaxEntries != 100_000 || pr.MaxExpandedBytes != 2<<30 ||
		release.MaxExpandedBytes != 10<<30 || pr.MaxFileBytes != 512<<20 || pr.MaxCompressionRatio != 1000 ||
		pr.Timeout != 15*time.Minute || release.Timeout != time.Hour || pr.MaxProcesses != 1 || release.MaxProcesses != 1 {
		t.Fatal("product profile drift")
	}
	pr.MaxEntries++
	if pr.Valid() {
		t.Fatal("expanded profile was accepted")
	}
}

func TestEveryDeclaredResourceBoundaryAtBelowAbove(t *testing.T) {
	for name, profile := range map[string]artifact.Limits{"pr": artifact.PRLimits(), "release": artifact.ReleaseLimits()} {
		t.Run(name, func(t *testing.T) {
			at := artifact.Usage{Depth: profile.MaxDepth, Entries: profile.MaxEntries, ExpandedBytes: profile.MaxExpandedBytes,
				FileBytes: profile.MaxFileBytes, Expanded: profile.MaxCompressionRatio, Compressed: 1,
				Elapsed: profile.Timeout, MemoryBytes: profile.MaxMemoryBytes, Processes: profile.MaxProcesses}
			if err := profile.CheckUsage(at); err != nil {
				t.Fatalf("exact boundary rejected: %v", err)
			}
			below := at
			below.Depth--
			below.Entries--
			below.ExpandedBytes--
			below.FileBytes--
			below.Expanded--
			below.Elapsed--
			below.MemoryBytes--
			if err := profile.CheckUsage(below); err != nil {
				t.Fatalf("below boundary rejected: %v", err)
			}
			cases := map[string]artifact.Usage{
				"depth": at, "entries": at, "expanded": at, "file": at,
				"ratio": at, "timeout": at, "memory": at, "processes": at,
			}
			value := cases["depth"]
			value.Depth++
			cases["depth"] = value
			value = cases["entries"]
			value.Entries++
			cases["entries"] = value
			value = cases["expanded"]
			value.ExpandedBytes++
			cases["expanded"] = value
			value = cases["file"]
			value.FileBytes++
			cases["file"] = value
			value = cases["ratio"]
			value.Expanded++
			cases["ratio"] = value
			value = cases["timeout"]
			value.Elapsed++
			cases["timeout"] = value
			value = cases["memory"]
			value.MemoryBytes++
			cases["memory"] = value
			value = cases["processes"]
			value.Processes++
			cases["processes"] = value
			for boundary, usage := range cases {
				if err := profile.CheckUsage(usage); !artifact.IsCode(err, artifact.RejectResource) {
					t.Errorf("%s over boundary got %v", boundary, err)
				}
			}
		})
	}
}

func TestUnboundOrReorderedProjectionCannotVerify(t *testing.T) {
	input := filepath.Join(t.TempDir(), "input")
	writeFile(t, input, bytes.Repeat([]byte("x"), 100_000))
	result := normalize(t, input, limits())
	writeFile(t, filepath.Join(result.RawRoot, "extra"), []byte("extra"))
	if artifact.Verify(result) == nil {
		t.Fatal("unbound raw material verified")
	}
	if err := os.Remove(filepath.Join(result.RawRoot, "extra")); err != nil {
		t.Fatal(err)
	}
	if len(result.ExpectedProbeFiles) < 2 {
		t.Fatal("fixture did not create multiple chunks")
	}
	result.ExpectedProbeFiles[0], result.ExpectedProbeFiles[1] = result.ExpectedProbeFiles[1], result.ExpectedProbeFiles[0]
	if artifact.Verify(result) == nil {
		t.Fatal("reordered coverage list verified")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := artifact.VerifyContext(ctx, result); !artifact.IsCode(err, artifact.RejectCancelled) {
		t.Fatalf("cancelled verification got %v", err)
	}
}
